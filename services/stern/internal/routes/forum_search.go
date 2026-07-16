package routes

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/repositories"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/services/stern/internal/server"
	"github.com/osuTitanic/titanic/services/stern/internal/templates"
)

const forumSearchTopicsPerPage = 25

func ForumSearch(ctx *server.Context) {
	query := ctx.Request.URL.Query()

	var currentUserId *int
	if ctx.CurrentUser != nil {
		currentUserId = &ctx.CurrentUser.Id
	}

	options := buildForumTopicSearchOptions(query, currentUserId)
	result, err := ctx.State.ForumTopics.SearchPage(
		options,
		"Forum", "Icon", "Creator", "Creator.Groups.Group",
	)
	if err != nil {
		ctx.Logger.Error("Failed to search forum topics", "options", options, "error", err)
		InternalServerError(ctx)
		return
	}

	currentUserIdValue := 0
	if currentUserId != nil {
		currentUserIdValue = *currentUserId
	}

	hasCustomIcons := topicsHaveCustomIcons(result.Topics)
	averageViews := 0.0
	if len(result.Topics) > 0 {
		averageViews = forumAverageTopicViews(ctx)
	}

	previews := buildTopicPreviews(
		result.Topics,
		fetchForumSearchPreviewPosts(ctx, result.Topics, result.Options.QueryString),
		forumTopicReadStatuses(ctx, result.Topics),
		averageViews,
		hasCustomIcons,
		currentUserIdValue,
		true,
	)
	page := result.Options.Offset/result.Options.Limit + 1
	normalizeForumSearchQuery(query, result.Options, page)

	selectedForumId := 0
	if result.Options.ForumId != nil {
		selectedForumId = *result.Options.ForumId
	}
	defaultView := buildDefaultView(ctx)
	defaultView.Query = query

	view := templates.ForumSearchView{
		DefaultView: defaultView,
		ForumJump:   buildForumJumpView(ctx, selectedForumId),
		Topics:      previews,
		SearchSort:  strconv.Itoa(int(result.Options.Sort)),
		SearchOrder: strconv.Itoa(int(result.Options.Order)),
		DefaultSort: strconv.Itoa(int(defaultForumTopicSearchSort(result.Options.QueryString))),
		Pagination: templates.NewPagination(templates.PaginationOptions{
			Path:        "/forum/search",
			Query:       query,
			CurrentPage: page,
			PageSize:    result.Options.Limit,
			Total:       int(result.Total),
		}),
	}
	ctx.RenderTemplate(http.StatusOK, "pages/forum/search", view)
}

func buildForumTopicSearchOptions(query url.Values, currentUserId *int) repositories.ForumTopicSearchOptions {
	options := repositories.ForumTopicSearchOptions{
		QueryString: query.Get("query"),
		Order:       constants.SearchOrderDescending,
		Sort:        defaultForumTopicSearchSort(query.Get("query")),
		Limit:       forumSearchTopicsPerPage,
	}

	if forumId, ok := parseInt(query.Get("forum")); ok && forumId > 0 {
		options.ForumId = &forumId
	}
	if username := strings.TrimSpace(query.Get("username")); username != "" {
		options.Creator = username
	}
	if sortValue, ok := parseInt(query.Get("sort")); ok {
		options.Sort = repositories.ForumTopicSearchSort(sortValue)
	}
	if order, ok := parseInt(query.Get("order")); ok {
		options.Order = constants.SearchOrder(order)
	}

	if currentUserId != nil {
		if query.Get("bookmarked") != "" {
			options.BookmarkedByUserId = currentUserId
		}
		if query.Get("subscribed") != "" {
			options.SubscribedByUserId = currentUserId
		}
	}

	page := 1
	if parsed, ok := parseInt(query.Get("page")); ok && parsed > 1 {
		page = parsed
	}
	options.Offset = (page - 1) * options.Limit

	return options
}

func fetchForumSearchPreviewPosts(
	ctx *server.Context,
	topics []*schemas.ForumTopic,
	textQuery string,
) map[int]*schemas.ForumPost {
	topicIds := make([]int, 0, len(topics))
	for _, topic := range topics {
		topicIds = append(topicIds, topic.Id)
	}

	var (
		posts map[int]*schemas.ForumPost
		err   error
	)
	if textQuery == "" {
		posts, err = ctx.State.ForumPosts.FetchLastForTopics(
			topicIds,
			"User", "User.Groups.Group",
		)
	} else {
		posts, err = ctx.State.ForumPosts.FetchSearchMatches(
			topicIds,
			textQuery,
			"User", "User.Groups.Group",
		)
	}
	if err != nil {
		ctx.Logger.Error("Failed to fetch preview posts for forum search", "error", err)
		return map[int]*schemas.ForumPost{}
	}
	return posts
}

func normalizeForumSearchQuery(query url.Values, options repositories.ForumTopicSearchOptions, page int) {
	setSearchQueryValue(query, "query", options.QueryString)
	setSearchQueryValue(query, "username", options.Creator)

	if options.ForumId == nil {
		query.Del("forum")
	} else {
		query.Set("forum", strconv.Itoa(*options.ForumId))
	}

	if options.BookmarkedByUserId == nil {
		query.Del("bookmarked")
	} else {
		query.Set("bookmarked", "1")
	}
	if options.SubscribedByUserId == nil {
		query.Del("subscribed")
	} else {
		query.Set("subscribed", "1")
	}

	if _, present := query["sort"]; present {
		query.Set("sort", strconv.Itoa(int(options.Sort)))
	}
	if _, present := query["order"]; present {
		query.Set("order", strconv.Itoa(int(options.Order)))
	}
	if page <= 1 {
		query.Del("page")
	} else {
		query.Set("page", strconv.Itoa(page))
	}
}

func setSearchQueryValue(query url.Values, key, value string) {
	if value == "" {
		query.Del(key)
		return
	}
	query.Set(key, value)
}

func defaultForumTopicSearchSort(query string) repositories.ForumTopicSearchSort {
	if strings.TrimSpace(query) == "" {
		return repositories.ForumTopicSearchSortCreated
	}
	return repositories.ForumTopicSearchSortRelevance
}
