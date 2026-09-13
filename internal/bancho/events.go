package bancho

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/redis/go-redis/v9"
)

const eventChannel = "bancho:events"

type Event struct {
	Event  string `json:"event"`
	Args   []any  `json:"args"`
	Kwargs any    `json:"kwargs"`
}

type EventDispatcher struct {
	client redis.Cmdable
}

func NewEventDispatcher(client redis.Cmdable) *EventDispatcher {
	return &EventDispatcher{client: client}
}

func (dispatcher *EventDispatcher) Dispatch(ctx context.Context, name string, args []any, kwargs any) error {
	if args == nil {
		args = []any{}
	}
	if kwargs == nil {
		kwargs = map[string]any{}
	}
	event := Event{
		Event:  name,
		Args:   args,
		Kwargs: kwargs,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("encode bancho event %q: %w", name, err)
	}
	if err := dispatcher.client.Publish(ctx, eventChannel, payload).Err(); err != nil {
		return fmt.Errorf("publish bancho event %q: %w", name, err)
	}
	return nil
}

func (dispatcher *EventDispatcher) Activity(
	ctx context.Context,
	userId int,
	mode *constants.Mode,
	activity constants.UserActivity,
	data map[string]any,
	isAnnouncement bool,
) error {
	if data == nil {
		data = map[string]any{}
	}
	return dispatcher.Dispatch(ctx, "bancho_event", nil, struct {
		UserId         int                    `json:"user_id"`
		Mode           *constants.Mode        `json:"mode"`
		Type           constants.UserActivity `json:"type"`
		Data           map[string]any         `json:"data"`
		IsAnnouncement bool                   `json:"is_announcement"`
	}{
		UserId:         userId,
		Mode:           mode,
		Type:           activity,
		Data:           data,
		IsAnnouncement: isAnnouncement,
	})
}

func (dispatcher *EventDispatcher) UserUpdate(ctx context.Context, userId int, mode constants.Mode) error {
	return dispatcher.Dispatch(ctx, "user_update", nil, struct {
		UserId int            `json:"user_id"`
		Mode   constants.Mode `json:"mode"`
	}{
		UserId: userId,
		Mode:   mode,
	})
}

// TODO: Add other bancho events
