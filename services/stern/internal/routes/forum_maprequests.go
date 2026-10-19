package routes

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	osuapi "github.com/Lekuruu/gosu/pkg/api"
)

const (
	mapRequestDateLayout = "2006/01/02"
	mapRequestMaxSetId   = 1000000000 // bss starting setid integer
)

var (
	errMapRequestInvalidLink        = errors.New("invalid bancho beatmapset link")
	errMapRequestBeatmapsetNotFound = errors.New("bancho beatmapset not found")
	errMapRequestIncompleteMetadata = errors.New("bancho beatmapset metadata is incomplete")
)

type MapRequest struct {
	BeatmapsetId int
	Title        string
	Content      string
}

type MapRequestMode struct {
	Ruleset osuapi.Ruleset
	Title   string
	Icon    string
}

var mapRequestModes = [...]MapRequestMode{
	{Ruleset: osuapi.RulesetOsu, Title: "Standard", Icon: "osu"},
	{Ruleset: osuapi.RulesetTaiko, Title: "Taiko", Icon: "taiko"},
	{Ruleset: osuapi.RulesetFruits, Title: "Catch", Icon: "ctb"},
	{Ruleset: osuapi.RulesetMania, Title: "Mania", Icon: "mania"},
}

func buildMapRequest(ctx context.Context, client *osuapi.Client, link string, comment string) (*MapRequest, error) {
	beatmapsetId, err := parseMapRequestLink(link)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, errors.New("build forum map request: beatmapset client is nil")
	}
	var apiError *osuapi.Error

	beatmapset, err := client.GetBeatmapset(
		ctx, int64(beatmapsetId),
	)
	if errors.As(err, &apiError) && apiError.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%w: %d", errMapRequestBeatmapsetNotFound, beatmapsetId)
	}
	if err != nil {
		return nil, fmt.Errorf("fetch bancho beatmapset %d: %w", beatmapsetId, err)
	}
	if beatmapset == nil {
		return nil, fmt.Errorf("%w: beatmapset %d was empty", errMapRequestIncompleteMetadata, beatmapsetId)
	}
	if beatmapset.ID != int64(beatmapsetId) {
		return nil, fmt.Errorf(
			"%w: requested beatmapset %d but received %d",
			errMapRequestIncompleteMetadata, beatmapsetId, beatmapset.ID,
		)
	}

	modes, beatmap, err := resolveMapRequestMetadata(beatmapset)
	if err != nil {
		return nil, err
	}

	return &MapRequest{
		BeatmapsetId: beatmapsetId,
		Title:        buildMapRequestTitle(beatmapset, modes),
		Content:      buildMapRequestContent(beatmapset, beatmap, modes, comment),
	}, nil
}

func parseMapRequestLink(rawLink string) (int, error) {
	rawLink = strings.TrimSpace(rawLink)
	parsed, err := url.Parse(rawLink)
	if err != nil {
		return 0, errMapRequestInvalidLink
	}

	isHttp := strings.EqualFold(parsed.Scheme, "http") || strings.EqualFold(parsed.Scheme, "https")
	isValidHostname := strings.EqualFold(parsed.Hostname(), "osu.ppy.sh")
	if !isHttp || !isValidHostname {
		return 0, errMapRequestInvalidLink
	}

	path := strings.Trim(parsed.Path, "/")
	parts := strings.Split(path, "/")
	if len(parts) != 2 {
		return 0, errMapRequestInvalidLink
	}
	// accept both /s/<id> and /beatmapsets/<id>#<mode>/<mapid>
	if parts[0] != "beatmapsets" && parts[0] != "s" {
		return 0, errMapRequestInvalidLink
	}

	beatmapsetId, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		return 0, errMapRequestInvalidLink
	}
	if beatmapsetId <= 0 || beatmapsetId > mapRequestMaxSetId {
		return 0, errMapRequestInvalidLink
	}
	return int(beatmapsetId), nil
}

func resolveMapRequestMetadata(beatmapset *osuapi.Beatmapset) ([]MapRequestMode, *osuapi.Beatmap, error) {
	if strings.TrimSpace(beatmapset.Artist) == "" || strings.TrimSpace(beatmapset.Title) == "" {
		return nil, nil, fmt.Errorf(
			"%w: beatmapset %d is missing its artist or title",
			errMapRequestIncompleteMetadata, beatmapset.ID,
		)
	}
	var presentModes = make(map[osuapi.Ruleset]bool, len(mapRequestModes))
	var longestBeatmap *osuapi.Beatmap

	for i := range beatmapset.Beatmaps {
		beatmap := &beatmapset.Beatmaps[i]
		if beatmap.Convert {
			continue
		}
		presentModes[beatmap.Mode] = true

		if longestBeatmap == nil {
			longestBeatmap = beatmap
			continue
		}
		if beatmap.TotalLength > longestBeatmap.TotalLength {
			longestBeatmap = beatmap
		}
	}
	if longestBeatmap == nil {
		return nil, nil, fmt.Errorf(
			"%w: beatmapset %d does not contain any beatmaps",
			errMapRequestIncompleteMetadata, beatmapset.ID,
		)
	}

	modes := make([]MapRequestMode, 0, len(presentModes))
	for _, mode := range mapRequestModes {
		if presentModes[mode.Ruleset] {
			modes = append(modes, mode)
		}
	}
	return modes, longestBeatmap, nil
}

func buildMapRequestTitle(beatmapset *osuapi.Beatmapset, modes []MapRequestMode) string {
	modeNames := make([]string, len(modes))
	for i, mode := range modes {
		modeNames[i] = mode.Title
	}

	title := fmt.Sprintf(
		"[%s] %s - %s",
		strings.Join(modeNames, "/"),
		strings.TrimSpace(beatmapset.Artist),
		strings.TrimSpace(beatmapset.Title),
	)
	return truncateTitle(title, forumTitleMaxLength)
}

func buildMapRequestContent(
	beatmapset *osuapi.Beatmapset,
	beatmap *osuapi.Beatmap,
	modes []MapRequestMode,
	comments string,
) string {
	var content strings.Builder

	content.WriteString("[centre][b][u]Map Information[/u][/b]\n")
	for _, mode := range modes {
		fmt.Fprintf(
			&content,
			"[img]https://osu.titanic.sh/images/icons/forum/%s.gif[/img]",
			mode.Icon,
		)
	}
	content.WriteByte('\n')

	fmt.Fprintf(
		&content,
		"https://osu.ppy.sh/beatmapsets/%d\n",
		beatmapset.ID,
	)
	fmt.Fprintf(
		&content,
		"Created: [b]%s[/b] | Last Updated: [b]%s[/b]\n",
		formatDate(beatmapset.SubmittedDate),
		formatDate(&beatmapset.LastUpdated),
	)
	fmt.Fprintf(
		&content,
		"BPM: [b]%s[/b] | Length: [b]%s (%s)[/b]\n",
		formatBPM(beatmapset.BPM),
		formatDuration(beatmap.TotalLength),
		formatDuration(beatmap.HitLength),
	)
	content.WriteString("[/centre]")

	if comments = strings.TrimSpace(comments); comments != "" {
		content.WriteString("\n\n")
		content.WriteString(comments)
	}
	return content.String()
}

func formatDate(value *time.Time) string {
	if value == nil {
		return "Unknown"
	}
	if value.IsZero() {
		return "Unknown"
	}
	return value.UTC().Format(mapRequestDateLayout)
}

func formatBPM(value float64) string {
	if value <= 0 {
		return "Unknown"
	}
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return "Unknown"
	}
	value = math.Round(value*100) / 100
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func formatDuration(seconds int) string {
	seconds = max(seconds, 0)
	return fmt.Sprintf("%d:%02d", seconds/60, seconds%60)
}

func truncateTitle(title string, limit int) string {
	if len(title) <= limit {
		return title
	}
	if limit <= 3 {
		return strings.Repeat(".", max(limit, 0))
	}

	end := limit - 3
	for end > 0 && !utf8.RuneStart(title[end]) {
		end--
	}
	return strings.TrimSpace(title[:end]) + "..."
}
