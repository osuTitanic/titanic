package routes

import (
	"cmp"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/repositories"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

type directDisplayMode int

const (
	directDisplayModeRanked directDisplayMode = iota
	directDisplayModeRankedStrict
	directDisplayModePending
	directDisplayModeQualified
	directDisplayModeAll
	directDisplayModeGraveyard
	directDisplayModeApproved
	directDisplayModeRankedPlayed
	directDisplayModeLoved
)

const directSearchLimit = 100
const directTimestampLayout = "2006-01-02 15:04:05.999999"

// /web/osu-search.php -> Search for beatmapsets through osu! direct
func DirectSearch(ctx *server.Context) {
	user, err := authenticateDirectUser(ctx)
	if err != nil {
		ctx.Logger.Error("Failed to authenticate user", "error", err)
		ctx.RenderText(http.StatusOK, directError("A server error occurred. Please try again!"))
		return
	}
	if user == nil && !ctx.State.Config.AllowUnauthenticatedDirect {
		ctx.RenderText(http.StatusOK, directError("This version of osu! is not supported."))
		return
	}

	var userId *int
	if user != nil {
		userId = new(user.Id)
	}

	query := ctx.QueryValue("q")
	options, supportsPageOffset, err := buildDirectSearchOptions(ctx, userId)
	if err != nil {
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return
	}

	beatmapsets, err := ctx.State.Beatmapsets.Search(options, "Beatmaps")
	if err != nil {
		ctx.Logger.Error("Failed to search beatmapsets", "error", err)
		ctx.RenderText(http.StatusOK, directError("A server error occurred. Please try again!"))
		return
	}

	// another go moment. we don't have sets. lmao.
	// TODO: add some kind of "Collect" method but with a filter function
	topicIds := make(map[int]struct{}, len(beatmapsets))
	for _, beatmapset := range beatmapsets {
		if beatmapset.TopicId != nil {
			topicIds[*beatmapset.TopicId] = struct{}{}
		}
	}

	// TODO: I want initial post IDs to be part of the beatmapset schema itself
	//		 one day, automatically updated by a postgres trigger
	initialPostIds, err := ctx.State.ForumPosts.FetchInitialIdsByTopicIds(
		slices.Collect(maps.Keys(topicIds)),
	)
	if err != nil {
		ctx.Logger.Error("Failed to resolve forum posts", "error", err)
		ctx.RenderText(http.StatusOK, directError("A server error occurred. Please try again!"))
		return
	}

	ctx.Logger.Info(
		"Processed osu!direct search request",
		"user_id", userId,
		"query", query,
		"results", len(beatmapsets),
	)
	ctx.RenderText(
		http.StatusOK,
		formatDirectSearchResponse(beatmapsets, initialPostIds, supportsPageOffset),
	)
}

func authenticateDirectUser(ctx *server.Context) (*schemas.User, error) {
	username := ctx.QueryValue("u")
	password := cmp.Or(ctx.QueryValue("h"), ctx.QueryValue("c"))
	if username == "" || password == "" {
		return nil, nil
	}

	user, err := ctx.AuthenticateUser(username, password, false)
	if errors.Is(err, server.ErrUserNotFound) || errors.Is(err, server.ErrInvalidPassword) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// TODO: Add permission check
	return user, nil
}

func buildDirectSearchOptions(ctx *server.Context, userId *int) (repositories.BeatmapsetSearchOptions, bool, error) {
	options := repositories.BeatmapsetSearchOptions{
		QueryString: ctx.QueryValue("q"),
		Category:    constants.BeatmapCategoryAny,
		Sort:        constants.BeatmapSortRelevance,
		UserId:      userId,
		Limit:       directSearchLimit,
	}

	displayMode := directDisplayModeAll
	if value := ctx.QueryValue("r"); value != "" {
		displayMode, _ = ctx.QueryValueEnum[directDisplayMode]("r")
	}

	// The display mode will determine which statuses are included
	// in the search results, e.g. ranked, pending, graveyard, ...
	applyDirectDisplayMode(&options, displayMode)

	if value := ctx.QueryValue("m"); value != "" {
		mode, _ := ctx.QueryValueEnum[constants.Mode]("m")
		if mode >= 0 {
			options.Mode = new(mode)
		}
	}

	supportsPageOffset := ctx.Request.URL.Query().Has("p")
	if supportsPageOffset {
		page, _ := ctx.QueryValueInt("p")
		options.Offset = page * directSearchLimit
	}

	// osu! direct doesn't have a separate query value that
	// specifies the sort order. It sets it through the query
	// string itself: "Newest", "Most Played" & "Top Rated"
	// Kinda dumb if you ask me, but hey, peppy code™
	switch options.QueryString {
	case "Newest":
		options.QueryString = ""
		switch displayMode {
		case directDisplayModePending, directDisplayModeAll, directDisplayModeGraveyard:
			// Makes more sense to update them by last update time rather
			// than ranked date, since they aren't ranked, duh!
			options.Sort = constants.BeatmapSortUpdated
		default:
			options.Sort = constants.BeatmapSortRanked
		}
	case "Most Played":
		options.QueryString = ""
		options.Sort = constants.BeatmapSortPlays
	case "Top Rated":
		options.QueryString = ""
		options.Sort = constants.BeatmapSortRating
	}

	// Search options are normalized by BeatmapsetSearchOptions.Normalize
	// so we don't need to validate them here
	return options, supportsPageOffset, nil
}

func applyDirectDisplayMode(options *repositories.BeatmapsetSearchOptions, displayMode directDisplayMode) {
	// Related:
	// https://github.com/osuTitanic/deck/pull/441
	// https://github.com/osuTitanic/common/pull/26
	switch displayMode {
	case directDisplayModeRanked:
		options.Statuses = []constants.BeatmapStatus{
			constants.BeatmapStatusRanked,
			constants.BeatmapStatusApproved,
		}
	case directDisplayModeRankedStrict:
		options.Statuses = []constants.BeatmapStatus{constants.BeatmapStatusRanked}
	case directDisplayModePending:
		options.Statuses = []constants.BeatmapStatus{
			constants.BeatmapStatusWIP,
			constants.BeatmapStatusPending,
		}
	case directDisplayModeQualified:
		options.Statuses = []constants.BeatmapStatus{constants.BeatmapStatusQualified}
	case directDisplayModeAll:
		options.Category = constants.BeatmapCategoryAny
	case directDisplayModeGraveyard:
		options.Statuses = []constants.BeatmapStatus{constants.BeatmapStatusGraveyard}
	case directDisplayModeApproved:
		options.Statuses = []constants.BeatmapStatus{constants.BeatmapStatusApproved}
	case directDisplayModeRankedPlayed:
		options.Statuses = []constants.BeatmapStatus{
			constants.BeatmapStatusRanked,
			constants.BeatmapStatusApproved,
		}
		if options.UserId == nil {
			options.UserId = new(0)
		}
		options.Played = true
	case directDisplayModeLoved:
		options.Statuses = []constants.BeatmapStatus{constants.BeatmapStatusLoved}
	}
}

func formatDirectSearchResponse(beatmapsets []*schemas.Beatmapset, initialPostIds map[int]int64, supportsPageOffset bool) string {
	resultCount := len(beatmapsets)
	if supportsPageOffset && resultCount >= directSearchLimit {
		resultCount++
	}

	var response strings.Builder
	response.WriteString(strconv.Itoa(resultCount))
	for _, beatmapset := range beatmapsets {
		// TODO: We eventually want the beatmapset itself to have the post ID
		//		 We could auto-update it with postgres triggers and stuff yuh
		postId := int64(0)
		if beatmapset.TopicId != nil {
			postId = initialPostIds[*beatmapset.TopicId]
		}
		response.WriteByte('\n')
		response.WriteString(formatDirectBeatmap(beatmapset, postId))
	}
	return response.String()
}

func formatDirectBeatmap(beatmapset *schemas.Beatmapset, postId int64) string {
	artist := pointerToString(beatmapset.Artist)
	title := pointerToString(beatmapset.Title)
	creator := pointerToString(beatmapset.Creator)
	filename := sanitizeDirectFilename(beatmapset.OszFilename())

	versions := make([]string, len(beatmapset.Beatmaps))
	for index, beatmap := range beatmapset.Beatmaps {
		versions[index] = fmt.Sprintf("%s@%d", beatmap.Version, beatmap.Mode)
	}

	topicId := 0
	if beatmapset.TopicId != nil {
		topicId = *beatmapset.TopicId
	}

	return strings.Join([]string{
		filename,
		artist,
		title,
		creator,
		strconv.Itoa(int(beatmapset.Status)),
		strconv.FormatFloat(beatmapset.RatingAverage, 'f', -1, 64),
		beatmapset.LastUpdate.Format(directTimestampLayout),
		strconv.Itoa(beatmapset.Id),
		strconv.Itoa(topicId),
		strconv.Itoa(integerBoolean(beatmapset.HasVideo)),
		strconv.Itoa(integerBoolean(beatmapset.HasStoryboard)),
		strconv.Itoa(beatmapset.OszFilesize),
		strconv.Itoa(beatmapset.OszFilesizeNovideo),
		strings.Join(versions, ","),
		strconv.FormatInt(postId, 10),
	}, "|")
}

func directError(message string) string {
	return "-1\n" + message
}

// having no standardized way of doing these sucks lol

func sanitizeDirectFilename(filename string) string {
	return strings.Map(func(character rune) rune {
		if character < 0x20 || strings.ContainsRune(`<>:"/\|?*`, character) {
			return -1
		}
		return character
	}, filename)
}

func pointerToString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func integerBoolean(value bool) int {
	if value {
		return 1
	}
	return 0
}

// TODO: /web/osu-search-set.php
