package tasks

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/osuTitanic/titanic/internal/state"
)

const QueryActivityCutoff = 90 * 24 * time.Hour

// UnlistUsers ensures inactive & restricted users are removed from rankings
func UnlistUsers(app *state.State, logger *slog.Logger) error {
	cutoff := time.Now().Add(-QueryActivityCutoff)

	users, err := app.Repositories.Users.FetchInactiveOrRestricted(cutoff)
	if err != nil {
		return fmt.Errorf("failed to fetch inactive/restricted users: %w", err)
	}
	logger.Info("Unlisting inactive/restricted users from rankings...", "total_users", len(users))

	for _, user := range users {
		if err := app.Rankings.Remove(user.Id, user.Country); err != nil {
			return fmt.Errorf("failed to remove user %d from rankings: %w", user.Id, err)
		}
		logger.Info(
			"Removed user from rankings",
			"user_id", user.Id, "username", user.Name,
		)
	}
	return nil
}
