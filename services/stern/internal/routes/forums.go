package routes

import (
	"net/http"

	"github.com/osuTitanic/titanic-go/internal/schemas"
	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
)

func ForumHome(ctx *server.Context) {
	mainForums, err := ctx.State.Forums.FetchMainForums()
	if err != nil {
		ctx.Logger.Error("Failed to fetch main forums", "error", err)
		InternalServerError(ctx)
		return
	}

	sections := make([]*templates.ForumSection, 0, len(mainForums))
	subForumIds := make([]int, 0)

	for _, mainForum := range mainForums {
		subForums, err := ctx.State.Forums.FetchSubForums(mainForum.Id)
		if err != nil {
			ctx.Logger.Error("Failed to fetch sub forums", "error", err, "forum", mainForum.Id)
			InternalServerError(ctx)
			return
		}

		sections = append(sections, &templates.ForumSection{
			Forum:     mainForum,
			Subforums: subForums,
		})

		for _, subForum := range subForums {
			subForumIds = append(subForumIds, subForum.Id)
		}
	}

	recent, err := ctx.State.ForumPosts.FetchLastForForums(
		subForumIds, "Topic", "User", "User.Groups.Group",
	)
	if err != nil {
		ctx.Logger.Error("Failed to fetch recent forum posts", "error", err)
		recent = map[int]*schemas.ForumPost{}
	}

	view := templates.ForumHomeView{
		DefaultView: buildDefaultView(ctx),
		Sections:    sections,
		Recent:      recent,
	}
	ctx.RenderTemplate(http.StatusOK, "pages/forum/home", view)
}
