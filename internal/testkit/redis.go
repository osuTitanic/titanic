//go:build integration

package testkit

import (
	"context"
	"testing"

	goredis "github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"
)

const redisImage = "redis:7-alpine"

func RedisClient(t testing.TB) *goredis.Client {
	t.Helper()

	ctx := context.Background()
	container, err := tcredis.Run(ctx, redisImage)
	testcontainers.CleanupContainer(t, container)
	if err != nil {
		t.Fatalf("failed to start redis container: %v", err)
	}

	dsn, err := container.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("failed to resolve redis connection string: %v", err)
	}

	options, err := goredis.ParseURL(dsn)
	if err != nil {
		t.Fatalf("failed to parse redis connection string: %v", err)
	}
	options.TLSConfig = container.TLSConfig()

	client := goredis.NewClient(options)
	t.Cleanup(func() {
		if err := client.Close(); err != nil {
			t.Fatalf("failed to close redis client: %v", err)
		}
	})

	if err := client.Ping(ctx).Err(); err != nil {
		t.Fatalf("failed to ping redis container: %v", err)
	}
	return client
}
