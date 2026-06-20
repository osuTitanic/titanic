package constants

import "strings"

type Mods uint32

const (
	NoMod       Mods = 0
	NoFail      Mods = 1 << 0
	Easy        Mods = 1 << 1
	NoVideo     Mods = 1 << 2
	Hidden      Mods = 1 << 3
	HardRock    Mods = 1 << 4
	SuddenDeath Mods = 1 << 5
	DoubleTime  Mods = 1 << 6
	Relax       Mods = 1 << 7
	HalfTime    Mods = 1 << 8
	Nightcore   Mods = 1 << 9
	Flashlight  Mods = 1 << 10
	Autoplay    Mods = 1 << 11
	SpunOut     Mods = 1 << 12
	Autopilot   Mods = 1 << 13
	Perfect     Mods = 1 << 14
	Key4        Mods = 1 << 15
	Key5        Mods = 1 << 16
	Key6        Mods = 1 << 17
	Key7        Mods = 1 << 18
	Key8        Mods = 1 << 19
	FadeIn      Mods = 1 << 20
	Random      Mods = 1 << 21
	Cinema      Mods = 1 << 22
	Target      Mods = 1 << 23
	Key9        Mods = 1 << 24
	KeyCoop     Mods = 1 << 25
	Key1        Mods = 1 << 26
	Key3        Mods = 1 << 27
	Key2        Mods = 1 << 28
	ScoreV2     Mods = 1 << 29
	Mirror      Mods = 1 << 30
)

var modsNames = map[Mods]string{
	NoFail:      "NF",
	Easy:        "EZ",
	NoVideo:     "NV",
	Hidden:      "HD",
	HardRock:    "HR",
	SuddenDeath: "SD",
	DoubleTime:  "DT",
	Relax:       "RX",
	HalfTime:    "HT",
	Nightcore:   "NC",
	Flashlight:  "FL",
	Autoplay:    "AP",
	SpunOut:     "SO",
	Autopilot:   "AT",
	Perfect:     "PF",
	Key4:        "4K",
	Key5:        "5K",
	Key6:        "6K",
	Key7:        "7K",
	Key8:        "8K",
	FadeIn:      "FI",
	Random:      "RD",
	Cinema:      "CN",
	Target:      "TP",
	Key9:        "9K",
	KeyCoop:     "CP",
	Key1:        "1K",
	Key3:        "3K",
	Key2:        "2K",
	ScoreV2:     "V2",
	Mirror:      "MR",
}

var modsNamesReverse = map[string]Mods{
	"NF": NoFail,
	"EZ": Easy,
	"NV": NoVideo,
	"HD": Hidden,
	"HR": HardRock,
	"SD": SuddenDeath,
	"DT": DoubleTime,
	"RX": Relax,
	"HT": HalfTime,
	"NC": Nightcore,
	"FL": Flashlight,
	"AP": Autoplay,
	"SO": SpunOut,
	"AT": Autopilot,
	"PF": Perfect,
	"4K": Key4,
	"5K": Key5,
	"6K": Key6,
	"7K": Key7,
	"8K": Key8,
	"FI": FadeIn,
	"RD": Random,
	"CN": Cinema,
	"TP": Target,
	"9K": Key9,
	"CP": KeyCoop,
	"1K": Key1,
	"3K": Key3,
	"2K": Key2,
	"V2": ScoreV2,
	"MR": Mirror,
}

// Has returns true if the given mod flag is enabled in mods.
func (m Mods) Has(flag Mods) bool {
	return m&flag != 0
}

// String returns the concatenated short representation of the enabled mods, in
// bit order, e.g. "HDDT". Returns "NM" when no mods are set.
func (m Mods) String() string {
	var result strings.Builder
	for mod, name := range modsNames {
		if m.Has(mod) {
			result.WriteString(name)
		}
	}
	if result.Len() <= 0 {
		return "NM"
	}
	return result.String()
}

// ModsFromString parses a mod string such as "HDDT" into a Mods bitfield.
func ModsFromString(s string) Mods {
	// Sanitize the input string
	s = strings.NewReplacer(",", "", " ", "").Replace(s)
	s = strings.ToUpper(s)
	mods := NoMod

	for idx := 0; idx+2 <= len(s); idx += 2 {
		// Try to parse 2 characters, e.g. "HD", "DT", "NC"
		if mod, ok := modsNamesReverse[s[idx:idx+2]]; ok {
			mods |= mod
		}
	}
	return mods
}
