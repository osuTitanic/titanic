package templates

import (
	"net/url"

	"github.com/osuTitanic/titanic-go/internal/config"
	"github.com/osuTitanic/titanic-go/internal/schemas"
)

type Statistics struct {
	TotalUsers  int
	OnlineUsers int
	TotalScores int
}

type DefaultView struct {
	Query       url.Values
	Config      *config.Config
	CurrentUser *schemas.User
	Stats       Statistics
	CSRFToken   string
	CurrentPath string
	CurrentURI  string
}

type HomeView struct {
	DefaultView
	// TODO: use pointers here
	News               []schemas.ForumPost
	ChatMessages       []schemas.Message
	MostPlayedBeatmaps map[int]*schemas.Beatmap
}

type LoginView struct {
	DefaultView
	Redirect     string
	ErrorMessage string
}

type DownloadView struct {
	DefaultView
	SelectedCategory string
	Categories       []*DownloadCategory
	Clients          []*schemas.Release
}

type DownloadCategory struct {
	Name     string
	Url      string
	Selected bool
}

type BeatmapSearchView struct {
	DefaultView
	Beatmapsets []*schemas.Beatmapset
	SearchSort  string
	SearchOrder string
	Pagination  PaginationView
}
