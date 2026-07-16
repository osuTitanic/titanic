package repositories

import (
	"strings"
	"unicode"
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
