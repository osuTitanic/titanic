package routes

import (
	"net/http"
	"strconv"

	"github.com/osuTitanic/titanic-go/internal/constants"
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
)

func Score(ctx *server.Context) {
	id := ctx.PathValue("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		NotFound(ctx)
		return
	}

	score, err := ctx.State.Repositories.Scores.ById(
		idInt, "User", "User.Stats", "Beatmap", "Beatmap.Beatmapset",
	)
	if err != nil {
		ctx.Logger.Error("Failed to fetch score", "id", idInt, "error", err)
		InternalServerError(ctx)
		return
	}

	// Only passed, non-hidden scores have a public page
	if score == nil || score.Hidden || !score.Passed() {
		NotFound(ctx)
		return
	}

	if score.Beatmap == nil || score.Beatmap.Beatmapset == nil || score.User == nil {
		NotFound(ctx)
		return
	}

	scoreRank, err := ctx.State.Repositories.Scores.FetchScoreIndexById(
		score.Id, score.BeatmapId, score.Mode,
	)
	if err != nil {
		ctx.Logger.Error("Failed to fetch score rank", "id", score.Id, "error", err)
		InternalServerError(ctx)
		return
	}

	view := templates.ScoreView{
		DefaultView: buildDefaultView(ctx),
		User:        score.User,
		UserStats:   resolveStatsForMode(score.User, score.Mode),
		Beatmapset:  score.Beatmap.Beatmapset,
		Beatmap:     score.Beatmap,
		Score:       score,
		ScoreRank:   scoreRank,
	}
	ctx.RenderTemplate(http.StatusOK, "pages/public/score", view)
}

func resolveStatsForMode(user *schemas.User, mode constants.Mode) *schemas.Stats {
	for _, stats := range user.Stats {
		if stats.Mode == mode {
			return stats
		}
	}
	return nil
}
