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
		DefaultView:        buildDefaultView(ctx),
		News:               fetchHomeNews(ctx),
		ChatMessages:       fetchHomeChatMessages(ctx),
		MostPlayedBeatmaps: fetchHomeMostPlayedBeatmaps(ctx),
	}
	ctx.RenderTemplate(http.StatusOK, "pages/public/home", view)
}

func HomeNewsPartial(ctx *server.Context) {
	view := any(fetchHomeNews(ctx))
	ctx.RenderTemplate(http.StatusOK, "partials/home_news", view)
}

func HomeChatPartial(ctx *server.Context) {
	view := any(fetchHomeChatMessages(ctx))
	ctx.RenderTemplate(http.StatusOK, "partials/home_chat", view)
}

func HomePlaysPartial(ctx *server.Context) {
	view := any(fetchHomeMostPlayedBeatmaps(ctx))
	ctx.RenderTemplate(http.StatusOK, "partials/home_plays", view)
}

func fetchHomeNews(ctx *server.Context) []*schemas.ForumPost {
	topics, err := ctx.State.ForumTopics.FetchAnnouncements(newsLimit, 0, "Creator")
	if err != nil {
		ctx.Logger.Error("Failed to fetch home news topics", "error", err)
		return []*schemas.ForumPost{}
	}

	topicIds := make([]int, 0, len(topics))
	for _, topic := range topics {
		topicIds = append(topicIds, topic.Id)
	}

	postsByTopic, err := ctx.State.ForumPosts.FetchInitialByTopicIds(topicIds, "Topic", "User")
	if err != nil {
		ctx.Logger.Error("Failed to fetch home news posts", "error", err)
		return []*schemas.ForumPost{}
	}

	posts := make([]*schemas.ForumPost, 0, len(topics))
	for _, topic := range topics {
		post, ok := postsByTopic[topic.Id]
		if ok {
			posts = append(posts, post)
			continue
		}
	}
	return posts
}

func fetchHomeChatMessages(ctx *server.Context) []*schemas.Message {
	messages, err := ctx.State.Messages.FetchRecent(chatMessagesChannel, chatMessagesLimit, 0)
	if err != nil {
		ctx.Logger.Error("Failed to fetch home chat messages", "error", err)
		return []*schemas.Message{}
	}

	// Reverse so the latest message is displayed at the bottom
	// Idk if there's a better way of doing this
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
	return messages
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
