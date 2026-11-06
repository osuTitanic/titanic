package ratelimiting

import (
	"context"
	"time"
)

// Limit specifies how many requests can be made in a given time period.
type Limit struct {
	// Burst is the maximum number of requests allowed without a refill.
	Burst int
	// RefillInterval is the time period after which the request slots are refilled.
	RefillInterval time.Duration
}

// Limiter enforces a rate limit.
type Limiter interface {
	// Allow consumes one request slot if one is available.
	Allow(ctx context.Context, key string, limit Limit) (bool, time.Duration, error)
}
