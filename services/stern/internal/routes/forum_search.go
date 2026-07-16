package routes

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"unicode"

	"github.com/osuTitanic/titanic/internal/bbcode"
	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/repositories"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/services/stern/internal/server"
	"github.com/osuTitanic/titanic/services/stern/internal/templates"
)

const forumSearchPostsPerPage = 25
const forumSearchExcerptLength = 200

func ForumSearch(ctx *server.Context) {
	query := ctx.Request.URL.Query()

	var currentUserId *int
	if ctx.CurrentUser != nil {
		currentUserId = &ctx.CurrentUser.Id
	}

	options := buildForumPostSearchOptions(query, currentUserId)
	result, err := ctx.State.ForumPosts.SearchPage(
		options,
		"Topic", "Topic.Forum", "User", "User.Groups.Group",
	)
	if err != nil {
		ctx.Logger.Error("Failed to search forum posts", "options", options, "error", err)
		InternalServerError(ctx)
		return
	}

	currentUserIdValue := 0
	if currentUserId != nil {
		currentUserIdValue = *currentUserId
	}
	previews := buildForumPostSearchPreviews(
		result.Posts,
		result.Options.QueryString,
		currentUserIdValue,
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
		Posts:       previews,
		SearchSort:  strconv.Itoa(int(result.Options.Sort)),
		SearchOrder: strconv.Itoa(int(result.Options.Order)),
		DefaultSort: strconv.Itoa(int(defaultSearchSort(result.Options.QueryString))),
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

func buildForumPostSearchOptions(query url.Values, currentUserId *int) repositories.ForumPostSearchOptions {
	options := repositories.ForumPostSearchOptions{
		QueryString: query.Get("query"),
		Order:       constants.SearchOrderDescending,
		Sort:        defaultSearchSort(query.Get("query")),
		Limit:       forumSearchPostsPerPage,
	}

	if forumId, ok := parseInt(query.Get("forum")); ok && forumId > 0 {
		options.ForumId = &forumId
	}
	if username := strings.TrimSpace(query.Get("username")); username != "" {
		options.Author = username
	}
	if sortValue, ok := parseInt(query.Get("sort")); ok {
		options.Sort = repositories.ForumSearchSort(sortValue)
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

func normalizeForumSearchQuery(query url.Values, options repositories.ForumPostSearchOptions, page int) {
	setSearchQueryValue(query, "query", options.QueryString)
	setSearchQueryValue(query, "username", options.Author)

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

func defaultSearchSort(query string) repositories.ForumSearchSort {
	if strings.TrimSpace(query) == "" {
		return repositories.ForumSearchSortCreated
	}
	return repositories.ForumSearchSortRelevance
}

func buildForumPostSearchPreviews(
	posts []*schemas.ForumPost,
	textQuery string,
	currentUserId int,
) []*templates.ForumSearchPostPreview {
	previews := make([]*templates.ForumSearchPostPreview, 0, len(posts))
	for _, post := range posts {
		if post.Topic == nil {
			continue
		}
		previews = append(previews, &templates.ForumSearchPostPreview{
			Post:          post,
			Excerpt:       forumPostSearchExcerpt(post.Content, textQuery),
			Index:         len(previews),
			CurrentUserId: currentUserId,
		})
	}
	return previews
}

func forumPostSearchExcerpt(content, textQuery string) []templates.ForumSearchExcerptPart {
	// We want to strip out bbcode & collapse whitespace
	text := strings.Join(strings.Fields(bbcode.Strip(content, false)), " ")
	if text == "" {
		return nil
	}

	terms := forumSearchTerms(textQuery)
	runes := []rune(text)

	if len(runes) <= forumSearchExcerptLength {
		// The text is short enough that we don't need to create an excerpt
		return forumPostSearchExcerptParts(text, terms)
	}

	// Find the earliest occurrence of any search term in the text, and create an excerpt around it
	start := forumPostSearchExcerptStart(runes, terms)
	end := min(start+forumSearchExcerptLength, len(runes))
	excerpt := string(runes[start:end])

	if start > 0 {
		excerpt = "..." + excerpt
	}
	if end < len(runes) {
		excerpt += "..."
	}
	return forumPostSearchExcerptParts(excerpt, terms)
}

func forumSearchTerms(textQuery string) map[string]struct{} {
	// Split on every non-letter & remove duplicate terms
	terms := make(map[string]struct{})
	for _, term := range strings.FieldsFunc(strings.ToLower(textQuery), func(r rune) bool {
		return !isForumSearchWordRune(r)
	}) {
		terms[term] = struct{}{}
	}
	return terms
}

func forumPostSearchExcerptStart(text []rune, terms map[string]struct{}) int {
	if len(terms) == 0 {
		return 0
	}

	for start := 0; start < len(text); {
		if !isForumSearchWordRune(text[start]) {
			// Skip non-word characters
			start++
			continue
		}
		end := start + 1
		for end < len(text) && isForumSearchWordRune(text[end]) {
			// Find the end of the current word
			end++
		}

		// Check if the current word matches any of the search terms
		if _, matched := terms[strings.ToLower(string(text[start:end]))]; matched {
			// We found a match, so we want to create an excerpt that starts a bit before the match
			excerptStart := max(0, start-forumSearchExcerptLength/3)
			return min(excerptStart, len(text)-forumSearchExcerptLength)
		}
		start = end
	}
	return 0
}

func forumPostSearchExcerptParts(text string, terms map[string]struct{}) []templates.ForumSearchExcerptPart {
	runes := []rune(text)
	parts := make([]templates.ForumSearchExcerptPart, 0)

	// We want to split the text into alternating sequences of "word" and "non-word"
	// characters, so we can highlight only the whole-word matches

	for start := 0; start < len(runes); {
		isWord := isForumSearchWordRune(runes[start])
		end := start + 1
		for end < len(runes) && isForumSearchWordRune(runes[end]) == isWord {
			// Keep going until we reach the end of the current word or non-word sequence
			end++
		}

		part := templates.ForumSearchExcerptPart{Text: string(runes[start:end])}
		if isWord {
			// Only exact, whole-word matches should be highlighted
			_, part.Matched = terms[strings.ToLower(part.Text)]
		}
		parts = append(parts, part)
		start = end
	}
	return parts
}

func isForumSearchWordRune(r rune) bool {
	// We only want to match letters and numbers as part of search terms
	return unicode.IsLetter(r) || unicode.IsNumber(r)
}
