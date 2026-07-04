package activity

import (
	"context"
	"encoding/json"

	"github.com/osuTitanic/titanic-go/internal/constants"
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"github.com/osuTitanic/titanic-go/internal/state"
)

const eventChannel = "bancho:events"

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
	publishEvent(app, userId, mode, activityType, data, isAnnouncement)

	if isHidden {
		// Hidden activities are broadcast only, never stored in db
		return nil
	}

	return storeEvent(app, userId, mode, activityType, data)
}

func publishEvent(
	app *state.State,
	userId int,
	mode *constants.Mode,
	activityType constants.UserActivity,
	data map[string]any,
	isAnnouncement bool,
) {
	payload, err := json.Marshal(map[string]any{
		"event": "bancho_event",
		"args":  []any{},
		"kwargs": map[string]any{
			"user_id":         userId,
			"mode":            mode,
			"type":            int(activityType),
			"data":            data,
			"is_announcement": isAnnouncement,
		},
	})
	if err != nil {
		app.Logger.Error("Failed to encode activity event", "error", err)
		return
	}

	if err := app.Redis.Publish(context.Background(), eventChannel, payload).Err(); err != nil {
		app.Logger.Warn("Failed to publish activity event", "error", err, "type", activityType)
	}
}

func storeEvent(
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
