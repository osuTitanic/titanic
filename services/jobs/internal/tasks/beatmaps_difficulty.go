package tasks

import (
	"fmt"
	"log/slog"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/internal/state"
	"github.com/osuTitanic/titanic/services/jobs/internal/workers"
)

const difficultyRecalculationWorkers = 2

// RecalculateDifficulty recalculates the NM star rating for every active beatmap.
func RecalculateDifficulty(app *state.State, logger *slog.Logger) error {
	if !app.PPv2.Available() {
		logger.Info("PPv2 service is unavailable, skipping difficulty recalculation")
		return nil
	}

	criteria := map[string]any{"status >= ?": constants.BeatmapStatusWIP}
	beatmapList, err := app.Beatmaps.Many(criteria)
	if err != nil {
		return fmt.Errorf("failed to fetch beatmaps for difficulty recalculation: %w", err)
	}

	logger.Info(
		"Recalculating beatmap difficulty...",
		"total_beatmaps", len(beatmapList),
	)
	workerCount := workers.TaskWorkerCount(len(beatmapList), difficultyRecalculationWorkers)

	if err := workers.RunWorkerPool(beatmapList, workerCount, func(beatmap *schemas.Beatmap) error {
		attributes, err := app.PPv2.CalculateDifficulty(beatmap.Id, beatmap.Mode, constants.NoMod)
		if err != nil {
			return fmt.Errorf("failed to recalculate difficulty for beatmap %d: %w", beatmap.Id, err)
		}

		if beatmap.Diff != attributes.StarRating {
			beatmap.Diff = attributes.StarRating
			if _, err := app.Beatmaps.Update(beatmap, "diff"); err != nil {
				return fmt.Errorf("failed to save recalculated difficulty for beatmap %d: %w", beatmap.Id, err)
			}
		}

		logger.Debug(
			"Recalculated beatmap difficulty",
			"beatmap_id", beatmap.Id,
			"stars", attributes.StarRating,
		)
		return nil
	}); err != nil {
		return err
	}

	logger.Info(
		"Recalculated beatmap difficulty",
		"total_beatmaps", len(beatmapList),
	)
	return nil
}
