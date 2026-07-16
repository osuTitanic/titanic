package repositories

import (
	"slices"
	"strings"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BeatmapsetSearchResult struct {
	Beatmapsets []*schemas.Beatmapset
	Total       int64
	Options     BeatmapsetSearchOptions
}

type BeatmapsetSearchOptions struct {
	QueryString string

	Order    constants.SearchOrder
	Category constants.BeatmapCategory
	Sort     constants.BeatmapSort

	Genre    *constants.BeatmapGenre
	Language *constants.BeatmapLanguage
	Mode     *constants.Mode

	HasStoryboard bool
	HasVideo      bool
	TitanicOnly   bool

	Played    bool
	Unplayed  bool
	Cleared   bool
	Uncleared bool
	UserId    *int

	Offset int
	Limit  int
}

// Normalize ensures that the search options are valid & sets defaults where necessary
func (options *BeatmapsetSearchOptions) Normalize() {
	defaulted := options.Limit == 0
	if defaulted && options.Category == constants.BeatmapCategoryAny {
		options.Category = constants.BeatmapCategoryLeaderboard
	}
	if defaulted && options.Sort == constants.BeatmapSortTitle {
		options.Sort = constants.BeatmapSortRanked
	}
	if options.Category < constants.BeatmapCategoryAny || options.Category > constants.BeatmapCategoryGraveyard {
		options.Category = constants.BeatmapCategoryLeaderboard
	}
	if options.Sort < constants.BeatmapSortTitle || options.Sort > constants.BeatmapSortRelevance {
		options.Sort = constants.BeatmapSortRanked
	}
	if options.Order != constants.SearchOrderAscending {
		options.Order = constants.SearchOrderDescending
	}
	if options.Limit < 1 {
		options.Limit = 50
	}
	if options.Limit > 50 {
		options.Limit = 50
	}
	if options.Offset < 0 {
		options.Offset = 0
	}
	if options.Played {
		options.Unplayed = false
	}
	if options.Cleared {
		options.Uncleared = false
	}
}

// searchJoins tracks which tables need to be joined in the search query
type searchJoins struct {
	beatmaps bool
	plays    bool
	scores   bool
}

func (joins *searchJoins) Apply(query *gorm.DB) *gorm.DB {
	if joins.beatmaps {
		query = query.Joins("JOIN beatmaps ON beatmaps.set_id = beatmapsets.id").Group("beatmapsets.id")
	}
	if joins.plays {
		query = query.Joins("JOIN plays ON plays.beatmap_id = beatmaps.id")
	}
	if joins.scores {
		query = query.Joins("JOIN scores ON scores.beatmap_id = beatmaps.id")
	}
	return query
}

// A set matches a mode when at least one of its beatmaps uses that mode.
const beatmapsetModeCondition = `EXISTS (
	SELECT 1 FROM beatmaps AS mode_maps
	WHERE mode_maps.set_id = beatmapsets.id
	AND mode_maps.mode = ?
)`

// A set is played when at least one of its beatmaps has been recorded as a "plays" db entry.
const beatmapsetPlayedCondition = `EXISTS (
	SELECT 1 FROM plays
	WHERE plays.set_id = beatmapsets.id
	AND plays.user_id = ?
)`

// Mode-specific variant of beatmapsetPlayedCondition.
const beatmapsetPlayedInModeCondition = `EXISTS (
	SELECT 1 FROM beatmaps AS played_maps
	WHERE played_maps.set_id = beatmapsets.id
	AND played_maps.mode = ?
	AND EXISTS (
		SELECT 1 FROM plays
		WHERE plays.beatmap_id = played_maps.id
		AND plays.user_id = ?
	)
)`

// A set is cleared when the user has a submitted score on each of its beatmaps.
const beatmapsetClearedCondition = `beatmapsets.id IN (
	SELECT scored_maps.set_id
	FROM scores AS clear_scores
	JOIN beatmaps AS scored_maps ON scored_maps.id = clear_scores.beatmap_id
	WHERE clear_scores.user_id = ?
	AND clear_scores.status >= ?
	GROUP BY scored_maps.set_id
	HAVING COUNT(DISTINCT scored_maps.id) = (
		SELECT COUNT(*) FROM beatmaps AS all_maps
		WHERE all_maps.set_id = scored_maps.set_id
	)
)`

// Mode-specific variant of beatmapsetClearedCondition.
const beatmapsetClearedInModeCondition = `beatmapsets.id IN (
	SELECT scored_maps.set_id
	FROM scores AS clear_scores
	JOIN beatmaps AS scored_maps ON scored_maps.id = clear_scores.beatmap_id
	WHERE scored_maps.mode = ?
	AND clear_scores.user_id = ?
	AND clear_scores.status >= ?
	GROUP BY scored_maps.set_id
	HAVING COUNT(DISTINCT scored_maps.id) = (
		SELECT COUNT(*) FROM beatmaps AS all_maps
		WHERE all_maps.set_id = scored_maps.set_id
		AND all_maps.mode = ?
	)
)`

var categoryToStatusMapping = map[constants.BeatmapCategory]constants.BeatmapStatus{
	constants.BeatmapCategoryRanked:    constants.BeatmapStatusRanked,
	constants.BeatmapCategoryApproved:  constants.BeatmapStatusApproved,
	constants.BeatmapCategoryQualified: constants.BeatmapStatusQualified,
	constants.BeatmapCategoryLoved:     constants.BeatmapStatusLoved,
	constants.BeatmapCategoryPending:   constants.BeatmapStatusPending,
	constants.BeatmapCategoryWIP:       constants.BeatmapStatusWIP,
	constants.BeatmapCategoryGraveyard: constants.BeatmapStatusGraveyard,
}

var beatmapsetSortExpressions = map[constants.BeatmapSort]string{
	constants.BeatmapSortTitle:      "beatmapsets.title",
	constants.BeatmapSortArtist:     "beatmapsets.artist",
	constants.BeatmapSortCreator:    "beatmapsets.creator",
	constants.BeatmapSortDifficulty: "beatmapsets.max_diff",
	constants.BeatmapSortPlays:      "beatmapsets.total_playcount",
	constants.BeatmapSortCreated:    "beatmapsets.id",
	// For ratings, we calculate a Bayesian average to provide a more balanced ranking
	// that accounts for both the average rating and the number of ratings
	// https://en.wikipedia.org/wiki/Bayesian_average
	constants.BeatmapSortRating: `
		((beatmapsets.rating_average * beatmapsets.rating_count) + (COALESCE((SELECT AVG(rating) FROM ratings), 0) * 10))
		/ (beatmapsets.rating_count + 10)
	`,
}

func (r *BeatmapsetRepository) Search(options BeatmapsetSearchOptions, preload ...string) ([]*schemas.Beatmapset, error) {
	result, err := r.SearchPage(options, preload...)
	if err != nil {
		return nil, err
	}
	return result.Beatmapsets, nil
}

func (r *BeatmapsetRepository) SearchPage(options BeatmapsetSearchOptions, preload ...string) (*BeatmapsetSearchResult, error) {
	var total int64
	options.Normalize()

	// Build the base query with all filters applied, but without pagination
	baseQuery, textQuery := r.buildBeatmapsetSearchQuery(r.db.Model(&schemas.Beatmapset{}), options)

	// Count total results before applying pagination
	countQuery := r.db.Table("(?) AS filtered_beatmapsets", baseQuery.Session(&gorm.Session{}).Select("beatmapsets.id"))
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	// If the offset is beyond the total number of results, clamp it to the last valid page
	options.Offset = clampSearchOffset(options.Offset, options.Limit, total)

	var beatmapsets []*schemas.Beatmapset
	resultQuery := Preloaded(baseQuery.Session(&gorm.Session{}), preload)
	resultQuery = applyBeatmapsetSearchSort(resultQuery, options, textQuery)

	// Apply pagination & execute the query
	err := resultQuery.Offset(options.Offset).Limit(options.Limit).Find(&beatmapsets).Error
	if err != nil {
		return nil, err
	}

	return &BeatmapsetSearchResult{
		Beatmapsets: beatmapsets,
		Total:       total,
		Options:     options,
	}, nil
}

func (r *BeatmapsetRepository) buildBeatmapsetSearchQuery(query *gorm.DB, options BeatmapsetSearchOptions) (*gorm.DB, string) {
	var joins searchJoins

	query = query.Where("EXISTS (SELECT 1 FROM beatmaps WHERE beatmaps.set_id = beatmapsets.id)")
	query = applyBeatmapsetCriteria(query, options)
	query = applyBeatmapsetSearchCategory(query, options.Category)

	query, updatedSearchQuery, filtersRequireBeatmaps := applyBeatmapsetSearchFilters(
		query,
		options.QueryString,
		options.Mode,
	)
	joins.beatmaps = joins.beatmaps || filtersRequireBeatmaps

	if strings.TrimSpace(updatedSearchQuery) != "" {
		query = applyBeatmapsetTextSearch(query, updatedSearchQuery)
	}

	query = applyBeatmapsetUserFilters(query, options)
	query = joins.Apply(query)
	return query, updatedSearchQuery
}

func applyBeatmapsetCriteria(query *gorm.DB, options BeatmapsetSearchOptions) *gorm.DB {
	if options.Genre != nil && *options.Genre != constants.BeatmapGenreAny {
		query = query.Where("beatmapsets.genre_id = ?", *options.Genre)
	}
	if options.Language != nil && *options.Language != constants.BeatmapLanguageAny {
		query = query.Where("beatmapsets.language_id = ?", *options.Language)
	}
	if options.HasStoryboard {
		query = query.Where("beatmapsets.has_storyboard = true")
	}
	if options.HasVideo {
		query = query.Where("beatmapsets.has_video = true")
	}
	if options.TitanicOnly {
		query = query.Where("(beatmapsets.server = ? OR beatmapsets.enhanced = true)", constants.BeatmapServerTitanic)
	}
	if options.Mode != nil {
		// Keep mode independent from any joined beatmap used by another filter
		query = query.Where(beatmapsetModeCondition, *options.Mode)
	}
	return query
}

func applyBeatmapsetUserFilters(query *gorm.DB, options BeatmapsetSearchOptions) *gorm.DB {
	if options.UserId == nil {
		return query
	}

	playedCondition := beatmapsetPlayedCondition
	playedArgs := []any{*options.UserId}
	clearedCondition := beatmapsetClearedCondition
	clearedArgs := []any{*options.UserId, constants.ScoreStatusSubmitted}

	// If a mode is specified, the cleared/uncleared conditions have to change as well
	// e.g. a beatmap has 2 osu!standard diffs and 3 osu!taiko diffs
	// 		a user plays only both osu!standard diffs & filters by uncleared + osu!standard
	// 		in that scenario, we want the set to be marked as "cleared"
	if options.Mode != nil {
		playedCondition = beatmapsetPlayedInModeCondition
		playedArgs = []any{*options.Mode, *options.UserId}
		clearedCondition = beatmapsetClearedInModeCondition
		clearedArgs = []any{*options.Mode, *options.UserId, constants.ScoreStatusSubmitted, *options.Mode}
	}

	if options.Played {
		query = query.Where(playedCondition, playedArgs...)
	}
	if options.Unplayed {
		query = query.Where("NOT ("+playedCondition+")", playedArgs...)
	}
	if options.Cleared {
		query = query.Where(clearedCondition, clearedArgs...)
	}
	if options.Uncleared {
		query = query.Where("NOT ("+clearedCondition+")", clearedArgs...)
	}
	return query
}

func applyBeatmapsetSearchCategory(query *gorm.DB, category constants.BeatmapCategory) *gorm.DB {
	if category == constants.BeatmapCategoryAny {
		// "anything but inactive" category
		return query.Where(
			"beatmapsets.submission_status != ?",
			constants.BeatmapStatusInactive,
		)
	}

	if status, ok := categoryToStatusMapping[category]; ok {
		return query.Where("beatmapsets.submission_status = ?", status)
	}

	// "has leaderboard" category
	return query.Where("beatmapsets.submission_status > 0")
}

func applyBeatmapsetSearchSort(query *gorm.DB, options BeatmapsetSearchOptions, textQuery string) *gorm.DB {
	desc := options.Order != constants.SearchOrderAscending
	textQuery = strings.TrimSpace(textQuery)

	if options.Sort == constants.BeatmapSortRelevance && textQuery != "" {
		// When sorting by relevance, we use the ts_rank of the full-text
		// search vector against the query as the primary sort key
		return applySearchRankOrder(
			query,
			"ts_rank(beatmapsets.search, plainto_tsquery('simple', ?))",
			[]any{textQuery},
			desc,
			"beatmapsets.id",
		)
	}

	// By default, we sort by approved/ranked date
	sortExpression := "beatmapsets.approved_date"
	if expression, ok := beatmapsetSortExpressions[options.Sort]; ok {
		sortExpression = expression
	}

	// For non-relevance sorts, we sort by the specified field
	return query.
		Order(clause.OrderByColumn{
			Column: clause.Column{Name: sortExpression, Raw: true},
			Desc:   desc,
		}).
		Order("beatmapsets.id DESC")
}

func applyBeatmapsetTextSearch(query *gorm.DB, textQuery string) *gorm.DB {
	// Start with a plain full-text search, which is fast and
	// effective for exact matches of the query terms
	setCondition := "sets.search @@ plainto_tsquery('simple', ?)"
	mapCondition := "maps.search @@ plainto_tsquery('simple', ?)"
	args := []any{textQuery}

	// Add an optional fuzzy full-text search, which can help find results
	// when some terms are slightly different from the indexed text
	if fuzzyQuery := fuzzyTsQuery(textQuery); fuzzyQuery != "" {
		setCondition += " OR sets.search @@ to_tsquery('simple', ?)"
		mapCondition += " OR maps.search @@ to_tsquery('simple', ?)"
		args = append(args, fuzzyQuery)
	}

	// Keeping each table in its own subquery lets postgres use the
	// GIN indexes on both "search" columns
	subquery := `beatmapsets.id IN (
		SELECT sets.id FROM beatmapsets AS sets WHERE ` + setCondition + `
		UNION
		SELECT maps.set_id FROM beatmaps AS maps WHERE ` + mapCondition + `
	)`

	return query.Where(subquery, slices.Concat(args, args)...)
}
