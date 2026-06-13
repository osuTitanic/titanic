package database

import (
	"time"

	"github.com/osuTitanic/titanic-go/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// CreateSession opens a postgres connection using values from `config.Config`
func CreateSession(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.PostgresDSN()), &gorm.Config{
		Logger:      NewGormLogger(),
		PrepareStmt: true, // TODO: Benchmark this change to see if it actually improves performance
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if !cfg.PostgresPoolEnabled {
		return db, nil
	}

	// Use configured pool values
	sqlDB.SetMaxOpenConns(cfg.PostgresPoolSizeOverflow)
	sqlDB.SetMaxIdleConns(cfg.PostgresPoolSize)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.PostgresPoolRecycle) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.PostgresPoolTimeout) * time.Second)

	// Warm the pool, if enabled
	if cfg.PostgresPoolPrePing {
		for i := 0; i < cfg.PostgresPoolSize; i++ {
			if err := sqlDB.Ping(); err != nil {
				return nil, err
			}
		}
	}
	return db, nil
}

// CloseSession closes the underlying database connection pool
func CloseSession(db *gorm.DB) error {
	if db == nil {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
