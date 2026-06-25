package activity

import (
	"encoding/json"
	"fmt"
	"html"

	"github.com/osuTitanic/titanic-go/internal/schemas"
)

// parseData unmarshals the activity json data.
func parseData(entry *schemas.Activity) (map[string]any, bool) {
	if entry == nil {
		return nil, false
	}

	var data map[string]any
	if err := json.Unmarshal(entry.Data, &data); err != nil {
		return nil, false
	}
	return data, true
}

// dataString reads a value from the activity data & returns it as a string.
func dataString(data map[string]any, key string) string {
	value, ok := data[key]
	if !ok || value == nil {
		return ""
	}

	switch value := value.(type) {
	case string:
		return value
	case float64:
		if value == float64(int64(value)) {
			return fmt.Sprintf("%d", int64(value))
		}
		return fmt.Sprintf("%v", value)
	default:
		return fmt.Sprintf("%v", value)
	}
}

// dataInt reads a value from the activity data & returns it as an integer.
func dataInt(data map[string]any, key string) int {
	if value, ok := data[key].(float64); ok {
		return int(value)
	}
	return 0
}

// htmlText reads & escapes a plain string value from the activity data.
func htmlText(data map[string]any, key string) string {
	return html.EscapeString(dataString(data, key))
}

// htmlLink builds an escaped a tag, returning "" when there is no link text.
func htmlLink(url string, text string) string {
	if text == "" {
		return ""
	}
	return fmt.Sprintf("<a href=\"%s\">%s</a>", html.EscapeString(url), html.EscapeString(text))
}
