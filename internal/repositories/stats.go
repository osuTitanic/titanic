package repositories

import (
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"gorm.io/gorm"
)

type StatsRepository struct {
	db *gorm.DB
}

func NewStatsRepository(db *gorm.DB) *StatsRepository {
	return &StatsRepository{db: db}
}

func (r *StatsRepository) Create(stats *schemas.Stats) error {
	return r.db.Create(stats).Error
}

func (r *StatsRepository) Delete(stats *schemas.Stats) error {
	return r.db.Delete(stats).Error
}

func (r *StatsRepository) Update(updates *schemas.Stats, columns ...string) (int64, error) {
	return CommonUpdate(
		r.db.Where("id = ? AND mode = ?", updates.UserId, updates.Mode),
		updates,
		columns...,
	)
}

func (r *StatsRepository) ByMode(userId int, mode int, preload ...string) (*schemas.Stats, error) {
	var stats schemas.Stats
	err := Preloaded(r.db, preload).Where("id = ? AND mode = ?", userId, mode).First(&stats).Error
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

func (r *StatsRepository) ManyByUserId(userId int, preload ...string) ([]*schemas.Stats, error) {
	var stats []*schemas.Stats
	err := Preloaded(r.db, preload).Where("id = ?", userId).Find(&stats).Error
	return stats, err
}
