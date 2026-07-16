//go:build integration

package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/osuTitanic/titanic/internal/authentication"
	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/internal/state"
	"github.com/osuTitanic/titanic/services/stern/internal/server"
	"github.com/osuTitanic/titanic/services/stern/internal/templates"
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
		{name: "forum search", path: "/forum/search"},
		{name: "forum search filters", path: fmt.Sprintf("/forum/search?forum=%d&username=%s&sort=1&order=1", data.forum.Id, url.QueryEscape(data.user.Name))},
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

	t.Run("forum search filters", func(t *testing.T) {
		matchingPost := fixtures.CreateForumPost(data.topic, data.friend, func(post *schemas.ForumPost) {
			post.Content = "Details about digital client"
		})
		unrelatedPost := fixtures.CreateForumPost(data.topic, data.friend, func(post *schemas.ForumPost) {
			post.Content = "An unrelated reply"
		})

		body := renderForumSearch(t, router, url.Values{"username": {data.friend.Name}})
		assertForumSearchPost(t, body, matchingPost, true)
		assertForumSearchPost(t, body, unrelatedPost, true)

		body = renderForumSearch(t, router, url.Values{
			"username": {data.friend.Name},
			"query":    {"digital client"},
		})
		assertForumSearchPost(t, body, matchingPost, true)
		assertForumSearchPost(t, body, unrelatedPost, false)
		for _, term := range []string{"digital", "client"} {
			if !strings.Contains(body, "<strong>"+term+"</strong>") {
				t.Errorf("forum search does not highlight %q", term)
			}
		}

		body = renderForumSearch(t, router, url.Values{"query": {data.topic.Title}})
		assertForumSearchPost(t, body, data.post, false)

		body = renderForumSearch(t, router, url.Values{"query": {"lorem ipsum"}})
		assertForumSearchPost(t, body, data.post, true)
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

	t.Run("post flows", func(t *testing.T) {
		assertWebsitePostFlows(t, app, router, fixtures, data)
	})
}

func renderForumSearch(t *testing.T, router http.Handler, query url.Values) string {
	t.Helper()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/forum/search?"+query.Encode(), nil)
	request.Header.Set("User-Agent", "Mozilla/5.0")
	router.ServeHTTP(recorder, request)

	assertStatus(t, recorder, http.StatusOK)
	return recorder.Body.String()
}

func assertForumSearchPost(t *testing.T, body string, post *schemas.ForumPost, want bool) {
	t.Helper()

	postUrl := fmt.Sprintf("/forum/%d/p/%d/", post.ForumId, post.Id)
	if got := strings.Contains(body, postUrl); got != want {
		t.Errorf("forum search contains post %d = %t, want %t", post.Id, got, want)
	}
}

func newWebsiteTestState(t *testing.T) *state.State {
	t.Helper()

	app := state.NewTestState(t, state.WithTestMigrations(
		&schemas.User{},
		&schemas.Stats{},
		&schemas.Login{},
		&schemas.Group{},
		&schemas.GroupEntry{},
		&schemas.UserPermission{},
		&schemas.GroupPermission{},
		&schemas.Notification{},
		&schemas.Verification{},
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
	enableForumPostSearchVector(t, app)
	return app
}

func enableForumPostSearchVector(t *testing.T, app *state.State) {
	t.Helper()

	// Gorm migrations don't do this automatically
	statement := `
		ALTER TABLE forum_posts ADD COLUMN search_vector tsvector
		GENERATED ALWAYS AS (to_tsvector('english', coalesce(content, ''))) STORED
	`
	if err := app.Database.Exec(statement).Error; err != nil {
		t.Fatalf("failed to add forum post search vector: %v", err)
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

func assertWebsitePostFlows(t *testing.T, app *state.State, router http.Handler, fixtures *state.TestData, data *websiteRenderData) {
	t.Helper()

	grantWebsitePostPermissions(t, app, data.user)
	populateRegistrationGroups(t, app)

	t.Run("login", func(t *testing.T) {
		before := countRows(t, app, &schemas.Login{}, "user_id = ? AND osu_version = ?", data.user.Id, "web")
		recorder := postWebsiteForm(t, router, "/account/login", url.Values{
			"username": {data.user.Name},
			"password": {"password"},
			"redirect": {"/download"},
		}, nil)

		assertRedirect(t, recorder, "/download")
		assertCookieSet(t, recorder, authentication.WebsiteSessionCookieName)

		// Ensure a login row was created
		after := countRows(t, app, &schemas.Login{}, "user_id = ? AND osu_version = ?", data.user.Id, "web")
		if after != before+1 {
			t.Fatalf("login rows = %d, want %d", after, before+1)
		}
	})

	t.Run("logout", func(t *testing.T) {
		recorder := postWebsiteForm(
			t, router, "/account/logout",
			url.Values{"redirect": {"/download"}},
			authenticateWebsiteRequest(t, app, fixtures, data.user),
		)
		assertRedirect(t, recorder, "/download")
		assertCookieSet(t, recorder, authentication.WebsiteSessionCookieName)
	})

	t.Run("register", func(t *testing.T) {
		recorder := postWebsiteForm(t, router, "/account/register", url.Values{
			"username": {"Registered User"},
			"email":    {"registereduser@example.com"},
			"password": {"testpassword"},
		}, func(request *http.Request) {
			request.Header.Set("CF-IPCountry", "US")
		})

		assertRedirectPrefix(t, recorder, "/account/verification?id=")
		user, err := app.Users.BySafeName("registered_user")
		if err != nil {
			t.Fatalf("failed to fetch registered user: %v", err)
		}
		if user == nil {
			t.Fatal("registered user was not created")
		}
		if user.Activated {
			t.Fatal("registered user should require activation")
		}
		assertRowCount(t, app, &schemas.GroupEntry{}, 1, "user_id = ? AND group_id = ?", user.Id, constants.GroupPlayers)
		assertRowCount(t, app, &schemas.GroupEntry{}, 1, "user_id = ? AND group_id = ?", user.Id, constants.GroupSupporter)
		assertRowCount(t, app, &schemas.Verification{}, 1, "user_id = ? AND type = ?", user.Id, constants.VerificationTypeActivation)
	})

	t.Run("password reset", func(t *testing.T) {
		resetUser := fixtures.CreateUser(func(user *schemas.User) {
			user.Name = "ResetUser"
			user.SafeName = "resetuser"
			user.Email = "reset@example.com"
		})
		recorder := postWebsiteForm(t, router, "/account/reset", url.Values{
			"email": {resetUser.Email},
		}, nil)

		assertRedirectPrefix(t, recorder, "/account/verification?id=")
		verification := latestVerification(t, app, resetUser.Id, constants.VerificationTypePassword)

		recorder = postWebsiteForm(t, router, "/account/reset", url.Values{
			"token":          {verification.Token},
			"password":       {"new-reset-password"},
			"password_match": {"new-reset-password"},
		}, nil)
		assertStatus(t, recorder, http.StatusOK)

		updated, err := app.Users.ById(resetUser.Id)
		if err != nil {
			t.Fatalf("failed to fetch reset user: %v", err)
		}
		if !authentication.VerifyPasswordHash("new-reset-password", updated.Bcrypt) {
			t.Fatal("password reset did not update the user's password")
		}
		assertRowCount(t, app, &schemas.Verification{}, 0, "token = ?", verification.Token)
	})

	t.Run("profile update", func(t *testing.T) {
		recorder := postWebsiteForm(t, router, "/account/profile", url.Values{
			"mode":      {"1"},
			"interests": {"testing"},
			"location":  {"osu hq"},
			"website":   {"https://example.com"},
			"discord":   {"tester_123"},
			"twitter":   {"@tester"},
		}, authenticateWebsiteRequest(t, app, fixtures, data.user))
		assertStatus(t, recorder, http.StatusOK)

		updated, err := app.Users.ById(data.user.Id)
		if err != nil {
			t.Fatalf("failed to fetch updated profile: %v", err)
		}
		assertStringPointer(t, "interests", updated.Interests, "testing")
		assertStringPointer(t, "location", updated.Location, "osu hq")
		assertStringPointer(t, "website", updated.Website, "https://example.com")
		assertStringPointer(t, "discord", updated.Discord, "tester_123")
		assertStringPointer(t, "twitter", updated.Twitter, "https://twitter.com/@tester")
		if updated.PreferredMode != constants.ModeTaiko {
			t.Fatalf("preferred mode = %v, want %v", updated.PreferredMode, constants.ModeTaiko)
		}
	})

	t.Run("userpage", func(t *testing.T) {
		recorder := postWebsiteForm(t, router, "/account/profile/userpage", url.Values{
			"user_id": {strconv.Itoa(data.user.Id)},
			"bbcode":  {"it still kiiillsss meeeeee"},
		}, authenticateWebsiteRequest(t, app, fixtures, data.user))

		assertRedirect(t, recorder, "/account/profile#userpage")
		updated, err := app.Users.ById(data.user.Id)
		if err != nil {
			t.Fatalf("failed to fetch updated userpage: %v", err)
		}
		assertStringPointer(t, "userpage", updated.Userpage, "it still kiiillsss meeeeee")
	})

	t.Run("signature", func(t *testing.T) {
		recorder := postWebsiteForm(t, router, "/account/profile/signature", url.Values{
			"user_id": {strconv.Itoa(data.user.Id)},
			"bbcode":  {"Test signature."},
		}, authenticateWebsiteRequest(t, app, fixtures, data.user))

		assertRedirect(t, recorder, "/account/profile#signature")
		updated, err := app.Users.ById(data.user.Id)
		if err != nil {
			t.Fatalf("failed to fetch updated signature: %v", err)
		}
		assertStringPointer(t, "signature", updated.Signature, "Test signature.")
	})

	t.Run("avatar", func(t *testing.T) {
		recorder := postWebsiteMultipart(
			t, router, "/account/profile/avatar", "avatar", "avatar.png",
			testPng(t), authenticateWebsiteRequest(t, app, fixtures, data.user),
		)
		assertRedirect(t, recorder, "/account/profile")

		updated, err := app.Users.ById(data.user.Id)
		if err != nil {
			t.Fatalf("failed to fetch updated avatar user: %v", err)
		}
		if updated.AvatarHash == nil || *updated.AvatarHash == "" {
			t.Fatal("avatar hash was not updated")
		}
		if !app.Storage.Exists(strconv.Itoa(data.user.Id), "avatars") {
			t.Fatal("avatar file was not saved")
		}
	})

	t.Run("security password", func(t *testing.T) {
		recorder := postWebsiteForm(t, router, "/account/security", url.Values{
			"current-password": {"password"},
			"new-password":     {"new-account-password"},
			"password-confirm": {"new-account-password"},
		}, authenticateWebsiteRequest(t, app, fixtures, data.user))

		assertStatus(t, recorder, http.StatusOK)
		updated, err := app.Users.ById(data.user.Id)
		if err != nil {
			t.Fatalf("failed to fetch password-updated user: %v", err)
		}
		if !authentication.VerifyPasswordHash("new-account-password", updated.Bcrypt) {
			t.Fatal("security password update did not update the user's password")
		}
	})

	t.Run("security email", func(t *testing.T) {
		emailUser := fixtures.CreateUser(func(user *schemas.User) {
			user.Name = "EmailChangeUser"
			user.SafeName = "emailchangeuser"
			user.Email = "email-change@example.com"
		})
		recorder := postWebsiteForm(t, router, "/account/security", url.Values{
			"current-password": {"password"},
			"new-email":        {"changed@example.com"},
			"email-confirm":    {"changed@example.com"},
		}, authenticateWebsiteRequest(t, app, fixtures, emailUser))

		assertRedirectPrefix(t, recorder, "/account/verification?id=")
		updated, err := app.Users.ById(emailUser.Id)
		if err != nil {
			t.Fatalf("failed to fetch email-updated user: %v", err)
		}
		if updated.Email != "changed@example.com" {
			t.Fatalf("email = %q, want %q", updated.Email, "changed@example.com")
		}
		if updated.Activated {
			t.Fatal("email change should deactivate the user")
		}
		assertRowCount(t, app, &schemas.Verification{}, 1, "user_id = ? AND type = ?", emailUser.Id, constants.VerificationTypeActivation)
	})
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

func populateRegistrationGroups(t *testing.T, app *state.State) {
	t.Helper()

	groups := []*schemas.Group{
		{
			Id:        constants.GroupSupporter,
			Name:      "Supporter",
			ShortName: "SUP",
			Color:     "#ff66aa",
		},
		{
			Id:        constants.GroupPlayers,
			Name:      "Players",
			ShortName: "PLY",
			Color:     "#66aa66",
		},
	}
	for _, group := range groups {
		if err := app.Database.Create(group).Error; err != nil {
			t.Fatalf("failed to seed registration group %d: %v", group.Id, err)
		}
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

func grantWebsitePostPermissions(t *testing.T, app *state.State, user *schemas.User) {
	t.Helper()

	permissions := []string{
		"users.profile.update",
		"forum.topics.create",
		"forum.posts.create",
		"forum.posts.edit",
	}
	for _, permission := range permissions {
		err := app.Repositories.Permissions.CreateUserPermission(&schemas.UserPermission{
			UserId:     user.Id,
			Permission: permission,
		})
		if err != nil {
			t.Fatalf("failed to grant test permission %q: %v", permission, err)
		}
	}
}

func postWebsiteForm(t *testing.T, router http.Handler, path string, values url.Values, prepare func(*http.Request)) *httptest.ResponseRecorder {
	t.Helper()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, path, strings.NewReader(values.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("User-Agent", "Mozilla/5.0")
	if prepare != nil {
		prepare(request)
	}
	router.ServeHTTP(recorder, request)
	return recorder
}

func postWebsiteMultipart(t *testing.T, router http.Handler, path string, field string, filename string, data []byte, prepare func(*http.Request)) *httptest.ResponseRecorder {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile(field, filename)
	if err != nil {
		t.Fatalf("failed to create multipart file: %v", err)
	}
	if _, err := part.Write(data); err != nil {
		t.Fatalf("failed to write multipart file: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close multipart writer: %v", err)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, path, &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("User-Agent", "Mozilla/5.0")
	if prepare != nil {
		prepare(request)
	}
	router.ServeHTTP(recorder, request)
	return recorder
}

func authenticateWebsiteRequest(t *testing.T, app *state.State, fixtures *state.TestData, user *schemas.User) func(*http.Request) {
	t.Helper()

	return func(request *http.Request) {
		request.AddCookie(fixtures.CreateWebsiteSessionCookie(user, request))
		token, err := app.CSRFStore.Upsert(request.Context(), user.Id)
		if err != nil {
			t.Fatalf("failed to create csrf token: %v", err)
		}
		request.Header.Set("X-CSRF-Token", token)
	}
}

func testPng(t *testing.T) []byte {
	t.Helper()

	// Create a simple 2x2 square
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	for y := range 2 {
		for x := range 2 {
			img.Set(x, y, color.RGBA{R: 0x44, G: 0x88, B: 0xcc, A: 0xff})
		}
	}

	var data bytes.Buffer
	if err := png.Encode(&data, img); err != nil {
		t.Fatalf("failed to encode test png: %v", err)
	}
	return data.Bytes()
}

func assertStatus(t *testing.T, recorder *httptest.ResponseRecorder, want int) {
	t.Helper()

	if recorder.Code != want {
		t.Fatalf("status = %d, want %d\n%s", recorder.Code, want, recorder.Body.String())
	}
}

func assertRedirect(t *testing.T, recorder *httptest.ResponseRecorder, location string) {
	t.Helper()

	assertStatus(t, recorder, http.StatusSeeOther)
	if got := recorder.Header().Get("Location"); got != location {
		t.Fatalf("Location = %q, want %q\n%s", got, location, recorder.Body.String())
	}
}

func assertRedirectPrefix(t *testing.T, recorder *httptest.ResponseRecorder, prefix string) {
	t.Helper()

	assertStatus(t, recorder, http.StatusSeeOther)
	if got := recorder.Header().Get("Location"); !strings.HasPrefix(got, prefix) {
		t.Fatalf("Location = %q, want prefix %q\n%s", got, prefix, recorder.Body.String())
	}
}

func assertCookieSet(t *testing.T, recorder *httptest.ResponseRecorder, name string) {
	t.Helper()

	for _, cookie := range recorder.Result().Cookies() {
		if cookie.Name == name {
			return
		}
	}
	t.Fatalf("response did not set cookie %q", name)
}

func assertStringPointer(t *testing.T, name string, got *string, want string) {
	t.Helper()

	if got == nil {
		t.Fatalf("%s = nil, want %q", name, want)
	}
	if *got != want {
		t.Fatalf("%s = %q, want %q", name, *got, want)
	}
}

func countRows(t *testing.T, app *state.State, model any, query string, args ...any) int64 {
	t.Helper()

	var count int64
	if err := app.Database.Model(model).Where(query, args...).Count(&count).Error; err != nil {
		t.Fatalf("failed to count %T rows: %v", model, err)
	}
	return count
}

func assertRowCount(t *testing.T, app *state.State, model any, want int64, query string, args ...any) {
	t.Helper()

	if got := countRows(t, app, model, query, args...); got != want {
		t.Fatalf("%T row count = %d, want %d", model, got, want)
	}
}

func latestVerification(t *testing.T, app *state.State, userId int, verificationType constants.VerificationType) *schemas.Verification {
	t.Helper()

	var verification schemas.Verification
	err := app.Database.
		Where("user_id = ? AND type = ?", userId, verificationType).
		Order("id DESC").
		First(&verification).Error
	if err != nil {
		t.Fatalf("failed to fetch latest verification: %v", err)
	}
	return &verification
}
