package templates

import (
	"reflect"

	"github.com/CloudyKit/jet/v6"
	"github.com/osuTitanic/titanic-go/internal/activity"
	"github.com/osuTitanic/titanic-go/internal/schemas"
)

func formatActivity(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("formatActivity", 1, 1)

	entry, ok := a.Get(0).Interface().(*schemas.Activity)
	if !ok || entry == nil {
		return reflect.ValueOf("")
	}

	return reflect.ValueOf(activity.RenderHtml(entry))
}
