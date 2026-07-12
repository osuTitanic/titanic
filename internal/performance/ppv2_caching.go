package performance

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/osuTitanic/titanic/internal/caching"
	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/redis/go-redis/v9"
)

const ppv2CacheExpiry = time.Hour * 24

type PPv2CacheLayer interface {
	ToCache(key PPv2CacheKey, data any) bool
	FromCache(key PPv2CacheKey) (any, bool)
}

type PPv2CacheKey struct {
	BeatmapId int
	Mode      constants.Mode
	Mods      constants.Mods
}

func (key *PPv2CacheKey) String() string {
	return fmt.Sprintf("%d-%d-%d", key.BeatmapId, key.Mode, key.Mods)
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

type PPv2RedisCache struct {
	client *redis.Client
}

func NewPPv2RedisCache(client *redis.Client) *PPv2RedisCache {
	return &PPv2RedisCache{client: client}
}

func (c *PPv2RedisCache) ToCache(key PPv2CacheKey, data any) bool {
	buffer := bytes.NewBuffer([]byte{})
	encoder := gob.NewEncoder(buffer)
	if err := encoder.Encode(data); err != nil {
		return false
	}

	status := c.client.Set(context.Background(), key.String(), buffer.Bytes(), ppv2CacheExpiry)
	return status.Err() != nil
}

func (c *PPv2RedisCache) FromCache(key PPv2CacheKey) (result any, ok bool) {
	data, err := c.client.Get(context.Background(), key.String()).Bytes()
	if err != nil {
		return nil, false
	}

	reader := bytes.NewReader(data)
	err = gob.NewDecoder(reader).Decode(result)
	return result, err != nil
}
