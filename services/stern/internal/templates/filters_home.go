package templates

import (
	"reflect"
	"regexp"

	"github.com/CloudyKit/jet/v6"
	"github.com/osuTitanic/titanic-go/internal/schemas"
)

var homeNewsIgnoredTags = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\[/?b\]`),
	regexp.MustCompile(`(?i)\[/?centre\]`),
	regexp.MustCompile(`(?i)\[size(?:=[^\]]+)?\]`),
	regexp.MustCompile(`(?i)\[/size\]`),
	regexp.MustCompile(`(?i)\[/?heading\]`),
	regexp.MustCompile(`(?i)\[/?img\]`),
}

func homeRenderNewsText(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("homeRenderNewsText", 1, 1)
	post := a.Get(0).Interface().(schemas.ForumPost)

	content := post.RenderForNews(homeNewsIgnoredTags...)
	return reflect.ValueOf(content)
}
