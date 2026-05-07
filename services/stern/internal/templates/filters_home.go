package templates

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/CloudyKit/jet/v6"
	"github.com/osuTitanic/titanic-go/internal/constants"
	"github.com/osuTitanic/titanic-go/internal/schemas"
)

type HomeMostPlayedRow struct {
	PlayCount  int
	Beatmapset schemas.Beatmapset
}

var homeNewsIgnoredTags = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\[/?b\]`),
	regexp.MustCompile(`(?i)\[/?centre\]`),
	regexp.MustCompile(`(?i)\[size(?:=[^\]]+)?\]`),
	regexp.MustCompile(`(?i)\[/size\]`),
}

/*
 * TODO: Move some of these filters into schema struct functions
 */

func homeNewsDate(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("homeNewsDate", 1, 1)

	post := a.Get(0).Interface().(schemas.ForumPost)
	return reflect.ValueOf(post.CreatedAt.Format("2.1.2006"))
}

func homeNewsLink(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("homeNewsLink", 1, 1)

	post := a.Get(0).Interface().(schemas.ForumPost)
	return reflect.ValueOf(fmt.Sprintf("/forum/%d/t/%d/", post.ForumId, post.TopicId))
}

func homeNewsTitle(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("homeNewsTitle", 1, 1)
	post := a.Get(0).Interface().(schemas.ForumPost)

	if post.Topic != nil && post.Topic.Title != "" {
		return reflect.ValueOf(post.Topic.Title)
	}
	return reflect.ValueOf("Untitled")
}

func homeNewsAuthor(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("homeNewsAuthor", 1, 1)
	post := a.Get(0).Interface().(schemas.ForumPost)

	if post.User != nil && post.User.Name != "" {
		return reflect.ValueOf(post.User.Name)
	}
	return reflect.ValueOf(fmt.Sprintf("User %d", post.UserId))
}

func homeNewsText(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("homeNewsText", 1, 1)
	post := a.Get(0).Interface().(schemas.ForumPost)

	for line := range strings.SplitSeq(post.Content, "\n") {
		lowercaseLine := strings.ToLower(line)
		if strings.Contains(lowercaseLine, "[heading]") || strings.Contains(lowercaseLine, "[img]") {
			continue
		}

		content := strings.TrimSpace(line)
		for _, regex := range homeNewsIgnoredTags {
			content = regex.ReplaceAllString(content, "")
		}

		// TODO: BBCode rendering
		content = strings.TrimSpace(content)
		if content != "" {
			return reflect.ValueOf(content)
		}
	}
	return reflect.ValueOf("")
}

func homeChatTime(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("homeChatTime", 1, 1)
	message := a.Get(0).Interface().(schemas.Message)

	return reflect.ValueOf(message.Time.Format("15:04:05"))
}

func homeMostPlayedRows(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("homeMostPlayedRows", 1, 1)

	beatmapsets := a.Get(0).Interface().(map[int]schemas.Beatmapset)
	rows := make([]HomeMostPlayedRow, 0, len(beatmapsets))

	for playCount, beatmapset := range beatmapsets {
		rows = append(rows, HomeMostPlayedRow{
			PlayCount:  playCount,
			Beatmapset: beatmapset,
		})
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].PlayCount > rows[j].PlayCount
	})
	return reflect.ValueOf(rows)
}

func homeBeatmapsetURL(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("homeBeatmapsetURL", 1, 1)
	beatmapset := a.Get(0).Interface().(schemas.Beatmapset)

	if len(beatmapset.Beatmaps) > 0 && beatmapset.Beatmaps[0] != nil && beatmapset.Beatmaps[0].Id > 0 {
		return reflect.ValueOf(fmt.Sprintf("/b/%d", beatmapset.Beatmaps[0].Id))
	}
	return reflect.ValueOf(fmt.Sprintf("/s/%d", beatmapset.Id))
}

func homeBeatmapsetName(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("homeBeatmapsetName", 1, 1)

	beatmapset := a.Get(0).Interface().(schemas.Beatmapset)
	name := beatmapset.Name()

	if len(beatmapset.Beatmaps) > 0 && beatmapset.Beatmaps[0] != nil && beatmapset.Beatmaps[0].Version != "" {
		name = fmt.Sprintf("%s [%s]", name, beatmapset.Beatmaps[0].Version)
	}
	return reflect.ValueOf(name)
}

func homeBeatmapsetCreatorName(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("homeBeatmapsetCreatorName", 1, 1)
	beatmapset := a.Get(0).Interface().(schemas.Beatmapset)

	if beatmapset.CreatorUser != nil && beatmapset.CreatorUser.Name != "" {
		return reflect.ValueOf(beatmapset.CreatorUser.Name)
	}
	if beatmapset.Creator != nil && *beatmapset.Creator != "" {
		return reflect.ValueOf(*beatmapset.Creator)
	}
	return reflect.ValueOf("Unknown")
}

func homeBeatmapsetCreatorUrl(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("homeBeatmapsetCreatorURL", 1, 1)
	beatmapset := a.Get(0).Interface().(schemas.Beatmapset)

	if beatmapset.CreatorId != nil {
		return reflect.ValueOf(fmt.Sprintf("/u/%d", *beatmapset.CreatorId))
	}
	if beatmapset.Server == constants.BeatmapServerBancho && beatmapset.Creator != nil && *beatmapset.Creator != "" {
		return reflect.ValueOf(fmt.Sprintf("https://osu.ppy.sh/u/%s", *beatmapset.Creator))
	}
	return reflect.ValueOf("#")
}
