package routes

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/osuTitanic/titanic-go/internal/constants"
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
	"gorm.io/gorm"
)

func BeatmapRedirect(ctx *server.Context) {
	id := ctx.PathValue("id")
	if id == "" {
		NotFound(ctx)
		return
	}
	ctx.Redirect(http.StatusFound, fmt.Sprintf("/b/%s", id))
}

func BeatmapsetRedirect(ctx *server.Context) {
	id := ctx.PathValue("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		NotFound(ctx)
		return
	}

	beatmapset, err := ctx.State.Repositories.Beatmapsets.ById(idInt, "Beatmaps")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			NotFound(ctx)
			return
		}
		ctx.Logger.Error("Failed to fetch beatmapset", "id", idInt, "error", err)
		InternalServerError(ctx)
		return
	}

	if len(beatmapset.Beatmaps) == 0 {
		NotFound(ctx)
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

func Beatmap(ctx *server.Context) {
	id := ctx.PathValue("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		NotFound(ctx)
		return
	}

	beatmap, err := ctx.State.Repositories.Beatmaps.ById(
		idInt, "Beatmapset", "Beatmapset.Beatmaps", "Beatmapset.CreatorUser",
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			NotFound(ctx)
			return
		}
		ctx.Logger.Error("Failed to fetch beatmap", "id", idInt, "error", err)
		InternalServerError(ctx)
		return
	}

	if beatmap.Status <= constants.BeatmapStatusInactive {
		// Beatmap is inactive / has not been submitted
		NotFound(ctx)
		return
	}

	mode := resolveMode(ctx, beatmap.Mode)
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

	view := templates.BeatmapView{
		DefaultView:      buildDefaultView(ctx),
		Beatmap:          beatmap,
		Beatmapset:       beatmap.Beatmapset,
		Mode:             mode,
		Mods:             modsString,
		Scores:           scores,
		PersonalBest:     personalBest,
		PersonalBestRank: personalBestRank,
		Favourites:       favourites,
		FavouritesCount:  beatmap.Beatmapset.FavouriteCount,
		Favourited:       favourited,
		Collaborations:   collaborations,
		Nominations:      nominations,
		Friends:          friends,
	}
	ctx.RenderTemplate(http.StatusOK, "pages/public/beatmap", view)
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

	friendIds, err := ctx.State.Repositories.Relationships.TargetIdsByStatus(ctx.CurrentUser.Id, 0)
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
	modeQuery := ctx.Request.URL.Query().Get("mode")
	if modeQuery != "" {
		if parsed, err := strconv.Atoi(modeQuery); err == nil {
			mode = constants.Mode(parsed)
		}
	}

	if mode < constants.ModeOsu || mode > constants.ModeMania {
		mode = fallback
	}
	return mode
}

func resolveMods(ctx *server.Context) (*constants.Mods, string) {
	mods := ctx.Request.URL.Query().Get("mods")
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
