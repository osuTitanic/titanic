package routes

import (
	"strconv"

	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
)

func BeatmapPacks(ctx *server.Context) {
	categories, err := ctx.State.Repositories.BeatmapPacks.FetchCategories()
	if err != nil {
		InternalServerError(ctx)
		return
	}

	category := ctx.Request.URL.Query().Get("category")
	if category == "" && len(categories) > 0 {
		category = categories[0]
	}

	packs, err := ctx.State.Repositories.BeatmapPacks.FetchByCategory(category, "Creator")
	if err != nil {
		InternalServerError(ctx)
		return
	}

	view := &templates.BeatmapPacksView{
		DefaultView:      buildDefaultView(ctx),
		Categories:       categories,
		CategorySelected: category,
		BeatmapPacks:     packs,
	}
	ctx.RenderTemplate(200, "pages/public/packs", view)
}

func BeatmapPackInfo(ctx *server.Context) {
	id := ctx.PathValue("id")
	if id == "" {
		NotFound(ctx)
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		NotFound(ctx)
		return
	}

	pack, err := ctx.State.Repositories.BeatmapPacks.FetchById(idInt, "Creator", "Entries", "Entries.Beatmapset")
	if err != nil {
		InternalServerError(ctx)
		return
	}

	view := any(pack)
	ctx.RenderTemplate(200, "partials/pack_info", view)
}
