package main

import (
	"log/slog"
	"os"

	"github.com/osuTitanic/titanic-go/internal/state"
	"github.com/osuTitanic/titanic-go/services/stern/internal/routes"
	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
	web "github.com/osuTitanic/titanic-go/services/stern/web"
)

func InitializeWebRoutes(server *server.Server) {
	server.Handle("GET /{$}", routes.Home)
	server.Handle("GET /partials/home/news", routes.HomeNewsPartial)
	server.Handle("GET /partials/home/chat", routes.HomeChatPartial)
	server.Handle("GET /partials/home/plays", routes.HomePlaysPartial)
	server.Handle("GET /partials/packs/{id}", routes.BeatmapPackInfo)
	server.Handle("GET /account/login", routes.AccountLoginPage)
	server.Handle("POST /account/login", routes.AccountLogin)
	server.Handle("POST /account/logout", routes.AccountLogout)
	server.Handle("GET /account/register", routes.AccountRegisterPage)
	server.Handle("POST /account/register", routes.AccountRegister)
	server.Handle("GET /account/register/check", routes.AccountRegisterCheck)
	server.Handle("GET /account/verification", routes.AccountVerification)
	server.Handle("GET /account/verification/resend", routes.AccountVerificationResend)
	server.Handle("GET /account/reset", routes.PasswordResetPage)
	server.Handle("POST /account/reset", routes.PasswordReset)
	server.Handle("GET /download", routes.Download)
	server.Handle("GET /download/{$}", routes.Download)
	server.Handle("GET /beatmapsets", routes.Search)
	server.Handle("GET /beatmapsets/", routes.Search)
	server.Handle("GET /beatmapsets/packs", routes.Search)
	server.Handle("GET /beatmapsets/packs/", routes.BeatmapPacks)
	server.Handle("GET /", routes.NotFound)
}

func InitializeStaticRoutes(server *server.Server) {
	css, err := web.StaticFS("/css")
	if err != nil {
		slog.Error("Failed to initialize static file system", "error", err)
		os.Exit(1)
	}
	js, err := web.StaticFS("/js")
	if err != nil {
		slog.Error("Failed to initialize static file system", "error", err)
		os.Exit(1)
	}
	images, err := web.StaticFS("/images")
	if err != nil {
		slog.Error("Failed to initialize static file system", "error", err)
		os.Exit(1)
	}
	lib, err := web.StaticFS("/lib")
	if err != nil {
		slog.Error("Failed to initialize static file system", "error", err)
		os.Exit(1)
	}
	webfonts, err := web.StaticFS("/webfonts")
	if err != nil {
		slog.Error("Failed to initialize static file system", "error", err)
		os.Exit(1)
	}
	robots, err := web.StaticFS("/robots.txt")
	if err != nil {
		slog.Error("Failed to initialize static file system", "error", err)
		os.Exit(1)
	}
	favicon, err := web.StaticFS("/favicon.ico")
	if err != nil {
		slog.Error("Failed to initialize static file system", "error", err)
		os.Exit(1)
	}

	server.HandleFileSystem("/css/", css)
	server.HandleFileSystem("/js/", js)
	server.HandleFileSystem("/images/", images)
	server.HandleFileSystem("/lib/", lib)
	server.HandleFileSystem("/webfonts/", webfonts)
	server.HandleFileSystem("/favicon.ico", favicon)
	server.HandleFileSystem("/robots.txt", robots)
}

func main() {
	app, err := state.NewState(".env")
	if err != nil {
		slog.Error("Failed to initialize application", "error", err)
		os.Exit(1)
	}
	defer app.Close()

	engine, err := templates.NewEngine(app.Config)
	if err != nil {
		slog.Error("Failed to initialize templates", "error", err)
		os.Exit(1)
	}

	server := server.NewServer(app.Config.FrontendHost, app.Config.FrontendPort, "stern", app, engine)
	InitializeWebRoutes(server)
	InitializeStaticRoutes(server)
	server.Serve()
}
