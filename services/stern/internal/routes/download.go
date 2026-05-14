package routes

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/osuTitanic/titanic-go/internal/schemas"
	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
)

func Download(ctx *server.Context) {
	releases, err := ctx.State.Repositories.ReleasesTitanic.FetchAll()
	if err != nil {
		ctx.Logger.Error("Failed to fetch download releases", "error", err)
		InternalServerError(ctx)
		return
	}

	selectedCategory := ctx.Request.URL.Query().Get("category")
	if selectedCategory == "" {
		selectedCategory = resolveDefaultCategory(releases)
	}

	view := &templates.DownloadView{
		DefaultView: buildDefaultView(ctx),
		Categories:  buildCategories(selectedCategory, releases),
		Clients:     buildReleases(selectedCategory, releases),
	}
	ctx.RenderTemplate(http.StatusOK, "pages/public/download", view)
}

func buildReleases(requestedCategory string, releases []*schemas.Release) []*schemas.Release {
	filteredReleases := make([]*schemas.Release, 0)

	for _, release := range releases {
		if !release.IsDisplayable() {
			continue
		}

		if release.Category == requestedCategory {
			filteredReleases = append(filteredReleases, release)
		}
	}

	return filteredReleases
}

func buildCategories(requestedCategory string, releases []*schemas.Release) []*templates.DownloadCategory {
	categories := make([]*templates.DownloadCategory, 0)
	seenCategories := make(map[string]bool)

	for _, release := range releases {
		if !release.IsDisplayable() {
			continue
		}
		if seenCategories[release.Category] {
			continue
		}

		categories = append(categories, &templates.DownloadCategory{
			Name:     release.Category,
			Selected: release.Category == requestedCategory,
			Url:      fmt.Sprintf("/download?category=%s", url.QueryEscape(release.Category)),
		})
		seenCategories[release.Category] = true
	}

	return categories
}

func resolveDefaultCategory(releases []*schemas.Release) string {
	for _, release := range releases {
		if !release.IsDisplayable() {
			continue
		}
		if release.Category == "" {
			continue
		}
		return release.Category
	}
	return ""
}
