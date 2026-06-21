package routes

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
	"gorm.io/gorm"
)

func UserProfileRedirect(ctx *server.Context) {
	query := strings.TrimSpace(ctx.PathValue("query"))
	if query == "" {
		NotFound(ctx)
		return
	}
	ctx.Redirect(http.StatusFound, fmt.Sprintf("/u/%s", query))
}

func UserProfile(ctx *server.Context) {
	query := strings.TrimSpace(ctx.PathValue("query"))
	if query == "" {
		NotFound(ctx)
		return
	}

	id, err := strconv.Atoi(query)
	if err != nil {
		resolveUserByName(ctx, query)
		return
	}

	user, err := ctx.State.Repositories.Users.ById(id, "Groups.Group", "Badges", "Names")
	if err != nil {
		// TODO: make use of LookupResult helper to reduce ErrRecordNotFound checks across the codebase
		if errors.Is(err, gorm.ErrRecordNotFound) {
			NotFound(ctx)
			return
		}
		ctx.Logger.Error("Failed to fetch user", "id", id, "error", err)
		InternalServerError(ctx)
		return
	}

	if !user.Activated {
		NotFound(ctx)
		return
	}

	online, err := ctx.State.Redis.Exists(
		ctx.Request.Context(),
		fmt.Sprintf("bancho:status:%d", user.Id),
	).Result()
	if err != nil {
		ctx.Logger.Error("Failed to fetch online status", "user", user.Id, "error", err)
	}

	followers, err := ctx.State.Repositories.Relationships.CountByTargetId(user.Id)
	if err != nil {
		ctx.Logger.Error("Failed to fetch follower count", "user", user.Id, "error", err)
		InternalServerError(ctx)
		return
	}

	totalPosts, err := ctx.State.Repositories.ForumPosts.CountByUserId(user.Id)
	if err != nil {
		ctx.Logger.Error("Failed to fetch post count", "user", user.Id, "error", err)
		InternalServerError(ctx)
		return
	}

	mode := resolveMode(ctx, user.PreferredMode)
	country := strings.ToUpper(user.Country)

	ppRank, err := ctx.State.Rankings.Rank(user.Id, mode, "performance", nil)
	if err != nil {
		ctx.Logger.Error("Failed to fetch global rank", "user", user.Id, "error", err)
	}
	ppRankCountry, err := ctx.State.Rankings.Rank(user.Id, mode, "performance", &country)
	if err != nil {
		ctx.Logger.Error("Failed to fetch country rank", "user", user.Id, "error", err)
	}

	currentAdded, targetAdded, isBlocked, err := resolveFriendStatus(ctx, user.Id)
	if err != nil {
		ctx.Logger.Error("Failed to resolve friend status", "user", user.Id, "error", err)
		InternalServerError(ctx)
		return
	}

	view := templates.UserProfileView{
		DefaultView:   buildDefaultViewWithPermissions(ctx),
		User:          user,
		Mode:          mode,
		IsOnline:      online == 1,
		Followers:     followers,
		TotalPosts:    totalPosts,
		PPRank:        ppRank,
		PPRankCountry: ppRankCountry,
		CurrentAdded:  currentAdded,
		TargetAdded:   targetAdded,
		IsBlocked:     isBlocked,
	}
	ctx.RenderTemplate(http.StatusOK, "pages/public/user", view)
}

func resolveFriendStatus(ctx *server.Context, targetId int) (currentAdded bool, targetAdded bool, isBlocked bool, err error) {
	if ctx.CurrentUser == nil || ctx.CurrentUser.Id == targetId {
		return false, false, false, nil
	}

	// The profile owner has blocked the current user
	if blocked, err := ctx.State.Repositories.Relationships.ByUserAndTarget(targetId, ctx.CurrentUser.Id); err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return false, false, false, err
		}
	} else if blocked.Status == 1 {
		return false, false, true, nil
	}

	currentFriends, err := ctx.State.Repositories.Relationships.TargetIdsByStatus(ctx.CurrentUser.Id, 0)
	if err != nil {
		return false, false, false, err
	}
	targetFriends, err := ctx.State.Repositories.Relationships.TargetIdsByStatus(targetId, 0)
	if err != nil {
		return false, false, false, err
	}

	currentAdded = slices.Contains(currentFriends, targetId)
	targetAdded = slices.Contains(targetFriends, ctx.CurrentUser.Id) ||
		slices.Contains(ctx.State.Config.SuperFriendlyUsers, targetId)
	return currentAdded, targetAdded, false, nil
}

func resolveUserByName(ctx *server.Context, query string) {
	// Try to find the user by their current name first
	if user, err := ctx.State.Repositories.Users.ByNameCaseInsensitive(query); err == nil {
		ctx.Redirect(http.StatusFound, fmt.Sprintf("/u/%d", user.Id))
		return
	}

	// Search the name history as a backup
	if name, err := ctx.State.Repositories.Names.ByName(query); err == nil {
		ctx.Redirect(http.StatusFound, fmt.Sprintf("/u/%d", name.UserId))
		return
	}

	// womp womp
	NotFound(ctx)
}
