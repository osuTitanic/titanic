
<p align="center">
  <img width="500" alt="logo" src="https://raw.githubusercontent.com/Lekuruu/titanic/main/.github/logo/logo_medium.png">
</p>

# Titanic

Titanic is a work in progress private server for osu! stable clients from 2008-2013.
The goal of this project was to gain deeper insights into the inner workings of bancho and how it changed over the years.

You can view the progress [here](https://github.com/users/osuTitanic/projects/2).

# Testing it out

You can test it out on my servers, using a pre-patched client: [osu.zip](https://github.com/osuTitanic/titanic/raw/main/.github/osu.zip)

I've also made a little song collection that you can download [here](https://eu2.contabostorage.com/6e40dbfbcaa94330a7e1a3f939ff105f:public/songs.zip).

For any legal issues with both of these files, please contact me on [contact@lekuru.xyz](mailto:contact@lekuru.xyz)!

Log in with the username `Anonymous` and the password `test`.
If you want a account for yourself, please message me on discord: `lekuru`.

Please keep in mind that only one player can be online with this account!

# Quick Start

To set up and use this project I would recommend to use [docker](https://www.docker.com/).

For manual setup, please view [this file](https://github.com/osuTitanic/titanic/blob/main/SETUP.md).

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

Please make sure that the folder in `bancho/app/common/` and `web/deck/app/common` is not empty!
If it is empty, then this command should fix it:

```shell
git submodule update --init
```

Rename the `.example_env` to `.env` and edit it.

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

You can generate a password, by hashing your password with **MD5** *and then* hash it again with **bcrypt**.
You can use the [create_password.py](https://github.com/osuTitanic/titanic/blob/main/tools/create_password.py) file, to do that.

## Adding beatmaps

To add beatmaps, you will *again* need to create them manually, inside the database.

**However**, I have a small collection of beatmaps and beatmapsets that you can import to your database with tools like pgAdmin:

- [beatmapsets.sql](https://github.com/osuTitanic/titanic/raw/main/migrations/beatmapsets.sql)
- [beatmaps.sql](https://github.com/osuTitanic/titanic/raw/main/migrations/beatmaps.sql)

They contain a total of 127226 beatmaps from 2007-2013.

## Patching the client

To actually use the client, you will need to patch it, and I would recommend using [dnSpy](https://github.com/dnSpy/dnSpy) for that.

Also, some older clients may be obfuscated.
As far as I know, [b2013606.1](https://osekai.net/snapshots/?version=179) is the latest non-obfuscated version that will work with this server.
Currently, there is support for clients from b2013716 to b487.

You will need to find a line inside `osu.Online.BanchoClient` that looks something like this:

![unpatched](https://raw.githubusercontent.com/osuTitanic/titanic/main/.github/images/unpatched.png)

and edit the ip address to match your setup:

![patched](https://raw.githubusercontent.com/osuTitanic/titanic/main/.github/images/patched.png)

You also may want to use a server switcher, like [ultimate-osu-server-switcher](https://github.com/minisbett/ultimate-osu-server-switcher),
to use features such as score submission, leaderboards, etc...

**Alternatively** you can patch every url in dnSpy, from `osu.ppy.sh` to match your domain, but that can be a bit annoying to do.

# Contributing

If you want to clean up the mess that I made, then feel free to make a pull request.
If somebody wants to make a frontend for this project, I would be very happy.

Also, feel free to contact me, if you have any questions:
[@Levi/Lekuru](https://www.github.com/lekuruu)

# Screenshots

![sanic](https://raw.githubusercontent.com/osuTitanic/titanic/main/.github/images/screenshot001.jpg)
![cool](https://raw.githubusercontent.com/osuTitanic/titanic/main/.github/images/screenshot002.jpg)
![nice](https://raw.githubusercontent.com/osuTitanic/titanic/main/.github/images/screenshot003.jpg)
![multiplayer](https://raw.githubusercontent.com/osuTitanic/titanic/main/.github/images/screenshot004.jpg)
