package repositories

import (
	"strings"
	"unicode"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func fuzzyTsQuery(query string) string {
	words := strings.FieldsFunc(strings.ToLower(query), func(r rune) bool {
		return !unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
	})

	for i := range words {
		words[i] += ":*"
	}
	return strings.Join(words, " & ")
}

// clampSearchOffset keeps pagination on the final valid page when results shrink
func clampSearchOffset(offset, limit int, total int64) int {
	if total == 0 || int64(offset) < total {
		return offset
	}
	last := total - 1
	return int(last - last%int64(limit))
}

// applySearchRankOrder orders results by the computed search rank & applies a descending tie-breaker
func applySearchRankOrder(query *gorm.DB, expression string, vars []any, descending bool, tieBreaker string) *gorm.DB {
	direction := "ASC"
	if descending {
		direction = "DESC"
	}

	return query.Order(clause.OrderBy{Expression: clause.Expr{
		SQL:  expression + " " + direction + ", " + tieBreaker + " DESC",
		Vars: vars,
	}})
}
