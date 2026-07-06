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
	activityPageSize       = 10
	activityRecentWindow   = 30 * 24 * time.Hour
	topPlaysPageSize       = 5
	topPlaysPageSizeExpand = 15
	historyRecentLimit     = 5
	historyMostPlayedLimit = 15
	kudosuHistoryLimit     = 30
)

func UserProfileRedirect(ctx *server.Context) {
	query := strings.TrimSpace(ctx.PathValue("query"))
	if query == "" {
		UserNotFound(ctx)
		return
	}
	ctx.Redirect(http.StatusFound, fmt.Sprintf("/u/%s", query))
}

func UserProfile(ctx *server.Context) {
	query := strings.TrimSpace(ctx.PathValue("query"))
	if query == "" {
		UserNotFound(ctx)
		return
	}

	id, err := strconv.Atoi(query)
	if err != nil {
		resolveUserByName(ctx, query)
		return
	}

	user, err := ctx.State.Repositories.Users.ById(id, "Groups.Group", "Badges", "Stamps", "Names")
	if err != nil {
		ctx.Logger.Error("Failed to fetch user", "id", id, "error", err)
		InternalServerError(ctx)
		return
	}

	if user == nil || !user.Activated {
		UserNotFound(ctx)
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
	user, ok := fetchProfileUser(ctx, "Stamps")
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
	if parsed, err := ctx.QueryValueInt("offset"); err == nil && parsed > 0 {
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

func UserTopPlaysPartial(ctx *server.Context) {
	user, ok := fetchProfileUser(ctx)
	if !ok {
		return
	}

	mode := resolveMode(ctx, user.PreferredMode)
	isOwner := ctx.CurrentUser != nil && ctx.CurrentUser.Id == user.Id

	firstsRank, err := ctx.State.Rankings.LeaderScoresRank(user.Id, mode)
	if err != nil {
		ctx.Logger.Error("Failed to fetch firsts rank", "user", user.Id, "error", err)
	}

	pinned, err := buildScorePage(ctx, user.Id, mode, "pinned", 0, isOwner)
	if err != nil {
		ctx.Logger.Error("Failed to fetch pinned scores", "user", user.Id, "error", err)
		InternalServerError(ctx)
		return
	}

	best, err := buildScorePage(ctx, user.Id, mode, "best", 0, isOwner)
	if err != nil {
		ctx.Logger.Error("Failed to fetch best scores", "user", user.Id, "error", err)
		InternalServerError(ctx)
		return
	}

	first, err := buildScorePage(ctx, user.Id, mode, "first", 0, isOwner)
	if err != nil {
		ctx.Logger.Error("Failed to fetch first place scores", "user", user.Id, "error", err)
		InternalServerError(ctx)
		return
	}

	tab := &templates.UserTopPlaysTab{
		UserId:     user.Id,
		Mode:       mode,
		IsOwner:    isOwner,
		FirstsRank: firstsRank,
		Pinned:     pinned,
		Best:       best,
		First:      first,
	}
	ctx.RenderTemplate(http.StatusOK, "partials/user_leader", tab)
}

func UserScoresPartial(ctx *server.Context) {
	user, ok := fetchProfileUser(ctx)
	if !ok {
		return
	}

	section := ctx.QueryValue("section")
	if section != "pinned" && section != "best" && section != "first" {
		NotFound(ctx)
		return
	}
	mode := resolveMode(ctx, user.PreferredMode)
	isOwner := ctx.CurrentUser != nil && ctx.CurrentUser.Id == user.Id

	offset := 0
	if parsed, err := ctx.QueryValueInt("offset"); err == nil && parsed > 0 {
		offset = parsed
	}

	page, err := buildScorePage(ctx, user.Id, mode, section, offset, isOwner)
	if err != nil {
		ctx.Logger.Error("Failed to fetch scores", "user", user.Id, "section", section, "error", err)
		InternalServerError(ctx)
		return
	}
	ctx.RenderTemplate(http.StatusOK, "partials/user_scores", page)
}

func UserHistoryPartial(ctx *server.Context) {
	user, ok := fetchProfileUser(ctx)
	if !ok {
		return
	}
	mode := resolveMode(ctx, user.PreferredMode)

	mostPlayed, err := ctx.State.Repositories.Plays.FetchMostPlayedByUser(
		user.Id, historyMostPlayedLimit, 0, "Beatmap.Beatmapset",
	)
	if err != nil {
		ctx.Logger.Error("Failed to fetch most played", "user", user.Id, "error", err)
		InternalServerError(ctx)
		return
	}

	recent, err := ctx.State.Repositories.Scores.FetchRecentByUser(
		user.Id, mode, historyRecentLimit, constants.ScoreStatusFailed, "Beatmap.Beatmapset",
	)
	if err != nil {
		ctx.Logger.Error("Failed to fetch recent plays", "user", user.Id, "error", err)
		InternalServerError(ctx)
		return
	}

	tab := &templates.UserHistoryTab{
		UserId:     user.Id,
		Mode:       mode,
		MostPlayed: mostPlayed,
		Recent:     recent,
	}
	ctx.RenderTemplate(http.StatusOK, "partials/user_history", tab)
}

func UserBeatmapsPartial(ctx *server.Context) {
	user, ok := fetchProfileUser(ctx)
	if !ok {
		return
	}
	isOwner := ctx.CurrentUser != nil && ctx.CurrentUser.Id == user.Id

	favourites, err := ctx.State.Repositories.Favourites.ManyByUserId(user.Id, "Beatmapset")
	if err != nil {
		ctx.Logger.Error("Failed to fetch favourites", "user", user.Id, "error", err)
		InternalServerError(ctx)
		return
	}
	slices.SortFunc(favourites, func(a, b *schemas.BeatmapFavourite) int {
		return a.CreatedAt.Compare(b.CreatedAt)
	})

	created, err := ctx.State.Repositories.Beatmapsets.FetchByCreator(user.Id)
	if err != nil {
		ctx.Logger.Error("Failed to fetch created beatmapsets", "user", user.Id, "error", err)
		InternalServerError(ctx)
		return
	}

	collaborations, err := ctx.State.Repositories.Collaborations.FetchByUser(user.Id, "Beatmap.Beatmapset")
	if err != nil {
		ctx.Logger.Error("Failed to fetch collaborations", "user", user.Id, "error", err)
		InternalServerError(ctx)
		return
	}

	nominations, err := ctx.State.Repositories.Nominations.FetchByUser(user.Id, "Beatmapset")
	if err != nil {
		ctx.Logger.Error("Failed to fetch nominations", "user", user.Id, "error", err)
		InternalServerError(ctx)
		return
	}

	tab := &templates.UserBeatmapsTab{
		UserId:         user.Id,
		IsOwner:        isOwner,
		Favourites:     favourites,
		Nominations:    nominations,
		Collaborations: collaborations,
		Created:        buildCreatedBeatmapGroups(created, isOwner),
	}
	ctx.RenderTemplate(http.StatusOK, "partials/user_beatmaps", tab)
}

func UserKudosuPartial(ctx *server.Context) {
	user, ok := fetchProfileUser(ctx)
	if !ok {
		return
	}

	totalKudosu, err := ctx.State.Repositories.Modding.TotalKudosuByUser(user.Id)
	if err != nil {
		ctx.Logger.Error("Failed to fetch kudosu total", "user", user.Id, "error", err)
		InternalServerError(ctx)
		return
	}

	mods, err := ctx.State.Repositories.Modding.FetchRangeByUser(
		user.Id, kudosuHistoryLimit, 0, "Sender", "Target", "Post.Topic",
	)
	if err != nil {
		ctx.Logger.Error("Failed to fetch kudosu history", "user", user.Id, "error", err)
		InternalServerError(ctx)
		return
	}

	tab := &templates.UserKudosuTab{
		UserId:      user.Id,
		TotalKudosu: totalKudosu,
		Entries:     buildKudosuEntries(user.Id, mods),
	}
	ctx.RenderTemplate(http.StatusOK, "partials/user_kudosu", tab)
}

func UserAchievementsPartial(ctx *server.Context) {
	user, ok := fetchProfileUser(ctx)
	if !ok {
		return
	}

	unlocked, err := ctx.State.Repositories.Achievements.ManyByUserId(user.Id)
	if err != nil {
		ctx.Logger.Error("Failed to fetch achievements", "user", user.Id, "error", err)
		InternalServerError(ctx)
		return
	}

	tab := &templates.UserAchievementsTab{
		UserId:        user.Id,
		UnlockedCount: len(unlocked),
		Categories:    buildAchievementCategories(unlocked),
	}
	ctx.RenderTemplate(http.StatusOK, "partials/user_achievements", tab)
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

	ppRank, ppCountryRank, scoreRank, totalScoreRank, ppv1Rank := fetchUserRanks(
		ctx, user.Id,
		mode, country,
	)

	return &templates.UserGeneralTab{
		User:           user,
		Mode:           mode,
		Stats:          stats,
		PPRank:         ppRank,
		PPRankCountry:  ppCountryRank,
		ScoreRank:      scoreRank,
		TotalScoreRank: totalScoreRank,
		PPv1Rank:       ppv1Rank,
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

func buildScorePage(ctx *server.Context, userId int, mode constants.Mode, section string, offset int, isOwner bool) (*templates.UserScorePage, error) {
	approvedRewards := ctx.State.Config.ApprovedMapRewards
	preload := "Beatmap.Beatmapset"

	// Show fewer scores on the initial page load, then load a
	// larger amount with "Show me more".
	pageSize := topPlaysPageSize
	if offset > 0 {
		pageSize = topPlaysPageSizeExpand
	}

	var scores []*schemas.Score
	var total int
	var err error

	switch section {
	case "pinned":
		scores, err = ctx.State.Repositories.Scores.FetchPinned(userId, mode, pageSize+1, offset, preload)
		if err == nil {
			total, err = ctx.State.Repositories.Scores.FetchPinnedCount(userId, mode)
		}
	case "best":
		scores, err = ctx.State.Repositories.Scores.FetchBestRange(userId, mode, !approvedRewards, pageSize+1, offset, preload)
		if err == nil {
			total, err = ctx.State.Repositories.Scores.FetchBestCount(userId, mode, !approvedRewards)
		}
	case "first":
		scores, err = ctx.State.Repositories.Scores.FetchLeaderScores(userId, mode, pageSize+1, offset, preload)
		if err == nil {
			total, err = ctx.State.Repositories.Scores.FetchLeaderCount(userId, mode)
		}
	}

	if err != nil {
		return nil, err
	}

	hasMore := len(scores) > pageSize
	if hasMore {
		scores = scores[:pageSize]
	}

	return &templates.UserScorePage{
		UserId:          userId,
		Mode:            mode,
		Section:         section,
		Scores:          scores,
		Offset:          offset,
		NextOffset:      offset + pageSize,
		HasMore:         hasMore,
		Total:           total,
		IsOwner:         isOwner,
		ApprovedRewards: approvedRewards,
	}, nil
}

func buildCreatedBeatmapGroups(sets []*schemas.Beatmapset, isOwner bool) []*templates.UserBeatmapGroup {
	groups := []*templates.UserBeatmapGroup{
		{Name: "Ranked", Key: "ranked"},
		{Name: "Loved", Key: "loved"},
		{Name: "Qualified", Key: "qualified"},
		{Name: "Pending", Key: "pending", CanEdit: isOwner},
		{Name: "WIP", Key: "wip", CanEdit: isOwner},
		{Name: "Graveyarded", Key: "graveyarded", CanEdit: isOwner, Revivable: isOwner},
	}

	groupIndex := func(status constants.BeatmapStatus) int {
		switch status {
		case constants.BeatmapStatusRanked, constants.BeatmapStatusApproved:
			return 0
		case constants.BeatmapStatusLoved:
			return 1
		case constants.BeatmapStatusQualified:
			return 2
		case constants.BeatmapStatusPending:
			return 3
		case constants.BeatmapStatusWIP:
			return 4
		default:
			return 5 // Graveyard & anything else
		}
	}

	for _, set := range sets {
		if set.Status == constants.BeatmapStatusInactive {
			continue
		}
		index := groupIndex(set.Status)
		groups[index].Beatmapsets = append(groups[index].Beatmapsets, set)
	}

	result := make([]*templates.UserBeatmapGroup, 0, len(groups))
	for _, group := range groups {
		if len(group.Beatmapsets) > 0 {
			result = append(result, group)
		}
	}
	return result
}

func buildKudosuEntries(userId int, mods []*schemas.BeatmapModding) []*templates.UserKudosuEntry {
	entries := make([]*templates.UserKudosuEntry, 0, len(mods))

	for _, mod := range mods {
		status := "gave"
		switch {
		case mod.Amount < 0:
			status = "revoked"
		case mod.TargetId == userId:
			status = "received"
		}

		// On "received" the profile owner is the target, so the sender is the
		// "other" party. Otherwise the sender is the actor.
		actor, other := mod.Sender, mod.Target
		if status == "received" {
			actor, other = mod.Target, mod.Sender
		}

		preposition := "from"
		if status == "gave" {
			preposition = "to"
		}

		amount := mod.Amount
		if amount < 0 {
			amount = -amount
		}

		entry := &templates.UserKudosuEntry{
			Time:        mod.Time,
			Status:      status,
			Preposition: preposition,
			Amount:      amount,
			PostId:      mod.PostId,
		}
		if actor != nil {
			entry.ActorId, entry.ActorName = actor.Id, actor.Name
		}
		if other != nil {
			entry.OtherId, entry.OtherName = other.Id, other.Name
		}
		if mod.Post != nil && mod.Post.Topic != nil {
			entry.TopicTitle = mod.Post.Topic.Title
		}
		entries = append(entries, entry)
	}

	return entries
}

func buildAchievementCategories(unlocked []*schemas.Achievement) []*templates.UserAchievementCategory {
	byName := make(map[string]*schemas.Achievement, len(unlocked))
	for _, achievement := range unlocked {
		byName[achievement.Name] = achievement
	}
	categories := make([]*templates.UserAchievementCategory, 0, len(constants.AchievementCategories))

	for _, catalog := range constants.AchievementCategories {
		category := &templates.UserAchievementCategory{Name: catalog.Name}

		// For each available achievement in the catalog, check if the user has unlocked it
		for _, name := range catalog.Achievements {
			entry := &templates.UserAchievement{
				Name:     name,
				Unlocked: false,
			}
			if achievement, ok := byName[name]; ok {
				entry.Unlocked = true
				entry.Filename = achievement.Filename
				entry.UnlockedAt = achievement.UnlockedAt
			}
			category.Achievements = append(category.Achievements, entry)
		}
		categories = append(categories, category)
	}
	return categories
}

func fetchProfileUser(ctx *server.Context, preload ...string) (*schemas.User, bool) {
	id, err := ctx.PathValueInt("id")
	if err != nil {
		NotFound(ctx)
		return nil, false
	}

	user, err := ctx.State.Repositories.Users.ById(id, preload...)
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

// fetchUserRanks resolves all relevant ranks for the selected user & mode
func fetchUserRanks(ctx *server.Context, userId int, mode constants.Mode, country string) (int, int, int, int, int) {
	keys := []string{
		ctx.State.Rankings.RankingKey(mode, "performance", nil),
		ctx.State.Rankings.RankingKey(mode, "performance", &country),
		ctx.State.Rankings.RankingKey(mode, "rscore", nil),
		ctx.State.Rankings.RankingKey(mode, "tscore", nil),
		ctx.State.Rankings.RankingKey(mode, "ppv1", nil),
	}

	values, err := ctx.State.Rankings.RanksByKeys(keys, userId)
	if err != nil {
		ctx.Logger.Error("Failed to fetch ranks", "user", userId, "error", err)
		return 0, 0, 0, 0, 0
	}

	return values[0], values[1], values[2], values[3], values[4]
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
	if user, err := ctx.State.Repositories.Users.ByNameExtended(query); err == nil && user != nil {
		ctx.Redirect(http.StatusFound, fmt.Sprintf("/u/%d", user.Id))
		return
	}

	// Search the name history as a backup
	if name, err := ctx.State.Repositories.Names.ByNameExtended(query); err == nil && name != nil {
		ctx.Redirect(http.StatusFound, fmt.Sprintf("/u/%d", name.UserId))
		return
	}

	// womp womp
	UserNotFound(ctx)
}
