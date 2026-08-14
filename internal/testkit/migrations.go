//go:build integration

package testkit

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/osuTitanic/titanic/internal/config"
)

// MigratePostgres applies this repository's migrations to a postgres test database.
func MigratePostgres(t testing.TB, cfg *config.Config) {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to resolve repository migrations path")
	}
	migrationsPath := filepath.Join(filepath.Dir(filename), "..", "..", "migrations")

	source, err := iofs.New(os.DirFS(migrationsPath), ".")
	if err != nil {
		t.Fatalf("failed to open repository migrations: %v", err)
	}

	migrator, err := migrate.NewWithSourceInstance(
		"iofs",
		source,
		cfg.PostgresDSN()+"?sslmode=disable",
	)
	if err != nil {
		source.Close()
		t.Fatalf("failed to initialize repository migrations: %v", err)
	}

	migrationErr := migrator.Up()
	sourceCloseErr, databaseCloseErr := migrator.Close()
	if migrationErr != nil && !errors.Is(migrationErr, migrate.ErrNoChange) {
		t.Fatalf("failed to apply repository migrations: %v", migrationErr)
	}
	if sourceCloseErr != nil {
		t.Fatalf("failed to close repository migrations: %v", sourceCloseErr)
	}
	if databaseCloseErr != nil {
		t.Fatalf("failed to close migration database connection: %v", databaseCloseErr)
	}
}
