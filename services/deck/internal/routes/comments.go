package routes

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/osuTitanic/titanic/internal/activity"
	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

type CommentRequest struct {
	Action    string
	BeatmapId int
	ReplayId  *int
	SetId     *int
	StartTime *int
	Content   string
	Color     string
	Target    constants.CommentTarget
}

func (r *CommentRequest) IsLegacy() bool {
	return r.SetId == nil || *r.SetId == 0 || r.ReplayId == nil || *r.ReplayId == 0
}

func (r *CommentRequest) Validate() error {
	if r.Action != "get" && r.Action != "post" {
		return fmt.Errorf("invalid action: %s", r.Action)
	}
	if r.BeatmapId <= 0 {
		return fmt.Errorf("invalid beatmap id: %d", r.BeatmapId)
	}
	if r.Content != "" && utf8.RuneCountInString(r.Content) > 80 {
		return fmt.Errorf("comment content too long: %d", utf8.RuneCountInString(r.Content))
	}
	r.Content = strings.NewReplacer("\t", "", "|", "").Replace(r.Content)
	return nil
}

// /web/osu-comment.php -> Retrieve or submit in-game beatmap comments
func Comments(ctx *server.Context) {
	request, err := decodeCommentRequest(ctx)
	if err != nil {
		ctx.Logger.Warn("Failed to decode comment request", "error", err)
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := request.Validate(); err != nil {
		ctx.Logger.Warn("Invalid comment request", "error", err)
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return
	}

	user, ok := ctx.HandleUserAuthenticationFormSimple("u", "p", true)
	if !ok {
		return
	}

	user.LatestActivity = time.Now()
	ctx.State.Users.Update(user, "latest_activity")

	switch request.Action {
	case "get":
		getComments(ctx, request)
	case "post":
		postComment(ctx, request, user)
	default:
		ctx.Response.WriteHeader(http.StatusBadRequest)
	}
}

func getComments(ctx *server.Context, request CommentRequest) {
	targets := []struct {
		id         *int
		targetType constants.CommentTarget
	}{
		{id: &request.BeatmapId, targetType: constants.CommentTargetMap},
		{id: request.ReplayId, targetType: constants.CommentTargetReplay},
		{id: request.SetId, targetType: constants.CommentTargetSong},
	}
	comments := make([]*schemas.BeatmapComment, 0)

	for _, target := range targets {
		if target.id == nil || *target.id <= 0 {
			continue
		}
		found, err := ctx.State.Comments.FetchByTarget(*target.id, target.targetType)
		if err != nil {
			ctx.Logger.Error("Failed to retrieve beatmap comments", "target", target.targetType, "target_id", *target.id, "error", err)
			ctx.Response.WriteHeader(http.StatusInternalServerError)
			return
		}
		comments = append(comments, found...)
	}

	legacy := request.IsLegacy()
	formatted := make([]string, len(comments))
	for i, comment := range comments {
		formatted[i] = formatBeatmapComment(comment, legacy)
	}
	ctx.RenderText(http.StatusOK, strings.Join(formatted, "\n"))
}

func postComment(ctx *server.Context, request CommentRequest, user *schemas.User) {
	if request.Content == "" {
		ctx.Logger.Warn("Failed to submit comment: no content", "user_id", user.Id)
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return
	}
	if request.StartTime == nil {
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return
	}

	targetId, ok := resolveCommentTargetId(request)
	if !ok {
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return
	}

	beatmap, err := ctx.State.Beatmaps.ById(request.BeatmapId, "Beatmapset")
	if err != nil {
		ctx.Logger.Error("Failed to retrieve beatmap", "beatmap_id", request.BeatmapId, "error", err)
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}
	if beatmap == nil {
		ctx.Logger.Warn("Failed to submit comment: beatmap not found", "beatmap_id", request.BeatmapId)
		ctx.Response.WriteHeader(http.StatusNotFound)
		return
	}

	permissions, err := ctx.State.Permissions.Resolve(user.Id)
	if err != nil {
		ctx.Logger.Error("Failed to resolve permissions", "user_id", user.Id, "error", err)
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}

	isBAT := permissions.IsBat()
	isDonator := permissions.IsDonator()
	isCreator := beatmap.Beatmapset != nil && *beatmap.Beatmapset.CreatorId == user.Id
	commentFormat := resolveCommentFormat(isCreator, isBAT, isDonator)

	comment := &schemas.BeatmapComment{
		TargetId:   targetId,
		TargetType: request.Target,
		UserId:     user.Id,
		Mode:       beatmap.Mode,
		Comment:    request.Content,
		Time:       *request.StartTime,
		Format:     &commentFormat,
		Color:      resolveCommentColor(request.Color, isDonator),
	}
	if err := ctx.State.Comments.Create(comment); err != nil {
		ctx.Logger.Error("Failed to create beatmap comment", "user_id", user.Id, "error", err)
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.Logger.Info(
		"Submitted beatmap comment",
		"user_id", user.Id, "target", request.Target, "target_id", targetId,
	)

	err = activity.Submit(
		ctx.State,
		user.Id,
		&beatmap.Mode,
		constants.ActivityBeatmapCommented,
		map[string]any{
			"username":     user.Name,
			"beatmap_id":   beatmap.Id,
			"beatmap_name": beatmap.Name(),
			"comment":      request.Content,
		},
		false, // should not be sent to #announce
		true,  // should be hidden in user profile
	)
	if err != nil {
		ctx.Logger.Warn("Failed to broadcast beatmap comment", "user_id", user.Id, "error", err)
	}

	ctx.RenderText(http.StatusOK, fmt.Sprintf("%d|%s\n", *request.StartTime, request.Content))
}

func decodeCommentRequest(ctx *server.Context) (CommentRequest, error) {
	beatmapId, err := ctx.FormValueInt("b")
	if err != nil {
		return CommentRequest{}, err
	}
	replayId, err := ctx.FormValueIntOptional("r")
	if err != nil {
		return CommentRequest{}, err
	}
	setId, err := ctx.FormValueIntOptional("s")
	if err != nil {
		return CommentRequest{}, err
	}
	startTime, err := ctx.FormValueIntOptional("starttime")
	if err != nil {
		return CommentRequest{}, err
	}
	target, err := ctx.FormValueEnum[constants.CommentTarget]("target")
	if err != nil {
		target = constants.CommentTargetMap
	}

	// NOTE: Client also provides the mode as parameter "m"
	// 		 We don't use it here though

	return CommentRequest{
		Action:    ctx.FormValue("a"),
		Target:    target,
		BeatmapId: beatmapId,
		ReplayId:  replayId,
		SetId:     setId,
		StartTime: startTime,
		Content:   ctx.FormValue("comment"),
		Color:     ctx.FormValue("f"),
	}, nil
}

func formatBeatmapComment(comment *schemas.BeatmapComment, legacy bool) string {
	if legacy {
		return fmt.Sprintf("%d|%s", comment.Time, comment.Comment)
	}

	commentFormat := ""
	if comment.Format != nil {
		commentFormat = *comment.Format
	}
	if comment.Color != nil && *comment.Color != "" {
		commentFormat += "|" + *comment.Color
	}

	return fmt.Sprintf(
		"%d\t%s\t%s\t%s",
		comment.Time, comment.TargetType,
		commentFormat, comment.Comment,
	)
}

func resolveCommentColor(color string, isDonator bool) *string {
	if color == "" || !isDonator {
		return nil
	}
	return new(color)
}

func resolveCommentFormat(isCreator, isBat, isDonator bool) string {
	switch {
	case isCreator:
		return "creator"
	case isBat:
		return "bat"
	case isDonator:
		return "subscriber"
	default:
		return "player"
	}
}

func resolveCommentTargetId(request CommentRequest) (int, bool) {
	switch request.Target {
	case constants.CommentTargetMap:
		if request.BeatmapId != 0 {
			return request.BeatmapId, true
		}
	case constants.CommentTargetReplay:
		if request.ReplayId != nil && *request.ReplayId != 0 {
			return *request.ReplayId, true
		}
	case constants.CommentTargetSong:
		if request.SetId != nil && *request.SetId != 0 {
			return *request.SetId, true
		}
	}
	return 0, false
}
