package routes

import (
	"net/http"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/services/stern/internal/charts"
	"github.com/osuTitanic/titanic/services/stern/internal/server"
)

func BeatmapRatingChart(ctx *server.Context) {
	setId, err := ctx.PathValueInt("id")
	if err != nil {
		ctx.Response.WriteHeader(http.StatusNotFound)
		return
	}

	beatmapset, err := ctx.State.Beatmapsets.ById(setId)
	if err != nil {
		ctx.Logger.Error("Failed to fetch beatmapset for rating chart", "set", setId, "error", err)
		InternalServerError(ctx)
		return
	}
	if beatmapset == nil {
		ctx.Response.WriteHeader(http.StatusNotFound)
		return
	}

	distribution, err := ctx.State.Ratings.Distribution(setId)
	if err != nil {
		ctx.Logger.Error("Failed to fetch beatmap rating distribution", "set", setId, "error", err)
		InternalServerError(ctx)
		return
	}

	data, err := charts.GenerateBeatmapRatingChart(distribution)
	if err != nil {
		ctx.Logger.Error("Failed to generate beatmap rating chart", "set", setId, "error", err)
		InternalServerError(ctx)
		return
	}
	writeBeatmapChart(ctx, data)
}

func BeatmapFailChart(ctx *server.Context) {
	beatmapId, err := ctx.PathValueInt("id")
	if err != nil {
		ctx.Response.WriteHeader(http.StatusNotFound)
		return
	}

	beatmap, err := ctx.State.Beatmaps.ById(beatmapId)
	if err != nil {
		ctx.Logger.Error("Failed to fetch beatmap for fail chart", "beatmap", beatmapId, "error", err)
		InternalServerError(ctx)
		return
	}
	if beatmap == nil || beatmap.Status <= constants.BeatmapStatusInactive {
		ctx.Response.WriteHeader(http.StatusNotFound)
		return
	}
	mode := resolveMode(ctx, beatmap.Mode)

	distribution, err := ctx.State.Scores.FetchFailtimeDistribution(
		beatmap.Id, mode,
		beatmap.TotalLength,
	)
	if err != nil {
		ctx.Logger.Error(
			"Failed to fetch beatmap failtime distribution",
			"beatmap", beatmapId, "mode", mode, "error", err,
		)
		InternalServerError(ctx)
		return
	}

	data, err := charts.GenerateBeatmapFailChart(distribution.Fails, distribution.Exits)
	if err != nil {
		ctx.Logger.Error("Failed to generate beatmap fail chart", "beatmap", beatmapId, "error", err)
		InternalServerError(ctx)
		return
	}
	writeBeatmapChart(ctx, data)
}

func writeBeatmapChart(ctx *server.Context, data []byte) {
	ctx.Response.Header().Set("Content-Type", "image/png")
	ctx.Response.Header().Set("Cache-Control", "public, max-age=3600")
	ctx.Response.WriteHeader(http.StatusOK)
	if _, err := ctx.Response.Write(data); err != nil {
		ctx.Logger.Error("Failed to write beatmap chart", "error", err)
	}
}
