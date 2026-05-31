package constants

import "strings"

type Playstyle uint8

const (
	PlaystyleNotSpecified Playstyle = 0
	PlaystyleMouse        Playstyle = 1 << 0
	PlaystyleTablet       Playstyle = 1 << 1
	PlaystyleKeyboard     Playstyle = 1 << 2
	PlaystyleTouch        Playstyle = 1 << 3
)

const (
	DefaultPlayerGroupId    = 999
	DefaultSupporterGroupId = 1000
)

func (p Playstyle) Has(flag Playstyle) bool {
	return p&flag != 0
}

func (p Playstyle) String() string {
	if p == PlaystyleNotSpecified {
		return "None"
	}

	parts := make([]string, 0, 4)
	if p.Has(PlaystyleMouse) {
		parts = append(parts, "Mouse")
	}
	if p.Has(PlaystyleTablet) {
		parts = append(parts, "Tablet")
	}
	if p.Has(PlaystyleKeyboard) {
		parts = append(parts, "Keyboard")
	}
	if p.Has(PlaystyleTouch) {
		parts = append(parts, "Touch")
	}

	if len(parts) == 0 {
		return "Unknown"
	}

	return strings.Join(parts, ",")
}

var DisallowedUsernameSubstrings = []string{
	"blow job",
	"blowjob",
	"cockhead",
	"cocksucker",
	"cunt",
	"cunts",
	"dildo",
	"fag1t",
	"faget",
	"fagg1t",
	"faggit",
	"faggot",
	"fagit",
	"fags",
	"massterbait",
	"masstrbait",
	"masstrbate",
	"masterbaiter",
	"masterbate",
	"masterbates",
	"n1gr",
	"nigger",
	"nigur",
	"niiger",
	"niigr",
	"penis",
	"pussy",
	"slut",
	"whore",
	"nigga",
}
