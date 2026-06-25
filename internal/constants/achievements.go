package constants

type AchievementCategory struct {
	Name         string
	Achievements []string
}

var AchievementCategories = []AchievementCategory{
	{
		Name: "Beatmap Packs",
		Achievements: []string{
			"Anime Pack vol.1", "Anime Pack vol.2", "Anime Pack vol.3", "Anime Pack vol.4",
			"Internet! Pack vol.1", "Internet! Pack vol.2", "Internet! Pack vol.3", "Internet! Pack vol.4",
			"Rhythm Game Pack vol.1", "Rhythm Game Pack vol.2", "Rhythm Game Pack vol.3", "Rhythm Game Pack vol.4",
			"Video Game Pack vol.1", "Video Game Pack vol.2", "Video Game Pack vol.3", "Video Game Pack vol.4",
		},
	},
	{
		Name: "Skill",
		Achievements: []string{
			"500 Combo  (any song)", "750 Combo  (any song)", "1000 Combo  (any song)", "2000 Combo  (any song)",
			"I can see the top", "The gradual rise", "Scaling up", "Approaching the summit",
		},
	},
	{
		Name: "Dedication",
		Achievements: []string{
			"5,000 Plays (osu! mode)", "15,000 Plays (osu! mode)", "25,000 Plays (osu! mode)", "50,000 Plays (osu! mode)",
			"30,000 Drum Hits", "300,000 Drum Hits", "3,000,000 Drum Hits",
			"Catch 20,000 fruits", "Catch 200,000 fruits", "Catch 2,000,000 fruits",
			"40,000 Keys", "400,000 Keys", "4,000,000 Keys",
		},
	},
	{
		Name: "Hush-Hush",
		Achievements: []string{
			"A meganekko approaches", "Challenge Accepted", "Consolation Prize",
			"Don't let the bunny distract you!", "Jack of All Trades", "Jackpot", "Most Improved",
			"Non-stop Dancer", "Nonstop", "Obsessed", "Quick Draw", "S-Ranker", "Stumbler",
		},
	},
}
