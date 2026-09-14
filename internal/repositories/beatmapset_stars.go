package repositories

import (
	"github.com/osuTitanic/titanic/internal/schemas"
	"gorm.io/gorm"
)

type BeatmapsetStarRepository struct {
	db *gorm.DB
}

func NewBeatmapsetStarRepository(db *gorm.DB) *BeatmapsetStarRepository {
	return &BeatmapsetStarRepository{db: db}
}

func (r *BeatmapsetStarRepository) Create(star *schemas.BeatmapsetStar) error {
	return r.db.Create(star).Error
}

func (r *BeatmapsetStarRepository) Delete(star *schemas.BeatmapsetStar) error {
	return r.db.Delete(star).Error
}

func (r *BeatmapsetStarRepository) BySetId(setId int, preload ...string) ([]*schemas.BeatmapsetStar, error) {
	var stars []*schemas.BeatmapsetStar
	err := Preloaded(r.db, preload).
		Where("set_id = ?", setId).
		Order("created_at ASC").
		Order("id ASC").
		Find(&stars).Error
	return stars, err
}
