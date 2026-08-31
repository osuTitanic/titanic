package repositories

import (
	"errors"
	"time"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
	"gorm.io/gorm"
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

func (r *ScoreRepository) Many(critera map[string]any, order string, offset, limit int, preload ...string) ([]*schemas.Score, error) {
	var scores []*schemas.Score
	query := Preloaded(r.db, preload)

	for key, value := range critera {
		query = query.Where(key, value)
	}
	if order != "" {
		query = query.Order(order)
	}

	if offset >= 0 {
		query = query.Offset(offset)
	}
	if limit >= 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&scores).Error
	return scores, err
}

func (r *ScoreRepository) Count(critera map[string]any) (int64, error) {
	var count int64
	query := r.db.Model(&schemas.Score{})

	for key, value := range critera {
		query = query.Where(key, value)
	}
	err := query.Count(&count).Error
	return count, err
}

// probably a bit weird to have Count and GetCount but who cares lol

func (r *ScoreRepository) GetCount() (int64, error) {
	var count int64
	err := r.db.Model(&schemas.Score{}).Count(&count).Error
	return count, err
}

func (r *ScoreRepository) FetchWithZeroPP(preload ...string) ([]*schemas.Score, error) {
	var scores []*schemas.Score
	err := Preloaded(r.db, preload).
		Where("pp = ?", 0).
		Where("hidden = ?", false).
		Where("status >= ?", constants.ScoreStatusSubmitted).
		Find(&scores).Error
	return scores, err
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

func (r *ScoreRepository) FetchPassed(userId int, mode constants.Mode, preload ...string) ([]*schemas.Score, error) {
	var scores []*schemas.Score
	err := Preloaded(r.db, preload).
		Where("user_id = ?", userId).
		Where("mode = ?", mode).
		Where("status > 1").
		Where("hidden = ?", false).
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
		Order("total_score DESC, submitted_at ASC, id ASC").
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
		Order("total_score DESC, submitted_at ASC, id ASC").
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

func (r *ScoreRepository) FetchScoreIndex(score *schemas.Score) (scoreRank int, err error) {
	if score.Id > 0 && score.StatusScore == constants.ScoreStatusBest && !score.Hidden {
		// Score exists & is a pb -> fetch its rank on the leaderboard by ID
		scoreRank, err = r.FetchScoreIndexById(
			score.Id, score.BeatmapId, score.Mode,
		)
	} else {
		// Score is not a pb / does not exist yet -> fetch the potential rank
		scoreRank, err = r.FetchScoreIndexByTscore(
			score.TotalScore, score.SubmittedAt,
			score.BeatmapId, score.Mode,
		)
	}
	return scoreRank, err
}

// FetchScoreIndexById fetches the rank of a score on the leaderboard, which is expected to be a pb score.
// It will return 0 if the score is not found on the leaderboard.
func (r *ScoreRepository) FetchScoreIndexById(scoreId int64, beatmapId int, mode constants.Mode) (int, error) {
	rankQuery := `
		SELECT rank
		FROM (
			SELECT
				id,
				ROW_NUMBER() OVER (
					ORDER BY total_score DESC, submitted_at ASC, id ASC
				) AS rank
			FROM scores
			WHERE beatmap_id = ?
				AND mode = ?
				AND hidden = FALSE
				AND status_score = ?
		) AS ranked
		WHERE id = ?
	`
	var rank int
	err := r.db.Raw(
		rankQuery,
		beatmapId,
		mode,
		constants.ScoreStatusBest,
		scoreId,
	).Scan(&rank).Error
	if err != nil {
		return 0, err
	}
	return rank, nil
}

// FetchScoreIndexByTscore fetches the theoretical rank of a score on the leaderboard based on its total score & submission time.
// It should be used for scores that are not yet on the leaderboard / not a pb.
func (r *ScoreRepository) FetchScoreIndexByTscore(totalScore int64, submittedAt time.Time, beatmapId int, mode constants.Mode) (int, error) {
	var precedingScores int64
	err := r.db.Model(&schemas.Score{}).
		Where("beatmap_id = ?", beatmapId).
		Where("mode = ?", mode).
		Where("hidden = ?", false).
		Where("status_score = ?", constants.ScoreStatusBest).
		Where(
			"total_score > ? OR (total_score = ? AND submitted_at <= ?)",
			totalScore,
			totalScore,
			submittedAt,
		).
		Count(&precedingScores).Error
	if err != nil {
		return 0, err
	}
	return int(precedingScores) + 1, nil
}

func (r *ScoreRepository) FetchLeaderScores(userId int, mode constants.Mode, limit, offset int, preload ...string) ([]*schemas.Score, error) {
	leaderSubquery := `
		scores.id = (
			SELECT leader.id
			FROM scores AS leader
			WHERE leader.beatmap_id = scores.beatmap_id
				AND leader.mode = scores.mode
				AND leader.status_score = ?
				AND leader.hidden = FALSE
			ORDER BY leader.total_score DESC, leader.submitted_at ASC, leader.id ASC
			LIMIT 1
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
		Where(leaderSubquery, constants.ScoreStatusBest).
		Order("scores.id DESC").
		Offset(offset).
		Limit(limit).
		Find(&scores).Error
	return scores, err
}

func (r *ScoreRepository) FetchLeaderCount(userId int, mode constants.Mode) (int, error) {
	query := `
		SELECT COUNT(DISTINCT s.beatmap_id)
		FROM scores s
		JOIN beatmaps ON beatmaps.id = s.beatmap_id
		WHERE beatmaps.status > ?
			AND s.user_id = ?
			AND s.mode = ?
			AND s.hidden = FALSE
			AND s.status_score = ?
			AND s.id = (
				SELECT leader.id
				FROM scores AS leader
				WHERE leader.beatmap_id = s.beatmap_id
					AND leader.mode = s.mode
					AND leader.hidden = FALSE
					AND leader.status_score = ?
				ORDER BY leader.total_score DESC, leader.submitted_at ASC, leader.id ASC
				LIMIT 1
			)
	`

	var count int
	err := r.db.Raw(
		query,
		constants.BeatmapStatusPending,
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

func (r *ScoreRepository) FetchMostWatchedByUser(userId int, mode constants.Mode, limit int, preload ...string) ([]*schemas.Score, error) {
	var scores []*schemas.Score
	err := Preloaded(r.db, preload).
		Where("user_id = ?", userId).
		Where("mode = ?", mode).
		Where("replay_views > 0").
		Where("hidden = ?", false).
		Where("failtime IS NULL").
		Order("replay_views DESC").
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
