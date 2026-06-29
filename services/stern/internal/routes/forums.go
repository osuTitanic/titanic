package routes

import (
	"fmt"
	"net/http"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/osuTitanic/titanic-go/internal/schemas"
	"github.com/osuTitanic/titanic-go/services/stern/internal/helpers"
	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
)

const forumTopicsPerPage = 25
const forumPostsPerPage = 15

var beatmapForumIds = map[int]bool{8: true, 9: true, 10: true, 12: true}
var topicPreloads = []string{"Icon", "Creator", "Creator.Groups.Group"}

func ForumHome(ctx *server.Context) {
	mainForums, err := ctx.State.Forums.FetchMainForums()
	if err != nil {
		ctx.Logger.Error("Failed to fetch main forums", "error", err)
		InternalServerError(ctx)
		return
	}

	currentUserId := 0
	if ctx.CurrentUser != nil {
		currentUserId = ctx.CurrentUser.Id
	}

	sections := make([]*templates.ForumSection, 0, len(mainForums))
	subForumIds := make([]int, 0)

	for _, mainForum := range mainForums {
		subForums, err := ctx.State.Forums.FetchSubForums(mainForum.Id)
		if err != nil {
			ctx.Logger.Error("Failed to fetch sub forums", "error", err, "forum", mainForum.Id)
			InternalServerError(ctx)
			return
		}

		subforums := make([]*templates.ForumSubforum, 0, len(subForums))
		for _, subForum := range subForums {
			subforums = append(subforums, &templates.ForumSubforum{
				Forum:         subForum,
				CurrentUserId: currentUserId,
			})
			subForumIds = append(subForumIds, subForum.Id)
		}

		sections = append(sections, &templates.ForumSection{
			Forum:     mainForum,
			Subforums: subforums,
		})
	}

	recent, err := ctx.State.ForumPosts.FetchLastForForums(
		subForumIds, "Topic", "User", "User.Groups.Group",
	)
	if err != nil {
		ctx.Logger.Error("Failed to fetch recent forum posts", "error", err)
		recent = map[int]*schemas.ForumPost{}
	}

	for _, section := range sections {
		for _, subforum := range section.Subforums {
			subforum.Recent = recent[subforum.Forum.Id]
		}
	}

	view := templates.ForumHomeView{
		DefaultView: buildDefaultView(ctx),
		Sections:    sections,
	}
	ctx.RenderTemplate(http.StatusOK, "pages/forum/home", view)
}

func ForumView(ctx *server.Context) {
	forumId, err := ctx.PathValueInt("id")
	if err != nil {
		NotFound(ctx)
		return
	}

	forum, err := ctx.State.Forums.ById(forumId)
	if err != nil {
		ctx.Logger.Error("Failed to fetch forum", "error", err, "forum", forumId)
		InternalServerError(ctx)
		return
	}
	if forum == nil || forum.Hidden {
		NotFound(ctx)
		return
	}

	if forum.ParentId == nil {
		// Main forums are only listed on the home page
		ctx.Redirect(http.StatusFound, "/forum")
		return
	}

	// Track the current user for "Users browsing this forum"
	helpers.ForumMarkUserActive(ctx, forum.Id)

	page := 1
	if parsed, err := strconv.Atoi(ctx.Request.URL.Query().Get("page")); err == nil && parsed > 1 {
		page = parsed
	}
	offset := (page - 1) * forumTopicsPerPage

	subForums, err := ctx.State.Forums.FetchSubForums(forum.Id)
	if err != nil {
		ctx.Logger.Error("Failed to fetch sub forums", "error", err, "forum", forum.Id)
		InternalServerError(ctx)
		return
	}

	topicCount, err := ctx.State.Forums.FetchTopicCount(forum.Id)
	if err != nil {
		ctx.Logger.Error("Failed to fetch topic count", "error", err, "forum", forum.Id)
		InternalServerError(ctx)
		return
	}

	recentTopics, err := ctx.State.ForumTopics.FetchRecentByLastPost(forum.Id, forumTopicsPerPage, offset, topicPreloads...)
	if err != nil {
		ctx.Logger.Error("Failed to fetch recent topics", "error", err, "forum", forum.Id)
		InternalServerError(ctx)
		return
	}

	pinnedTopics, err := ctx.State.ForumTopics.FetchPinnedByForumId(forum.Id, topicPreloads...)
	if err != nil {
		ctx.Logger.Error("Failed to fetch pinned topics", "error", err, "forum", forum.Id)
		InternalServerError(ctx)
		return
	}

	announcements, err := ctx.State.ForumTopics.FetchAnnouncementsByForumId(forum.Id, 3, 0, topicPreloads...)
	if err != nil {
		ctx.Logger.Error("Failed to fetch announcements", "error", err, "forum", forum.Id)
		InternalServerError(ctx)
		return
	}

	// Merge pinned topics into the recent listing so they float to the top
	topics := mergeForumTopics(pinnedTopics, recentTopics)

	topicIds := make([]int, 0, len(topics)+len(announcements))
	for _, topic := range topics {
		topicIds = append(topicIds, topic.Id)
	}
	for _, topic := range announcements {
		topicIds = append(topicIds, topic.Id)
	}

	lastPosts, err := ctx.State.ForumPosts.FetchLastForTopics(topicIds, "User", "User.Groups.Group")
	if err != nil {
		ctx.Logger.Error("Failed to fetch last posts", "error", err, "forum", forum.Id)
		lastPosts = map[int]*schemas.ForumPost{}
	}

	readStatuses := helpers.ForumTopicReadStatuses(ctx, slices.Concat(announcements, topics))
	averageViews := helpers.ForumAverageTopicViews(ctx)

	hasCustomIcons := false
	for _, topic := range topics {
		if topic.IconId != nil {
			hasCustomIcons = true
			break
		}
	}

	subForumIds := make([]int, 0, len(subForums))
	for _, subForum := range subForums {
		subForumIds = append(subForumIds, subForum.Id)
	}

	subforumRecent, err := ctx.State.ForumPosts.FetchLastForForums(
		subForumIds, "Topic", "User", "User.Groups.Group",
	)
	if err != nil {
		ctx.Logger.Error("Failed to fetch sub forum recent posts", "error", err, "forum", forum.Id)
		subforumRecent = map[int]*schemas.ForumPost{}
	}

	currentUserId := 0
	if ctx.CurrentUser != nil {
		currentUserId = ctx.CurrentUser.Id
	}

	view := templates.ForumView{
		DefaultView:    buildDefaultView(ctx),
		Forum:          forum,
		Subforums:      subForums,
		SubforumRecent: subforumRecent,
		Parents:        fetchForumParents(ctx, forum),
		Announcements:  buildTopicPreviews(announcements, lastPosts, readStatuses, averageViews, forum.Id, hasCustomIcons, currentUserId),
		Topics:         buildTopicPreviews(topics, lastPosts, readStatuses, averageViews, forum.Id, hasCustomIcons, currentUserId),
		ActiveUsers:    fetchActiveForumUsers(ctx, forum.Id),
		CanCreateTopic: canCreateForumTopic(ctx, forum),
		HasCustomIcons: hasCustomIcons,
		TopicCount:     topicCount,
		Pagination: templates.NewPagination(templates.PaginationOptions{
			Path:        fmt.Sprintf("/forum/%d", forum.Id),
			Query:       ctx.Request.URL.Query(),
			CurrentPage: page,
			Total:       topicCount,
			PageSize:    forumTopicsPerPage,
		}),
	}
	ctx.RenderTemplate(http.StatusOK, "pages/forum/forum", view)
}

func mergeForumTopics(pinned, recent []*schemas.ForumTopic) []*schemas.ForumTopic {
	seen := make(map[int]bool, len(pinned)+len(recent))
	merged := make([]*schemas.ForumTopic, 0, len(pinned)+len(recent))

	for _, topic := range pinned {
		if !seen[topic.Id] {
			seen[topic.Id] = true
			merged = append(merged, topic)
		}
	}
	for _, topic := range recent {
		if !seen[topic.Id] {
			seen[topic.Id] = true
			merged = append(merged, topic)
		}
	}

	// Sort the merged topics so that pinned topics appear first
	sort.SliceStable(merged, func(i, j int) bool {
		a, b := merged[i], merged[j]
		if a.Pinned != b.Pinned {
			return a.Pinned
		}
		return a.LastPostAt.After(b.LastPostAt)
	})
	return merged
}

// long ass function definition incoming, beware

func buildTopicPreviews(
	topics []*schemas.ForumTopic,
	lastPosts map[int]*schemas.ForumPost,
	readStatuses map[int]bool,
	averageViews float64,
	forumId int,
	hasCustomIcons bool,
	currentUserId int,
) []*templates.ForumTopicPreview {
	previews := make([]*templates.ForumTopicPreview, 0, len(topics))
	for index, topic := range topics {
		previews = append(previews, &templates.ForumTopicPreview{
			Topic:          topic,
			LastPost:       lastPosts[topic.Id],
			StatusIcon:     topicStatusIcon(topic, readStatuses[topic.Id], averageViews),
			PageCount:      (topic.PostCount + forumPostsPerPage - 1) / forumPostsPerPage,
			Index:          index,
			ForumId:        forumId,
			HasCustomIcons: hasCustomIcons,
			CurrentUserId:  currentUserId,
		})
	}
	return previews
}

func topicStatusIcon(topic *schemas.ForumTopic, read bool, averageViews float64) string {
	state := "unread"
	if read {
		state = "read"
	}

	if topic.Pinned || topic.Announcement {
		if topic.LockedAt != nil {
			return fmt.Sprintf("/images/icons/topics/announce_%s_locked.gif", state)
		}
		return fmt.Sprintf("/images/icons/topics/announce_%s.gif", state)
	}

	if topic.LockedAt != nil {
		return fmt.Sprintf("/images/icons/topics/topic_%s_locked.gif", state)
	}

	// A topic is considered "hot" if it has more than half
	// the average views and was created within the last 7 days
	age := time.Since(topic.CreatedAt)
	isHot := float64(topic.Views) > averageViews/2 && age < 7*24*time.Hour

	if isHot {
		return fmt.Sprintf("/images/icons/topics/topic_%s_hot.gif", state)
	}
	return fmt.Sprintf("/images/icons/topics/topic_%s.gif", state)
}

func fetchForumParents(ctx *server.Context, forum *schemas.Forum) []*schemas.Forum {
	parents := make([]*schemas.Forum, 0)
	current := forum

	for current.ParentId != nil {
		parent, err := ctx.State.Forums.ById(*current.ParentId)
		if err != nil || parent == nil {
			break
		}
		parents = slices.Concat([]*schemas.Forum{parent}, parents)
		current = parent
	}
	return parents
}

func fetchActiveForumUsers(ctx *server.Context, forumId int) []*templates.ForumActiveUser {
	activeIds := helpers.ForumGetActiveUsers(ctx, forumId)
	if len(activeIds) == 0 {
		return nil
	}

	users, err := ctx.State.Users.ManyById(activeIds)
	if err != nil {
		ctx.Logger.Error("Failed to fetch active forum users", "error", err, "forum", forumId)
		return nil
	}

	nameById := make(map[int]string, len(users))
	for _, user := range users {
		nameById[user.Id] = user.Name
	}

	activeUsers := make([]*templates.ForumActiveUser, 0, len(activeIds))
	for _, id := range activeIds {
		if name, ok := nameById[id]; ok {
			activeUsers = append(activeUsers, &templates.ForumActiveUser{Id: id, Name: name})
		}
	}
	return activeUsers
}

func canCreateForumTopic(ctx *server.Context, forum *schemas.Forum) bool {
	if ctx.CurrentUser == nil {
		return false
	}
	if !ctx.HasPermission("forum.topics.create") {
		return false
	}
	if beatmapForumIds[forum.Id] {
		// Inside the beatmap forums, we want to check for a special permission
		return ctx.HasPermission("forum.topics.create_beatmap")
	}
	return true
}
