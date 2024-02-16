
<p align="center">
  <img width="500" alt="logo" src="https://raw.githubusercontent.com/Lekuruu/titanic/main/.github/logo/logo_medium.png">
</p>

# Titanic

Titanic is a work in progress private server for all osu! stable clients (2008-2024).
The goal of this project was to gain deeper insights into the inner workings of bancho and how it changed over the years.

The main goal of this project was achieved. There are still some features that I want to add, which you can view [here](https://github.com/users/osuTitanic/projects/2).

# Testing it out

You can test it out on my servers, by [registering on the website](https://osu.lekuru.xyz/account/register), and [downloading](https://osu.lekuru.xyz/download) the client from the website.
**Please note that it will only work with the clients downloaded from the website.**

For anyone interested, we also got a discord server: https://discord.gg/qupv72e7YH

# Server Setup

To set up and use this project I would recommend to use [docker](https://www.docker.com/). Otherwise here are some instructions for [manual setup](https://github.com/osuTitanic/titanic/blob/main/SETUP.md), which is not recommended but still possible to do.

Verify that docker is installed:

```shell
docker --version
  Docker version X.X.X, build ...
```

Clone this project onto your machine:

```shell
git clone --recurse-submodules --shallow-submodules https://github.com/osuTitanic/titanic.git
```

Rename the `.example_env` to `.env` and **edit it**.

Start the server:

```shell
docker compose up -d
```

and hope that nothing goes wrong 😅

## Creating a user

You can create users by simply registering on the website. New users will appear in the `users` table, inside the database.

## Adding beatmaps

To add beatmaps, you will need to create them manually, inside the database.

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
- Everyone that donated & contributed to this project (you guys are a big help!)

# Contributing

If you want to clean up the mess that I made, then feel free to fork and make a pull request.

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
