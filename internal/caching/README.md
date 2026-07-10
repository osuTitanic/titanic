# Caching

This module provides in-process TTL caches for data that is inexpensive to keep in memory but expensive to load or generate repeatedly.

Two cache types are available:

- `Cache[K, V]` stores multiple values by key
- `Value[V]` is a convenience wrapper for caching one value

Both types combine simultaneous cache misses for the same key. Only one "loader" runs, while the other callers wait for its result.

## Keyed cache

Create a cache with a TTL, then use `GetOrLoad` to populate it on demand:

```go
users := caching.New[int, *schemas.User](5 * time.Minute)

user, err := users.GetOrLoad(userId, func() (*schemas.User, error) {
	return app.Users.ById(userId)
})
if err != nil {
	return err
}
```

Loader errors are returned to all callers waiting for that user to load. Errors are not retained though, so subsequent calls will attempt to load the value again.

Use `GetOrCompute` when generating a value cannot fail:

```go
labels := caching.New[string, []string](time.Hour)

result := labels.GetOrCompute(language, func() []string {
	return buildLabels(language)
})
```

## Single-value cache

Use `Value` when there is no meaningful cache key:

```go
chart := caching.NewValue[[]byte](time.Minute)

image, err := chart.GetOrLoad(func() ([]byte, error) {
	return generateChart()
})
if err != nil {
	return err
}
```

`Value` similarly provides `Get`, `GetOrCompute`, `Set`, and `Invalidate`.

## Expiration and invalidation

The TTL begins after a value loads successfully. A non-positive TTL creates entries that remain valid indefinitely, until explicitly removed.

Expired entries are removed automatically when their key is accessed. For caches with many short-lived or rarely reused keys, call `DeleteExpired` periodically to release inactive entries.
