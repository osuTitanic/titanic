package routes

import (
	"net/http"

	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

// /web/check-updates.php -> Latest update endpoint for osu! clients (after late-2014)
func CheckUpdates(ctx *server.Context) {
	// We don't really do much with this endpoint as we
	// don't want clients to auto-update themselves
	ctx.Response.Header().Set("Content-Type", "application/json; charset=utf-8")
	ctx.Response.WriteHeader(http.StatusOK)
	ctx.Response.Write([]byte("[]"))
}
