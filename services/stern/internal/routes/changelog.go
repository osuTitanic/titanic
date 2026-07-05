package routes

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/osuTitanic/titanic-go/internal/schemas"
	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
)

var (
	changelogClientCutoff      = time.Date(2015, time.December, 30, 0, 0, 0, 0, time.UTC)
	changelogClientCutoffOsume = time.Date(2014, time.August, 11, 0, 0, 0, 0, time.UTC)
)

// OsuChangelog emulates the osu! in-game updater changelog endpoint.
// The behaviour is different depending on the "updater" query parameter:
//
//	1 -> osume changelog page
//	2 -> osume changelog page (test stream)
//	3 -> plain-text client changelog stream (e.g. https://i.ibb.co/kgQRJ4fK/image.png)
//	_ -> redirect to the modern changelog page (currently a todo)
func OsuChangelog(ctx *server.Context) {
	updater, _ := ctx.QueryValueInt("updater")

	switch updater {
	case 3:
		clientChangelog(ctx)
	case 2, 1:
		osumeChangelog(ctx)
	default:
		ctx.Redirect(http.StatusFound, "/changelog")
	}
}

func osumeChangelog(ctx *server.Context) {
	startDate := resolveChangelogTargetDate(ctx.Request.URL.Query(), changelogClientCutoffOsume)
	entries, err := ctx.State.Repositories.ReleasesOfficial.FetchChangelogRangeDesc(startDate, 50, 0)
	if err != nil {
		ctx.Logger.Error("Failed to fetch changelog entries", "error", err)
		InternalServerError(ctx)
		return
	}

	view := &templates.ChangelogOsumeView{
		Config: ctx.State.Config,
		Groups: groupChangelogByDate(entries),
	}
	ctx.RenderTemplate(http.StatusOK, "pages/public/changelog_osume", view)
}

func clientChangelog(ctx *server.Context) {
	startDate := resolveChangelogTargetDate(ctx.Request.URL.Query(), changelogClientCutoff)
	entries, err := ctx.State.Repositories.ReleasesOfficial.FetchChangelogRangeDesc(startDate, 100, 0)
	if err != nil {
		ctx.Logger.Error("Failed to fetch changelog entries", "error", err)
		InternalServerError(ctx)
		return
	}

	var builder strings.Builder
	for i, group := range groupChangelogByDate(entries) {
		if i > 0 {
			builder.WriteByte('\n')
		}

		// The date header preserves the "(m/d/yyyy)" format
		date := group.Entries[0].CreatedAt
		fmt.Fprintf(&builder, "(%d/%d/%d)\n", int(date.Month()), date.Day(), date.Year())

		for j, entry := range group.Entries {
			if j > 0 {
				builder.WriteByte('\n')
			}
			fmt.Fprintf(&builder, "%s\t%s\t%s", entry.TypeSymbol(), entry.Author, entry.PrefixedText())
		}
	}

	ctx.Response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	ctx.Response.WriteHeader(http.StatusOK)
	ctx.Response.Write([]byte(builder.String()))
}

func groupChangelogByDate(entries []*schemas.ReleaseChangelog) []*templates.ChangelogGroup {
	groups := make([]*templates.ChangelogGroup, 0)
	var current *templates.ChangelogGroup

	for _, entry := range entries {
		key := entry.CreatedAt.Format("Jan 2, 2006")
		if current == nil || current.Date != key {
			current = &templates.ChangelogGroup{Date: key}
			groups = append(groups, current)
		}
		current.Entries = append(current.Entries, entry)
	}

	return groups
}

func resolveChangelogTargetDate(query url.Values, fallback time.Time) time.Time {
	// from & current are both in the format YYYYMMDD
	from, _ := strconv.Atoi(query.Get("from"))
	current, _ := strconv.Atoi(query.Get("current"))

	// current takes precedence over from
	version := current
	if version == 0 {
		version = from
	}

	versionString := strconv.Itoa(version)
	if len(versionString) < 8 {
		return fallback
	}

	year, _ := strconv.Atoi(versionString[0:4])
	month, _ := strconv.Atoi(versionString[4:6])
	day, _ := strconv.Atoi(versionString[6:8])
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
