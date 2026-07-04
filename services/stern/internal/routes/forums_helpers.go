package routes

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/osuTitanic/titanic-go/internal/constants"
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
)

const forumTitleMaxLength = 128
const forumPostMaxLength = 1 << 14

var (
	quoteStripPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?s)\[quote(?:=[^\]]+)?\].*?\[/quote\]`),
		regexp.MustCompile(`(?s)\[img\].*?\[/img\]`),
		regexp.MustCompile(`(?s)\[(?:video|youtube)\].*?\[/(?:video|youtube)\]`),
	}
	quoteBlankLines = regexp.MustCompile(`(?:\r\n){2,}`)
)

func isPostingRejected(ctx *server.Context) bool {
	if ctx.CurrentUser == nil {
		return false
	}
	if ctx.CurrentUser.SilenceEnd != nil && time.Now().Before(*ctx.CurrentUser.SilenceEnd) {
		RenderErrorPage(ctx, http.StatusForbidden, "You are silenced!", "You cannot post while you are silenced.")
		return true
	}
	if ctx.CurrentUser.Restricted {
		RenderErrorPage(ctx, http.StatusForbidden, "You are restricted!", "You cannot post while your account is restricted.")
		return true
	}

	canBypassLength := ctx.HasPermission("forum.moderation.posts.bypass_length")
	postLength := len(ctx.Request.FormValue("bbcode"))
	if postLength > forumPostMaxLength && !canBypassLength {
		RenderErrorPage(ctx, http.StatusForbidden, "Post too long", "Please limit your post to 15000 characters or less.")
		return true
	}

	titleLength := len(ctx.Request.FormValue("title"))
	if titleLength > forumTitleMaxLength && !canBypassLength {
		RenderErrorPage(ctx, http.StatusForbidden, "Title too long", "Please limit your title to 128 characters or less.")
		return true
	}

	return false
}

func fetchForumIcons(ctx *server.Context) []*schemas.ForumIcon {
	icons, err := ctx.State.ForumIcons.FetchAll()
	if err != nil {
		ctx.Logger.Error("Failed to fetch forum icons", "error", err)
		return nil
	}
	return icons
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

func canEditForumIcon(ctx *server.Context, allowIcons bool) bool {
	if ctx.CurrentUser == nil {
		return false
	}
	if ctx.HasPermission("forum.moderation.topics.edit_icon") {
		return true
	}
	return allowIcons && ctx.HasPermission("forum.topics.edit_icon")
}

func canEditPostIcon(ctx *server.Context, topic *schemas.ForumTopic, action string, editingLatestPost bool) bool {
	if action == forumActionEdit && !editingLatestPost {
		return false
	}
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

func buildEditorIcons(icons []*schemas.ForumIcon, selectedId int) []*templates.ForumEditorIcon {
	result := make([]*templates.ForumEditorIcon, 0, len(icons))
	for _, icon := range icons {
		result = append(result, &templates.ForumEditorIcon{
			Id:       int(icon.Id),
			Name:     icon.Name,
			Location: icon.Location,
			Selected: int(icon.Id) == selectedId,
		})
	}
	return result
}

func resolveEditorContent(ctx *server.Context, topic *schemas.ForumTopic, action string, post *schemas.ForumPost) string {
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
	if !editingInitialPost {
		return
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

	if title := strings.TrimSpace(ctx.FormValue("title")); title != "" {
		update := &schemas.ForumTopic{
			Id:    topic.Id,
			Title: title,
		}
		if _, err := ctx.State.ForumTopics.Update(update, "title"); err != nil {
			ctx.Logger.Error("Failed to update topic title", "error", err, "topic", topic.Id)
		}
	}
}

func resolveTopicType(ctx *server.Context) (pinned, announcement bool) {
	if !ctx.HasPermission("forum.moderation.topics.set_options") {
		return false, false
	}

	switch ctx.Request.FormValue("type") {
	case "announcement":
		return false, true
	case "pinned":
		return true, false
	default:
		return false, false
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

func updateForumSubscription(ctx *server.Context, topicId int, notify bool) {
	subscriber := &schemas.ForumSubscriber{
		UserId:  ctx.CurrentUser.Id,
		TopicId: topicId,
	}
	if !notify {
		if err := ctx.State.ForumSubscribers.Delete(subscriber); err != nil {
			ctx.Logger.Warn("Failed to remove forum subscription", "error", err, "topic", topicId)
		}
		return
	}

	if exists, _ := ctx.State.ForumSubscribers.Exists(topicId, ctx.CurrentUser.Id); exists {
		return
	}
	if err := ctx.State.ForumSubscribers.Create(subscriber); err != nil {
		ctx.Logger.Warn("Failed to add forum subscription", "error", err, "topic", topicId)
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
	}
}

func updateBeatmapTopicStatus(ctx *server.Context, topic *schemas.ForumTopic, beatmapset *schemas.Beatmapset) {
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

	editor.ShowKudosuHint = true
	editor.BeatmapsetId = beatmapset.Id
	editor.KudosuReward = 2
	if time.Since(topic.LastPostAt) < 7*24*time.Hour {
		editor.KudosuReward = 1
	}
	editor.ShowKudosuIconNote = canEditForumIcon(ctx, topic.CanChangeIcon)
}
