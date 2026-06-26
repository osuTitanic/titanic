package resources

import (
	"bytes"
	"io"
	"log/slog"
	"strconv"

	"github.com/osuTitanic/titanic-go/internal/media"
	"github.com/osuTitanic/titanic-go/internal/storage"
)

// Dimensions for the beatmap thumbnails
const (
	thumbnailSmallWidth  = 80
	thumbnailSmallHeight = 60
	thumbnailLargeWidth  = 160
	thumbnailLargeHeight = 120
)

// StorageResolver receives beatmap resources directly from our own storage backend.
// It is used for beatmaps where the download server is hosted by us (constants.BeatmapServerTitanic).
type StorageResolver struct {
	logger  *slog.Logger
	storage storage.Storage
}

func NewStorageResolver(store storage.Storage) *StorageResolver {
	return &StorageResolver{
		logger:  slog.Default().With("component", "BeatmapStorage"),
		storage: store,
	}
}

func (resolver *StorageResolver) Setup() error {
	// guh
	return nil
}

func (resolver *StorageResolver) Osz(setId int, noVideo bool) (io.ReadCloser, int64, error) {
	resolver.logger.Debug("Reading osz from storage...", "set_id", setId, "no_video", noVideo)

	if !noVideo {
		return resolver.ReadStreamAndSize(strconv.Itoa(setId), "osz")
	}

	stream, size, err := StreamWithoutVideo(resolver.storage, strconv.Itoa(setId))
	if err != nil {
		resolver.logger.Warn(
			"Failed to start no-video osz stream, serving original anyway",
			"set_id", setId, "error", err.Error(),
		)
		return resolver.ReadStreamAndSize(strconv.Itoa(setId), "osz")
	}

	return stream, size, nil
}

func (resolver *StorageResolver) Osu(beatmapId int) (io.ReadCloser, error) {
	resolver.logger.Debug("Reading beatmap from storage...", "beatmap_id", beatmapId)
	return resolver.ReadStream(strconv.Itoa(beatmapId), "beatmaps")
}

func (resolver *StorageResolver) Preview(setId int) (io.ReadCloser, error) {
	resolver.logger.Debug("Reading preview from storage...", "set_id", setId)
	return resolver.ReadStream(strconv.Itoa(setId), "audio")
}

func (resolver *StorageResolver) Background(setId int, large bool) (io.ReadCloser, error) {
	resolver.logger.Debug(
		"Reading background from storage...",
		"set_id", setId,
		"large", large,
	)

	stream, err := resolver.ReadStream(strconv.Itoa(setId), "thumbnails")
	if err != nil {
		return nil, err
	}

	// We only keep a single (large) thumbnail per set in storage
	if large {
		return stream, nil
	}

	defer stream.Close()
	data, err := io.ReadAll(stream)
	if err != nil {
		return nil, err
	}

	resized, err := media.ResizeImage(data, thumbnailSmallWidth, thumbnailSmallHeight)
	if err != nil {
		resolver.logger.Warn(
			"Failed to resize thumbnail, serving original anyway",
			"set_id", setId, "error", err.Error(),
		)
		return io.NopCloser(bytes.NewReader(data)), nil
	}

	// We have successfully resized the thumbnail
	return io.NopCloser(bytes.NewReader(resized)), nil
}

func (resolver *StorageResolver) ReadStream(key string, bucket string) (io.ReadSeekCloser, error) {
	stream, err := resolver.storage.ReadStream(key, bucket)
	if err != nil {
		return nil, ErrResourceNotFound
	}
	return stream, nil
}

func (resolver *StorageResolver) ReadStreamAndSize(key string, bucket string) (io.ReadSeekCloser, int64, error) {
	stream, err := resolver.storage.ReadStream(key, bucket)
	if err != nil {
		return nil, 0, ErrResourceNotFound
	}

	size, err := stream.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, 0, err
	}
	if _, err := stream.Seek(0, io.SeekStart); err != nil {
		return nil, 0, err
	}

	return stream, size, nil
}
