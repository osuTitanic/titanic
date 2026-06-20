package repositories

import (
	"errors"
	"time"

	"github.com/osuTitanic/titanic-go/internal/constants"
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ScoreRepository struct {
	db *gorm.DB
}

func NewScoreRepository(db *gorm.DB) *ScoreRepository {
	return &ScoreRepository{db: db}
}

func (r *ScoreRepository) Create(score *schemas.Score) error {
	return r.db.Create(score).Error
}

func (r *ScoreRepository) Delete(score *schemas.Score) error {
	return r.db.Delete(score).Error
}

func (r *ScoreRepository) Update(updates *schemas.Score, columns ...string) (int64, error) {
	return CommonUpdate(r.db, updates, columns...)
}

func (r *ScoreRepository) UpdateByBeatmapId(updates *schemas.Score, columns ...string) (int64, error) {
	if len(columns) == 0 {
		return 0, errors.New("at least one column must be specified")
	}
	result := r.db.Model(&schemas.Score{}).Where("beatmap_id = ?", updates.BeatmapId).Select(columns).Updates(updates)
	return result.RowsAffected, result.Error
}

func (r *ScoreRepository) ById(id int64, preload ...string) (*schemas.Score, error) {
	var score schemas.Score
	err := Preloaded(r.db, preload).Where("id = ?", id).First(&score).Error
	if err != nil {
		return nil, err
	}
	return &score, nil
}

func (r *ScoreRepository) ManyById(ids []int64, preload ...string) ([]*schemas.Score, error) {
	if len(ids) == 0 {
		return []*schemas.Score{}, nil
	}

	var scores []*schemas.Score
	err := Preloaded(r.db, preload).Where("id IN ?", ids).Find(&scores).Error
	return scores, err
}

func (r *ScoreRepository) GetCount() (int64, error) {
	var count int64
	err := r.db.Model(&schemas.Score{}).Count(&count).Error
	return count, err
}

func (r *ScoreRepository) FetchScoreIndexById(scoreId int64, beatmapId int, mode constants.Mode) (int, error) {
	rankQuery := `
		SELECT ranked.rank
		FROM (
			SELECT
				id,
				RANK() OVER (ORDER BY total_score DESC) AS rank
			FROM scores
			WHERE beatmap_id = ?
				AND mode = ?
				AND hidden = FALSE
				AND status_score = 3
		) AS ranked
		WHERE ranked.id = ?
		LIMIT 1
	`
	var rank int
	err := r.db.Raw(rankQuery, beatmapId, mode, scoreId).Scan(&rank).Error
	if err != nil {
		return 0, err
	}

	return rank, nil
}

func (r *ScoreRepository) FetchScoreIndexByTscore(totalScore int64, beatmapId int, mode constants.Mode) (int, error) {
	var closestScore schemas.Score
	err := r.db.Model(&schemas.Score{}).
		Where("total_score > ?", totalScore).
		Where("beatmap_id = ?", beatmapId).
		Where("mode = ?", mode).
		Where("status_score = 3").
		Where("hidden = FALSE").
		Order(clause.Expr{SQL: "ABS(total_score - ?)", Vars: []any{totalScore}}).
		Order("id ASC").
		First(&closestScore).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 1, nil
		}
		return 0, err
	}

	rank, err := r.FetchScoreIndexById(closestScore.Id, beatmapId, mode)
	if err != nil {
		return 0, err
	}
	if rank == 0 {
		return 1, nil
	}
	return rank + 1, nil
}

func (r *ScoreRepository) FetchBest(userId int, mode constants.Mode, excludeApproved bool, preload ...string) ([]*schemas.Score, error) {
	allowedStatus := []constants.BeatmapStatus{
		constants.BeatmapStatusRanked,
		constants.BeatmapStatusApproved,
	}

	if !excludeApproved {
		allowedStatus = append(allowedStatus,
			constants.BeatmapStatusQualified,
			constants.BeatmapStatusLoved,
		)
	}

	var scores []*schemas.Score
	err := Preloaded(r.db, preload).
		Joins("JOIN beatmaps ON beatmaps.id = scores.beatmap_id").
		Where("beatmaps.status IN ?", allowedStatus).
		Where("scores.user_id = ?", userId).
		Where("scores.mode = ?", mode).
		Where("scores.status = ?", constants.ScoreStatusBest).
		Where("scores.hidden = ?", false).
		Order("scores.pp DESC").
		Find(&scores).Error

	return scores, err
}

func (r *ScoreRepository) FetchRangeScores(beatmapId int, mode constants.Mode, limit, offset int, preload ...string) ([]*schemas.Score, error) {
	var scores []*schemas.Score
	err := Preloaded(r.db, preload).
		Where("beatmap_id = ?", beatmapId).
		Where("mode = ?", mode).
		Where("status_score = ?", constants.ScoreStatusBest).
		Where("hidden = ?", false).
		Order("total_score DESC").
		Order("id ASC").
		Offset(offset).
		Limit(limit).
		Find(&scores).Error
	return scores, err
}

func (r *ScoreRepository) FetchRangeScoresMods(beatmapId int, mode constants.Mode, mods constants.Mods, limit, offset int, preload ...string) ([]*schemas.Score, error) {
	var scores []*schemas.Score
	err := Preloaded(r.db, preload).
		Where("beatmap_id = ?", beatmapId).
		Where("mode = ?", mode).
		Where("status_score IN ?", []constants.ScoreStatus{constants.ScoreStatusBest, constants.ScoreStatusMods}).
		Where("hidden = ?", false).
		Where("mods = ?", mods).
		Order("total_score DESC").
		Order("id ASC").
		Offset(offset).
		Limit(limit).
		Find(&scores).Error
	return scores, err
}

func (r *ScoreRepository) FetchPersonalBest(beatmapId, userId int, mode constants.Mode, preload ...string) (*schemas.Score, error) {
	var score schemas.Score
	err := Preloaded(r.db, preload).
		Where("beatmap_id = ?", beatmapId).
		Where("user_id = ?", userId).
		Where("mode = ?", mode).
		Where("status_score = ?", constants.ScoreStatusBest).
		Where("hidden = ?", false).
		First(&score).Error
	return LookupResult(&score, err)
}

func (r *ScoreRepository) FetchLeaderCount(userId int, mode constants.Mode) (int, error) {
	leaderCountQuery := `
		SELECT COUNT(*)
		FROM (
			SELECT
				user_id,
				RANK() OVER (PARTITION BY beatmap_id ORDER BY total_score DESC) AS rank
			FROM scores
			WHERE mode = ?
				AND hidden = FALSE
				AND status_score = ?
		) AS ranked
		WHERE ranked.user_id = ?
			AND ranked.rank = 1
	`

	var count int
	err := r.db.Raw(leaderCountQuery, mode, constants.ScoreStatusBest, userId).Scan(&count).Error
	return count, err
}

func (r *ScoreRepository) FetchSubmittedTimestamps(userId int, mode constants.Mode) ([]time.Time, error) {
	timestamps := make([]time.Time, 0)
	err := r.db.Model(&schemas.Score{}).
		Where("hidden = ?", false).
		Where("user_id = ?", userId).
		Where("mode = ?", mode).
		Pluck("submitted_at", &timestamps).Error
	return timestamps, err
}
