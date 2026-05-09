package repositories

import (
	"github.com/osuTitanic/titanic-go/internal/constants"
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"gorm.io/gorm"
)

type BeatmapsetRepository struct {
	db *gorm.DB
}

func NewBeatmapsetRepository(db *gorm.DB) *BeatmapsetRepository {
	return &BeatmapsetRepository{db: db}
}

func (r *BeatmapsetRepository) Create(beatmapset *schemas.Beatmapset) error {
	return r.db.Create(beatmapset).Error
}

func (r *BeatmapsetRepository) Delete(beatmapset *schemas.Beatmapset) error {
	return r.db.Delete(beatmapset).Error
}

func (r *BeatmapsetRepository) Update(updates *schemas.Beatmapset, columns ...string) (int64, error) {
	return CommonUpdate(r.db, updates, columns...)
}

func (r *BeatmapsetRepository) ById(id int, preload ...string) (*schemas.Beatmapset, error) {
	var beatmapset schemas.Beatmapset
	err := Preloaded(r.db, preload).Where("id = ?", id).First(&beatmapset).Error
	if err != nil {
		return nil, err
	}
	return &beatmapset, nil
}

func (r *BeatmapsetRepository) ManyById(ids []int, preload ...string) ([]*schemas.Beatmapset, error) {
	if len(ids) == 0 {
		return []*schemas.Beatmapset{}, nil
	}

	var beatmapsets []*schemas.Beatmapset
	err := Preloaded(r.db, preload).Where("id IN ?", ids).Find(&beatmapsets).Error
	return beatmapsets, err
}

func (r *BeatmapsetRepository) GetCount() (int, error) {
	var count int64
	err := r.db.Model(&schemas.Beatmapset{}).Count(&count).Error
	return int(count), err
}

func (r *BeatmapsetRepository) FetchByStatus(status constants.BeatmapStatus, preload ...string) ([]*schemas.Beatmapset, error) {
	var beatmapsets []*schemas.Beatmapset
	err := Preloaded(r.db, preload).Where("submission_status = ?", status).Find(&beatmapsets).Error
	return beatmapsets, err
}
