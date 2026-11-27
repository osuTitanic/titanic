package bancho

import (
	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/services/anchor/internal/models"
)

type Client interface {
	Id() int32
	Name() string

	models.UserPresence
	models.UserStatus
	models.UserStats

	ApplyStatus(models.UserStatus)
	// ...

	// Different capabilities will be attached to each client based on what the protocol can offer
	// such as capabilities.Chat for IRC clients, or capabilities.Spectator for osu! clients
}

type ClientBase struct {
	Client

	user  *schemas.User
	stats [constants.ModeCount]*schemas.Stats
	mode  constants.Mode

	timezone  int8
	country   string // 2-letter country code
	city      string // only populated if client enabled it
	longitude float32
	latitude  float32

	statusAction          uint8 // TODO: Would be nice to have strict types for this, but chio won't like that I think ?
	statusText            string
	statusMods            uint32
	statusMode            uint8
	statusBeatmapChecksum string
	statusBeatmapId       int32
}

func (client *ClientBase) Id() int32 {
	if client.user == nil {
		return -1
	}
	return int32(client.user.Id)
}

func (client *ClientBase) Name() string {
	if client.user == nil {
		return ""
	}
	return client.user.Name
}

/*
 * Presence
 */

func (client *ClientBase) IsIrc() bool {
	return false
}

func (client *ClientBase) Permissions() uint8 {
	// TODO: Permissions should be inferred from the user struct
	return 1
}

func (client *ClientBase) Timezone() int8 {
	return client.timezone
}

func (client *ClientBase) Country() string {
	return client.country
}

func (client *ClientBase) CountryIndex() uint8 {
	return uint8(constants.GetCountryIndexFromCode(client.Country()))
}

func (client *ClientBase) City() string {
	return client.city
}

func (client *ClientBase) Longitude() float32 {
	return client.longitude // TODO: automatically reduce precision based on if city was provided
}

func (client *ClientBase) Latitude() float32 {
	return client.latitude // TODO: automatically reduce precision based on if city was provided
}

/*
 * Stats
 */

func (client *ClientBase) CurrentStats() *schemas.Stats {
	if len(client.stats) < constants.ModeCount {
		return nil
	}
	if !client.mode.Valid() {
		return nil
	}
	return client.stats[client.mode]
}

func (client *ClientBase) Rank() int32 {
	stats := client.CurrentStats()
	if stats == nil {
		return 0
	}
	return int32(stats.Rank)
}

func (client *ClientBase) RankedScore() uint64 {
	stats := client.CurrentStats()
	if stats == nil {
		return 0
	}
	return uint64(stats.Rscore)
}

func (client *ClientBase) TotalScore() uint64 {
	stats := client.CurrentStats()
	if stats == nil {
		return 0
	}
	return uint64(stats.Tscore)
}

func (client *ClientBase) Accuracy() float64 {
	stats := client.CurrentStats()
	if stats == nil {
		return 0
	}
	return stats.Acc
}

func (client *ClientBase) Playcount() int32 {
	stats := client.CurrentStats()
	if stats == nil {
		return 0
	}
	return int32(stats.Playcount)
}

func (client *ClientBase) PP() uint16 {
	stats := client.CurrentStats()
	if stats == nil {
		return 0
	}
	return uint16(stats.PP)
}

/*
 * Status
 */

func (client *ClientBase) ApplyStatus(status models.UserStatus) {
	client.statusAction = status.Action()
	client.statusText = status.Text()
	client.statusMods = status.Mods()
	client.statusMode = status.Mode()
	client.statusBeatmapChecksum = status.BeatmapChecksum()
	client.statusBeatmapId = status.BeatmapId()
}

func (client *ClientBase) Action() uint8 {
	return client.statusAction
}

func (client *ClientBase) Text() string {
	return client.statusText
}

func (client *ClientBase) Mods() uint32 {
	return client.statusMods
}

func (client *ClientBase) Mode() uint8 {
	return client.statusMode
}

func (client *ClientBase) BeatmapChecksum() string {
	return client.statusBeatmapChecksum
}

func (client *ClientBase) BeatmapId() int32 {
	return client.statusBeatmapId
}
