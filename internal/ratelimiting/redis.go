package ratelimiting

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// Adapted code from https://github.com/go-redis/redis_rate/tree/v10
// Copyright (c) 2013 The github.com/go-redis/redis_rate Authors. All rights reserved.

type RedisLimiter struct {
	client redis.Scripter
}

func NewRedisLimiter(client redis.Scripter) *RedisLimiter {
	return &RedisLimiter{client: client}
}

func (limiter *RedisLimiter) Allow(ctx context.Context, key string, limit Limit) (bool, time.Duration, error) {
	if limit.Burst <= 0 {
		return false, 0, errors.New("rate limit requires a positive burst")
	}

	rate := limit.RefillRate
	if rate < 0 {
		return false, 0, errors.New("rate limit requires a non-negative refill rate")
	}
	if rate == 0 {
		rate = 1
	}
	if limit.RefillInterval/time.Duration(rate) < time.Microsecond {
		return false, 0, errors.New("rate limit requires at least 1 microsecond between slot refills")
	}

	result, err := allowN.Run(
		ctx,
		limiter.client,
		[]string{"ratelimiting:burst:" + key},
		limit.Burst,                    // ARGV[1] (burst) is the maximum number of available slots
		rate,                           // ARGV[2] (rate) is the number of slots restored per period
		limit.RefillInterval.Seconds(), // ARGV[3] (period) is the refill interval in seconds
		1,                              // ARGV[4] (cost) consumes one slot per request
	).Slice()
	if err != nil {
		return false, 0, fmt.Errorf("apply redis rate limit: %w", err)
	}
	if len(result) != 4 {
		return false, 0, fmt.Errorf("unexpected redis rate limit result length: %d", len(result))
	}

	allowed, ok := result[0].(int64)
	if !ok {
		return false, 0, fmt.Errorf("unexpected redis rate limit allowance type: %T", result[0])
	}
	if allowed != 0 && allowed != 1 {
		return false, 0, fmt.Errorf("unexpected redis rate limit allowance: %v", result[0])
	}
	if allowed == 1 {
		return true, 0, nil
	}

	retry, ok := result[2].(string)
	if !ok {
		return false, 0, fmt.Errorf("unexpected redis rate limit retry delay: %v", result[2])
	}

	seconds, err := strconv.ParseFloat(retry, 64)
	if err != nil {
		return false, 0, fmt.Errorf("invalid redis rate limit retry delay: %q", retry)
	}
	if seconds <= 0 || math.IsInf(seconds, 0) {
		return false, 0, fmt.Errorf("invalid redis rate limit retry delay: %q", retry)
	}

	nanoseconds := math.Ceil(seconds * float64(time.Second))
	if nanoseconds >= float64(math.MaxInt64) {
		return false, time.Duration(math.MaxInt64), nil
	}

	return false, time.Duration(nanoseconds), nil
}
