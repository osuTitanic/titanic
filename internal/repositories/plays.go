package repositories

import (
	"github.com/osuTitanic/titanic/internal/schemas"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *BeatmapPlaysRepository) Increment(plays *schemas.BeatmapPlays) error {
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}, {Name: "beatmap_id"}},
		DoUpdates: clause.Assignments(map[string]any{
			"count":        gorm.Expr("plays.count + 1"),
			"set_id":       plays.SetId,
			"beatmap_file": plays.BeatmapFile,
		}),
	}).Create(plays).Error
}

func (r *BeatmapPlaysRepository) FetchMostPlayedByUser(userId int, limit, offset int, preload ...string) ([]*schemas.BeatmapPlays, error) {
	var plays []*schemas.BeatmapPlays
	err := Preloaded(r.db, preload).
		Where("user_id = ?", userId).
		Order("count DESC").
		Offset(offset).
		Limit(limit).
		Find(&plays).Error
	return plays, err
}
