package constants

import "fmt"

type Mode int8

const (
	ModeOsu   Mode = 0
	ModeTaiko Mode = 1
	ModeCatch Mode = 2
	ModeMania Mode = 3
)

var Modes = []Mode{ModeOsu, ModeTaiko, ModeCatch, ModeMania}

func (m Mode) Value() int {
	return int(m)
}

func (m Mode) String() string {
	switch m {
	case ModeOsu:
		return "osu!"
	case ModeTaiko:
		return "Taiko"
	case ModeCatch:
		return "CatchTheBeat"
	case ModeMania:
		return "osu!mania"
	default:
		return fmt.Sprintf("%d", m)
	}
}

func (m Mode) Alias() string {
	switch m {
	case ModeOsu:
		return "osu"
	case ModeTaiko:
		return "taiko"
	case ModeCatch:
		return "catch"
	case ModeMania:
		return "mania"
	default:
		return fmt.Sprintf("%d", m)
	}
}

func NewModeFromAlias(alias string) (Mode, bool) {
	switch alias {
	case "osu":
		return ModeOsu, true
	case "taiko":
		return ModeTaiko, true
	case "catch":
		return ModeCatch, true
	case "mania":
		return ModeMania, true
	default:
		return 0, false
	}
}
