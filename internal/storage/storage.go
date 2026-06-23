package storage

import (
	"io"

	_ "github.com/osuTitanic/titanic-go/internal/logging"
)

// Storage defines the interface for a storage backend
type Storage interface {
	Setup() error
	Save(key string, directory string, data []byte) error
	SaveStream(key string, directory string, stream io.Reader) error
	Read(key string, directory string) ([]byte, error)
	ReadStream(key string, directory string) (io.ReadSeekCloser, error)
	Remove(key string, directory string) error
	Exists(key string, directory string) bool
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
	"logs",
}
