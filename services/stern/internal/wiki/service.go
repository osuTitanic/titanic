package wiki

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/osuTitanic/titanic-go/internal/config"
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"github.com/osuTitanic/titanic-go/internal/state"
)

const markdownCacheTTL = 5 * time.Minute
const markdownFetchLimit = 10 << 20 // ~10 MB

type Service struct {
	cfg    *config.Config
	repos  *state.Repositories
	logger *slog.Logger
	client *http.Client
	urls   URLs
}

type PageResult struct {
	Page    *schemas.WikiPage
	Content *schemas.WikiContent
}

type CachedMarkdown struct {
	content string
	found   bool
	expires time.Time
}

var markdownCache = struct {
	sync.Mutex
	entries map[string]CachedMarkdown
}{
	entries: make(map[string]CachedMarkdown),
}

func NewService(cfg *config.Config, repos *state.Repositories, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{
		cfg:    cfg,
		repos:  repos,
		logger: logger.With("component", "wiki"),
		client: &http.Client{Timeout: 15 * time.Second},
		urls:   BuildURLs(cfg),
	}
}

func (s *Service) URLs() URLs {
	return s.urls
}

func (s *Service) DefaultLanguage() string {
	return DefaultLanguageFromConfig(s.cfg)
}

// FetchPage retrieves a wiki page and its content for the specified path and language.
// It automatically creates and updates pages and content as needed.
func (s *Service) FetchPage(path, language string) (*PageResult, error) {
	if s == nil || s.repos == nil {
		return nil, fmt.Errorf("wiki: service is not configured")
	}

	language = NormalizeLanguage(language)
	path = PagePath(path)

	page, err := s.repos.WikiPages.ByPath(path)
	if err != nil {
		return nil, err
	}
	if page == nil {
		s.logger.Info("Page not found in database, creating...", "path", path)
		return s.createPage(path, language)
	}

	defaultLanguage := s.DefaultLanguage()
	defaultContent, err := s.repos.WikiContents.ByPageLanguage(page.Id, defaultLanguage)
	if err != nil {
		return nil, err
	}
	if defaultContent == nil {
		s.logger.Error("Default content missing, recreating page...", "path", path)
		if err := s.DeletePageData(page.Id); err != nil {
			return nil, err
		}
		return s.createPage(path, language)
	}

	if language == defaultLanguage {
		// Update the default content if it has changed, and return it
		content, err := s.updateContent(page.Path, defaultContent)
		if err != nil {
			return nil, err
		}
		return &PageResult{Page: page, Content: content}, nil
	}

	// Find the content for the requested language, or create it if it doesn't exist
	content, err := s.repos.WikiContents.ByPageLanguage(page.Id, language)
	if err != nil {
		return nil, err
	}
	if content == nil {
		// Content for the requested language doesn't exist, create it if available
		contentMarkdown, found := s.FetchMarkdownCached(page.Path, language)
		if !found {
			// Content for this language doesn't exist, so we return the default content
			return &PageResult{Page: page, Content: defaultContent}, nil
		}

		content = &schemas.WikiContent{
			PageId:   page.Id,
			Language: language,
			Title:    titleOrName(contentMarkdown, page.Name),
			Content:  contentMarkdown,
		}
		if err := s.repos.WikiContents.Create(content); err != nil {
			return nil, err
		}
		return &PageResult{Page: page, Content: content}, nil
	}

	// Update the content for this language if it has changed, and return it
	content, err = s.updateContent(page.Path, content)
	if err != nil {
		return nil, err
	}
	return &PageResult{Page: page, Content: content}, nil
}

func (s *Service) FetchMarkdownCached(path, language string) (string, bool) {
	cacheKey := s.MarkdownUrl(path, language)
	now := time.Now()

	// Check the cache first & if the entry has expired
	markdownCache.Lock()
	entry, ok := markdownCache.entries[cacheKey]
	if ok && now.Before(entry.expires) {
		markdownCache.Unlock()
		return entry.content, entry.found
	}
	markdownCache.Unlock()

	// Entry was not found or has expired -> fetch the markdown from the source
	content, found := s.FetchMarkdown(path, language)

	// Persist the result in the cache
	markdownCache.Lock()
	markdownCache.entries[cacheKey] = CachedMarkdown{
		content: content,
		found:   found,
		expires: now.Add(markdownCacheTTL),
	}
	markdownCache.Unlock()

	return content, found
}

func (s *Service) FetchMarkdown(path, language string) (string, bool) {
	targetURL := s.MarkdownUrl(path, language)
	request, err := http.NewRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		s.logger.Error("Failed to create markdown request", "url", targetURL, "error", err)
		return "", false
	}
	request.Header.Set("User-Agent", "osuTitanic/wiki")

	response, err := s.client.Do(request)
	if err != nil {
		s.logger.Error("Failed to fetch markdown", "url", targetURL, "error", err)
		return "", false
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		s.logger.Error("Failed to fetch markdown", "url", response.Request.URL.String(), "status", response.StatusCode)
		return "", false
	}

	// Limit the content to prevent excessive memory usage
	body, err := io.ReadAll(io.LimitReader(response.Body, markdownFetchLimit))
	if err != nil {
		s.logger.Error("Failed to read markdown response", "url", targetURL, "error", err)
		return "", false
	}

	return SanitizeMarkdown(string(body)), true
}

func (s *Service) createPage(path, language string) (*PageResult, error) {
	s.logger.Debug(
		"Creating new page",
		"path", path, "language", language,
	)

	defaultLanguage := s.DefaultLanguage()
	defaultContentMarkdown, found := s.FetchMarkdownCached(path, defaultLanguage)
	if !found {
		s.logger.Error("Page not found in default language", "path", path, "language", defaultLanguage)
		return nil, nil
	}

	page := &schemas.WikiPage{
		Name: PageName(path),
		Path: PagePath(path),
	}
	if err := s.repos.WikiPages.Create(page); err != nil {
		return nil, err
	}

	defaultContent := &schemas.WikiContent{
		PageId:   page.Id,
		Language: defaultLanguage,
		Title:    titleOrName(defaultContentMarkdown, page.Name),
		Content:  defaultContentMarkdown,
	}
	if err := s.repos.WikiContents.Create(defaultContent); err != nil {
		return nil, err
	}

	if err := s.createOutlinks(page.Id, defaultContentMarkdown); err != nil {
		return nil, err
	}

	if language == defaultLanguage {
		return &PageResult{Page: page, Content: defaultContent}, nil
	}

	contentMarkdown, found := s.FetchMarkdownCached(path, language)
	if !found {
		s.logger.Info("Page only available in default language", "path", path, "language", language)
		return &PageResult{Page: page, Content: defaultContent}, nil
	}

	content := &schemas.WikiContent{
		PageId:   page.Id,
		Language: language,
		Title:    titleOrName(contentMarkdown, page.Name),
		Content:  contentMarkdown,
	}
	if err := s.repos.WikiContents.Create(content); err != nil {
		return nil, err
	}

	return &PageResult{Page: page, Content: content}, nil
}

func (s *Service) updateContent(path string, entry *schemas.WikiContent) (*schemas.WikiContent, error) {
	s.logger.Debug(
		"Updating content",
		"path", path, "language", entry.Language, "last_updated", entry.LastUpdated,
	)

	contentMarkdown, found := s.FetchMarkdownCached(path, entry.Language)
	if !found {
		s.logger.Warn(
			"Content no longer exists, deleting page data",
			"page_id", entry.PageId,
		)
		if err := s.DeletePageData(entry.PageId); err != nil {
			return nil, err
		}
		return entry, nil
	}

	// Content has not changed, no need to update
	if contentMarkdown == entry.Content {
		return entry, nil
	}

	entry.Content = contentMarkdown
	entry.Title = titleOrName(contentMarkdown, entry.Title)
	entry.LastUpdated = time.Now()

	if _, err := s.repos.WikiContents.Update(entry, "content", "title", "last_updated"); err != nil {
		return nil, err
	}
	if err := s.createOutlinks(entry.PageId, contentMarkdown); err != nil {
		return nil, err
	}
	return entry, nil
}

func (s *Service) createOutlinks(pageId int, content string) error {
	s.logger.Debug("Creating outlinks", "page_id", pageId)

	links := ExtractOutlinks(content)
	if len(links) <= 0 {
		// Delete any existing outlinks for this page if there are no links in the content
		_, err := s.repos.WikiOutlinks.DeleteByPageId(pageId)
		return err
	}

	// Ensure that we only create unique outlinks by using a
	// map to track target IDs and visited paths
	targetIds := make(map[int]struct{}, len(links))
	visitedPaths := make(map[string]struct{}, len(links))

	for _, link := range links {
		pathKey := strings.ToLower(PagePath(link.Path))
		if _, ok := visitedPaths[pathKey]; ok {
			continue
		}
		visitedPaths[pathKey] = struct{}{}

		// Recursively fetch the target page to ensure it exists and get its ID
		result, err := s.FetchPage(link.Path, s.DefaultLanguage())
		if err != nil {
			return err
		}
		if result == nil || result.Page == nil {
			continue
		}
		targetIds[result.Page.Id] = struct{}{}
	}

	// Delete old entries & create new ones
	if _, err := s.repos.WikiOutlinks.DeleteByPageId(pageId); err != nil {
		return err
	}
	s.logger.Debug("Creating new outlinks", "page_id", pageId, "outlinks", targetIds)

	outlinks := make([]*schemas.WikiOutlink, 0, len(targetIds))
	for targetId := range targetIds {
		outlinks = append(outlinks, &schemas.WikiOutlink{
			PageId:   pageId,
			TargetId: targetId,
		})
	}
	return s.repos.WikiOutlinks.CreateMany(outlinks)
}

func (s *Service) DeletePageData(pageId int) error {
	s.logger.Debug("Deleting page data", "page_id", pageId)

	if _, err := s.repos.WikiOutlinks.DeleteByPageId(pageId); err != nil {
		return err
	}
	if _, err := s.repos.WikiContents.DeleteByPageId(pageId); err != nil {
		return err
	}
	_, err := s.repos.WikiPages.DeleteById(pageId)
	return err
}

func (s *Service) MarkdownUrl(path, language string) string {
	path = strings.TrimSuffix(GitHubPath(path), "/")
	language = NormalizeLanguage(language)
	return fmt.Sprintf("%s/%s/%s.md", strings.TrimRight(s.urls.Content, "/"), path, language)
}

func titleOrName(markdown, fallback string) string {
	if title := ParseTitle(markdown); title != "" {
		return title
	}
	return fallback
}
