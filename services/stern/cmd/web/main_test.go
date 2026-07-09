//go:build integration

package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/osuTitanic/titanic-go/internal/constants"
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"github.com/osuTitanic/titanic-go/internal/state"
	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
	"github.com/redis/go-redis/v9"
)

type websiteRouteTest struct {
	name string
	path string
}

type websiteRenderData struct {
	user       *schemas.User
	friend     *schemas.User
	group      *schemas.Group
	forum      *schemas.Forum
	topic      *schemas.ForumTopic
	post       *schemas.ForumPost
	beatmapset *schemas.Beatmapset
	beatmap    *schemas.Beatmap
	score      *schemas.Score
	pack       *schemas.BeatmapPack
}

func TestWebsiteRoutesRender(t *testing.T) {
	app := newWebsiteTestState(t)
	router := newTestRouter(t, app)
	populateWebsiteStats(t, app)

	fixtures := state.NewTestData(t, app)
	data := populateWebsiteData(t, app, fixtures)

	publicRoutes := []websiteRouteTest{
		{name: "home", path: "/"},
		{name: "home news partial", path: "/partials/home/news"},
		{name: "home chat partial", path: "/partials/home/chat"},
		{name: "home plays partial", path: "/partials/home/plays"},
		{name: "download", path: "/download"},
		{name: "login", path: "/account/login"},
		{name: "register", path: "/account/register"},
		{name: "reset", path: "/account/reset"},
		{name: "events", path: "/events"},
		{name: "peppy", path: "/p/doyoureallywanttoaskpeppy"},
		{name: "forum home", path: "/forum"},
		{name: "forum view", path: fmt.Sprintf("/forum/%d", data.forum.Id)},
		{name: "forum topic", path: fmt.Sprintf("/forum/%d/t/%d", data.forum.Id, data.topic.Id)},
		{name: "group", path: fmt.Sprintf("/g/%d", data.group.Id)},
		{name: "user profile", path: fmt.Sprintf("/u/%d", data.user.Id)},
		{name: "beatmap search", path: "/beatmapsets"},
		{name: "beatmap packs", path: "/beatmapsets/packs"},
		{name: "beatmap pack partial", path: fmt.Sprintf("/partials/packs/%d", data.pack.Id)},
		{name: "beatmap", path: fmt.Sprintf("/b/%d", data.beatmap.Id)},
		{name: "score", path: fmt.Sprintf("/scores/%d", data.score.Id)},
		{name: "rankings performance", path: "/rankings/osu/performance"},
		{name: "rankings country", path: "/rankings/osu/country"},
		{name: "rankings kudosu", path: "/rankings/kudosu"},
	}

	t.Run("public", func(t *testing.T) {
		assertWebsiteRoutesRender(t, router, publicRoutes, nil)
	})

	fixtures.CreateNotification(data.user)
	fixtures.CreateForumBookmark(data.user, data.topic)
	fixtures.CreateLogin(data.user)

	authenticatedRoutes := []websiteRouteTest{
		{name: "account overview", path: "/account"},
		{name: "account profile", path: "/account/profile"},
		{name: "account security", path: "/account/security"},
		{name: "account friends", path: "/account/friends"},
		{name: "account chat", path: "/account/chat"},
		{name: "authenticated forum topic", path: fmt.Sprintf("/forum/%d/t/%d", data.forum.Id, data.topic.Id)},
		{name: "authenticated beatmap", path: fmt.Sprintf("/b/%d", data.beatmap.Id)},
	}

	t.Run("authenticated", func(t *testing.T) {
		assertWebsiteRoutesRender(t, router, authenticatedRoutes, func(request *http.Request) {
			request.AddCookie(fixtures.CreateWebsiteSessionCookie(data.user, request))
		})
	})

	partialRoutes := []websiteRouteTest{
		{name: "general", path: fmt.Sprintf("/partials/users/%d/general", data.user.Id)},
		{name: "activity", path: fmt.Sprintf("/partials/users/%d/activity", data.user.Id)},
		{name: "leader", path: fmt.Sprintf("/partials/users/%d/leader", data.user.Id)},
		{name: "best scores", path: fmt.Sprintf("/partials/users/%d/scores?section=best", data.user.Id)},
		{name: "first scores", path: fmt.Sprintf("/partials/users/%d/scores?section=first", data.user.Id)},
		{name: "history", path: fmt.Sprintf("/partials/users/%d/history", data.user.Id)},
		{name: "beatmaps", path: fmt.Sprintf("/partials/users/%d/beatmaps", data.user.Id)},
		{name: "kudosu", path: fmt.Sprintf("/partials/users/%d/kudosu", data.user.Id)},
		{name: "achievements", path: fmt.Sprintf("/partials/users/%d/achievements", data.user.Id)},
	}

	t.Run("partials", func(t *testing.T) {
		assertWebsiteRoutesRender(t, router, partialRoutes, nil)
	})
}

func newWebsiteTestState(t *testing.T) *state.State {
	t.Helper()

	return state.NewTestState(t, state.WithTestMigrations(
		&schemas.User{},
		&schemas.Stats{},
		&schemas.Login{},
		&schemas.Group{},
		&schemas.GroupEntry{},
		&schemas.UserPermission{},
		&schemas.GroupPermission{},
		&schemas.Notification{},
		&schemas.Relationship{},
		&schemas.Badge{},
		&schemas.Stamp{},
		&schemas.Name{},
		&schemas.Activity{},
		&schemas.Achievement{},
		&schemas.Forum{},
		&schemas.ForumIcon{},
		&schemas.ForumTopic{},
		&schemas.ForumPost{},
		&schemas.ForumBookmark{},
		&schemas.ForumSubscriber{},
		&schemas.Message{},
		&schemas.Beatmapset{},
		&schemas.Beatmap{},
		&schemas.Score{},
		&schemas.BeatmapFavourite{},
		&schemas.BeatmapCollaboration{},
		&schemas.BeatmapNomination{},
		&schemas.BeatmapModding{},
		&schemas.BeatmapPack{},
		&schemas.BeatmapPackEntry{},
		&schemas.BeatmapPlays{},
		&schemas.Release{},
	))
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

func assertWebsiteRoutesRender(t *testing.T, router http.Handler, tests []websiteRouteTest, prepare func(*http.Request)) {
	t.Helper()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, test.path, nil)
			request.Header.Set("User-Agent", "Mozilla/5.0")
			if prepare != nil {
				prepare(request)
			}
			router.ServeHTTP(recorder, request)

			if recorder.Code != http.StatusOK {
				t.Fatalf("GET %s -> %d, want %d\n%s", test.path, recorder.Code, http.StatusOK, recorder.Body.String())
			}
			t.Logf("GET %s -> %d, body length: %d", test.path, recorder.Code, len(recorder.Body.String()))
		})
	}
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

func populateWebsiteData(t *testing.T, app *state.State, fixtures *state.TestData) *websiteRenderData {
	t.Helper()

	user := fixtures.CreateUser(func(user *schemas.User) {
		user.Name = "TestPlayer"
		user.SafeName = "testplayer"
		user.Email = "testplayer@example.com"
		user.Country = "US"
	})
	stats := fixtures.CreateStats(user)
	friend := fixtures.CreateUser(func(user *schemas.User) {
		user.Name = "TestFriend"
		user.SafeName = "testfriend"
		user.Email = "testfriend@example.com"
		user.Country = "CA"
	})
	friendStats := fixtures.CreateStats(friend, func(stats *schemas.Stats) {
		stats.PP = 900.5
		stats.Rscore = 654321
		stats.Tscore = 98765432
	})
	fixtures.CreateRelease()
	fixtures.CreateMessage(func(message *schemas.Message) {
		message.Sender = user.Name
	})
	fixtures.CreateRelationship(user, friend)
	fixtures.CreateRelationship(friend, user)

	group := fixtures.CreateGroup(func(group *schemas.Group) {
		group.Name = "Testers"
		group.ShortName = "TST"
	})
	fixtures.CreateGroupEntry(group, user)

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
		topic.Title = "Website integration test"
	})
	post := fixtures.CreateForumPost(topic, user, func(post *schemas.ForumPost) {
		post.Content = "Lorem ipsum dollar sitting america"
	})

	beatmapset := fixtures.CreateBeatmapset(user, func(beatmapset *schemas.Beatmapset) {
		beatmapset.TopicId = &topic.Id
	})
	beatmap := fixtures.CreateBeatmap(beatmapset)
	score := fixtures.CreateScore(user, beatmap, func(score *schemas.Score) {
		score.Pinned = true
		score.TotalScore = 9876543
	})
	pack := fixtures.CreateBeatmapPack(user)
	fixtures.CreateBeatmapPackEntry(pack, beatmapset)
	fixtures.CreateBeatmapModding(user, friend, beatmapset, post)

	updateWebsiteRankings(t, app, stats, user.Country)
	updateWebsiteRankings(t, app, friendStats, friend.Country)
	populateWebsiteRankingCountries(t, app)
	if err := app.Rankings.UpdateKudosu(user.Id, user.Country, app.Repositories.Modding); err != nil {
		t.Fatalf("failed to seed kudosu rankings: %v", err)
	}

	return &websiteRenderData{
		user:       user,
		friend:     friend,
		group:      group,
		forum:      subForum,
		topic:      topic,
		post:       post,
		beatmapset: beatmapset,
		beatmap:    beatmap,
		score:      score,
		pack:       pack,
	}
}

func updateWebsiteRankings(t *testing.T, app *state.State, stats *schemas.Stats, country string) {
	t.Helper()

	if err := app.Rankings.Update(stats, country); err != nil {
		t.Fatalf("failed to seed rankings: %v", err)
	}
	if err := app.Rankings.UpdateLeaderScores(stats, country, app.Repositories.Scores); err != nil {
		t.Fatalf("failed to seed leader rankings: %v", err)
	}
}

func populateWebsiteRankingCountries(t *testing.T, app *state.State) {
	t.Helper()

	pipe := app.Redis.Pipeline()
	seeded := 0
	for _, country := range constants.CountryCodes {
		if country == "XX" || country == "US" || country == "CA" {
			continue
		}

		// Seed some rankings for the country to ensure it appears in the country selector
		// They aren't real users, but they will be enough to make the country appear in the selector

		score := float64(1000 - seeded)
		member := 100000 + seeded
		country := country

		for _, rankingType := range []string{"performance", "rscore", "tscore"} {
			pipe.ZAdd(context.Background(), app.Rankings.RankingKey(constants.ModeOsu, rankingType, &country), redis.Z{
				Score:  score,
				Member: member,
			})
		}

		seeded++
		if seeded == 20 {
			// We only need to seed 20 ish countries for the test
			break
		}
	}

	if _, err := pipe.Exec(context.Background()); err != nil {
		t.Fatalf("failed to seed ranking country selector: %v", err)
	}
}
