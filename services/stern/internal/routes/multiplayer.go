package routes

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/services/stern/internal/server"
	"github.com/osuTitanic/titanic/services/stern/internal/templates"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var matchNumberPrinter = message.NewPrinter(language.English)

func Match(ctx *server.Context) {
	id := ctx.PathValue("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		NotFound(ctx)
		return
	}

	match, err := ctx.State.Repositories.Matches.ById(idInt, "Creator")
	if err != nil {
		ctx.Logger.Error("Failed to fetch match", "id", idInt, "error", err)
		InternalServerError(ctx)
		return
	}
	if match == nil || match.Creator == nil {
		NotFound(ctx)
		return
	}

	events, err := ctx.State.Repositories.Matches.EventsByMatchId(match.Id)
	if err != nil {
		ctx.Logger.Error("Failed to fetch match events", "id", match.Id, "error", err)
		InternalServerError(ctx)
		return
	}

	view := templates.MatchView{
		DefaultView: buildDefaultView(ctx),
		Match:       match,
		Events:      buildMatchEvents(events),
	}
	ctx.RenderTemplate(http.StatusOK, "pages/public/match", view)
}

func buildMatchEvents(events []*schemas.MatchEvent) []*templates.MatchEventView {
	views := make([]*templates.MatchEventView, 0, len(events))

	for _, event := range events {
		if view := buildMatchEvent(event); view != nil {
			views = append(views, view)
		}
	}
	return views
}

func buildMatchEvent(event *schemas.MatchEvent) *templates.MatchEventView {
	view := &templates.MatchEventView{
		Time: event.Time,
		Type: event.Type,
	}

	switch event.Type {
	case constants.MatchEventJoin, constants.MatchEventLeave, constants.MatchEventKick:
		data, err := event.UserData()
		if err != nil || data.Name == "" {
			return nil
		}
		view.User = &templates.MatchEventUser{Id: data.UserId, Name: data.Name}

	case constants.MatchEventHost:
		data, err := event.HostData()
		if err != nil || data.New.Name == "" {
			return nil
		}
		view.User = &templates.MatchEventUser{Id: data.New.Id, Name: data.New.Name}

	case constants.MatchEventDisband, constants.MatchEventAbort:
		// No event data required

	case constants.MatchEventResult:
		data, err := event.GameData()
		if err != nil {
			return nil
		}
		view.Game = buildMatchGame(data)

	default:
		// Start events are not displayed
		return nil
	}

	return view
}

func buildMatchGame(data *schemas.MatchEventGameData) *templates.MatchGameView {
	game := &templates.MatchGameView{
		BeatmapId:   data.BeatmapId,
		BeatmapText: data.BeatmapText,
		Mode:        data.Mode,
		TeamMode:    data.TeamMode,
		ScoringMode: data.ScoringMode,
		Duration:    formatMatchDuration(data.Duration()),
		Results:     make([]*templates.MatchGameResult, 0, len(data.Results)),
	}

	for _, result := range data.Results {
		game.Results = append(game.Results, buildMatchGameResult(result, data.Mods))
	}

	if game.TeamMode.HasTeams() {
		game.TeamResult = buildMatchTeamResult(game.ScoringMode, game.Results)
	}
	return game
}

func buildMatchGameResult(result schemas.MatchEventResultData, matchMods constants.Mods) *templates.MatchGameResult {
	return &templates.MatchGameResult{
		Place:     result.Place,
		UserId:    result.Player.Id,
		Username:  result.Player.Name,
		Country:   result.Player.Country,
		Team:      result.Player.Team,
		Mods:      matchMods | result.Player.Mods,
		Score:     result.Score.Total,
		Accuracy:  result.Score.Accuracy,
		MaxCombo:  result.Score.MaxCombo,
		Count300:  result.Score.Count300,
		Count100:  result.Score.Count100,
		Count50:   result.Score.Count50,
		CountMiss: result.Score.CountMiss,
		Failed:    result.Score.Failed,
	}
}

func buildMatchTeamResult(scoring constants.MatchScoringType, results []*templates.MatchGameResult) *templates.MatchTeamResult {
	blue := matchTeamValue(scoring, results, constants.SlotTeamBlue)
	red := matchTeamValue(scoring, results, constants.SlotTeamRed)

	winner := constants.SlotTeamNeutral
	if blue > red {
		winner = constants.SlotTeamBlue
	}
	if red > blue {
		winner = constants.SlotTeamRed
	}

	return &templates.MatchTeamResult{
		Blue:   formatMatchTeamValue(scoring, blue),
		Red:    formatMatchTeamValue(scoring, red),
		Winner: winner,
		Margin: formatMatchTeamMargin(scoring, math.Abs(blue-red)),
	}
}

func matchTeamValue(scoring constants.MatchScoringType, results []*templates.MatchGameResult, team constants.SlotTeam) float64 {
	var total float64
	var count int

	for _, result := range results {
		if result.Team != team {
			continue
		}

		switch scoring {
		case constants.MatchScoringAccuracy:
			total += result.Accuracy
		case constants.MatchScoringCombo:
			total += float64(result.MaxCombo)
		default:
			total += float64(result.Score)
		}
		count++
	}

	if scoring == constants.MatchScoringAccuracy && count > 0 {
		// For accuracy scoring, we want the average accuracy of the team, not the total
		total /= float64(count)
	}
	return total
}

func formatMatchTeamValue(scoring constants.MatchScoringType, value float64) string {
	if scoring == constants.MatchScoringAccuracy {
		return fmt.Sprintf("%.2f%%", value)
	}
	return matchNumberPrinter.Sprintf("%d", int64(value))
}

func formatMatchTeamMargin(scoring constants.MatchScoringType, difference float64) string {
	switch scoring {
	case constants.MatchScoringAccuracy:
		return fmt.Sprintf("%.2f%%", difference)
	case constants.MatchScoringCombo:
		return matchNumberPrinter.Sprintf("%d combo", int64(difference))
	default:
		return matchNumberPrinter.Sprintf("%d points", int64(difference))
	}
}

// formatMatchDuration formats the game duration, e.g. "3m 42s"
func formatMatchDuration(duration time.Duration) string {
	if duration <= 0 {
		return ""
	}
	return fmt.Sprintf("%dm %ds", int(duration.Minutes()), int(duration.Seconds())%60)
}
