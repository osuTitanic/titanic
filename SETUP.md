
## Setup

Be aware that this project is typically not recommended for use on your own private server, as it is largely customized for our specific needs and can be challenging to modify. **Additionally, if you plan to publicly host this project yourself, you *must* rebrand it under a name distinct from "Titanic!"**.

To set up and use this project, it is advisable to use [Docker](https://www.docker.com/), as it is much simpler in most cases. If you do not feel comfortable using docker, here are some instructions for the [manual setup](https://github.com/osuTitanic/titanic/blob/main/MANUAL.md), which is not recommended but still possible to do.

Verify that docker is installed:

```
docker --version
```

Clone this project onto your machine:

```
git clone --recurse-submodules --shallow-submodules https://github.com/osuTitanic/titanic.git
```

Rename the `.example.env` to `.env` and **edit it**.

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

## Connecting with osu!

To connect with osu! stable you will have to set up an ssl certificate for your reverse proxy. Please look up instructions to do this online! I would personally recommend the guide from [PEACE](https://peace.osu.icu/docs/guide#2-generate-test-ssl-certificate).

If you are using a local setup environment, the easiest way to get a connection working is by editing your hosts file.  
Under Windows: `C:\Windows\System32\drivers\etc\hosts`
Under Linux: `/etc/hosts`

Add the following entries, depending on your domain name:

```
127.0.0.1 osu.bancho.local
127.0.0.1 c.bancho.local
127.0.0.1 a.bancho.local
127.0.0.1 bancho.local
127.0.0.1 ce.bancho.local
127.0.0.1 c1.bancho.local
127.0.0.1 c2.bancho.local
127.0.0.1 c3.bancho.local
127.0.0.1 c4.bancho.local
127.0.0.1 c5.bancho.local
127.0.0.1 c6.bancho.local
127.0.0.1 s.bancho.local
127.0.0.1 i.bancho.local
```

Finally, connect to your server by using the `-devserver` argument:

```shell
osu!.exe -devserver bancho.local
```

### Using older clients

The purpose of Titanic! is to be able to use older clients.
For this them to work, we made a special [patcher](https://github.com/osuTitanic/hook/releases), that will automatically change all server URLs from `ppy.sh` to your specified domain.

Simply download the patcher, put it in your osu! installation folder and run it.  
After the first run, you'll see a configuration file `Titanic!.cfg`, which you can edit to use a custom `ServerName`.
