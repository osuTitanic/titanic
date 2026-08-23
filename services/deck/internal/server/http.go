package server

import (
	"github.com/osuTitanic/titanic/internal/server"
	"github.com/osuTitanic/titanic/internal/state"
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
	State *state.State
}

func NewContextFactory(state *state.State) server.HttpContextFactory[*Context] {
	return func(base server.HttpContext) *Context {
		return &Context{
			HttpContext: base,
			State:       state,
		}
	}
}
