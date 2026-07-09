package routes

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/osuTitanic/titanic-go/internal/media"
	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
)

var allowedAvatarSizes = map[int]struct{}{
	25:  {},
	128: {},
	256: {},
}

const defaultAvatarSize = 128
const avatarCacheTTL = time.Hour * 24

func Avatar(ctx *server.Context) {
	avatarFilename := ctx.PathValue("filename")

	// Workaround for older clients that use file extensions
	userIdString := strings.SplitN(avatarFilename, "_", 2)[0]
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
