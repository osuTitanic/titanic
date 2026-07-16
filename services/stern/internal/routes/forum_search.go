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
	options, page := buildForumTopicSearchOptions(query, currentUserId)

	result, err := ctx.State.ForumTopics.SearchPage(
		options,
		"Forum", "Icon", "Creator", "Creator.Groups.Group",
	)
	if err != nil {
		ctx.Logger.Error("Failed to search forum topics", "options", options, "error", err)
		InternalServerError(ctx)
		return
	}

	topicIds := make([]int, 0, len(result.Topics))
	for _, topic := range result.Topics {
		topicIds = append(topicIds, topic.Id)
	}
	lastPosts, err := ctx.State.ForumPosts.FetchLastForTopics(topicIds, "User", "User.Groups.Group")
	if err != nil {
		ctx.Logger.Error("Failed to fetch last posts for forum search", "error", err)
		lastPosts = map[int]*schemas.ForumPost{}
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
		lastPosts,
		forumTopicReadStatuses(ctx, result.Topics),
		averageViews,
		hasCustomIcons,
		currentUserIdValue,
		true,
	)
	page = result.Options.Offset/result.Options.Limit + 1
	view := templates.ForumSearchView{
		DefaultView:    buildDefaultView(ctx),
		Topics:         previews,
		HasCustomIcons: hasCustomIcons,
		SearchSort:     strconv.Itoa(int(result.Options.Sort)),
		SearchOrder:    strconv.Itoa(int(result.Options.Order)),
		DefaultSort:    strconv.Itoa(int(defaultForumTopicSearchSort(query.Get("query")))),
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

func buildForumTopicSearchOptions(query url.Values, currentUserId *int) (repositories.ForumTopicSearchOptions, int) {
	options := repositories.ForumTopicSearchOptions{
		QueryString: query.Get("query"),
		Order:       constants.SearchOrderDescending,
		Sort:        defaultForumTopicSearchSort(query.Get("query")),
		Limit:       forumSearchTopicsPerPage,
	}

	if forumId, ok := parseInt(query.Get("forum")); ok && forumId > 0 {
		options.ForumId = &forumId
	}
	if userId, ok := parseInt(query.Get("user")); ok && userId > 0 {
		options.CreatorId = &userId
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

	return options, page
}

func defaultForumTopicSearchSort(query string) repositories.ForumTopicSearchSort {
	if strings.TrimSpace(query) == "" {
		return repositories.ForumTopicSearchSortCreated
	}
	return repositories.ForumTopicSearchSortRelevance
}
