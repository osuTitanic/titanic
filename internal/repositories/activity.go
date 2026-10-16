package repositories

import (
	"time"

	"github.com/osuTitanic/titanic/internal/schemas"
	"gorm.io/gorm"
)

type BanchoActivityRepository struct {
	db *gorm.DB
}

func NewBanchoActivityRepository(db *gorm.DB) *BanchoActivityRepository {
	return &BanchoActivityRepository{db: db}
}

func (r *BanchoActivityRepository) Create(activity *schemas.BanchoActivity) error {
	return r.db.Create(activity).Error
}

func (r *BanchoActivityRepository) Delete(activity *schemas.BanchoActivity) error {
	return r.db.Delete(activity).Error
}

func (r *BanchoActivityRepository) DeleteOlderThan(cutoff time.Time) (int64, error) {
	result := r.db.Where("time < ?", cutoff).Delete(&schemas.BanchoActivity{})
	return result.RowsAffected, result.Error
}

func (r *BanchoActivityRepository) Update(updates *schemas.BanchoActivity, columns ...string) (int64, error) {
	return CommonUpdate(r.db, updates, columns...)
}

func (r *BanchoActivityRepository) FetchRange(start time.Time, end time.Time) ([]*schemas.BanchoActivity, error) {
	var entries []*schemas.BanchoActivity
	err := r.db.
		Where("time >= ? AND time <= ?", start, end).
		Order("time DESC").
		Find(&entries).Error
	return entries, err
}
