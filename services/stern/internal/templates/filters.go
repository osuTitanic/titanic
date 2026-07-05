package templates

import (
	"fmt"
	"html"
	"math"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/osuTitanic/titanic-go/internal/bbcode"
	"github.com/osuTitanic/titanic-go/internal/constants"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var markdownLinkPattern = regexp.MustCompile(`\[([^\]]+)\]\((https?://[^\s)]+)\)`)
var printer = message.NewPrinter(language.English)

func formatNumber(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("formatNumber", 1, 1)

	var result string
	var value any = a.Get(0).Interface()

	switch value := reflect.ValueOf(value); value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		result = printer.Sprintf("%d", value.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		result = printer.Sprintf("%d", value.Uint())
	default:
		result = fmt.Sprint(value.Interface())
	}

	return reflect.ValueOf(result)
}

func round(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("round", 1, 1)

	rounded := math.Round(reflectFloat(a.Get(0)))
	return reflect.ValueOf(int64(rounded))
}

func floor(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("floor", 1, 1)

	floored := math.Floor(reflectFloat(a.Get(0)))
	return reflect.ValueOf(int64(floored))
}

func formatFloat(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("formatFloat", 2, 2)

	value := reflectFloat(a.Get(0))
	places := reflectFloat(a.Get(1))
	return reflect.ValueOf(printer.Sprintf("%.*f", int(places), value))
}

func formatDateShort(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("formatDateShort", 1, 1)

	value := a.Get(0).Interface()
	switch value := value.(type) {
	case time.Time:
		return reflect.ValueOf(value.Format("Jan 2, 2006"))
	case *time.Time:
		if value == nil {
			return reflect.ValueOf("")
		}
		return reflect.ValueOf(value.Format("Jan 2, 2006"))
	default:
		return reflect.ValueOf("")
	}
}

func countryName(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("countryName", 1, 1)

	value := a.Get(0).Interface()
	if code, ok := value.(string); ok {
		return reflect.ValueOf(constants.GetCountryNameFromCode(code))
	}

	return reflect.ValueOf("")
}

func renderBBCode(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("bbcode", 1, 1)

	input, ok := a.Get(0).Interface().(string)
	if !ok {
		return reflect.ValueOf("")
	}
	return reflect.ValueOf(bbcode.RenderHtml(input))
}

func shortMods(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("shortMods", 1, 1)

	mods, ok := a.Get(0).Interface().(constants.Mods)
	if !ok {
		return reflect.ValueOf("None")
	}
	if short := mods.String(); short != "NM" {
		return reflect.ValueOf(short)
	}
	return reflect.ValueOf("None")
}

func markdownUrls(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("markdownUrls", 1, 1)

	input, ok := a.Get(0).Interface().(string)
	if !ok {
		return reflect.ValueOf("")
	}

	var builder strings.Builder
	lastEnd := 0

	for _, match := range markdownLinkPattern.FindAllStringSubmatchIndex(input, -1) {
		start, end := match[0], match[1]

		// Extract the text and href from the match
		text := input[match[2]:match[3]]
		href := input[match[4]:match[5]]

		// Write the text before the match and the link itself
		builder.WriteString(html.EscapeString(input[lastEnd:start]))

		// Write the link in HTML format
		builder.WriteString(`<a href="`)
		builder.WriteString(html.EscapeString(href))
		builder.WriteString(`">`)
		builder.WriteString(html.EscapeString(text))
		builder.WriteString(`</a>`)
		lastEnd = end
	}

	// Write the remaining text after the last match
	builder.WriteString(html.EscapeString(input[lastEnd:]))
	return reflect.ValueOf(builder.String())
}

func reflectFloat(value reflect.Value) float64 {
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(value.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return float64(value.Uint())
	case reflect.Float32, reflect.Float64:
		return value.Float()
	default:
		return 0
	}
}
