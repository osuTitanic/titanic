package routes

import (
	"io"
	"strconv"
	"strings"

	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
)

func Avatar(ctx *server.Context) {
	avatarFilename := ctx.PathValue("filename")

	// Workaround for older clients that use file extensions
	userIdString := strings.SplitN(avatarFilename, "_", 2)[0]
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		DefaultAvatar(ctx)
		return
	}
	// TODO: Caching & Resizing

	avatar, err := ctx.State.Storage.ReadStream(strconv.Itoa(userId), "avatars")
	if err != nil {
		ctx.Logger.Error("Failed to read avatar", "userId", userId, "error", err)
		DefaultAvatar(ctx)
		return
	}

	ctx.Response.Header().Set("Content-Type", "image/png")
	ctx.Response.WriteHeader(200)
	io.Copy(ctx.Response, avatar)
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
