package templates

import (
	"bytes"

	"github.com/CloudyKit/jet/v6"
	"github.com/CloudyKit/jet/v6/loaders/embedfs"
	"github.com/osuTitanic/titanic-go/internal/config"
	web "github.com/osuTitanic/titanic-go/services/stern/web"
)

type Engine struct {
	Set *jet.Set
}

func NewEngine(cfg *config.Config) (*Engine, error) {
	set := jet.NewSet(
		embedfs.NewLoader("template", web.Templates),
		jet.DevelopmentMode(cfg.Reload),
	)
	set.AddGlobalFunc("formatNumber", formatNumber)
	set.AddGlobalFunc("homeNewsDate", homeNewsDate)
	set.AddGlobalFunc("homeNewsLink", homeNewsLink)
	set.AddGlobalFunc("homeNewsTitle", homeNewsTitle)
	set.AddGlobalFunc("homeNewsAuthor", homeNewsAuthor)
	set.AddGlobalFunc("homeNewsText", homeNewsText)
	set.AddGlobalFunc("homeChatTime", homeChatTime)
	set.AddGlobalFunc("homeMostPlayedRows", homeMostPlayedRows)
	set.AddGlobalFunc("homeBeatmapsetURL", homeBeatmapsetURL)
	set.AddGlobalFunc("homeBeatmapsetName", homeBeatmapsetName)
	set.AddGlobalFunc("homeBeatmapsetCreatorName", homeBeatmapsetCreatorName)
	set.AddGlobalFunc("homeBeatmapsetCreatorURL", homeBeatmapsetCreatorUrl)
	return &Engine{Set: set}, nil
}

func (e *Engine) Render(name string, data any) ([]byte, error) {
	view, err := e.Set.GetTemplate(name)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := view.Execute(&buf, nil, data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
