
# Manual Setup

To use [bancho](https://github.com/osuTitanic/anchor), the [api server](https://github.com/osuTitanic/deck)
and the [website](https://github.com/osuTitanic/stern), you need to have a [PostgreSQL](https://www.postgresql.org/)
and [Redis](https://redis.io/) server set up and also have [Python](https://www.python.org/) with pip installed.

For pp calculations you will also need to install the [rust toolchain](https://rustup.rs/).
You can check if its installed it by running:

```shell
cargo --version
```

### Applying migrations

After that, log in to your postgres server and apply/run this [base.sql](https://github.com/osuTitanic/titanic/blob/main/migrations/base.sql) file, that contains all the tables needed.

### Changing the default user

You also might want to change the default user `peppy` inside the `users` table, if you want to have a different username/password.
Speaking of password... you might want to use the [create_password.py](https://github.com/osuTitanic/titanic/blob/main/.github/create_password.py) file, to generate a password. *Or*, alternatively you can generate one by yourself, by hashing your password with **MD5** *and then* hash it again with **bcrypt**.

### Set up the repositories

You can now clone all of the repositories

```shell
git clone --recursive https://github.com/osuTitanic/anchor.git
```

```shell
git clone --recursive https://github.com/osuTitanic/deck.git
```

```shell
git clone --recursive https://github.com/osuTitanic/stern.git
```

After that run this command in both folders to install the dependencies:

```shell
python3 -m pip install -r requirements.txt
```

> *If pip returns an error saying that "py3rijndael" could not be found, you may need to downgrade your python version.*

And lastly you need to rename the `.example_env` files, to `.env` and edit them, to match your setup.

After that you should be ready to run the all the servers:

```shell
python3 main.py
```

You may now follow the rest of the setup [here](https://github.com/osuTitanic/titanic#adding-beatmaps).
