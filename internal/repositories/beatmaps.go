package repositories

import (
	"errors"
	"time"

	"github.com/osuTitanic/titanic-go/internal/schemas"
	"gorm.io/gorm"
)

type BeatmapRepository struct {
	db *gorm.DB
}

func NewBeatmapRepository(db *gorm.DB) *BeatmapRepository {
	return &BeatmapRepository{db: db}
}

func (r *BeatmapRepository) Create(beatmap *schemas.Beatmap) error {
	return r.db.Create(beatmap).Error
}

func (r *BeatmapRepository) Delete(beatmap *schemas.Beatmap) error {
	return r.db.Delete(beatmap).Error
}

func (r *BeatmapRepository) Update(updates *schemas.Beatmap, columns ...string) (int64, error) {
	return CommonUpdate(r.db, updates, columns...)
}

func (r *BeatmapRepository) UpdateBySetId(updates *schemas.Beatmap, columns ...string) (int64, error) {
	if len(columns) == 0 {
		return 0, errors.New("at least one column must be specified")
	}
	result := r.db.Model(&schemas.Beatmap{}).Where("set_id = ?", updates.SetId).Select(columns).Updates(updates)
	return result.RowsAffected, result.Error
}

func (r *BeatmapRepository) ById(id int, preload ...string) (*schemas.Beatmap, error) {
	var beatmap schemas.Beatmap
	err := Preloaded(r.db, preload).Where("id = ?", id).First(&beatmap).Error
	if err != nil {
		return nil, err
	}
	return &beatmap, nil
}

func (r *BeatmapRepository) ManyById(ids []int, preload ...string) ([]*schemas.Beatmap, error) {
	if len(ids) == 0 {
		return []*schemas.Beatmap{}, nil
	}

	var beatmaps []*schemas.Beatmap
	err := Preloaded(r.db, preload).Where("id IN ?", ids).Find(&beatmaps).Error
	return beatmaps, err
}

func (r *BeatmapRepository) GetCount() (int, error) {
	var count int64
	err := r.db.Model(&schemas.Beatmap{}).Count(&count).Error
	return int(count), err
}

func (r *BeatmapRepository) GetCountGroupedByStatus(mode int) (map[int]int, error) {
	var results []struct {
		Status int
		Count  int
	}

	err := r.db.Model(&schemas.Beatmap{}).
		Select("status, count(*) as count").
		Where("mode = ?", mode).
		Group("status").
		Scan(&results).Error

	counts := make(map[int]int)
	for _, res := range results {
		counts[res.Status] = res.Count
	}

	return counts, err
}

func (r *BeatmapRepository) FetchMostPlayedSince(since time.Time, limit int, preload ...string) (map[int]*schemas.Beatmap, error) {
	type result struct {
		BeatmapId int
		PlayCount int
	}

	var results []result
	err := r.db.Model(&schemas.Score{}).
		Select("scores.beatmap_id, COUNT(scores.id) AS play_count").
		Where("scores.submitted_at >= ?", since).
		Where("scores.hidden = ?", false).
		Group("scores.beatmap_id").
		Order("play_count DESC").
		Limit(limit).
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	beatmapIds := make([]int, 0, len(results))
	for _, result := range results {
		beatmapIds = append(beatmapIds, result.BeatmapId)
	}

	beatmaps, err := r.ManyById(beatmapIds, preload...)
	if err != nil {
		return nil, err
	}

	beatmapsById := make(map[int]*schemas.Beatmap, len(beatmaps))
	for _, beatmap := range beatmaps {
		beatmapsById[beatmap.Id] = beatmap
	}

	mostPlayed := make(map[int]*schemas.Beatmap, len(results))
	for _, result := range results {
		beatmap, ok := beatmapsById[result.BeatmapId]
		if !ok {
			continue
		}
		mostPlayed[result.PlayCount] = beatmap
	}

	return mostPlayed, nil
}
