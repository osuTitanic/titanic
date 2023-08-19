
# Manual Setup

To use [bancho](https://github.com/osuTitanic/anchor) and the [api server](https://github.com/osuTitanic/deck),
you need to have a [PostgreSQL](https://www.postgresql.org/) and [Redis](https://redis.io/) server set up and also
have [Python](https://www.python.org/) and pip installed.

### Applying migrations

After that, log in to your postgres server and apply/run this [base.sql](https://github.com/osuTitanic/titanic/blob/main/migrations/base.sql) file, that contains all the tables needed.

### Changing the default user

You also might want to change the default user `peppy` inside the `users` table, if you want to have a different username/password.
Speaking of password... you might want to use the [create_password.py](https://github.com/osuTitanic/titanic/blob/main/tools/create_password.py) file, to generate a password. *Or*, alternatively you can generate one by yourself, by hashing your password with **MD5** *and then* hash it again with **bcrypt**.

### Set up the repositories

You can now clone both of the repositories

```shell
git clone --recursive https://github.com/osuTitanic/anchor.git
```

```shell
git clone --recursive https://github.com/osuTitanic/deck.git
```

Please make sure that the folder in `app/common/` is not empty!
If it is empty, then this command should fix it:

```shell
git submodule update --init
```

After that run this command in both folders to install the dependencies:

```shell
python3 -m pip install -r requirements.txt
```

And lastly you need to rename the `.example_env` files, to `.env` and edit them, to match your setup.

After that you should be ready to run the server:

```shell
python3 main.py
```

You may now want to follow the rest of the setup [here](https://github.com/osuTitanic/titanic#adding-beatmaps).
