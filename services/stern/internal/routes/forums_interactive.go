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

const forumTitleMaxLength = 128
const forumPostMaxLength = 1 << 14

func ForumCreateTopicView(ctx *server.Context) {
	if !requireLogin(ctx) {
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
	if !requireLogin(ctx) {
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

	ctx.Logger.Info(
		"Created a new forum topic",
		"user", ctx.CurrentUser.Id, "topic", topic.Id, "title", topic.Title,
	)
	ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/forum/%d/t/%d", forum.Id, topic.Id))
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

func requireLogin(ctx *server.Context) bool {
	if ctx.IsAuthenticated() {
		return true
	}
	ctx.Redirect(
		http.StatusSeeOther,
		"/account/login?redirect="+ctx.Request.URL.RequestURI(),
	)
	return false
}

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

func canEditForumIcon(ctx *server.Context, allowIcons bool) bool {
	if ctx.CurrentUser == nil {
		return false
	}
	if ctx.HasPermission("forum.moderation.topics.edit_icon") {
		return true
	}
	return allowIcons && ctx.HasPermission("forum.topics.edit_icon")
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

func fetchForumIcons(ctx *server.Context) []*schemas.ForumIcon {
	icons, err := ctx.State.ForumIcons.FetchAll()
	if err != nil {
		ctx.Logger.Error("Failed to fetch forum icons", "error", err)
		return nil
	}
	return icons
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
