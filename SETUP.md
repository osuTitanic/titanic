
# Manual Setup

### Dependencies

- [PostgreSQL](https://www.postgresql.org/)
- [Redis](https://redis.io/)
- [Python](https://www.python.org/) with pip (3.11 is recommended)
- [Rust Toolchain](https://rustup.rs/) for pp calculations

### Applying migrations

Follow the installation guide for [golang-migrate](https://github.com/golang-migrate/migrate/blob/master/cmd/migrate/README.md#installation).

Please open the `migrations.sh` file in the `scripts` folder, and change your database credentials.
After that run the script:

```shell
> ./migrations.sh up

(output)
0/u BaseTables (5.471535524s)
1/u Users (5.660987277s)
2/u Groups (5.861020717s)
3/u Channels (6.114019534s)
4/u VerifiedClients (6.324159415s)
5/u Forums (6.576804694s)
6/u Mirrors (6.787458535s)
7/u Indexes (7.797943227s)
8/u BeatmapOffset (8.029588468s)
9/u Beatmapsets (10.478337961s)
10/u Beatmaps (17.158129164s)
```

If you run this project on windows, you may need to enter the command manually:

```shell
migrate -database "postgres://<USER>:<PASSWORD>@<HOST>:<PORT>/<DATABASE>?sslmode=disable" -path ../migrations up
```

### Set up the repositories

You can now clone all of the repositories:

```shell
git clone --recursive https://github.com/osuTitanic/anchor.git
```

```shell
git clone --recursive https://github.com/osuTitanic/deck.git
```

```shell
git clone --recursive https://github.com/osuTitanic/stern.git
```

After that install the dependencies for each project with:

```shell
python3 -m pip install -r requirements.txt
```

You might want to use a [virtual environment](https://docs.python.org/3/tutorial/venv.html) for that.

Rename the `.example_env` files, to `.env` and edit them, to match your setup.
After that you should be ready to run the all the servers:

```shell
python3 main.py
```

You may now want to follow the rest of the setup [here](https://github.com/osuTitanic/titanic#adding-beatmaps).
