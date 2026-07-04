package routes

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/osuTitanic/titanic-go/internal/schemas"
	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/redis/go-redis/v9"
)

const forumReadTimestampExpiry = 14 * 24 * time.Hour
const forumActivityExpiry = 5 * time.Minute
const forumAverageViewsTTL = 5 * time.Minute
const forumViewLockExpiry = time.Minute

// forumSessionIdentifier returns a key used to track
// which topics the current visitor has already read.
func forumSessionIdentifier(ctx *server.Context) string {
	var seed string = ctx.IP()

	if ctx.CurrentUser != nil {
		// If user is authenticated, we'll use the user ID
		// to make read status persist across sessions
		seed = strconv.Itoa(ctx.CurrentUser.Id)
	}

	sum := md5.Sum([]byte(seed))
	return hex.EncodeToString(sum[:])
}

// forumMarkUserActive records that the current user is browsing a forum.
func forumMarkUserActive(ctx *server.Context, forumId int) {
	if ctx.CurrentUser == nil {
		return
	}

	key := forumActiveUsersKey(forumId)
	now := float64(time.Now().Unix())
	context := ctx.Request.Context()

	if err := ctx.State.Redis.ZAdd(context, key, redis.Z{
		Score:  now,
		Member: strconv.Itoa(ctx.CurrentUser.Id),
	}).Err(); err != nil {
		ctx.Logger.Error("Failed to mark forum user active", "error", err, "forum", forumId)
		return
	}

	ctx.State.Redis.Expire(context, key, forumActivityExpiry)
}

// forumGetActiveUsers returns the ids of users that recently browsed the forum.
func forumGetActiveUsers(ctx *server.Context, forumId int) []int {
	key := forumActiveUsersKey(forumId)
	cutoff := time.Now().Add(-forumActivityExpiry).Unix()
	context := ctx.Request.Context()

	// First, remove any users that haven't been active in the last 5 minutes
	ctx.State.Redis.ZRemRangeByScore(context, key, "-inf", strconv.FormatInt(cutoff, 10))

	// Then, fetch the remaining users that have been active in the last 5 minutes
	members, err := ctx.State.Redis.ZRangeArgs(context, redis.ZRangeArgs{
		Key:     key,
		Start:   fmt.Sprintf("(%d", cutoff),
		Stop:    "+inf",
		ByScore: true,
		Rev:     true,
	}).Result()
	if err != nil {
		ctx.Logger.Error("Failed to fetch active forum users", "error", err, "forum", forumId)
		return nil
	}

	userIds := make([]int, 0, len(members))
	for _, member := range members {
		if id, err := strconv.Atoi(member); err == nil {
			userIds = append(userIds, id)
		}
	}
	return userIds
}

// forumAverageTopicViews returns the average view count across all
// topics, used to decide whether a topic is "hot".
func forumAverageTopicViews(ctx *server.Context) float64 {
	context := ctx.Request.Context()

	// Check if value is cached first
	if cached, err := ctx.State.Redis.Get(context, "forums:average_topic_views").Result(); err == nil {
		if value, err := strconv.ParseFloat(cached, 64); err == nil {
			return value
		}
	}

	// If not, we'll query it from the database & cache it afterwards
	average, err := ctx.State.ForumTopics.AverageViews()
	if err != nil {
		ctx.Logger.Error("Failed to fetch average topic views", "error", err)
		return 0
	}

	ctx.State.Redis.Set(
		context, "forums:average_topic_views",
		strconv.FormatFloat(average, 'f', -1, 64),
		forumAverageViewsTTL,
	)
	return average
}

// forumTopicReadStatuses resolves whether each given topic has been read by the current visitor.
func forumTopicReadStatuses(ctx *server.Context, topics []*schemas.ForumTopic) map[int]bool {
	statuses := make(map[int]bool, len(topics))
	if len(topics) == 0 {
		return statuses
	}

	// De-duplicate while preserving the topics we need to look up
	unique := make([]*schemas.ForumTopic, 0, len(topics))
	for _, topic := range topics {
		if _, seen := statuses[topic.Id]; seen {
			continue
		}
		statuses[topic.Id] = false
		unique = append(unique, topic)
	}

	key := fmt.Sprintf("forums:topic_read_timestamps:%s", forumSessionIdentifier(ctx))
	context := ctx.Request.Context()

	fields := make([]string, len(unique))
	for i, topic := range unique {
		fields[i] = strconv.Itoa(topic.Id)
	}

	timestamps, err := ctx.State.Redis.HMGet(context, key, fields...).Result()
	if err != nil {
		ctx.Logger.Error("Failed to fetch topic read timestamps", "error", err)
		return statuses
	}

	now := time.Now()
	readUpdates := make(map[string]any)

	// Check each topic's stored timestamp against its
	// last post time to determine if it's read
	for i, topic := range unique {
		if i < len(timestamps) && timestamps[i] != nil {
			raw, ok := timestamps[i].(string)
			if !ok {
				ctx.Logger.Error("Failed to parse topic read timestamp", "topic_id", topic.Id, "raw_value", timestamps[i])
				continue
			}

			timestamp, err := strconv.ParseFloat(raw, 64)
			if err != nil {
				ctx.Logger.Error("Failed to parse topic read timestamp", "topic_id", topic.Id, "raw_value", raw, "error", err)
				continue
			}

			statuses[topic.Id] = timestamp >= float64(topic.LastPostAt.Unix())
			continue
		}

		// No stored timestamp -> We treat topics older than two days as "read"
		read := now.Sub(topic.CreatedAt) >= 48*time.Hour
		statuses[topic.Id] = read
		if read {
			readUpdates[strconv.Itoa(topic.Id)] = strconv.FormatFloat(float64(now.Unix()), 'f', -1, 64)
		}
	}

	if len(readUpdates) > 0 {
		if err := ctx.State.Redis.HSet(context, key, readUpdates).Err(); err != nil {
			ctx.Logger.Error("Failed to persist topic read timestamps", "error", err)
		} else {
			ctx.State.Redis.Expire(context, key, forumReadTimestampExpiry)
		}
	}

	return statuses
}

// forumUpdateTopicReadState marks a topic as read by the current visitor.
func forumUpdateTopicReadState(ctx *server.Context, topicId int) {
	key := fmt.Sprintf("forums:topic_read_timestamps:%s", forumSessionIdentifier(ctx))
	now := strconv.FormatFloat(float64(time.Now().Unix()), 'f', -1, 64)

	context := ctx.Request.Context()
	if err := ctx.State.Redis.HSet(context, key, strconv.Itoa(topicId), now).Err(); err != nil {
		ctx.Logger.Error("Failed to update topic read state", "topic_id", topicId, "error", err)
		return
	}

	ctx.State.Redis.Expire(context, key, forumReadTimestampExpiry)
}

// forumUpdateViews increments a topic's view counter.
func forumUpdateViews(ctx *server.Context, topicId int) {
	key := fmt.Sprintf("forums:viewlock:%d:%s", topicId, ctx.IP())
	context := ctx.Request.Context()

	// View updates are limited per-ip
	// TODO: We should probably use the session identifier here though...?
	if locked, _ := ctx.State.Redis.Exists(context, key).Result(); locked > 0 {
		return
	}

	if err := ctx.State.ForumTopics.IncrementViews(topicId); err != nil {
		ctx.Logger.Error("Failed to increment topic views", "topic_id", topicId, "error", err)
		return
	}

	ctx.State.Redis.Set(context, key, 1, forumViewLockExpiry)
}

func forumActiveUsersKey(forumId int) string {
	return fmt.Sprintf("forum:%d:active", forumId)
}
