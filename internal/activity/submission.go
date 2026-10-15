package activity

import (
	"context"
	"encoding/json"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/internal/state"
)

// Submit broadcasts the activity on the redis pubsub queue
// then stores it for profile timelines unless it is hidden.
func Submit(
	app *state.State,
	userId int,
	mode *constants.Mode, // nil -> "no specific mode"
	activityType constants.UserActivity,
	data map[string]any,
	isAnnouncement bool,
	isHidden bool,
) error {
	publishActivity(app, userId, mode, activityType, data, isAnnouncement)

	if isHidden {
		// Hidden activities are broadcast only, never stored in db
		return nil
	}

	return storeActivity(app, userId, mode, activityType, data)
}

func publishActivity(
	app *state.State,
	userId int,
	mode *constants.Mode,
	activityType constants.UserActivity,
	data map[string]any,
	isAnnouncement bool,
) {
	if err := app.BanchoEvents.Activity(
		context.Background(),
		userId,
		mode,
		activityType,
		data,
		isAnnouncement,
	); err != nil {
		app.Logger.Warn("Failed to publish activity event", "error", err, "type", activityType)
	}
}

func storeActivity(
	app *state.State,
	userId int,
	mode *constants.Mode,
	activityType constants.UserActivity,
	data map[string]any,
) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return app.Activities.Create(&schemas.Activity{
		UserId: userId,
		Mode:   mode,
		Type:   int(activityType),
		Data:   payload,
		Hidden: false,
	})
}
