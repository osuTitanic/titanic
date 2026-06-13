package routes

import (
	"fmt"
	"slices"
	"strings"

	"github.com/osuTitanic/titanic-go/internal/constants"
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
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

	country := query.Get("country")
	countryName := ""
	location := "All Locations"

	if country != "" {
		countryName = constants.GetCountryNameFromCode(country)
		location = countryName
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

	pagination := templates.NewPagination(templates.PaginationOptions{
		Path:        fmt.Sprintf("/rankings/%s/%s", mode.Alias(), rankingTypeString),
		Query:       query,
		CurrentPage: page,
		Total:       total,
		PageSize:    RankingsEntriesPerPage,
	})
	view := templates.RankingsView{
		DefaultView:  buildDefaultView(ctx),
		Pagination:   pagination,
		CountryName:  countryName,
		Country:      strings.ToLower(country),
		Location:     location,
		Type:         rankingType,
		Mode:         mode,
		TopCountries: topCountries,
		Entries:      entries,
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
