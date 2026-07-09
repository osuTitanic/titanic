package tasks

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/internal/state"
)

// Retain activity metrics for 7 days
const ActivityRetentionDuration = 7 * 24 * time.Hour

// UpdateActivityStats stores active user counts in history for the website.
func UpdateActivityStats(app *state.State, logger *slog.Logger) error {
	ctx := context.Background()

	values, err := app.Redis.MGet(ctx,
		"bancho:activity:osu",
		"bancho:activity:irc",
		"bancho:activity:mp",
	).Result()
	if err != nil {
		return fmt.Errorf("failed to fetch activity from redis: %w", err)
	}

	osuCount, err := redisValueAsInt(values, 0)
	if err != nil {
		return fmt.Errorf("failed to parse osu activity count: %w", err)
	}

	ircCount, err := redisValueAsInt(values, 1)
	if err != nil {
		return fmt.Errorf("failed to parse irc activity count: %w", err)
	}

	mpCount, err := redisValueAsInt(values, 2)
	if err != nil {
		return fmt.Errorf("failed to parse mp activity count: %w", err)
	}

	entry := &schemas.BanchoActivity{
		OsuCount: osuCount,
		IrcCount: ircCount,
		MpCount:  mpCount,
	}
	if err := app.Repositories.Activity.Create(entry); err != nil {
		return fmt.Errorf("failed to create usercount history entry: %w", err)
	}
	logger.Info("Created usercount entry.", "players", osuCount+ircCount+mpCount)

	rows, err := app.Repositories.Activity.DeleteOlderThan(time.Now().Add(-ActivityRetentionDuration))
	if err != nil {
		return fmt.Errorf("failed to delete old usercount history entries: %w", err)
	}
	if rows > 0 {
		logger.Info("Deleted old usercount entries.", "rows_affected", rows)
	}

	return nil
}

func redisValueAsInt(values []any, idx int) (int, error) {
	if idx >= len(values) || values[idx] == nil {
		return 0, nil
	}

	strValue, ok := values[idx].(string)
	if !ok {
		return 0, fmt.Errorf("unexpected redis value type %T", values[idx])
	}

	parsed, err := strconv.Atoi(strValue)
	if err != nil {
		return 0, err
	}
	if parsed < 0 {
		return 0, nil
	}

	return parsed, nil
}
