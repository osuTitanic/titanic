package storage

import (
	"io"

	_ "github.com/osuTitanic/titanic-go/internal/logging"
)

// ReaderAtCloser provides random-access reads of an object
type ReaderAtCloser interface {
	io.ReaderAt
	io.Closer
}

// Storage defines the interface for a storage backend
type Storage interface {
	Setup() error
	Save(key string, directory string, data []byte) error
	SaveStream(key string, directory string, stream io.Reader) error
	SaveUrl(key string, directory string, url string) error
	Read(key string, directory string) ([]byte, error)
	ReadStream(key string, directory string) (io.ReadSeekCloser, error)
	ReadStreamAt(key string, directory string) (ReaderAtCloser, int64, error)
	Remove(key string, directory string) error
	Exists(key string, directory string) bool

	// TODO: Add context.Context to all methods
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
