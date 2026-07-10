package routes

import (
	"fmt"
	"net/url"
	"slices"
	"strings"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/services/stern/internal/server"
	"github.com/osuTitanic/titanic/services/stern/internal/templates"
)

const RankingsEntriesPerPage = 50

func RankingsGlobal(ctx *server.Context) {
	modeString := ctx.PathValue("mode")
	mode, ok := constants.NewModeFromAlias(modeString)
	if !ok {
		NotFound(ctx)
		return
	}

	rankingTypeString := strings.ToLower(ctx.PathValue("rankingType"))
	rankingTypeString = strings.TrimSpace(rankingTypeString)
	rankingType, ok := constants.NewRankingTypeFromAlias(rankingTypeString)
	if !ok {
		NotFound(ctx)
		return
	}

	query := ctx.Request.URL.Query()
	page, _ := parseInt(query.Get("page"))
	if page <= 0 {
		page = 1
	}

	country, countryName, ok := resolveRankingCountry(query.Get("country"))
	if !ok {
		NotFound(ctx)
		return
	}

	location := "All Locations"
	if country != "" {
		location = countryName
	}

	countryPointer := &country
	if country == "" {
		countryPointer = nil
	}

	jumpTo, page, err := resolveJumpTo(ctx, query, page, func(userId int) (int, error) {
		return ctx.State.Rankings.Rank(userId, mode, string(rankingType), countryPointer)
	})
	if err != nil {
		ctx.Logger.Error("Failed to resolve jumpto target", "error", err)
		InternalServerError(ctx)
		return
	}

	entries, total, err := resolveRankingEntries(ctx, mode, rankingType, country, page)
	if err != nil {
		ctx.Logger.Error("Failed to resolve ranking entries", "error", err)
		InternalServerError(ctx)
		return
	}

	topCountries, err := resolveTopCountries(ctx, mode, rankingType)
	if err != nil {
		ctx.Logger.Error("Failed to resolve top countries", "error", err)
		InternalServerError(ctx)
		return
	}
	totalBeatmaps := 0

	if rankingType == constants.RankingTypeClears {
		totalBeatmaps, err = resolveTotalBeatmaps(ctx, mode)
		if err != nil {
			ctx.Logger.Error("Failed to resolve total ranked beatmaps", "error", err)
			InternalServerError(ctx)
			return
		}
	}

	pagination := templates.NewPagination(templates.PaginationOptions{
		Path:        fmt.Sprintf("/rankings/%s/%s", mode.Alias(), rankingTypeString),
		Query:       filterQuery(query, "jumpto", "jumpto_id"),
		CurrentPage: page,
		Total:       total,
		PageSize:    RankingsEntriesPerPage,
	})
	view := templates.RankingsView{
		DefaultView:   buildDefaultView(ctx),
		Pagination:    pagination,
		CountryName:   countryName,
		Country:       strings.ToLower(country),
		Location:      location,
		Type:          rankingType,
		Mode:          mode,
		TopCountries:  topCountries,
		Entries:       entries,
		JumpTo:        jumpTo,
		TotalBeatmaps: totalBeatmaps,
	}
	ctx.RenderTemplate(200, "pages/public/rankings", view)
}

func resolveRankingEntries(ctx *server.Context, mode constants.Mode, rankingType constants.RankingType, country string, page int) ([]*templates.RankingEntry, int, error) {
	offset := max((page-1)*RankingsEntriesPerPage, 0)

	rankingTypeString := string(rankingType)
	countryPointer := &country
	if country == "" {
		countryPointer = nil
	}

	entries, err := ctx.State.Rankings.TopPlayers(
		mode, int64(offset), RankingsEntriesPerPage,
		rankingTypeString, countryPointer,
	)
	if err != nil {
		return nil, 0, err
	}

	totalPlayers, err := ctx.State.Rankings.PlayerCount(mode, rankingTypeString, countryPointer)
	if err != nil {
		return nil, 0, err
	}

	userIds := make([]int, len(entries))
	for i, entry := range entries {
		userIds[i] = entry.UserId
	}

	users, err := ctx.State.Repositories.Users.ManyById(userIds, "Stats")
	if err != nil {
		return nil, 0, err
	}

	userMapping := make(map[int]*schemas.User, len(users))
	for _, user := range users {
		userMapping[user.Id] = user
	}

	friendIds, err := resolveFriendIds(ctx)
	if err != nil {
		return nil, 0, err
	}

	rankingEntries := make([]*templates.RankingEntry, len(entries))
	for i, entry := range entries {
		rankingEntries[i] = &templates.RankingEntry{
			User:     userMapping[entry.UserId],
			Rank:     offset + i + 1,
			Score:    int(entry.Score),
			IsFriend: slices.Contains(friendIds, entry.UserId),
		}
		rankingEntries[i].User.SortStats()
	}

	return rankingEntries, int(totalPlayers), nil
}

func resolveTopCountries(ctx *server.Context, mode constants.Mode, rankingType constants.RankingType) ([]string, error) {
	countries, err := ctx.State.Rankings.TopCountriesForType(mode, string(rankingType))
	if err != nil {
		return nil, err
	}

	topCountries := make([]string, len(countries))
	for i, country := range countries {
		topCountries[i] = country.Name
	}
	return topCountries, nil
}

func resolveFriendIds(ctx *server.Context) ([]int, error) {
	if !ctx.IsAuthenticated() {
		return []int{}, nil
	}

	targetIds, err := ctx.State.Repositories.Relationships.TargetIdsByStatus(ctx.CurrentUser.Id, 0)
	if err != nil {
		return nil, err
	}

	return targetIds, nil
}

func resolveRankingCountry(value string) (string, string, bool) {
	country := strings.ToLower(strings.TrimSpace(value))
	if country == "" {
		return "", "", true
	}
	if country == "xx" {
		return "", "", false
	}

	countryName := constants.GetCountryNameFromCode(country)
	if countryName == "" {
		return "", "", false
	}
	return country, countryName, true
}

func resolveJumpTo(ctx *server.Context, query url.Values, page int, rankFor func(userId int) (int, error)) (string, int, error) {
	jumpTo := query.Get("jumpto")
	jumpToId := query.Get("jumpto_id")
	userId, _ := parseInt(jumpToId)

	// Attempt to resolve by username
	if jumpTo != "" {
		userIdByName, _ := ctx.State.Repositories.Users.GetUserIdCaseInsensitive(jumpTo)
		if userIdByName != 0 {
			userId = userIdByName
		}
	}

	// Attempt to resolve by user ID
	if userId != 0 {
		page, err := resolveJumpPage(rankFor, userId, page)
		if err != nil {
			return jumpTo, page, err
		}
		jumpTo, _ = ctx.State.Users.GetUsername(userId)
		return jumpTo, page, nil
	}

	return "", page, nil
}

func resolveJumpPage(rankFor func(userId int) (int, error), userId int, fallbackPage int) (int, error) {
	rank, err := rankFor(userId)
	if err != nil {
		return fallbackPage, err
	}
	if rank <= 0 {
		return fallbackPage, nil
	}
	return (rank-1)/RankingsEntriesPerPage + 1, nil
}

func resolveTotalBeatmaps(ctx *server.Context, mode constants.Mode) (int, error) {
	statuses := []constants.BeatmapStatus{
		constants.BeatmapStatusRanked,
		constants.BeatmapStatusApproved,
		constants.BeatmapStatusQualified,
		constants.BeatmapStatusLoved,
	}
	keys := make([]string, len(statuses))
	for i, status := range statuses {
		keys[i] = fmt.Sprintf("bancho:totalbeatmaps:%d:%d", mode.Value(), status.Value())
	}

	counts, err := ctx.State.Redis.MGet(ctx.Request.Context(), keys...).Result()
	if err != nil {
		return 0, err
	}

	total := 0
	for i, value := range counts {
		if value == nil {
			continue
		}

		valueString, ok := value.(string)
		if !ok {
			return 0, fmt.Errorf("unexpected value for %s", keys[i])
		}
		count, ok := parseInt(valueString)
		if !ok {
			return 0, fmt.Errorf("invalid value for %s: %q", keys[i], valueString)
		}
		total += count
	}
	return total, nil
}

func filterQuery(query url.Values, exclude ...string) url.Values {
	clone := make(url.Values, len(query))
	for key, values := range query {
		if slices.Contains(exclude, key) {
			continue
		}
		clone[key] = values
	}
	return clone
}
