package tasks

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/osuTitanic/titanic-go/internal/schemas"
	"github.com/osuTitanic/titanic-go/internal/state"
)

const UserAutodeleteAge = 7 * 24 * time.Hour

// AutoDeleteUsers "deletes" inactive users & clears their verifications
func AutoDeleteUsers(app *state.State, logger *slog.Logger) error {
	cutoff := time.Now().UTC().Add(-UserAutodeleteAge)
	deletedEmailDomain := resolveEmailDomain(app)

	users, err := app.Repositories.Users.ManyAutoDeleteCandidates(cutoff)
	if err != nil {
		return fmt.Errorf("failed to fetch users eligible for auto-delete: %w", err)
	}

	logger.Info(
		"Auto-deleting inactive unactivated users...",
		"total_users", len(users), "cutoff", cutoff,
	)

	for _, user := range users {
		if err := deleteUser(app, logger, user.Id, cutoff, deletedEmailDomain); err != nil {
			return err
		}
	}

	return nil
}

func deleteUser(
	app *state.State,
	logger *slog.Logger,
	userId int,
	cutoff time.Time,
	deletedEmailDomain string,
) error {
	return app.DatabaseTransaction(func(repos *state.Repositories) error {
		user, err := repos.Users.ById(userId)
		if err != nil {
			return fmt.Errorf("failed to fetch user %d: %w", userId, err)
		}
		if user == nil || user.Activated || !user.CreatedAt.Before(cutoff) {
			return nil
		}

		statsEntries, err := repos.Stats.ManyByUserId(user.Id)
		if err != nil {
			return fmt.Errorf("failed to fetch stats for user %d: %w", user.Id, err)
		}
		if len(statsEntries) > 0 {
			return nil
		}

		previousName := user.Name
		user.Name = formatDeletedUserName(user.Id)
		user.Email = formatDeletedUserEmail(user.Id, deletedEmailDomain)
		user.SafeName = schemas.ResolveSafeName(user.Name)

		if _, err := repos.Users.Update(user, "name", "safe_name", "email"); err != nil {
			return fmt.Errorf("failed to auto-delete user %d: %w", user.Id, err)
		}

		deletedVerifications, err := repos.Verifications.DeleteByUserId(user.Id)
		if err != nil {
			return fmt.Errorf("failed to clear verifications for user %d: %w", user.Id, err)
		}

		logger.Info(
			"Auto-deleted user",
			"user_id", user.Id,
			"previous_name", previousName,
			"new_name", user.Name,
			"deleted_verifications", deletedVerifications,
		)
		return nil
	})
}

func formatDeletedUserName(userId int) string {
	return fmt.Sprintf("DeletedUser%d", userId)
}

func formatDeletedUserEmail(userId int, domain string) string {
	return fmt.Sprintf("deleteduser%d@%s", userId, domain)
}

func resolveEmailDomain(app *state.State) string {
	if app.Config == nil {
		return "example.com"
	}
	if domain := app.Config.EmailDomain(); domain != nil && *domain != "" {
		return *domain
	}
	if app.Config.DomainName != "" {
		return app.Config.DomainName
	}
	return "example.com"
}
