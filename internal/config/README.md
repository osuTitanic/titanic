# Configuration

This module loads the app configuration from environment variables and dotenv files, using [caarlos0/env](https://github.com/caarlos0/env) & [joho/godotenv](https://github.com/joho/godotenv).

## Usage

Load the default `.env` file and parse the current process environment into a `Config`.

```go
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
}
```

You can pass explicit dotenv filepaths when a service needs a different environment file.
Files are loaded in order before parsing the environment.

```go
cfg, err := config.LoadConfig(".env", ".env.local")
if err != nil {
	return err
}
```

## Helpers

Make use of available helper methods, e.g.

- `PostgresDSN()` for the PostgreSQL connection string
- `OsuBaseUrl()`, `ApiBaseUrl()`, ... for public service URLs
- `GetAllowInsecureCookies()` for cookie security decisions
- `ValidImageServices()` for a list of image hosts that can bypass the proxy
- ...

## Custom Types

`StringSlice` and `IntSlice` accept either JSON arrays or comma-separated values:

```env
AUTOJOIN_CHANNELS=["#osu", "#announce"]
BANCHO_TCP_PORTS=13380,13381,13382,13383
```

`DynamicTime` accepts RFC3339 timestamps, datetime strings, and date-only strings:

```env
BEGINNING_ENDED_AT=2023-12-31T06:00:00Z
```
