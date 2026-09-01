package tasks

import (
	"cmp"
	"fmt"
	"log/slog"
	"slices"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/internal/state"
	"github.com/osuTitanic/titanic/services/jobs/internal/workers"
)

type StatusRecalculationOptions struct {
	UserId int
	Mode   constants.Mode
}

func (options StatusRecalculationOptions) Validate() error {
	if options.UserId < 0 {
		return fmt.Errorf("userId must be greater than zero")
	}
	return nil
}

const statusRecalculationWorkers = 4

func RecalculateScoreStatusUser(app *state.State, logger *slog.Logger, options StatusRecalculationOptions) error {
	return performStatusRecalculationForUser(
		app, logger, options, "score",
		compareScoreLeaderboardPosition,
	)
}

func RecalculatePPStatusUser(app *state.State, logger *slog.Logger, options StatusRecalculationOptions) error {
	return performStatusRecalculationForUser(
		app, logger, options, "pp",
		compareScorePP,
	)
}

func RecalculateScoreStatusAll(app *state.State, logger *slog.Logger) error {
	criteria := map[string]any{
		"restricted = ?": false,
		"activated = ?":  true,
	}
	userList, err := app.Repositories.Users.Many(criteria)
	if err != nil {
		return fmt.Errorf("failed to fetch users: %w", err)
	}

	logger.Info("Recalculating score status for all users", "users", len(userList))
	workerCount := workers.TaskWorkerCount(len(userList), statusRecalculationWorkers)

	return workers.RunWorkerPool(userList, workerCount, func(user *schemas.User) error {
		err := RecalculateScoreStatusUser(app, logger, StatusRecalculationOptions{
			UserId: user.Id,
			Mode:   constants.ModeOsu,
		})
		if err != nil {
			logger.Error("Failed to recalculate score status for user", "userId", user.Id, "error", err)
			return nil
		}
		err = RecalculateScoreStatusUser(app, logger, StatusRecalculationOptions{
			UserId: user.Id,
			Mode:   constants.ModeTaiko,
		})
		if err != nil {
			logger.Error("Failed to recalculate score status for user", "userId", user.Id, "error", err)
			return nil
		}
		err = RecalculateScoreStatusUser(app, logger, StatusRecalculationOptions{
			UserId: user.Id,
			Mode:   constants.ModeCatch,
		})
		if err != nil {
			logger.Error("Failed to recalculate score status for user", "userId", user.Id, "error", err)
			return nil
		}
		err = RecalculateScoreStatusUser(app, logger, StatusRecalculationOptions{
			UserId: user.Id,
			Mode:   constants.ModeMania,
		})
		if err != nil {
			logger.Error("Failed to recalculate score status for user", "userId", user.Id, "error", err)
			return nil
		}
		return nil
	})
}

func RecalculatePPStatusAll(app *state.State, logger *slog.Logger) error {
	criteria := map[string]any{
		"restricted = ?": false,
		"activated = ?":  true,
	}
	userList, err := app.Repositories.Users.Many(criteria)
	if err != nil {
		return fmt.Errorf("failed to fetch users: %w", err)
	}

	logger.Info("Recalculating pp status for all users", "users", len(userList))
	workerCount := workers.TaskWorkerCount(len(userList), statusRecalculationWorkers)

	return workers.RunWorkerPool(userList, workerCount, func(user *schemas.User) error {
		err := RecalculatePPStatusUser(app, logger, StatusRecalculationOptions{
			UserId: user.Id,
			Mode:   constants.ModeOsu,
		})
		if err != nil {
			logger.Error("Failed to recalculate pp status for user", "userId", user.Id, "error", err)
			return nil
		}
		err = RecalculatePPStatusUser(app, logger, StatusRecalculationOptions{
			UserId: user.Id,
			Mode:   constants.ModeTaiko,
		})
		if err != nil {
			logger.Error("Failed to recalculate pp status for user", "userId", user.Id, "error", err)
			return nil
		}
		err = RecalculatePPStatusUser(app, logger, StatusRecalculationOptions{
			UserId: user.Id,
			Mode:   constants.ModeCatch,
		})
		if err != nil {
			logger.Error("Failed to recalculate pp status for user", "userId", user.Id, "error", err)
			return nil
		}
		err = RecalculatePPStatusUser(app, logger, StatusRecalculationOptions{
			UserId: user.Id,
			Mode:   constants.ModeMania,
		})
		if err != nil {
			logger.Error("Failed to recalculate pp status for user", "userId", user.Id, "error", err)
			return nil
		}
		return nil
	})
}

func performStatusRecalculationForUser(
	app *state.State, logger *slog.Logger, options StatusRecalculationOptions,
	statusType string, compare func(a *schemas.Score, b *schemas.Score) int,
) error {
	if err := options.Validate(); err != nil {
		return fmt.Errorf("invalid status recalculation options: %w", err)
	}

	user, err := app.Repositories.Users.ById(options.UserId)
	if err != nil {
		return fmt.Errorf("failed to resolve user '%d': %v", options.UserId, err)
	}
	if user == nil {
		return fmt.Errorf("user '%d' not found", options.UserId)
	}

	passedScores, err := app.Repositories.Scores.FetchPassed(user.Id, options.Mode)
	if err != nil {
		return fmt.Errorf("failed to fetch passed scores: %v", err)
	}
	if len(passedScores) <= 0 {
		logger.Info(
			"Skipping status recalculation",
			"username", user.Name, "mode", options.Mode, "scores", len(passedScores),
		)
		return nil
	}

	logger.Info(
		"Performing status recalculation",
		"type", statusType, "userId", options.UserId,
		"mode", options.Mode.String(), "scores", len(passedScores),
	)

	updateStatus := func(score *schemas.Score, status constants.ScoreStatus) error {
		var column string

		if statusType == "pp" {
			score.StatusPP = status
			column = "status"
		} else {
			score.StatusScore = status
			column = "status_score"
		}

		if _, err := app.Repositories.Scores.Update(score, column); err != nil {
			return fmt.Errorf("failed to update score '%d' to status %d: %w", score.Id, status, err)
		}
		return nil
	}

	scoresByBeatmap := make(map[int][]*schemas.Score)

	for _, score := range passedScores {
		if score.Relaxing() {
			// Exclude rx / ap from global rankings
			if err := updateStatus(score, constants.ScoreStatusSubmitted); err != nil {
				return err
			}
			continue
		}

		scoresByBeatmap[score.BeatmapId] = append(
			scoresByBeatmap[score.BeatmapId],
			score,
		)
	}

	// TODO: Use database transaction for each beatmap

	for beatmapId, scores := range scoresByBeatmap {
		// Sort scores by criteria, best first
		slices.SortFunc(scores, compare)

		bestScore := scores[0]
		if err := updateStatus(bestScore, constants.ScoreStatusBest); err != nil {
			return err
		}

		logger.Debug(
			"Determined best score",
			"userId", user.Id, "mode", bestScore.Mode,
			"beatmapId", beatmapId, "scoreId", bestScore.Id, "clears", len(scores),
		)

		// Sort scores by mods & select best scores for each mod combination
		// They will be assigned the status 4
		scoresByMods := make(map[constants.Mods][]*schemas.Score)

		for _, score := range scores {
			scoresByMods[score.Mods] = append(
				scoresByMods[score.Mods],
				score,
			)
		}

		for mods, moddedScores := range scoresByMods {
			slices.SortFunc(moddedScores, compare)

			bestScoreWithMods := moddedScores[0]

			for _, score := range moddedScores[1:] {
				if err := updateStatus(score, constants.ScoreStatusSubmitted); err != nil {
					return err
				}
			}

			if mods == bestScore.Mods {
				// Don't update the status of the best score
				continue
			}

			if err := updateStatus(bestScoreWithMods, constants.ScoreStatusMods); err != nil {
				return err
			}
		}
	}
	return nil
}

func compareScorePP(a, b *schemas.Score) int {
	return cmp.Or(
		// Sorted by pp, then leaderboard position if equal
		cmp.Compare(b.PP, a.PP),
		compareScoreLeaderboardPosition(a, b),
	)
}

func compareScoreLeaderboardPosition(a, b *schemas.Score) int {
	return cmp.Or(
		// Sorted by total_score DESC, submitted_at ASC, id ASC
		cmp.Compare(b.TotalScore, a.TotalScore),
		a.SubmittedAt.Compare(b.SubmittedAt),
		cmp.Compare(a.Id, b.Id),
	)
}
