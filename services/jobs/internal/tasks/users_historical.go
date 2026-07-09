package tasks

import (
	"fmt"
	"log/slog"
	"sort"
	"time"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/internal/state"
)

// FixHistoricalData ensures there are no missing monthly entries in replay & play histories.
func FixHistoricalData(app *state.State, logger *slog.Logger) error {
	criteria := map[string]any{
		"restricted = ?": false,
		"activated = ?":  true,
	}
	userList, err := app.Repositories.Users.Many(criteria)
	if err != nil {
		return fmt.Errorf("failed to fetch users: %w", err)
	}

	for _, user := range userList {
		if err := fixHistoricalDataForUser(app, logger, user.Id); err != nil {
			return err
		}
	}
	return nil
}

func fixHistoricalDataForUser(app *state.State, logger *slog.Logger, userId int) error {
	for _, mode := range constants.Modes {
		if err := fixReplayHistoryForUser(app, logger, userId, mode); err != nil {
			return err
		}
		if err := fixPlayHistoryForUser(app, logger, userId, mode); err != nil {
			return err
		}
	}
	return nil
}

func fixReplayHistoryForUser(app *state.State, logger *slog.Logger, userId int, mode constants.Mode) error {
	user, err := app.Repositories.Users.ById(userId)
	if err != nil {
		return fmt.Errorf("failed to fetch user %d: %w", userId, err)
	}
	if user == nil {
		logger.Warn("User was not found.", "user_id", userId)
		return nil
	}

	replayEntries, err := app.Repositories.Histories.FetchReplayHistoryAll(user.Id, mode)
	if err != nil {
		return fmt.Errorf("failed to fetch replay history for user %d mode %d: %w", user.Id, mode, err)
	}

	// Sort replay entries by date
	sort.Slice(replayEntries, func(i, j int) bool {
		if replayEntries[i].Year == replayEntries[j].Year {
			return replayEntries[i].Month < replayEntries[j].Month
		}
		return replayEntries[i].Year < replayEntries[j].Year
	})

	if len(replayEntries) == 0 {
		return nil
	}
	lastEntry := replayEntries[0]

	for _, entry := range replayEntries[1:] {
		lastEntryDate := lastEntry.Date()
		thisEntryDate := entry.Date()

		// Iterate through each month between lastEntry and thisEntry, adding missing entries
		for nextDate := lastEntryDate.AddDate(0, 1, 0); nextDate.Before(thisEntryDate); nextDate = nextDate.AddDate(0, 1, 0) {
			missingEntry := &schemas.ReplayHistory{
				UserId:    user.Id,
				Mode:      mode,
				Year:      nextDate.Year(),
				Month:     int(nextDate.Month()),
				CreatedAt: time.Date(nextDate.Year(), nextDate.Month(), 1, 0, 0, 0, 0, time.UTC),
			}

			if err := app.Repositories.Histories.CreateReplayHistory(missingEntry); err != nil {
				return fmt.Errorf("failed to add missing replay history entry for user %d mode %d: %w", user.Id, mode, err)
			}

			logger.Info(
				"Added missing replay history entry.",
				"user_id", user.Id,
				"username", user.Name,
				"mode", mode,
				"year", nextDate.Year(),
				"month", int(nextDate.Month()),
			)
		}

		lastEntry = entry
	}

	return nil
}

func fixPlayHistoryForUser(app *state.State, logger *slog.Logger, userId int, mode constants.Mode) error {
	user, err := app.Repositories.Users.ById(userId)
	if err != nil {
		return fmt.Errorf("failed to fetch user %d: %w", userId, err)
	}
	if user == nil {
		logger.Warn("User was not found.", "user_id", userId)
		return nil
	}

	playEntries, err := app.Repositories.Histories.FetchPlaysHistoryAll(user.Id, mode)
	if err != nil {
		return fmt.Errorf("failed to fetch play history for user %d mode %d: %w", user.Id, mode, err)
	}

	// Sort play entries by date
	sort.Slice(playEntries, func(i, j int) bool {
		if playEntries[i].Year == playEntries[j].Year {
			return playEntries[i].Month < playEntries[j].Month
		}
		return playEntries[i].Year < playEntries[j].Year
	})

	if len(playEntries) == 0 {
		return nil
	}
	lastEntry := playEntries[0]

	for _, entry := range playEntries[1:] {
		lastEntryDate := lastEntry.Date()
		thisEntryDate := entry.Date()

		// Iterate through each month between lastEntry and thisEntry, adding missing entries
		for nextDate := lastEntryDate.AddDate(0, 1, 0); nextDate.Before(thisEntryDate); nextDate = nextDate.AddDate(0, 1, 0) {
			missingEntry := &schemas.PlayHistory{
				UserId:    user.Id,
				Mode:      mode,
				Year:      nextDate.Year(),
				Month:     int(nextDate.Month()),
				CreatedAt: time.Date(nextDate.Year(), nextDate.Month(), 1, 0, 0, 0, 0, time.UTC),
			}

			if err := app.Repositories.Histories.CreatePlayHistory(missingEntry); err != nil {
				return fmt.Errorf("failed to add missing play history entry for user %d mode %d: %w", user.Id, mode, err)
			}

			logger.Info(
				"Added missing play history entry.",
				"user_id", user.Id,
				"username", user.Name,
				"mode", mode,
				"year", nextDate.Year(),
				"month", int(nextDate.Month()),
			)
		}

		lastEntry = entry
	}

	return nil
}
