package repositories

import (
	"github.com/osuTitanic/titanic/internal/schemas"
	"gorm.io/gorm"
)

type BenchmarkRepository struct {
	db *gorm.DB
}

func NewBenchmarkRepository(db *gorm.DB) *BenchmarkRepository {
	return &BenchmarkRepository{db: db}
}

func (r *BenchmarkRepository) Create(benchmark *schemas.Benchmark) error {
	return r.db.Create(benchmark).Error
}

func (r *BenchmarkRepository) Delete(benchmark *schemas.Benchmark) error {
	return r.db.Delete(benchmark).Error
}

func (r *BenchmarkRepository) Update(updates *schemas.Benchmark, columns ...string) (int64, error) {
	return CommonUpdate(r.db, updates, columns...)
}
