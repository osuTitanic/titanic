# HTTP Server

This module provides a shared http server used by Titanic's web services. It allows us to define shared http logic once, and extend it with each service's individual requirements.  
In the future, this module should contain standardized TCP & WS logic as well, primarily for the `anchor` service.

## Usage

Define a service context that embeds `HttpContext`, then provide a factory that creates it for each request.

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
	"127.0.0.1", 8080,   // host & port
	"server",            // logger component and Server header
	NewContextFactory(app),
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

Register a panic handler when the service should turn recovered panics into a response or notification:

```go
httpServer.OnPanic(func(ctx *Context, recovered any) {
	ctx.Logger.Error("Request panicked", "error", recovered)

	if ctx.StatusExists() {
		// Status was already set :/
		// We can't respond with 500 anymore
		panic(http.ErrAbortHandler)
	}
	
	ctx.RenderText(http.StatusInternalServerError, "internal server error")
})
```

## Request Context

`HttpContext` lets you access the response writer, request & component logger. It also provides helper functions for common handler operations:

- `PathValue`, `QueryValue`, `QueryValues`, `FormValue` & `PostFormValue` read request values
- `Int`, `Int64`, `Float`, `Optional`, `Default` & `Enum` variants parse or normalize values
- `LimitBody`, `EnsureMultipartForm` & `FormFile` handle bounded uploads
- `Redirect`, `RenderText` & `RenderJson` write common responses
- `Status` & `StatusExists` return the response status if set
- `IP` to resolve the client IP address
- `Country` to resolve Cloudflare's `CF-IPCountry` header (`XX` if not found)

Services can extend `HttpContext` to add other dependencies, authentication, templates, etc..
