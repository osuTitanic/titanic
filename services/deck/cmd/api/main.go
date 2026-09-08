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
	server.Handle("GET /a/", routes.DefaultAvatar)
	server.Handle("GET /a/{filename}", routes.Avatar)
	server.Handle("GET /forum/download.php", routes.AvatarForum)
	server.Handle("GET /d/{filename}", routes.BeatmapDownload)
	server.Handle("GET /bss/{filename}", routes.BeatmapDownload)
	server.Handle("GET /preview/{filename}", routes.BeatmapAudioPreview)
	server.Handle("GET /mp3/preview/{filename}", routes.BeatmapAudioPreview)
	server.Handle("GET /mt/{filename}", routes.BeatmapThumbnail)
	server.Handle("GET /thumb/{filename}", routes.BeatmapThumbnail)
	server.Handle("GET /images/map-thumb/{filename}", routes.BeatmapThumbnail)
	server.Handle("GET /ss/{id}", routes.ScreenshotRedirect)
	server.Handle("GET /ss/{id}/{checksum}", routes.ScreenshotImage)
	server.Handle("GET /osu/{query}", routes.BeatmapFile)
	server.Handle("GET /assets/menu-content.json", routes.MenuContent)

	server.Handle("GET /release/filter.txt", routes.ReleaseFilter)
	server.Handle("GET /release/Localisation/{filename}", routes.ReleaseLocalisation)
	server.Handle("GET /release/update", routes.ReleaseUpdate)
	server.Handle("GET /release/patches.php", routes.ReleasePatches)
	server.Handle("GET /release/update.php", routes.ReleaseUpdateCheck)
	server.Handle("GET /release/update2.php", routes.ReleaseUpdateV2)
	server.Handle("GET /release/{language}/{filename}", routes.ReleaseLegacyLocalisation)
	server.Handle("GET /release/{filename}", routes.ReleaseFile)

	server.Handle("GET /web/maps/{query}", routes.BeatmapFile)
	server.Handle("GET /web/bancho_connect.php", routes.BanchoConnect)
	server.Handle("GET /web/check-updates.php", routes.CheckUpdates)
	server.Handle("GET /web/coins.php", routes.Coins)
	server.Handle("GET /rating/ingame-rate.php", routes.IngameRate)
	server.Handle("GET /rating/ingame-rate2.php", routes.IngameRate2)
	server.Handle("GET /web/osu-addfavourite.php", routes.AddFavourite)
	server.Handle("POST /web/osu-benchmark.php", routes.Benchmark)
	server.Handle("POST /web/osu-comment.php", routes.Comments)
	server.Handle("GET /web/osu-checktweets.php", routes.CheckTweets)
	server.Handle("POST /web/osu-error.php", routes.ErrorReport)
	server.Handle("POST /web/osu-getbeatmapinfo.php", routes.BeatmapInfo)
	server.Handle("GET /web/osu-getfavourites.php", routes.GetFavourites)
	server.Handle("GET /web/osu-getfriends.php", routes.Friends)
	server.Handle("GET /web/osu-getreplay.php", routes.Replay)
	server.Handle("GET /web/osu-getscores.php", routes.GetScores)
	server.Handle("GET /web/osu-getscores2.php", routes.GetScores2)
	server.Handle("GET /web/osu-getscores3.php", routes.GetScores3)
	server.Handle("GET /web/osu-getscores4.php", routes.GetScores4)
	server.Handle("GET /web/osu-getscores5.php", routes.GetScores5)
	server.Handle("GET /web/osu-getscores6.php", routes.GetScores6)
	server.Handle("GET /web/osu-osz2-getscores.php", routes.GetScoresOsz2)
	server.Handle("GET /web/osu-getseasonal.php", routes.SeasonalBackgrounds)
	server.Handle("GET /web/osu-getstatus.php", routes.BeatmapStatus)
	server.Handle("GET /web/osu-login.php", routes.LegacyLogin)
	server.Handle("GET /web/osu-markasread.php", routes.MarkAsRead)
	server.Handle("GET /web/osu-rate.php", routes.OsuRate)
	server.Handle("GET /web/osu-search.php", routes.DirectSearch)
	server.Handle("GET /web/osu-search-set.php", routes.DirectSearchSet)
	server.Handle("POST /web/osu-screenshot.php", routes.Screenshot)
	server.Handle("GET /web/osu-stat.php", routes.UserStats)
	server.Handle("GET /web/osu-statoth.php", routes.UserStatsOther)
	server.Handle("GET /web/osu-title-image.php", routes.MenuIcon)
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
