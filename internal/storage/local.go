package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type FileStorage struct {
	dataPath string
}

func NewFileStorage(dataPath string) *FileStorage {
	return &FileStorage{dataPath: dataPath}
}

func (storage *FileStorage) Setup(ctx context.Context) error {
	err := storage.CreateDefaultFolders(ctx)
	if err != nil {
		return err
	}

	err = storage.DownloadDefaultAssets(ctx)
	if err != nil {
		return err
	}
	return nil
}

// TODO: Use os.Root to ensure all file operations stay in dataPath
//       This does not need immediate attention since we only store /
// 		 request IDs right now, which are validated beforehand

// TODO: Idk if there's a standardized way to use context.Context w/
// 		 os-level reads / writes. If there is, we should implement that.

func (storage *FileStorage) Read(ctx context.Context, key string, folder string) ([]byte, error) {
	stream, err := storage.ReadStream(ctx, key, folder)
	if err != nil {
		return nil, err
	}
	defer stream.Close()
	return io.ReadAll(stream)
}

func (storage *FileStorage) ReadStream(ctx context.Context, key string, folder string) (io.ReadSeekCloser, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("%s/%s/%s", storage.dataPath, folder, key)
	return os.Open(path)
}

func (storage *FileStorage) ReadStreamAt(ctx context.Context, key string, folder string) (ReaderAtCloser, int64, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	path := fmt.Sprintf("%s/%s/%s", storage.dataPath, folder, key)
	file, err := os.Open(path)
	if err != nil {
		return nil, 0, err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, 0, err
	}

	return file, info.Size(), nil
}

func (storage *FileStorage) Save(ctx context.Context, key string, folder string, data []byte) error {
	return storage.SaveStream(ctx, key, folder, bytes.NewReader(data))
}

func (storage *FileStorage) SaveStream(ctx context.Context, key string, folder string, stream io.Reader) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	filePath := fmt.Sprintf("%s/%s/%s", storage.dataPath, folder, key)
	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	buffer := make([]byte, 256*1024)
	_, err = io.CopyBuffer(file, stream, buffer)
	return err
}

func (storage *FileStorage) SaveUrl(ctx context.Context, key string, directory string, url string) error {
	stream, err := downloadStream(ctx, url)
	if err != nil {
		return fmt.Errorf("failed to download url content %q: %w", url, err)
	}
	if stream == nil {
		return fmt.Errorf("no download stream found for url %q", url)
	}
	defer stream.Close()

	err = storage.SaveStream(ctx, key, directory, stream)
	if err != nil {
		return err
	}
	return nil
}

func (storage *FileStorage) Exists(ctx context.Context, key string, folder string) bool {
	if ctx.Err() != nil {
		return false
	}
	path := fmt.Sprintf("%s/%s/%s", storage.dataPath, folder, key)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func (storage *FileStorage) Remove(ctx context.Context, key string, folder string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	path := fmt.Sprintf("%s/%s/%s", storage.dataPath, folder, key)
	return os.Remove(path)
}

// CreateDefaultFolders ensures that all required storage folders exist
func (storage *FileStorage) CreateDefaultFolders(ctx context.Context) error {
	for _, directory := range RequiredDirectories {
		if err := ctx.Err(); err != nil {
			return err
		}
		folder := fmt.Sprintf(
			"%s/%s",
			storage.dataPath, directory,
		)

		if _, err := os.Stat(folder); !os.IsNotExist(err) {
			slog.Debug("Storage directory already exists, skipping creation", slog.String("path", folder))
			continue
		}

		err := os.MkdirAll(folder, 0755)
		if err != nil {
			return fmt.Errorf("failed to create storage directory %s: %w", folder, err)
		}

		slog.Info("Created storage directory", slog.String("path", folder))
	}
	return nil
}

// DownloadDefaultAssets downloads all default assets if they do not already exist
func (storage *FileStorage) DownloadDefaultAssets(ctx context.Context) error {
	for assetUrl := range defaultAssetUrls {
		if err := ctx.Err(); err != nil {
			return err
		}
		parts := strings.SplitN(assetUrl, "/", 2)
		if len(parts) != 2 {
			slog.Warn("Invalid asset URL, skipping download", slog.String("url", assetUrl))
			continue
		}
		folder := parts[0]
		key := parts[1]

		// Check if asset already exists
		if storage.Exists(ctx, key, folder) {
			slog.Debug("Asset already exists, skipping download", slog.String("path", assetUrl))
			continue
		}

		stream, err := downloadAssetStream(ctx, assetUrl)
		if err != nil {
			return fmt.Errorf("failed to get download stream for %s: %w", assetUrl, err)
		}
		if stream == nil {
			return fmt.Errorf("no download stream found for %s", assetUrl)
		}
		defer stream.Close()

		err = storage.SaveStream(ctx, key, folder, stream)
		if err != nil {
			return fmt.Errorf("failed to save download %s to storage: %w", assetUrl, err)
		}
		slog.Info("Downloaded asset to storage", slog.String("path", assetUrl))
	}
	return nil
}
