package performance

import (
	"time"

	"github.com/osuTitanic/titanic/internal/caching"
)

const ppv2CacheExpiry = time.Hour * 24

type PPv2CacheLayer interface {
	ToCache(beatmapId int, data any) bool
	FromCache(beatmapId int) (any, bool)
}

type PPv2MemoryCache struct {
	cache *caching.Cache[int, any]
}

func NewPPv2MemoryCache() *PPv2MemoryCache {
	return &PPv2MemoryCache{cache: caching.New[int, any](ppv2CacheExpiry)}
}

func (c *PPv2MemoryCache) ToCache(beatmapId int, data any) bool {
	c.cache.Set(beatmapId, data)
	return true // in-memory cache always succeeds
}

func (c *PPv2MemoryCache) FromCache(beatmapId int) (any, bool) {
	return c.cache.Get(beatmapId)
}
