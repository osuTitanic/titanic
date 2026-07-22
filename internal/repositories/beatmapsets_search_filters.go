package repositories

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/osuTitanic/titanic/internal/constants"
	"gorm.io/gorm"
)

var (
	// A list of filters that can be applied through the search query
	// Example: bpm>180, ar<=9, status=ranked, ...
	searchFilters = map[string]searchFilter{
		"status":     statusSearchFilter,
		"artist":     stringSearchFilter("beatmapsets.artist", "beatmapsets.artist_unicode"),
		"title":      stringSearchFilter("beatmapsets.title", "beatmapsets.title_unicode"),
		"creator":    stringSearchFilter("beatmapsets.creator"),
		"source":     stringSearchFilter("beatmapsets.source", "beatmapsets.source_unicode"),
		"ranked":     dateSearchFilter("beatmapsets.approved_date"),
		"created":    dateSearchFilter("beatmapsets.submission_date"),
		"year":       yearSearchFilter("EXTRACT(YEAR FROM beatmapsets.approved_date)"),
		"difficulty": numberSearchFilter("beatmaps.diff"),
		"diff":       numberSearchFilter("beatmaps.diff"),
		"stars":      numberSearchFilter("beatmaps.diff"),
		"sr":         numberSearchFilter("beatmaps.diff"),
		"length":     numberSearchFilter("beatmaps.total_length"),
		"drain":      numberSearchFilter("beatmaps.drain_length"),
		"bpm":        numberSearchFilter("beatmaps.bpm"),
		"ar":         numberSearchFilter("beatmaps.ar"),
		"cs":         numberSearchFilter("beatmaps.cs"),
		"od":         numberSearchFilter("beatmaps.od"),
		"hp":         numberSearchFilter("beatmaps.hp"),
	}

	// Mapping of beatmap status names to their corresponding constants, used in `parseBeatmapStatus`
	beatmapStatusByName = map[string]constants.BeatmapStatus{
		"inactive":         constants.BeatmapStatusInactive,
		"graveyard":        constants.BeatmapStatusGraveyard,
		"graveyarded":      constants.BeatmapStatusGraveyard,
		"wip":              constants.BeatmapStatusWIP,
		"work_in_progress": constants.BeatmapStatusWIP,
		"pending":          constants.BeatmapStatusPending,
		"ranked":           constants.BeatmapStatusRanked,
		"approved":         constants.BeatmapStatusApproved,
		"qualified":        constants.BeatmapStatusQualified,
		"loved":            constants.BeatmapStatusLoved,
	}

	// Supported operators for number-based filters
	searchNumberOperators = map[string]string{
		"=":  "=",
		"!=": "<>",
		">":  ">",
		"<":  "<",
		">=": ">=",
		"<=": "<=",
	}

	// Regex pattern to identify filters in the search query
	searchFilterPattern = regexp.MustCompile(`(?i)(^|\s)([a-z]+)\s*(>=|<=|!=|=|>|<)\s*("[^"]*"|[^\s]+)`)
)

// searchFilterCondition represents the parsed components of a search filter,
// such as "bpm>180" or "status=ranked"
type searchFilterCondition struct {
	Operator string
	Value    string
	Mode     *constants.Mode
}

// searchFilter is a function type that applies a specific filter condition to a gorm query
type searchFilter func(*gorm.DB, searchFilterCondition) (*gorm.DB, bool, error)

// applyBeatmapsetSearchFilters parses the search query for any filters and applies them to the gorm query
func applyBeatmapsetSearchFilters(query *gorm.DB, queryString string, mode *constants.Mode) (*gorm.DB, string, bool) {
	joinBeatmaps := false

	textQuery := searchFilterPattern.ReplaceAllStringFunc(queryString, func(match string) string {
		name, condition, leadingSpace, ok := parseSearchFilterMatch(match, mode)
		if !ok {
			return match
		}
		filter, ok := searchFilters[name]
		if !ok {
			return match
		}

		nextQuery, needsBeatmaps, err := filter(query, condition)
		if err != nil {
			return match
		}
		query = nextQuery
		joinBeatmaps = joinBeatmaps || needsBeatmaps
		return leadingSpace
	})

	return query, textQuery, joinBeatmaps
}

func parseSearchFilterMatch(match string, mode *constants.Mode) (string, searchFilterCondition, string, bool) {
	parts := searchFilterPattern.FindStringSubmatch(match)
	if len(parts) != 5 {
		return "", searchFilterCondition{}, "", false
	}

	condition := searchFilterCondition{
		Operator: parts[3],
		Value:    strings.Trim(parts[4], `"`),
		Mode:     mode,
	}
	return strings.ToLower(parts[2]), condition, parts[1], true
}

func stringSearchFilter(columns ...string) searchFilter {
	return func(query *gorm.DB, condition searchFilterCondition) (*gorm.DB, bool, error) {
		var operator string

		switch condition.Operator {
		case "=":
			operator = "ILIKE"
		case "!=":
			operator = "NOT ILIKE"
		default:
			return query, false, fmt.Errorf("invalid string operator")
		}

		pattern := "%" + condition.Value + "%"
		clauses := make([]string, len(columns))
		args := make([]any, len(columns))

		for i, column := range columns {
			clauses[i] = column + " " + operator + " ?"
			args[i] = pattern
		}

		return query.Where(
			// Combine multiple columns with OR
			"("+strings.Join(clauses, " OR ")+")",
			args...,
		), false, nil
	}
}

func statusSearchFilter(query *gorm.DB, condition searchFilterCondition) (*gorm.DB, bool, error) {
	status, ok := parseBeatmapStatus(condition.Value)
	if !ok {
		return query, false, fmt.Errorf("invalid status")
	}
	return applyNumberCondition(query, "beatmapsets.submission_status", condition.Operator, status, false)
}

func yearSearchFilter(expression string) searchFilter {
	return func(query *gorm.DB, condition searchFilterCondition) (*gorm.DB, bool, error) {
		year, err := strconv.Atoi(condition.Value)
		if err != nil {
			return query, false, err
		}
		return applyNumberCondition(query, expression, condition.Operator, year, false)
	}
}

func dateSearchFilter(column string) searchFilter {
	return func(query *gorm.DB, condition searchFilterCondition) (*gorm.DB, bool, error) {
		// Parse date as "YYYY", "YYYY-MM", or "YYYY-MM-DD" and
		// get the start & end of that range
		start, end, err := parseSearchDateRange(condition.Value)
		if err != nil {
			return query, false, err
		}

		switch condition.Operator {
		case "=":
			// Between start & end
			query = query.Where(column+" >= ? AND "+column+" < ?", start, end)
		case "!=":
			// Outside of start & end range
			query = query.Where("("+column+" < ? OR "+column+" >= ?)", start, end)
		case ">":
			// After end
			query = query.Where(column+" >= ?", end)
		case ">=":
			// After start
			query = query.Where(column+" >= ?", start)
		case "<":
			// Before start
			query = query.Where(column+" < ?", start)
		case "<=":
			// Before end
			query = query.Where(column+" < ?", end)
		default:
			return query, false, fmt.Errorf("invalid date operator")
		}

		return query, false, nil
	}
}

func numberSearchFilter(column string) searchFilter {
	return func(query *gorm.DB, condition searchFilterCondition) (*gorm.DB, bool, error) {
		value, err := strconv.ParseFloat(condition.Value, 64)
		if err != nil {
			return query, false, err
		}

		query, _, err = applyNumberCondition(query, column, condition.Operator, value, true)
		if err != nil {
			return query, false, err
		}

		if condition.Mode != nil && isModeSpecificBeatmapStat(column) {
			// If the filter is AR, CS, OD, HP & a mode is specified, we need to filter
			// beatmaps by that mode as well to ensure the stats are relevant
			query = query.Where("beatmaps.mode = ?", *condition.Mode)
		}
		return query, true, nil
	}
}

func isModeSpecificBeatmapStat(column string) bool {
	switch column {
	case "beatmaps.ar", "beatmaps.cs", "beatmaps.od", "beatmaps.hp":
		return true
	default:
		return false
	}
}

func applyNumberCondition(query *gorm.DB, expression, operator string, value any, needsBeatmaps bool) (*gorm.DB, bool, error) {
	sqlOperator, ok := searchNumberOperators[operator]
	if !ok {
		return query, false, fmt.Errorf("invalid operator")
	}
	return query.Where(expression+" "+sqlOperator+" ?", value), needsBeatmaps, nil
}

func parseBeatmapStatus(value string) (constants.BeatmapStatus, bool) {
	// First we try to parse by integer
	if status, err := strconv.Atoi(value); err == nil {
		return constants.BeatmapStatus(status), true
	}

	// Its more likely for users to search by the status name, so we support that too
	status, ok := beatmapStatusByName[strings.ToLower(value)]
	return status, ok
}

func parseSearchDateRange(value string) (time.Time, time.Time, error) {
	var layout string
	var end func(time.Time) time.Time

	switch len(value) {
	// We support three date formats: YYYY, YYYY-MM, and YYYY-MM-DD
	case len("2006"):
		layout = "2006"
		end = func(date time.Time) time.Time { return date.AddDate(1, 0, 0) }
	case len("2006-01"):
		layout = "2006-01"
		end = func(date time.Time) time.Time { return date.AddDate(0, 1, 0) }
	case len("2006-01-02"):
		layout = time.DateOnly
		end = func(date time.Time) time.Time { return date.AddDate(0, 0, 1) }
	default:
		return time.Time{}, time.Time{}, fmt.Errorf("invalid date %q", value)
	}

	start, err := time.Parse(layout, value)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return start, end(start), nil
}
