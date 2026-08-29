package tasks

import (
	"fmt"
	"log/slog"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/state"
)

// RecalculateEyupStars recalculates the eyup star rating for every beatmap
func RecalculateEyupStars(app *state.State, logger *slog.Logger) error {
	criteria := map[string]any{"status > ?": constants.BeatmapStatusInactive}
	beatmaps, err := app.Repositories.Beatmaps.Many(criteria)
	if err != nil {
		return fmt.Errorf("failed to fetch beatmaps for eyup recalculation: %w", err)
	}

	logger.Info(
		"Recalculating eyup star ratings...",
		"total_beatmaps", len(beatmaps),
	)

	for _, beatmap := range beatmaps {
		stars, err := app.PPv1.RecalculateEyupStarRating(beatmap)
		if err != nil {
			return fmt.Errorf("failed to recalculate eyup stars for beatmap %d: %w", beatmap.Id, err)
		}
		logger.Debug(
			"Recalculated eyup star rating",
			"beatmap_id", beatmap.Id, "stars", stars,
		)
	}

	logger.Info(
		"Recalculated eyup star ratings",
		"total_beatmaps", len(beatmaps),
	)
	return nil
}
