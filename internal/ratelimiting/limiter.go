package ratelimiting

import (
	"context"
	"time"
)

// Limit specifies burst capacity and when / how much to refill.
type Limit struct {
	// Burst is the maximum number of requests allowed without a refill.
	Burst int
	// RefillInterval is the period over which RefillRate slots are restored.
	RefillInterval time.Duration
	// RefillRate is the number of slots restored evenly over RefillInterval.
	RefillRate int
}

// Limiter enforces a rate limit.
type Limiter interface {
	// Allow consumes one request slot if one is available.
	Allow(ctx context.Context, key string, limit Limit) (bool, time.Duration, error)
}
