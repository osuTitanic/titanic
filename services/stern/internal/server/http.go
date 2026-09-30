package server

import (
	"errors"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/osuTitanic/titanic/internal/authentication"
	"github.com/osuTitanic/titanic/internal/permissions"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/internal/server"
	"github.com/osuTitanic/titanic/internal/state"
	"github.com/osuTitanic/titanic/services/stern/internal/templates"
)

type Server struct {
	*server.HttpServer[*Context]
	State *state.State
}

func NewServer(host string, port int, name string, state *state.State, engine *templates.Engine) *Server {
	return &Server{
		HttpServer: server.NewHttpServer(host, port, name, NewContextFactory(state, engine)),
		State:      state,
	}
}

// IsDebug returns true if the server is running in debug mode
func (server *Server) IsDebug() bool {
	return server.State != nil && server.State.Config != nil && server.State.Config.Debug
}

// HandleFileSystem registers a static file handler under the provided prefix.
func (server *Server) HandleFileSystem(prefix string, instance fs.FS) {
	if strings.HasSuffix(prefix, "/") {
		server.serveDirectory(prefix, instance)
		return
	}
	server.serveFile(prefix, instance)
}

func (server *Server) serveDirectory(prefix string, instance fs.FS) {
	fileServer := http.FileServerFS(instance)

	handler := http.StripPrefix(prefix, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(
			path.Clean(r.URL.Path), // path traversal prevention
			"/",                    // strip leading slash
		)

		// Go's FileServerFS exposes the directory structure
		// by default which is not what we want
		if info, err := fs.Stat(instance, name); err == nil && info.IsDir() {
			http.NotFound(w, r)
			return
		}
		fileServer.ServeHTTP(w, r)
	}))

	server.Router.Handle("GET "+prefix, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.setCacheHeaders(w.Header(), r)
		handler.ServeHTTP(w, r)
	}))
}

func (server *Server) serveFile(prefix string, instance fs.FS) {
	filename := path.Base(prefix)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, instance, filename)
	})

	server.Router.Handle("GET "+prefix, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.setCacheHeaders(w.Header(), r)
		handler.ServeHTTP(w, r)
	}))
}

func (server *Server) setCacheHeaders(header http.Header, request *http.Request) {
	if server.IsDebug() {
		// No caching in debug mode pretty please
		return
	}
	hasChecksum := request.URL.Query().Has("c")

	if strings.HasPrefix(request.URL.Path, "/images/") && !hasChecksum {
		// Images basically won't change so we can cache them for a week
		header.Set("Cache-Control", "public, max-age=604800")
		return
	}

	// Only cache the following paths if we have a "c" parameter
	// This ensures that we can deploy new versions of static assets
	// without worrying about users having stale cached versions
	if !hasChecksum {
		return
	}

	cacheableStaticPaths := [...]string{
		"/js/",
		"/css/",
		"/lib/",
		"/images/",
		"/webfonts/",
	}
	for _, prefix := range cacheableStaticPaths {
		if strings.HasPrefix(request.URL.Path, prefix) {
			header.Set("Cache-Control", "public, max-age=31536000, immutable")
			return
		}
	}
}

// Context is a struct that holds the request context for each endpoint call.
type Context struct {
	server.HttpContext
	State          *state.State
	Templates      *templates.Engine
	CurrentUser    *schemas.User
	CurrentSession *authentication.WebsiteSession
	CSRFToken      string

	resolvedPermissions *permissions.Set
}

func NewContextFactory(state *state.State, engine *templates.Engine) server.HttpContextFactory[*Context] {
	return func(base server.HttpContext) *Context {
		ctx := &Context{
			HttpContext: base,
			State:       state,
			Templates:   engine,
		}
		ctx.ResolveAuthentication()
		return ctx
	}
}

// RequireLogin redirects the user to the login page if they are not authenticated.
func (ctx *Context) RequireLogin() bool {
	if ctx.IsAuthenticated() {
		return true
	}
	ctx.Redirect(
		http.StatusSeeOther,
		"/account/login?redirect="+ctx.Request.URL.RequestURI(),
	)
	return false
}

// HasPermission returns true if the user has the specified permission.
func (ctx *Context) HasPermission(permission string) bool {
	return ctx.Permissions().Has(permission)
}

// Permissions resolves & memoizes the current user's permission set for this request
func (ctx *Context) Permissions() *permissions.Set {
	if ctx.resolvedPermissions != nil {
		return ctx.resolvedPermissions
	}

	ctx.resolvedPermissions = &permissions.Set{}
	if ctx.CurrentUser == nil {
		return ctx.resolvedPermissions
	}

	set, err := ctx.State.Permissions.Resolve(ctx.CurrentUser.Id)
	if err != nil {
		ctx.Logger.Error("Failed to resolve permissions", "user", ctx.CurrentUser.Id, "error", err)
		return ctx.resolvedPermissions
	}

	ctx.resolvedPermissions = set
	return ctx.resolvedPermissions
}

// RenderTemplate renders the specified template with the given context & writes the result to the response.
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
