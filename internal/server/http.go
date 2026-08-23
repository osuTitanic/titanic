package server

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/osuTitanic/titanic/internal/constants"
)

const (
	httpReadHeaderTimeout = 5 * time.Second
	httpIdleTimeout       = 2 * time.Minute
	httpShutdownTimeout   = 30 * time.Second
)

// HttpContext is a struct that holds the request context for each endpoint call.
type HttpContext struct {
	Response http.ResponseWriter
	Request  *http.Request
	Logger   *slog.Logger
}

func (ctx *HttpContext) IP() string {
	return GetRequestIP(ctx.Request)
}

func (ctx *HttpContext) Country() string {
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

// PathValue is a helper function to get path variables from the request context.
// e.g. if the route is "/users/{id}", you can get the "id" variable by calling ctx.PathValue("id").
func (ctx *HttpContext) PathValue(name string) string {
	return ctx.Request.PathValue(name)
}

// PathValueInt does the same thing as PathValue, but tries to parse the query as an integer.
func (ctx *HttpContext) PathValueInt(name string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(ctx.PathValue(name)))
}

// PathValueInt64 returns a path variable as an int64.
func (ctx *HttpContext) PathValueInt64(name string) (int64, error) {
	return strconv.ParseInt(strings.TrimSpace(ctx.PathValue(name)), 10, 64)
}

// QueryValue is a helper function to get query parameters from the request context.
func (ctx *HttpContext) QueryValue(name string) string {
	return ctx.Request.URL.Query().Get(name)
}

// QueryValueInt returns a query parameter as an integer.
func (ctx *HttpContext) QueryValueInt(name string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(ctx.QueryValue(name)))
}

// QueryValueInt64 returns a query parameter as an int64.
func (ctx *HttpContext) QueryValueInt64(name string) (int64, error) {
	return strconv.ParseInt(strings.TrimSpace(ctx.QueryValue(name)), 10, 64)
}

// QueryValueDefault attempts to get a query parameter from
// the request while falling back to the given if not present.
func (ctx *HttpContext) QueryValueDefault(name, fallback string) string {
	if value := ctx.QueryValue(name); value != "" {
		return value
	}
	return fallback
}

// FormValue is a helper function to get form values from the request body.
func (ctx *HttpContext) FormValue(name string) string {
	return ctx.Request.FormValue(name)
}

// FormValueInt returns a form value as an integer.
func (ctx *HttpContext) FormValueInt(name string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(ctx.FormValue(name)))
}

func (ctx *HttpContext) Redirect(status int, location string) {
	http.Redirect(ctx.Response, ctx.Request, location, status)
}

func (ctx *HttpContext) RenderText(status int, data string) error {
	ctx.Response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	ctx.Response.WriteHeader(status)
	_, err := ctx.Response.Write([]byte(data))
	return err
}

func (ctx *HttpContext) RenderJson(status int, data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	ctx.Response.Header().Set("Content-Type", "application/json; charset=utf-8")
	ctx.Response.WriteHeader(status)
	_, err = ctx.Response.Write(payload)
	return err
}

// HttpContextFactory turns the shared http context into a service context.
type HttpContextFactory[C any] func(HttpContext) C

// HttpHandlerFunc handles a request using a service's individual context.
type HttpHandlerFunc[C any] func(C)

// HttpServer provides http routing & lifecycle management for a service.
type HttpServer[C any] struct {
	Host   string
	Port   int
	Name   string
	Logger *slog.Logger
	Router *http.ServeMux

	contextFactory HttpContextFactory[C]
}

func NewHttpServer[C any](host string, port int, name string, createContext HttpContextFactory[C]) *HttpServer[C] {
	return &HttpServer[C]{
		Host:           host,
		Port:           port,
		Name:           name,
		Logger:         slog.Default().With("component", name),
		Router:         http.NewServeMux(),
		contextFactory: createContext,
	}
}

// Handle registers a stdlib route pattern on the server.
func (server *HttpServer[C]) Handle(pattern string, handler HttpHandlerFunc[C]) {
	server.Router.HandleFunc(pattern, server.ContextMiddleware(handler))
}

// ContextMiddleware creates a new Context struct for each request.
func (server *HttpServer[C]) ContextMiddleware(handler HttpHandlerFunc[C]) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set(
			"Server",
			server.Name,
		)
		base := HttpContext{
			Response: response,
			Request:  request,
			Logger:   server.Logger,
		}
		handler(server.contextFactory(base))
	}
}

// Serve starts the server and gracefully shuts it down when ctx is cancelled.
func (server *HttpServer[C]) Serve(ctx context.Context) error {
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

// HttpResponseContext is a wrapper around http.ResponseWriter that
// allows us to capture the status code of a response.
type HttpResponseContext struct {
	w http.ResponseWriter
	s int
}

func (ctx *HttpResponseContext) Header() http.Header {
	return ctx.w.Header()
}

func (ctx *HttpResponseContext) Write(data []byte) (int, error) {
	ctx.WriteImplicitStatus()
	return ctx.w.Write(data)
}

func (ctx *HttpResponseContext) WriteHeader(status int) {
	// "Informational" responses do not contain the actual response status
	if status >= 100 && status < 200 && status != http.StatusSwitchingProtocols {
		ctx.w.WriteHeader(status)
		return
	}
	if ctx.s != 0 {
		return
	}
	ctx.s = status
	ctx.w.WriteHeader(status)
}

func (ctx *HttpResponseContext) Unwrap() http.ResponseWriter {
	return ctx.w
}

func (ctx *HttpResponseContext) Flush() {
	ctx.FlushError()
}

func (ctx *HttpResponseContext) FlushError() error {
	ctx.WriteImplicitStatus()
	return http.NewResponseController(ctx.w).Flush()
}

func (ctx *HttpResponseContext) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return http.NewResponseController(ctx.w).Hijack()
}

func (ctx *HttpResponseContext) Push(target string, options *http.PushOptions) error {
	writer := ctx.w
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

func (ctx *HttpResponseContext) ReadFrom(reader io.Reader) (int64, error) {
	ctx.WriteImplicitStatus()
	if readerFrom, ok := ctx.w.(io.ReaderFrom); ok {
		return readerFrom.ReadFrom(reader)
	}

	// Hide ReadFrom from io.Copy to avoid recursively calling this method
	return io.Copy(struct{ io.Writer }{ctx}, reader)
}

func (ctx *HttpResponseContext) WriteImplicitStatus() {
	if ctx.s == 0 {
		ctx.s = http.StatusOK
	}
}

func (ctx *HttpResponseContext) Status() int {
	if ctx.s == 0 {
		return http.StatusOK
	}
	return ctx.s
}

// LoggingMiddleware logs the details of each request.
func (server *HttpServer[C]) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		ctx := &HttpResponseContext{w: response}
		start := time.Now()
		next.ServeHTTP(ctx, request)

		server.Logger.Info(
			fmt.Sprintf("%s %s", request.Method, request.RequestURI),
			"ip", GetRequestIP(request),
			"status", ctx.Status(),
			"duration", time.Since(start).String(),
		)
	})
}
