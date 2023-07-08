
<p align="center">
  <img width="500" alt="logo" src="https://raw.githubusercontent.com/Lekuruu/titanic/main/.github/logo/logo_medium.png">
</p>

# Titanic

Titanic is a work in progress private server for osu! stable clients from 2010-2013.
The goal of this project was to gain deeper insights into the inner workings of bancho and how it changed over the years.

You can view the progress in each repository:

- https://github.com/Lekuruu/deck
- https://github.com/Lekuruu/anchor

# Testing it out

You can test it out on my servers, using a pre-patched executable: [osu.zip](https://github.com/Lekuruu/titanic/raw/main/.github/osu.zip)

Log in with the username `Anonymous` and the password `test`.
If you want a account for yourself, please message me on discord: `lekuru`.

Please keep in mind that only one player can be online with this account!

# Quick Start

To set up and use this project I would recommend to use [docker](https://www.docker.com/).

Verify that both docker *and* docker-compose are installed:
```shell
docker --version
  Docker version X.X.X

docker-compose --version
  Docker Compose version vX.X.X
```

Clone this project onto your machine:
```shell
git clone --recursive https://github.com/Lekuruu/titanic.git
```

Rename the `.example_env` to `.env` and edit it.

You may also want to edit the `client.py` file inside the `bancho` folder, if you
have the `DISABLE_CLIENT_VERIFICATION` set to `False`.

Start the server:
```shell
docker-compose up -d
```

and hope that nothing goes wrong 😅

## Creating a user

To create a user you will need to edit the database manually, because the old clients don't support registrations
and I currently don't have a website/frontend.

Inside the `users` table, you will need to create a new row, with these attributes:

- name
- safe_name
- email
- pw (bcrypt)
- activated (true)

The password should be a bcrypt hash of a md5 hash in hex form.

## Adding beatmaps

To add beatmaps, you will *again* need to create them manually, inside the database.

**However**, I have a small collection of beatmaps and beatmapsets that you can import to your database:

- [beatmapsets.csv](https://github.com/Lekuruu/titanic/raw/main/migrations/beatmapsets.csv)
- [beatmaps.csv](https://github.com/Lekuruu/titanic/raw/main/migrations/beatmaps.csv)

They contain a total of 127226 beatmaps from 2007-2013.

## Patching the client

To actually use the client, you will need to patch it, and I would recommend using [dnSpy](https://github.com/dnSpy/dnSpy) for that.

Also, some older clients may be obfuscated.
As far as I know, [b2013606.1](https://osekai.net/snapshots/?version=179) is the latest non-obfuscated version that will work with this server.
Currently, there is support for clients from b20130716 to b1807.

You will need to find a line inside `osu.Online.BanchoClient` that looks something like this:

![unpatched](https://raw.githubusercontent.com/lekuruu/titanic/main/.github/images/unpatched.png)

and edit the ip address to match your setup:

![patched](https://raw.githubusercontent.com/lekuruu/titanic/main/.github/images/patched.png)

You also may want to use a server switcher, like [ultimate-osu-server-switcher](https://github.com/minisbett/ultimate-osu-server-switcher),
to use features such as score submission, leaderboards, etc...

**Alternatively** you can patch every url in dnSpy, from `osu.ppy.sh` to match your domain, but that can be a bit annoying.

# Contributing

If you want to clean up the mess that I made, then feel free to make a pull request.
If somebody wants to make a frontend for this project, I would be very happy.

Feel free to contact me, if you have any questions:
[@Levi/Lekuru](https://www.github.com/lekuruu)

# Screenshots

![sanic](https://raw.githubusercontent.com/lekuruu/titanic/main/.github/images/screenshot001.jpg)
![cool](https://raw.githubusercontent.com/lekuruu/titanic/main/.github/images/screenshot002.jpg)
![nice](https://raw.githubusercontent.com/lekuruu/titanic/main/.github/images/screenshot003.jpg)
![multiplayer](https://raw.githubusercontent.com/lekuruu/titanic/main/.github/images/screenshot004.jpg)
