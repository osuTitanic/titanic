
## Setup

Be aware that this project is typically not recommended for use on your own private server, as it is largely customized for our specific needs and can be challenging to modify. However, you are still welcome to use the project as you see fit.

To set up and use this project, it is advisable to use [Docker](https://www.docker.com/), as it is much simpler in most cases. If you do not feel comfortable using docker, here are some instructions for the [manual setup](https://github.com/osuTitanic/titanic/blob/main/MANUAL.md), which is not recommended but still possible to do.

Verify that docker is installed:

```
docker --version
```

Clone this project onto your machine:

```
git clone --recurse-submodules --shallow-submodules https://github.com/osuTitanic/titanic.git
```

Rename the `.example_env` to `.env` and **edit it**.

Start the server:

```
docker compose up -d
```

("-d" argument means detached, meaning that containers will run in background)

To turn off the server, from the titanic root folder, execute:

```
docker compose stop
```

If you experience issues on the first run, you may need to restart your containers:

```
docker compose restart
```

If you changed some files around, and don't see your changes applied, execute:

```
(rebuild)
docker compose build
(apply changes & restart affected containers)
docker compose up -d
```

After the setup is done, you should have a PostgreSQL database instance, which you can access using your database management system of choice.
By default, it contains the user `peppy` with the password `recorderinthesandybridge`.

## Updating

Titanic will get updates from time to time, so it's a good idea to apply them once in a while.

Start by first pulling all pending changes into your root folder:

```
git pull
```

After that update all of your submodules:

```
git submodule update --recursive
```

Finally, rebuild and restart all of your containers:

```
docker compose build
docker compose up -d
```

## Patching the client

You can view the instructions for patching the client [here](https://github.com/osuTitanic/clients/blob/main/PATCHING.md).