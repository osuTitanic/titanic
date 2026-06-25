package repositories

import (
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"gorm.io/gorm"
)

type BeatmapNominationRepository struct {
	db *gorm.DB
}

func NewNominationRepository(db *gorm.DB) *BeatmapNominationRepository {
	return &BeatmapNominationRepository{db: db}
}

func (r *BeatmapNominationRepository) Create(nomination *schemas.BeatmapNomination) error {
	return r.db.Create(nomination).Error
}

func (r *BeatmapNominationRepository) Delete(nomination *schemas.BeatmapNomination) error {
	return r.db.Delete(nomination).Error
}

func (r *BeatmapNominationRepository) DeleteAll(setId int) error {
	return r.db.Where("set_id = ?", setId).Delete(&schemas.BeatmapNomination{}).Error
}

func (r *BeatmapNominationRepository) FetchBySet(setId int, preload ...string) ([]*schemas.BeatmapNomination, error) {
	var nominations []*schemas.BeatmapNomination
	err := Preloaded(r.db, preload).Where("set_id = ?", setId).Find(&nominations).Error
	return nominations, err
}

func (r *BeatmapNominationRepository) FetchByUser(userId int, preload ...string) ([]*schemas.BeatmapNomination, error) {
	var nominations []*schemas.BeatmapNomination
	err := Preloaded(r.db, preload).
		Where("user_id = ?", userId).
		Order("time DESC").
		Find(&nominations).Error
	return nominations, err
}
