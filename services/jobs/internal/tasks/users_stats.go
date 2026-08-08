package tasks

import (
	"fmt"
	"log/slog"

	"github.com/osuTitanic/titanic/internal/performance"
	"github.com/osuTitanic/titanic/internal/state"
)

// UpdateUsersStats rebuilds pp, accuracy, and global rank for all eligible users.
func UpdateUsersStats(app *state.State, logger *slog.Logger) error {
	criteria := map[string]any{
		"restricted = ?": false,
		"activated = ?":  true,
	}
	userList, err := app.Repositories.Users.Many(criteria)
	if err != nil {
		return fmt.Errorf("failed to fetch users: %w", err)
	}

	logger.Info(
		"Updating user stats...",
		"total_users", len(userList),
	)

	for _, user := range userList {
		userStats, err := app.Repositories.Stats.ManyByUserId(user.Id)
		if err != nil {
			return fmt.Errorf("failed to fetch stats for user %d: %w", user.Id, err)
		}

		for _, stats := range userStats {
			bestScores, err := app.Repositories.Scores.FetchBest(
				user.Id, stats.Mode,
				!app.Config.ApprovedMapRewards,
			)
			if err != nil {
				return fmt.Errorf("failed to fetch best scores for user %d mode %d: %w", user.Id, stats.Mode, err)
			}
			if len(bestScores) == 0 {
				continue
			}

			stats.PP = performance.CalculateWeightedPPv2(bestScores)
			stats.Acc = performance.CalculateWeightedAccuracy(bestScores)

			if err := app.Rankings.Update(stats, user.Country); err != nil {
				return fmt.Errorf("failed to update rankings for user %d mode %d: %w", user.Id, stats.Mode, err)
			}

			stats.Rank, err = app.Rankings.GlobalRank(user.Id, stats.Mode)
			if err != nil {
				return fmt.Errorf("failed to fetch global rank for user %d mode %d: %w", user.Id, stats.Mode, err)
			}

			if _, err := app.Repositories.Stats.Update(stats, "pp", "acc", "rank"); err != nil {
				return fmt.Errorf("failed to update stats for user %d mode %d: %w", user.Id, stats.Mode, err)
			}

			logger.Info(
				"Updated user stats",
				"user_id", user.Id,
				"mode", stats.Mode,
				"pp", stats.PP,
				"accuracy", stats.Acc,
				"rank", stats.Rank,
			)
		}
	}
	return nil
}

// TODO: add task to recalculate *full* stats of a user
