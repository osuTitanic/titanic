package routes

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/osuTitanic/titanic/internal/media"
	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

// NOTE: These are static endpoints, also included in the stern service
//       Included as completeness for deck as a standalone service

var allowedAvatarSizes = map[int]struct{}{
	25:  {},
	128: {},
	256: {},
}

const defaultAvatarSize = 128
const avatarCacheTTL = time.Hour * 24

func BeatmapDownload(ctx *server.Context) {
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

	oszStream, oszSize, err := ctx.State.Resources.Osz(setId, noVideo)
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

func BeatmapThumbnail(ctx *server.Context) {
	filename := ctx.PathValue("filename")

	// Handle filenames such as "1.jpg" (small) and "1l.jpg" (large)
	key, _, _ := strings.Cut(filename, ".")
	large := strings.Contains(key, "l")

	setId, err := strconv.Atoi(strings.ReplaceAll(key, "l", ""))
	if err != nil {
		ctx.Response.WriteHeader(404)
		return
	}

	stream, err := ctx.State.Resources.Background(setId, large)
	if err != nil {
		ctx.Response.WriteHeader(404)
		return
	}
	defer stream.Close()

	// If a cache key is provided, the thumbnail may be cached by the client
	if ctx.QueryValue("c") != "" {
		ctx.Response.Header().Set("Cache-Control", "public, max-age=604800, immutable")
	}

	ctx.Response.Header().Set("Content-Type", "image/jpeg")
	ctx.Response.WriteHeader(200)
	io.Copy(ctx.Response, stream)
}

func BeatmapAudioPreview(ctx *server.Context) {
	filename := ctx.PathValue("filename")

	// Handle filenames such as "1.mp3"
	key, _, _ := strings.Cut(filename, ".")

	setId, err := strconv.Atoi(key)
	if err != nil {
		ctx.Response.WriteHeader(404)
		return
	}

	stream, err := ctx.State.Resources.Preview(setId)
	if err != nil {
		ctx.Response.WriteHeader(404)
		return
	}
	defer stream.Close()

	// If a cache key is provided, the preview may be cached by the client
	if ctx.QueryValue("c") != "" {
		ctx.Response.Header().Set("Cache-Control", "public, max-age=604800, immutable")
	}

	ctx.Response.Header().Set("Content-Type", "audio/mpeg")
	ctx.Response.WriteHeader(200)
	io.Copy(ctx.Response, stream)
}

func ScreenshotImage(ctx *server.Context) {
	id, err := ctx.PathValueInt("id")
	if err != nil {
		ctx.Response.WriteHeader(http.StatusNotFound)
		return
	}

	screenshot, err := ctx.State.Screenshots.ById(id)
	if err != nil || screenshot.Hidden {
		ctx.Response.WriteHeader(http.StatusNotFound)
		return
	}

	if ctx.PathValue("checksum") != screenshot.Checksum() {
		ctx.Response.WriteHeader(http.StatusNotFound)
		return
	}

	image, err := ctx.State.Storage.Read(strconv.Itoa(id), "screenshots")
	if err != nil {
		ctx.Response.WriteHeader(http.StatusNotFound)
		return
	}

	contentType := http.DetectContentType(image)
	if !strings.HasPrefix(contentType, "image/") {
		ctx.Response.WriteHeader(http.StatusNotFound)
		return
	}

	fileExtension := strings.TrimPrefix(contentType, "image/")
	filename := fmt.Sprintf(
		"ss (%s).%s", // inspired by puush
		screenshot.CreatedAt.Format("2006-01-02 at 15.04.05"), fileExtension,
	)

	ctx.Response.Header().Set("Content-Type", contentType)
	ctx.Response.Header().Set("Cache-Control", "public, max-age=1209600, immutable")
	ctx.Response.Header().Set("Date", screenshot.CreatedAt.Format("Mon, 02 Jan 2006 15:04:05 GMT"))
	ctx.Response.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filename))
	ctx.Response.Header().Set("Content-Length", strconv.Itoa(len(image)))
	ctx.Response.WriteHeader(http.StatusOK)
	ctx.Response.Write(image)
}

func ScreenshotRedirect(ctx *server.Context) {
	id, err := ctx.PathValueInt("id")
	if err != nil {
		ctx.Response.WriteHeader(http.StatusNotFound)
		return
	}

	screenshot, err := ctx.State.Screenshots.ById(id)
	if err != nil || screenshot.Hidden {
		ctx.Response.WriteHeader(http.StatusNotFound)
		return
	}

	if time.Since(screenshot.CreatedAt) > 7*24*time.Hour {
		ctx.Response.WriteHeader(http.StatusNotFound)
		return
	}

	ctx.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/ss/%d/%s", id, screenshot.Checksum()))
}

func Avatar(ctx *server.Context) {
	avatarFilename := ctx.PathValue("filename")

	// Workaround for older clients that use file extensions
	userIdString, _, _ := strings.Cut(avatarFilename, "_")
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		DefaultAvatar(ctx)
		return
	}
	size := resolveAvatarSize(ctx)

	// If a cache key is provided, the avatar may be cached by the client
	if ctx.QueryValue("c") != "" {
		ctx.Response.Header().Set("Cache-Control", "public, max-age=604800, immutable")
	}

	// Serve a previously resized avatar straight from the cache when available
	cacheKey := fmt.Sprintf("avatar:%d:%d", userId, size)
	cached, err := ctx.State.Redis.Get(context.Background(), cacheKey).Bytes()

	if err == nil && len(cached) > 0 {
		writeAvatar(ctx, cached)
		return
	}

	avatar, err := ctx.State.Storage.Read(strconv.Itoa(userId), "avatars")
	if err != nil {
		ctx.Logger.Error("Failed to read avatar", "userId", userId, "error", err)
		DefaultAvatar(ctx)
		return
	}

	// Only resize & cache the avatar for the allowed sizes
	if _, allowed := allowedAvatarSizes[size]; !allowed {
		writeAvatar(ctx, avatar)
		return
	}

	resized, err := media.ResizeImage(avatar, size, size)
	if err != nil {
		ctx.Logger.Warn("Failed to resize avatar", "userId", userId, "size", size, "error", err)
		writeAvatar(ctx, avatar)
		return
	}

	ctx.State.Redis.Set(context.Background(), cacheKey, resized, avatarCacheTTL)
	writeAvatar(ctx, resized)
}

func DefaultAvatar(ctx *server.Context) {
	defaultAvatar, err := ctx.State.Storage.ReadStream("unknown", "avatars")
	if err != nil {
		ctx.Logger.Error("Failed to read default avatar", "error", err)
		ctx.Response.WriteHeader(404)
		return
	}
	defer defaultAvatar.Close()

	ctx.Response.Header().Set("Content-Type", "image/png")
	ctx.Response.WriteHeader(200)
	io.Copy(ctx.Response, defaultAvatar)
}

func writeAvatar(ctx *server.Context, avatar []byte) {
	ctx.Response.Header().Set("Content-Type", "image/png")
	ctx.Response.WriteHeader(200)
	ctx.Response.Write(avatar)
}

func resolveAvatarSize(ctx *server.Context) int {
	size, err := ctx.QueryValueInt("s")
	if err != nil {
		return defaultAvatarSize
	}
	return size
}
