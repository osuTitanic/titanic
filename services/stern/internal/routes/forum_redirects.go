package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
)

// ForumPostRedirect handles /forum/p/{postId} -> /forum/{forumId}/t/{topicId}/p/{postId} redirects
func ForumPostRedirect(ctx *server.Context) {
	postId, err := ctx.PathValueInt64("postId")
	if err != nil {
		NotFound(ctx)
		return
	}
	redirectToForumPost(ctx, postId)
}

// ForumPostPermalink redirects to a #post-{postId} anchor
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

// ForumShortlinkRedirect handles the short topic & post urls, i.e.
// /forum/t/{topicId}, /forum/t/{topicId}/p/{postId} and /forum/p/{postId}.
// These cannot be registered as separate patterns, since they would
// conflict with the /forum/{id}/... routes inside of net/http.
func ForumShortlinkRedirect(ctx *server.Context) {
	kind := ctx.PathValue("kind")
	rest := strings.Trim(ctx.PathValue("rest"), "/")
	segments := strings.Split(rest, "/")

	switch {
	// /forum/t/{topicId}
	case kind == "t" && len(segments) == 1:
		topicId, err := strconv.Atoi(segments[0])
		if err != nil {
			NotFound(ctx)
			return
		}
		redirectToForumTopic(ctx, topicId)
	// /forum/t/{topicId}/p/{postId}
	// /forum/p/{postId}
	case kind == "t" && len(segments) == 3 && segments[1] == "p",
		kind == "p" && len(segments) == 1:
		postId, err := strconv.ParseInt(segments[len(segments)-1], 10, 64)
		if err != nil {
			NotFound(ctx)
			return
		}
		redirectToForumPost(ctx, postId)
	default:
		NotFound(ctx)
	}
}

// ForumQuickReplyRedirect handles phpBB-style /forum/posting.php?t={topicId} urls
func ForumQuickReplyRedirect(ctx *server.Context) {
	topicId, err := ctx.QueryValueInt("t")
	if err != nil || topicId <= 0 {
		RenderErrorPage(ctx, http.StatusNotFound, "Topic Not Found", "The topic you are looking for could not be found.")
		return
	}

	topic, err := ctx.State.ForumTopics.ById(topicId)
	if err != nil || topic == nil {
		RenderErrorPage(ctx, http.StatusNotFound, "Topic Not Found", "The topic you are looking for could not be found.")
		return
	}

	ctx.Redirect(http.StatusFound, fmt.Sprintf(
		"/forum/%d/t/%d/post", topic.ForumId, topic.Id,
	))
}

// ForumViewTopicRedirect handles phpBB-style /forum/viewtopic.php?t={topicId}&p={postId} urls
func ForumViewTopicRedirect(ctx *server.Context) {
	if postId, err := ctx.QueryValueInt64("p"); err == nil && postId > 0 {
		redirectToForumPost(ctx, postId)
		return
	}

	if topicId, err := ctx.QueryValueInt("t"); err == nil && topicId > 0 {
		redirectToForumTopic(ctx, topicId)
		return
	}

	RenderErrorPage(ctx, http.StatusNotFound, "Topic Not Found", "The topic you are looking for could not be found.")
}

// ForumViewForumRedirect handles phpBB-style /forum/viewforum.php?f={forumId} urls
func ForumViewForumRedirect(ctx *server.Context) {
	forumId, err := ctx.QueryValueInt("f")
	if err != nil || forumId <= 0 {
		RenderErrorPage(ctx, http.StatusNotFound, "Forum Not Found", "The forum you are looking for could not be found.")
		return
	}

	ctx.Redirect(http.StatusFound, fmt.Sprintf("/forum/%d", forumId))
}

// ForumControlPanelRedirect handles phpBB-style /forum/ucp.php?mode={mode} urls
func ForumControlPanelRedirect(ctx *server.Context) {
	redirects := map[string]string{
		"register":     "/account/register",
		"sendpassword": "/account/reset",
		"avatar":       "/account/profile#avatar",
	}
	if target, ok := redirects[ctx.QueryValue("mode")]; ok {
		ctx.Redirect(http.StatusFound, target)
		return
	}
	ctx.Redirect(http.StatusFound, "/account")
}

// ForumIndexRedirect handles phpBB-style /forum/index.php urls
func ForumIndexRedirect(ctx *server.Context) {
	ctx.Redirect(http.StatusFound, "/forum")
}

func redirectToForumTopic(ctx *server.Context, topicId int) {
	topic, err := ctx.State.ForumTopics.ById(topicId)
	if err != nil || topic == nil {
		RenderErrorPage(ctx, http.StatusNotFound, "Topic Not Found", "The topic you are looking for could not be found.")
		return
	}

	ctx.Redirect(http.StatusFound, fmt.Sprintf(
		"/forum/%d/t/%d", topic.ForumId, topic.Id,
	))
}

func redirectToForumPost(ctx *server.Context, postId int64) {
	post, err := ctx.State.ForumPosts.ById(postId)
	if err != nil || post == nil {
		RenderErrorPage(ctx, http.StatusNotFound, "Post Not Found", "The post you are looking for could not be found.")
		return
	}

	ctx.Redirect(http.StatusFound, fmt.Sprintf(
		"/forum/%d/t/%d/p/%d", post.ForumId, post.TopicId, post.Id,
	))
}
