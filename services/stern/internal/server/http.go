package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/osuTitanic/titanic-go/internal/authentication"
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"github.com/osuTitanic/titanic-go/internal/state"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
)

// Server is the main struct that holds the state for an http server.
type Server struct {
	Host      string
	Port      int
	Name      string
	State     *state.State
	Logger    *slog.Logger
	Router    *http.ServeMux
	Templates *templates.Engine
}

// Handle registers a stdlib route pattern on the server.
func (server *Server) Handle(pattern string, handler func(*Context)) {
	server.Router.HandleFunc(pattern, server.ContextMiddleware(handler))
}

// HandleFileSystem registers a static file handler under the provided prefix.
func (server *Server) HandleFileSystem(prefix string, instance fs.FS) {
	if strings.HasSuffix(prefix, "/") {
		server.Router.Handle("GET "+prefix, http.StripPrefix(prefix, http.FileServerFS(instance)))
		return
	}
	filename := path.Base(prefix)

	server.Router.HandleFunc("GET "+prefix, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, instance, filename)
	})
}

func NewServer(host string, port int, name string, state *state.State, engine *templates.Engine) *Server {
	return &Server{
		Host:      host,
		Port:      port,
		Name:      name,
		State:     state,
		Templates: engine,
		Logger:    slog.Default().With("component", name),
		Router:    http.NewServeMux(),
	}
}

// Context is a struct that holds the request context for each endpoint call.
type Context struct {
	Response       http.ResponseWriter
	Request        *http.Request
	State          *state.State
	Templates      *templates.Engine
	Logger         *slog.Logger
	CurrentUser    *schemas.User
	CurrentSession *authentication.WebsiteSession
	CSRFToken      string
}

func (ctx *Context) IP() string {
	return GetRequestIP(ctx.Request)
}

// PathValue is a helper function to get path variables from the request context.
// e.g. if the route is "/users/{id}", you can get the "id" variable by calling ctx.PathValue("id").
func (ctx *Context) PathValue(name string) string {
	return ctx.Request.PathValue(name)
}

func (ctx *Context) Redirect(status int, location string) {
	http.Redirect(ctx.Response, ctx.Request, location, status)
}

func (ctx *Context) RenderTemplate(status int, name string, data any) error {
	if ctx.Templates == nil {
		err := errors.New("templates engine is not configured")
		ctx.Logger.Error("Failed to render template", "template", name, "error", err)
		templates.InternalServerErrorFallback(ctx.Response)
		return err
	}

	body, err := ctx.Templates.Render(name, data)
	if err != nil {
		ctx.Logger.Error("Failed to render template", "template", name, "error", err)
		templates.InternalServerErrorFallback(ctx.Response)
		return err
	}

	ctx.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx.Response.WriteHeader(status)
	_, err = ctx.Response.Write(body)
	return err
}

func (ctx *Context) RenderJson(status int, data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	ctx.Response.Header().Set("Content-Type", "application/json; charset=utf-8")
	ctx.Response.WriteHeader(status)
	_, err = ctx.Response.Write(payload)
	return err
}

// Serve starts the HTTP server and listens for incoming requests.
func (server *Server) Serve() {
	bind := fmt.Sprintf(
		"%s:%d",
		server.Host,
		server.Port,
	)
	server.Logger.Info(
		"Listening for requests",
		"host", server.Host,
		"port", server.Port,
	)

	err := http.ListenAndServe(bind, server.LoggingMiddleware(server.Router))
	if err != nil {
		server.Logger.Error("Failed to start server", "error", err)
		return
	}
}

// ResponseContext is a wrapper around http.ResponseWriter that
// allows us to capture the status code of a response.
type ResponseContext struct {
	w http.ResponseWriter
	s int
}

func (rc *ResponseContext) Header() http.Header {
	return rc.w.Header()
}

func (rc *ResponseContext) Write(b []byte) (int, error) {
	return rc.w.Write(b)
}

func (rc *ResponseContext) WriteHeader(status int) {
	rc.s = status
	rc.w.WriteHeader(status)
}

func (rc *ResponseContext) Status() int {
	if rc.s == 0 {
		return http.StatusOK
	}
	return rc.s
}

// ContextMiddleware creates a new Context struct for each request.
func (server *Server) ContextMiddleware(handler func(*Context)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		context := &Context{
			Response:  w,
			Request:   r,
			State:     server.State,
			Templates: server.Templates,
			Logger:    server.Logger,
		}

		w.Header().Set("Server", server.Name)
		context.ResolveAuthentication()
		handler(context)
	}
}

// LoggingMiddleware logs the details of each request.
func (server *Server) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rc := &ResponseContext{w: w}
		start := time.Now()
		next.ServeHTTP(rc, r)

		server.Logger.Info(
			fmt.Sprintf("%s %s", r.Method, r.RequestURI),
			"ip", GetRequestIP(r),
			"status", rc.Status(),
			"duration", time.Since(start).String(),
		)
	})
}
