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
	server.Handle("GET /a/", routes.DefaultAvatar)
	server.Handle("GET /a/{filename}", routes.Avatar)
	server.Handle("GET /mt/{filename}", routes.BeatmapThumbnail)
	server.Handle("GET /thumb/{filename}", routes.BeatmapThumbnail)
	server.Handle("GET /images/map-thumb/{filename}", routes.BeatmapThumbnail)
	server.Handle("GET /preview/{filename}", routes.BeatmapAudioPreview)
	server.Handle("GET /mp3/preview/{filename}", routes.BeatmapAudioPreview)
	server.Handle("GET /d/{filename}", routes.BeatmapDownload)
	server.Handle("GET /beatmapsets/{filename}/download", routes.BeatmapDownload)
	server.Handle("GET /ss/{id}", routes.ScreenshotRedirect)
	server.Handle("GET /ss/{id}/{checksum}", routes.Screenshot)

	server.Handle("GET /{$}", routes.Home)
	server.Handle("GET /partials/home/news", routes.HomeNewsPartial)
	server.Handle("GET /partials/home/chat", routes.HomeChatPartial)
	server.Handle("GET /partials/home/plays", routes.HomePlaysPartial)
	server.Handle("GET /partials/packs/{id}", routes.BeatmapPackInfo)
	server.Handle("GET /partials/users/{id}/general", routes.UserGeneralPartial)
	server.Handle("GET /partials/users/{id}/activity", routes.UserActivityPartial)
	server.Handle("GET /partials/users/{id}/leader", routes.UserTopPlaysPartial)
	server.Handle("GET /partials/users/{id}/scores", routes.UserScoresPartial)
	server.Handle("GET /partials/users/{id}/history", routes.UserHistoryPartial)
	server.Handle("GET /partials/users/{id}/beatmaps", routes.UserBeatmapsPartial)
	server.Handle("GET /partials/users/{id}/kudosu", routes.UserKudosuPartial)
	server.Handle("GET /partials/users/{id}/achievements", routes.UserAchievementsPartial)
	server.Handle("GET /charts/activity", routes.ActivityChart)
	server.Handle("GET /account", routes.AccountOverview)
	server.Handle("GET /account/{$}", routes.AccountOverview)
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
	server.Handle("GET /account/profile", routes.AccountProfile)
	server.Handle("POST /account/profile", routes.AccountProfileUpdate)
	server.Handle("POST /account/profile/userpage", routes.AccountProfileUserpage)
	server.Handle("POST /account/profile/signature", routes.AccountProfileSignature)
	server.Handle("GET /account/profile/avatar", routes.AccountAvatarRedirect)
	server.Handle("POST /account/profile/avatar", routes.AccountAvatarUpdate)
	server.Handle("GET /account/security", routes.AccountSecurity)
	server.Handle("POST /account/security", routes.AccountSecurityUpdate)
	server.Handle("GET /account/friends", routes.AccountFriends)
	server.Handle("GET /account/chat", routes.AccountChat)
	server.Handle("GET /download", routes.Download)
	server.Handle("GET /download/{$}", routes.Download)
	server.Handle("GET /p/doyoureallywanttoaskpeppy", routes.PeppySkillIssue)
	server.Handle("GET /p/doyoureallywanttoaskpeppy/", routes.PeppySkillIssue)
	server.Handle("GET /events", routes.Events)
	server.Handle("GET /p/changelog", routes.OsuChangelog)
	server.Handle("GET /u/{query}", routes.UserProfile)
	server.Handle("GET /users/{query}", routes.UserProfileRedirect)
	server.Handle("GET /g/{id}", routes.Group)
	server.Handle("GET /forum", routes.ForumHome)
	server.Handle("GET /forum/{$}", routes.ForumHome)
	server.Handle("GET /forum/{id}", routes.ForumView)
	server.Handle("GET /forum/{id}/create", routes.ForumCreateTopicView)
	server.Handle("POST /forum/{id}/create", routes.ForumCreateTopicAction)
	server.Handle("GET /forum/{id}/t/{topicId}", routes.ForumTopicView)
	server.Handle("GET /forum/{id}/t/{topicId}/", routes.ForumTopicView)
	server.Handle("GET /forum/{id}/t/{topicId}/post", routes.ForumPostEditorView)
	server.Handle("GET /forum/{id}/t/{topicId}/post/", routes.ForumPostEditorView)
	server.Handle("POST /forum/{id}/t/{topicId}/post", routes.ForumPostAction)
	server.Handle("POST /forum/{id}/t/{topicId}/draft", routes.ForumDraftAction)
	server.Handle("GET /forum/{id}/t/{topicId}/p/{postId}", routes.ForumPostPermalink)
	server.Handle("GET /forum/{id}/t/{topicId}/p/{postId}/", routes.ForumPostPermalink)
	server.Handle("GET /forum/{id}/p/{postId}", routes.ForumPostRedirect)
	server.Handle("GET /forum/{id}/p/{postId}/", routes.ForumPostRedirect)
	server.Handle("GET /forum/{kind}/{rest...}", routes.ForumShortlinkRedirect)
	server.Handle("GET /forum/posting.php", routes.ForumQuickReplyRedirect)
	server.Handle("GET /forum/viewtopic.php", routes.ForumViewTopicRedirect)
	server.Handle("GET /forum/viewforum.php", routes.ForumViewForumRedirect)
	server.Handle("GET /forum/ucp.php", routes.ForumControlPanelRedirect)
	server.Handle("GET /forum/index.php", routes.ForumIndexRedirect)
	server.Handle("GET /s/{id}", routes.BeatmapsetRedirect)
	server.Handle("GET /mp/{id}", routes.Match)
	server.Handle("GET /mp/{id}/{$}", routes.Match)
	server.Handle("GET /b/{id}", routes.Beatmap)
	server.Handle("GET /scores/{id}", routes.Score)
	server.Handle("GET /scores/{id}/{$}", routes.Score)
	server.Handle("GET /scores/{id}/download", routes.ScoreReplayDownload)
	server.Handle("GET /beatmaps/{id}", routes.BeatmapRedirect)
	server.Handle("GET /beatmapsets", routes.Search)
	server.Handle("GET /beatmapsets/", routes.Search)
	server.Handle("GET /beatmapsets/packs", routes.Search)
	server.Handle("GET /beatmapsets/packs/{$}", routes.BeatmapPacks)
	server.Handle("GET /rankings/kudosu", routes.RankingsKudosu)
	server.Handle("GET /rankings/{mode}/country", routes.RankingsCountry)
	server.Handle("GET /rankings/{mode}/{rankingType}", routes.RankingsGlobal)
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
