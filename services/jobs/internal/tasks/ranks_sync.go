package tasks

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/osuTitanic/titanic-go/internal/constants"
	"github.com/osuTitanic/titanic-go/internal/state"
)

// UpdateRanks updates rank history for all active users.
func UpdateRanks(app *state.State, logger *slog.Logger) error {
	if app.Config.FrozenRankUpdates {
		logger.Info("Rank updates are frozen, skipping...")
		return nil
	}

	criteria := map[string]any{
		"restricted = ?": false,
		"activated = ?":  true,
	}
	userList, err := app.Repositories.Users.Many(criteria, "Stats")
	if err != nil {
		return fmt.Errorf("failed to fetch users: %w", err)
	}
	logger.Info("Updating rank history...", "total_users", len(userList))

	for _, user := range userList {
		for _, userStats := range user.Stats {
			if userStats.Playcount <= 0 {
				continue
			}

			globalRank, err := app.Rankings.GlobalRank(user.Id, userStats.Mode)
			if err != nil {
				return fmt.Errorf("failed to fetch global rank for user %d mode %d: %w", user.Id, userStats.Mode, err)
			}

			// Check if rank has desynced from redis -> db & update if necessary
			rankChanged := userStats.Rank != globalRank
			if rankChanged {
				userStats.Rank = globalRank
				if _, err := app.Repositories.Stats.Update(userStats, "rank"); err != nil {
					return fmt.Errorf("failed to update current rank for user %d mode %d: %w", user.Id, userStats.Mode, err)
				}
			}

			// We want at least 1 rank history update per day
			needsHistoryUpdate, err := userRequiresHistoryUpdate(
				user.Id,
				userStats.Mode,
				app,
			)
			if err != nil {
				return fmt.Errorf("failed to check rank history for user %d mode %d: %w", user.Id, userStats.Mode, err)
			}

			// If the rank changed, we always want a history update
			needsHistoryUpdate = needsHistoryUpdate || rankChanged

			if needsHistoryUpdate {
				inserted, err := app.Repositories.Histories.UpdateRank(userStats, user.Country, app.Rankings)
				if err != nil {
					return fmt.Errorf("failed to update rank history for user %d mode %d: %w", user.Id, userStats.Mode, err)
				}
				if inserted {
					logger.Info(
						"Added rank history entry",
						"user_id", user.Id, "mode", userStats.Mode,
						"rank_changed", rankChanged, "daily_update", needsHistoryUpdate,
					)
				}
			}

			// Update peak rank if current rank is better than peak
			if userStats.PeakRank <= 0 || globalRank < userStats.PeakRank {
				userStats.PeakRank = globalRank
				if _, err := app.Repositories.Stats.Update(userStats, "peak_rank"); err != nil {
					return fmt.Errorf("failed to update peak rank for user %d mode %d: %w", user.Id, userStats.Mode, err)
				}
			}

			// Keep only the most recent rank history entry per day to save storage space
			err = cleanupRankHistory(
				user.Id,
				userStats.Mode,
				app,
			)
			if err != nil {
				return fmt.Errorf("failed to clean up rank history for user %d mode %d: %w", user.Id, userStats.Mode, err)
			}
		}

		logger.Info(
			"Updated user",
			"user_id", user.Id, "username", user.Name,
		)
	}
	return nil
}

func userRequiresHistoryUpdate(userId int, mode constants.Mode, app *state.State) (bool, error) {
	lastUpdate, err := app.Repositories.Histories.FetchLastRankHistoryEntry(userId, mode)
	if err != nil {
		return false, err
	}
	if lastUpdate == nil {
		return true, nil
	}
	return rankHistoryNeedsDailyUpdate(lastUpdate.Time, time.Now()), nil
}

func cleanupRankHistory(userId int, mode constants.Mode, app *state.State) error {
	// We only keep one rank history entry per day
	dayStart, dayEnd := rankHistoryDayBounds(time.Now())
	entries, err := app.Repositories.Histories.FetchRankHistoryEntriesBetween(userId, mode, dayStart, dayEnd)
	if err != nil {
		return fmt.Errorf("failed to fetch recent rank history entries for user %d: %w", userId, err)
	}
	if len(entries) <= 1 {
		return nil
	}

	// Delete all but the most recent entry
	for _, entry := range entries[1:] {
		_, err := app.Repositories.Histories.DeleteRankHistoryEntry(userId, mode, entry.Time)
		if err != nil {
			return fmt.Errorf("failed to delete rank history entry %v: %w", entry.Time, err)
		}
		// TODO: Should this delete all but the highest ranked entry instead?
		//       Currently not sure, but worth to consider for the future
	}
	return nil
}

func rankHistoryNeedsDailyUpdate(lastUpdate, now time.Time) bool {
	dayStart, _ := rankHistoryDayBounds(now)
	return lastUpdate.In(dayStart.Location()).Before(dayStart)
}

func rankHistoryDayBounds(t time.Time) (time.Time, time.Time) {
	dayStart := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return dayStart, dayStart.AddDate(0, 0, 1)
}
