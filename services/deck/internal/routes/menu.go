package routes

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/osuTitanic/titanic/internal/caching"
	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

var menuContentCache = caching.NewValue[[]byte](time.Hour)
var menuContentClient = &http.Client{Timeout: 8 * time.Second}

type MenuContentResponse struct {
	Images []MenuContentImage `json:"images"`
}

type MenuContentImage struct {
	Image     string     `json:"image"`
	URL       string     `json:"url"`
	IsCurrent bool       `json:"IsCurrent"`
	Begins    *time.Time `json:"begins"`
	Expires   *time.Time `json:"expires"`
}

// /menu-content.json -> Modern way of displaying menu content
func MenuContent(ctx *server.Context) {
	if ctx.State.Config.MenuIconImage == "" {
		ctx.RenderJson(http.StatusOK, MenuContentResponse{})
		return
	}

	response := MenuContentResponse{
		Images: []MenuContentImage{{
			Image:     ctx.State.Config.MenuIconImage,
			URL:       ctx.State.Config.MenuIconUrl,
			IsCurrent: true,
		}},
	}
	if err := ctx.RenderJson(http.StatusOK, response); err != nil {
		ctx.Logger.Error("Failed to render menu content", "error", err)
	}
}

// /web/osu-title-image.php -> Menu icon image & on-click redirect
func MenuIcon(ctx *server.Context) {
	// Clients additionally send "c" as an image checksum,
	// which we don't use at the moment

	if ctx.QueryValue("l") != "" {
		// Client (usually the browser) wants to be redirected to
		// the target link after clicking on the menu icon
		location := ctx.State.Config.MenuIconUrl
		if location == "" {
			location = ctx.State.Config.OsuBaseUrl()
		}
		ctx.Redirect(http.StatusTemporaryRedirect, location)
		return
	}

	// Client expects an image as response, so we cannot simply redirect to the image
	// -> Download & cache image, then return it to the client

	if ctx.State.Config.MenuIconImage == "" {
		// We don't have anything configured for the menu icon
		ctx.Response.WriteHeader(http.StatusOK)
		return
	}

	image, err := menuContentCache.GetOrLoad(func() ([]byte, error) {
		// We'll assume the configured url is trusted and doesn't blow up the cache with a huge image
		return fetchTitleImageContents(ctx)
	})
	if err != nil {
		ctx.Logger.Error("Failed to fetch title image", "error", err)
		ctx.Response.WriteHeader(http.StatusOK)
		return
	}

	ctx.Response.Header().Set("Content-Type", http.DetectContentType(image))
	ctx.Response.Header().Set("Cache-Control", "public, max-age=60")
	ctx.Response.WriteHeader(http.StatusOK)
	ctx.Response.Write(image)
}

func fetchTitleImageContents(ctx *server.Context) ([]byte, error) {
	request, err := http.NewRequestWithContext(
		ctx.Request.Context(),
		http.MethodGet,
		ctx.State.Config.MenuIconImage,
		nil,
	)
	if err != nil {
		return nil, err
	}
	request.Header.Set(
		"User-Agent",
		fmt.Sprintf("osuTitanic/deck (%s)", ctx.State.Config.DomainName),
	)

	response, err := menuContentClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	// TODO: We should defer http operations somewhere more central for reusability

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("unexpected response status: %s", response.Status)
	}
	return io.ReadAll(response.Body)
}
