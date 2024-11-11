
# Manual Setup

## Dependencies

- [PostgreSQL](https://www.postgresql.org/)
- [Redis](https://redis.io/)
- [Python](https://www.python.org/) with pip (3.11 is recommended)
- [Rust Toolchain](https://rustup.rs/) for pp calculations

## Cloning the main repository

Just like in the docker setup, clone the main repository onto your machine:

```
git clone --recurse-submodules --shallow-submodules https://github.com/osuTitanic/titanic.git
```

## Applying database migrations

Follow the installation guide for [golang-migrate](https://github.com/golang-migrate/migrate/blob/master/cmd/migrate/README.md#installation).

Please open the `migrations.sh` file in the `scripts` folder, and change your database credentials.
After that run the script:

```shell
> ./migrations.sh up

(output)
0/u PgCrypto (236.668784ms)
1/u Users (1.775676093s)
2/u Forums (2.650836088s)
3/u Beatmaps (4.20043261s)
4/u Rankings (5.085522413s)
5/u Bancho (6.024367129s)
6/u Profile (6.920077113s)
7/u Clients (7.4258269s)
8/u BeatmapsetData (9.9139888s)
9/u BeatmapData (16.864791013s)
```

If you run this project on windows, you may need to enter the command manually:

```shell
migrate -database "postgres://<USER>:<PASSWORD>@<HOST>:<PORT>/<DATABASE>?sslmode=disable" -path ../migrations up
```

## Run the project

Here is a list of all projects you want to run:

- `/bancho` (Bancho Server)
- `/web/deck` (Score Server)
- `/web/deck` (Website)

First, install the dependencies for each project with:

```shell
python3 -m pip install -r requirements.txt
```

You might want to use a [virtual environment](https://docs.python.org/3/tutorial/venv.html) for that, if any dependencies conflict with each other.

Rename the `.example.env` files, to `.env` and edit them, to match your setup.
After that you should be ready to run the all the servers, by running this command for all projects:

```shell
python3 main.py
```

You may now want to patch a client, using [this guide](https://github.com/osuTitanic/clients/blob/main/PATCHING.md). It is also recommended to run the website, as well as the score server behind a reverse proxy, so that both servers can be accessed under `osu.yourdomain.com`. However this guide won't contain any instructions for that.