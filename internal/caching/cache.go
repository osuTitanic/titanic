package caching

import (
	"sync"
	"time"
)

type entry[V any] struct {
	value     V
	err       error
	expiresAt time.Time
	loading   bool
	cached    bool
	done      chan struct{}
}

// Loader resolves a value after a cache miss.
type Loader[V any] func() (V, error)

// Cache is an in-process TTL cache.
type Cache[K comparable, V any] struct {
	entries map[K]*entry[V]
	now     func() time.Time
	mu      sync.Mutex
	ttl     time.Duration
}

func New[K comparable, V any](ttl time.Duration) *Cache[K, V] {
	return newWithClock[K, V](ttl, time.Now)
}

func newWithClock[K comparable, V any](ttl time.Duration, now func() time.Time) *Cache[K, V] {
	return &Cache[K, V]{
		entries: make(map[K]*entry[V]),
		ttl:     ttl,
		now:     now,
	}
}

// Get returns a cached value when it exists and has not expired.
func (cache *Cache[K, V]) Get(key K) (V, bool) {
	var zero V
	now := cache.now()

	cache.mu.Lock()
	defer cache.mu.Unlock()

	entry, ok := cache.entries[key]
	if !ok || entry.loading {
		return zero, false
	}
	if !entryValid(entry, now) {
		delete(cache.entries, key)
		return zero, false
	}
	return entry.value, true
}

// GetOrLoad returns a cached value or invokes loader to resolve it.
func (cache *Cache[K, V]) GetOrLoad(key K, loader Loader[V]) (V, error) {
	var zero V

	for {
		now := cache.now()
		cache.mu.Lock()

		if current, ok := cache.entries[key]; ok {
			// If the current entry is loading, wait for it to finish and check if the value is valid
			// Otherwise, if the current entry is valid, return it
			// If not, delete it and continue to load a new value
			if current.loading {
				cache.mu.Unlock()
				<-current.done

				if current.err != nil {
					return zero, current.err
				}
				if current.cached && entryValid(current, cache.now()) {
					return current.value, nil
				}
				continue
			}

			if entryValid(current, now) {
				value := current.value
				cache.mu.Unlock()
				return value, nil
			}
			delete(cache.entries, key)
		}

		// Create a new entry for the key and mark it as loading
		// This will prevent other goroutines from trying to load the same key concurrently
		pending := &entry[V]{
			loading: true,
			done:    make(chan struct{}),
		}
		cache.entries[key] = pending
		cache.mu.Unlock()

		// Load the value (outside of the lock)
		value, err := loader()

		cache.mu.Lock()
		pending.loading = false
		pending.value = value
		pending.err = err

		// If the entry is still the same one we created, mark it as cached and set the expiration time
		if cache.entries[key] == pending {
			if err == nil {
				pending.cached = true
				pending.expiresAt = cache.expiration()
			} else {
				delete(cache.entries, key)
			}
		}
		close(pending.done)
		cache.mu.Unlock()

		return value, err
	}
}

// GetOrCompute is an error-free form of GetOrLoad.
func (cache *Cache[K, V]) GetOrCompute(key K, compute func() V) V {
	value, _ := cache.GetOrLoad(key, func() (V, error) {
		return compute(), nil
	})
	return value
}

// Set stores a value, replacing any existing entry for the key.
func (cache *Cache[K, V]) Set(key K, value V) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.entries[key] = &entry[V]{
		value:     value,
		expiresAt: cache.expiration(),
		cached:    true,
	}
}

// Delete removes a key from the cache.
func (cache *Cache[K, V]) Delete(key K) {
	cache.mu.Lock()
	delete(cache.entries, key)
	cache.mu.Unlock()
}

// Clear removes all cached values.
func (cache *Cache[K, V]) Clear() {
	cache.mu.Lock()
	cache.entries = make(map[K]*entry[V])
	cache.mu.Unlock()
}

// DeleteExpired removes expired entries and returns the number removed.
func (cache *Cache[K, V]) DeleteExpired() int {
	now := cache.now()
	removed := 0

	cache.mu.Lock()
	defer cache.mu.Unlock()

	for key, entry := range cache.entries {
		if entry.loading || entryValid(entry, now) {
			continue
		}
		delete(cache.entries, key)
		removed++
	}
	return removed
}

func (cache *Cache[K, V]) expiration() time.Time {
	if cache.ttl <= 0 {
		return time.Time{}
	}
	return cache.now().Add(cache.ttl)
}

func entryValid[V any](entry *entry[V], now time.Time) bool {
	return entry.cached && (entry.expiresAt.IsZero() || now.Before(entry.expiresAt))
}
