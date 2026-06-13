package routes

import (
	"fmt"

	"github.com/osuTitanic/titanic-go/internal/constants"
	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
)

func RankingsCountry(ctx *server.Context) {
	modeString := ctx.PathValue("mode")
	mode, ok := constants.NewModeFromAlias(modeString)
	if !ok {
		NotFound(ctx)
		return
	}

	query := ctx.Request.URL.Query()
	page, _ := parseInt(query.Get("page"))
	if page < 1 {
		page = 1
	}

	countries, err := ctx.State.Rankings.TopCountries(mode)
	if err != nil {
		ctx.Logger.Error("Failed to resolve country rankings", "error", err)
		InternalServerError(ctx)
		return
	}

	pagination := templates.NewPagination(templates.PaginationOptions{
		Path:        fmt.Sprintf("/rankings/%s/country", mode.Alias()),
		Query:       query,
		CurrentPage: page,
		Total:       len(countries),
		PageSize:    RankingsEntriesPerPage,
	})

	// Cap countries to pagination window
	start := min((page-1)*RankingsEntriesPerPage, len(countries))
	end := min(start+RankingsEntriesPerPage, len(countries))
	countries = countries[start:end]

	view := templates.CountryRankingsView{
		DefaultView: buildDefaultView(ctx),
		Pagination:  pagination,
		Type:        constants.RankingTypeCountry,
		Mode:        mode,
		Entries:     countries,
	}
	ctx.RenderTemplate(200, "pages/public/rankings", view)
}
