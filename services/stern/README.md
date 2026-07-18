# Stern

Stern is the website frontend for Titanic!, inspired by the design of the [old osu! website](https://old.ppy.sh), internally called "osu! web 1.0".

## Usage

Run stern from the repository root so it can load the `.env` & resolve the default static directory:

```sh
go run ./services/stern/cmd/web
```

The server listens on `FRONTEND_HOST` and `FRONTEND_PORT`, which default to `localhost:8080`.

### Static Assets

Jet templates are embedded into the binary. However, static assets are served from `services/stern/web/static` by default, or from the directory set through `STERN_STATIC_DIR`.

Use the `static_embedded` build tag when a self-contained binary should include the static assets:

```sh
go build -tags static_embedded -o stern ./services/stern/cmd/web
```

Keep in mind that this will increase the binary size by a lot. You will also need to rebuild the binary whenever static assets are changed. The option is there though, if you really want that.

The Docker image does not embed or copy static assets. It expects them to be mounted at `/app/static`.

## Adding a Page

First, define the data exposed to the template in `internal/templates/views.go`.
Full pages should embed `DefaultView`, since the base layout uses its fields to populate the header and footer.

Lets make an "about" page as an example:

```go
type AboutView struct {
	DefaultView
	Heading string
	Message string
}
```

Add a handler in `internal/routes`, e.g. `about.go`:

```go
package routes

import (
	"net/http"

	"github.com/osuTitanic/titanic/services/stern/internal/server"
	"github.com/osuTitanic/titanic/services/stern/internal/templates"
)

func About(ctx *server.Context) {
	view := templates.AboutView{
		DefaultView: buildDefaultView(ctx),
		Heading:     "About Titanic",
		Message:     "we're crashing big time",
	}
	ctx.RenderTemplate(http.StatusOK, "pages/public/about", view)
}
```

Register the route in `InitializeWebRoutes` in `cmd/web/main.go`. Route patterns use Go's standard `http.ServeMux` syntax and include the HTTP method.

```go
server.Handle("GET /about", routes.About)
```

Path parameters are available through `ctx.PathValue("name")`, and query parameters through `ctx.QueryValue("name")`. For an authenticated page, return when `ctx.RequireLogin()` fails:

```go
if !ctx.RequireLogin() {
	return
}
```

### Adding a Template

Stern uses the [jet template engine](https://github.com/CloudyKit/jet). Refer to the [jet syntax reference](https://github.com/CloudyKit/jet/blob/master/docs/syntax.md) for more details. The documentation is very straightforward and serves as a good reference cheat sheet.

Following our example, create a `web/template/pages/public/about.jet`.

```jet
{{ extends "/layout.jet" }}

{{ block head() }}
<title>{{ .Heading }}</title>
{{ end }}

{{ block body() }}
<div class="heading">
    <h1>{{ .Heading }}</h1>
    <p>{{ .Message }}</p>
</div>
{{ end }}
```

Page-specific styles and scripts belong in `web/static/css` and `web/static/js`. Include them from the `head` block with `cachedUrl`:

```jet
<link rel="stylesheet" href="{{ cachedUrl("/css/about.css") }}">
```

### Error Codes

Error helpers are defined in `internal/routes/errors.go`. They render the error page with the provided HTTP status, but do not stop the handler, so always return after calling one.

Use `NotFound` when a route parameter is invalid or the requested resource does not exist:

```go
userId, err := ctx.PathValueInt("id")
if err != nil {
	NotFound(ctx)
	return
}
```

For unexpected failures, log the original error before rendering the shared internal server error page:

```go
user, err := ctx.State.Repositories.Users.ById(userId)
if err != nil {
	ctx.Logger.Error("Failed to fetch user", "user_id", userId, "error", err)
	InternalServerError(ctx)
	return
}
if user == nil {
	UserNotFound(ctx)
	return
}
```

Use `RenderError` for another status code with a route-specific heading and message:

```go
RenderError(
	ctx,
	http.StatusBadRequest,
	"Invalid Request",
	"The provided user ID is invalid.",
)
```

Common errors can have their own template in `web/template/errors/custom`. `RenderErrorPage` resolves names relative to that directory:

```go
RenderErrorPage(ctx, http.StatusForbidden, "user_restricted")
return
```

Prefer an existing helper such as `UserNotFound`, `BeatmapNotFound`, `ForumNotFound`, `TopicLocked`, or `PostingTooQuickly` when it matches the failure. You may add new reusable helpers to `internal/routes/errors.go`.

### Components

Components are reusable pieces rendered as part of another template. They live in `web/template/components` and do not have their own route.

Use `include` when the component renders its data directly, without a `block`:

```jet
{{ include "/components/pagination.jet" .Pagination }}
```

If the component defines blocks, import it first and render a block with `yield`:

```jet
{{ import "/components/editor.jet" }}
{{ yield editor() .Editor }}
```

### Partials

"Partials" are standalone HTML fragments, used for lazy loading or periodically refreshing part of a page. They live in `web/template/partials` and should *not* extend `layout.jet`.

To reuse HTML from a full page, define it as a block in the page template.
The partial can then import that page and render only the block:

```jet
{* web/template/partials/home_news.jet *}
{{ import "/pages/public/home.jet" }}
{{ yield homeNews() . }}
```

Add a handler that prepares only the data needed by the partial:

```go
func HomeNewsPartial(ctx *server.Context) {
	view := any(fetchHomeNews(ctx))
	ctx.RenderTemplate(http.StatusOK, "partials/home_news", view)
}
```

Register the partial route in `InitializeWebRoutes`:

```go
server.Handle("GET /partials/home/news", routes.HomeNewsPartial)
```

Partials can then be loaded on demand with a `data-partial` attribute & page JavaScript, or refreshed with `registerPartial`:

```jet
<script defer>registerPartial("#chat", "/partials/home/chat #chat > *", 30000)</script>
```
