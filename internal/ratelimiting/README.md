# Rate Limiting

This module provides a system for (burst) ratelimits.
The idea is that a user has a set of "slots" that refill over time. If no slots are available, the ratelimit has been reached.

## Usage

Create a `Limiter` with an appropriate redis client (e.g. `app.Redis`), then call `Allow`:

```go
limiter := ratelimiting.NewRedisLimiter(app.Redis)
limit := ratelimiting.Limit{
	Burst:          12,
	RefillRate:     1,
	RefillInterval: time.Second,
}

key := "ratelimits:example:user:" + strconv.Itoa(user.Id)
allowed, retryAfter, err := limiter.Allow(ctx, key, limit)
if err != nil {
	return err
}
if !allowed {
	return fmt.Errorf("rate limit reached, retry in %s", retryAfter)
}
```

`Allow` uses up one "slot" on success. If no slot is available, it returns the time until the next slot is available.  
An example implementation is stern's [beatmap downloads](../../services/stern/internal/routes/static_osz.go), which uses separate keys for authenticated users & guest IP addresses.
