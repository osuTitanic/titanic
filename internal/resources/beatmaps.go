package resources

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/osuTitanic/titanic/internal/config"
	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/repositories"
	"github.com/osuTitanic/titanic/internal/storage"
	"github.com/redis/go-redis/v9"
)

// Cache durations for the different resource types
const (
	cacheTTLBeatmap    = 4 * time.Hour
	cacheTTLPreview    = 6 * time.Hour
	cacheTTLBackground = 12 * time.Hour
)

type BeatmapProvider struct {
	logger      *slog.Logger
	cache       *redis.Client
	beatmapsets *repositories.BeatmapsetRepository
	resolvers   map[constants.BeatmapServer]BeatmapResourceProvider
	fallback    BeatmapResourceProvider
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
		cache:       cache,
		beatmapsets: beatmapsets,
		fallback:    mirrorResolver,
		resolvers: map[constants.BeatmapServer]BeatmapResourceProvider{
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

func (provider *BeatmapProvider) Osz(setId int, noVideo bool) (io.ReadCloser, int64, error) {
	return provider.ResolverForSet(setId).Osz(setId, noVideo)
}

func (provider *BeatmapProvider) Osu(beatmapId int) (io.ReadCloser, error) {
	key := fmt.Sprintf("osu:%d", beatmapId)
	return provider.Cached(key, cacheTTLBeatmap, func() (io.ReadCloser, error) {
		return provider.ResolverForBeatmap(beatmapId).Osu(beatmapId)
	})
}

func (provider *BeatmapProvider) Preview(setId int) (io.ReadCloser, error) {
	key := fmt.Sprintf("mp3:%d", setId)
	return provider.Cached(key, cacheTTLPreview, func() (io.ReadCloser, error) {
		return provider.ResolverForSet(setId).Preview(setId)
	})
}

func (provider *BeatmapProvider) Background(setId int, large bool) (io.ReadCloser, error) {
	key := fmt.Sprintf("mt:%d", setId)
	if large {
		key += "l"
	}
	return provider.Cached(key, cacheTTLBackground, func() (io.ReadCloser, error) {
		return provider.ResolverForSet(setId).Background(setId, large)
	})
}

func (provider *BeatmapProvider) ResolverForServer(server constants.BeatmapServer) BeatmapResourceProvider {
	if resolver, ok := provider.resolvers[server]; ok {
		return resolver
	}
	return provider.fallback
}

// ResolverForSet picks the resolver responsible for a beatmap set.
func (provider *BeatmapProvider) ResolverForSet(setId int) BeatmapResourceProvider {
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
func (provider *BeatmapProvider) ResolverForBeatmap(beatmapId int) BeatmapResourceProvider {
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

// Cached returns a stream for the given key, serving it from redis when available.
// On a cache miss it calls `fetch`, caches the result under the key with the
// provided ttl & returns a stream over the fetched bytes.
func (provider *BeatmapProvider) Cached(key string, ttl time.Duration, fetch func() (io.ReadCloser, error)) (io.ReadCloser, error) {
	ctx := context.Background()

	if data, err := provider.cache.Get(ctx, key).Bytes(); err == nil && len(data) > 0 {
		provider.logger.Debug("Serving resource from cache", "key", key)
		return io.NopCloser(bytes.NewReader(data)), nil
	}

	stream, err := fetch()
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	data, err := io.ReadAll(stream)
	if err != nil {
		return nil, err
	}

	if err := provider.cache.Set(ctx, key, data, ttl).Err(); err != nil {
		provider.logger.Warn("Failed to cache resource", "key", key, "error", err.Error())
	}

	return io.NopCloser(bytes.NewReader(data)), nil
}
