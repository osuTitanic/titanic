package repositories

import (
	"errors"
	"time"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
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
	return LookupResult(&score, err)
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

func (r *ScoreRepository) FetchBest(userId int, mode constants.Mode, excludeApproved bool, preload ...string) ([]*schemas.Score, error) {
	// A negative limit cancels the limit clause entirely -> it will fetch every pb
	return r.FetchBestRange(userId, mode, excludeApproved, -1, 0, preload...)
}

func (r *ScoreRepository) FetchBestRange(userId int, mode constants.Mode, excludeApproved bool, limit, offset int, preload ...string) ([]*schemas.Score, error) {
	var scores []*schemas.Score
	err := bestScoresQuery(userId, mode, excludeApproved, Preloaded(r.db, preload)).
		Order("scores.pp DESC").
		Offset(offset).
		Limit(limit).
		Find(&scores).Error
	return scores, err
}

func (r *ScoreRepository) FetchBestCount(userId int, mode constants.Mode, excludeApproved bool) (int, error) {
	var count int64
	err := bestScoresQuery(userId, mode, excludeApproved, r.db.Model(&schemas.Score{})).
		Count(&count).Error
	return int(count), err
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

func (r *ScoreRepository) FetchLeaderScores(userId int, mode constants.Mode, limit, offset int, preload ...string) ([]*schemas.Score, error) {
	// Only keep scores where no other pb score on the same beatmap has
	// a higher total score, i.e. the user is in first place.
	notBeatenSubquery := `
		NOT EXISTS (
			SELECT 1 FROM scores AS other
			WHERE other.beatmap_id = scores.beatmap_id
				AND other.mode = scores.mode
				AND other.status_score = ?
				AND other.hidden = FALSE
				AND other.total_score > scores.total_score
		)
	`

	var scores []*schemas.Score
	err := Preloaded(r.db, preload).
		Joins("JOIN beatmaps ON beatmaps.id = scores.beatmap_id").
		Where("beatmaps.status > 0").
		Where("scores.user_id = ?", userId).
		Where("scores.mode = ?", mode).
		Where("scores.status_score = ?", constants.ScoreStatusBest).
		Where("scores.hidden = ?", false).
		Where(notBeatenSubquery, constants.ScoreStatusBest).
		Order("scores.id DESC").
		Offset(offset).
		Limit(limit).
		Find(&scores).Error
	return scores, err
}

func (r *ScoreRepository) FetchLeaderCount(userId int, mode constants.Mode) (int, error) {
	// We want to base off the user's own scores instead of ranking all scores
	// in the mode, which is what I was doing previously.
	// For each score, check that no visible best score on the same beatmap/mode
	// has a higher total_score.
	query := `
		SELECT COUNT(DISTINCT s.beatmap_id)
		FROM scores s
		WHERE s.user_id = ?
			AND s.mode = ?
			AND s.hidden = FALSE
			AND s.status_score = ?
			AND NOT EXISTS (
				SELECT 1
				FROM scores better
				WHERE better.beatmap_id = s.beatmap_id
					AND better.mode = s.mode
					AND better.hidden = FALSE
					AND better.status_score = ?
					AND better.total_score > s.total_score
			)
	`

	var count int
	err := r.db.Raw(
		query,
		userId,
		mode,
		constants.ScoreStatusBest,
		constants.ScoreStatusBest,
	).Scan(&count).Error
	return count, err
}

func (r *ScoreRepository) FetchRecentByUser(userId int, mode constants.Mode, limit int, minStatus constants.ScoreStatus, preload ...string) ([]*schemas.Score, error) {
	var scores []*schemas.Score
	err := Preloaded(r.db, preload).
		Where("user_id = ?", userId).
		Where("mode = ?", mode).
		Where("status >= ?", minStatus).
		Where("hidden = ?", false).
		Order("id DESC").
		Limit(limit).
		Find(&scores).Error
	return scores, err
}

func (r *ScoreRepository) FetchPinned(userId int, mode constants.Mode, limit, offset int, preload ...string) ([]*schemas.Score, error) {
	var scores []*schemas.Score
	err := pinnedQuery(userId, mode, Preloaded(r.db, preload)).
		Order("pp DESC").
		Offset(offset).
		Limit(limit).
		Find(&scores).Error
	return scores, err
}

func (r *ScoreRepository) FetchPinnedCount(userId int, mode constants.Mode) (int, error) {
	var count int64
	err := pinnedQuery(userId, mode, r.db.Model(&schemas.Score{})).
		Count(&count).Error
	return int(count), err
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

func bestScoresQuery(userId int, mode constants.Mode, excludeApproved bool, query *gorm.DB) *gorm.DB {
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

	return query.
		Joins("JOIN beatmaps ON beatmaps.id = scores.beatmap_id").
		Where("beatmaps.status IN ?", allowedStatus).
		Where("scores.user_id = ?", userId).
		Where("scores.mode = ?", mode).
		Where("scores.status = ?", constants.ScoreStatusBest).
		Where("scores.hidden = ?", false)
}

func pinnedQuery(userId int, mode constants.Mode, query *gorm.DB) *gorm.DB {
	return query.
		Where("user_id = ?", userId).
		Where("mode = ?", mode).
		Where("status > ?", 1).
		Where("hidden = ?", false).
		Where("pinned = ?", true)
}
