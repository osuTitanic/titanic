package scoring

import (
	"fmt"
	"strings"

	"github.com/osuTitanic/titanic/internal/achievements"
	"github.com/osuTitanic/titanic/internal/activity"
	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
)

func (processor *Processor) unlockAchievements() error {
	submission := processor.submission
	if !submission.Passed {
		return nil
	}
	if submission.Id <= 0 {
		return nil
	}
	if submission.Relaxing() {
		return nil
	}

	context, err := achievements.LoadContext(
		processor.repositories,
		submission.Score,
		submission.CurrentStats,
	)
	if err != nil {
		return err
	}

	// Concurrently check for new achievment matches
	matched := context.Evaluate()

	newAchievements := make([]*schemas.Achievement, 0, len(matched))
	for _, definition := range matched {
		achievement := &schemas.Achievement{
			UserId:   submission.UserId,
			Name:     definition.Name,
			Category: definition.Category,
			Filename: definition.Filename,
		}
		created, err := processor.repositories.Achievements.CreateIfMissing(achievement)
		if err != nil {
			processor.AddWarning("unlock achievement %q: %w", definition.Name, err)
			continue
		}
		if created {
			newAchievements = append(newAchievements, achievement)
		}
		processor.context.Logger.Info(
			"Achievement unlocked",
			"user_id", submission.UserId, "name", definition.Name,
		)
	}

	if len(newAchievements) == 0 {
		return nil
	}
	submission.Achievements = newAchievements
	submission.Charts.Overall.Achievements = achievementFilenames(newAchievements)

	processor.publishAchievementActivities(newAchievements)
	processor.createAchievementNotification(newAchievements)
	return nil
}

func (processor *Processor) publishAchievementActivities(newAchievements []*schemas.Achievement) {
	submission := processor.submission

	for _, achievement := range newAchievements {
		err := activity.Submit(
			processor.context.State,
			submission.UserId,
			&submission.Mode,
			constants.ActivityAchievementUnlocked,
			map[string]any{
				"username":    submission.User.Name,
				"achievement": achievement.Name,
				"beatmap":     submission.Beatmap.Name(),
				"beatmap_id":  submission.BeatmapId,
			},
			true,  // Should be sent to #announce
			false, // Should be visible on profile page
		)
		if err != nil {
			processor.AddWarning("publish activity for achievement %q: %w", achievement.Name, err)
		}
	}
}

func (processor *Processor) createAchievementNotification(newAchievements []*schemas.Achievement) {
	header, content := achievementNotification(newAchievements)

	err := processor.repositories.Notifications.Create(&schemas.Notification{
		UserId:  processor.submission.UserId,
		Type:    constants.NotificationTypeAchievement,
		Header:  header,
		Content: content,
		Link: fmt.Sprintf(
			"%s/u/%d#achievements",
			processor.context.State.Config.OsuBaseUrl(),
			processor.submission.UserId,
		),
	})
	if err != nil {
		processor.AddWarning("create achievement notification: %w", err)
	}
}

func achievementNotification(newAchievements []*schemas.Achievement) (header, content string) {
	names := make([]string, 0, len(newAchievements))
	for _, achievement := range newAchievements {
		if achievement != nil {
			names = append(names, fmt.Sprintf("%q", achievement.Name))
		}
	}

	if len(names) == 1 {
		return "Achievement Unlocked!", fmt.Sprintf(
			"Congratulations for unlocking the %s achievement!",
			names[0],
		)
	}

	return "Achievements Unlocked!", fmt.Sprintf(
		"Congratulations for unlocking the %s and %s achievements!",
		strings.Join(names[:len(names)-1], ", "),
		names[len(names)-1],
	)
}

func achievementFilenames(newAchievements []*schemas.Achievement) []string {
	filenames := make([]string, 0, len(newAchievements))
	for _, achievement := range newAchievements {
		if achievement != nil {
			filenames = append(filenames, achievement.Filename)
		}
	}
	return filenames
}
