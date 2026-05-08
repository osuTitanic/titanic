# Database Sessions

This module manages the PostgreSQL database session used by the app, housing on [GORM](https://gorm.io/) with the PostgreSQL driver.

## Usage with state system

Most services should use `state.NewState(...)` instead of creating a database session directly.
`State` manages the database lifecycle and calls `database.CloseSession` from `state.Close()`.

```go
app, err := state.NewState()
if err != nil {
	return err
}
defer app.Close()
// Proceed to use `app.Database` ...
```

## Usage without state system

Create a database session from the config & close it when the service shuts down.

```go
cfg, err := config.LoadConfig()
if err != nil {
	panic(err)
}

db, err := database.CreateSession(cfg)
if err != nil {
	panic(err)
}
defer database.CloseSession(db)
```

---

Internally, `CreateSession` uses `cfg.PostgresDSN()` for the connection string & installs the custom `slog`-backed gorm logger.

```go
db, err := gorm.Open(postgres.Open(cfg.PostgresDSN()), &gorm.Config{
	Logger: database.NewGormLogger(),
})
```

## Pooling

Connection pooling is controlled by the postgres pool fields in `config.Config`.

- `POSTGRES_POOL_ENABLED` enables pool configuration
- `POSTGRES_POOL_SIZE_OVERFLOW` sets the maximum number of open connections
- `POSTGRES_POOL_SIZE` sets the maximum number of idle connections
- `POSTGRES_POOL_RECYCLE` sets the maximum connection lifetime in seconds
- `POSTGRES_POOL_TIMEOUT` sets the maximum idle time in seconds
- `POSTGRES_POOL_PRE_PING` pings the database during startup
