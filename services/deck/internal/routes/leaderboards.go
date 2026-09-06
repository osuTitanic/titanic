package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/repositories"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

type LeaderboardType int

const (
	LeaderboardLocal LeaderboardType = iota
	LeaderboardGlobal
	LeaderboardMods
	LeaderboardFriends
	LeaderboardCountry
)

func (t LeaderboardType) Valid() bool {
	return t >= LeaderboardLocal && t <= LeaderboardCountry
}

type LeaderboardResponseType int

const (
	LeaderboardBeatmapNotSubmitted = -1
	LeaderboardBeatmapPending      = 0
	LeaderboardBeatmapNeedsUpdate  = 1
	LeaderboardBeatmapRanked       = 2
	LeaderboardBeatmapApproved     = 3
	LeaderboardBeatmapQualified    = 4
	LeaderboardBeatmapLoved        = 5
)

type LeaderboardResponse struct {
	Type    LeaderboardResponseType
	Beatmap *schemas.Beatmap

	Scores     []*schemas.Score
	ScoreCount int

	PersonalBest      *schemas.Score
	PersonalBestIndex int
}

type LeaderboardRequest struct {
	RequestType    LeaderboardType
	RequestVersion int

	BeatmapNeedsUpdate bool
	IsAuthenticated    bool
	SkipScores         bool

	Beatmap *schemas.Beatmap
	User    *schemas.User
	Mods    *constants.Mods
	Mode    constants.Mode
}

func (r *LeaderboardRequest) FriendsOf() *int {
	if r.User == nil {
		return nil
	}
	if r.RequestType != LeaderboardFriends {
		return nil
	}
	return &r.User.Id
}

func (r *LeaderboardRequest) Country() *string {
	if r.User == nil {
		return nil
	}
	if r.RequestType != LeaderboardCountry {
		return nil
	}
	return &r.User.Country
}

// /web/osu-osz2-getscores.php -> Latest leaderboard endpoint / introduced filtered leaderboards
func GetScoresOsz2(ctx *server.Context) {
	// TODO
}

// /web/osu-getscores6.php -> Last pre-osz2 leaderboard / added beatmap ratings
func GetScores6(ctx *server.Context) {
	// TODO
}

// /web/osu-getscores5.php -> Added offset & display title metadata
func GetScores5(ctx *server.Context) {
	// TODO
}

// /web/osu-getscores4.php -> Added user's personal best to the response
func GetScores4(ctx *server.Context) {
	request, err := NewLeaderboardRequest(ctx)
	if err != nil {
		ctx.RenderText(http.StatusBadRequest, "-1")
		return
	}

	response, err := processLeaderboardRequest(request, ctx)
	if err != nil {
		ctx.RenderText(http.StatusInternalServerError, "-1")
		return
	}
	if response.Type <= LeaderboardBeatmapNeedsUpdate {
		ctx.RenderText(http.StatusOK, fmt.Sprint(response.Type))
		return
	}

	// This endpoint does not support qualified or loved beatmaps
	// Instead, we'll show them as approved beatmaps
	response.Type = min(response.Type, LeaderboardBeatmapApproved)

	if request.SkipScores {
		ctx.RenderText(http.StatusOK, fmt.Sprint(response.Type))
		return
	}

	lines := []string{
		fmt.Sprint(response.Type),
		// Personal best is always the first line after the status code
		// It may be empty if the user has no submitted score on this map
		formatScore(
			response.PersonalBest,
			response.PersonalBestIndex,
			request.RequestVersion,
		),
	}
	for index, score := range response.Scores {
		lines = append(lines, formatScore(score, index+1, request.RequestVersion))
	}

	ctx.RenderText(http.StatusOK, strings.Join(lines, "\n"))
}

// /web/osu-getscores3.php -> Added "approved" status to the response
func GetScores3(ctx *server.Context) {
	request, err := NewLeaderboardRequest(ctx)
	if err != nil {
		ctx.RenderText(http.StatusBadRequest, "-1")
		return
	}

	response, err := processLeaderboardRequest(request, ctx)
	if err != nil {
		ctx.RenderText(http.StatusInternalServerError, "-1")
		return
	}
	if response.Type <= LeaderboardBeatmapNeedsUpdate {
		ctx.RenderText(http.StatusOK, fmt.Sprint(response.Type))
		return
	}

	// This endpoint does not support qualified or loved beatmaps
	// Instead, we'll show them as approved beatmaps
	response.Type = min(response.Type, LeaderboardBeatmapApproved)

	if request.SkipScores {
		ctx.RenderText(http.StatusOK, fmt.Sprint(response.Type))
		return
	}

	formatter := func(score *schemas.Score) string {
		return formatScoreLegacy(score, "|")
	}
	lines := []string{fmt.Sprint(response.Type)}
	lines = append(lines, format(response.Scores, formatter)...)

	ctx.RenderText(http.StatusOK, strings.Join(lines, "\n"))
}

// /web/osu-getscores2.php -> Added beatmap updating logic & pending/ranked status distinction
func GetScores2(ctx *server.Context) {
	request, err := NewLeaderboardRequest(ctx)
	if err != nil {
		ctx.RenderText(http.StatusBadRequest, "-1")
		return
	}

	response, err := processLeaderboardRequest(request, ctx)
	if err != nil {
		ctx.RenderText(http.StatusInternalServerError, "-1")
		return
	}
	if response.Type <= LeaderboardBeatmapNeedsUpdate {
		// Beatmap is pending, needs update, or not submitted
		ctx.RenderText(http.StatusOK, fmt.Sprint(response.Type))
		return
	}

	if response.Type > LeaderboardBeatmapRanked {
		// This endpoint does not support approved, qualified, or loved beatmaps
		response.Type = LeaderboardBeatmapRanked
	}

	// Note that if we respond with any status at all,
	// the client will skip reading the scores, since
	// statuses -1, 0 and 1 are unranked, and 2 is only
	// being responded with when using SkipScores

	if request.SkipScores {
		// Later iterations of this endpoint added the "s" parameter to skip scores
		// This is used in the editor & multiplayer song-select screens,
		// since it doesn't display any scores there
		ctx.RenderText(http.StatusOK, fmt.Sprint(response.Type))
		return
	}

	formatter := func(score *schemas.Score) string {
		return formatScoreLegacy(score, "|")
	}
	ctx.RenderText(http.StatusOK, strings.Join(format(response.Scores, formatter), "\n"))
}

// /web/osu-getscores.php -> The most barebones way to get a leaderboard
func GetScores(ctx *server.Context) {
	request, err := NewLeaderboardRequest(ctx)
	if err != nil {
		ctx.RenderText(http.StatusBadRequest, "-1")
		return
	}

	response, err := processLeaderboardRequest(request, ctx)
	if err != nil {
		ctx.RenderText(http.StatusInternalServerError, "-1")
		return
	}

	// This endpoint does not have any handling of beatmaps with "update available"
	// since the only thing it provides is a checksum of the map
	// It also can't distinguish between pending and ranked maps
	if request.Beatmap == nil {
		ctx.RenderText(http.StatusOK, "-1")
		return
	}
	if response.Type <= LeaderboardBeatmapNotSubmitted {
		ctx.RenderText(http.StatusOK, "-1")
		return
	}

	formatter := func(score *schemas.Score) string {
		return formatScoreLegacy(score, ":")
	}
	ctx.RenderText(http.StatusOK, strings.Join(format(response.Scores, formatter), "\n"))
}

func NewLeaderboardRequest(ctx *server.Context) (*LeaderboardRequest, error) {
	beatmapFilename := ctx.QueryValue("f")
	beatmapChecksum := ctx.QueryValue("c")
	if beatmapChecksum == "" {
		// This is required in every leaderboard endpoint
		return nil, fmt.Errorf("missing beatmap checksum")
	}

	userId, err := ctx.QueryValueIntOptional("u")
	if err != nil {
		// Value was provided but was not valid
		return nil, fmt.Errorf("invalid user id: %w", err)
	}
	username := ctx.QueryValueOptional("us")
	password := ctx.QueryValueOptional("ha")

	leaderboardType, err := ctx.QueryValueEnum[LeaderboardType]("v")
	if err != nil {
		leaderboardType = LeaderboardGlobal
	}

	mode, err := ctx.QueryValueEnum[constants.Mode]("m")
	if err != nil {
		mode = constants.ModeOsu
	}
	mods, err := ctx.QueryValueEnumOptional[constants.Mods]("mods")
	if err != nil {
		mods = nil
	}
	if leaderboardType != LeaderboardMods {
		mods = nil
	}

	requestVersion, err := ctx.QueryValueInt("vv")
	if err != nil {
		requestVersion = 0
	}
	skipScores := ctx.QueryValue("s") == "1"

	request := &LeaderboardRequest{
		RequestVersion:  requestVersion,
		RequestType:     leaderboardType,
		IsAuthenticated: username != nil && password != nil,
		SkipScores:      skipScores,
		Mode:            mode,
		Mods:            mods,
	}

	if username != nil && password != nil {
		// This endpoint requires user authentication for pb's & certain
		// leaderboard types (e.g. friends, country)
		user, err := ctx.AuthenticateUser(*username, *password, true)
		if err != nil {
			return nil, fmt.Errorf("failed to authenticate user: %w", err)
		}
		request.User = user
	}

	if request.User == nil && userId != nil {
		// This endpoint does not require user authentication, but
		// still needs user-related information for pb's
		user, err := ctx.State.Users.ById(*userId)
		if err == nil {
			request.User = user
		}
	}

	beatmap, ok := resolveBeatmap(beatmapChecksum, beatmapFilename, ctx)
	if !ok {
		// Let endpoints handle non-submitted beatmaps
		return request, nil
	}
	request.Beatmap = beatmap
	request.BeatmapNeedsUpdate = beatmap.Checksum != beatmapChecksum
	return request, nil
}

func processLeaderboardRequest(request *LeaderboardRequest, ctx *server.Context) (response LeaderboardResponse, err error) {
	if request.Beatmap == nil {
		return LeaderboardResponse{Type: LeaderboardBeatmapNotSubmitted}, nil
	}
	if request.Beatmap.Status == constants.BeatmapStatusInactive {
		return LeaderboardResponse{Type: LeaderboardBeatmapNotSubmitted}, nil
	}
	if request.BeatmapNeedsUpdate {
		return LeaderboardResponse{Type: LeaderboardBeatmapNeedsUpdate}, nil
	}
	response = LeaderboardResponse{
		Beatmap: request.Beatmap,
	}

	switch request.Beatmap.Status {
	case constants.BeatmapStatusPending, constants.BeatmapStatusWIP, constants.BeatmapStatusGraveyard:
		response.Type = LeaderboardBeatmapPending
	case constants.BeatmapStatusRanked:
		response.Type = LeaderboardBeatmapRanked
	case constants.BeatmapStatusApproved:
		response.Type = LeaderboardBeatmapApproved
	case constants.BeatmapStatusQualified:
		response.Type = LeaderboardBeatmapQualified
	case constants.BeatmapStatusLoved:
		response.Type = LeaderboardBeatmapLoved
	default:
		return response, fmt.Errorf("invalid beatmap status: %d", request.Beatmap.Status)
	}

	if response.Type <= LeaderboardBeatmapPending {
		return response, nil
	}

	filter := repositories.BeatmapLeaderboardFilter{
		BeatmapId: request.Beatmap.Id,
		Mode:      request.Mode,
		Mods:      request.Mods,
		Country:   request.Country(),
		FriendsOf: request.FriendsOf(),
	}

	response.ScoreCount, err = ctx.State.Scores.FetchLeaderboardCount(filter)
	if err != nil {
		return response, fmt.Errorf("fetch leaderboard count: %w", err)
	}
	if request.SkipScores {
		return response, nil
	}

	if request.User != nil {
		response.PersonalBest, err = ctx.State.Scores.FetchLeaderboardPersonalBest(filter, request.User.Id, "User")
		if err != nil {
			return response, fmt.Errorf("fetch personal best: %w", err)
		}
		response.PersonalBestIndex, err = ctx.State.Scores.FetchLeaderboardScoreIndex(filter, response.PersonalBest)
		if err != nil {
			return response, fmt.Errorf("fetch personal best index: %w", err)
		}
	}

	response.Scores, err = ctx.State.Scores.FetchLeaderboardScores(filter, ctx.State.Config.ScoreResponseLimit, "User")
	if err != nil {
		return response, fmt.Errorf("fetch leaderboard scores: %w", err)
	}
	return response, nil
}

func resolveBeatmap(checksum string, filename string, ctx *server.Context) (*schemas.Beatmap, bool) {
	// TODO: maybe return an error for this function
	beatmap, err := ctx.State.Beatmaps.ByChecksum(checksum, "Beatmapset")
	if err != nil {
		return nil, false
	}
	if beatmap != nil {
		return beatmap, true
	}
	if filename == "" {
		return nil, false
	}
	beatmap, err = ctx.State.Beatmaps.ByFilename(filename, "Beatmapset")
	if err != nil {
		return nil, false
	}
	if beatmap != nil {
		return beatmap, true
	}
	return nil, false
}

func formatScore(score *schemas.Score, index int, requestVersion int) string {
	if score == nil {
		return ""
	}
	submittedAt := score.SubmittedAt.Format("2006-01-02 15:04:05")

	if requestVersion >= 2 {
		// This was changed to a unix timestamp in request version 2
		submittedAt = strconv.FormatInt(score.SubmittedAt.Unix(), 10)
	}

	// Filter out the characters we use to separate fields, just to be sure
	replacer := strings.NewReplacer("|", "", "\n", "")

	fields := []string{
		strconv.FormatInt(score.Id, 10),
		replacer.Replace(score.User.Name),
		strconv.FormatInt(score.TotalScore, 10),
		strconv.Itoa(score.MaxCombo),
		strconv.Itoa(score.Count50),
		strconv.Itoa(score.Count100),
		strconv.Itoa(score.Count300),
		strconv.Itoa(score.CountMiss),
		strconv.Itoa(score.CountKatu),
		strconv.Itoa(score.CountGeki),
		strconv.Itoa(integerBoolean(score.Perfect)),
		strconv.Itoa(int(score.Mods)),
		strconv.Itoa(score.UserId),
		strconv.Itoa(index),
		submittedAt,
	}
	if requestVersion >= 4 {
		// "Has Replay", added in request version 4
		fields = append(fields, "1")
	}
	// TODO: NC filter

	return strings.Join(fields, "|")
}

func formatScoreLegacy(score *schemas.Score, separator string) string {
	// Filter out the characters we use to separate fields, just to be sure
	replacer := strings.NewReplacer(separator, "", "\n", "")

	return strings.Join([]string{
		strconv.FormatInt(score.Id, 10),
		replacer.Replace(score.User.Name),
		strconv.FormatInt(score.TotalScore, 10),
		strconv.Itoa(score.MaxCombo),
		strconv.Itoa(score.Count50),
		strconv.Itoa(score.Count100),
		strconv.Itoa(score.Count300),
		strconv.Itoa(score.CountMiss),
		strconv.Itoa(score.CountKatu),
		strconv.Itoa(score.CountGeki),
		strconv.Itoa(integerBoolean(score.Perfect)),
		strconv.Itoa(int(score.Mods)),
		strconv.Itoa(score.UserId),
		score.User.AvatarFilename(),
		score.SubmittedAt.Format("2006-01-02 15:04:05"),
	}, separator)
}

func format[T any, R any](values []T, fn func(T) R) []R {
	result := make([]R, len(values))
	for i, value := range values {
		result[i] = fn(value)
	}
	return result
}
