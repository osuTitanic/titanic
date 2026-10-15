package bancho

import (
	"context"
	"fmt"
	"strconv"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/redis/go-redis/v9"
)

type Status struct {
	UserId          int                    `redis:"-"`
	Action          constants.ClientStatus `redis:"action"`
	Text            string                 `redis:"text"`
	Mods            constants.Mods         `redis:"mods"`
	Mode            constants.Mode         `redis:"mode"`
	BeatmapChecksum string                 `redis:"beatmap_checksum"`
	BeatmapId       int                    `redis:"beatmap_id"`
	ClientHash      string                 `redis:"hash"`
	ClientVersion   int                    `redis:"version"`
	ClientString    string                 `redis:"version_string"`
}

var requiredStatusFields = [...]string{
	"action",
	"text",
	"mods",
	"mode",
	"beatmap_checksum",
	"beatmap_id",
	"hash",
	"version",
	"version_string",
}

type StatusStore struct {
	client redis.Cmdable
}

func NewStatusStore(client redis.Cmdable) *StatusStore {
	return &StatusStore{client: client}
}

func (store *StatusStore) Exists(ctx context.Context, userId int) (bool, error) {
	exists, err := store.client.Exists(ctx, statusKey(userId)).Result()
	if err != nil {
		return false, fmt.Errorf("check bancho status for user %d: %w", userId, err)
	}
	return exists > 0, nil
}

func (store *StatusStore) Get(ctx context.Context, userId int) (*Status, error) {
	cmd := store.client.HGetAll(ctx, statusKey(userId))

	values, err := cmd.Result()
	if err != nil {
		return nil, fmt.Errorf("get bancho status for user %d: %w", userId, err)
	}
	if len(values) == 0 {
		return nil, nil
	}

	// Ensure fields are not missing (scan doesn't care if they are)
	for _, field := range requiredStatusFields {
		if _, ok := values[field]; !ok {
			return nil, invalidStatus("missing %s", field)
		}
	}

	status := &Status{
		UserId: userId,
	}
	if err := cmd.Scan(status); err != nil {
		return nil, invalidStatus("%w", err)
	}

	if !status.Action.Valid() {
		return nil, invalidStatus("invalid action %d", status.Action)
	}
	if !status.Mods.Valid() {
		return nil, invalidStatus("invalid mods %d", status.Mods)
	}
	if !status.Mode.Valid() {
		return nil, invalidStatus("invalid mode %d", status.Mode)
	}

	return status, nil
}

func statusKey(userId int) string {
	return "bancho:status:" + strconv.Itoa(userId)
}

func invalidStatus(format string, a ...any) error {
	return fmt.Errorf("invalid bancho status: %w", fmt.Errorf(format, a...))
}
