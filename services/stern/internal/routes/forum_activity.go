package routes

import (
	"github.com/osuTitanic/titanic/internal/activity"
	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/services/stern/internal/server"
)

func broadcastForumTopicActivity(ctx *server.Context, forum *schemas.Forum, topic *schemas.ForumTopic, post *schemas.ForumPost) {
	err := activity.Submit(
		ctx.State, ctx.CurrentUser.Id, nil,
		constants.ActivityForumTopicCreated,
		map[string]any{
			"username":   ctx.CurrentUser.Name,
			"topic_name": topic.Title,
			"forum_name": forum.Name,
			"forum_id":   topic.ForumId,
			"topic_id":   topic.Id,
			"topic_icon": forumTopicIconLocation(topic),
			"content":    truncateActivityContent(post.Content),
		},
		true,
		false,
	)
	if err != nil {
		ctx.Logger.Warn("Failed to record forum activity", "error", err, "type", constants.ActivityForumTopicCreated)
	}
}

func broadcastForumPostActivity(ctx *server.Context, topic *schemas.ForumTopic, post *schemas.ForumPost) {
	err := activity.Submit(
		ctx.State, ctx.CurrentUser.Id, nil,
		constants.ActivityForumPostCreated,
		map[string]any{
			"username":   ctx.CurrentUser.Name,
			"post_id":    post.Id,
			"topic_name": topic.Title,
			"topic_id":   topic.Id,
			"topic_icon": forumTopicIconLocation(topic),
			"content":    truncateActivityContent(post.Content),
		},
		true,
		true,
	)
	if err != nil {
		ctx.Logger.Warn("Failed to record forum activity", "error", err, "type", constants.ActivityForumPostCreated)
	}
}

func truncateActivityContent(content string) string {
	runes := []rune(content)
	truncated := content
	if len(runes) > 512 {
		truncated = string(runes[:512])
	}
	if len(runes) > 1024 {
		truncated += "..."
	}
	return truncated
}

func forumTopicIconLocation(topic *schemas.ForumTopic) string {
	if topic.Icon != nil {
		return topic.Icon.Location
	}
	return ""
}
