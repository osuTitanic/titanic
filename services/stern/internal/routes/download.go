package routes

import (
	"net/http"

	"github.com/osuTitanic/titanic/services/stern/internal/server"
	"github.com/osuTitanic/titanic/services/stern/internal/templates"
)

func Download(ctx *server.Context) {
	updateUrl := ctx.State.Config.UpdateBaseUrl()
	downloadUrl := updateUrl + "/update?type=net40"
	downloadUrlFallback := updateUrl + "/update?type=net20"

	view := &templates.DownloadView{
		DefaultView:         buildDefaultView(ctx),
		DownloadUrl:         downloadUrl,
		DownloadUrlFallback: downloadUrlFallback,
	}
	ctx.RenderTemplate(http.StatusOK, "pages/public/download", view)
}
