package server

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/osuTitanic/titanic/internal/authentication"
	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/permissions"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/internal/state"
	"github.com/osuTitanic/titanic/services/stern/internal/templates"
)

const (
	httpReadHeaderTimeout = 5 * time.Second
	httpIdleTimeout       = 2 * time.Minute
	httpShutdownTimeout   = 30 * time.Second
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

// IsDebug returns true if the server is running in debug mode
func (server *Server) IsDebug() bool {
	return server.State != nil && server.State.Config != nil && server.State.Config.Debug
}

// Handle registers a stdlib route pattern on the server.
func (server *Server) Handle(pattern string, handler func(*Context)) {
	server.Router.HandleFunc(pattern, server.ContextMiddleware(handler))
}

// HandleFileSystem registers a static file handler under the provided prefix.
func (server *Server) HandleFileSystem(prefix string, instance fs.FS) {
	// Check if we are serving a directory or a single file
	if strings.HasSuffix(prefix, "/") {
		handler := http.StripPrefix(prefix, http.FileServerFS(instance))
		server.Router.Handle("GET "+prefix, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			server.SetCacheHeaders(w.Header(), r)
			handler.ServeHTTP(w, r)
		}))
		return
	}
	filename := path.Base(prefix)

	server.Router.HandleFunc("GET "+prefix, func(w http.ResponseWriter, r *http.Request) {
		server.SetCacheHeaders(w.Header(), r)
		http.ServeFileFS(w, r, instance, filename)
	})
}

// SetCacheHeaders sets the appropriate cache headers for static assets based on the request path & query parameters.
func (server *Server) SetCacheHeaders(header http.Header, request *http.Request) {
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

	resolvedPermissions *permissions.Set
}

func (ctx *Context) IP() string {
	return GetRequestIP(ctx.Request)
}

func (ctx *Context) Country() string {
	// TODO: Add geoip fallback lookup when a geolocation service exists
	// 		 For now we only trust cloudflare headers & otherwise return XX
	country := ctx.Request.Header.Get("CF-IPCountry")
	country = strings.ToUpper(strings.TrimSpace(country))

	if country == "" || country == "XX" || country == "T1" {
		// "XX" -> Unknown country
		// "T1" -> Most likely a tor exit node
		return "XX"
	}
	if constants.GetCountryIndexFromCode(country) == 0 {
		// This country does not exist in our country list
		country = "XX"
	}
	return country
}

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

// PathValue is a helper function to get path variables from the request context.
// e.g. if the route is "/users/{id}", you can get the "id" variable by calling ctx.PathValue("id").
func (ctx *Context) PathValue(name string) string {
	return ctx.Request.PathValue(name)
}

// PathValueInt does the same thing as PathValue, but tries to parse the query as an integer.
func (ctx *Context) PathValueInt(name string) (int, error) {
	pathValue := strings.TrimSpace(ctx.PathValue(name))
	return strconv.Atoi(pathValue)
}

// PathValueInt64 returns a path variable as an int64.
func (ctx *Context) PathValueInt64(name string) (int64, error) {
	pathValue := strings.TrimSpace(ctx.PathValue(name))
	return strconv.ParseInt(pathValue, 10, 64)
}

// QueryValue is a helper function to get query parameters from the request context.
func (ctx *Context) QueryValue(name string) string {
	return ctx.Request.URL.Query().Get(name)
}

// QueryValueInt returns a query parameter as an integer.
func (ctx *Context) QueryValueInt(name string) (int, error) {
	queryValue := strings.TrimSpace(ctx.QueryValue(name))
	return strconv.Atoi(queryValue)
}

// QueryValueInt64 returns a query parameter as an int64.
func (ctx *Context) QueryValueInt64(name string) (int64, error) {
	queryValue := strings.TrimSpace(ctx.QueryValue(name))
	return strconv.ParseInt(queryValue, 10, 64)
}

// QueryValueDefault attempts to get a query parameter from
// the request while falling back to the given if not present.
func (ctx *Context) QueryValueDefault(name, fallback string) string {
	if queryValue := ctx.QueryValue(name); queryValue != "" {
		return queryValue
	}
	return fallback
}

// FormValue is a helper function to get form values from the request body.
func (ctx *Context) FormValue(name string) string {
	return ctx.Request.FormValue(name)
}

// FormValueInt returns a form value as an integer.
func (ctx *Context) FormValueInt(name string) (int, error) {
	formValue := strings.TrimSpace(ctx.FormValue(name))
	return strconv.Atoi(formValue)
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

// Serve starts the server and gracefully shuts it down when ctx is cancelled.
func (server *Server) Serve(ctx context.Context) error {
	httpServer := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", server.Host, server.Port),
		Handler:           server.LoggingMiddleware(server.Router),
		ReadHeaderTimeout: httpReadHeaderTimeout,
		IdleTimeout:       httpIdleTimeout,
	}
	server.Logger.Info(
		"Listening for requests",
		"host", server.Host,
		"port", server.Port,
	)

	serveErrors := make(chan error, 1)
	go func() {
		serveErrors <- httpServer.ListenAndServe()
	}()

	select {
	case err := <-serveErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		server.Logger.Info("Shutting down server")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), httpShutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		// Shutdown leaves active connections open when its deadline expires
		// -> close them before returning
		httpServer.Close()
		<-serveErrors
		return fmt.Errorf("gracefully shut down HTTP server: %w", err)
	}

	err := <-serveErrors
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
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
	rc.WriteImplicitStatus()
	return rc.w.Write(b)
}

func (rc *ResponseContext) WriteHeader(status int) {
	// "Informational" responses do not contain the actual response status
	if status >= 100 && status < 200 && status != http.StatusSwitchingProtocols {
		rc.w.WriteHeader(status)
		return
	}
	if rc.s != 0 {
		return
	}
	rc.s = status
	rc.w.WriteHeader(status)
}

func (rc *ResponseContext) Unwrap() http.ResponseWriter {
	return rc.w
}

func (rc *ResponseContext) Flush() {
	rc.FlushError()
}

func (rc *ResponseContext) FlushError() error {
	rc.WriteImplicitStatus()
	return http.NewResponseController(rc.w).Flush()
}

func (rc *ResponseContext) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return http.NewResponseController(rc.w).Hijack()
}

func (rc *ResponseContext) Push(target string, options *http.PushOptions) error {
	writer := rc.w
	for {
		if pusher, ok := writer.(http.Pusher); ok {
			return pusher.Push(target, options)
		}
		unwrapper, ok := writer.(interface{ Unwrap() http.ResponseWriter })
		if !ok {
			return http.ErrNotSupported
		}
		writer = unwrapper.Unwrap()
	}
}

func (rc *ResponseContext) ReadFrom(reader io.Reader) (int64, error) {
	rc.WriteImplicitStatus()
	if readerFrom, ok := rc.w.(io.ReaderFrom); ok {
		return readerFrom.ReadFrom(reader)
	}

	// Hide ReadFrom from io.Copy to avoid recursively calling this method
	return io.Copy(struct{ io.Writer }{rc}, reader)
}

func (rc *ResponseContext) WriteImplicitStatus() {
	if rc.s == 0 {
		rc.s = http.StatusOK
	}
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
