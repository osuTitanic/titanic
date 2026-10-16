package repositories

import (
	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
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

func (r *BeatmapNominationRepository) Exists(setId int, userId int) (bool, error) {
	var count int64
	err := r.db.Model(&schemas.BeatmapNomination{}).
		Where("set_id = ? AND user_id = ?", setId, userId).
		Count(&count).Error
	return count > 0, err
}

func (r *BeatmapNominationRepository) ExistsForSet(setId int) (bool, error) {
	var count int64
	err := r.db.Model(&schemas.BeatmapNomination{}).
		Where("set_id = ?", setId).
		Count(&count).Error
	return count > 0, err
}

func (r *BeatmapNominationRepository) FetchBySet(setId int, preload ...string) ([]*schemas.BeatmapNomination, error) {
	var nominations []*schemas.BeatmapNomination
	err := Preloaded(r.db, preload).Where("set_id = ?", setId).Find(&nominations).Error
	return nominations, err
}

func (r *BeatmapNominationRepository) FetchByUserApproved(userId int, preload ...string) ([]*schemas.BeatmapNomination, error) {
	var nominations []*schemas.BeatmapNomination
	err := Preloaded(r.db, preload).
		Joins("JOIN beatmapsets ON beatmapsets.id = beatmap_nominations.set_id").
		Where("beatmap_nominations.user_id = ?", userId).
		Where("beatmapsets.submission_status IN ?", []constants.BeatmapStatus{
			constants.BeatmapStatusRanked,
			constants.BeatmapStatusLoved,
			constants.BeatmapStatusQualified,
			constants.BeatmapStatusApproved,
		}).
		Order("beatmap_nominations.time DESC").
		Find(&nominations).Error
	return nominations, err
}

func (r *BeatmapNominationRepository) FetchByUserPending(userId int, preload ...string) ([]*schemas.BeatmapNomination, error) {
	var nominations []*schemas.BeatmapNomination
	err := Preloaded(r.db, preload).
		Joins("JOIN beatmapsets ON beatmapsets.id = beatmap_nominations.set_id").
		Where("beatmap_nominations.user_id = ?", userId).
		Where("beatmapsets.submission_status NOT IN ?", []constants.BeatmapStatus{
			constants.BeatmapStatusRanked,
			constants.BeatmapStatusLoved,
			constants.BeatmapStatusQualified,
			constants.BeatmapStatusApproved,
		}).
		Order("beatmap_nominations.time DESC").
		Find(&nominations).Error
	return nominations, err
}
