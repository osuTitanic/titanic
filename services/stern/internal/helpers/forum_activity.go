package helpers

import (
	"encoding/json"

	"github.com/osuTitanic/titanic-go/internal/constants"
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
)

func BroadcastForumTopicActivity(ctx *server.Context, forum *schemas.Forum, topic *schemas.ForumTopic, post *schemas.ForumPost) {
	createForumActivity(ctx, constants.ActivityForumTopicCreated, false, map[string]any{
		"username":   ctx.CurrentUser.Name,
		"topic_name": topic.Title,
		"forum_name": forum.Name,
		"forum_id":   topic.ForumId,
		"topic_id":   topic.Id,
		"topic_icon": forumTopicIconLocation(topic),
		"content":    truncateActivityContent(post.Content),
	})
}

func BroadcastForumPostActivity(ctx *server.Context, topic *schemas.ForumTopic, post *schemas.ForumPost) {
	createForumActivity(ctx, constants.ActivityForumPostCreated, true, map[string]any{
		"username":   ctx.CurrentUser.Name,
		"post_id":    post.Id,
		"topic_name": topic.Title,
		"topic_id":   topic.Id,
		"topic_icon": forumTopicIconLocation(topic),
		"content":    truncateActivityContent(post.Content),
	})
}

func createForumActivity(ctx *server.Context, activityType constants.UserActivity, hidden bool, data map[string]any) {
	payload, err := json.Marshal(data)
	if err != nil {
		ctx.Logger.Error("Failed to encode forum activity", "error", err)
		return
	}

	activity := &schemas.Activity{
		UserId: ctx.CurrentUser.Id,
		Type:   int(activityType),
		Data:   payload,
		Hidden: hidden,
	}
	if err := ctx.State.Activities.Create(activity); err != nil {
		ctx.Logger.Warn("Failed to record forum activity", "error", err, "type", activityType)
	}
	// TODO: Add redis event queue broadcast
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
