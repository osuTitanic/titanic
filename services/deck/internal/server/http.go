package server

import (
	"github.com/osuTitanic/titanic/internal/server"
	"github.com/osuTitanic/titanic/internal/state"
	"github.com/osuTitanic/titanic/services/deck/internal/bss"
)

type Server struct {
	*server.HttpServer[*Context]
}

func NewServer(host string, port int, name string, state *state.State) *Server {
	return &Server{
		HttpServer: server.NewHttpServer(host, port, name, NewContextFactory(state)),
	}
}

type Context struct {
	server.HttpContext
	State             *state.State
	StagedSubmissions *bss.StagedSubmissionStore
}

func NewContextFactory(state *state.State) server.HttpContextFactory[*Context] {
	stagedSubmissions := bss.NewStagedSubmissionStore(state.Redis)

	return func(base server.HttpContext) *Context {
		return &Context{
			HttpContext:       base,
			State:             state,
			StagedSubmissions: stagedSubmissions,
		}
	}
}
