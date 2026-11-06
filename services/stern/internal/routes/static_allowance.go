package routes

import (
	"encoding/hex"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/osuTitanic/titanic/internal/authentication"
	"github.com/osuTitanic/titanic/internal/ratelimiting"
	"github.com/osuTitanic/titanic/services/stern/internal/server"
)

var limitUnauthenticated = ratelimiting.Limit{Burst: 3, RefillInterval: 10 * time.Second, RefillRate: 1}
var limitAuthenticated = ratelimiting.Limit{Burst: 10, RefillInterval: 2 * time.Second, RefillRate: 1}

func allowBeatmapDownload(ctx *server.Context) bool {
	key, limit, err := resolveRatelimitConfiguration(ctx)
	if err != nil {
		ctx.Logger.Error("Failed to resolve ratelimit configuration", "error", err)
		ctx.Response.Header().Set("Cache-Control", "no-store")
		http.Error(ctx.Response, "Service Unavailable", http.StatusServiceUnavailable)
		return false
	}

	ctx.Logger.Debug(
		"Using ratelimit configuration",
		"key", key, "limit", limit,
	)

	limiter := ratelimiting.NewRedisLimiter(ctx.State.Redis) // TODO: Move this to ctx
	allowed, retryAfter, err := limiter.Allow(ctx.Context(), key, limit)
	if err != nil {
		ctx.Logger.Error("Failed to rate limit beatmap download", "error", err)
		ctx.Response.Header().Set("Cache-Control", "no-store")
		http.Error(ctx.Response, "Service Unavailable", http.StatusServiceUnavailable)
		return false
	}
	if !allowed {
		retryAfterSeconds := int64(math.Ceil(retryAfter.Seconds()))
		ctx.Response.Header().Set("Cache-Control", "no-store")
		ctx.Response.Header().Set("Retry-After", strconv.FormatInt(retryAfterSeconds, 10))
		http.Error(ctx.Response, "Too Many Requests", http.StatusTooManyRequests)
	}
	return allowed
}

func resolveRatelimitConfiguration(ctx *server.Context) (string, ratelimiting.Limit, error) {
	userId, err := resolveUser(ctx)
	if err != nil {
		ctx.Logger.Error("Failed to resolve beatmap download user", "error", err)
		http.Error(ctx.Response, "Service Unavailable", http.StatusServiceUnavailable)
		return "", limitUnauthenticated, err
	}

	if userId != 0 {
		return "ratelimits:osz:user:" + strconv.Itoa(userId), limitAuthenticated, nil
	}
	return "ratelimits:osz:ip:" + ctx.IP(), limitUnauthenticated, nil
}

func resolveUser(ctx *server.Context) (int, error) {
	if ctx.IsAuthenticated() {
		return ctx.CurrentUser.Id, nil
	}

	query := ctx.Request.URL.Query()
	username, password := query.Get("u"), query.Get("h") // auth query from osu! clients
	if username == "" || len(password) != 32 {
		return 0, nil
	}
	if _, err := hex.DecodeString(password); err != nil {
		return 0, nil
	}

	user, err := ctx.State.Users.ByName(username)
	if err != nil {
		return 0, err
	}
	if user == nil {
		return 0, nil
	}
	if !authentication.VerifyPasswordHashFromMd5(password, user.Bcrypt) {
		return 0, nil
	}
	return user.Id, nil
}
