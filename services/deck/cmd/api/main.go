package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/osuTitanic/titanic/internal/state"
	"github.com/osuTitanic/titanic/services/deck/internal/routes"
	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

func InitializeRoutes(server *server.Server) {
	server.Handle("GET /web/check-updates.php", routes.CheckUpdates)
}

func main() {
	// TODO: Healthcheck

	app, err := state.NewState(".env")
	if err != nil {
		slog.Error("Failed to initialize application", "error", err)
		os.Exit(1)
	}
	defer app.Close()

	deck := server.NewServer(
		app.Config.ApiHost,
		app.Config.ApiPort,
		"deck", app,
	)
	InitializeRoutes(deck)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := deck.Serve(ctx); err != nil {
		slog.Error("HTTP server stopped unexpectedly", "error", err)
	}
}
