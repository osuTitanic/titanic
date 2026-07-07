<p align="center">
  <img width="300" alt="logo" src="https://raw.githubusercontent.com/Lekuruu/titanic/main/.github/images/logo-vector.min.svg">
</p>

# Titanic! (Go Rewrite)

This repository contains the work-in-progress Go rewrite of the [Titanic!](https://github.com/osuTitanic/titanic) services originally written in python.

## Services

It is planned to port every service from the original python codebase, with *Jobs* (background tasks) and *Stern* (frontend) being the first to be rewritten. *Anchor* (bancho game server), *Deck* (score server), *Keel* (API) and *BanchoBot* (discord bot) will follow sooner or later.

Most shared application code lives in [State](internal/state/README.md). Services should generally start with `state.NewState(".env")` instead of creating config, database, redis, storage and repositories manually.
For more focused examples of what the individual shared components do, see the package-level READMEs inside `internal/`.

## Development

To set up the development environment, you will (of course) need to have [Go](https://go.dev/doc/install) installed. A configured postgresql & redis instance is also required to run the services. For these, follow the setup instructions in the [titanic repository](https://github.com/osuTitanic/titanic/blob/main/SETUP.md).

Build the current service binaries from the **root** repository folder:

```sh
go build ./services/stern/cmd/web
go build ./services/jobs/cmd/cli
```

```sh
go run ./services/stern/cmd/web
go run ./services/jobs/cmd/cli
```

## Docker

Service images are also built from the repository root.

```sh
docker build -f ./services/stern/Dockerfile -t osutitanic/stern .
docker build -f ./services/jobs/Dockerfile -t osutitanic/jobs .
```

The images only contain the compiled service binaries to keep the image size small. I'd expect ~20 to 30 MB for them.
Runtime configuration still comes from environment variables, and persistent data should be mounted at the configured data path.

Pre-built images are also available on ghcr.io:

- [ghcr.io/osutitanic/stern](https://ghcr.io/osutitanic/stern)
- [ghcr.io/osutitanic/jobs](https://ghcr.io/osutitanic/jobs)
