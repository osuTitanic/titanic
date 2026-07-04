package schemas

import (
	"encoding/json"
	"time"

	"github.com/osuTitanic/titanic-go/internal/constants"
)

type Match struct {
	Id        int        `gorm:"column:id;primaryKey;autoIncrement"`
	BanchoId  int        `gorm:"column:bancho_id"`
	Name      string     `gorm:"column:name"`
	CreatorId int        `gorm:"column:creator_id"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	EndedAt   *time.Time `gorm:"column:ended_at"`

	Creator *User         `gorm:"foreignKey:CreatorId;references:Id"`
	Events  []*MatchEvent `gorm:"foreignKey:MatchId;references:Id"`
}

func (Match) TableName() string {
	return "mp_matches"
}

type MatchEvent struct {
	MatchId int                      `gorm:"column:match_id;primaryKey"`
	Time    time.Time                `gorm:"column:time;primaryKey;autoCreateTime"`
	Type    constants.MatchEventType `gorm:"column:type"`
	Data    json.RawMessage          `gorm:"column:data;type:jsonb"`

	Match *Match `gorm:"foreignKey:MatchId;references:Id"`
}

func (MatchEvent) TableName() string {
	return "mp_events"
}

/* Views of the "data" column in mp_events, as written by bancho */

type MatchEventUserData struct {
	UserId int    `json:"user_id"`
	Name   string `json:"name"`
}

type MatchEventUserRef struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type MatchEventHostData struct {
	New      MatchEventUserRef `json:"new"`
	Previous MatchEventUserRef `json:"previous"`
}

type MatchEventGameData struct {
	BeatmapId   int                        `json:"beatmap_id"`
	BeatmapText string                     `json:"beatmap_text"`
	Mode        constants.Mode             `json:"mode"`
	TeamMode    constants.MatchTeamType    `json:"team_mode"`
	ScoringMode constants.MatchScoringType `json:"scoring_mode"`
	Mods        constants.Mods             `json:"mods"`
	StartTime   string                     `json:"start_time"`
	EndTime     string                     `json:"end_time"`
	Results     []MatchEventResultData     `json:"results"`
}

type MatchEventResultData struct {
	Place  int                    `json:"place"`
	Player MatchEventResultPlayer `json:"player"`
	Score  MatchEventResultScore  `json:"score"`
}

type MatchEventResultPlayer struct {
	Id      int                `json:"id"`
	Name    string             `json:"name"`
	Country string             `json:"country"`
	Team    constants.SlotTeam `json:"team"`
	Mods    constants.Mods     `json:"mods"`
}

type MatchEventResultScore struct {
	Total     int     `json:"score"`
	Count300  int     `json:"c300"`
	Count100  int     `json:"c100"`
	Count50   int     `json:"c50"`
	CountMiss int     `json:"cMiss"`
	Accuracy  float64 `json:"accuracy"`
	MaxCombo  int     `json:"max_combo"`
	Failed    bool    `json:"failed"`
}

func (e *MatchEvent) UserData() (*MatchEventUserData, error) {
	var data MatchEventUserData
	err := json.Unmarshal(e.Data, &data)
	return &data, err
}

func (e *MatchEvent) HostData() (*MatchEventHostData, error) {
	var data MatchEventHostData
	err := json.Unmarshal(e.Data, &data)
	return &data, err
}

func (e *MatchEvent) GameData() (*MatchEventGameData, error) {
	var data MatchEventGameData
	err := json.Unmarshal(e.Data, &data)
	return &data, err
}

// Duration returns the duration of the game, or 0 if it cannot be resolved
func (d *MatchEventGameData) Duration() time.Duration {
	const layout = "2006-01-02 15:04:05"

	start, err := time.Parse(layout, d.StartTime)
	if err != nil {
		return 0
	}
	end, err := time.Parse(layout, d.EndTime)
	if err != nil {
		return 0
	}

	if end.Before(start) {
		return 0
	}
	return end.Sub(start)
}
