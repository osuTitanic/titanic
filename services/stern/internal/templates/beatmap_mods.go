package templates

import "github.com/osuTitanic/titanic/internal/constants"

type BeatmapModOption struct {
	Acronym  string
	Label    string
	Selected bool
}

type BeatmapModGroup struct {
	Name string
	Mods []BeatmapModOption
}

func BuildBeatmapModGroups(mode constants.Mode, selected constants.Mods) []BeatmapModGroup {
	groups := []BeatmapModGroup{
		{
			Name: "Difficulty Reduction",
			Mods: []BeatmapModOption{
				beatmapModOption("EZ", constants.Easy, selected),
				beatmapModOption("NF", constants.NoFail, selected),
				beatmapModOption("HT", constants.HalfTime, selected),
			},
		},
		{
			Name: "Difficulty Increase",
			Mods: []BeatmapModOption{
				beatmapModOption("HR", constants.HardRock, selected),
				beatmapModOption("DT", constants.DoubleTime, selected),
				beatmapModOption("NC", constants.Nightcore, selected),
				beatmapModOption("HD", constants.Hidden, selected),
				beatmapModOption("FL", constants.Flashlight, selected),
			},
		},
		{
			Name: "Special",
			Mods: []BeatmapModOption{beatmapModOption("NM", constants.NoMod, selected)},
		},
	}

	switch mode {
	case constants.ModeOsu:
		groups[2].Mods = []BeatmapModOption{
			beatmapModOption("RX", constants.Relax, selected),
			beatmapModOption("AP", constants.Autopilot, selected),
			beatmapModOption("SO", constants.SpunOut, selected),
		}
	case constants.ModeTaiko, constants.ModeCatch:
		groups[2].Mods = []BeatmapModOption{
			beatmapModOption("RX", constants.Relax, selected),
		}
	case constants.ModeMania:
		groups[1].Mods = append(
			groups[1].Mods, beatmapModOptionWithLabel("FI", "FadeIn", constants.FadeIn, selected),
		)
		groups[2].Mods = []BeatmapModOption{
			beatmapModOptionWithLabel("RD", "Random", constants.Random, selected),
		}
	}

	return groups
}

func beatmapModOptionWithLabel(acronym, label string, mod constants.Mods, selected constants.Mods) BeatmapModOption {
	return BeatmapModOption{
		Acronym:  acronym,
		Label:    label,
		Selected: selected.Has(mod),
	}
}

func beatmapModOption(acronym string, mod constants.Mods, selected constants.Mods) BeatmapModOption {
	return beatmapModOptionWithLabel(acronym, acronym, mod, selected)
}
