package routes

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/osuTitanic/titanic-go/internal/constants"
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
)

const (
	forumActionPost  = "post"
	forumActionEdit  = "edit"
	forumActionQuote = "quote"
)

var (
	quoteStripPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?s)\[quote(?:=[^\]]+)?\].*?\[/quote\]`),
		regexp.MustCompile(`(?s)\[img\].*?\[/img\]`),
		regexp.MustCompile(`(?s)\[(?:video|youtube)\].*?\[/(?:video|youtube)\]`),
	}
	quoteBlankLines = regexp.MustCompile(`(?:\r\n){2,}`)
)

func ForumPostEditorView(ctx *server.Context) {
	if !ctx.RequireLogin() {
		return
	}

	forumId, err := ctx.PathValueInt("id")
	if err != nil {
		ForumNotFound(ctx)
		return
	}

	topicId, err := ctx.PathValueInt("topicId")
	if err != nil {
		TopicNotFound(ctx)
		return
	}

	topic, err := ctx.State.ForumTopics.ById(topicId, "Forum", "Icon")
	if err != nil {
		ctx.Logger.Error("Failed to fetch topic", "error", err, "topic", topicId)
		InternalServerError(ctx)
		return
	}
	if topic == nil || topic.Hidden {
		TopicNotFound(ctx)
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
		ShowStatusInput:  (action != forumActionEdit || editingInitialPost) && ctx.HasPermission("forum.moderation.topics.set_status"),
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
		ForumNotFound(ctx)
		return
	}

	topicId, err := ctx.PathValueInt("topicId")
	if err != nil {
		TopicNotFound(ctx)
		return
	}

	topic, err := ctx.State.ForumTopics.ById(topicId, "Forum", "Icon")
	if err != nil {
		ctx.Logger.Error("Failed to fetch topic", "error", err, "topic", topicId)
		InternalServerError(ctx)
		return
	}
	if topic == nil || topic.Hidden || topic.ForumId != forumId {
		TopicNotFound(ctx)
		return
	}

	if valid, err := ctx.ValidateCSRF(); err != nil || !valid {
		RenderError(ctx, http.StatusForbidden, "Invalid Request", "Your session has expired, please try again.")
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
		ForumNotFound(ctx)
		return
	}

	topicId, err := ctx.PathValueInt("topicId")
	if err != nil {
		TopicNotFound(ctx)
		return
	}

	topic, err := ctx.State.ForumTopics.ById(topicId)
	if err != nil {
		ctx.Logger.Error("Failed to fetch topic", "error", err, "topic", topicId)
		InternalServerError(ctx)
		return
	}
	if topic == nil || topic.Hidden || topic.ForumId != forumId {
		TopicNotFound(ctx)
		return
	}

	if valid, err := ctx.ValidateCSRF(); err != nil || !valid {
		RenderError(ctx, http.StatusForbidden, "Invalid Request", "Your session has expired, please try again.")
		return
	}

	if isPostingRejected(ctx) {
		return
	}

	if !ctx.HasPermission("forum.posts.create") {
		RenderError(ctx, http.StatusForbidden, "Forbidden", "You are not allowed to post here.")
		return
	}

	topicUrl := fmt.Sprintf("/forum/%d/t/%d/", topic.ForumId, topic.Id)
	content := ctx.Request.FormValue("bbcode")
	if content == "" {
		ctx.Redirect(http.StatusSeeOther, topicUrl)
		return
	}

	// Only a single draft is kept per user & topic
	clearForumDrafts(ctx, topic.Id)

	draft := &schemas.ForumPost{
		TopicId:   topic.Id,
		ForumId:   topic.ForumId,
		UserId:    ctx.CurrentUser.Id,
		Content:   content,
		Draft:     true,
		Hidden:    true,
		CreatedAt: time.Now(),
	}
	if err := ctx.State.ForumPosts.Create(draft); err != nil {
		ctx.Logger.Error("Failed to save draft", "error", err, "topic", topic.Id)
		InternalServerError(ctx)
		return
	}

	ctx.Logger.Info("Saved a forum draft", "user", ctx.CurrentUser.Id, "topic", topic.Id, "draft", draft.Id)
	ctx.Redirect(http.StatusSeeOther, topicUrl)
}

func clearForumDrafts(ctx *server.Context, topicId int) {
	drafts, err := ctx.State.ForumPosts.FetchDrafts(ctx.CurrentUser.Id, topicId)
	if err != nil {
		ctx.Logger.Warn("Failed to fetch drafts for cleanup", "error", err, "topic", topicId)
		return
	}
	for _, draft := range drafts {
		if err := ctx.State.ForumPosts.Delete(draft); err != nil {
			ctx.Logger.Warn("Failed to delete draft", "error", err, "draft", draft.Id)
		}
	}
}

func handleForumReply(ctx *server.Context, topic *schemas.ForumTopic) {
	if !ctx.HasPermission("forum.posts.create") {
		RenderError(ctx, http.StatusForbidden, "Forbidden", "You are not allowed to post here.")
		return
	}
	if topic.LockedAt != nil && !ctx.HasPermission("forum.moderation.topics.bypass_lock") {
		TopicLocked(ctx)
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
			PostingTooQuickly(ctx)
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
		CreatedAt: time.Now(),
	}

	wasChanged, iconUpdate := resolveTopicIconChange(ctx, topic)
	if wasChanged {
		post.IconId = iconUpdate
	}

	if err := ctx.State.ForumPosts.Create(post); err != nil {
		ctx.Logger.Error("Failed to create reply", "error", err, "topic", topic.Id)
		InternalServerError(ctx)
		return
	}

	// The reply was posted, so any saved drafts for this topic are now obsolete
	clearForumDrafts(ctx, topic.Id)

	notify := ctx.Request.FormValue("notify") != ""
	notifyForumSubscribers(ctx, topic, post)
	updateForumSubscription(ctx, topic.Id, notify)

	// Assemble the topic updates caused by this reply
	topicUpdates := &schemas.ForumTopic{Id: topic.Id, LastPostAt: time.Now()}
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
	go broadcastForumPostActivity(ctx, topic, post)

	ctx.Logger.Info("Created a forum post", "user", ctx.CurrentUser.Id, "topic", topic.Id, "post", post.Id)
	ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/forum/%d/t/%d/p/%d", topic.ForumId, topic.Id, post.Id))
}

func handleForumPostEdit(ctx *server.Context, topic *schemas.ForumTopic) {
	if !ctx.HasPermission("forum.posts.edit") {
		RenderError(ctx, http.StatusForbidden, "Forbidden", "You are not allowed to edit posts.")
		return
	}

	if topic.LockedAt != nil && !ctx.HasPermission("forum.moderation.topics.bypass_lock") {
		TopicLocked(ctx)
		return
	}

	postId, _ := strconv.ParseInt(strings.TrimSpace(ctx.Request.FormValue("id")), 10, 64)
	post, err := ctx.State.ForumPosts.ById(postId)
	if err != nil || post == nil {
		PostNotFound(ctx)
		return
	}

	if post.EditLocked && !ctx.HasPermission("forum.moderation.posts.bypass_lock") {
		PostLocked(ctx)
		return
	}

	isOwnPost := post.UserId == ctx.CurrentUser.Id
	if !isOwnPost && !ctx.HasPermission("forum.moderation.posts.edit") {
		RenderError(ctx, http.StatusForbidden, "Forbidden", "You are not allowed to edit this post.")
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

func canEditPostIcon(ctx *server.Context, topic *schemas.ForumTopic, action string, editingLatestPost bool) bool {
	// We only allow editing the post icon when editing the latest post in the topic.
	// The other way to do it would be to make a new post that changes the icon.
	if action == forumActionEdit && !editingLatestPost {
		return false
	}

	// We only want certain users to be able to edit the post icon.
	// Either we have the permission to change the icon anywhere, or the topic
	// allows for it and we also have the permission to change it in that case.
	if ctx.HasPermission("forum.moderation.topics.edit_icon") {
		return true
	}
	return topic.CanChangeIcon && ctx.HasPermission("forum.topics.edit_icon")
}

func resolveTopicIconChange(ctx *server.Context, topic *schemas.ForumTopic) (bool, *constants.ForumIcon) {
	// Either the user has the permission to change the icon anywhere, or the topic allows
	// for it and the user also has the permission to change it in that case.
	// (we definitively want a better permission system for forums lol)
	canChange := ctx.HasPermission("forum.moderation.topics.edit_icon") ||
		(topic.CanChangeIcon && ctx.HasPermission("forum.topics.edit_icon"))
	if !canChange {
		return false, nil
	}

	iconId, err := ctx.FormValueInt("icon")
	if err != nil {
		iconId = -1
	}

	current := -1
	if topic.IconId != nil {
		current = int(*topic.IconId)
	}

	if iconId == current {
		// Nothing changed
		return false, nil
	}
	if iconId < 0 {
		// Icon was removed (-1)
		return true, nil
	}
	// Icon was changed to a new value
	icon := constants.ForumIcon(iconId)
	return true, &icon
}

func resolveEditorContent(ctx *server.Context, topic *schemas.ForumTopic, action string, post *schemas.ForumPost) string {
	// Determine the placeholder content for the editor
	// For edit: use the post content
	// For quote: use the post content wrapped in a quote tag
	// For new post: use the last saved draft, if any
	switch action {
	case forumActionEdit:
		if post != nil && !post.Deleted {
			return post.Content
		}
	case forumActionQuote:
		if post != nil && !post.Deleted {
			return buildQuote(post)
		}
	case forumActionPost:
		return restoreDraft(ctx, topic.Id)
	}
	return ""
}

func buildQuote(post *schemas.ForumPost) string {
	// The quote is built by stripping out any existing quotes, images, or videos from the post content.
	// To be entirely honest, this is a bit janky right now... It would probably be better to use
	// bbcode.Strip(...), which would remove all bbcode tags.
	content := post.Content
	for _, pattern := range quoteStripPatterns {
		content = pattern.ReplaceAllString(content, "")
	}

	content = quoteBlankLines.ReplaceAllString(content, "\r\n\r\n")

	// Drop the separator in beatmap submission posts
	if parts := strings.Split(content, "---------------\n"); len(parts) > 1 {
		content = parts[len(parts)-1]
	}

	author := ""
	if post.User != nil {
		author = strings.NewReplacer("[", "", "]", "").Replace(post.User.Name)
	}
	return fmt.Sprintf("[quote=%s]%s[/quote]", author, strings.TrimSpace(content))
}

func restoreDraft(ctx *server.Context, topicId int) string {
	drafts, err := ctx.State.ForumPosts.FetchDrafts(ctx.CurrentUser.Id, topicId)
	if err != nil || len(drafts) == 0 {
		return ""
	}
	if err := ctx.State.ForumPosts.Delete(drafts[0]); err != nil {
		ctx.Logger.Warn("Failed to delete restored draft", "error", err, "draft", drafts[0].Id)
	}
	return drafts[0].Content
}

func applyEditedTopicOptions(ctx *server.Context, topic *schemas.ForumTopic, editingInitialPost bool) {
	// Topic "options" refers to the topic type (pinned, announcement, etc.), the topic title and status.
	// We can only set those if we're editing the initial post of the topic, and if we have the appropriate permissions.
	if !editingInitialPost {
		return
	}

	if title := strings.TrimSpace(ctx.FormValue("title")); title != "" {
		update := &schemas.ForumTopic{
			Id:    topic.Id,
			Title: title,
		}
		if _, err := ctx.State.ForumTopics.Update(update, "title"); err != nil {
			ctx.Logger.Error("Failed to update topic title", "error", err, "topic", topic.Id)
		}
	}

	if ctx.HasPermission("forum.moderation.topics.set_options") {
		pinned, announcement := resolveTopicType(ctx)
		update := &schemas.ForumTopic{
			Id:           topic.Id,
			Pinned:       pinned,
			Announcement: announcement,
		}
		if _, err := ctx.State.ForumTopics.Update(update, "pinned", "announcement"); err != nil {
			ctx.Logger.Error("Failed to update topic type", "error", err, "topic", topic.Id)
		}
	}

	if ctx.HasPermission("forum.moderation.topics.set_status") {
		status := strings.TrimSpace(ctx.FormValue("topic-status"))

		// Check if the status text has changed, and if so, update it
		if status != topic.StatusTextValue() {
			update := &schemas.ForumTopic{Id: topic.Id}
			if status != "" {
				update.StatusText = &status
			}
			if _, err := ctx.State.ForumTopics.Update(update, "status_text"); err != nil {
				ctx.Logger.Error("Failed to update topic status", "error", err, "topic", topic.Id)
			}
		}
	}
}

func topicTypeString(topic *schemas.ForumTopic) string {
	switch {
	case topic == nil:
		return "global"
	case topic.Announcement:
		return "announcement"
	case topic.Pinned:
		return "pinned"
	default:
		return "global"
	}
}

func notifyForumSubscribers(ctx *server.Context, topic *schemas.ForumTopic, post *schemas.ForumPost) {
	subscribers, err := ctx.State.ForumSubscribers.FetchByTopic(topic.Id)
	if err != nil {
		ctx.Logger.Warn("Failed to fetch topic subscribers", "error", err, "topic", topic.Id)
		return
	}

	for _, subscriber := range subscribers {
		if subscriber.UserId == ctx.CurrentUser.Id {
			// Don't notify the author of the post
			continue
		}

		notification := &schemas.Notification{
			UserId:  subscriber.UserId,
			Type:    constants.NotificationTypeForum,
			Header:  "New Post",
			Content: fmt.Sprintf("%s posted something in \"%s\". Click here to view it!", ctx.CurrentUser.Name, topic.Title),
			Link:    fmt.Sprintf("/forum/%d/p/%d", topic.ForumId, post.Id),
		}
		if err := ctx.State.Notifications.Create(notification); err != nil {
			ctx.Logger.Warn("Failed to create forum notification", "error", err, "user", subscriber.UserId)
		}
		// TODO: In the future, we should probably have a system to send
		// 		 notifications to multiple sources e.g. email, discord, etc.
	}
}

func updateBeatmapTopicStatus(ctx *server.Context, topic *schemas.ForumTopic, beatmapset *schemas.Beatmapset) {
	// The topic status text is dynamically resolved based on the beatmapset status
	// Not every topic has a linked beatmapset though, so we only do it if we have one
	status := resolveBeatmapTopicStatus(ctx, topic, beatmapset)
	update := &schemas.ForumTopic{Id: topic.Id, StatusText: nil}
	if status != "" {
		update.StatusText = &status
	}
	if _, err := ctx.State.ForumTopics.Update(update, "status_text"); err != nil {
		ctx.Logger.Error("Failed to update beatmap topic status", "error", err, "topic", topic.Id)
	}
}

func resolveBeatmapTopicStatus(ctx *server.Context, topic *schemas.ForumTopic, beatmapset *schemas.Beatmapset) string {
	if beatmapset.IsApproved() || beatmapset.Status == constants.BeatmapStatusGraveyard {
		// Ranked or graveyarded sets carry no status text
		return ""
	}

	if nominations, _ := ctx.State.Nominations.FetchBySet(beatmapset.Id); len(nominations) > 0 {
		// We already have at least one nomination, so the set is waiting for approval
		return "Waiting for approval..."
	}

	lastBatPost, _ := ctx.State.ForumPosts.FetchLastBatPost(topic.Id)
	if lastBatPost == nil {
		// No nominations & no BAT posts -> the set is still waiting for modding
		return "Needs modding"
	}

	var lastCreatorPost *schemas.ForumPost
	if beatmapset.CreatorId != nil {
		lastCreatorPost, _ = ctx.State.ForumPosts.FetchLastByUserInTopic(topic.Id, *beatmapset.CreatorId)
	}
	if lastCreatorPost == nil || lastBatPost.Id > lastCreatorPost.Id {
		// Last post was by a BAT member, so the set is waiting for a response from the creator
		return "Waiting for creator's response..."
	}

	return "Waiting for further modding..."
}

func applyKudosuHint(ctx *server.Context, editor *templates.ForumEditorContext, topic *schemas.ForumTopic, beatmapset *schemas.Beatmapset) {
	if beatmapset == nil || beatmapset.CreatorId == nil {
		return
	}
	if *beatmapset.CreatorId == ctx.CurrentUser.Id || beatmapset.Status >= constants.BeatmapStatusRanked {
		return
	}

	// The kudosu hint is shown to a user if we have a linked beatmapset, stating
	// that you can receive kudosu for posting in the topic. The amount of kudosu
	// is determined by whether the topic has been posted in within the last 7 days or not.
	// This should reward users for modding beatmaps that have been inactive for a while.
	editor.ShowKudosuHint = true
	editor.BeatmapsetId = beatmapset.Id
	editor.KudosuReward = 2
	if time.Since(topic.LastPostAt) < 7*24*time.Hour {
		editor.KudosuReward = 1
	}
	editor.ShowKudosuIconNote = canEditForumIcon(ctx, topic.CanChangeIcon)
}
