package resources

import (
	"io"
	"log/slog"
	"strconv"

	"github.com/osuTitanic/titanic-go/internal/storage"
)

// StorageResolver receives beatmap resources directly from our own storage backend.
// It is used for beatmaps where the download server is hosted by us (constants.BeatmapServerTitanic).
type StorageResolver struct {
	logger  *slog.Logger
	storage storage.Storage
}

func NewStorageResolver(store storage.Storage) *StorageResolver {
	return &StorageResolver{
		logger:  slog.Default(),
		storage: store,
	}
}

func (resolver *StorageResolver) Setup() error {
	// guh
	return nil
}

func (resolver *StorageResolver) Osz(setId int, noVideo bool) (io.ReadCloser, error) {
	resolver.logger.Debug("Reading osz from storage...", slog.Int("set_id", setId))
	return resolver.ReadStream(strconv.Itoa(setId), "osz")
}

func (resolver *StorageResolver) Osu(beatmapId int) (io.ReadCloser, error) {
	resolver.logger.Debug("Reading beatmap from storage...", slog.Int("beatmap_id", beatmapId))
	return resolver.ReadStream(strconv.Itoa(beatmapId), "beatmaps")
}

func (resolver *StorageResolver) Preview(setId int) (io.ReadCloser, error) {
	resolver.logger.Debug("Reading preview from storage...", slog.Int("set_id", setId))
	return resolver.ReadStream(strconv.Itoa(setId), "audio")
}

func (resolver *StorageResolver) Background(setId int, large bool) (io.ReadCloser, error) {
	resolver.logger.Debug("Reading background from storage...", slog.Int("set_id", setId))
	// TODO: We only keep a single thumbnail per set, regardless of size.
	// 		 We'd have to implement a solution to resize them here.
	//       It could be a shared logic between avatars & thumbnails, since they both require resizing.
	return resolver.ReadStream(strconv.Itoa(setId), "thumbnails")
}

func (resolver *StorageResolver) ReadStream(key string, bucket string) (io.ReadCloser, error) {
	stream, err := resolver.storage.ReadStream(key, bucket)
	if err != nil {
		return nil, ErrResourceNotFound
	}
	return stream, nil
}
