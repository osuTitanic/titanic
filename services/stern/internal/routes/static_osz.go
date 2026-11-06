package routes

import (
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/osuTitanic/titanic/internal/authentication"
	"github.com/osuTitanic/titanic/internal/ratelimiting"
	"github.com/osuTitanic/titanic/services/stern/internal/server"
)

var limitUnauthenticated = ratelimiting.Limit{Burst: 3, RefillInterval: 10 * time.Second, RefillRate: 1}
var limitAuthenticated = ratelimiting.Limit{Burst: 10, RefillInterval: 2 * time.Second, RefillRate: 1}

func BeatmapDownload(ctx *server.Context) {
	if !allowBeatmapDownload(ctx) {
		return
	}

	filename := ctx.PathValue("filename")

	// Handle filenames such as "1 Kenji Ninuma - DISCO PRINCE.osz"
	filename = strings.SplitN(filename, " ", 2)[0]
	noVideo := strings.Contains(filename, "n")

	setIdString := strings.TrimSuffix(filename, "n")
	setId, err := strconv.Atoi(setIdString)
	if err != nil {
		ctx.Response.WriteHeader(404)
		return
	}

	beatmapset, err := ctx.State.Beatmapsets.ById(setId)
	if err != nil || beatmapset == nil {
		ctx.Response.WriteHeader(404)
		return
	}
	if !beatmapset.Available {
		ctx.Response.WriteHeader(451)
		return
	}

	// noVideo can only be true if the beatmapset has videos
	noVideo = noVideo && beatmapset.HasVideo

	oszStream, oszSize, err := ctx.State.Resources.Osz(ctx.Context(), setId, noVideo)
	if err != nil {
		ctx.Response.WriteHeader(404)
		return
	}
	defer oszStream.Close()

	oszFilename := fmt.Sprintf("%d %s.osz", beatmapset.Id, beatmapset.Name())
	if noVideo {
		oszFilename = fmt.Sprintf("%d %s (no video).osz", beatmapset.Id, beatmapset.Name())
	}

	ctx.Response.Header().Set(
		"Content-Disposition",
		mime.FormatMediaType("attachment", map[string]string{"filename": oszFilename}),
	)
	ctx.Response.Header().Set("Content-Type", "application/octet-stream")
	ctx.Response.Header().Set("Last-Modified", beatmapset.LastUpdate.Format("Mon, 02 Jan 2006 15:04:05 GMT"))

	if oszSize > 0 {
		// Set content length if we can determine it, otherwise we'll use chunked transfer encoding
		ctx.Response.Header().Set("Content-Length", strconv.FormatInt(oszSize, 10))
	}

	ctx.Response.WriteHeader(200)
	io.Copy(ctx.Response, oszStream)
}

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
