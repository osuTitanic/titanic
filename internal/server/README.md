# HTTP Server

This module provides a shared http server used by Titanic's web services. It allows us to define shared http logic once, and extend it with each service's individual requirements.  
In the future, this module should contain standardized TCP & WS logic as well, primarily for the `anchor` service.

## Usage

Define a service context that includes `HttpContext`, then provide a factory that creates it for every request, passing it on to every request handler.

```go
type Context struct {
	server.HttpContext
	State *state.State
}

func NewContextFactory(app *state.State) server.HttpContextFactory[*Context] {
	return func(base server.HttpContext) *Context {
		return &Context{
			HttpContext: base,
			State:       app,
		}
	}
}

httpServer := server.NewHttpServer(
	"127.0.0.1", 8080, 		// Host & Port
	"server", 				// Server name to use for logging & "server" header
	NewContextFactory(app), // Context generator / factory
)

httpServer.Handle("GET /users/{id}", func(ctx *Context) {
	id, err := ctx.PathValueInt64("id")
	if err != nil {
		ctx.RenderText(http.StatusBadRequest, "invalid user id")
		return
	}

	ctx.RenderJson(http.StatusOK, map[string]any{"id": id})
})

if err := httpServer.Serve(ctx); err != nil {
	return err
}
```

`Handle` accepts the usual standard library patterns. `Serve` listens until the context is cancelled, then gives active requests up to 30 seconds to finish before closing their connections.

Each response receives a `Server` header containing the configured server name. Furthermore, requests are logged by default with their method, path, client ip, status code & processing duration.

## Request Context

`HttpContext` lets you access the response writer, request & component logger. It also provides helper functions for common handler operations:

- `PathValue`, `QueryValue` & `FormValue` to read request values
	- their `Int` and `Int64` variants parse numeric values
- `Redirect`, `RenderText` & `RenderJson` for common responses
- `IP` to resolve the client address
- `Country` to resolve Cloudflare's `CF-IPCountry` header (`XX` if not found)

Services can extend `HttpContext` to add other dependencies, authentication, templates, etc..
