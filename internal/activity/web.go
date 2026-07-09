package activity

import (
	"fmt"
	"html"
	"strings"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
)

// htmlFormatter renders an activity entry into HTML.
type htmlFormatter func(entry *schemas.Activity, data map[string]any) string

var htmlFormatters = map[constants.UserActivity]htmlFormatter{
	constants.ActivityRanksGained:            formatRanksGained,
	constants.ActivityNumberOne:              formatNumberOne,
	constants.ActivityBeatmapLeaderboardRank: formatLeaderboardRank,
	constants.ActivityLostFirstPlace:         formatLostFirstPlace,
	constants.ActivityPPRecord:               formatPPRecord,
	constants.ActivityTopPlay:                formatTopPlay,
	constants.ActivityAchievementUnlocked:    formatAchievement,
	constants.ActivityBeatmapUploaded:        formatBeatmapUpload,
	constants.ActivityBeatmapUpdated:         formatBeatmapUpdate,
	constants.ActivityBeatmapRevived:         formatBeatmapRevival,
	constants.ActivityBeatmapStatusUpdated:   formatBeatmapStatusUpdate,
	constants.ActivityBeatmapNominated:       formatBeatmapNomination,
	constants.ActivityBeatmapNuked:           formatBeatmapNuke,
	constants.ActivityForumTopicCreated:      formatTopicCreated,
	constants.ActivityForumPostCreated:       formatPostCreated,
}

// RenderHtml renders an activity entry to a HTML string.
func RenderHtml(entry *schemas.Activity) string {
	data, ok := parseData(entry)
	if !ok {
		return ""
	}

	formatter, ok := htmlFormatters[constants.UserActivity(entry.Type)]
	if !ok {
		return ""
	}
	return formatter(entry, data)
}

func formatRanksGained(entry *schemas.Activity, data map[string]any) string {
	plural := "s"
	if dataString(data, "ranks_gained") == "1" {
		plural = ""
	}

	return fmt.Sprintf(
		"%s has risen %s rank%s, now placed #%s overall in %s.",
		userLink(entry, data), dataString(data, "ranks_gained"), plural,
		dataString(data, "rank"), htmlText(data, "mode"),
	)
}

func formatNumberOne(entry *schemas.Activity, data map[string]any) string {
	return fmt.Sprintf(
		"%s has taken the lead as the top-ranked %s player.",
		userLink(entry, data), htmlText(data, "mode"),
	)
}

func formatLeaderboardRank(entry *schemas.Activity, data map[string]any) string {
	builder := strings.Builder{}

	// <user> achieved rank #<rank> on <beatmap>
	fmt.Fprintf(
		&builder, "%s achieved rank #%s on %s ",
		userLink(entry, data), dataString(data, "beatmap_rank"), beatmapLink(data),
	)

	// with <mods>
	if mods := dataString(data, "mods"); mods != "" {
		fmt.Fprintf(&builder, "with %s ", html.EscapeString(mods))
	}

	// <mode>
	fmt.Fprintf(&builder, "<%s>", htmlText(data, "mode"))

	// (pp)
	if pp := dataString(data, "pp"); pp != "" {
		fmt.Fprintf(&builder, " (%spp)", pp)
	}
	return builder.String()
}

func formatLostFirstPlace(entry *schemas.Activity, data map[string]any) string {
	return fmt.Sprintf(
		"%s has lost first place on %s <%s>",
		userLink(entry, data), beatmapLink(data), htmlText(data, "mode"),
	)
}

func formatPPRecord(entry *schemas.Activity, data map[string]any) string {
	return fmt.Sprintf(
		"%s has set the new pp record on %s with %spp <%s>",
		userLink(entry, data), beatmapLink(data), dataString(data, "pp"), htmlText(data, "mode"),
	)
}

func formatTopPlay(entry *schemas.Activity, data map[string]any) string {
	return fmt.Sprintf(
		"%s got a new top play on %s with %spp <%s>",
		userLink(entry, data), beatmapLink(data), dataString(data, "pp"), htmlText(data, "mode"),
	)
}

func formatAchievement(entry *schemas.Activity, data map[string]any) string {
	return fmt.Sprintf(
		"%s unlocked an achievement: %s",
		userLink(entry, data), htmlText(data, "achievement"),
	)
}

func formatBeatmapUpload(entry *schemas.Activity, data map[string]any) string {
	return fmt.Sprintf("%s uploaded a new beatmapset \"%s\"", userLink(entry, data), beatmapsetLink(data))
}

func formatBeatmapUpdate(entry *schemas.Activity, data map[string]any) string {
	return fmt.Sprintf("%s has updated the beatmap \"%s\"", userLink(entry, data), beatmapsetLink(data))
}

func formatBeatmapRevival(entry *schemas.Activity, data map[string]any) string {
	return fmt.Sprintf("%s has been revived from eternal slumber.", beatmapsetLink(data))
}

func formatBeatmapStatusUpdate(entry *schemas.Activity, data map[string]any) string {
	status := constants.BeatmapStatus(dataInt(data, "status")).String()

	return fmt.Sprintf(
		"%s was set to \"%s\" by %s",
		beatmapsetLink(data), html.EscapeString(status), userLink(entry, data),
	)
}

func formatBeatmapNomination(entry *schemas.Activity, data map[string]any) string {
	if dataString(data, "type") != "reset" {
		return fmt.Sprintf("%s was nominated by %s", beatmapsetLink(data), userLink(entry, data))
	}
	return fmt.Sprintf("%s popped the bubble for \"%s\"", userLink(entry, data), beatmapsetLink(data))
}

func formatBeatmapNuke(entry *schemas.Activity, data map[string]any) string {
	return fmt.Sprintf("%s has nuked \"%s\"", userLink(entry, data), beatmapsetLink(data))
}

func formatTopicCreated(entry *schemas.Activity, data map[string]any) string {
	topic := htmlLink(
		fmt.Sprintf("/forum/t/%s", dataString(data, "topic_id")),
		dataString(data, "topic_name"),
	)
	return fmt.Sprintf("%s created a new topic \"%s\"", userLink(entry, data), topic)
}

func formatPostCreated(entry *schemas.Activity, data map[string]any) string {
	post := htmlLink(
		fmt.Sprintf("/forum/t/%s/p/%s", dataString(data, "topic_id"), dataString(data, "post_id")),
		dataString(data, "topic_name"),
	)
	return fmt.Sprintf("%s created a post in \"%s\"", userLink(entry, data), post)
}

func userLink(entry *schemas.Activity, data map[string]any) string {
	return htmlLink(fmt.Sprintf("/u/%d", entry.UserId), dataString(data, "username"))
}

func beatmapLink(data map[string]any) string {
	return htmlLink(fmt.Sprintf("/b/%s", dataString(data, "beatmap_id")), dataString(data, "beatmap"))
}

func beatmapsetLink(data map[string]any) string {
	return htmlLink(fmt.Sprintf("/s/%s", dataString(data, "beatmapset_id")), dataString(data, "beatmapset_name"))
}
