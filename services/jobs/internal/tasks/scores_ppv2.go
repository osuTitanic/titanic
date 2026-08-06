package tasks

import (
	"fmt"
	"log/slog"

	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/internal/state"
	"github.com/osuTitanic/titanic/services/jobs/internal/workers"
)

const ppv2RecalculationWorkers = 8

// RecalculatePPv2Failed recalculates scores whose ppv2 calculation previously failed.
func RecalculatePPv2Failed(app *state.State, logger *slog.Logger) error {
	if !app.PPv2.Available() {
		logger.Info("PPv2 service is unavailable, skipping recalculation")
		return nil
	}

	scores, err := app.Repositories.Scores.FetchWithZeroPP()
	if err != nil {
		return fmt.Errorf("failed to fetch scores with zero pp: %w", err)
	}

	workerCount := workers.TaskWorkerCount(app, len(scores), ppv2RecalculationWorkers)
	logger.Info(
		"Recalculating failed ppv2 calculations",
		"total_scores", len(scores), "workers", workerCount,
	)

	return workers.RunWorkerPool(scores, workerCount, func(score *schemas.Score) error {
		pp, err := app.PPv2.CalculatePerformance(score)
		if err != nil {
			logger.Error(
				"Failed to recalculate ppv2",
				"score_id", score.Id, "beatmap_id", score.BeatmapId, "error", err,
			)
			return nil
		}

		score.PP = pp
		if _, err := app.Repositories.Scores.Update(score, "pp"); err != nil {
			logger.Error(
				"Failed to save recalculated ppv2",
				"score_id", score.Id, "beatmap_id", score.BeatmapId, "error", err,
			)
			return nil
		}

		logger.Debug(
			"Recalculated ppv2",
			"score_id", score.Id, "beatmap_id", score.BeatmapId, "pp", pp,
		)
		return nil
	})
}
