package resources

import (
	"io"
	"log/slog"

	"github.com/osuTitanic/titanic-go/internal/config"
	"github.com/osuTitanic/titanic-go/internal/constants"
	"github.com/osuTitanic/titanic-go/internal/repositories"
	"github.com/osuTitanic/titanic-go/internal/storage"
	"github.com/redis/go-redis/v9"
)

type BeatmapProvider struct {
	logger      *slog.Logger
	beatmapsets *repositories.BeatmapsetRepository
	resolvers   map[constants.BeatmapServer]BeatmapResourceResolver
	fallback    BeatmapResourceResolver
}

func NewBeatmapProvider(
	cfg *config.Config,
	cache *redis.Client,
	store storage.Storage,
	mirrors *repositories.ResourceMirrorRepository,
	beatmapsets *repositories.BeatmapsetRepository,
) *BeatmapProvider {
	// For bancho maps we need to reach out for 3rd-party beatmap mirrors
	// For titanic sourced maps we can directly use our storage solution
	mirrorResolver := NewMirrorResolver(constants.BeatmapServerBancho, cfg, cache, mirrors)
	storageResolver := NewStorageResolver(store)

	return &BeatmapProvider{
		logger:      slog.Default().With("component", "BeatmapProvider"),
		beatmapsets: beatmapsets,
		fallback:    mirrorResolver,
		resolvers: map[constants.BeatmapServer]BeatmapResourceResolver{
			constants.BeatmapServerBancho:  mirrorResolver,
			constants.BeatmapServerTitanic: storageResolver,
		},
	}
}

func (provider *BeatmapProvider) Setup() error {
	for _, resolver := range provider.resolvers {
		if err := resolver.Setup(); err != nil {
			return err
		}
	}
	return nil
}

func (provider *BeatmapProvider) Osz(setId int, noVideo bool) (io.ReadCloser, error) {
	return provider.ResolverForSet(setId).Osz(setId, noVideo)
}

func (provider *BeatmapProvider) Osu(beatmapId int) (io.ReadCloser, error) {
	// TODO: redis caching
	return provider.ResolverForBeatmap(beatmapId).Osu(beatmapId)
}

func (provider *BeatmapProvider) Preview(setId int) (io.ReadCloser, error) {
	// TODO: redis caching
	return provider.ResolverForSet(setId).Preview(setId)
}

func (provider *BeatmapProvider) Background(setId int, large bool) (io.ReadCloser, error) {
	// TODO: redis caching
	return provider.ResolverForSet(setId).Background(setId, large)
}

func (provider *BeatmapProvider) ResolverForServer(server constants.BeatmapServer) BeatmapResourceResolver {
	if resolver, ok := provider.resolvers[server]; ok {
		return resolver
	}
	return provider.fallback
}

// ResolverForSet picks the resolver responsible for a beatmap set.
func (provider *BeatmapProvider) ResolverForSet(setId int) BeatmapResourceResolver {
	server, err := provider.beatmapsets.FetchDownloadServer(setId)
	if err != nil {
		provider.logger.Warn(
			"Failed to determine download server, using fallback",
			"set_id", setId, "error", err.Error(),
		)
		return provider.fallback
	}
	return provider.ResolverForServer(server)
}

// ResolverForBeatmap picks the resolver responsible for the set a beatmap belongs to.
func (provider *BeatmapProvider) ResolverForBeatmap(beatmapId int) BeatmapResourceResolver {
	server, err := provider.beatmapsets.FetchDownloadServerByBeatmap(beatmapId)
	if err != nil {
		provider.logger.Warn(
			"Failed to determine download server, using fallback",
			"beatmap_id", beatmapId, "error", err.Error(),
		)
		return provider.fallback
	}
	return provider.ResolverForServer(server)
}
