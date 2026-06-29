package routes

import (
	"fmt"
	"net/http"
	"slices"
	"sort"
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
	if parsed, err := ctx.QueryValueInt("page"); err == nil && parsed > 1 {
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

func ForumTopicView(ctx *server.Context) {
	forumId, err := ctx.PathValueInt("id")
	if err != nil {
		NotFound(ctx)
		return
	}

	topicId, err := ctx.PathValueInt("topicId")
	if err != nil {
		NotFound(ctx)
		return
	}

	topic, err := ctx.State.ForumTopics.ById(topicId, "Forum", "Icon")
	if err != nil {
		ctx.Logger.Error("Failed to fetch topic", "error", err, "topic", topicId)
		InternalServerError(ctx)
		return
	}
	if topic == nil || topic.Hidden {
		NotFound(ctx)
		return
	}

	if topic.ForumId != forumId {
		// Redirect to the correct url for this topic
		ctx.Redirect(http.StatusFound, fmt.Sprintf("/forum/%d/t/%d/", topic.ForumId, topic.Id))
		return
	}

	page := 1
	if parsed, err := ctx.QueryValueInt("page"); err == nil && parsed > 1 {
		page = parsed
	}
	offset := (page - 1) * forumPostsPerPage

	posts, err := ctx.State.ForumPosts.FetchRangeByTopic(
		topic.Id, forumPostsPerPage, offset,
		"Icon", "User", "User.Groups.Group",
	)
	if err != nil {
		ctx.Logger.Error("Failed to fetch topic posts", "error", err, "topic", topic.Id)
		InternalServerError(ctx)
		return
	}

	postCount, err := ctx.State.ForumPosts.CountByTopic(topic.Id)
	if err != nil {
		ctx.Logger.Error("Failed to fetch post count", "error", err, "topic", topic.Id)
		InternalServerError(ctx)
		return
	}

	// Mark the topic as read & bump its view counter for this visitor
	helpers.ForumUpdateTopicReadState(ctx, topic.Id)
	helpers.ForumUpdateViews(ctx, topic.Id)
	helpers.ForumMarkUserActive(ctx, topic.ForumId)

	// Resolve the very first post in the topic, used for meta tags
	initialPost := new(schemas.ForumPost)

	if page == 1 && len(posts) > 0 {
		// We're on the first page so we don't need to explicitly query for it
		initialPost = posts[0]
	} else {
		initialPost, err = ctx.State.ForumPosts.FetchInitialByTopic(topic.Id, "User")
		if err != nil {
			ctx.Logger.Error("Failed to fetch initial post", "error", err, "topic", topic.Id)
		}
	}

	if initialPost == nil {
		NotFound(ctx)
		return
	}

	// Override the icon of the initial post with the topic's icon
	if page == 1 && len(posts) > 0 {
		posts[0].Icon = topic.Icon
	}

	postUserIds := make([]int, 0, len(posts))
	for _, post := range posts {
		postUserIds = append(postUserIds, post.UserId)
	}

	postCounts, err := ctx.State.ForumPosts.PostCountsByUsers(postUserIds)
	if err != nil {
		ctx.Logger.Error("Failed to fetch user post counts", "error", err, "topic", topic.Id)
		postCounts = map[int]int{}
	}

	isSubscribed := false
	isBookmarked := false
	if ctx.CurrentUser != nil {
		isSubscribed, _ = ctx.State.ForumSubscribers.Exists(topic.Id, ctx.CurrentUser.Id)
		isBookmarked, _ = ctx.State.ForumBookmarks.Exists(topic.Id, ctx.CurrentUser.Id)
	}

	// Resolve the permissions that gate the topic & post actions
	authenticated := ctx.CurrentUser != nil
	canCreatePosts := authenticated && ctx.HasPermission("forum.posts.create")
	canEditOwn := authenticated && ctx.HasPermission("forum.posts.edit")
	canDeleteOwn := authenticated && ctx.HasPermission("forum.posts.delete")
	canEditOthers := authenticated && ctx.HasPermission("forum.moderation.posts.edit")
	canDeleteOthers := authenticated && ctx.HasPermission("forum.moderation.posts.delete")
	canBypassTopicLock := authenticated && ctx.HasPermission("forum.moderation.topics.bypass_lock")
	canBypassPostLock := authenticated && ctx.HasPermission("forum.moderation.posts.bypass_lock")

	topicLocked := topic.LockedAt != nil
	showActions := authenticated && (!topicLocked || canBypassTopicLock)

	previews := make([]*templates.ForumPostPreview, 0, len(posts))
	for _, post := range posts {
		isOwn := authenticated && ctx.CurrentUser.Id == post.UserId

		editable := !post.EditLocked || canBypassPostLock
		canModify := showActions && !post.Deleted && editable

		canDelete := (isOwn && canDeleteOwn) || (!isOwn && canDeleteOthers)
		canEdit := (isOwn && canEditOwn) || (!isOwn && canEditOthers)

		previews = append(previews, &templates.ForumPostPreview{
			Post:        post,
			Icon:        post.Icon,
			AuthorTitle: forumUserTitle(post.User, postCounts[post.UserId]),
			PostCount:   postCounts[post.UserId],
			CanDelete:   canModify && canDelete,
			CanEdit:     canModify && canEdit,
			CanQuote:    showActions && canCreatePosts,
		})
	}

	// Just to be safe here (this should actually never be nil)
	metaImage := ""
	if initialPost.User != nil {
		metaImage = ctx.State.Config.OsuBaseUrl() + initialPost.User.AvatarUrl()
	}

	view := templates.ForumTopicView{
		DefaultView:     buildDefaultViewWithPermissions(ctx),
		Forum:           topic.Forum,
		Topic:           topic,
		Parents:         fetchForumParents(ctx, topic.Forum),
		Posts:           previews,
		ActiveUsers:     fetchActiveForumUsers(ctx, topic.ForumId),
		PostCount:       postCount,
		IsSubscribed:    isSubscribed,
		IsBookmarked:    isBookmarked,
		CanCreatePosts:  canCreatePosts,
		CanReply:        canCreatePosts && (!topicLocked || canBypassTopicLock),
		ReplyLocked:     topicLocked && !canBypassTopicLock,
		MetaDescription: strings.SplitN(initialPost.Content, "\n", 2)[0],
		MetaImage:       metaImage,
		Pagination: templates.NewPagination(templates.PaginationOptions{
			Path:        fmt.Sprintf("/forum/%d/t/%d/", topic.ForumId, topic.Id),
			Query:       ctx.Request.URL.Query(),
			CurrentPage: page,
			Total:       postCount,
			PageSize:    forumPostsPerPage,
		}),
	}
	ctx.RenderTemplate(http.StatusOK, "pages/forum/topic", view)
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

func forumUserTitle(user *schemas.User, postCount int) string {
	if user != nil {
		if title := user.TitleText(); title != "" {
			return title
		}
	}

	switch {
	case postCount < 5:
		return "Rhythm Rookie"
	case postCount < 15:
		return "Tempo Trainee"
	case postCount < 30:
		return "Whistle Blower"
	case postCount < 50:
		return "Cymbal Sounder"
	case postCount < 80:
		return "Beat Clicker"
	case postCount < 120:
		return "Slider Savant"
	case postCount < 180:
		return "Spinner Sage"
	case postCount < 260:
		return "Star Shooter"
	case postCount < 500:
		return "Combo Commander"
	default:
		return "Rhythm Incarnate"
	}
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
	if forum.IsBeatmapForum() {
		return ctx.HasPermission("forum.topics.create_beatmap")
	}
	return true
}
