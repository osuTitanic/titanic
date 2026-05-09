package repositories

import (
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"gorm.io/gorm"
)

type BeatmapRatingRepository struct {
	db *gorm.DB
}

func NewBeatmapRatingRepository(db *gorm.DB) *BeatmapRatingRepository {
	return &BeatmapRatingRepository{db: db}
}

func (r *BeatmapRatingRepository) Create(rating *schemas.BeatmapRating) error {
	return r.db.Create(rating).Error
}

func (r *BeatmapRatingRepository) Delete(rating *schemas.BeatmapRating) error {
	return r.db.Delete(rating).Error
}

func (r *BeatmapRatingRepository) Update(updates *schemas.BeatmapRating, columns ...string) (int64, error) {
	return CommonUpdate(
		r.db.Where("user_id = ? AND map_checksum = ?", updates.UserId, updates.MapChecksum),
		updates,
		columns...,
	)
}
