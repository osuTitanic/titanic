package scoring

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/osuTitanic/titanic/internal/constants"
)

type ChartWriter struct {
	strings.Builder
	fieldCount int
}

func (writer *ChartWriter) Field(name string, value any) {
	if writer.fieldCount > 0 {
		writer.WriteByte('|')
	}
	writer.WriteString(name)
	writer.WriteByte(':')
	fmt.Fprint(&writer.Builder, value)
	writer.fieldCount++
}

func (writer *ChartWriter) Change[T any](name string, change Change[T]) {
	writer.Field(name+"Before", optionalValue(change.Before))
	writer.Field(name+"After", optionalValue(change.After))
}

func (writer *ChartWriter) ChangeFormatted[T any](name string, change Change[T], format func(T) string) {
	before := ""
	if change.Before != nil {
		before = format(*change.Before)
	}
	after := ""
	if change.After != nil {
		after = format(*change.After)
	}
	writer.Field(name+"Before", before)
	writer.Field(name+"After", after)
}

func FormatSubmissionResponse(endpoint Endpoint, result *SubmissionContext) string {
	if endpoint.UsesLegacyResponse() {
		if !result.Passed {
			return ""
		}
		return FormatSubmissionResponseLegacy(result)
	}
	return FormatSubmissionResponseModular(result, endpoint)
}

func FormatSubmissionResponseModular(result *SubmissionContext, endpoint Endpoint) string {
	charts := []string{formatBeatmapInfoChart(result.Charts.Beatmap)}

	if endpoint.IncludesBeatmapRankingChart() {
		// /web/osu-submit-modular-selector.php introduced per-beatmap
		// statistics in addition to the overall user statistics
		charts = append(charts, formatBeatmapRankingChart(result.Charts.Ranking))
	}

	charts = append(charts, formatOverallChart(endpoint, result.Charts.Overall))
	return strings.Join(charts, "\n")
}

func FormatSubmissionResponseLegacy(result *SubmissionContext) string {
	beatmapRank := 0
	if result.StatusScore == constants.ScoreStatusBest {
		beatmapRank = result.NewBeatmapRank
	}

	lines := []string{
		// "Congratulations - you achieved online rank #<rank>!" or
		// "You have previously set a higher online record."
		strconv.Itoa(beatmapRank),
		// "You need another <x> points to reach the next rank!"
		strconv.FormatInt(result.Charts.Overall.ToNextRank, 10),
	}

	// Achievements are loaded from `http://osu.<domain>/images/achievements/<filename>`
	achievements := make([]string, 0, len(result.Achievements))
	for _, achievement := range result.Achievements {
		if achievement != nil {
			achievements = append(achievements, achievement.Filename)
		}
	}
	if len(achievements) > 0 {
		lines = append(lines, strings.Join(achievements, " "))
	}

	return strings.Join(lines, "\n")
}

func formatBeatmapInfoChart(chart BeatmapInfoChart) string {
	var response ChartWriter
	response.Field("beatmapId", chart.BeatmapId)
	response.Field("beatmapSetId", chart.BeatmapSetId)
	response.Field("beatmapPlaycount", chart.BeatmapPlaycount)
	response.Field("beatmapPasscount", chart.BeatmapPasscount)
	response.Field("approvedDate", formatApprovedDate(chart.ApprovedDate))
	return response.String()
}

func formatOverallChart(endpoint Endpoint, chart OverallChart) string {
	var response ChartWriter
	response.Field("chartId", chart.ChartId)
	response.Field("chartName", chart.ChartName)
	response.Field("chartUrl", chart.ChartUrl)
	response.Field("chartEndDate", chart.ChartEndDate)
	response.Field("achievements", strings.Join(chart.Achievements, " "))
	response.Field("achievements-new", strings.Join(chart.AchievementsNew, " "))
	response.Change("rank", chart.Rank)
	response.Change("rankedScore", chart.RankedScore)
	response.Change("totalScore", chart.TotalScore)
	response.Change("playCount", chart.PlayCount)
	response.Change("maxCombo", chart.MaxCombo)
	response.ChangeFormatted("pp", chart.PP, formatPP)

	accuracyFormatter := formatAccuracy
	if endpoint == EndpointModular {
		accuracyFormatter = formatAccuracyNormalized
	}
	response.ChangeFormatted("accuracy", chart.Accuracy, accuracyFormatter)

	response.Field("onlineScoreId", chart.OnlineScoreId)
	response.Field("toNextRankUser", chart.ToNextRankUser)
	response.Field("toNextRank", chart.ToNextRank)
	response.Change("beatmapRanking", chart.BeatmapRanking)
	return response.String()
}

func formatBeatmapRankingChart(chart BeatmapRankingChart) string {
	var response ChartWriter
	response.Field("chartId", chart.ChartId)
	response.Field("chartName", chart.ChartName)
	response.Field("chartUrl", chart.ChartUrl)
	response.Change("rank", chart.Rank)
	response.Change("rankedScore", chart.RankedScore)
	response.Change("totalScore", chart.TotalScore)
	response.Change("maxCombo", chart.MaxCombo)
	response.ChangeFormatted("accuracy", chart.Accuracy, formatAccuracy)
	response.ChangeFormatted("pp", chart.PP, formatPP)
	response.Field("toNextRank", chart.ToNextRank)
	response.Field("toNextRankUser", chart.ToNextRankUser)
	return response.String()
}

func optionalValue[T any](value *T) any {
	if value == nil {
		return ""
	}
	return *value
}

func formatPP(value float64) string {
	return strconv.FormatFloat(math.RoundToEven(value), 'f', -1, 64)
}

func formatAccuracy(value float64) string {
	value = math.Round(value*100) / 100
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func formatAccuracyNormalized(value float64) string {
	normalized := value / 100
	normalized = math.Round(normalized*10000) / 10000
	return strconv.FormatFloat(normalized, 'f', -1, 64)
}

func formatApprovedDate(value *time.Time) string {
	if value == nil {
		return ""
	}
	return value.Format("2006-01-02 15:04:05.999999-07:00")
}
