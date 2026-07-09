//go:build integration

package database_test

import (
	"testing"

	"github.com/osuTitanic/titanic/internal/testkit"
)

func TestCreateSessionWithPostgresContainer(t *testing.T) {
	// `PostgresInstance` will call `CreateSession`, which confirms that
	// the session can be created successfully.
	db := testkit.PostgresInstance(t)

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to access sql db: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("failed to ping postgres: %v", err)
	}

	if err := db.Exec("CREATE TABLE test_table (id int PRIMARY KEY, name text NOT NULL)").Error; err != nil {
		t.Fatalf("failed to create test table: %v", err)
	}
	if err := db.Exec("INSERT INTO test_table (id, name) VALUES (?, ?)", 1, "postgres").Error; err != nil {
		t.Fatalf("failed to insert test row: %v", err)
	}

	var name string
	if err := db.Raw("SELECT name FROM test_table WHERE id = ?", 1).Scan(&name).Error; err != nil {
		t.Fatalf("failed to query test row: %v", err)
	}
	if name != "postgres" {
		t.Fatalf("name = %q, want %q", name, "postgres")
	}
}
