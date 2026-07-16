package repositories

import (
	"strings"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ForumSearchSort int

const (
	ForumSearchSortRelevance ForumSearchSort = iota
	ForumSearchSortCreated
)

type ForumPostSearchOptions struct {
	QueryString string
	Order       constants.SearchOrder
	Sort        ForumSearchSort

	ForumId *int
	Author  string

	BookmarkedByUserId *int
	SubscribedByUserId *int

	Offset int
	Limit  int
}

type ForumPostSearchResult struct {
	Posts   []*schemas.ForumPost
	Total   int64
	Options ForumPostSearchOptions
}

// Normalize validates search options & applies defaults
func (options *ForumPostSearchOptions) Normalize() {
	options.QueryString = strings.TrimSpace(options.QueryString)
	options.Author = strings.TrimSpace(options.Author)

	if options.Order != constants.SearchOrderAscending {
		options.Order = constants.SearchOrderDescending
	}
	if options.Sort < ForumSearchSortRelevance || options.Sort > ForumSearchSortCreated {
		options.Sort = ForumSearchSortRelevance
	}
	if options.Sort == ForumSearchSortRelevance && options.QueryString == "" {
		options.Sort = ForumSearchSortCreated
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

func (r *ForumPostRepository) SearchPage(options ForumPostSearchOptions, preload ...string) (*ForumPostSearchResult, error) {
	options.Normalize()
	baseQuery := r.buildForumPostSearchQuery(r.db.Model(&schemas.ForumPost{}), options)

	var total int64
	countQuery := r.db.Table(
		"(?) AS filtered_forum_posts",
		baseQuery.Session(&gorm.Session{}).Select("forum_posts.id"),
	)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	options.Offset = clampSearchOffset(options.Offset, options.Limit, total)

	var posts []*schemas.ForumPost
	resultQuery := Preloaded(baseQuery.Session(&gorm.Session{}), preload)
	resultQuery = applyForumPostSearchSort(resultQuery, options)
	if err := resultQuery.Offset(options.Offset).Limit(options.Limit).Find(&posts).Error; err != nil {
		return nil, err
	}

	return &ForumPostSearchResult{
		Posts:   posts,
		Total:   total,
		Options: options,
	}, nil
}

func (r *ForumPostRepository) buildForumPostSearchQuery(query *gorm.DB, options ForumPostSearchOptions) *gorm.DB {
	query = query.
		Select("forum_posts.*").
		Joins("JOIN forum_topics AS search_topics ON search_topics.id = forum_posts.topic_id").
		Joins("JOIN forums AS search_forums ON search_forums.id = search_topics.forum_id").
		Where("forum_posts.hidden = ?", false).
		Where("forum_posts.draft = ?", false).
		Where("forum_posts.deleted = ?", false).
		Where("search_topics.hidden = ?", false).
		Where("search_forums.hidden = ?", false)

	if options.QueryString != "" {
		query = applyForumPostTextSearch(query, options.QueryString)
	}
	if options.ForumId != nil {
		query = query.Where("search_topics.forum_id = ?", *options.ForumId)
	}
	if options.Author != "" {
		query = query.Where(`EXISTS (
			SELECT 1 FROM users AS search_authors
			WHERE search_authors.id = forum_posts.user_id
			AND search_authors.safe_name = ?
		)`, schemas.ResolveSafeName(options.Author))
	}

	if options.BookmarkedByUserId != nil {
		query = query.Where(`EXISTS (
			SELECT 1 FROM forum_bookmarks
			WHERE forum_bookmarks.topic_id = forum_posts.topic_id
			AND forum_bookmarks.user_id = ?
		)`, *options.BookmarkedByUserId)
	}
	if options.SubscribedByUserId != nil {
		query = query.Where(`EXISTS (
			SELECT 1 FROM forum_subscribers
			WHERE forum_subscribers.topic_id = forum_posts.topic_id
			AND forum_subscribers.user_id = ?
		)`, *options.SubscribedByUserId)
	}

	return query
}

func applyForumPostTextSearch(query *gorm.DB, textQuery string) *gorm.DB {
	condition, args := forumSearchVectorMatchCondition("forum_posts.search_vector", textQuery)
	return query.Where("("+condition+")", args...)
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

func applyForumPostSearchSort(query *gorm.DB, options ForumPostSearchOptions) *gorm.DB {
	descending := options.Order != constants.SearchOrderAscending

	if options.Sort == ForumSearchSortRelevance && options.QueryString != "" {
		expression, args := forumSearchVectorRankExpression("forum_posts.search_vector", options.QueryString)
		return applySearchRankOrder(query, expression, args, descending, "forum_posts.id")
	}

	return query.
		Order(clause.OrderByColumn{
			Column: clause.Column{Name: "forum_posts.created_at", Raw: true},
			Desc:   descending,
		}).
		Order("forum_posts.id DESC")
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
