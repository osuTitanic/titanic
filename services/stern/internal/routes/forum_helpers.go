package routes

import (
	"net/http"
	"slices"
	"time"

	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/services/stern/internal/server"
	"github.com/osuTitanic/titanic/services/stern/internal/templates"
)

const forumTitleMaxLength = 128
const forumPostMaxLength = 1 << 14

func isPostingRejected(ctx *server.Context) bool {
	if ctx.CurrentUser == nil {
		return false
	}
	if ctx.CurrentUser.SilenceEnd != nil && time.Now().Before(*ctx.CurrentUser.SilenceEnd) {
		UserSilenced(ctx)
		return true
	}
	if ctx.CurrentUser.Restricted {
		UserRestricted(ctx)
		return true
	}

	canBypassLength := ctx.HasPermission("forum.moderation.posts.bypass_length")
	postLength := len(ctx.Request.FormValue("bbcode"))
	if postLength > forumPostMaxLength && !canBypassLength {
		PostTooLong(ctx)
		return true
	}

	titleLength := len(ctx.Request.FormValue("title"))
	if titleLength > forumTitleMaxLength && !canBypassLength {
		RenderError(ctx, http.StatusForbidden, "Title too long", "Please limit your title to 128 characters or less.")
		return true
	}

	return false
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
	activeIds := forumGetActiveUsers(ctx, forumId)
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

func fetchForumIcons(ctx *server.Context) []*schemas.ForumIcon {
	icons, err := ctx.State.ForumIcons.FetchAll()
	if err != nil {
		ctx.Logger.Error("Failed to fetch forum icons", "error", err)
		return nil
	}
	return icons
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
