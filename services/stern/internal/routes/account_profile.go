package routes

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/osuTitanic/titanic-go/internal/constants"
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
)

const maxUserpageLength = 1 << 13 // 8192

func AccountProfile(ctx *server.Context) {
	if !ctx.RequireLogin() {
		return
	}
	renderProfileSettings(ctx, "", "")
}

func AccountProfileUpdate(ctx *server.Context) {
	if !ctx.RequireLogin() {
		return
	}
	if err := ctx.Request.ParseForm(); err != nil {
		ctx.Logger.Warn("Failed to parse profile form", "error", err)
		renderProfileSettings(ctx, "", "Failed to update profile. Please try again!")
		return
	}

	interests := formValueOrNil(ctx, "interests")
	location := formValueOrNil(ctx, "location")
	website := formValueOrNil(ctx, "website")
	discord := formValueOrNil(ctx, "discord")
	twitter := formValueOrNil(ctx, "twitter")

	if discord != nil {
		stripped := strings.TrimPrefix(*discord, "@")
		discord = &stripped

		if !constants.DiscordUsername.MatchString(stripped) {
			renderProfileSettings(ctx, "", "Invalid discord username. Please try again!")
			return
		}
	}

	if interests != nil && utf8.RuneCountInString(*interests) > 30 {
		renderProfileSettings(ctx, "", "Please keep your interests short!")
		return
	}

	if location != nil && utf8.RuneCountInString(*location) > 30 {
		renderProfileSettings(ctx, "", "Please keep your location short!")
		return
	}

	if twitter != nil && utf8.RuneCountInString(*twitter) > 64 {
		renderProfileSettings(ctx, "", "Please type in a valid twitter handle or url!")
		return
	}

	if website != nil && utf8.RuneCountInString(*website) > 64 {
		renderProfileSettings(ctx, "", "Please keep your website url short!")
		return
	}

	if website != nil && !matchesURL(*website) {
		renderProfileSettings(ctx, "", "Please enter in a valid url!")
		return
	}

	if status := checkProfileAccountStatus(ctx); status != "" {
		renderProfileSettings(ctx, "", status)
		return
	}

	mode := ctx.CurrentUser.PreferredMode
	if parsed, err := ctx.FormValueInt("mode"); err == nil && parsed >= 0 && parsed <= 3 {
		mode = constants.Mode(parsed)
	}

	var twitterUrl *string
	if twitter != nil {
		formatted := "https://twitter.com/" + twitterHandle(*twitter)
		twitterUrl = &formatted
	}

	updates := &schemas.User{
		Id:            ctx.CurrentUser.Id,
		PreferredMode: mode,
		Interests:     interests,
		Location:      location,
		Website:       website,
		Discord:       discord,
		Twitter:       twitterUrl,
	}

	_, err := ctx.State.Users.Update(
		updates,
		"preferred_mode",
		"userpage_interests",
		"userpage_location",
		"userpage_website",
		"userpage_discord",
		"userpage_twitter",
	)
	if err != nil {
		ctx.Logger.Error("Failed to update profile", "user", ctx.CurrentUser.Id, "error", err)
		InternalServerError(ctx)
		return
	}

	// Reflect the changes on the in-memory user for the re-rendered form
	ctx.CurrentUser.PreferredMode = mode
	ctx.CurrentUser.Interests = interests
	ctx.CurrentUser.Location = location
	ctx.CurrentUser.Website = website
	ctx.CurrentUser.Discord = discord
	ctx.CurrentUser.Twitter = twitterUrl

	renderProfileSettings(ctx, "Successfully updated profile.", "")
}

func AccountProfileUserpage(ctx *server.Context) {
	if !ctx.RequireLogin() {
		return
	}
	if err := ctx.Request.ParseForm(); err != nil {
		ctx.Redirect(http.StatusSeeOther, "/account/profile")
		return
	}

	if !ctx.Request.PostForm.Has("bbcode") {
		ctx.Redirect(http.StatusSeeOther, "/account/profile")
		return
	}

	bbcode := ctx.FormValue("bbcode")
	userId, _ := ctx.FormValueInt("user_id")

	if ctx.CurrentUser.Id != userId && !ctx.Permissions().IsModerator() {
		ctx.Redirect(http.StatusSeeOther, "/account/profile")
		return
	}

	if status := checkProfileAccountStatus(ctx); status != "" {
		renderProfileSettings(ctx, "", status)
		return
	}

	if utf8.RuneCountInString(bbcode) > maxUserpageLength {
		renderProfileSettings(ctx, "", "Your userpage is too long!")
		return
	}

	updates := &schemas.User{Id: userId, Userpage: &bbcode}
	if _, err := ctx.State.Users.Update(updates, "userpage_about"); err != nil {
		ctx.Logger.Error("Failed to update userpage", "user", userId, "error", err)
		InternalServerError(ctx)
		return
	}

	if userId != ctx.CurrentUser.Id {
		ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/u/%d", userId))
		return
	}

	ctx.CurrentUser.Userpage = &bbcode
	ctx.Redirect(http.StatusSeeOther, "/account/profile#userpage")
}

func AccountProfileSignature(ctx *server.Context) {
	if !ctx.RequireLogin() {
		return
	}
	if err := ctx.Request.ParseForm(); err != nil {
		ctx.Redirect(http.StatusSeeOther, "/account/profile")
		return
	}

	if !ctx.Request.PostForm.Has("bbcode") {
		ctx.Redirect(http.StatusSeeOther, "/account/profile")
		return
	}

	bbcode := ctx.FormValue("bbcode")
	userId, _ := ctx.FormValueInt("user_id")

	if ctx.CurrentUser.Id != userId && !ctx.Permissions().IsAdmin() {
		ctx.Redirect(http.StatusSeeOther, "/account/profile")
		return
	}

	if status := checkProfileAccountStatus(ctx); status != "" {
		renderProfileSettings(ctx, "", status)
		return
	}

	if utf8.RuneCountInString(bbcode) > maxUserpageLength {
		renderProfileSettings(ctx, "", "Your signature is too long!")
		return
	}

	updates := &schemas.User{Id: userId, Signature: &bbcode}
	if _, err := ctx.State.Users.Update(updates, "userpage_signature"); err != nil {
		ctx.Logger.Error("Failed to update signature", "user", userId, "error", err)
		InternalServerError(ctx)
		return
	}

	if userId != ctx.CurrentUser.Id {
		ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/u/%d", userId))
		return
	}

	ctx.CurrentUser.Signature = &bbcode
	ctx.Redirect(http.StatusSeeOther, "/account/profile#signature")
}

func renderProfileSettings(ctx *server.Context, info, errorMessage string) {
	view := templates.SettingsProfileView{
		DefaultView:     buildDefaultView(ctx),
		InfoMessage:     info,
		ErrorMessage:    errorMessage,
		UserpageEditor:  templates.ForumEditorContext{Content: ctx.CurrentUser.UserpageText(), SubmitText: "Update"},
		SignatureEditor: templates.ForumEditorContext{Content: ctx.CurrentUser.SignatureText(), SubmitText: "Update"},
	}
	ctx.RenderTemplate(http.StatusOK, "pages/account/settings_profile", view)
}

func checkProfileAccountStatus(ctx *server.Context) string {
	user := ctx.CurrentUser

	switch {
	case user.Restricted:
		return "Your account was restricted."
	case user.SilenceEnd != nil && user.SilenceEnd.After(time.Now()):
		return "Your account was silenced."
	case !user.Activated:
		return "Your account is not activated."
	case !ctx.HasPermission("users.profile.update"):
		return "You are not allowed to update your profile."
	}
	return ""
}

func formValueOrNil(ctx *server.Context, name string) *string {
	value := ctx.FormValue(name)
	if value == "" {
		return nil
	}
	return &value
}

func matchesURL(value string) bool {
	loc := constants.URL.FindStringIndex(value)
	return loc != nil && loc[0] == 0
}

func twitterHandle(value string) string {
	if match := constants.TwitterHandle.FindStringSubmatch(value); match != nil {
		return match[3]
	}
	if !strings.HasPrefix(value, "@") {
		return "@" + value
	}
	return value
}
