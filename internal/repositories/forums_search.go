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

	ForumId *int
	Creator string

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
	options.Creator = strings.TrimSpace(options.Creator)

	if options.Order != constants.SearchOrderAscending {
		options.Order = constants.SearchOrderDescending
	}
	if options.Sort < ForumTopicSearchSortRelevance || options.Sort > ForumTopicSearchSortPosts {
		options.Sort = ForumTopicSearchSortRelevance
	}
	if options.Sort == ForumTopicSearchSortRelevance && options.QueryString == "" {
		options.Sort = ForumTopicSearchSortCreated
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

func (r *ForumPostRepository) FetchSearchPreviews(topicIds []int, textQuery, creator string, preload ...string) (map[int]*schemas.ForumPost, error) {
	postsByTopic := make(map[int]*schemas.ForumPost, len(topicIds))
	textQuery = strings.TrimSpace(textQuery)
	if len(topicIds) == 0 {
		return postsByTopic, nil
	}

	query := Preloaded(r.db.Model(&schemas.ForumPost{}), preload).
		Select("DISTINCT ON (forum_posts.topic_id) forum_posts.*").
		Where("forum_posts.topic_id IN ?", topicIds).
		Where("forum_posts.hidden = ?", false).
		Where("forum_posts.draft = ?", false).
		Where("forum_posts.deleted = ?", false)
	query = applyForumPostPreviewCreatorFilter(query, creator, textQuery)

	orderExpression := "forum_posts.topic_id ASC, forum_posts.id DESC"
	var orderArgs []any
	if textQuery != "" {
		rankExpression, rankArgs := forumSearchVectorRankExpression(
			"forum_posts.search_vector",
			textQuery,
		)
		orderExpression = "forum_posts.topic_id ASC, " + rankExpression + " DESC, forum_posts.id ASC"
		orderArgs = rankArgs
	}

	var posts []*schemas.ForumPost
	err := query.Order(clause.OrderBy{Expression: clause.Expr{
		SQL:  orderExpression,
		Vars: orderArgs,
	}}).Find(&posts).Error
	if err != nil {
		return nil, err
	}
	for _, post := range posts {
		postsByTopic[post.TopicId] = post
	}
	return postsByTopic, nil
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
	if options.Creator != "" {
		query = applyForumTopicCreatorFilter(query, options.Creator, options.QueryString)
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

func applyForumTopicCreatorFilter(query *gorm.DB, creator, textQuery string) *gorm.DB {
	postMatchCondition := ""
	args := []any{schemas.ResolveSafeName(creator)}
	if textQuery != "" {
		condition, conditionArgs := forumSearchVectorMatchCondition("search_user_posts.search_vector", textQuery)
		postMatchCondition = " AND (" + condition + ")"
		args = append(args, conditionArgs...)
	}

	condition := `EXISTS (
		SELECT 1 FROM users AS search_users
		WHERE search_users.safe_name = ?
		AND (
			search_users.id = forum_topics.creator_id
			OR EXISTS (
				SELECT 1 FROM forum_posts AS search_user_posts
				WHERE search_user_posts.topic_id = forum_topics.id
				AND search_user_posts.user_id = search_users.id
				AND search_user_posts.hidden = false
				AND search_user_posts.draft = false
				AND search_user_posts.deleted = false` + postMatchCondition + `
			)
		)
	)`
	return query.Where(condition, args...)
}

func applyForumPostPreviewCreatorFilter(query *gorm.DB, creator, textQuery string) *gorm.DB {
	creator = strings.TrimSpace(creator)
	if creator == "" {
		return query
	}

	postMatchCondition := ""
	args := []any{schemas.ResolveSafeName(creator)}
	if textQuery != "" {
		condition, conditionArgs := forumSearchVectorMatchCondition("forum_posts.search_vector", textQuery)
		postMatchCondition = " AND (" + condition + ")"
		args = append(args, conditionArgs...)
	}

	condition := `EXISTS (
		SELECT 1
		FROM users AS search_preview_users
		JOIN forum_topics AS search_preview_topics
		ON search_preview_topics.id = forum_posts.topic_id
		WHERE search_preview_users.safe_name = ?
		AND (
			search_preview_users.id = search_preview_topics.creator_id
			OR (
				search_preview_users.id = forum_posts.user_id` + postMatchCondition + `
			)
		)
	)`
	return query.Where(condition, args...)
}

func applyForumTopicTextSearch(query *gorm.DB, textQuery string) *gorm.DB {
	topicCondition, topicArgs := forumSearchVectorMatchCondition("search_topics.search_vector", textQuery)
	postCondition, postArgs := forumSearchVectorMatchCondition("search_posts.search_vector", textQuery)

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

	return query.Where(condition, slices.Concat(topicArgs, postArgs)...)
}

func forumSearchVectorMatchCondition(searchVector, textQuery string) (string, []any) {
	condition := searchVector + " @@ plainto_tsquery('english', ?)"
	args := []any{textQuery}
	if fuzzyQuery := fuzzyTsQuery(textQuery); fuzzyQuery != "" {
		condition += " OR " + searchVector + " @@ to_tsquery('english', ?)"
		args = append(args, fuzzyQuery)
	}
	return condition, args
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
	topicRank, topicArgs := forumSearchVectorRankExpression("forum_topics.search_vector", textQuery)
	postRank, postArgs := forumSearchVectorRankExpression("search_rank_posts.search_vector", textQuery)

	return `(
		2 * (` + topicRank + `)
		+ COALESCE((
			SELECT MAX(` + postRank + `)
			FROM forum_posts AS search_rank_posts
			WHERE search_rank_posts.topic_id = forum_topics.id
			AND search_rank_posts.hidden = false
			AND search_rank_posts.draft = false
			AND search_rank_posts.deleted = false
		), 0)
	)`, slices.Concat(topicArgs, postArgs)
}

func forumSearchVectorRankExpression(searchVector, textQuery string) (string, []any) {
	fuzzyQuery := fuzzyTsQuery(textQuery)
	if fuzzyQuery == "" {
		return "ts_rank(" + searchVector + ", plainto_tsquery('english', ?))", []any{textQuery}
	}

	return `GREATEST(
		ts_rank(` + searchVector + `, plainto_tsquery('english', ?)),
		ts_rank(` + searchVector + `, to_tsquery('english', ?)) * 0.5
	)`, []any{textQuery, fuzzyQuery}
}
