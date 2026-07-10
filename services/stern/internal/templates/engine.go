package templates

import (
	"bytes"

	"github.com/CloudyKit/jet/v6"
	"github.com/CloudyKit/jet/v6/loaders/embedfs"
	"github.com/osuTitanic/titanic/internal/bbcode"
	"github.com/osuTitanic/titanic/internal/config"
	web "github.com/osuTitanic/titanic/services/stern/web"
)

type Engine struct {
	Set *jet.Set
}

func NewEngine(cfg *config.Config) (*Engine, error) {
	bbcode.ConfigureDefault(bbcode.Options{
		BaseUrl:            cfg.OsuBaseUrl(),
		ValidImageServices: cfg.ValidImageServices(),
		ImageProxyBaseUrl:  cfg.ImageProxyBaseUrl,
		ImageProxySecret:   cfg.FrontendSecretKey,
	})

	set := jet.NewSet(
		embedfs.NewLoader("template", web.Templates),
		jet.DevelopmentMode(cfg.Reload),
	)
	staticFiles, err := web.StaticFS("/")
	if err != nil {
		return nil, err
	}
	staticURLs := newStaticUrlCache(staticFiles)

	set.AddGlobalFunc("cachedUrl", staticURLs.cachedUrl)
	set.AddGlobalFunc("formatNumber", formatNumber)
	set.AddGlobalFunc("formatFloat", formatFloat)
	set.AddGlobalFunc("formatDateShort", formatDateShort)
	set.AddGlobalFunc("round", round)
	set.AddGlobalFunc("floor", floor)
	set.AddGlobalFunc("countryName", countryName)
	set.AddGlobalFunc("homeRenderNewsText", homeRenderNewsText)
	set.AddGlobalFunc("beatmapDifficultySort", beatmapDifficultySort)
	set.AddGlobalFunc("beatmapGenres", beatmapGenres)
	set.AddGlobalFunc("beatmapLanguages", beatmapLanguages)
	set.AddGlobalFunc("beatmapStatuses", beatmapStatuses)
	set.AddGlobalFunc("shortMods", shortMods)
	set.AddGlobalFunc("scoreWeight", scoreWeight)
	set.AddGlobalFunc("bbcode", renderBBCode)
	set.AddGlobalFunc("markdownUrls", markdownUrls)
	set.AddGlobalFunc("forumUserLink", forumUserLink)
	set.AddGlobalFunc("formatActivity", formatActivity)
	set.AddGlobalFunc("beatmapRatingWidth", beatmapRatingWidth)
	set.AddGlobalFunc("searchParamUrl", searchParamUrl)
	set.AddGlobalFunc("searchFlagUrl", searchFlagUrl)
	set.AddGlobalFunc("searchSortUrl", searchSortUrl)
	set.AddGlobalFunc("searchHiddenInputs", searchHiddenInputs)
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
