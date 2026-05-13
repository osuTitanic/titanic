package routes

import (
	"net/http"
	"time"

	"github.com/osuTitanic/titanic-go/internal/schemas"
	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
)

const mostPlayedDelta = 7 * 24 * time.Hour
const mostPlayedLimit = 5
const chatMessagesChannel = "#osu"
const chatMessagesLimit = 10
const newsLimit = 4

func Home(ctx *server.Context) {
	view := templates.HomeView{
		// TODO: Refactor home view such that we don't have to dereference pointers
		DefaultView:        buildDefaultView(ctx),
		News:               fetchHomeNews(ctx),
		ChatMessages:       fetchHomeChatMessages(ctx),
		MostPlayedBeatmaps: fetchHomeMostPlayedBeatmaps(ctx),
	}
	ctx.RenderTemplate(http.StatusOK, "pages/public/home", view)
}

func fetchHomeNews(ctx *server.Context) []schemas.ForumPost {
	topics, err := ctx.State.ForumTopics.FetchAnnouncements(newsLimit, 0, "Creator")
	if err != nil {
		ctx.Logger.Error("Failed to fetch home news topics", "error", err)
		return []schemas.ForumPost{}
	}

	topicIds := make([]int, 0, len(topics))
	for _, topic := range topics {
		topicIds = append(topicIds, topic.Id)
	}

	postsByTopic, err := ctx.State.ForumPosts.FetchInitialByTopicIds(topicIds, "Topic", "User")
	if err != nil {
		ctx.Logger.Error("Failed to fetch home news posts", "error", err)
		return []schemas.ForumPost{}
	}

	posts := make([]schemas.ForumPost, 0, len(topics))
	for _, topic := range topics {
		post, ok := postsByTopic[topic.Id]
		if ok {
			posts = append(posts, *post)
			continue
		}

		posts = append(posts, schemas.ForumPost{
			TopicId:   topic.Id,
			ForumId:   topic.ForumId,
			UserId:    topic.CreatorId,
			CreatedAt: topic.CreatedAt,
			Topic:     topic,
			User:      topic.Creator,
		})
	}
	return posts
}

func fetchHomeChatMessages(ctx *server.Context) []schemas.Message {
	messages, err := ctx.State.Messages.FetchRecent(chatMessagesChannel, chatMessagesLimit, 0)
	if err != nil {
		ctx.Logger.Error("Failed to fetch home chat messages", "error", err)
		return []schemas.Message{}
	}

	values := make([]schemas.Message, 0, len(messages))
	for i := len(messages) - 1; i >= 0; i-- {
		message := messages[i]
		values = append(values, *message)
	}
	return values
}

func fetchHomeMostPlayedBeatmaps(ctx *server.Context) map[int]*schemas.Beatmap {
	beatmaps, err := ctx.State.Beatmaps.FetchMostPlayedSince(
		time.Now().Add(-mostPlayedDelta),
		mostPlayedLimit,
		"Beatmapset",
		"Beatmapset.CreatorUser",
	)
	if err != nil {
		ctx.Logger.Error("Failed to fetch home most played beatmaps", "error", err)
		return map[int]*schemas.Beatmap{}
	}
	return beatmaps
}
