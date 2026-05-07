package repositories

import (
	"time"

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

func (r *BeatmapsetRepository) FetchMostPlayedSince(since time.Time, limit int, preload ...string) (map[int]schemas.Beatmapset, error) {
	type result struct {
		SetId     int
		PlayCount int
	}

	var results []result
	err := r.db.Model(&schemas.Score{}).
		Select("beatmaps.set_id, COUNT(scores.id) AS play_count").
		Joins("JOIN beatmaps ON beatmaps.id = scores.beatmap_id").
		Where("scores.submitted_at >= ?", since).
		Where("scores.hidden = ?", false).
		Group("beatmaps.set_id").
		Order("play_count DESC").
		Limit(limit).
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	setIds := make([]int, 0, len(results))
	for _, result := range results {
		setIds = append(setIds, result.SetId)
	}

	beatmapsets, err := r.ManyById(setIds, preload...)
	if err != nil {
		return nil, err
	}

	beatmapsetsById := make(map[int]schemas.Beatmapset, len(beatmapsets))
	for _, beatmapset := range beatmapsets {
		beatmapsetsById[beatmapset.Id] = *beatmapset
	}

	mostPlayed := make(map[int]schemas.Beatmapset, len(results))
	for _, result := range results {
		beatmapset, ok := beatmapsetsById[result.SetId]
		if !ok {
			continue
		}
		mostPlayed[result.PlayCount] = beatmapset
	}

	// TODO: Check if we can put this into a subquery
	return mostPlayed, nil
}
