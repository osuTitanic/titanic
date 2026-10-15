package routes

import (
	"net/http"

	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

// /web/osu-getseasonal.php -> Seasonal background images for osu! clients
func SeasonalBackgrounds(ctx *server.Context) {
	backgrounds := append([]string{}, ctx.State.Config.SeasonalBackgrounds...)
	ctx.RenderJson(http.StatusOK, backgrounds)
}
