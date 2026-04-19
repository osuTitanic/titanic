package templates

import (
	"github.com/osuTitanic/titanic-go/internal/config"
	"github.com/osuTitanic/titanic-go/internal/schemas"
)

type Statistics struct {
	TotalUsers  int
	OnlineUsers int
	TotalScores int
}

type DefaultView struct {
	Config      *config.Config
	CurrentUser *schemas.User
	Stats       Statistics
	CSRFToken   string
	CurrentPath string
	CurrentURI  string
}

type HomeView struct {
	DefaultView
	// TODO: Add news, chat, most played maps, ...
}

type LoginView struct {
	DefaultView
	Redirect     string
	ErrorMessage string
}
