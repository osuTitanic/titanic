//go:build integration

package main

import (
	"context"
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
		&schemas.Group{},
		&schemas.GroupEntry{},
		&schemas.UserPermission{},
		&schemas.GroupPermission{},
		&schemas.Notification{},
		&schemas.Forum{},
		&schemas.ForumIcon{},
		&schemas.ForumTopic{},
		&schemas.ForumPost{},
		&schemas.Message{},
		&schemas.Beatmapset{},
		&schemas.Beatmap{},
		&schemas.Score{},
		&schemas.Release{},
	))
	setWebsiteStats(t, app)

	tests := []struct {
		name string
		path string
	}{
		{name: "home", path: "/"},
		{name: "download", path: "/download"},
		{name: "login", path: "/account/login"},
		{name: "register", path: "/account/register"},
		// TODO: Add more pages to test, once we have a nice way of creating test data in the db
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

func setWebsiteStats(t *testing.T, app *state.State) {
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
