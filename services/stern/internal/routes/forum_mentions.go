package routes

import (
	"fmt"

	"github.com/osuTitanic/titanic/internal/bbcode"
	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/services/stern/internal/server"
)

func notifyForumMentions(ctx *server.Context, topic *schemas.ForumTopic, post *schemas.ForumPost, previousContent string) {
	// Collect previous user IDs to avoid sending duplicate notifications
	previousUserIds := resolveMentionedUserIds(ctx, previousContent)

	for userId := range resolveMentionedUserIds(ctx, post.Content) {
		if userId == ctx.CurrentUser.Id {
			continue
		}
		if _, alreadyMentioned := previousUserIds[userId]; alreadyMentioned {
			continue
		}

		// We don't want to notify users who have blocked the user
		blocked, err := ctx.State.Relationships.IsBlockedBetween(ctx.CurrentUser.Id, userId)
		if err != nil {
			ctx.Logger.Warn("Failed to check mention block status", "error", err, "user", userId)
			continue
		}
		if blocked {
			continue
		}

		notification := &schemas.Notification{
			UserId:  userId,
			Type:    constants.NotificationTypeForum,
			Header:  "You got mentioned",
			Content: fmt.Sprintf("%s mentioned you in \"%s\". Click here to view it!", ctx.CurrentUser.Name, topic.Title),
			Link:    fmt.Sprintf("/forum/%d/p/%d", topic.ForumId, post.Id),
		}
		if err := ctx.State.Notifications.Create(notification); err != nil {
			ctx.Logger.Warn("Failed to create forum mention notification", "error", err, "user", userId)
		}
	}
}

func resolveMentionedUserIds(ctx *server.Context, content string) map[int]struct{} {
	mentionedUsernames := bbcode.MentionedUsernames(content)
	userIds, err := ctx.State.Users.GetUserIdsCaseInsensitive(mentionedUsernames)
	if err != nil {
		ctx.Logger.Warn("Failed to resolve mentioned users", "error", err)
		return map[int]struct{}{}
	}

	result := make(map[int]struct{})
	for _, userId := range userIds {
		result[userId] = struct{}{}
	}
	return result
}
