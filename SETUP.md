
# Manual Setup

### Dependencies

- [PostgreSQL](https://www.postgresql.org/)
- [Redis](https://redis.io/)
- [Python](https://www.python.org/) with pip (3.11 is recommended)
- [Rust Toolchain](https://rustup.rs/) for pp calculations

### Applying migrations

Log in to your postgres server with your database management tool of choice, and apply/run this [base.sql](https://github.com/osuTitanic/titanic/blob/main/migrations/base.sql) file, that contains all the tables needed.

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
