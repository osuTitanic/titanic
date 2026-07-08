# Test Kit

This module contains helpers for tests that need external services. It currently starts temporary PostgreSQL and Redis instances through [testcontainers](https://golang.testcontainers.org/).

## Integration Tests

Integration tests are only available behind the `integration` build tag.
Normal test runs don't use docker:

```bash
go test -tags=integration ./...
```

This means that Docker must be available to the test process when using these helpers.

## PostgreSQL

Use `PostgresInstance` when a test needs a `*gorm.DB` session. The helper starts a fresh PostgreSQL container and registers binds with the test.

```go
db := testkit.PostgresInstance(t)

if err := db.Exec("SELECT 1").Error; err != nil {
	t.Fatal(err)
}
```

Use `PostgresConfig` when a test needs a `config.Config` pointing at the container instead of an already-open session.

## Redis

Use `RedisClient` when a test needs a `*redis.Client`. The helper starts a fresh Redis container, verifies it with `PING`, and binds cleanup with the test.

```go
client := testkit.RedisClient(t)

if err := client.Set(context.Background(), "key", "value", 0).Err(); err != nil {
	t.Fatal(err)
}
```
