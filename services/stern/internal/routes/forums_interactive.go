package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/osuTitanic/titanic-go/internal/schemas"
	"github.com/osuTitanic/titanic-go/services/stern/internal/helpers"
	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
)

const (
	forumActionPost  = "post"
	forumActionEdit  = "edit"
	forumActionQuote = "quote"
)

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
	go helpers.BroadcastForumTopicActivity(ctx, forum, topic, post)

	ctx.Logger.Info(
		"Created a new forum topic",
		"user", ctx.CurrentUser.Id, "topic", topic.Id, "title", topic.Title,
	)
	ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/forum/%d/t/%d", forum.Id, topic.Id))
}

func ForumPostEditorView(ctx *server.Context) {
	if !ctx.RequireLogin() {
		return
	}

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
		ctx.Redirect(http.StatusFound, fmt.Sprintf("/forum/%d/t/%d/post", topic.ForumId, topic.Id))
		return
	}

	action := ctx.QueryValueDefault("action", forumActionPost)
	if action != forumActionPost && action != forumActionEdit && action != forumActionQuote {
		NotFound(ctx)
		return
	}
	actionId, _ := ctx.QueryValueInt64("id")

	// The post being edited or quoted, if any
	var referencedPost *schemas.ForumPost
	if action == forumActionEdit || action == forumActionQuote {
		referencedPost, err = ctx.State.ForumPosts.ById(actionId, "User")
		if err != nil {
			referencedPost = nil
		}
	}

	initialPost, _ := ctx.State.ForumPosts.FetchInitialByTopic(topic.Id)
	latestPost, _ := ctx.State.ForumPosts.FetchLastByTopic(topic.Id)
	beatmapset, err := ctx.State.Beatmapsets.ByTopicId(topic.Id)
	if err != nil {
		beatmapset = nil
	}

	editingInitialPost := action == forumActionEdit && initialPost != nil && initialPost.Id == actionId
	editingLatestPost := action == forumActionEdit && latestPost != nil && latestPost.Id == actionId

	isSubscribed, _ := ctx.State.ForumSubscribers.Exists(topic.Id, ctx.CurrentUser.Id)
	content := resolveEditorContent(ctx, topic, action, referencedPost)
	canEditIcon := canEditPostIcon(ctx, topic, action, editingLatestPost)

	selectedIcon := -1
	if topic.IconId != nil {
		selectedIcon = int(*topic.IconId)
	}

	editor := templates.ForumEditorContext{
		Content:          content,
		SubmitText:       strings.Title(action),
		CancelUrl:        fmt.Sprintf("/forum/%d/t/%d/", topic.ForumId, topic.Id),
		FocusBody:        true,
		Subject:          topic.Title,
		ShowSubject:      editingInitialPost,
		ShowIcons:        canEditIcon,
		NoneIconSelected: topic.IconId == nil,
		Icons:            buildEditorIcons(fetchForumIcons(ctx), selectedIcon),
		ShowControls:     true,
		StatusText:       topic.StatusTextValue(),
		NotifyChecked:    isSubscribed,
		TopicLocked:      topic.LockedAt != nil,
		PostLocked:       referencedPost != nil && referencedPost.EditLocked,
		ShowStatusInput:  action != forumActionEdit && ctx.HasPermission("forum.moderation.topics.set_status"),
		ShowLockTopic:    action != forumActionEdit && ctx.HasPermission("forum.moderation.topics.lock"),
		ShowLockPost:     action == forumActionEdit && ctx.HasPermission("forum.moderation.posts.lock"),
		ShowTopicTypes:   editingInitialPost && ctx.HasPermission("forum.moderation.topics.set_options"),
		TopicType:        topicTypeString(topic),
	}
	if action == forumActionPost {
		editor.DraftUrl = fmt.Sprintf("/forum/%d/t/%d/draft", topic.ForumId, topic.Id)
	}
	applyKudosuHint(ctx, &editor, topic, beatmapset)

	view := templates.ForumPostEditorView{
		DefaultView: buildDefaultViewWithPermissions(ctx),
		Forum:       topic.Forum,
		Topic:       topic,
		Parents:     fetchForumParents(ctx, topic.Forum),
		Editor:      editor,
		Action:      action,
		ActionId:    actionId,
	}
	ctx.RenderTemplate(http.StatusOK, "pages/forum/post", view)
}

func ForumPostAction(ctx *server.Context) {
	if !ctx.RequireLogin() {
		return
	}

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
	if topic == nil || topic.Hidden || topic.ForumId != forumId {
		NotFound(ctx)
		return
	}

	if valid, err := ctx.ValidateCSRF(); err != nil || !valid {
		RenderErrorPage(ctx, http.StatusForbidden, "Invalid Request", "Your session has expired, please try again.")
		return
	}

	if isPostingRejected(ctx) {
		return
	}

	if ctx.Request.FormValue("action") == forumActionEdit {
		handleForumPostEdit(ctx, topic)
		return
	}
	handleForumReply(ctx, topic)
}

func ForumDraftAction(ctx *server.Context) {
	if !ctx.RequireLogin() {
		return
	}

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

	topic, err := ctx.State.ForumTopics.ById(topicId)
	if err != nil {
		ctx.Logger.Error("Failed to fetch topic", "error", err, "topic", topicId)
		InternalServerError(ctx)
		return
	}
	if topic == nil || topic.Hidden || topic.ForumId != forumId {
		NotFound(ctx)
		return
	}

	if valid, err := ctx.ValidateCSRF(); err != nil || !valid {
		RenderErrorPage(ctx, http.StatusForbidden, "Invalid Request", "Your session has expired, please try again.")
		return
	}

	if isPostingRejected(ctx) {
		return
	}

	if !ctx.HasPermission("forum.posts.create") {
		RenderErrorPage(ctx, http.StatusForbidden, "Forbidden", "You are not allowed to post here.")
		return
	}

	topicUrl := fmt.Sprintf("/forum/%d/t/%d/", topic.ForumId, topic.Id)
	content := ctx.Request.FormValue("bbcode")
	if content == "" {
		ctx.Redirect(http.StatusSeeOther, topicUrl)
		return
	}

	// Only a single draft is kept per user & topic
	if drafts, _ := ctx.State.ForumPosts.FetchDrafts(ctx.CurrentUser.Id, topic.Id); len(drafts) > 0 {
		for _, draft := range drafts {
			if err := ctx.State.ForumPosts.Delete(draft); err != nil {
				ctx.Logger.Warn("Failed to delete old draft", "error", err, "draft", draft.Id)
			}
		}
	}

	draft := &schemas.ForumPost{
		TopicId:   topic.Id,
		ForumId:   topic.ForumId,
		UserId:    ctx.CurrentUser.Id,
		Content:   content,
		Draft:     true,
		Hidden:    true,
		CreatedAt: time.Now().UTC(),
	}
	if err := ctx.State.ForumPosts.Create(draft); err != nil {
		ctx.Logger.Error("Failed to save draft", "error", err, "topic", topic.Id)
		InternalServerError(ctx)
		return
	}

	ctx.Logger.Info("Saved a forum draft", "user", ctx.CurrentUser.Id, "topic", topic.Id, "draft", draft.Id)
	ctx.Redirect(http.StatusSeeOther, topicUrl)
}

func handleForumReply(ctx *server.Context, topic *schemas.ForumTopic) {
	if !ctx.HasPermission("forum.posts.create") {
		RenderErrorPage(ctx, http.StatusForbidden, "Forbidden", "You are not allowed to post here.")
		return
	}
	if topic.LockedAt != nil && !ctx.HasPermission("forum.moderation.topics.bypass_lock") {
		RenderErrorPage(ctx, http.StatusForbidden, "Topic Locked", "The topic you are trying to post in is locked!")
		return
	}

	// Prevent accidental duplicate submissions
	// This can happen if the user clicks the submit button multiple times, somehow ...
	if last, _ := ctx.State.ForumPosts.FetchLastByUserInTopic(topic.Id, ctx.CurrentUser.Id); last != nil {
		delta := time.Since(last.CreatedAt)
		if delta <= 2*time.Second {
			// Most likely a double submit -> land on the existing post
			ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/forum/%d/t/%d/p/%d", topic.ForumId, topic.Id, last.Id))
			return
		}
		if delta < 8*time.Second {
			// Let the user know to chill out a lil
			RenderErrorPage(ctx, http.StatusTooManyRequests, "Slow Down!", "You are posting too quickly, slow down!")
			return
		}
	}

	content := ctx.Request.FormValue("bbcode")
	if content == "" {
		ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/forum/%d/t/%d", topic.ForumId, topic.Id))
		return
	}

	post := &schemas.ForumPost{
		TopicId:   topic.Id,
		ForumId:   topic.ForumId,
		UserId:    ctx.CurrentUser.Id,
		Content:   content,
		CreatedAt: time.Now().UTC(),
	}
	// TODO: we eventually want to migrate to timezone-aware timestamps
	// 	     unfortunately i made the stupid decision to use non-timezoned timestamps
	//       for literally everything inside the database, which sucks

	wasChanged, iconUpdate := resolveTopicIconChange(ctx, topic)
	if wasChanged {
		post.IconId = iconUpdate
	}

	if err := ctx.State.ForumPosts.Create(post); err != nil {
		ctx.Logger.Error("Failed to create reply", "error", err, "topic", topic.Id)
		InternalServerError(ctx)
		return
	}

	notify := ctx.Request.FormValue("notify") != ""
	notifyForumSubscribers(ctx, topic, post)
	updateForumSubscription(ctx, topic.Id, notify)

	// Assemble the topic updates caused by this reply
	topicUpdates := &schemas.ForumTopic{Id: topic.Id, LastPostAt: time.Now().UTC()}
	columns := []string{"last_post_at"}

	if wasChanged {
		topicUpdates.IconId = iconUpdate
		columns = append(columns, "icon")
	}
	if ctx.HasPermission("forum.moderation.topics.lock") {
		if ctx.Request.FormValue("locked") != "" {
			now := time.Now()
			topicUpdates.LockedAt = &now
		}
		columns = append(columns, "locked_at")
	}

	statusWasSet := false
	if ctx.HasPermission("forum.moderation.topics.set_status") {
		status := strings.TrimSpace(ctx.FormValue("topic-status"))

		// Check if the status text has changed, and if so, update it
		if status != topic.StatusTextValue() {
			if status != "" {
				topicUpdates.StatusText = &status
				statusWasSet = true
			}
			columns = append(columns, "status_text")
		}
	}

	if _, err := ctx.State.ForumTopics.Update(topicUpdates, columns...); err != nil {
		ctx.Logger.Error("Failed to update topic after reply", "error", err, "topic", topic.Id)
	}

	// Update the status text of the beatmap topic unless it was set manually
	// A beatmap topic status can be e.g. "Needs modding", "Waiting for approval...", ...
	if !statusWasSet {
		if beatmapset, _ := ctx.State.Beatmapsets.ByTopicId(topic.Id); beatmapset != nil {
			updateBeatmapTopicStatus(ctx, topic, beatmapset)
		}
	}

	// Broadcast to activity feed (discord, #announce, profile, ...)
	go helpers.BroadcastForumPostActivity(ctx, topic, post)

	ctx.Logger.Info("Created a forum post", "user", ctx.CurrentUser.Id, "topic", topic.Id, "post", post.Id)
	ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/forum/%d/t/%d/p/%d", topic.ForumId, topic.Id, post.Id))
}

func handleForumPostEdit(ctx *server.Context, topic *schemas.ForumTopic) {
	if !ctx.HasPermission("forum.posts.edit") {
		RenderErrorPage(ctx, http.StatusForbidden, "Forbidden", "You are not allowed to edit posts.")
		return
	}

	if topic.LockedAt != nil && !ctx.HasPermission("forum.moderation.topics.bypass_lock") {
		RenderErrorPage(ctx, http.StatusForbidden, "Topic Locked", "The topic you are trying to post in is locked!")
		return
	}

	postId, _ := strconv.ParseInt(strings.TrimSpace(ctx.Request.FormValue("id")), 10, 64)
	post, err := ctx.State.ForumPosts.ById(postId)
	if err != nil || post == nil {
		RenderErrorPage(ctx, http.StatusNotFound, "Post Not Found", "The post you are trying to edit could not be found.")
		return
	}

	if post.EditLocked && !ctx.HasPermission("forum.moderation.posts.bypass_lock") {
		RenderErrorPage(ctx, http.StatusForbidden, "Post Locked", "The post you are trying to edit is locked!")
		return
	}

	isOwnPost := post.UserId == ctx.CurrentUser.Id
	if !isOwnPost && !ctx.HasPermission("forum.moderation.posts.edit") {
		RenderErrorPage(ctx, http.StatusForbidden, "Forbidden", "You are not allowed to edit this post.")
		return
	}

	content := ctx.Request.FormValue("bbcode")
	if content == "" {
		ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/forum/%d/t/%d/p/%d", topic.ForumId, topic.Id, post.Id))
		return
	}

	notify := ctx.Request.FormValue("notify") != ""
	updateForumSubscription(ctx, topic.Id, notify)

	initialPost, _ := ctx.State.ForumPosts.FetchInitialByTopic(topic.Id)
	latestPost, _ := ctx.State.ForumPosts.FetchLastByTopic(topic.Id)
	editingInitialPost := initialPost != nil && initialPost.Id == post.Id
	editingLatestPost := latestPost != nil && latestPost.Id == post.Id

	updates := &schemas.ForumPost{Id: post.Id, Content: content}
	columns := []string{"content"}

	// Icons may only be changed while editing the latest post
	if editingLatestPost {
		if wasChanged, iconUpdate := resolveTopicIconChange(ctx, topic); wasChanged {
			updates.IconId = iconUpdate
			columns = append(columns, "icon_id")

			topicIcon := &schemas.ForumTopic{Id: topic.Id, IconId: iconUpdate}
			if _, err := ctx.State.ForumTopics.Update(topicIcon, "icon"); err != nil {
				ctx.Logger.Error("Failed to update topic icon", "error", err, "topic", topic.Id)
			}
		}
	}

	if ctx.HasPermission("forum.moderation.posts.lock") {
		updates.EditLocked = ctx.Request.FormValue("edit-locked") != ""
		columns = append(columns, "edit_locked")
	}

	if isOwnPost {
		updates.EditCount = post.EditCount + 1
		updates.EditTime = time.Now()
		columns = append(columns, "edit_count", "edit_time")
	}

	if _, err := ctx.State.ForumPosts.Update(updates, columns...); err != nil {
		ctx.Logger.Error("Failed to update post", "error", err, "post", post.Id)
		InternalServerError(ctx)
		return
	}

	// If allowed, update the topic type & title when the user edits the initial post
	applyEditedTopicOptions(ctx, topic, editingInitialPost)

	ctx.Logger.Info("Edited a forum post", "user", ctx.CurrentUser.Id, "post", post.Id)
	ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/forum/%d/t/%d/p/%d", topic.ForumId, topic.Id, post.Id))
}
