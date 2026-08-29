package routes

import (
	"io"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/osuTitanic/titanic/internal/activity"
	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

const (
	screenshotMaxSize     = 4 << 20
	screenshotMaxBodySize = screenshotMaxSize + 64<<10
)

var allowedScreenshotFilenames = []string{"jpg", "png", "ss"}

// /web/osu-screenshot.php -> Upload an in-game screenshot (Shift+F12)
func Screenshot(ctx *server.Context) {
	user, ok := ctx.HandleUserAuthenticationSimple("u", "p", true)
	if !ok {
		// Response already set by function
		return
	}

	image, ok := readScreenshot(ctx)
	if !ok {
		// Response already set by function
		return
	}

	screenshot := &schemas.Screenshot{UserId: user.Id}
	if err := ctx.State.Screenshots.Create(screenshot); err != nil {
		ctx.Logger.Error("Failed to create screenshot", "user_id", user.Id, "error", err)
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}

	key := strconv.Itoa(screenshot.Id)
	if err := ctx.State.Storage.Save(key, "screenshots", image); err != nil {
		ctx.Logger.Error("Failed to store screenshot", "id", screenshot.Id, "error", err)
		deleteScreenshotRecord(ctx, screenshot)
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}

	user.LatestActivity = time.Now()
	ctx.State.Users.Update(user, "latest_activity")

	err := activity.Submit(
		ctx.State,
		user.Id,
		nil, // mode independent
		constants.ActivityScreenshotUploaded,
		map[string]any{"username": user.Name},
		false, // should not be sent to #announce
		true,  // should be hidden in user profile
	)
	if err != nil {
		ctx.Logger.Warn("Failed to publish screenshot activity", "id", screenshot.Id, "error", err)
	}

	ctx.Logger.Info("Screenshot uploaded", "user_id", user.Id, "id", screenshot.Id)
	ctx.RenderText(http.StatusOK, key)
}

func readScreenshot(ctx *server.Context) ([]byte, bool) {
	ctx.Request.Body = http.MaxBytesReader(
		ctx.Response,
		ctx.Request.Body,
		screenshotMaxBodySize,
	)

	file, header, err := ctx.Request.FormFile("ss")
	if err != nil {
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return nil, false
	}
	defer file.Close()

	if !slices.Contains(allowedScreenshotFilenames, header.Filename) {
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return nil, false
	}

	image, err := io.ReadAll(file)
	if err != nil || len(image) > screenshotMaxSize {
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return nil, false
	}

	switch http.DetectContentType(image) {
	case "image/jpeg", "image/png":
		return image, true
	default:
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return nil, false
	}
}

func deleteScreenshotRecord(ctx *server.Context, screenshot *schemas.Screenshot) {
	if err := ctx.State.Screenshots.Delete(screenshot); err != nil {
		ctx.Logger.Error("Failed to remove screenshot record", "id", screenshot.Id, "error", err)
	}
}
