package tasks

import (
	"fmt"
	"log/slog"
	"sort"
	"strings"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/internal/state"
)

// UpdateUnreadDmNotifications creates unread DM notifications for users.
func UpdateUnreadDmNotifications(app *state.State, logger *slog.Logger) error {
	criteria := map[string]any{
		"restricted = ?": false,
		"activated = ?":  true,
	}
	users, err := app.Repositories.Users.Many(criteria)
	if err != nil {
		return fmt.Errorf("failed to fetch users: %w", err)
	}
	logger.Info("Updating unread chat message notifications...", "total_users", len(users))

	for _, user := range users {
		message, link, err := generateUnreadDmNotification(app, user.Id)
		if err != nil {
			return fmt.Errorf("failed to generate unread chat notification for user %d: %w", user.Id, err)
		}
		if message == "" {
			continue
		}

		// Delete old chat notifications to avoid duplicates
		if _, err := app.Repositories.Notifications.DeleteByType(user.Id, constants.NotificationTypeChat); err != nil {
			return fmt.Errorf("failed to delete old chat notifications for user %d: %w", user.Id, err)
		}

		notification := &schemas.Notification{
			UserId:  user.Id,
			Type:    constants.NotificationTypeChat,
			Header:  "New Direct Messages",
			Content: message,
			Link:    link,
		}
		if err := app.Repositories.Notifications.Create(notification); err != nil {
			return fmt.Errorf("failed to create chat notification for user %d: %w", user.Id, err)
		}

		logger.Info(
			"Created unread chat notification",
			"user_id", user.Id, "username", user.Name, "message", message,
		)
	}
	return nil
}

type UnreadDmEntry struct {
	userId int
	name   string
	count  int
}

func generateUnreadDmNotification(app *state.State, userId int) (string, string, error) {
	unreadBySender, err := app.Repositories.Messages.FetchDMsUnreadCountAll(userId)
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch unread dm counts: %w", err)
	}
	entries := make([]UnreadDmEntry, 0, len(unreadBySender))
	totalMessages := 0

	for senderId, count := range unreadBySender {
		username, err := app.Repositories.Users.GetUsername(senderId)
		if err != nil {
			return "", "", fmt.Errorf("failed to fetch username for sender %d: %w", senderId, err)
		}
		entries = append(entries, UnreadDmEntry{
			userId: senderId,
			name:   username,
			count:  count,
		})
		totalMessages += count
	}

	if totalMessages <= 0 || len(entries) == 0 {
		return "", "", nil
	}

	// Sort users by count of unread messages in descending order
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].count > entries[j].count
	})
	link := fmt.Sprintf("/account/chat?target=%d", entries[0].userId)

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.name)
	}

	if len(entries) == 1 {
		suffix := "s"
		if totalMessages == 1 {
			suffix = ""
		}
		msg := fmt.Sprintf(
			"You have %d unread message%s from %s",
			totalMessages, suffix, names[0],
		)
		return msg, link, nil
	}

	msg := fmt.Sprintf(
		"You have %d unread messages from %s",
		totalMessages, joinNamesWithAnd(names),
	)
	return msg, link, nil
}

func joinNamesWithAnd(values []string) string {
	if len(values) <= 1 {
		// e.g. "Alice"
		return strings.Join(values, "")
	}
	if len(values) == 2 {
		// e.g. "Alice and Bob"
		return values[0] + " and " + values[1]
	}
	// e.g. "Alice, Bob, and Charlie"
	firstNames := values[:len(values)-1]
	lastName := values[len(values)-1]
	return strings.Join(firstNames, ", ") + " and " + lastName
}
