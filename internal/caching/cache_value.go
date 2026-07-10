package caching

import "time"

type valueKey struct{}

// Value is a cache for a single value.
type Value[V any] struct {
	cache *Cache[valueKey, V]
}

func NewValue[V any](ttl time.Duration) *Value[V] {
	return &Value[V]{cache: New[valueKey, V](ttl)}
}

// Get returns the cached value when it exists and has not expired.
func (value *Value[V]) Get() (V, bool) {
	return value.cache.Get(valueKey{})
}

// GetOrLoad returns the cached value or invokes loader to resolve it.
func (value *Value[V]) GetOrLoad(loader Loader[V]) (V, error) {
	return value.cache.GetOrLoad(valueKey{}, loader)
}

// GetOrCompute is an error-free form of GetOrLoad.
func (value *Value[V]) GetOrCompute(compute func() V) V {
	return value.cache.GetOrCompute(valueKey{}, compute)
}

// Set replaces the cached value.
func (value *Value[V]) Set(cached V) {
	value.cache.Set(valueKey{}, cached)
}

// Invalidate removes the cached value.
func (value *Value[V]) Invalidate() {
	value.cache.Delete(valueKey{})
}
