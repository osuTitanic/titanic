package routes

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/internal/state"
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

func buildForumJumpView(ctx *server.Context, currentForumId int) templates.ForumJumpView {
	forums, err := ctx.State.Forums.FetchAllVisible()
	if err != nil {
		ctx.Logger.Error("Failed to fetch forums for jump box", "error", err)
		return templates.ForumJumpView{CurrentForumId: currentForumId}
	}

	return templates.ForumJumpView{
		CurrentForumId: currentForumId,
		Options:        buildForumJumpOptions(forums),
	}
}

func buildForumJumpOptions(forums []*schemas.Forum) []templates.ForumJumpOption {
	// Build a tree of forums to represent the hierarchy
	children := make(map[int][]*schemas.Forum)
	roots := make([]*schemas.Forum, 0)
	for _, forum := range forums {
		if forum.ParentId == nil {
			roots = append(roots, forum)
			continue
		}
		children[*forum.ParentId] = append(children[*forum.ParentId], forum)
	}

	options := make([]templates.ForumJumpOption, 0, len(forums))
	visited := make(map[int]bool, len(forums))

	// Recursive function to append forums and their children to the options list
	var appendBranch func(*schemas.Forum, int)
	appendBranch = func(forum *schemas.Forum, depth int) {
		if visited[forum.Id] {
			return
		}
		visited[forum.Id] = true

		options = append(options, templates.ForumJumpOption{
			Id:    forum.Id,
			Label: strings.Repeat("\u00a0 \u00a0", depth) + forum.Name,
		})
		for _, child := range children[forum.Id] {
			appendBranch(child, depth+1)
		}
	}
	for _, root := range roots {
		appendBranch(root, 0)
	}
	return options
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

func canMoveForumTopic(ctx *server.Context) bool {
	return ctx.CurrentUser != nil && ctx.HasPermission("forum.moderation.topics.move")
}

func resolveForumMoveTarget(ctx *server.Context, topic *schemas.ForumTopic, eligible bool) (*schemas.Forum, bool) {
	if !eligible || !canMoveForumTopic(ctx) {
		return nil, true
	}
	if strings.TrimSpace(ctx.FormValue("forum_id")) == "" {
		return nil, true
	}

	forumId, err := ctx.FormValueInt("forum_id")
	if err != nil {
		RenderError(ctx, http.StatusBadRequest, "Invalid Forum", "Please select a valid destination forum.")
		return nil, false
	}
	if forumId == topic.ForumId {
		// We already are in that forum, no need to move anything
		return nil, true
	}

	forum, err := ctx.State.Forums.ById(forumId)
	if err != nil {
		ctx.Logger.Error("Failed to fetch forum move target", "error", err, "forum", forumId, "topic", topic.Id)
		InternalServerError(ctx)
		return nil, false
	}
	if forum == nil || forum.Hidden {
		RenderError(ctx, http.StatusBadRequest, "Invalid Forum", "Please select a valid destination forum.")
		return nil, false
	}
	return forum, true
}

func moveForumTopic(ctx *server.Context, topic *schemas.ForumTopic, target *schemas.Forum) bool {
	if target == nil {
		return true
	}
	previousForumId := topic.ForumId

	// Use a database transaction to ensure that both the topic and its posts are moved at the same time
	err := ctx.State.DatabaseTransaction(func(repositories *state.Repositories) error {
		topicUpdate := &schemas.ForumTopic{Id: topic.Id, ForumId: target.Id}
		if _, err := repositories.ForumTopics.Update(topicUpdate, "forum_id"); err != nil {
			return fmt.Errorf("update topic forum: %w", err)
		}

		postUpdate := &schemas.ForumPost{TopicId: topic.Id, ForumId: target.Id}
		if _, err := repositories.ForumPosts.UpdateByTopic(postUpdate, "forum_id"); err != nil {
			return fmt.Errorf("update post forums: %w", err)
		}
		return nil
	})
	if err != nil {
		ctx.Logger.Error("Failed to move forum topic", "error", err, "topic", topic.Id, "forum", target.Id)
		InternalServerError(ctx)
		return false
	}

	topic.ForumId = target.Id
	topic.Forum = target
	ctx.Logger.Info(
		"Moved a forum topic",
		"user", ctx.CurrentUser.Id, "topic", topic.Id,
		"from_forum", previousForumId, "to_forum", target.Id,
	)
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
