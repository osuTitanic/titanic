package routes

import (
	"net/http"
	"time"

	"github.com/osuTitanic/titanic/internal/caching"
	"github.com/osuTitanic/titanic/services/stern/internal/charts"
	"github.com/osuTitanic/titanic/services/stern/internal/server"
)

const (
	activityWidth  = 470
	activityHeight = 70
	activityWindow = 24 * time.Hour
)

// We want to cache the activity chart for a short period.
// The operation itself is not very expensive, but its still a good idea nonetheless to avoid generating it on every request.

var activityChartCache = caching.NewValue[[]byte](time.Minute)

func ActivityChart(ctx *server.Context) {
	chart, err := activityChartCache.GetOrLoad(func() ([]byte, error) {
		now := time.Now()
		entries, err := ctx.State.Activity.FetchRange(now.Add(-activityWindow), now)
		if err != nil {
			ctx.Logger.Error("Failed to fetch user activity", "error", err)
			return nil, err
		}

		chart, err := charts.GenerateActivityChart(entries, activityWidth, activityHeight)
		if err != nil {
			ctx.Logger.Error("Failed to generate activity chart", "error", err)
			return nil, err
		}
		return chart, nil
	})
	if err != nil {
		InternalServerError(ctx)
		return
	}

	ctx.Response.Header().Set("Content-Type", "image/png")
	ctx.Response.Header().Set("Cache-Control", "public, max-age=60") // similarly, also cache browser-sided
	ctx.Response.WriteHeader(http.StatusOK)
	ctx.Response.Write(chart)
}
