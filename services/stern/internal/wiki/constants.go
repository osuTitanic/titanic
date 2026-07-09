package wiki

import (
	"fmt"
	"strings"

	"github.com/osuTitanic/titanic/internal/config"
)

const DefaultLanguage = "en"

var Languages = []string{
	"en",
	"ru",
	"de",
	"pl",
	"fi",
	"et",
	"fr",
	"es",
}
var LanguageNames = map[string]string{
	"en": "English",
	"ru": "Русский",
	"de": "Deutsch",
	"pl": "Polski",
	"fi": "Suomi",
	"et": "Eesti",
	"fr": "Français",
	"es": "Español",
}

type LanguageLink struct {
	Code string
	Name string
}

type URLs struct {
	GitHubBase string // https://github.com/osuTitanic/wiki
	BlobBase   string // https://github.com/osuTitanic/wiki/blob/main
	History    string // https://github.com/osuTitanic/wiki/commits/main
	Create     string // https://github.com/osuTitanic/wiki/new/main
	Content    string // https://raw.githubusercontent.com/osuTitanic/wiki/refs/heads/main
}

// IsSupportedLanguage checks if the given language is supported
func IsSupportedLanguage(language string) bool {
	_, ok := LanguageNames[strings.ToLower(language)]
	return ok
}

// NormalizeLanguage normalizes the language code to lowercase and trims whitespace
func NormalizeLanguage(language string) string {
	return strings.ToLower(strings.TrimSpace(language))
}

// DefaultLanguageFromConfig returns the default language from the config,
// or the default language if not set or unsupported
func DefaultLanguageFromConfig(cfg *config.Config) string {
	if cfg != nil && IsSupportedLanguage(cfg.WikiDefaultLanguage) {
		return NormalizeLanguage(cfg.WikiDefaultLanguage)
	}
	return DefaultLanguage
}

// AvailableLanguagesExcept returns a list of available languages except the specified one
func AvailableLanguagesExcept(language string) []LanguageLink {
	language = NormalizeLanguage(language)
	links := make([]LanguageLink, 0, len(Languages)-1)
	for _, code := range Languages {
		if code == language {
			continue
		}
		links = append(links, LanguageLink{
			Code: code,
			Name: LanguageNames[code],
		})
	}
	return links
}

func GitHubPath(path string) string {
	return strings.ReplaceAll(strings.TrimSuffix(path, "/"), " ", "_")
}

func BuildURLs(cfg *config.Config) URLs {
	owner := "osuTitanic"
	name := "wiki"
	branch := "main"
	repositoryPath := "wiki"

	if cfg != nil {
		if value := strings.Trim(cfg.WikiRepositoryOwner, "/"); value != "" {
			owner = value
		}
		if value := strings.Trim(cfg.WikiRepositoryName, "/"); value != "" {
			name = value
		}
		if value := strings.Trim(cfg.WikiRepositoryBranch, "/"); value != "" {
			branch = value
		}
		if value := strings.Trim(cfg.WikiRepositoryPath, "/"); value != "" {
			repositoryPath = value
		}
	}

	githubBase := fmt.Sprintf("https://github.com/%s/%s", owner, name)
	repositoryBase := strings.TrimRight(fmt.Sprintf("%s/%s", branch, repositoryPath), "/")

	return URLs{
		GitHubBase: githubBase,
		BlobBase:   fmt.Sprintf("%s/blob/%s", githubBase, repositoryBase),
		History:    fmt.Sprintf("%s/commits/%s", githubBase, repositoryBase),
		Create:     fmt.Sprintf("%s/new/%s", githubBase, repositoryBase),
		Content:    fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/refs/heads/%s", owner, name, repositoryBase),
	}
}
