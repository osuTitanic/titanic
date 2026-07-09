package repositories

import (
	"time"

	"github.com/osuTitanic/titanic/internal/schemas"
	"gorm.io/gorm"
)

type ActivityRepository struct {
	db *gorm.DB
}

func NewActivityRepository(db *gorm.DB) *ActivityRepository {
	return &ActivityRepository{db: db}
}

func (r *ActivityRepository) Create(activity *schemas.Activity) error {
	return r.db.Create(activity).Error
}

func (r *ActivityRepository) Delete(activity *schemas.Activity) error {
	return r.db.Delete(activity).Error
}

func (r *ActivityRepository) Update(updates *schemas.Activity, columns ...string) (int64, error) {
	return CommonUpdate(r.db, updates, columns...)
}

func (r *ActivityRepository) FetchRecentByUser(userId int, mode int, limit int, offset int, until time.Duration, preload ...string) ([]*schemas.Activity, error) {
	var activities []*schemas.Activity
	err := Preloaded(r.db, preload).
		Where("user_id = ?", userId).
		Where("mode = ?", mode).
		Where("time > ?", time.Now().Add(-until)).
		Where("hidden = ?", false).
		Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&activities).Error
	return activities, err
}
