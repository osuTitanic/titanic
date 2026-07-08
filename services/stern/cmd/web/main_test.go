//go:build integration

package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/osuTitanic/titanic-go/internal/schemas"
	"github.com/osuTitanic/titanic-go/internal/state"
	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
)

func TestPublicPagesRender(t *testing.T) {
	app := state.NewTestState(t, state.WithTestMigrations(
		&schemas.User{},
		&schemas.Stats{},
		&schemas.Login{},
		&schemas.Group{},
		&schemas.GroupEntry{},
		&schemas.UserPermission{},
		&schemas.GroupPermission{},
		&schemas.Notification{},
		&schemas.Forum{},
		&schemas.ForumIcon{},
		&schemas.ForumTopic{},
		&schemas.ForumPost{},
		&schemas.ForumBookmark{},
		&schemas.Message{},
		&schemas.Beatmapset{},
		&schemas.Beatmap{},
		&schemas.Score{},
		&schemas.Release{},
	))
	populateWebsiteStats(t, app)

	fixtures := state.NewTestData(t, app)
	_, forum, _ := populateWebsiteRenderData(t, fixtures)

	tests := []struct {
		name string
		path string
	}{
		{name: "home", path: "/"},
		{name: "download", path: "/download"},
		{name: "login", path: "/account/login"},
		{name: "register", path: "/account/register"},
		{name: "forum home", path: "/forum"},
		{name: "forum view", path: fmt.Sprintf("/forum/%d", forum.Id)},
	}
	router := newTestRouter(t, app)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, test.path, nil)
			request.Header.Set("User-Agent", "Mozilla/5.0")
			router.ServeHTTP(recorder, request)

			if recorder.Code != http.StatusOK {
				t.Fatalf("GET %s -> %d, want %d\n%s", test.path, recorder.Code, http.StatusOK, recorder.Body.String())
			}
			t.Logf("GET %s -> %d, body length: %d", test.path, recorder.Code, len(recorder.Body.String()))
		})
	}
}

func TestAuthenticatedPagesRender(t *testing.T) {
	app := state.NewTestState(t, state.WithTestMigrations(
		&schemas.User{},
		&schemas.Stats{},
		&schemas.Login{},
		&schemas.Group{},
		&schemas.GroupEntry{},
		&schemas.UserPermission{},
		&schemas.GroupPermission{},
		&schemas.Notification{},
		&schemas.Forum{},
		&schemas.ForumIcon{},
		&schemas.ForumTopic{},
		&schemas.ForumPost{},
		&schemas.ForumBookmark{},
		&schemas.Message{},
		&schemas.Beatmapset{},
		&schemas.Beatmap{},
		&schemas.Score{},
		&schemas.Release{},
	))
	populateWebsiteStats(t, app)

	fixtures := state.NewTestData(t, app)
	user, _, topic := populateWebsiteRenderData(t, fixtures)
	fixtures.CreateNotification(user)
	fixtures.CreateForumBookmark(user, topic)
	fixtures.CreateLogin(user)

	tests := []struct {
		name string
		path string
	}{
		{name: "account overview", path: "/account"},
		{name: "account profile", path: "/account/profile"},
		{name: "account security", path: "/account/security"},
		{name: "account chat", path: "/account/chat"},
	}
	router := newTestRouter(t, app)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, test.path, nil)
			request.Header.Set("User-Agent", "Mozilla/5.0")
			request.AddCookie(fixtures.CreateWebsiteSessionCookie(user, request))
			router.ServeHTTP(recorder, request)

			if recorder.Code != http.StatusOK {
				t.Fatalf("GET %s -> %d, want %d\n%s", test.path, recorder.Code, http.StatusOK, recorder.Body.String())
			}
			t.Logf("GET %s -> %d, body length: %d", test.path, recorder.Code, len(recorder.Body.String()))
		})
	}
}

func newTestRouter(t *testing.T, app *state.State) http.Handler {
	t.Helper()

	engine, err := templates.NewEngine(app.Config)
	if err != nil {
		t.Fatalf("failed to initialize templates: %v", err)
	}

	server := server.NewServer("localhost", 0, "stern-test", app, engine)
	InitializeWebRoutes(server)
	return server.Router
}

func populateWebsiteStats(t *testing.T, app *state.State) {
	t.Helper()

	if err := app.Redis.MSet(context.Background(),
		"bancho:totalusers", "0",
		"bancho:totalscores", "0",
		"bancho:activity:osu", "0",
		"bancho:activity:irc", "0",
	).Err(); err != nil {
		t.Fatalf("failed to seed redis stats: %v", err)
	}
}

func populateWebsiteRenderData(t *testing.T, fixtures *state.TestData) (*schemas.User, *schemas.Forum, *schemas.ForumTopic) {
	t.Helper()

	user := fixtures.CreateUser()
	fixtures.CreateStats(user)
	fixtures.CreateRelease()
	fixtures.CreateMessage(func(message *schemas.Message) {
		message.Sender = user.Name
	})

	mainForum := fixtures.CreateForum(func(forum *schemas.Forum) {
		forum.Name = "Announcements"
	})
	subForum := fixtures.CreateForum(func(forum *schemas.Forum) {
		parentId := mainForum.Id
		forum.ParentId = &parentId
		forum.Name = "Development"
	})
	topic := fixtures.CreateForumTopic(subForum, user, func(topic *schemas.ForumTopic) {
		topic.Announcement = true
		topic.Title = "i just farted"
	})
	fixtures.CreateForumPost(topic, user, func(post *schemas.ForumPost) {
		post.Content = "idk what to write here"
	})
	return user, subForum, topic
}
