package tasks

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/osuTitanic/titanic/internal/state"
)

// UpdateWebsiteStats updates the stats required for the website & api stats
func UpdateWebsiteStats(app *state.State, logger *slog.Logger) error {
	ctx := context.Background()

	userCount, err := app.Repositories.Users.GetCount()
	if err != nil {
		return fmt.Errorf("failed to fetch user count: %w", err)
	}
	app.Redis.Set(ctx, "bancho:totalusers", userCount, 0)

	scoreCount, err := app.Repositories.Scores.GetCount()
	if err != nil {
		return fmt.Errorf("failed to fetch scores count: %w", err)
	}
	app.Redis.Set(ctx, "bancho:totalscores", scoreCount, 0)

	beatmapCount, err := app.Repositories.Beatmaps.GetCount()
	if err != nil {
		return fmt.Errorf("failed to fetch beatmap count: %w", err)
	}
	app.Redis.Set(ctx, "bancho:totalbeatmaps", beatmapCount, 0)

	setCount, err := app.Repositories.Beatmapsets.GetCount()
	if err != nil {
		return fmt.Errorf("failed to fetch beatmapset count: %w", err)
	}
	app.Redis.Set(ctx, "bancho:totalbeatmapsets", setCount, 0)

	for mode := range 4 {
		counts, err := app.Repositories.Beatmaps.GetCountGroupedByStatus(mode)
		if err != nil {
			return fmt.Errorf("failed to fetch beatmap stats %d: %w", mode, err)
		}

		for status, count := range counts {
			key := fmt.Sprintf("bancho:totalbeatmaps:%d:%d", mode, status)
			app.Redis.Set(ctx, key, count, 0)
		}
	}

	logger.Info("Updated website stats", "users", userCount, "scores", scoreCount, "beatmaps", beatmapCount, "beatmapsets", setCount)
	return nil
}
