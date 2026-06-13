package repositories

import (
	"errors"

	"gorm.io/gorm"
)

// Preloaded applies a list of preload strings to a query
func Preloaded(db *gorm.DB, preload []string) *gorm.DB {
	for _, p := range preload {
		db = db.Preload(p)
	}
	return db
}

// CommonUpdate performs an update on the given updates struct, optionally
// selecting specific columns to update.
func CommonUpdate[T any](db *gorm.DB, updates *T, columns ...string) (int64, error) {
	if len(columns) == 0 {
		result := db.Save(updates)
		return result.RowsAffected, result.Error
	}
	result := db.Model(updates).Select(columns).Updates(updates)
	return result.RowsAffected, result.Error
}

// LookupResult converts gorm's ErrRecordNotFound into a nil result
// with no error, which is a common pattern in our repositories.
func LookupResult[T any](schema *T, err error) (*T, error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return schema, err
}
