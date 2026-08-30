package tasks

import (
	"cmp"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/performance"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/internal/state"
	"github.com/osuTitanic/titanic/services/jobs/internal/workers"
)

type PPv1RecalculationOptions struct {
	Workers   int
	BatchSize int
}

func DefaultPPv1RecalculationOptions() PPv1RecalculationOptions {
	return PPv1RecalculationOptions{
		BatchSize: 500,
		Workers:   4,
	}
}

func (o PPv1RecalculationOptions) Validate() error {
	if o.Workers < 0 {
		return fmt.Errorf("workers must be greater than or equal to zero")
	}
	if o.BatchSize < 1 {
		return fmt.Errorf("batch size must be greater than zero")
	}
	return nil
}

var ppv1UpdateWorkers = 8

// UpdatePPv1 updates ppv1 calculations for all users
func UpdatePPv1(app *state.State, logger *slog.Logger) error {
	if app.Config.FrozenPPv1Updates {
		logger.Info("ppv1 updates are disabled, skipping...")
		return nil
	}

	criteria := map[string]any{
		"restricted = ?": false,
		"activated = ?":  true,
	}
	userList, err := app.Repositories.Users.Many(criteria, "Stats")
	if err != nil {
		return err
	}
	logger.Info("Updating ppv1 calculations...", "total_users", len(userList))

	slices.SortFunc(userList, func(a, b *schemas.User) int {
		return cmp.Compare(resolveUserPPv1(b), resolveUserPPv1(a))
	})

	logger.Info(
		"Starting ppv1 update workers",
		"workers", ppv1UpdateWorkerCount(len(userList)),
	)
	return updatePPv1ForUsers(app, logger, userList)
}

// UpdatePPv1All recalculates ppv1 for all submitted scores
func UpdatePPv1All(app *state.State, logger *slog.Logger, options PPv1RecalculationOptions) error {
	if err := options.Validate(); err != nil {
		return fmt.Errorf("invalid ppv1 recalculation options: %w", err)
	}

	criteria := map[string]any{
		"status >= ?": constants.ScoreStatusSubmitted,
		"hidden = ?":  false,
	}

	totalScores, err := app.Repositories.Scores.Count(criteria)
	if err != nil {
		return fmt.Errorf("failed to count scores for ppv1 recalculation: %w", err)
	}

	logger.Info(
		"Starting ppv1 recalculation",
		"total_scores", totalScores,
		"batch_size", options.BatchSize,
		"workers", options.Workers,
	)
	var totalBatchDuration time.Duration

	offset := 0
	completedBatches := 0

	for {
		scores, err := app.Repositories.Scores.Many(
			criteria, "id ASC",
			offset, options.BatchSize,
			"Beatmap",
		)
		if err != nil {
			return fmt.Errorf("failed to fetch scores for ppv1 recalculation: %w", err)
		}
		if len(scores) == 0 {
			break
		}

		workerCount := workers.TaskWorkerCount(len(scores), options.Workers)
		batchStarted := time.Now()

		if err := workers.RunWorkerPool(scores, workerCount, func(score *schemas.Score) error {
			ppv1, err := app.PPv1.CalculatePerformance(score)
			if err != nil {
				return fmt.Errorf("failed to recalculate ppv1 for score %d: %w", score.Id, err)
			}

			if score.PPv1 != ppv1 {
				score.PPv1 = ppv1
				if _, err := app.Repositories.Scores.Update(score, "ppv1"); err != nil {
					return fmt.Errorf("failed to save recalculated ppv1 for score %d: %w", score.Id, err)
				}
			}

			logger.Debug(
				"Recalculated ppv1",
				"score_id", score.Id,
				"beatmap_id", score.BeatmapId,
				"ppv1", ppv1,
			)
			return nil
		}); err != nil {
			return fmt.Errorf("failed to recalculate ppv1: %w", err)
		}

		batchDuration := time.Since(batchStarted)
		totalBatchDuration += batchDuration
		offset += len(scores)
		completedBatches++

		averageBatchDuration := totalBatchDuration / time.Duration(completedBatches)
		scoresLeft := max(int(totalScores)-offset, 0)
		batchesLeft := (scoresLeft + options.BatchSize - 1) / options.BatchSize
		estimatedRemaining := time.Duration(batchesLeft) * averageBatchDuration

		logger.Info(
			"Completed ppv1 recalculation batch",
			"completed_batches", completedBatches,
			"batches_left", batchesLeft,
			"average_batch_duration", averageBatchDuration,
			"estimated_remaining", estimatedRemaining,
			"estimated_completion", time.Now().Add(estimatedRemaining).Format(time.RFC1123),
		)
	}

	return nil
}

func updatePPv1ForUsers(app *state.State, logger *slog.Logger, users []*schemas.User) error {
	workerCount := ppv1UpdateWorkerCount(len(users))
	return workers.RunWorkerPool(users, workerCount, func(user *schemas.User) error {
		if err := updatePPv1ForUser(app, logger, user); err != nil {
			logger.Error("Failed to update user", "id", user.Id, "error", err)
		}
		return nil
	})
}

func ppv1UpdateWorkerCount(userCount int) int {
	return workers.TaskWorkerCount(userCount, ppv1UpdateWorkers)
}

func updatePPv1ForUser(app *state.State, logger *slog.Logger, user *schemas.User) error {
	return app.DatabaseTransaction(func(repositories *state.Repositories) error {
		ppv1 := performance.NewPPv1Service(
			repositories.Scores,
			repositories.Beatmaps,
		)

		for _, stats := range user.Stats {
			if stats.Playcount <= 0 {
				continue
			}

			bestScores, err := repositories.Scores.FetchBest(
				user.Id,
				stats.Mode,
				!app.Config.ApprovedMapRewards,
				"Beatmap",
			)
			if err != nil {
				return err
			}
			if len(bestScores) == 0 {
				continue
			}

			stats.PPv1, err = ppv1.RecalculateWeightFromScores(bestScores)
			if err != nil {
				return err
			}

			repositories.Stats.Update(stats, "ppv1")
			app.Rankings.Update(stats, user.Country)

			if !app.Config.FrozenRankUpdates {
				repositories.Histories.UpdateRank(stats, user.Country, app.Rankings)
			}

			logger.Debug(
				"ppv1 update",
				"id", user.Id, "name", user.Name,
				"mode", stats.Mode, "ppv1", stats.PPv1,
			)
		}

		logger.Info("Updated ppv1 for user", "name", user.Name, "id", user.Id)
		return nil
	})
}

func resolveUserPPv1(user *schemas.User) float64 {
	var totalPPv1 float64
	for _, stats := range user.Stats {
		totalPPv1 += stats.PPv1
	}
	return totalPPv1
}
