
<p align="center">
  <img width="500" alt="logo" src="https://raw.githubusercontent.com/Lekuruu/titanic/main/.github/logo/logo_medium.png">
</p>

# Titanic

Titanic is a work in progress private server for osu! stable clients from 2008-2013.
The goal of this project was to gain deeper insights into the inner workings of bancho and how it changed over the years.

You can view the progress [here](https://github.com/users/osuTitanic/projects/2).

# Testing it out

For anyone interested, we got a discord server: https://discord.gg/qupv72e7YH

You can test it out on my servers using a pre-patched client, which you can find [here](https://github.com/osuTitanic/clients). Log in with the username `Anonymous` and the password `test`.
If you want a account for yourself, please join the discord server!

Please keep in mind that only one player can be online with this account!

# Quick Start

To set up and use this project I would recommend to use [docker](https://www.docker.com/). Otherwise here are some instructions for [manual setup](https://github.com/osuTitanic/titanic/blob/main/SETUP.md).

Verify that both docker *and* docker-compose are installed:

```shell
docker --version
  Docker version X.X.X

docker-compose --version
  Docker Compose version vX.X.X
```

Clone this project onto your machine:

```shell
git clone  --recurse-submodules --shallow-submodules https://github.com/Lekuruu/titanic.git
```

**Please make sure** that the folder in `bancho/app/common/` and `web/deck/app/common` is not empty!
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
You can use the [create_password.py](https://github.com/osuTitanic/titanic/blob/main/.github/create_password.py) file, to do that.

## Adding beatmaps

To add beatmaps, you will *again* need to create them manually, inside the database.

**However**, I have a small collection of beatmaps and beatmapsets that you can import to your database with tools like pgAdmin:

- [beatmapsets.sql](https://github.com/osuTitanic/titanic/raw/main/migrations/beatmapsets.sql)
- [beatmaps.sql](https://github.com/osuTitanic/titanic/raw/main/migrations/beatmaps.sql)

They contain a total of 140291 beatmaps from 2007-2013.

## Patching the client

You can view the instructions [here](https://github.com/osuTitanic/clients/blob/main/PATCHING.md).

# Credits

I would like to thank...

- [kanaarima](https://github.com/kanaarima/) for helping with the discord bot
- [rory](https://github.com/TheArcaneBrony) for [funny playtesting](https://raw.githubusercontent.com/osuTitanic/titanic/main/.github/images/screenshot022.jpg)
- [osu.direct](https://osu.direct/) & [nerinyan.moe](https://nerinyan.moe/) for providing beatmap resources

# Contributing

If you want to clean up the mess that I made, then feel free to fork and make a pull request.
I am also working on frontend, which will be available soon™ and I would appreciate any help or suggestions.

Contact me on the discord server, if you have any questions.

# Screenshots

#### b20130606.1

![sanic](https://raw.githubusercontent.com/osuTitanic/titanic/main/.github/images/screenshot007.jpg)
![cool](https://raw.githubusercontent.com/osuTitanic/titanic/main/.github/images/screenshot008.jpg)

#### b20130303

![wow](https://raw.githubusercontent.com/osuTitanic/titanic/main/.github/images/screenshot023.jpg)

#### b1700

![nice](https://raw.githubusercontent.com/osuTitanic/titanic/main/.github/images/screenshot005.jpg)
![multiplayer](https://raw.githubusercontent.com/osuTitanic/titanic/main/.github/images/screenshot006.jpg)
