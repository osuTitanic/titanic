package repositories

import (
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"gorm.io/gorm"
)

type BeatmapPackRepository struct {
	db *gorm.DB
}

func NewBeatmapPackRepository(db *gorm.DB) *BeatmapPackRepository {
	return &BeatmapPackRepository{db: db}
}

func (r *BeatmapPackRepository) Create(pack *schemas.BeatmapPack) error {
	return r.db.Create(pack).Error
}

func (r *BeatmapPackRepository) Delete(pack *schemas.BeatmapPack) error {
	return r.db.Delete(pack).Error
}

func (r *BeatmapPackRepository) Update(updates *schemas.BeatmapPack, columns ...string) (int64, error) {
	return CommonUpdate(r.db, updates, columns...)
}

func (r *BeatmapPackRepository) CreateEntry(entry *schemas.BeatmapPackEntry) error {
	return r.db.Create(entry).Error
}

func (r *BeatmapPackRepository) DeleteEntry(entry *schemas.BeatmapPackEntry) error {
	return r.db.Delete(entry).Error
}

func (r *BeatmapPackRepository) UpdateEntry(updates *schemas.BeatmapPackEntry, columns ...string) (int64, error) {
	return CommonUpdate(
		r.db.Where("pack_id = ? AND beatmapset_id = ?", updates.PackId, updates.BeatmapsetId),
		updates,
		columns...,
	)
}
