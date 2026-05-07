package tasks

import (
	"log/slog"
	"sort"

	"github.com/osuTitanic/titanic-go/internal/performance"
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"github.com/osuTitanic/titanic-go/internal/state"
	"github.com/osuTitanic/titanic-go/services/jobs/internal/workers"
)

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

	sort.Slice(userList, func(i, j int) bool {
		return resolveUserPPv1(userList[i]) > resolveUserPPv1(userList[j])
	})

	logger.Info(
		"Starting ppv1 update workers",
		"workers", ppv1UpdateWorkerCount(app, len(userList)),
	)
	return updatePPv1ForUsers(app, logger, userList)
}

func updatePPv1ForUsers(app *state.State, logger *slog.Logger, users []*schemas.User) error {
	workerCount := ppv1UpdateWorkerCount(app, len(users))
	return workers.RunWorkerPool(users, workerCount, func(user *schemas.User) error {
		if err := updatePPv1ForUser(app, logger, user); err != nil {
			logger.Error("Failed to update user", "id", user.Id, "error", err)
		}
		return nil
	})
}

func ppv1UpdateWorkerCount(app *state.State, userCount int) int {
	return workers.TaskWorkerCount(app, userCount, ppv1UpdateWorkers)
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
