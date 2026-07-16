package repositories

import (
	"slices"
	"strings"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ForumTopicSearchSort int

const (
	ForumTopicSearchSortRelevance ForumTopicSearchSort = iota
	ForumTopicSearchSortCreated
	ForumTopicSearchSortViews
	ForumTopicSearchSortPosts
)

type ForumTopicSearchOptions struct {
	QueryString string
	Order       constants.SearchOrder
	Sort        ForumTopicSearchSort

	ForumId   *int
	CreatorId *int
	// TODO: Excluded IDs would be kinda sick

	BookmarkedByUserId *int
	SubscribedByUserId *int

	Offset int
	Limit  int
}

type ForumTopicSearchResult struct {
	Topics  []*schemas.ForumTopic
	Total   int64
	Options ForumTopicSearchOptions
}

// Normalize validates search options & applies defaults
func (options *ForumTopicSearchOptions) Normalize() {
	options.QueryString = strings.TrimSpace(options.QueryString)

	if options.Order != constants.SearchOrderAscending {
		options.Order = constants.SearchOrderDescending
	}
	if options.Sort < ForumTopicSearchSortRelevance || options.Sort > ForumTopicSearchSortPosts {
		options.Sort = ForumTopicSearchSortRelevance
	}
	if options.Limit < 1 {
		options.Limit = 25
	}
	if options.Limit > 50 {
		options.Limit = 50
	}
	if options.Offset < 0 {
		options.Offset = 0
	}

	if options.ForumId != nil && *options.ForumId < 1 {
		options.ForumId = nil
	}
	if options.CreatorId != nil && *options.CreatorId < 1 {
		options.CreatorId = nil
	}
	if options.BookmarkedByUserId != nil && *options.BookmarkedByUserId < 1 {
		options.BookmarkedByUserId = nil
	}
	if options.SubscribedByUserId != nil && *options.SubscribedByUserId < 1 {
		options.SubscribedByUserId = nil
	}
}

func (r *ForumTopicRepository) SearchPage(options ForumTopicSearchOptions, preload ...string) (*ForumTopicSearchResult, error) {
	options.Normalize()
	baseQuery := r.buildForumTopicSearchQuery(r.db.Model(&schemas.ForumTopic{}), options)

	var total int64
	countQuery := r.db.Table(
		"(?) AS filtered_forum_topics",
		baseQuery.Session(&gorm.Session{}).Select("forum_topics.id"),
	)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	options.Offset = clampSearchOffset(options.Offset, options.Limit, total)

	var topics []*schemas.ForumTopic
	resultQuery := Preloaded(baseQuery.Session(&gorm.Session{}), preload)
	resultQuery = applyForumTopicSearchSort(resultQuery, options)
	if err := resultQuery.Offset(options.Offset).Limit(options.Limit).Find(&topics).Error; err != nil {
		return nil, err
	}

	return &ForumTopicSearchResult{
		Topics:  topics,
		Total:   total,
		Options: options,
	}, nil
}

func (r *ForumTopicRepository) buildForumTopicSearchQuery(query *gorm.DB, options ForumTopicSearchOptions) *gorm.DB {
	query = query.
		Where("forum_topics.hidden = ?", false).
		Where(`EXISTS (
			SELECT 1 FROM forums
			WHERE forums.id = forum_topics.forum_id
			AND forums.hidden = ?
		)`, false)

	if options.QueryString != "" {
		query = applyForumTopicTextSearch(query, options.QueryString)
	}
	if options.ForumId != nil {
		query = query.Where("forum_topics.forum_id = ?", *options.ForumId)
	}
	if options.CreatorId != nil {
		query = query.Where("forum_topics.creator_id = ?", *options.CreatorId)
	}

	if options.BookmarkedByUserId != nil {
		query = query.Where(`EXISTS (
			SELECT 1 FROM forum_bookmarks
			WHERE forum_bookmarks.topic_id = forum_topics.id
			AND forum_bookmarks.user_id = ?
		)`, *options.BookmarkedByUserId)
	}
	if options.SubscribedByUserId != nil {
		query = query.Where(`EXISTS (
			SELECT 1 FROM forum_subscribers
			WHERE forum_subscribers.topic_id = forum_topics.id
			AND forum_subscribers.user_id = ?
		)`, *options.SubscribedByUserId)
	}

	return query
}

func applyForumTopicTextSearch(query *gorm.DB, textQuery string) *gorm.DB {
	topicCondition := "search_topics.search_vector @@ plainto_tsquery('english', ?)"
	postCondition := "search_posts.search_vector @@ plainto_tsquery('english', ?)"
	args := []any{textQuery}

	if fuzzyQuery := fuzzyTsQuery(textQuery); fuzzyQuery != "" {
		topicCondition += " OR search_topics.search_vector @@ to_tsquery('english', ?)"
		postCondition += " OR search_posts.search_vector @@ to_tsquery('english', ?)"
		args = append(args, fuzzyQuery)
	}

	// Separate subqueries let PostgreSQL use both GIN indexes before
	// combining post and topic matches into one result
	condition := `forum_topics.id IN (
		SELECT search_topics.id
		FROM forum_topics AS search_topics
		WHERE (` + topicCondition + `)
		UNION
		SELECT search_posts.topic_id
		FROM forum_posts AS search_posts
		WHERE search_posts.hidden = false
		AND search_posts.draft = false
		AND search_posts.deleted = false
		AND (` + postCondition + `)
	)`

	return query.Where(condition, slices.Concat(args, args)...)
}

var forumTopicSortExpressions = map[ForumTopicSearchSort]string{
	ForumTopicSearchSortCreated: "forum_topics.created_at",
	ForumTopicSearchSortViews:   "forum_topics.views",
	ForumTopicSearchSortPosts:   "forum_topics.post_count",
}

func applyForumTopicSearchSort(query *gorm.DB, options ForumTopicSearchOptions) *gorm.DB {
	descending := options.Order != constants.SearchOrderAscending

	// Only allow relevance if there's a search query present
	if options.Sort == ForumTopicSearchSortRelevance && options.QueryString != "" {
		expression, args := forumTopicRankExpression(options.QueryString)
		return applySearchRankOrder(query, expression, args, descending, "forum_topics.id")
	}

	// Use "created" by default and allow specified sort expressions to override it
	expression := forumTopicSortExpressions[ForumTopicSearchSortCreated]
	if configured, ok := forumTopicSortExpressions[options.Sort]; ok {
		expression = configured
	}
	return query.
		Order(clause.OrderByColumn{
			Column: clause.Column{Name: expression, Raw: true},
			Desc:   descending,
		}).
		Order("forum_topics.id DESC")
}

func forumTopicRankExpression(textQuery string) (string, []any) {
	// This allows a partial term such as "serv" to match words beginning with
	// that prefix, such as "server" or "servers"
	fuzzyQuery := fuzzyTsQuery(textQuery)

	// fuzzyTsQuery returns an empty string when the input contains no usable words
	// In that case, avoid passing an empty query to to_tsquery()
	if fuzzyQuery == "" {
		return `(
			-- Rank matches found directly in the topic's search vector
			2 * ts_rank(forum_topics.search_vector, plainto_tsquery('english', ?))
			-- Add the score of the single best matching post in the topic
			+ COALESCE((
				SELECT MAX(ts_rank(search_rank_posts.search_vector, plainto_tsquery('english', ?)))
				FROM forum_posts AS search_rank_posts
				-- Restrict the subquery to posts belonging to the current topic
				WHERE search_rank_posts.topic_id = forum_topics.id
				-- Exclude posts that should not be visible inside the search results
				AND search_rank_posts.hidden = false
				AND search_rank_posts.draft = false
				AND search_rank_posts.deleted = false
			), 0)
		)`, []any{textQuery, textQuery}
	}

	// When a prefix query is available, calculate both:
	// - A normal full-text search rank using plainto_tsquery()
	// - A prefix-search rank using to_tsquery()
	return `(
		-- Calculate the topic's own relevance score (word-based & prefix-based)
		2 * GREATEST(
			ts_rank(forum_topics.search_vector, plainto_tsquery('english', ?)),
			ts_rank(forum_topics.search_vector, to_tsquery('english', ?)) * 0.5
		)
		-- Add the relevance score of the topic's best matching eligible post
		+ COALESCE((
			SELECT MAX(GREATEST(
				ts_rank(search_rank_posts.search_vector, plainto_tsquery('english', ?)),
				ts_rank(search_rank_posts.search_vector, to_tsquery('english', ?)) * 0.5
			))
			FROM forum_posts AS search_rank_posts
			-- Restrict the subquery to posts belonging to the current topic
			WHERE search_rank_posts.topic_id = forum_topics.id
			-- Exclude posts that should not be visible inside the search results
			AND search_rank_posts.hidden = false
			AND search_rank_posts.draft = false
			AND search_rank_posts.deleted = false
		), 0)
	)`, []any{textQuery, fuzzyQuery, textQuery, fuzzyQuery}
}
