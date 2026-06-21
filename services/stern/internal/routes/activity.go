package routes

import (
	"net/http"
	"time"

	"github.com/osuTitanic/titanic-go/services/stern/internal/charts"
	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
)

const (
	activityWidth  = 470
	activityHeight = 70
	activityWindow = 24 * time.Hour
)

func ActivityChart(ctx *server.Context) {
	now := time.Now()
	entries, err := ctx.State.Activity.FetchRange(now.Add(-activityWindow), now)
	if err != nil {
		ctx.Logger.Error("Failed to fetch user activity", "error", err)
		InternalServerError(ctx)
		return
	}

	chart, err := charts.GenerateActivityChart(entries, activityWidth, activityHeight)
	if err != nil {
		ctx.Logger.Error("Failed to generate activity chart", "error", err)
		InternalServerError(ctx)
		return
	}

	ctx.Response.Header().Set("Content-Type", "image/png")
	ctx.Response.WriteHeader(http.StatusOK)
	ctx.Response.Write(chart)
}
