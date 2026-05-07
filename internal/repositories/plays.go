package repositories

import (
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"gorm.io/gorm"
)

type BeatmapPlaysRepository struct {
	db *gorm.DB
}

func NewBeatmapPlaysRepository(db *gorm.DB) *BeatmapPlaysRepository {
	return &BeatmapPlaysRepository{db: db}
}

func (r *BeatmapPlaysRepository) Create(plays *schemas.BeatmapPlays) error {
	return r.db.Create(plays).Error
}

func (r *BeatmapPlaysRepository) Delete(plays *schemas.BeatmapPlays) error {
	return r.db.Delete(plays).Error
}

func (r *BeatmapPlaysRepository) Update(updates *schemas.BeatmapPlays, columns ...string) (int64, error) {
	return CommonUpdate(
		r.db.Where("user_id = ? AND beatmap_id = ?", updates.UserId, updates.BeatmapId),
		updates,
		columns...,
	)
}
