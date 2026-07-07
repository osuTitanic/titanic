package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/osuTitanic/titanic-go/internal/schemas"
	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
	"github.com/osuTitanic/titanic-go/services/stern/internal/wiki"
)

// TODO: Add redirect that redirects /wiki/<path> to /wiki/<language>/<path>

func WikiRedirect(ctx *server.Context) {
	language := wiki.DefaultLanguageFromConfig(ctx.State.Config)
	ctx.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/wiki/%s/", language))
}

func WikiLanguageRedirect(ctx *server.Context) {
	language, ok := wikiRouteLanguage(ctx)
	if !ok {
		return
	}
	ctx.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/wiki/%s/", language))
}

func WikiHome(ctx *server.Context) {
	language, ok := wikiRouteLanguage(ctx)
	if !ok {
		return
	}

	// TODO: Persist wiki service somewhere, not sure where yet
	service := wiki.NewService(ctx.State.Config, ctx.State.Repositories, ctx.Logger)
	urls := service.URLs()

	pageCount, err := ctx.State.WikiPages.Count()
	if err != nil {
		ctx.Logger.Error("Failed to fetch wiki page count", "error", err)
		InternalServerError(ctx)
		return
	}

	categories, err := wikiHomeCategories(ctx, language)
	if err != nil {
		ctx.Logger.Error("Failed to fetch wiki categories", "error", err)
		InternalServerError(ctx)
		return
	}

	view := templates.WikiView{
		DefaultView:        buildDefaultView(ctx),
		Title:              "Wiki",
		SiteTitle:          "Wiki - Titanic!",
		SiteDescription:    "Browse through Titanic! wiki articles",
		CanonicalUrl:       fmt.Sprintf("/wiki/%s/", language),
		RequestedLanguage:  language,
		Language:           language,
		SourceUrl:          urls.GitHubBase,
		DiscussionUrl:      urls.GitHubBase + "/pulls",
		HistoryUrl:         urls.History,
		PageCount:          pageCount,
		CurrentDate:        wiki.CurrentDate(language, time.Now()),
		AvailableLanguages: wikiTemplateLanguageLinks(language),
		Categories:         categories,
	}
	ctx.RenderTemplate(http.StatusOK, "pages/wiki/home/"+language, view)
}

func WikiSearchRedirect(ctx *server.Context) {
	language, ok := wikiRouteLanguage(ctx)
	if !ok {
		return
	}

	target := fmt.Sprintf("/wiki/%s/search/", language)
	if rawQuery := ctx.Request.URL.RawQuery; rawQuery != "" {
		target += "?" + rawQuery
	}
	ctx.Redirect(http.StatusMovedPermanently, target)
}

func WikiSearch(ctx *server.Context) {
	language, ok := wikiRouteLanguage(ctx)
	if !ok {
		return
	}

	query := ctx.QueryValue("query")
	title := query
	if title == "" {
		title = "Search"
	}
	// TODO: Implement search functionality

	view := templates.WikiView{
		DefaultView:       buildDefaultView(ctx),
		Title:             title + " - Wiki",
		SiteTitle:         title + " - Titanic! Wiki",
		SiteDescription:   "Search the Titanic! wiki",
		CanonicalUrl:      fmt.Sprintf("/wiki/%s/search/", language),
		RequestedLanguage: language,
		Language:          language,
		SearchQuery:       query,
		SourceUrl:         "#",
		HistoryUrl:        "#",
		DiscussionUrl:     "#",
	}
	ctx.RenderTemplate(http.StatusOK, "pages/wiki/search/"+language, view)
}

func WikiArticle(ctx *server.Context) {
	language, ok := wikiRouteLanguage(ctx)
	if !ok {
		return
	}

	path := strings.TrimSuffix(ctx.PathValue("path"), "/")
	if path == "" {
		NotFound(ctx)
		return
	}

	// TODO: Persist wiki service somewhere, not sure where yet
	service := wiki.NewService(ctx.State.Config, ctx.State.Repositories, ctx.Logger)
	result, err := service.FetchPage(path, language)
	if err != nil {
		ctx.Logger.Error("Failed to fetch wiki page", "path", path, "language", language, "error", err)
		InternalServerError(ctx)
		return
	}
	if result == nil || result.Page == nil || result.Content == nil {
		NotFound(ctx)
		return
	}

	formattedPath := wiki.FormatPath(path, result.Page.Name)
	if formattedPath != path {
		ctx.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/wiki/%s/%s", language, formattedPath))
		return
	}

	renderedContent, err := wiki.RenderMarkdown(result.Content.Content, language)
	if err != nil {
		ctx.Logger.Error("Failed to render wiki markdown", "path", path, "language", language, "error", err)
		InternalServerError(ctx)
		return
	}

	githubPath := wiki.GitHubPath(path)
	urls := service.URLs()
	siteURL := fmt.Sprintf("/wiki/%s/%s", result.Content.Language, path)

	view := templates.WikiView{
		DefaultView:       buildDefaultView(ctx),
		Title:             result.Content.Title + " - Wiki",
		SiteTitle:         result.Content.Title + " - Titanic! Wiki",
		SiteDescription:   "Titanic » Wiki » " + result.Content.Title,
		SiteUrl:           siteURL,
		CanonicalUrl:      siteURL,
		Content:           renderedContent,
		RequestedLanguage: language,
		Language:          result.Content.Language,
		TranslationUrl:    fmt.Sprintf("%s/%s", urls.Create, githubPath),
		SourceUrl:         fmt.Sprintf("%s/%s/%s.md", urls.BlobBase, githubPath, result.Content.Language),
		HistoryUrl:        fmt.Sprintf("%s/%s/%s.md", urls.History, githubPath, result.Content.Language),
		DiscussionUrl:     fmt.Sprintf("%s/pulls?q=%s", urls.GitHubBase, url.QueryEscape(result.Page.Name)),
	}
	ctx.RenderTemplate(http.StatusOK, "pages/wiki/content/"+language, view)
}

func wikiRouteLanguage(ctx *server.Context) (string, bool) {
	language := wiki.NormalizeLanguage(ctx.PathValue("language"))
	if !wiki.IsSupportedLanguage(language) {
		NotFound(ctx)
		return "", false
	}
	return language, true
}

func wikiTemplateLanguageLinks(language string) []templates.WikiLanguageLink {
	links := wiki.AvailableLanguagesExcept(language)
	result := make([]templates.WikiLanguageLink, 0, len(links))
	for _, link := range links {
		result = append(result, templates.WikiLanguageLink{
			Code: link.Code,
			Name: link.Name,
		})
	}
	return result
}

func wikiHomeCategories(ctx *server.Context, language string) ([]*templates.WikiCategoryView, error) {
	categories, err := ctx.State.WikiCategories.MainCategories()
	if err != nil {
		return nil, err
	}

	result := make([]*templates.WikiCategoryView, 0, len(categories))
	for _, category := range categories {
		view := &templates.WikiCategoryView{
			Name:  wikiCategoryName(category, language),
			Pages: make([]*templates.WikiPageLink, 0, len(category.Pages)),
		}

		// Fetch all available pages for this category & translate their titles if necessary
		for _, page := range category.Pages {
			title := page.Name
			if language != wiki.DefaultLanguageFromConfig(ctx.State.Config) {
				translatedTitle, err := ctx.State.WikiContents.TranslatedTitle(page.Id, language)
				if err != nil {
					return nil, err
				}
				if translatedTitle != "" {
					title = translatedTitle
				}
				// i love nested code
			}

			view.Pages = append(view.Pages, &templates.WikiPageLink{
				Path:  page.Path,
				Title: title,
			})
		}
		result = append(result, view)
	}

	return result, nil
}

func wikiCategoryName(category *schemas.WikiCategory, language string) string {
	if category == nil {
		return ""
	}
	if language == wiki.DefaultLanguage {
		return category.Name
	}

	var translations map[string]string
	if len(category.Translations) > 0 {
		// Translations are stored as json.RawMessage, which requires us to unmarshal them
		// We could actually make a query to fetch the translation directly from the database,
		// but uhhh.... i frogot that this is a thing and i'm lazy and yes deal with it lol
		if err := json.Unmarshal(category.Translations, &translations); err == nil {
			if name := translations[language]; name != "" {
				return name
			}
		}
	}
	return category.Name
}
