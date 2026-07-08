//go:build integration

package state

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/osuTitanic/titanic-go/internal/authentication"
	"github.com/osuTitanic/titanic-go/internal/constants"
	"github.com/osuTitanic/titanic-go/internal/schemas"
)

type TestData struct {
	state    *State
	t        testing.TB
	mu       sync.Mutex
	sequence int
}
type FixtureOption[T any] func(*T)

func NewTestData(t testing.TB, state *State) *TestData {
	t.Helper()
	if state == nil {
		t.Fatal("test data requires state")
	}
	return &TestData{t: t, state: state}
}

func (data *TestData) CreateUser(opts ...FixtureOption[schemas.User]) *schemas.User {
	data.t.Helper()

	sequence := data.next()
	now := testDataTime(sequence)
	passwordHash, err := authentication.CreatePasswordHash("password")
	if err != nil {
		data.t.Fatalf("failed to create test user password hash: %v", err)
	}

	user := &schemas.User{
		Name:             fmt.Sprintf("TestUser%d", sequence),
		SafeName:         fmt.Sprintf("testuser%d", sequence),
		Email:            fmt.Sprintf("testuser%d@example.com", sequence),
		Bcrypt:           passwordHash,
		Country:          "XX",
		CreatedAt:        now,
		LatestActivity:   now,
		Activated:        true,
		PreferredMode:    constants.ModeOsu,
		PreferredRanking: constants.RankingTypePerformance,
		IrcToken:         fmt.Sprintf("token%d", sequence),
		AvatarLastUpdate: now,
	}
	applyFixtureOptions(user, opts)
	data.create(user)
	return user
}

func (data *TestData) CreateStats(user *schemas.User, opts ...FixtureOption[schemas.Stats]) *schemas.Stats {
	data.t.Helper()
	data.requireUser(user)

	stats := &schemas.Stats{
		UserId:    user.Id,
		Mode:      constants.ModeOsu,
		Rank:      1,
		Tscore:    123456789,
		Rscore:    123456,
		PP:        1234.56,
		PPv1:      727.27,
		Playcount: 69,
		Acc:       0.9876,
		TotalHits: 12345,
	}
	applyFixtureOptions(stats, opts)
	data.create(stats)
	return stats
}

func (data *TestData) CreateLogin(user *schemas.User, opts ...FixtureOption[schemas.Login]) *schemas.Login {
	data.t.Helper()
	data.requireUser(user)

	login := &schemas.Login{
		UserId:  user.Id,
		Time:    testDataTime(data.next()),
		Ip:      "1.2.3.4",
		Version: "web",
	}
	applyFixtureOptions(login, opts)
	data.create(login)
	return login
}

func (data *TestData) CreateWebsiteSessionCookie(user *schemas.User, request *http.Request) *http.Cookie {
	data.t.Helper()
	data.requireUser(user)

	ttl := 24 * time.Hour
	session, err := data.state.SessionStore.Create(context.Background(), user.Id, time.Now(), ttl)
	if err != nil {
		data.t.Fatalf("failed to create website session: %v", err)
	}
	return authentication.NewWebsiteSessionCookie(data.state.Config, request, session.Id, ttl)
}

func (data *TestData) CreateForum(opts ...FixtureOption[schemas.Forum]) *schemas.Forum {
	data.t.Helper()

	sequence := data.next()
	forum := &schemas.Forum{
		Name:        fmt.Sprintf("Test Forum %d", sequence),
		Description: "A forum™",
		CreatedAt:   testDataTime(sequence),
		AllowIcons:  true,
	}
	applyFixtureOptions(forum, opts)
	data.create(forum)
	return forum
}

func (data *TestData) CreateForumTopic(forum *schemas.Forum, creator *schemas.User, opts ...FixtureOption[schemas.ForumTopic]) *schemas.ForumTopic {
	data.t.Helper()
	data.requireForum(forum)
	data.requireUser(creator)

	sequence := data.next()
	now := testDataTime(sequence)
	topic := &schemas.ForumTopic{
		ForumId:       forum.Id,
		CreatorId:     creator.Id,
		Title:         fmt.Sprintf("Test Topic %d", sequence),
		CreatedAt:     now,
		LastPostAt:    now,
		CanChangeIcon: true,
	}
	applyFixtureOptions(topic, opts)
	data.create(topic)
	return topic
}

func (data *TestData) CreateForumPost(topic *schemas.ForumTopic, author *schemas.User, opts ...FixtureOption[schemas.ForumPost]) *schemas.ForumPost {
	data.t.Helper()
	data.requireTopic(topic)
	data.requireUser(author)

	sequence := data.next()
	post := &schemas.ForumPost{
		TopicId:   topic.Id,
		ForumId:   topic.ForumId,
		UserId:    author.Id,
		Content:   fmt.Sprintf("Test post %d", sequence),
		CreatedAt: testDataTime(sequence),
		EditTime:  testDataTime(sequence),
	}
	applyFixtureOptions(post, opts)
	data.create(post)
	return post
}

func (data *TestData) CreateForumBookmark(user *schemas.User, topic *schemas.ForumTopic, opts ...FixtureOption[schemas.ForumBookmark]) *schemas.ForumBookmark {
	data.t.Helper()
	data.requireUser(user)
	data.requireTopic(topic)

	bookmark := &schemas.ForumBookmark{
		UserId:  user.Id,
		TopicId: topic.Id,
	}
	applyFixtureOptions(bookmark, opts)
	data.create(bookmark)
	return bookmark
}

func (data *TestData) CreateNotification(user *schemas.User, opts ...FixtureOption[schemas.Notification]) *schemas.Notification {
	data.t.Helper()
	data.requireUser(user)

	notification := &schemas.Notification{
		UserId:  user.Id,
		Type:    constants.NotificationTypeWelcome,
		Header:  "Test notification",
		Content: "This notification was sponsored by nordvpn",
		Link:    "/",
		Time:    testDataTime(data.next()),
	}
	applyFixtureOptions(notification, opts)
	data.create(notification)
	return notification
}

func (data *TestData) CreateRelease(opts ...FixtureOption[schemas.Release]) *schemas.Release {
	data.t.Helper()

	sequence := data.next()
	release := &schemas.Release{
		Name:        fmt.Sprintf("test-release-%d", sequence),
		Version:     sequence,
		Description: "osu!auth.dll spyware update)))",
		Category:    "Unstable",
		Supported:   true,
		Downloads:   pq.StringArray{fmt.Sprintf("https://example.com/releases/%d.zip", sequence)},
		Screenshots: pq.StringArray{},
		Hashes:      json.RawMessage("[]"),
		CreatedAt:   testDataTime(sequence),
	}
	applyFixtureOptions(release, opts)
	data.create(release)
	return release
}

func (data *TestData) CreateMessage(opts ...FixtureOption[schemas.Message]) *schemas.Message {
	data.t.Helper()

	sequence := data.next()
	message := &schemas.Message{
		Target:  "#osu",
		Sender:  fmt.Sprintf("TestUser%d", sequence),
		Message: fmt.Sprintf("Test chat message %d", sequence),
		Time:    testDataTime(sequence),
	}
	applyFixtureOptions(message, opts)
	data.create(message)
	return message
}

func (data *TestData) create(value any) {
	data.t.Helper()
	if err := data.state.Database.Create(value).Error; err != nil {
		data.t.Fatalf("failed to create test data %T: %v", value, err)
	}
}

func (data *TestData) next() int {
	data.mu.Lock()
	defer data.mu.Unlock()

	data.sequence++
	return data.sequence
}

func (data *TestData) requireUser(user *schemas.User) {
	data.t.Helper()
	if user == nil || user.Id == 0 {
		data.t.Fatal("test data requires a persisted user")
	}
}

func (data *TestData) requireForum(forum *schemas.Forum) {
	data.t.Helper()
	if forum == nil || forum.Id == 0 {
		data.t.Fatal("test data requires a persisted forum")
	}
}

func (data *TestData) requireTopic(topic *schemas.ForumTopic) {
	data.t.Helper()
	if topic == nil || topic.Id == 0 {
		data.t.Fatal("test data requires a persisted forum topic")
	}
}

func applyFixtureOptions[T any](value *T, opts []FixtureOption[T]) {
	for _, opt := range opts {
		if opt != nil {
			opt(value)
		}
	}
}

func testDataTime(sequence int) time.Time {
	return time.Date(2012, 1, 1, 12, 0, 0, 0, time.UTC).Add(time.Duration(sequence) * time.Minute)
}
