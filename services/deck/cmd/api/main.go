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
	server.Handle("GET /web/maps/{query}", routes.BeatmapFile)
	server.Handle("GET /web/bancho_connect.php", routes.BanchoConnect)
	server.Handle("GET /web/check-updates.php", routes.CheckUpdates)
	server.Handle("GET /web/coins.php", routes.Coins)
	server.Handle("GET /web/osu-addfavourite.php", routes.AddFavourite)
	server.Handle("POST /web/osu-comment.php", routes.Comments)
	server.Handle("GET /web/osu-checktweets.php", routes.CheckTweets)
	server.Handle("POST /web/osu-getbeatmapinfo.php", routes.BeatmapInfo)
	server.Handle("GET /web/osu-getfavourites.php", routes.GetFavourites)
	server.Handle("GET /web/osu-getfriends.php", routes.Friends)
	server.Handle("GET /web/osu-getseasonal.php", routes.SeasonalBackgrounds)
	server.Handle("GET /web/osu-getstatus.php", routes.BeatmapStatus)
	server.Handle("GET /web/osu-login.php", routes.LegacyLogin)
	server.Handle("GET /web/osu-markasread.php", routes.MarkAsRead)
	server.Handle("GET /web/osu-search.php", routes.DirectSearch)
	server.Handle("POST /web/osu-screenshot.php", routes.Screenshot)
	server.Handle("GET /web/osu-stat.php", routes.UserStats)
	server.Handle("GET /web/osu-statoth.php", routes.UserStatsOther)
	server.Handle("GET /web/osu-title-image.php", routes.TitleImage)
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
		app.Config.WebHost,
		app.Config.WebPort,
		"deck", app,
	)
	InitializeRoutes(deck)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := deck.Serve(ctx); err != nil {
		slog.Error("HTTP server stopped unexpectedly", "error", err)
	}
}
