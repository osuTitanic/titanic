package scoring

import (
	"net/http"
	"time"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
)

type Result struct {
	Type       ResultType
	Warnings   []error
	Submission *SubmissionContext
}

func (r *Result) Rejected() bool {
	return r.Type != ResultAccepted
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
	ResultBanchoUnavailable // Will respond with a 503 error so client can retry
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

func (resultType ResultType) StatusCode() int {
	// Used for /web/osu-submit.php which has
	// no error response strings
	switch resultType {
	case ResultUserNotFound,
		ResultInvalidPassword,
		ResultInactive,
		ResultBanned:
		return http.StatusUnauthorized
	case ResultBeatmapUnavailable:
		return http.StatusNotFound
	case ResultBanchoUnavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusBadRequest
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

func (endpoint Endpoint) FormatResponse(result *SubmissionContext) string {
	return FormatSubmissionResponse(endpoint, result)
}

func (endpoint Endpoint) FormatError(result ResultType) (int, string) {
	return FormatSubmissionError(endpoint, result)
}

func (endpoint Endpoint) RequestValue(request *http.Request, name string) string {
	if endpoint.UsesLegacyResponse() {
		return request.URL.Query().Get(name)
	} else {
		return request.PostFormValue(name)
	}
}

type SubmissionContext struct {
	*schemas.Score

	Endpoint Endpoint

	BeatmapChecksum string
	Username        string
	Replay          []byte
	FunSpoiler      string
	ClientHash      string
	Processes       string
	Flags           constants.IntegrityFlags
	Exited          bool
	Passed          bool

	PersonalBestPP    *schemas.Score
	PersonalBestScore *schemas.Score

	CurrentStats  *schemas.Stats
	PreviousStats *schemas.Stats

	OldBeatmapRank int
	NewBeatmapRank int

	Charts       *SubmissionCharts
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
		Charts:   NewSubmissionCharts(),
	}
}
