package tasks

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/osuTitanic/titanic/internal/state"
)

const LoginRetention = 90 * 24 * time.Hour

// AnonymizeOldLogins removes IP addresses from logins older than 90 days
func AnonymizeOldLogins(app *state.State, logger *slog.Logger) error {
	cutoff := time.Now().UTC().Add(-LoginRetention)

	rows, err := app.Repositories.Logins.AnonymizeOlderThan(cutoff)
	if err != nil {
		return fmt.Errorf("failed to anonymize old logins: %w", err)
	}

	logger.Info(
		"Anonymized old logins",
		"rows_affected", rows, "cutoff", cutoff,
	)
	return nil
}
