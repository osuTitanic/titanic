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
	"github.com/osuTitanic/titanic/internal/authentication"
	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
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

func (data *TestData) CreateBeatmapset(creator *schemas.User, opts ...FixtureOption[schemas.Beatmapset]) *schemas.Beatmapset {
	data.t.Helper()
	data.requireUser(creator)

	sequence := data.next()
	now := testDataTime(sequence)

	title := fmt.Sprintf("Test Title %d", sequence)
	artist := fmt.Sprintf("Test Artist %d", sequence)
	description := "I can't pass this map. Hell diff. this is easy diff but hell difficult."
	tags := "real stream"

	beatmapset := &schemas.Beatmapset{
		Title:          new(title),
		Artist:         new(artist),
		Creator:        new(creator.Name),
		Description:    new(description),
		Tags:           new(tags),
		Status:         constants.BeatmapStatusRanked,
		Server:         constants.BeatmapServerTitanic,
		DownloadServer: constants.BeatmapServerTitanic,
		CreatorId:      &creator.Id,
		Available:      true,
		CreatedAt:      now,
		ApprovedAt:     &now,
		LastUpdate:     now,
		AddedAt:        &now,
		LanguageId:     constants.BeatmapLanguageEnglish,
		GenreId:        constants.BeatmapGenrePop,
	}
	applyFixtureOptions(beatmapset, opts)
	data.create(beatmapset)
	return beatmapset
}

func (data *TestData) CreateBeatmap(beatmapset *schemas.Beatmapset, opts ...FixtureOption[schemas.Beatmap]) *schemas.Beatmap {
	data.t.Helper()
	data.requireBeatmapset(beatmapset)

	sequence := data.next()
	now := testDataTime(sequence)
	beatmap := &schemas.Beatmap{
		SetId:            beatmapset.Id,
		Mode:             constants.ModeOsu,
		Status:           constants.BeatmapStatusRanked,
		Checksum:         fmt.Sprintf("test-md5-%d", sequence),
		Version:          "Insane",
		Filename:         fmt.Sprintf("%d Test Artist - Test Title.osu", sequence),
		CreatedAt:        now,
		LastUpdate:       now,
		TotalLength:      180,
		DrainLength:      170,
		CountNormal:      300,
		CountSlider:      120,
		MaxCombo:         543,
		BPM:              180,
		CS:               4,
		AR:               9,
		OD:               8,
		HP:               6,
		Diff:             4.5,
		SliderMultiplier: 1.4,
	}
	applyFixtureOptions(beatmap, opts)
	data.create(beatmap)
	return beatmap
}

func (data *TestData) CreateScore(user *schemas.User, beatmap *schemas.Beatmap, opts ...FixtureOption[schemas.Score]) *schemas.Score {
	data.t.Helper()
	data.requireUser(user)
	data.requireBeatmap(beatmap)

	sequence := data.next()
	score := &schemas.Score{
		UserId:        user.Id,
		BeatmapId:     beatmap.Id,
		ClientVersion: 20121212,
		ClientString:  "b20121212",
		Checksum:      fmt.Sprintf("score-%d", sequence),
		Mode:          beatmap.Mode,
		PP:            250.25,
		PPv1:          180.25,
		Acc:           0.9876,
		TotalScore:    987654,
		MaxCombo:      beatmap.MaxCombo,
		Mods:          constants.NoMod,
		Perfect:       true,
		Count300:      450,
		Count100:      12,
		Count50:       1,
		Grade:         constants.GradeS,
		StatusPP:      constants.ScoreStatusBest,
		StatusScore:   constants.ScoreStatusBest,
		SubmittedAt:   testDataTime(sequence),
	}
	applyFixtureOptions(score, opts)
	data.create(score)
	return score
}

func (data *TestData) CreateRelationship(user *schemas.User, target *schemas.User, opts ...FixtureOption[schemas.Relationship]) *schemas.Relationship {
	data.t.Helper()
	data.requireUser(user)
	data.requireUser(target)

	relationship := &schemas.Relationship{
		UserId:   user.Id,
		TargetId: target.Id,
		Status:   constants.RelationshipStatusFriend,
	}
	applyFixtureOptions(relationship, opts)
	data.create(relationship)
	return relationship
}

func (data *TestData) CreateGroup(opts ...FixtureOption[schemas.Group]) *schemas.Group {
	data.t.Helper()

	sequence := data.next()
	group := &schemas.Group{
		Name:        fmt.Sprintf("Test Group %d", sequence),
		ShortName:   fmt.Sprintf("TG%d", sequence),
		Description: new("cool people"),
		Color:       "#3366cc",
	}
	applyFixtureOptions(group, opts)
	data.create(group)
	return group
}

func (data *TestData) CreateGroupEntry(group *schemas.Group, user *schemas.User, opts ...FixtureOption[schemas.GroupEntry]) *schemas.GroupEntry {
	data.t.Helper()
	data.requireGroup(group)
	data.requireUser(user)

	entry := &schemas.GroupEntry{
		GroupId: group.Id,
		UserId:  user.Id,
	}
	applyFixtureOptions(entry, opts)
	data.create(entry)
	return entry
}

func (data *TestData) CreateBeatmapPack(creator *schemas.User, opts ...FixtureOption[schemas.BeatmapPack]) *schemas.BeatmapPack {
	data.t.Helper()
	data.requireUser(creator)

	sequence := data.next()
	pack := &schemas.BeatmapPack{
		Name:         fmt.Sprintf("test-pack-%d", sequence),
		Category:     "category",
		DownloadLink: fmt.Sprintf("https://example.com/packs/%d.zip", sequence),
		Description:  "description",
		CreatorId:    creator.Id,
		CreatedAt:    testDataTime(sequence),
		UpdatedAt:    testDataTime(sequence),
	}
	applyFixtureOptions(pack, opts)
	data.create(pack)
	return pack
}

func (data *TestData) CreateBeatmapPackEntry(pack *schemas.BeatmapPack, beatmapset *schemas.Beatmapset, opts ...FixtureOption[schemas.BeatmapPackEntry]) *schemas.BeatmapPackEntry {
	data.t.Helper()
	data.requireBeatmapPack(pack)
	data.requireBeatmapset(beatmapset)

	entry := &schemas.BeatmapPackEntry{
		PackId:       pack.Id,
		BeatmapsetId: beatmapset.Id,
		CreatedAt:    testDataTime(data.next()),
	}
	applyFixtureOptions(entry, opts)
	data.create(entry)
	return entry
}

func (data *TestData) CreateBeatmapModding(target *schemas.User, sender *schemas.User, beatmapset *schemas.Beatmapset, post *schemas.ForumPost, opts ...FixtureOption[schemas.BeatmapModding]) *schemas.BeatmapModding {
	data.t.Helper()
	data.requireUser(target)
	data.requireUser(sender)
	data.requireBeatmapset(beatmapset)
	data.requireForumPost(post)

	modding := &schemas.BeatmapModding{
		TargetId: target.Id,
		SenderId: sender.Id,
		SetId:    beatmapset.Id,
		PostId:   post.Id,
		Amount:   1,
		Time:     testDataTime(data.next()),
	}
	applyFixtureOptions(modding, opts)
	data.create(modding)
	return modding
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

func (data *TestData) requireForumPost(post *schemas.ForumPost) {
	data.t.Helper()
	if post == nil || post.Id == 0 {
		data.t.Fatal("test data requires a persisted forum post")
	}
}

func (data *TestData) requireBeatmapset(beatmapset *schemas.Beatmapset) {
	data.t.Helper()
	if beatmapset == nil || beatmapset.Id == 0 {
		data.t.Fatal("test data requires a persisted beatmapset")
	}
}

func (data *TestData) requireBeatmap(beatmap *schemas.Beatmap) {
	data.t.Helper()
	if beatmap == nil || beatmap.Id == 0 {
		data.t.Fatal("test data requires a persisted beatmap")
	}
}

func (data *TestData) requireGroup(group *schemas.Group) {
	data.t.Helper()
	if group == nil || group.Id == 0 {
		data.t.Fatal("test data requires a persisted group")
	}
}

func (data *TestData) requireBeatmapPack(pack *schemas.BeatmapPack) {
	data.t.Helper()
	if pack == nil || pack.Id == 0 {
		data.t.Fatal("test data requires a persisted beatmap pack")
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
