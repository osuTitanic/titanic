package repositories

import (
	"errors"
	"fmt"
	"time"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
	"gorm.io/gorm"
)

type BeatmapRepository struct {
	db *gorm.DB
}

type BeatmapPlayCount struct {
	Beatmap *schemas.Beatmap
	Count   int
}

type BeatmapStatusResult struct {
	Id       int
	SetId    int
	TopicId  *int
	Checksum string
	Status   constants.BeatmapStatus
}

type BeatmapInfoResult struct {
	Id         int
	SetId      int
	Filename   string
	Checksum   string
	Status     constants.BeatmapStatus
	GradeOsu   constants.Grade
	GradeTaiko constants.Grade
	GradeCatch constants.Grade
	GradeMania constants.Grade
}

func (result BeatmapInfoResult) Format(index int) string {
	return fmt.Sprintf(
		"%d|%d|%d|%s|%d|%s|%s|%s|%s",
		index,
		result.Id,
		result.SetId,
		result.Checksum,
		max(result.Status.Value(), 0),
		result.GradeOsu,
		result.GradeTaiko,
		result.GradeCatch,
		result.GradeMania,
	)
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

func (r *BeatmapRepository) UpdateByCriteria(criteria map[string]any, updates *schemas.Beatmap, columns ...string) (int64, error) {
	if len(columns) == 0 {
		return 0, errors.New("at least one column must be specified")
	}
	result := r.db.Model(&schemas.Beatmap{}).Where(criteria).Select(columns).Updates(updates)
	return result.RowsAffected, result.Error
}

func (r *BeatmapRepository) ById(id int, preload ...string) (*schemas.Beatmap, error) {
	var beatmap schemas.Beatmap
	err := Preloaded(r.db, preload).Where("id = ?", id).First(&beatmap).Error
	return LookupResult(&beatmap, err)
}

func (r *BeatmapRepository) ByFilename(filename string, preload ...string) (*schemas.Beatmap, error) {
	var beatmap schemas.Beatmap
	err := Preloaded(r.db, preload).Where("filename = ?", filename).First(&beatmap).Error
	return LookupResult(&beatmap, err)
}

func (r *BeatmapRepository) ByChecksum(checksum string, preload ...string) (*schemas.Beatmap, error) {
	var beatmap schemas.Beatmap
	err := Preloaded(r.db, preload).Where("md5 = ?", checksum).First(&beatmap).Error
	return LookupResult(&beatmap, err)
}

func (r *BeatmapRepository) Many(critera map[string]any, preload ...string) ([]*schemas.Beatmap, error) {
	var beatmaps []*schemas.Beatmap
	query := Preloaded(r.db, preload)

	for key, value := range critera {
		query = query.Where(key, value)
	}

	err := query.Find(&beatmaps).Error
	return beatmaps, err
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

func (r *BeatmapRepository) FetchMostPlayedSince(since time.Time, limit int, preload ...string) ([]BeatmapPlayCount, error) {
	var results []struct {
		BeatmapId int
		PlayCount int
	}
	err := r.db.Model(&schemas.Score{}).
		Select("scores.beatmap_id, COUNT(scores.id) AS play_count").
		Where("scores.submitted_at >= ?", since).
		Where("scores.hidden = ?", false).
		Group("scores.beatmap_id").
		Order("play_count DESC").
		Order("scores.beatmap_id ASC").
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

	mostPlayed := make([]BeatmapPlayCount, 0, len(results))
	for _, result := range results {
		beatmap, ok := beatmapsById[result.BeatmapId]
		if !ok {
			continue
		}
		mostPlayed = append(mostPlayed, BeatmapPlayCount{
			Beatmap: beatmap,
			Count:   result.PlayCount,
		})
	}

	return mostPlayed, nil
}

func (r *BeatmapRepository) FetchInfoByFilenamesOrIds(userId int, filenames []string, ids []int) ([]BeatmapInfoResult, error) {
	if len(filenames) == 0 && len(ids) == 0 {
		return []BeatmapInfoResult{}, nil
	}

	query := r.db.Model(&schemas.Beatmap{}).
		Select(`
			beatmaps.id,
			beatmaps.set_id,
			beatmaps.filename,
			beatmaps.md5 AS checksum,
			beatmaps.status,
			COALESCE(MAX(scores.grade) FILTER (WHERE scores.mode = ?), ?) AS grade_osu,
			COALESCE(MAX(scores.grade) FILTER (WHERE scores.mode = ?), ?) AS grade_taiko,
			COALESCE(MAX(scores.grade) FILTER (WHERE scores.mode = ?), ?) AS grade_catch,
			COALESCE(MAX(scores.grade) FILTER (WHERE scores.mode = ?), ?) AS grade_mania
		`,
			constants.ModeOsu, constants.GradeN,
			constants.ModeTaiko, constants.GradeN,
			constants.ModeCatch, constants.GradeN,
			constants.ModeMania, constants.GradeN,
		).
		Joins(`
			LEFT JOIN scores ON scores.beatmap_id = beatmaps.id
				AND scores.user_id = ?
				AND scores.status = ?
				AND scores.hidden = FALSE
		`, userId, constants.ScoreStatusBest).
		Where("beatmaps.status > ?", constants.BeatmapStatusInactive).
		Group("beatmaps.id")

	switch {
	case len(filenames) == 0:
		query = query.Where("beatmaps.id IN ?", ids)
	case len(ids) == 0:
		query = query.Where("beatmaps.filename IN ?", filenames)
	default:
		query = query.Where("(beatmaps.filename IN ? OR beatmaps.id IN ?)", filenames, ids)
	}

	var results []BeatmapInfoResult
	err := query.Scan(&results).Error
	return results, err
}

func (r *BeatmapRepository) FetchStatusesByChecksums(checksums []string) ([]BeatmapStatusResult, error) {
	if len(checksums) == 0 {
		return []BeatmapStatusResult{}, nil
	}

	var results []BeatmapStatusResult
	err := r.db.Model(&schemas.Beatmap{}).
		Select("beatmaps.id, beatmaps.set_id, beatmaps.md5 AS checksum, beatmaps.status, beatmapsets.topic_id").
		Joins("JOIN beatmapsets ON beatmapsets.id = beatmaps.set_id").
		Where("beatmaps.md5 IN ?", checksums).
		Scan(&results).Error
	return results, err
}
