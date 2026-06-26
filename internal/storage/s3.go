package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Config struct {
	// Endpoint without the schema, e.g. "s3.amazonaws.com"
	Endpoint string

	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string

	BucketName string
	Region     string
	UseSSL     bool
}

type S3Storage struct {
	client              *minio.Client
	bucketName          string
	region              string
	requiredDirectories []string
}

func NewS3Storage(cfg S3Config) (Storage, error) {
	if strings.TrimSpace(cfg.Endpoint) == "" {
		return nil, errors.New("s3 endpoint must not be empty")
	}
	if strings.TrimSpace(cfg.BucketName) == "" {
		return nil, errors.New("s3 bucket name must not be empty")
	}

	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			cfg.SessionToken,
		),
		Secure: cfg.UseSSL,
		Region: cfg.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("create minio client: %w", err)
	}

	return &S3Storage{
		client:              client,
		bucketName:          cfg.BucketName,
		region:              cfg.Region,
		requiredDirectories: RequiredDirectories,
	}, nil
}

func (s *S3Storage) Setup() error {
	ctx := context.Background()

	// Setup main bucket where all objects will be stored
	exists, err := s.client.BucketExists(ctx, s.bucketName)
	if err != nil {
		return fmt.Errorf("check bucket %q: %w", s.bucketName, err)
	}

	if !exists {
		err = s.client.MakeBucket(ctx, s.bucketName, minio.MakeBucketOptions{
			Region: s.region,
		})
		if err != nil {
			resp := minio.ToErrorResponse(err)

			if resp.Code != "BucketAlreadyOwnedByYou" &&
				resp.Code != "BucketAlreadyExists" {
				return fmt.Errorf("create bucket %q: %w", s.bucketName, err)
			}
		}
	}

	// Setup directories well... not really "directories", since S3 doesn't have them, but
	// rather directory markers, which are just empty objects with a trailing slash in their name
	for _, directory := range s.requiredDirectories {
		prefix := cleanPrefix(directory)
		if prefix == "" {
			continue
		}
		objectName := prefix + "/"

		_, err := s.client.PutObject(
			ctx,
			s.bucketName,
			objectName,
			bytes.NewReader(nil),
			0,
			minio.PutObjectOptions{
				ContentType: "application/x-directory",
			},
		)
		if err != nil {
			return fmt.Errorf("create directory marker %q: %w", prefix, err)
		}
	}

	// TODO: Download default avatars, I suppose...
	return nil
}

func (s *S3Storage) Save(key string, directory string, data []byte) error {
	objectName, err := makeObjectName(directory, key)
	if err != nil {
		return err
	}

	_, err = s.client.PutObject(
		context.Background(),
		s.bucketName,
		objectName,
		bytes.NewReader(data),
		int64(len(data)),
		minio.PutObjectOptions{
			ContentType: "application/octet-stream",
		},
	)
	if err != nil {
		return fmt.Errorf("save object %q: %w", objectName, err)
	}
	return nil
}

func (s *S3Storage) SaveStream(key string, directory string, stream io.Reader) error {
	objectName, err := makeObjectName(directory, key)
	if err != nil {
		return err
	}

	_, err = s.client.PutObject(
		context.Background(),
		s.bucketName,
		objectName,
		stream,
		-1,
		minio.PutObjectOptions{
			ContentType: "application/octet-stream",
			PartSize:    10 * 1024 * 1024,
		},
	)
	if err != nil {
		return fmt.Errorf("save stream object %q: %w", objectName, err)
	}
	return nil
}

func (s *S3Storage) SaveUrl(key string, directory string, url string) error {
	stream, err := downloadStream(url)
	if err != nil {
		return fmt.Errorf("failed to download url content %q: %w", url, err)
	}
	if stream == nil {
		return fmt.Errorf("no download stream found for url %q", url)
	}
	defer stream.Close()

	err = s.SaveStream(key, directory, stream)
	if err != nil {
		return err
	}
	return nil
}

func (s *S3Storage) Read(key string, directory string) ([]byte, error) {
	stream, err := s.ReadStream(key, directory)
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	data, err := io.ReadAll(stream)
	if err != nil {
		return nil, fmt.Errorf("read object stream: %w", err)
	}
	return data, nil
}

func (s *S3Storage) ReadStream(key string, directory string) (io.ReadSeekCloser, error) {
	objectName, err := makeObjectName(directory, key)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()

	object, err := s.client.GetObject(
		ctx,
		s.bucketName,
		objectName,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("get object %q: %w", objectName, err)
	}

	// GetObject may defer errors until we first read from the stream
	// We can use Stat() to validate that the object exists before returning the stream
	if _, err := object.Stat(); err != nil {
		object.Close()
		return nil, fmt.Errorf("stat object %q: %w", objectName, err)
	}

	return object, nil
}

func (s *S3Storage) ReadStreamAt(key string, directory string) (ReaderAtCloser, int64, error) {
	objectName, err := makeObjectName(directory, key)
	if err != nil {
		return nil, 0, err
	}
	ctx := context.Background()

	object, err := s.client.GetObject(
		ctx,
		s.bucketName,
		objectName,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, 0, fmt.Errorf("get object %q: %w", objectName, err)
	}

	info, err := object.Stat()
	if err != nil {
		object.Close()
		return nil, 0, fmt.Errorf("stat object %q: %w", objectName, err)
	}

	objectBuffer := newBufferedReader(object, info.Size, 6*1024*1024)
	return objectBuffer, info.Size, nil
}

func (s *S3Storage) Remove(key string, directory string) error {
	objectName, err := makeObjectName(directory, key)
	if err != nil {
		return err
	}

	err = s.client.RemoveObject(
		context.Background(),
		s.bucketName,
		objectName,
		minio.RemoveObjectOptions{},
	)
	if err != nil {
		return fmt.Errorf("remove object %q: %w", objectName, err)
	}
	return nil
}

func (s *S3Storage) Exists(key string, directory string) bool {
	objectName, err := makeObjectName(directory, key)
	if err != nil {
		return false
	}

	_, err = s.client.StatObject(
		context.Background(),
		s.bucketName,
		objectName,
		minio.StatObjectOptions{},
	)
	return err == nil
}

func makeObjectName(directory string, key string) (string, error) {
	directory = cleanPrefix(directory)
	key = cleanPrefix(key)

	if key == "" {
		return "", errors.New("s3 object key must not be empty")
	}
	if directory == "" {
		return key, nil
	}
	return directory + "/" + key, nil
}

// e.g. "//foo/bar/../baz" -> "foo/baz"
func cleanPrefix(value string) string {
	value = strings.Trim(value, "/")
	if value == "" {
		return ""
	}

	cleaned := path.Clean("/" + value)
	return strings.TrimPrefix(cleaned, "/")
}
