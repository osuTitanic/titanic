package routes

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/osuTitanic/titanic-go/internal/constants"
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
)

const (
	activityPageSize     = 10
	activityRecentWindow = 30 * 24 * time.Hour
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
		ctx.Logger.Error("Failed to fetch user", "id", id, "error", err)
		InternalServerError(ctx)
		return
	}

	if user == nil || !user.Activated {
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

	currentAdded, targetAdded, isBlocked, err := resolveFriendStatus(ctx, user.Id)
	if err != nil {
		ctx.Logger.Error("Failed to resolve friend status", "user", user.Id, "error", err)
		InternalServerError(ctx)
		return
	}

	mode := resolveMode(ctx, user.PreferredMode)

	// The general tab is always pre-loaded, while the other
	// tabs are fetched on demand via their partials.
	general := buildUserGeneralTab(ctx, user, mode)

	view := templates.UserProfileView{
		DefaultView:   buildDefaultViewWithPermissions(ctx),
		User:          user,
		Mode:          mode,
		IsOnline:      online == 1,
		Followers:     followers,
		TotalPosts:    totalPosts,
		PPRank:        general.PPRank,
		PPRankCountry: general.PPRankCountry,
		CurrentAdded:  currentAdded,
		TargetAdded:   targetAdded,
		IsBlocked:     isBlocked,
		SuperFriendly: slices.Contains(ctx.State.Config.SuperFriendlyUsers, user.Id),
		General:       general,
	}
	ctx.RenderTemplate(http.StatusOK, "pages/public/user", view)
}

func UserGeneralPartial(ctx *server.Context) {
	user, ok := fetchProfileUser(ctx)
	if !ok {
		return
	}

	mode := resolveMode(ctx, user.PreferredMode)
	general := buildUserGeneralTab(ctx, user, mode)
	ctx.RenderTemplate(http.StatusOK, "partials/user_general", general)
}

func UserActivityPartial(ctx *server.Context) {
	user, ok := fetchProfileUser(ctx)
	if !ok {
		return
	}

	mode := resolveMode(ctx, user.PreferredMode)
	offset := 0
	if parsed, err := strconv.Atoi(ctx.Request.URL.Query().Get("offset")); err == nil && parsed > 0 {
		offset = parsed
	}

	page, err := buildActivityPage(ctx, user.Id, mode, offset)
	if err != nil {
		ctx.Logger.Error("Failed to fetch activity", "user", user.Id, "error", err)
		InternalServerError(ctx)
		return
	}
	ctx.RenderTemplate(http.StatusOK, "partials/user_activity", page)
}

func buildUserGeneralTab(ctx *server.Context, user *schemas.User, mode constants.Mode) *templates.UserGeneralTab {
	country := strings.ToUpper(user.Country)

	stats, err := ctx.State.Repositories.Stats.ByMode(user.Id, mode.Value())
	if err != nil {
		ctx.Logger.Error("Failed to fetch stats", "user", user.Id, "mode", mode, "error", err)
	}

	totalKudosu, err := ctx.State.Repositories.Modding.TotalKudosuByUser(user.Id)
	if err != nil {
		ctx.Logger.Error("Failed to fetch kudosu total", "user", user.Id, "error", err)
	}

	activity, err := buildActivityPage(ctx, user.Id, mode, 0)
	if err != nil {
		ctx.Logger.Error("Failed to fetch activity", "user", user.Id, "error", err)
		activity = &templates.UserActivityPage{UserId: user.Id, Mode: mode}
	}

	return &templates.UserGeneralTab{
		User:           user,
		Mode:           mode,
		Stats:          stats,
		PPRank:         fetchRankOrDefault(ctx, user.Id, mode, "performance", nil),
		PPRankCountry:  fetchRankOrDefault(ctx, user.Id, mode, "performance", &country),
		ScoreRank:      fetchRankOrDefault(ctx, user.Id, mode, "rscore", nil),
		TotalScoreRank: fetchRankOrDefault(ctx, user.Id, mode, "tscore", nil),
		PPv1Rank:       fetchRankOrDefault(ctx, user.Id, mode, "ppv1", nil),
		TotalKudosu:    totalKudosu,
		Activity:       activity,
	}
}

func buildActivityPage(ctx *server.Context, userId int, mode constants.Mode, offset int) (*templates.UserActivityPage, error) {
	rows, err := ctx.State.Repositories.Activities.FetchRecentByUser(
		userId, mode.Value(), activityPageSize+1, offset, activityRecentWindow,
	)
	if err != nil {
		return nil, err
	}

	hasMore := len(rows) > activityPageSize
	if hasMore {
		rows = rows[:activityPageSize]
	}

	return &templates.UserActivityPage{
		UserId:     userId,
		Mode:       mode,
		Rows:       rows,
		Offset:     offset,
		NextOffset: offset + activityPageSize,
		HasMore:    hasMore,
	}, nil
}

func fetchProfileUser(ctx *server.Context) (*schemas.User, bool) {
	id, err := strconv.Atoi(ctx.PathValue("id"))
	if err != nil {
		NotFound(ctx)
		return nil, false
	}

	user, err := ctx.State.Repositories.Users.ById(id)
	if err != nil {
		ctx.Logger.Error("Failed to fetch user", "id", id, "error", err)
		InternalServerError(ctx)
		return nil, false
	}

	if user == nil || !user.Activated {
		NotFound(ctx)
		return nil, false
	}
	return user, true
}

func fetchRankOrDefault(ctx *server.Context, userId int, mode constants.Mode, rankType string, country *string) int {
	rank, err := ctx.State.Rankings.Rank(userId, mode, rankType, country)
	if err != nil {
		ctx.Logger.Error("Failed to fetch rank", "user", userId, "type", rankType, "error", err)
		return 0
	}
	return rank
}

func resolveFriendStatus(ctx *server.Context, targetId int) (currentAdded bool, targetAdded bool, isBlocked bool, err error) {
	if ctx.CurrentUser == nil || ctx.CurrentUser.Id == targetId {
		return false, false, false, nil
	}

	// The profile owner has blocked the current user
	blocked, err := ctx.State.Repositories.Relationships.ByUserAndTarget(targetId, ctx.CurrentUser.Id)
	if err != nil {
		return false, false, false, err
	}
	if blocked != nil && blocked.Status == 1 {
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
	if user, err := ctx.State.Repositories.Users.ByNameCaseInsensitive(query); err == nil && user != nil {
		ctx.Redirect(http.StatusFound, fmt.Sprintf("/u/%d", user.Id))
		return
	}

	// Search the name history as a backup
	if name, err := ctx.State.Repositories.Names.ByName(query); err == nil && name != nil {
		ctx.Redirect(http.StatusFound, fmt.Sprintf("/u/%d", name.UserId))
		return
	}

	// womp womp
	NotFound(ctx)
}
