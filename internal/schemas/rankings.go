package schemas

import (
	"time"

	"github.com/osuTitanic/titanic/internal/constants"
)

type Achievement struct {
	UserId     int       `gorm:"column:user_id;primaryKey"`
	Name       string    `gorm:"column:name;primaryKey"`
	Category   string    `gorm:"column:category"` // TODO: Add constant for categories
	Filename   string    `gorm:"column:filename"`
	UnlockedAt time.Time `gorm:"column:unlocked_at;autoCreateTime"`

	User *User `gorm:"foreignKey:UserId;references:Id"`
}

func (Achievement) TableName() string {
	return "achievements"
}

type Score struct {
	Id            int64                 `gorm:"column:id;primaryKey;autoIncrement"`
	UserId        int                   `gorm:"column:user_id"`
	BeatmapId     int                   `gorm:"column:beatmap_id"`
	ClientVersion int                   `gorm:"column:client_version"`
	ClientString  string                `gorm:"column:client_version_string;default:"`
	Checksum      string                `gorm:"column:score_checksum"`
	Mode          constants.Mode        `gorm:"column:mode"`
	PP            float64               `gorm:"column:pp"`
	PPv1          float64               `gorm:"column:ppv1"`
	Acc           float64               `gorm:"column:acc"`
	TotalScore    int64                 `gorm:"column:total_score"`
	MaxCombo      int                   `gorm:"column:max_combo"`
	Mods          constants.Mods        `gorm:"column:mods"`
	Perfect       bool                  `gorm:"column:perfect"`
	Count300      int                   `gorm:"column:n300"`
	Count100      int                   `gorm:"column:n100"`
	Count50       int                   `gorm:"column:n50"`
	CountMiss     int                   `gorm:"column:nmiss"`
	CountGeki     int                   `gorm:"column:ngeki"`
	CountKatu     int                   `gorm:"column:nkatu"`
	Grade         constants.Grade       `gorm:"column:grade;type:varchar(2);default:N"`
	StatusPP      constants.ScoreStatus `gorm:"column:status;default:-1"`
	StatusScore   constants.ScoreStatus `gorm:"column:status_score;default:-1"`
	Pinned        bool                  `gorm:"column:pinned;default:false"`
	Hidden        bool                  `gorm:"column:hidden;default:false"`
	SubmittedAt   time.Time             `gorm:"column:submitted_at;autoCreateTime"`
	Failtime      *int                  `gorm:"column:failtime"`
	ReplayMd5     *string               `gorm:"column:replay_md5"`
	ReplayViews   int                   `gorm:"column:replay_views;default:0"`

	User    *User    `gorm:"foreignKey:UserId;references:Id"`
	Beatmap *Beatmap `gorm:"foreignKey:BeatmapId;references:Id"`
}

func (Score) TableName() string {
	return "scores"
}

func (score *Score) Passed() bool {
	return score.Failtime == nil
}

func (score *Score) Relaxing() bool {
	return score.Mods.Has(constants.Relax) || score.Mods.Has(constants.Autopilot)
}

func (score *Score) RequiresPPv1Update() bool {
	if score.PPv1 <= 0 {
		return true
	}
	timeSinceSubmission := time.Since(score.SubmittedAt)
	// TODO: Add column that determines the last ppv1 update time
	// 		 For now we'll use the submission time

	// Every 10 days: the score loses ~1% of its pp
	// Every 24 hours: the score loses ~0.1% of its pp
	if score.StatusPP >= constants.ScoreStatusBest {
		// For personal best's we want to update scores every 24 hours
		return timeSinceSubmission > 24*time.Hour
	} else {
		// For everything else we can update it every 10 days
		return timeSinceSubmission > 24*time.Hour*10
	}
}

// TotalHits returns the number of successful note hits in the score.
func (score *Score) TotalHits() int {
	if score == nil {
		return 0
	}
	// For osu! & catch modes
	total := score.Count50 + score.Count100 + score.Count300

	// For mania & taiko modes
	if score.Mode == constants.ModeMania || score.Mode == constants.ModeTaiko {
		total += score.CountGeki + score.CountKatu
	}
	return total
}

// TotalObjects returns the number of passed objects used to calculate accuracy.
func (score *Score) TotalObjects() int {
	if score == nil {
		return 0
	}
	// Base total, for taiko & osu!
	total := score.Count50 + score.Count100 + score.Count300 + score.CountMiss

	switch score.Mode {
	case constants.ModeOsu, constants.ModeTaiko:
		return total
	case constants.ModeCatch:
		return total + score.CountKatu
	default:
		return total + score.CountGeki + score.CountKatu
	}
}

// Accuracy returns the score's normalized accuracy in the range [0, 1].
func (score *Score) Accuracy() float64 {
	if score == nil {
		return 0
	}

	totalObjects := float64(score.TotalObjects())
	if totalObjects == 0 {
		return 0
	}

	count300 := float64(score.Count300)
	count100 := float64(score.Count100)
	count50 := float64(score.Count50)
	countGeki := float64(score.CountGeki)
	countKatu := float64(score.CountKatu)

	switch score.Mode {
	case constants.ModeOsu:
		return (300*count300 + 100*count100 + 50*count50) / (300 * totalObjects)
	case constants.ModeTaiko:
		return (count300 + 0.5*count100) / totalObjects
	case constants.ModeCatch:
		return (count300 + count100 + count50) / totalObjects
	case constants.ModeMania:
		return (300*(countGeki+count300) + 200*countKatu + 100*count100 + 50*count50) / (300 * totalObjects)
	default:
		return 0
	}
}

type RankHistory struct {
	UserId      int            `gorm:"column:user_id;primaryKey"`
	Time        time.Time      `gorm:"column:time;primaryKey;autoCreateTime"`
	Mode        constants.Mode `gorm:"column:mode"`
	Rscore      int64          `gorm:"column:rscore"`
	PP          int            `gorm:"column:pp"`
	PPv1        int            `gorm:"column:ppv1"`
	GlobalRank  int            `gorm:"column:global_rank"`
	CountryRank int            `gorm:"column:country_rank"`
	ScoreRank   int            `gorm:"column:score_rank"`
	PPv1Rank    int            `gorm:"column:ppv1_rank"`

	User *User `gorm:"foreignKey:UserId;references:Id"`
}

func (RankHistory) TableName() string {
	return "profile_rank_history"
}

type PlayHistory struct {
	UserId    int            `gorm:"column:user_id;primaryKey"`
	Mode      constants.Mode `gorm:"column:mode;primaryKey"`
	Year      int            `gorm:"column:year;primaryKey"`
	Month     int            `gorm:"column:month;primaryKey"`
	Plays     int            `gorm:"column:plays;default:0"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`

	User *User `gorm:"foreignKey:UserId;references:Id"`
}

func (PlayHistory) TableName() string {
	return "profile_play_history"
}

func (p *PlayHistory) Date() time.Time {
	return time.Date(p.Year, time.Month(p.Month), 1, 0, 0, 0, 0, time.UTC)
}

type ReplayHistory struct {
	UserId      int            `gorm:"column:user_id;primaryKey"`
	Mode        constants.Mode `gorm:"column:mode;primaryKey"`
	Year        int            `gorm:"column:year;primaryKey"`
	Month       int            `gorm:"column:month;primaryKey"`
	ReplayViews int            `gorm:"column:replay_views;default:0"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime"`

	User *User `gorm:"foreignKey:UserId;references:Id"`
}

func (ReplayHistory) TableName() string {
	return "profile_replay_history"
}

func (r *ReplayHistory) Date() time.Time {
	return time.Date(r.Year, time.Month(r.Month), 1, 0, 0, 0, 0, time.UTC)
}
