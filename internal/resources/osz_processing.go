package resources

import (
	"archive/zip"
	"io"
	"path"
	"strings"

	"github.com/Lekuruu/zipstream"
	"github.com/osuTitanic/titanic/internal/storage"
)

const noVideoBufferSize = 256 * 1024

var videoFileExtensions = map[string]struct{}{
	".wmv":  {},
	".flv":  {},
	".mp4":  {},
	".avi":  {},
	".m4v":  {},
	".mpg":  {},
	".mov":  {},
	".webm": {},
	".mkv":  {},
	".ogv":  {},
	".mpeg": {},
	".3gp":  {},
}

func IsVideoFile(name string) bool {
	_, ok := videoFileExtensions[strings.ToLower(path.Ext(name))]
	return ok
}

func StreamWithoutVideo(store storage.Storage, key string) (io.ReadCloser, int64, error) {
	source, size, err := store.ReadStreamAt(key, "osz")
	if err != nil {
		return nil, 0, err
	}

	stream, streamSize, err := zipstream.StreamFiltered(
		source,
		size,
		func(file *zip.File) bool {
			// Filter out video files from the archive
			return !IsVideoFile(file.Name)
		},
		zipstream.WithBufferSize(noVideoBufferSize),
	)
	if err != nil {
		return nil, 0, err
	}

	return stream, streamSize, nil
}
