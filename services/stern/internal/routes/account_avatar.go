package routes

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/osuTitanic/titanic/internal/media"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/services/stern/internal/server"
)

const (
	maxAvatarBytes     = int64(2.5 * 1024 * 1024)
	maxAvatarDimension = 8000
	avatarResizeSize   = 256
)

func AccountAvatarRedirect(ctx *server.Context) {
	if !ctx.RequireLogin() {
		return
	}
	ctx.Redirect(http.StatusSeeOther, "/account/profile")
}

func AccountAvatarUpdate(ctx *server.Context) {
	if !ctx.RequireLogin() {
		return
	}

	user := ctx.CurrentUser
	file, _, err := ctx.Request.FormFile("avatar")
	if err != nil {
		renderProfileSettings(ctx, "", "Please provide a valid image!")
		return
	}
	defer file.Close()

	switch {
	case user.Restricted:
		renderProfileSettings(ctx, "", "Your account was restricted.")
		return
	case user.SilenceEnd != nil && user.SilenceEnd.After(time.Now()):
		renderProfileSettings(ctx, "", "Your account was silenced.")
		return
	case !user.Activated:
		renderProfileSettings(ctx, "", "Your account is not activated.")
		return
	}

	// Read at most one byte over the limit, so we can detect oversized uploads
	data, err := io.ReadAll(io.LimitReader(file, maxAvatarBytes+1))
	if err != nil {
		renderProfileSettings(ctx, "", "Please provide a valid image!")
		return
	}
	if int64(len(data)) > maxAvatarBytes {
		renderProfileSettings(ctx, "", "This image is too large. Please upload an image below 2.5mb!")
		return
	}

	config, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		ctx.Logger.Warn("Failed to read avatar image", "user", user.Id, "error", err)
		renderProfileSettings(ctx, "", "Please provide a valid image!")
		return
	}

	if config.Width > maxAvatarDimension || config.Height > maxAvatarDimension {
		renderProfileSettings(ctx, "", "This image is too large. Please lower the resolution!")
		return
	}

	resized, err := media.ResizeImageToPng(data, avatarResizeSize, avatarResizeSize)
	if err != nil {
		ctx.Logger.Warn("Failed to resize avatar image", "user", user.Id, "error", err)
		renderProfileSettings(ctx, "", "Please provide a valid image!")
		return
	}

	if err := ctx.State.Storage.Save(strconv.Itoa(user.Id), "avatars", resized); err != nil {
		ctx.Logger.Error("Failed to upload avatar", "user", user.Id, "error", err)
		InternalServerError(ctx)
		return
	}

	checksum := md5.Sum(resized)
	hash := hex.EncodeToString(checksum[:])

	updates := &schemas.User{Id: user.Id, AvatarHash: &hash, AvatarLastUpdate: time.Now()}
	user.AvatarHash = &hash
	user.AvatarLastUpdate = updates.AvatarLastUpdate

	if _, err := ctx.State.Users.Update(updates, "avatar_hash", "avatar_last_changed"); err != nil {
		ctx.Logger.Error("Failed to update avatar hash", "user", user.Id, "error", err)
		InternalServerError(ctx)
		return
	}
	invalidateAvatarCaches(ctx, user.Id)

	ctx.Logger.Info("User changed their avatar", "user", user.Id)
	ctx.Redirect(http.StatusSeeOther, "/account/profile")
}

func invalidateAvatarCaches(ctx *server.Context, userId int) {
	keys := []string{fmt.Sprintf("bancho:avatar_hash:%d", userId)}
	for size := range allowedAvatarSizes {
		keys = append(keys, fmt.Sprintf("avatar:%d:%d", userId, size))
	}
	ctx.State.Redis.Del(context.Background(), keys...)
}
