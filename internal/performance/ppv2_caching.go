package performance

import (
	"time"

	"github.com/osuTitanic/titanic/internal/caching"
	"github.com/osuTitanic/titanic/internal/constants"
)

const ppv2CacheExpiry = time.Hour * 24

type PPv2CacheKey struct {
	BeatmapId int
	Mode      constants.Mode
	Mods      constants.Mods
}

type PPv2CacheLayer interface {
	ToCache(key PPv2CacheKey, data any) bool
	FromCache(key PPv2CacheKey) (any, bool)
}

type PPv2MemoryCache struct {
	cache *caching.Cache[PPv2CacheKey, any]
}

func NewPPv2MemoryCache() *PPv2MemoryCache {
	return &PPv2MemoryCache{cache: caching.New[PPv2CacheKey, any](ppv2CacheExpiry)}
}

func (c *PPv2MemoryCache) ToCache(key PPv2CacheKey, data any) bool {
	c.cache.Set(key, data)
	return true // in-memory cache always succeeds
}

func (c *PPv2MemoryCache) FromCache(key PPv2CacheKey) (any, bool) {
	return c.cache.Get(key)
}
