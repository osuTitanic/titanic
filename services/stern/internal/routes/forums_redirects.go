package routes

import (
	"fmt"
	"net/http"

	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
)

func ForumTopicRedirect(ctx *server.Context) {
	topicId, err := ctx.PathValueInt("topicId")
	if err != nil {
		NotFound(ctx)
		return
	}

	topic, err := ctx.State.ForumTopics.ById(topicId)
	if err != nil || topic == nil {
		RenderErrorPage(ctx, http.StatusNotFound, "Topic Not Found", "The topic you are looking for could not be found.")
		return
	}

	ctx.Redirect(http.StatusFound, fmt.Sprintf(
		"/forum/%d/t/%d", topic.ForumId, topic.Id,
	))
}

func ForumPostRedirect(ctx *server.Context) {
	postId, err := ctx.PathValueInt64("postId")
	if err != nil {
		NotFound(ctx)
		return
	}

	post, err := ctx.State.ForumPosts.ById(postId)
	if err != nil || post == nil {
		RenderErrorPage(ctx, http.StatusNotFound, "Post Not Found", "The post you are looking for could not be found.")
		return
	}

	ctx.Redirect(http.StatusFound, fmt.Sprintf(
		"/forum/%d/t/%d/p/%d", post.ForumId, post.TopicId, post.Id,
	))
}

func ForumPostPermalink(ctx *server.Context) {
	postId, err := ctx.PathValueInt64("postId")
	if err != nil {
		NotFound(ctx)
		return
	}

	post, err := ctx.State.ForumPosts.ById(postId)
	if err != nil || post == nil {
		RenderErrorPage(ctx, http.StatusNotFound, "Post Not Found", "The post you are looking for could not be found.")
		return
	}

	countBefore, err := ctx.State.ForumPosts.CountBeforePost(post.Id, post.TopicId)
	if err != nil {
		ctx.Logger.Error("Failed to count posts before permalink", "error", err, "post", post.Id)
		InternalServerError(ctx)
		return
	}
	page := (countBefore / forumPostsPerPage) + 1

	ctx.Redirect(http.StatusFound, fmt.Sprintf(
		"/forum/%d/t/%d/?page=%d#post-%d",
		post.ForumId, post.TopicId, page, post.Id,
	))
}
