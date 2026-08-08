package tasks

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/internal/state"
	"github.com/osuTitanic/titanic/services/jobs/internal/workers"
)

// TODO: make task config
const ppv2RecalculationWorkers = 12
const ppv2RecalculationBatchSize = 500

// RecalculatePPv2 recalculates ppv2 for all passed scores.
func RecalculatePPv2(app *state.State, logger *slog.Logger) error {
	if !app.PPv2.Available() {
		logger.Info("PPv2 service is unavailable, skipping recalculation")
		return nil
	}

	criteria := map[string]any{
		"status >= ?": constants.ScoreStatusSubmitted,
		"hidden = ?":  false,
	}

	totalScores, err := app.Repositories.Scores.Count(criteria)
	if err != nil {
		return fmt.Errorf("failed to count scores for ppv2 recalculation: %w", err)
	}

	logger.Info(
		"Starting ppv2 recalculation",
		"total_scores", totalScores,
		"batch_size", ppv2RecalculationBatchSize,
		"workers", ppv2RecalculationWorkers,
	)

	currentOffset := 0
	completedBatches := 0
	var totalBatchDuration time.Duration

	for {
		scores, err := app.Repositories.Scores.Many(
			criteria, "pp ASC, id ASC",
			currentOffset, ppv2RecalculationBatchSize,
		)
		if err != nil {
			return fmt.Errorf("failed to fetch scores for ppv2 recalculation: %w", err)
		}
		if len(scores) == 0 {
			break
		}

		workerCount := workers.TaskWorkerCount(app, len(scores), ppv2RecalculationWorkers)
		batchStarted := time.Now()

		if err := workers.RunWorkerPool(scores, workerCount, func(score *schemas.Score) error {
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
		}); err != nil {
			return fmt.Errorf("failed to recalculate ppv2: %w", err)
		}

		totalBatchDuration += time.Since(batchStarted)
		completedBatches++

		batchesLeft := (int(totalScores) - currentOffset - len(scores)) / ppv2RecalculationBatchSize
		estimatedRemaining := time.Duration(batchesLeft) * (totalBatchDuration / time.Duration(completedBatches))
		estimatedCompletion := time.Now().Add(estimatedRemaining)

		logger.Info(
			"Completed ppv2 recalculation batch",
			"completed_batches", completedBatches,
			"batches_left", batchesLeft,
			"average_batch_duration", totalBatchDuration/time.Duration(completedBatches),
			"estimated_remaining", estimatedRemaining,
			"estimated_completion", estimatedCompletion.Format(time.RFC1123),
		)
		currentOffset += len(scores)
	}

	return nil
}

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
