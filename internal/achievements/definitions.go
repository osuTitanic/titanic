package achievements

import (
	"slices"
	"strconv"
	"strings"

	"github.com/osuTitanic/titanic/internal/constants"
)

// References:
// https://www.reddit.com/r/osugame/comments/4fnkgo/osu_achievementsmedals_thread/
// https://osu.ppy.sh/community/forums/topics/494188

type AchievementDefinition struct {
	Name     string
	Category string
	Filename string
	Check    func(*AchievementContext) bool
}

var Definitions = slices.Concat(
	definitionsBeatmapPacks,
	definitionsCombo,
	definitionsDedication,
	definitionsHushHush,
	definitionsRanking,
)

var definitionsBeatmapPacks = beatmapPackDefinitions()

var definitionsCombo = []AchievementDefinition{
	// Get a 500 combo on any map
	combo("500 Combo  (any song)", "combo500.png", 500),
	// Get a 750 combo on any map
	combo("750 Combo  (any song)", "combo750.png", 750),
	// Get a 1000 combo on any map
	combo("1000 Combo  (any song)", "combo1000.png", 1000),
	// Get a 2000 combo on any map
	combo("2000 Combo  (any song)", "combo2000.png", 2000),
}

var definitionsDedication = []AchievementDefinition{
	// Get a Play Count of 5,000 in osu!standard
	dedication("5,000 Plays (osu! mode)", "plays1.png", constants.ModeOsu, 5000, false),
	// Get a Play Count of 15,000 in osu!standard
	dedication("15,000 Plays (osu! mode)", "plays2.png", constants.ModeOsu, 15000, false),
	// Get a Play Count of 25,000 in osu!standard
	dedication("25,000 Plays (osu! mode)", "plays3.png", constants.ModeOsu, 25000, false),
	// Get a Play Count of 50,000 in osu!standard
	dedication("50,000 Plays (osu! mode)", "plays4.png", constants.ModeOsu, 50000, false),
	// Hit 30,000 notes in osu!taiko
	dedication("30,000 Drum Hits", "taiko1.png", constants.ModeTaiko, 30000, true),
	// Hit 300,000 notes in osu!taiko
	dedication("300,000 Drum Hits", "taiko2.png", constants.ModeTaiko, 300000, true),
	// Hit 3,000,000 notes in osu!taiko
	dedication("3,000,000 Drum Hits", "taiko3.png", constants.ModeTaiko, 3000000, true),
	// Catch 20,000 fruits in osu!catch
	dedication("Catch 20,000 fruits", "fruitsalad.png", constants.ModeCatch, 20000, true),
	// Catch 200,000 fruits in osu!catch
	dedication("Catch 200,000 fruits", "fruitplatter.png", constants.ModeCatch, 200000, true),
	// Catch 2,000,000 fruits in osu!catch
	dedication("Catch 2,000,000 fruits", "fruitod.png", constants.ModeCatch, 2000000, true),
	// Hit 40,000 keys in osu!mania
	dedication("40,000 Keys", "maniahits1.png", constants.ModeMania, 40000, true),
	// Hit 400,000 keys in osu!mania
	dedication("400,000 Keys", "maniahits2.png", constants.ModeMania, 400000, true),
	// Hit 4,000,000 keys in osu!mania
	dedication("4,000,000 Keys", "maniahits3.png", constants.ModeMania, 4000000, true),
}

var definitionsRanking = []AchievementDefinition{
	// Reach a profile rank of at least 500 in any osu! mode
	// NOTE: Used to be 50,000
	ranking("I can see the top", "high-ranker-1.png", 500),
	// Reach a profile rank of at least 100 in any osu! mode
	// NOTE: Used to be 10,000
	ranking("The gradual rise", "high-ranker-2.png", 100),
	// Reach a profile rank of at least 50 in any osu! mode
	// NOTE: Used to be 5,000
	ranking("Scaling up", "high-ranker-3.png", 50),
	// Reach a profile rank of at least 10 in any osu! mode
	// NOTE: Used to be 1,000
	ranking("Approaching the summit", "high-ranker-4.png", 10),
}

var definitionsHushHush = []AchievementDefinition{
	// Get a 371 out of 371 combo in the normal or a 447 out of 447 combo in the
	// hard difficulty of the beatmap "Chatmonchy - Make Up! Make Up!" by peppy
	definition("Don't let the bunny distract you!", "Hush-Hush", "bunny.png", func(context *AchievementContext) bool {
		return hasScoreAndBeatmap(context) &&
			strings.HasPrefix(context.Score.Beatmap.Filename, "Chatmonchy - Make Up! Make Up! (peppy)") &&
			context.Score.Perfect
	}),
	// Get an S rank (or higher) on 5 different beatmaps in a row
	definition("S-Ranker", "Hush-Hush", "s-ranker.png", func(context *AchievementContext) bool {
		if context == nil || len(context.RecentScores) != 5 {
			return false
		}
		beatmaps := make(map[int]struct{}, 5)

		for _, score := range context.RecentScores {
			if score == nil {
				return false
			}
			if score.Grade.Value() > constants.GradeS.Value() {
				return false
			}
			beatmaps[score.BeatmapId] = struct{}{} // hey golang dev team, please add sets
		}
		// Can't be on same beatmap
		return len(beatmaps) == 5
	}),
	// Set a D Rank then A rank (or higher), in the last day
	definition("Most Improved", "Hush-Hush", "improved.png", func(context *AchievementContext) bool {
		return hasScore(context) && context.Score.StatusPP == constants.ScoreStatusBest &&
			// Check if player has set a D Rank in the last 24 hours ...
			context.HadDOnMapLastDay &&
			// ... and the grade being A or higher
			context.Score.Grade.Value() <= constants.GradeA.Value()
	}),
	// Pass Yoko Ishida - paraparaMAX I without No Fail
	definition("Non-stop Dancer", "Hush-Hush", "dancer.png", func(context *AchievementContext) bool {
		return hasScoreAndBeatmap(context) &&
			context.Score.Beatmap.Filename == "Yoko Ishida - paraparaMAX I (chan) [marathon].osu" &&
			!context.Score.Mods.Has(constants.NoFail)
	}),
	// Pass any difficulty of any ranked mapset with below 75% accuracy without no-fail and/or easy mods
	definition("Consolation Prize", "Hush-Hush", "consolationprize.png", func(context *AchievementContext) bool {
		return hasScore(context) &&
			context.Score.Grade == constants.GradeD &&
			!context.Score.Mods.Has(constants.Easy) &&
			!context.Score.Mods.Has(constants.NoFail)
	}),
	// Complete an approved map
	definition("Challenge Accepted", "Hush-Hush", "challengeaccepted.png", func(context *AchievementContext) bool {
		return hasScoreAndBeatmap(context) && context.Score.Beatmap.Status == constants.BeatmapStatusApproved
	}),
	// Full Combo a map with less than 85% accuracy
	definition("Stumbler", "Hush-Hush", "stumbler.png", func(context *AchievementContext) bool {
		return hasScore(context) && context.Score.Perfect && context.Score.Acc <= 0.85
	}),
	// Complete a map with a score of at least 6 recurring numbers (i.e. 222,222 or 6,666,666)
	definition("Jackpot", "Hush-Hush", "jackpot.png", func(context *AchievementContext) bool {
		if !hasScore(context) {
			return false
		}
		digits := strconv.FormatInt(context.Score.TotalScore, 10)
		for _, digit := range digits {
			if strings.Count(digits, string(digit)) >= 6 {
				return true
			}
		}
		return false
	}),
	// Be the first person to pass a ranked or qualified map
	definition("Quick Draw", "Hush-Hush", "quickdraw.png", func(context *AchievementContext) bool {
		return hasScoreAndBeatmap(context) &&
			context.Score.Beatmap.AwardsScore() &&
			context.BeatmapLeaderboardSize <= 1
	}),
	// Play the same map over 100 times in a day, retries included
	definition("Obsessed", "Hush-Hush", "obsessed.png", func(context *AchievementContext) bool {
		return context != nil && context.SameMapPlaysLastDay >= 100
	}),
	// Get a Max Combo on a map with over 10 minutes of drain time
	definition("Nonstop", "Hush-Hush", "nonstop.png", func(context *AchievementContext) bool {
		return hasScoreAndBeatmap(context) &&
			context.Score.MaxCombo >= context.Score.Beatmap.MaxCombo &&
			context.Score.Beatmap.DrainLength >= 600
	}),
	// Reach a play count of at least 5,000 in all osu!standard, osu!taiko, osu!catch and osu!mania
	definition("Jack of All Trades", "Hush-Hush", "jack.png", func(context *AchievementContext) bool {
		if context == nil {
			return false
		}
		for _, mode := range constants.Modes {
			stats := context.StatsByMode[mode]
			if stats == nil || stats.Playcount < 5000 {
				return false
			}
		}
		return true
	}),
	// Meet Maria, the osu!mania mascot (finish an osu!mania map with at least 100 combo)
	definition("A meganekko approaches", "Hush-Hush", "meganekko.png", func(context *AchievementContext) bool {
		return hasScore(context) && context.Score.Mode == constants.ModeMania && context.Score.MaxCombo >= 100
	}),
}

func definition(name, category, filename string, check func(*AchievementContext) bool) AchievementDefinition {
	return AchievementDefinition{Name: name, Category: category, Filename: filename, Check: check}
}

func combo(name, filename string, required int) AchievementDefinition {
	return definition(name, "Skill", filename, func(context *AchievementContext) bool {
		return hasScore(context) && context.Score.MaxCombo >= required
	})
}

func ranking(name, filename string, maximum int) AchievementDefinition {
	return definition(name, "Skill", filename, func(context *AchievementContext) bool {
		return context != nil && context.GlobalRank > 0 && context.GlobalRank <= maximum
	})
}

func dedication(name, filename string, mode constants.Mode, required int, totalHits bool) AchievementDefinition {
	return definition(name, "Dedication", filename, func(context *AchievementContext) bool {
		if !hasScore(context) || context.Score.Mode != mode {
			return false
		}
		stats := context.StatsByMode[mode]
		if stats == nil {
			return false
		}
		if totalHits {
			return stats.TotalHits >= required
		}
		return stats.Playcount >= int64(required)
	})
}

func beatmapPackDefinition(pack BeatmapPack) AchievementDefinition {
	return definition(pack.name, "Beatmap Packs", pack.filename, func(context *AchievementContext) bool {
		if !hasScoreAndBeatmap(context) {
			return false
		}
		if !contains(pack.beatmapsetIds, context.Score.Beatmap.SetId) {
			// Score was not set inside this pack
			return false
		}
		for _, setId := range pack.beatmapsetIds {
			if _, completed := context.CompletedBeatmapsetIds[setId]; !completed {
				// User has not completed this beatmap
				return false
			}
		}
		return true
	})
}

func beatmapPackDefinitions() (result []AchievementDefinition) {
	for _, pack := range BeatmapPacks {
		result = append(result, beatmapPackDefinition(pack))
	}
	return result
}

func hasScore(context *AchievementContext) bool {
	return context != nil && context.Score != nil
}

func hasScoreAndBeatmap(context *AchievementContext) bool {
	return hasScore(context) && context.Score.Beatmap != nil
}

func contains(values []int, value int) bool {
	for _, candidate := range values {
		if candidate == value {
			return true
		}
	}
	return false
}
