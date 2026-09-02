package routes

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

// defaultClientVersion is the latest osu! client version that supports TCP bancho connections
const defaultClientVersion = "b20130815"

// /web/bancho_connect.php -> Used to record statistics about bancho connections & more
func BanchoConnect(ctx *server.Context) {
	version := defaultClientVersion
	if ctx.Request.URL.Query().Has("v") {
		version = ctx.QueryValue("v")
	}
	// The client may also send retry, u, h, fx, and fail as query parameters
	// We don't really do anything with them though, since they primarily serve as analytical data

	match := constants.OsuVersion.FindStringSubmatch(version)
	if match == nil {
		ctx.RenderText(http.StatusOK, "xx")
		return
	}

	date, err := strconv.Atoi(match[1])
	if err != nil {
		ctx.RenderText(http.StatusOK, "xx")
		return
	}
	if date <= 20130915 {
		// TODO: forgot if this should be b20130815 or b20130915
		// Either way, TCP clients use this endpoint to resolve the host ip address of the bancho server
		// This was very quickly unused though, due to the switch to HTTP-based bancho connections
		ctx.RenderText(http.StatusOK, ctx.State.Config.BanchoIp)
		return
	}

	// It's possible to respond with "420" here to indicate that the server is busy
	// osu! will proceed to show: "Server is busy, please wait..."

	// It's also possible to respond with "error: verify" to initiate the client/2fa verification
	// osu! will open /p/verify + client hash as a parameter, which the user can approve in their browser

	// osu! can switch between bancho servers internally based on the returned country code
	// This feature seems to be unused though, with the only countries used being cn & tw
	ctx.RenderText(http.StatusOK, resolveBanchoCountry(ctx))
}

func resolveBanchoCountry(ctx *server.Context) string {
	if country := ctx.Request.Header.Get("CF-IPCountry"); country != "" {
		return strings.ToLower(country)
	}

	location, err := ctx.State.Location.Resolve(ctx.IP())
	if err != nil {
		ctx.Logger.Warn("Failed to resolve country for Bancho connection", "error", err)
		return "xx"
	}
	return strings.ToLower(location.CountryCode)
}
