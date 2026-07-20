package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/services/stern/internal/server"
	"github.com/osuTitanic/titanic/services/stern/internal/templates"
)

func BeatmapRedirect(ctx *server.Context) {
	id := ctx.PathValue("id")
	if id == "" {
		BeatmapNotFound(ctx)
		return
	}
	ctx.Redirect(http.StatusFound, fmt.Sprintf("/b/%s", id))
}

func BeatmapsetRedirect(ctx *server.Context) {
	id, err := ctx.PathValueInt("id")
	if err != nil {
		BeatmapNotFound(ctx)
		return
	}

	beatmapset, err := ctx.State.Repositories.Beatmapsets.ById(id, "Beatmaps")
	if err != nil {
		ctx.Logger.Error("Failed to fetch beatmapset", "id", id, "error", err)
		InternalServerError(ctx)
		return
	}
	if beatmapset == nil || len(beatmapset.Beatmaps) == 0 {
		BeatmapNotFound(ctx)
		return
	}

	// Default to the requested mode, or osu! standard when none is given
	mode := resolveMode(ctx, constants.ModeOsu)

	// Resolve beatmap that matches the mode
	beatmap := beatmapset.Beatmaps[0]
	for _, candidate := range beatmapset.Beatmaps {
		if candidate.Mode == mode {
			beatmap = candidate
			break
		}
	}

	location := fmt.Sprintf("/b/%d", beatmap.Id)
	if ctx.Request.URL.RawQuery != "" {
		location += "?" + ctx.Request.URL.RawQuery
	}
	ctx.Redirect(http.StatusMovedPermanently, location)
}

func RedirectToBeatmapset(ctx *server.Context) {
	id := ctx.PathValue("id")
	if id == "" {
		BeatmapNotFound(ctx)
		return
	}
	ctx.Redirect(http.StatusFound, fmt.Sprintf("/s/%s", id))
}

func RedirectToDiscussion(ctx *server.Context) {
	id, err := ctx.PathValueInt("setId")
	if err != nil {
		BeatmapNotFound(ctx)
		return
	}

	beatmapset, err := ctx.State.Repositories.Beatmapsets.ById(id)
	if err != nil {
		ctx.Logger.Error("Failed to fetch beatmapset", "id", id, "error", err)
		InternalServerError(ctx)
		return
	}
	if beatmapset == nil {
		BeatmapNotFound(ctx)
		return
	}

	// Beatmapsets without a forum topic fall back to the set page
	if beatmapset.TopicId == nil {
		ctx.Redirect(http.StatusFound, fmt.Sprintf("/s/%d", beatmapset.Id))
		return
	}

	ctx.Redirect(http.StatusFound, fmt.Sprintf("/forum/t/%d", *beatmapset.TopicId))
}

func Beatmap(ctx *server.Context) {
	id, err := ctx.PathValueInt("id")
	if err != nil {
		BeatmapNotFound(ctx)
		return
	}

	beatmap, err := ctx.State.Repositories.Beatmaps.ById(
		id, "Beatmapset", "Beatmapset.Beatmaps", "Beatmapset.CreatorUser",
	)
	if err != nil {
		ctx.Logger.Error("Failed to fetch beatmap", "id", id, "error", err)
		InternalServerError(ctx)
		return
	}
	if beatmap == nil || beatmap.Status <= constants.BeatmapStatusInactive {
		// Beatmap not found / inactive / has not been submitted
		BeatmapNotFound(ctx)
		return
	}

	mode := resolveBeatmapMode(ctx, beatmap.Mode)
	mods, modsString := resolveMods(ctx)

	personalBest, personalBestRank, err := fetchUserScore(ctx, beatmap, mode)
	if err != nil {
		ctx.Logger.Error("Failed to fetch user pb", "beatmap", beatmap.Id, "error", err)
		InternalServerError(ctx)
		return
	}

	friends, err := fetchUserFriends(ctx)
	if err != nil {
		ctx.Logger.Error("Failed to fetch friend IDs", "user", ctx.CurrentUser.Id, "error", err)
		InternalServerError(ctx)
		return
	}

	scores, err := fetchBeatmapScores(ctx, beatmap.Id, mode, mods)
	if err != nil {
		ctx.Logger.Error("Failed to fetch beatmap scores", "beatmap", beatmap.Id, "error", err)
		InternalServerError(ctx)
		return
	}

	favourites, err := ctx.State.Repositories.Favourites.ManyBySetId(beatmap.SetId, 5, "User")
	if err != nil {
		ctx.Logger.Error("Failed to fetch favourites", "set", beatmap.SetId, "error", err)
		InternalServerError(ctx)
		return
	}

	favourited := false
	if ctx.CurrentUser != nil {
		favourited, _ = ctx.State.Repositories.Favourites.ExistsForUser(
			ctx.CurrentUser.Id,
			beatmap.SetId,
		)
	}

	collaborations, err := ctx.State.Repositories.Collaborations.FetchByBeatmap(beatmap.Id, "User")
	if err != nil {
		ctx.Logger.Error("Failed to fetch collaborations", "beatmap", beatmap.Id, "error", err)
		InternalServerError(ctx)
		return
	}

	nominations, err := ctx.State.Repositories.Nominations.FetchBySet(beatmap.SetId, "User")
	if err != nil {
		ctx.Logger.Error("Failed to fetch nominations", "set", beatmap.SetId, "error", err)
		InternalServerError(ctx)
		return
	}

	collaborationRequests, err := fetchCollaborationRequests(ctx, beatmap)
	if err != nil {
		ctx.Logger.Error("Failed to fetch collaboration requests", "beatmap", beatmap.Id, "error", err)
		InternalServerError(ctx)
		return
	}

	view := templates.BeatmapView{
		DefaultView:           buildDefaultViewWithPermissions(ctx),
		Beatmap:               beatmap,
		Beatmapset:            beatmap.Beatmapset,
		Mode:                  mode,
		Mods:                  modsString,
		Scores:                scores,
		PersonalBest:          personalBest,
		PersonalBestRank:      personalBestRank,
		Favourites:            favourites,
		FavouritesCount:       beatmap.Beatmapset.FavouriteCount,
		Favourited:            favourited,
		Collaborations:        collaborations,
		Nominations:           nominations,
		Friends:               friends,
		CollaborationRequests: collaborationRequests,
		IsBeatmapAuthor:       isBeatmapAuthor(ctx, beatmap.Beatmapset),
		Invite:                findInvite(ctx, collaborationRequests),
		BatNominated:          hasNominated(ctx, nominations),
	}
	ctx.RenderTemplate(http.StatusOK, "pages/public/beatmap", view)
}

func fetchCollaborationRequests(ctx *server.Context, beatmap *schemas.Beatmap) ([]*schemas.BeatmapCollaborationRequest, error) {
	if ctx.CurrentUser == nil || beatmap.Status > constants.BeatmapStatusPending {
		return nil, nil
	}
	return ctx.State.Repositories.Collaborations.FetchRequestsByBeatmap(beatmap.Id, "User", "Target")
}

func findInvite(ctx *server.Context, requests []*schemas.BeatmapCollaborationRequest) *schemas.BeatmapCollaborationRequest {
	if ctx.CurrentUser == nil {
		return nil
	}
	for _, request := range requests {
		if request.TargetId == ctx.CurrentUser.Id {
			return request
		}
	}
	return nil
}

func hasNominated(ctx *server.Context, nominations []*schemas.BeatmapNomination) bool {
	if ctx.CurrentUser == nil {
		return false
	}
	for _, nomination := range nominations {
		if nomination.UserId == ctx.CurrentUser.Id {
			return true
		}
	}
	return false
}

func isBeatmapAuthor(ctx *server.Context, beatmapset *schemas.Beatmapset) bool {
	if ctx.CurrentUser == nil || beatmapset.CreatorId == nil {
		return false
	}
	return beatmapset.Server != constants.BeatmapServerBancho &&
		*beatmapset.CreatorId == ctx.CurrentUser.Id
}

func fetchBeatmapScores(ctx *server.Context, beatmapId int, mode constants.Mode, mods *constants.Mods) ([]*schemas.Score, error) {
	limit := ctx.State.Config.ScoreResponseLimit

	if mods != nil {
		return ctx.State.Repositories.Scores.FetchRangeScoresMods(
			beatmapId, mode, *mods, limit, 0, "User",
		)
	}

	return ctx.State.Repositories.Scores.FetchRangeScores(
		beatmapId, mode, limit, 0, "User",
	)
}

func fetchUserScore(ctx *server.Context, beatmap *schemas.Beatmap, mode constants.Mode) (personalBest *schemas.Score, personalBestRank int, err error) {
	personalBest = new(schemas.Score)
	personalBestRank = 0

	if ctx.CurrentUser == nil {
		return personalBest, personalBestRank, nil
	}

	personalBest, err = ctx.State.Repositories.Scores.FetchPersonalBest(
		beatmap.Id, ctx.CurrentUser.Id, mode, "User",
	)
	if err != nil {
		return nil, 0, err
	}
	if personalBest == nil {
		return nil, 0, nil
	}

	personalBestRank, err = ctx.State.Repositories.Scores.FetchScoreIndexById(
		personalBest.Id, beatmap.Id, mode,
	)
	if err != nil {
		return nil, 0, err
	}

	return personalBest, personalBestRank, nil
}

func fetchUserFriends(ctx *server.Context) (friends map[int]bool, err error) {
	friends = make(map[int]bool)
	if ctx.CurrentUser == nil {
		return friends, nil
	}

	friendIds, err := ctx.State.Repositories.Relationships.TargetIdsByStatus(
		ctx.CurrentUser.Id,
		constants.RelationshipStatusFriend,
	)
	if err != nil {
		ctx.Logger.Error("Failed to fetch friends", "user", ctx.CurrentUser.Id, "error", err)
		InternalServerError(ctx)
		return
	}
	for _, friendId := range friendIds {
		friends[friendId] = true
	}

	return friends, nil
}

func resolveMode(ctx *server.Context, fallback constants.Mode) constants.Mode {
	mode := fallback
	if parsed, err := ctx.QueryValueInt("mode"); err == nil {
		mode = constants.Mode(parsed)
	}

	if mode < constants.ModeOsu || mode > constants.ModeMania {
		mode = fallback
	}
	return mode
}

func resolveBeatmapMode(ctx *server.Context, beatmapMode constants.Mode) constants.Mode {
	fallback := beatmapMode
	if beatmapMode != constants.ModeOsu {
		// Beatmap doesn't support converts so we can't fall back to the user's preferred mode
		return resolveMode(ctx, fallback)
	}
	if ctx.IsAuthenticated() {
		// If the `mode` parameter is not given, fall back to the user's preferred mode
		fallback = ctx.CurrentUser.PreferredMode
	}
	return resolveMode(ctx, fallback)
}

func resolveMods(ctx *server.Context) (*constants.Mods, string) {
	mods := ctx.QueryValue("mods")
	mods = strings.TrimPrefix(strings.TrimSpace(mods), "+")
	if mods == "" {
		return nil, ""
	}

	if parsed, err := strconv.ParseUint(mods, 10, 32); err == nil {
		result := constants.Mods(parsed)
		return &result, mods
	}

	result := constants.ModsFromString(mods)
	return &result, mods
}
