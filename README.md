<p align="center">
  <img width="350" alt="logo" src="https://raw.githubusercontent.com/Lekuruu/titanic/main/.github/images/logo-vector.min.svg">
</p>

# Titanic

Titanic is an osu! private server designed to be compatible with every osu! stable client out there (2007-2025).
The goal of this project was to gain deeper insights into the inner workings of Bancho and how it has evolved over the years.

You can play on it, by [registering on our website](https://osu.titanic.sh/account/register), and [downloading](https://osu.titanic.sh/download/) a client.
For more questions, feel free to join our Discord server: https://discord.gg/qupv72e7YH

## Services

This repository combines the main deployment, database migrations, Python service submodules, and Go service rewrites.

- `services/bancho`: Bancho game server, from the `anchor` submodule.
- `services/deck`: Score server submodule.
- `services/bot`: Discord bot submodule.
- `services/keel`: API submodule.
- `services/stern`: Go website/frontend.
- `services/jobs`: Go background task runner.

Most shared Go application code lives in [State](internal/state/README.md). Go services should generally start with `state.NewState(".env")` instead of creating config, database, Redis, storage, and repositories manually.

## Development

Docker remains the recommended way to run the full stack. See [SETUP.md](SETUP.md) for the normal Docker setup and [MANUAL.md](MANUAL.md) for manual service startup notes.

To work on the Go services, install [Go](https://go.dev/doc/install), configure PostgreSQL and Redis through the root `.env`, and run commands from the repository root:

```sh
go test ./...
go build ./services/stern/cmd/web
go build ./services/jobs/cmd/cli
```

```sh
go run ./services/stern/cmd/web
go run ./services/jobs/cmd/cli
```

Service images are also built from the repository root:

```sh
docker build -f ./services/stern/Dockerfile -t osutitanic/stern .
docker build -f ./services/jobs/Dockerfile -t osutitanic/jobs .
```

## Contributing

You are very welcome to make any kinds of suggestions or contributions to this project.
I would suggest starting by familiarizing yourself with the different kinds of external resources and APIs, that the codebase offers.
I have written a document about it [here](https://github.com/osuTitanic/common/blob/main/USAGE.md).
For large-scale changes you want to make, I'd suggest to contact me for further help and coordination.

## Special Thanks

- Everyone that is playing, donating & contributing to this project
- All early [Testers](https://osu.titanic.sh/g/8) that somehow found this project through GitHub and decided to register through Discord DMs
- [Adachi](https://osu.titanic.sh/u/39), [BlueChinchompa](https://osu.titanic.sh/u/40) & [Meru](https://osu.titanic.sh/u/41) for community management
- Beatmap mirrors, such as [osu.direct](https://osu.direct/), [nerinyan](https://nerinyan.moe/) & [mino](https://catboy.best)

## Screenshots

![image](https://raw.githubusercontent.com/osuTitanic/titanic/main/.github/images/screenshot001.png)
![image](https://raw.githubusercontent.com/osuTitanic/titanic/main/.github/images/screenshot002.png)
![image](https://raw.githubusercontent.com/osuTitanic/titanic/main/.github/images/screenshot003.png)
![image](https://raw.githubusercontent.com/osuTitanic/titanic/main/.github/images/screenshot004.png)
![image](https://raw.githubusercontent.com/osuTitanic/titanic/main/.github/images/screenshot005.png)
![image](https://raw.githubusercontent.com/osuTitanic/titanic/main/.github/images/screenshot006.png)
![image](https://raw.githubusercontent.com/osuTitanic/titanic/main/.github/images/screenshot007.png)
![image](https://raw.githubusercontent.com/osuTitanic/titanic/main/.github/images/screenshot008.png)
