package routes

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/osuTitanic/titanic-go/internal/constants"
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
)

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
	forumUpdateTopicReadState(ctx, topic.Id)
	forumUpdateViews(ctx, topic.Id)
	forumMarkUserActive(ctx, topic.ForumId)

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
	postIds := make([]int, 0, len(posts))
	for _, post := range posts {
		postUserIds = append(postUserIds, post.UserId)
		postIds = append(postIds, int(post.Id))
	}

	postCounts, err := ctx.State.ForumPosts.PostCountsByUsers(postUserIds)
	if err != nil {
		ctx.Logger.Error("Failed to fetch user post counts", "error", err, "topic", topic.Id)
		postCounts = map[int]int{}
	}

	linkedBeatmapset, err := ctx.State.Beatmapsets.ByTopicId(topic.Id)
	if err != nil {
		ctx.Logger.Error("Failed to fetch beatmapset by topic", "error", err, "topic", topic.Id)
		linkedBeatmapset = nil
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

		kudosuTotals, latestKudosu := fetchKudosuForPosts(
			postIds, linkedBeatmapset, ctx,
		)
		preview := &templates.ForumPostPreview{
			Post:         post,
			Icon:         post.Icon,
			AuthorTitle:  forumUserTitle(post.User, postCounts[post.UserId]),
			PostCount:    postCounts[post.UserId],
			CanDelete:    canModify && canDelete,
			CanEdit:      canModify && canEdit,
			CanQuote:     showActions && canCreatePosts,
			KudosuTotal:  kudosuTotals[post.Id],
			LatestKudosu: latestKudosu[post.Id],
		}

		if linkedBeatmapset != nil && linkedBeatmapset.CreatorId != nil {
			canForceRewardKudosu := authenticated && ctx.HasPermission("beatmaps.moderation.force_nominate")
			isBeatmapsetCreator := authenticated && *linkedBeatmapset.CreatorId == ctx.CurrentUser.Id

			preview.BeatmapsetId = linkedBeatmapset.Id
			preview.CanResetKudosu = authenticated && ctx.HasPermission("forum.kudosu.reset")   // && !linkedBeatmapset.IsApproved()
			preview.CanRevokeKudosu = authenticated && ctx.HasPermission("forum.kudosu.revoke") // && !linkedBeatmapset.IsApproved()
			preview.ShowKudosuBox = post.UserId != *linkedBeatmapset.CreatorId && !preview.HasKudosuExcludedIcon()

			// When the user is owner of the set -> allow kudsou awards while set is unranked
			// When the user is a BAT member -> allow deny / reset actions even when set it ranked
			// When the user is a BAT manager -> allow all kudosu actions even when set it ranked

			preview.CanManageKudosu = (isBeatmapsetCreator && !linkedBeatmapset.IsApproved()) || canForceRewardKudosu
			preview.CanManageKudosu = preview.CanManageKudosu || (preview.CanResetKudosu || preview.CanRevokeKudosu)
		}

		previews = append(previews, preview)
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
		Beatmapset:      linkedBeatmapset,
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

func ForumCreateTopicView(ctx *server.Context) {
	if !ctx.RequireLogin() {
		return
	}

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

	if !canCreateForumTopic(ctx, forum) {
		RenderErrorPage(ctx, http.StatusForbidden, "Forbidden", "You are not allowed to create topics in this forum.")
		return
	}

	editor := templates.ForumEditorContext{
		SubmitText:     "Create Topic",
		CancelUrl:      fmt.Sprintf("/forum/%d", forum.Id),
		ShowSubject:    true,
		ShowIcons:      canEditForumIcon(ctx, forum.AllowIcons),
		Icons:          buildEditorIcons(fetchForumIcons(ctx), -1),
		ShowControls:   true,
		ShowTopicTypes: ctx.HasPermission("forum.moderation.topics.set_options"),
		TopicType:      "global", // TODO: perhaps add an enum for this
	}
	editor.NoneIconSelected = true

	view := templates.ForumCreateTopicView{
		DefaultView: buildDefaultViewWithPermissions(ctx),
		Forum:       forum,
		Parents:     fetchForumParents(ctx, forum),
		Editor:      editor,
	}
	ctx.RenderTemplate(http.StatusOK, "pages/forum/create", view)
}

func ForumCreateTopicAction(ctx *server.Context) {
	if !ctx.RequireLogin() {
		return
	}

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

	if valid, err := ctx.ValidateCSRF(); err != nil || !valid {
		RenderErrorPage(ctx, http.StatusForbidden, "Invalid Request", "Your session has expired, please try again.")
		return
	}

	if !canCreateForumTopic(ctx, forum) {
		RenderErrorPage(ctx, http.StatusForbidden, "Forbidden", "You are not allowed to create topics in this forum.")
		return
	}
	if isPostingRejected(ctx) {
		return
	}

	title := strings.TrimSpace(ctx.Request.FormValue("title"))
	content := ctx.Request.FormValue("bbcode")
	if title == "" || content == "" {
		ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/forum/%d", forum.Id))
		return
	}

	pinned, announcement := resolveTopicType(ctx)
	canEditIcon := canEditForumIcon(ctx, forum.AllowIcons)

	topic := &schemas.ForumTopic{
		ForumId:       forum.Id,
		CreatorId:     ctx.CurrentUser.Id,
		IconId:        resolveSubmittedIcon(ctx, canEditIcon),
		CanChangeIcon: forum.AllowIcons,
		Title:         title,
		Pinned:        pinned,
		Announcement:  announcement,
		CreatedAt:     time.Now().UTC(),
	}
	if err := ctx.State.ForumTopics.Create(topic); err != nil {
		ctx.Logger.Error("Failed to create topic", "error", err, "forum", forum.Id)
		InternalServerError(ctx)
		return
	}

	post := &schemas.ForumPost{
		TopicId: topic.Id,
		ForumId: forum.Id,
		UserId:  ctx.CurrentUser.Id,
		Content: content,
		IconId:  topic.IconId,
	}
	if err := ctx.State.ForumPosts.Create(post); err != nil {
		ctx.Logger.Error("Failed to create initial post", "error", err, "topic", topic.Id)
		InternalServerError(ctx)
		return
	}

	shouldNotify := ctx.Request.FormValue("notify") != ""
	updateForumSubscription(ctx, topic.Id, shouldNotify)

	// Broadcast to activity feed (discord, #announce, profile, ...)
	go broadcastForumTopicActivity(ctx, forum, topic, post)

	ctx.Logger.Info(
		"Created a new forum topic",
		"user", ctx.CurrentUser.Id, "topic", topic.Id, "title", topic.Title,
	)
	ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/forum/%d/t/%d", forum.Id, topic.Id))
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

func fetchKudosuForPosts(postIds []int, linkedBeatmapset *schemas.Beatmapset, ctx *server.Context) (map[int64]int, map[int64]*schemas.BeatmapModding) {
	kudosuTotals := map[int64]int{}
	latestKudosu := map[int64]*schemas.BeatmapModding{}
	if linkedBeatmapset == nil {
		return kudosuTotals, latestKudosu
	}

	mods, err := ctx.State.Repositories.Modding.FetchByPosts(postIds)
	if err != nil {
		ctx.Logger.Error("Failed to fetch post kudosu", "error", err)
		return kudosuTotals, latestKudosu
	}

	for _, mod := range mods {
		postId := mod.PostId
		kudosuTotals[postId] += mod.Amount
		if _, exists := latestKudosu[postId]; !exists {
			latestKudosu[postId] = mod
		}
	}

	return kudosuTotals, latestKudosu
}

func resolveSubmittedIcon(ctx *server.Context, canEdit bool) *constants.ForumIcon {
	if !canEdit {
		return nil
	}
	iconId, err := ctx.FormValueInt("icon")
	if err != nil || iconId < 0 {
		return nil
	}
	icon := constants.ForumIcon(iconId)
	return &icon
}
