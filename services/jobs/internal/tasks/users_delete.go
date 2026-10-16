package tasks

import (
	"cmp"
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/internal/state"
	"github.com/osuTitanic/titanic/internal/storage"
)

const UserAutodeleteAge = 7 * 24 * time.Hour

// AutoDeleteUsers anonymizes inactive users & clears their verifications
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
			logger.Error(
				"Failed to delete user",
				"user_id", user.Id, "error", err,
			)
		}
	}

	return nil
}

// DeletePendingUsers anonymizes users whose pending deletion time has passed
func DeletePendingUsers(app *state.State, logger *slog.Logger) error {
	cutoff := time.Now().UTC()
	deletedEmailDomain := resolveEmailDomain(app)

	users, err := app.Repositories.Users.ManyPendingDeletionCandidates(cutoff)
	if err != nil {
		return fmt.Errorf("failed to fetch users with expired pending deletions: %w", err)
	}

	logger.Info(
		"Deleting users with expired pending deletions...",
		"total_users", len(users), "cutoff", cutoff,
	)

	for _, user := range users {
		if err := deletePendingUser(app, logger, user.Id, cutoff, deletedEmailDomain); err != nil {
			logger.Error(
				"Failed to delete user",
				"user_id", user.Id, "error", err,
			)
		}
	}

	return nil
}

func deleteUser(
	app *state.State,
	logger *slog.Logger,
	userId int,
	cutoff time.Time,
	emailDomain string,
) error {
	return app.DatabaseTransaction(func(repos *state.Repositories) error {
		user, err := repos.Users.ByIdWithLock(userId)
		if err != nil {
			return fmt.Errorf("failed to fetch user %d: %w", userId, err)
		}
		if user == nil {
			return nil
		}
		if user.Activated || !user.CreatedAt.Before(cutoff) {
			return nil
		}

		statsEntries, err := repos.Stats.ManyByUserId(user.Id)
		if err != nil {
			return fmt.Errorf("failed to fetch stats for user %d: %w", user.Id, err)
		}
		if len(statsEntries) > 0 {
			return nil
		}

		return anonymizeUser(user, repos, app.Storage, logger, emailDomain)
	})
}

func deletePendingUser(
	app *state.State,
	logger *slog.Logger,
	userId int,
	cutoff time.Time,
	deletedEmailDomain string,
) error {
	return app.DatabaseTransaction(func(repos *state.Repositories) error {
		user, err := repos.Users.ByIdWithLock(userId)
		if err != nil {
			return fmt.Errorf("failed to fetch user %d: %w", userId, err)
		}
		if user == nil {
			return nil
		}
		if user.PendingDeletion == nil {
			return nil
		}
		if user.PendingDeletion.After(cutoff) {
			return nil
		}
		return anonymizeUser(user, repos, app.Storage, logger, deletedEmailDomain)
	})
}

func anonymizeUser(
	user *schemas.User,
	transaction *state.Repositories,
	store storage.Storage,
	logger *slog.Logger,
	emailDomain string,
) error {
	previousName := user.Name
	user.Name = formatDeletedUserName(user.Id)
	user.Email = formatDeletedUserEmail(user.Id, emailDomain)
	user.SafeName = schemas.ResolveSafeName(user.Name)
	user.Bcrypt = strings.Repeat("-", 60)
	user.Country = "XX"
	user.DiscordId = nil
	user.PendingDeletion = nil
	user.SilenceEnd = nil
	user.Activated = false
	user.Restricted = true
	user.AvatarHash = nil
	user.AvatarLastUpdate = time.Now()

	user.Userpage = nil
	user.Signature = nil
	user.Title = nil
	user.Banner = nil
	user.Website = nil
	user.Discord = nil
	user.Twitter = nil
	user.Location = nil
	user.Interests = nil

	if _, err := transaction.Users.Update(user, "name", "safe_name", "email", "pending_deletion"); err != nil {
		return fmt.Errorf("failed to auto-delete user %d: %w", user.Id, err)
	}

	if _, err := transaction.Logins.AnonymizeForUser(user.Id); err != nil {
		return fmt.Errorf("failed to anonymize logins for user %d: %w", user.Id, err)
	}
	if err := transaction.Names.DeleteByUserId(user.Id); err != nil {
		return fmt.Errorf("failed to delete names for user %d: %w", user.Id, err)
	}
	if _, err := transaction.Activities.DeleteByUserId(user.Id); err != nil {
		return fmt.Errorf("failed to delete bancho activity for user %d: %w", user.Id, err)
	}
	if _, err := transaction.Verifications.DeleteByUserId(user.Id); err != nil {
		return fmt.Errorf("failed to clear verifications for user %d: %w", user.Id, err)
	}
	if err := transaction.Groups.DeleteByUserId(user.Id); err != nil {
		return fmt.Errorf("failed to delete groups for user %d: %w", user.Id, err)
	}
	if err := transaction.Permissions.DeleteUserPermissionByUserId(user.Id); err != nil {
		return fmt.Errorf("failed to delete user permissions for user %d: %w", user.Id, err)
	}
	if err := transaction.Favourites.DeleteByUserId(user.Id); err != nil {
		return fmt.Errorf("failed to delete beatmap favourites for user %d: %w", user.Id, err)
	}
	if err := transaction.Screenshots.DeleteByUserId(user.Id); err != nil {
		return fmt.Errorf("failed to delete screenshots for user %d: %w", user.Id, err)
	}

	store.Remove(context.TODO(), strconv.Itoa(user.Id), "avatars")
	// TODO: Remove screenshots

	// TODO: How do we handle beatmapsets?
	// TODO: Store messages in db by user ID and not username so that they get auto-anonymized

	logger.Info(
		"Auto-deleted user",
		"user_id", user.Id,
		"previous_name", previousName,
		"new_name", user.Name,
	)
	return nil
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
	return cmp.Or(app.Config.DomainName, "example.com")
}
