package main

import (
	"log/slog"
	"os"

	"github.com/osuTitanic/titanic/internal/state"
)

func main() {
	// TODO: Healthcheck

	app, err := state.NewState(".env")
	if err != nil {
		slog.Error("Failed to initialize application", "error", err)
		os.Exit(1)
	}
	defer app.Close()

	slog.Info("gaming")
	// TODO: HTTP server
}
