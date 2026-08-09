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

type PPv2RecalculationOptions struct {
	Workers   int
	BatchSize int
}

func DefaultPPv2RecalculationOptions() PPv2RecalculationOptions {
	return PPv2RecalculationOptions{
		BatchSize: 500,
		Workers:   4,
	}
}

func (o PPv2RecalculationOptions) Validate() error {
	if o.Workers < 0 {
		return fmt.Errorf("workers must be greater than or equal to zero")
	}
	if o.BatchSize < 1 {
		return fmt.Errorf("batch size must be greater than zero")
	}
	return nil
}

// RecalculatePPv2 recalculates ppv2 for all passed scores.
func RecalculatePPv2(app *state.State, logger *slog.Logger, options PPv2RecalculationOptions) error {
	if err := options.Validate(); err != nil {
		return fmt.Errorf("invalid ppv2 recalculation options: %w", err)
	}

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
		"batch_size", options.BatchSize,
		"workers", options.Workers,
	)
	var totalBatchDuration time.Duration

	return recalculatePPv2Batch(
		app, logger, options,
		// Next batch -> fetch scores from the database in batches
		func(offset, limit int) ([]*schemas.Score, error) {
			return app.Repositories.Scores.Many(
				criteria, "pp ASC",
				offset, limit,
			)
		},
		// On batch completion
		func(completedBatches, processedScores int, batchDuration time.Duration) {
			totalBatchDuration += batchDuration
			averageBatchDuration := totalBatchDuration / time.Duration(completedBatches)

			scoresLeft := max(int(totalScores)-processedScores, 0)
			batchesLeft := (scoresLeft + options.BatchSize - 1) / options.BatchSize
			estimatedRemaining := time.Duration(batchesLeft) * averageBatchDuration

			logger.Info(
				"Completed ppv2 recalculation batch",
				"completed_batches", completedBatches,
				"batches_left", batchesLeft,
				"average_batch_duration", averageBatchDuration,
				"estimated_remaining", estimatedRemaining,
				"estimated_completion", time.Now().Add(estimatedRemaining).Format(time.RFC1123),
			)
		},
	)
}

// RecalculatePPv2Failed recalculates scores whose ppv2 calculation previously failed.
func RecalculatePPv2Failed(app *state.State, logger *slog.Logger, options PPv2RecalculationOptions) error {
	if err := options.Validate(); err != nil {
		return fmt.Errorf("invalid ppv2 recalculation options: %w", err)
	}

	if !app.PPv2.Available() {
		logger.Info("PPv2 service is unavailable, skipping recalculation")
		return nil
	}

	scores, err := app.Repositories.Scores.FetchWithZeroPP()
	if err != nil {
		return fmt.Errorf("failed to fetch scores with zero pp: %w", err)
	}

	logger.Info(
		"Recalculating failed ppv2 calculations",
		"total_scores", len(scores),
		"batch_size", options.BatchSize,
		"workers", options.Workers,
	)

	return recalculatePPv2Batch(
		app, logger, options,
		// Next batch -> in this case we already have all the scores in memory, so we just slice them
		func(offset, limit int) ([]*schemas.Score, error) {
			if offset >= len(scores) {
				return nil, nil
			}
			end := min(offset+limit, len(scores))
			return scores[offset:end], nil
		},
		nil,
	)
}

func recalculatePPv2Batch(
	app *state.State,
	logger *slog.Logger,
	options PPv2RecalculationOptions,
	next func(offset, limit int) ([]*schemas.Score, error),
	onBatchDone func(completedBatches, processedScores int, batchDuration time.Duration),
) error {
	offset := 0
	completedBatches := 0

	for {
		scores, err := next(offset, options.BatchSize)
		if err != nil {
			return fmt.Errorf("failed to fetch scores for ppv2 recalculation: %w", err)
		}
		if len(scores) == 0 {
			break
		}

		workerCount := workers.TaskWorkerCount(app, len(scores), options.Workers)
		batchStarted := time.Now()

		if err := workers.RunWorkerPool(scores, workerCount, func(score *schemas.Score) error {
			pp, err := app.PPv2.CalculatePerformance(score)
			if err != nil {
				logger.Error(
					"Failed to recalculate ppv2",
					"score_id", score.Id,
					"beatmap_id", score.BeatmapId,
					"error", err,
				)
				return nil
			}

			score.PP = pp
			if _, err := app.Repositories.Scores.Update(score, "pp"); err != nil {
				logger.Error(
					"Failed to save recalculated ppv2",
					"score_id", score.Id,
					"beatmap_id", score.BeatmapId,
					"error", err,
				)
				return nil
			}

			logger.Debug(
				"Recalculated ppv2",
				"score_id", score.Id,
				"beatmap_id", score.BeatmapId,
				"pp", pp,
			)
			return nil
		}); err != nil {
			return fmt.Errorf("failed to recalculate ppv2: %w", err)
		}

		batchDuration := time.Since(batchStarted)

		offset += len(scores)
		completedBatches++

		if onBatchDone != nil {
			onBatchDone(completedBatches, offset, batchDuration)
		}
	}

	return nil
}
