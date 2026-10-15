package storage

import (
	"context"
	"io"

	_ "github.com/osuTitanic/titanic/internal/logging"
)

// ReaderAtCloser provides random-access reads of an object
type ReaderAtCloser interface {
	io.ReaderAt
	io.Closer
}

// Storage defines the interface for a storage backend
type Storage interface {
	Setup(ctx context.Context) error
	Save(ctx context.Context, key string, directory string, data []byte) error
	SaveStream(ctx context.Context, key string, directory string, stream io.Reader) error
	SaveUrl(ctx context.Context, key string, directory string, url string) error
	Read(ctx context.Context, key string, directory string) ([]byte, error)
	ReadStream(ctx context.Context, key string, directory string) (io.ReadSeekCloser, error)
	ReadStreamAt(ctx context.Context, key string, directory string) (ReaderAtCloser, int64, error)
	Remove(ctx context.Context, key string, directory string) error
	Exists(ctx context.Context, key string, directory string) bool
}

var RequiredDirectories = []string{
	"audio",
	"avatars",
	"beatmaps",
	"osz",
	"osz2",
	"release",
	"replays",
	"screenshots",
	"thumbnails",
}
