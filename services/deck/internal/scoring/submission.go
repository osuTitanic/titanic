package scoring

import (
	"time"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
)

type Result struct {
	Type       ResultType
	Warnings   []error
	Submission *SubmissionContext
}

type ResultType uint8

const (
	ResultAccepted ResultType = iota
	ResultUserNotFound
	ResultInvalidPassword
	ResultInactive
	ResultBanned
	ResultRejected
	ResultBeatmapUnavailable
)

func (resultType ResultType) String() string {
	switch resultType {
	case ResultAccepted:
		return ""
	case ResultUserNotFound:
		return "error: nouser"
	case ResultInvalidPassword:
		return "error: pass"
	case ResultInactive:
		return "error: inactive"
	case ResultBanned:
		return "error: ban"
	case ResultBeatmapUnavailable:
		return "error: beatmap"
	case ResultRejected:
		return "error: no"
	default:
		return "error: no"
	}
}

type Endpoint uint8

const (
	EndpointLegacy Endpoint = iota
	EndpointModular
	EndpointModularSelector
)

func (endpoint Endpoint) UsesEncryptedFormScore() bool {
	return endpoint == EndpointModular || endpoint == EndpointModularSelector
}

func (endpoint Endpoint) IncludesBeatmapRankingChart() bool {
	return endpoint == EndpointModularSelector
}

func (endpoint Endpoint) UsesLegacyResponse() bool {
	return endpoint == EndpointLegacy
}

type SubmissionContext struct {
	*schemas.Score

	Endpoint Endpoint

	Username        string
	BeatmapChecksum string
	Replay          []byte
	FunSpoiler      string
	ClientHash      *string
	Processes       *string
	Flags           constants.IntegrityFlags
	Exited          bool
	ClientPassed    bool

	PersonalBestPP    *schemas.Score
	PersonalBestScore *schemas.Score

	CurrentStats  *schemas.Stats
	PreviousStats *schemas.Stats

	OldBeatmapRank int
	NewBeatmapRank int

	Charts       SubmissionCharts
	Achievements []*schemas.Achievement
}

func NewSubmissionContext(endpoint Endpoint, score *schemas.Score) *SubmissionContext {
	if score == nil {
		score = new(schemas.Score)
	}

	score.SubmittedAt = time.Now().UTC()
	score.StatusPP = constants.ScoreStatusSubmitted
	score.StatusScore = constants.ScoreStatusSubmitted

	return &SubmissionContext{
		Score:    score,
		Endpoint: endpoint,
		Charts:   newSubmissionCharts(),
	}
}
