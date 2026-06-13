package routes

import (
	"slices"

	"github.com/osuTitanic/titanic-go/internal/schemas"
	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
)

func RankingsKudosu(ctx *server.Context) {
	query := ctx.Request.URL.Query()
	page, _ := parseInt(query.Get("page"))
	if page < 1 {
		page = 1
	}

	entries, total, err := resolveKudosuEntries(ctx, page)
	if err != nil {
		ctx.Logger.Error("Failed to resolve kudosu rankings", "error", err)
		InternalServerError(ctx)
		return
	}

	pagination := templates.NewPagination(templates.PaginationOptions{
		Path:        "/rankings/kudosu",
		Query:       query,
		CurrentPage: page,
		Total:       total,
		PageSize:    RankingsEntriesPerPage,
	})
	view := templates.KudosuView{
		DefaultView: buildDefaultView(ctx),
		Pagination:  pagination,
		Entries:     entries,
		JumpTo:      query.Get("jumpto"),
	}
	ctx.RenderTemplate(200, "pages/public/kudosu", view)
}

func resolveKudosuEntries(ctx *server.Context, page int) ([]*templates.KudosuEntry, int, error) {
	offset := max((page-1)*RankingsEntriesPerPage, 0)

	players, err := ctx.State.Rankings.TopKudosu(int64(offset), RankingsEntriesPerPage)
	if err != nil {
		return nil, 0, err
	}

	total, err := ctx.State.Rankings.PlayerCountKudosu()
	if err != nil {
		return nil, 0, err
	}

	userIds := make([]int, len(players))
	for i, player := range players {
		userIds[i] = player.UserId
	}

	users, err := ctx.State.Repositories.Users.ManyById(userIds)
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

	entries := make([]*templates.KudosuEntry, 0, len(players))
	for i, player := range players {
		user, ok := userMapping[player.UserId]
		if !ok {
			continue
		}
		entries = append(entries, &templates.KudosuEntry{
			User:     user,
			Kudosu:   int64(player.Score),
			Rank:     offset + i + 1,
			IsFriend: slices.Contains(friendIds, player.UserId),
		})
	}

	return entries, int(total), nil
}
