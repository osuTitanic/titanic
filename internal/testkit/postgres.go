//go:build integration

package testkit

import (
	"context"
	"net"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/osuTitanic/titanic/internal/config"
	"github.com/osuTitanic/titanic/internal/database"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"gorm.io/gorm"
)

const (
	postgresImage    = "postgres:17-alpine"
	postgresDB       = "bancho"
	postgresUser     = "bancho"
	postgresPassword = "examplePassword"
)

func PostgresConfig(t testing.TB) *config.Config {
	t.Helper()

	ctx := context.Background()
	container, err := tcpostgres.Run(
		ctx,
		postgresImage,
		tcpostgres.WithDatabase(postgresDB),
		tcpostgres.WithUsername(postgresUser),
		tcpostgres.WithPassword(postgresPassword),
		tcpostgres.BasicWaitStrategies(),
	)
	testcontainers.CleanupContainer(t, container)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to resolve postgres connection string: %v", err)
	}

	parsed, err := url.Parse(dsn)
	if err != nil {
		t.Fatalf("failed to parse postgres connection string: %v", err)
	}

	host, portText, err := net.SplitHostPort(parsed.Host)
	if err != nil {
		t.Fatalf("failed to parse postgres host and port: %v", err)
	}
	port, err := strconv.Atoi(portText)
	if err != nil {
		t.Fatalf("failed to parse postgres port: %v", err)
	}
	password, _ := parsed.User.Password()

	return &config.Config{
		PostgresHost:        host,
		PostgresPort:        port,
		PostgresUser:        parsed.User.Username(),
		PostgresPassword:    password,
		PostgresDatabase:    strings.TrimPrefix(parsed.Path, "/"),
		PostgresPoolEnabled: false,
	}
}

func PostgresInstance(t testing.TB) *gorm.DB {
	t.Helper()

	db, err := database.CreateSession(PostgresConfig(t))
	if err != nil {
		t.Fatalf("failed to create postgres session: %v", err)
	}
	t.Cleanup(func() {
		if err := database.CloseSession(db); err != nil {
			t.Fatalf("failed to close postgres session: %v", err)
		}
	})
	return db
}
