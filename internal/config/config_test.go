package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("POSTGRES_PASSWORD", "testpass")
	defer os.Unsetenv("POSTGRES_PASSWORD")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.PostgresPassword != "testpass" {
		t.Errorf("PostgresPassword = %q, want %q", cfg.PostgresPassword, "testpass")
	}
}

func TestLoadConfig_MissingRequired(t *testing.T) {
	os.Unsetenv("POSTGRES_PASSWORD")

	_, err := LoadConfig()
	if err == nil {
		t.Error("LoadConfig() expected error for missing required field")
	}
}
