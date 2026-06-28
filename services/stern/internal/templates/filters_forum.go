package templates

import (
	"fmt"
	"html"
	"reflect"

	"github.com/CloudyKit/jet/v6"
	"github.com/osuTitanic/titanic-go/internal/schemas"
)

// forumUserLink renders a colored username link for forum listings
func forumUserLink(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("forumUserLink", 1, 2)

	user, ok := a.Get(0).Interface().(*schemas.User)
	if !ok || user == nil {
		return reflect.ValueOf("")
	}

	color := user.DisplayColor()
	if a.NumOfArguments() > 1 {
		currentUserId := int(reflectFloat(a.Get(1)))
		if color == "" && currentUserId == user.Id {
			// Highlight the visitor's own username
			color = "#310e7a"
		}
	}

	name := html.EscapeString(user.Name)
	if color != "" {
		return reflect.ValueOf(fmt.Sprintf(
			`<a href="/u/%d" class="username-colored" style="color: %s;">%s</a>`,
			user.Id, color, name,
		))
	}
	return reflect.ValueOf(fmt.Sprintf(`<a href="/u/%d">%s</a>`, user.Id, name))
}
