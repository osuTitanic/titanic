<p align="center">
  <img width="300" alt="logo" src="https://raw.githubusercontent.com/Lekuruu/titanic/main/.github/images/logo-vector.min.svg">
</p>

# Titanic

Titanic is an osu! private server designed to be compatible with every osu! stable client out there (2007-2026).
The goal of this project was to gain deeper insights into the inner workings of Bancho and how it has evolved over the years.

You can play on it, by [registering on our website](https://osu.titanic.sh/account/register), and [downloading](https://osu.titanic.sh/download/) a client. You may also use clients directly from [Osekai Snapshots](https://osekai.net/snapshots/) with the Titanic! [Patcher](https://osu.titanic.sh/wiki/en/Patcher).
For more questions, feel free to join our [Discord server](https://discord.gg/qupv72e7YH).

## Services

This repository combines the docker deployment, database migrations, Python service submodules, and Go service rewrites. This makes the project a bit of a weird amalgamation of different languages. The goal is to have everything in Go eventually, but that is going to take some time since I am not use the "make no mistakes" approach here.

- `services/bancho`: Bancho game server (Python)
- `services/deck`: Score server (Python)
- `services/stern`: Website/Frontend (Go)
- `services/keel`: API (Python)
- `services/bot`: Discord bot (Python)
- `services/jobs`: Background task runner (Go)

## Contributing

You are very welcome to make any kinds of suggestions or contributions to this project. However, please note the unusual split between Python and Go services, that I have mentioned above. For large-scale changes you want to make, I'd suggest to contact me for further help and coordination.

## Screenshots

![image](https://raw.githubusercontent.com/osuTitanic/titanic/main/.github/images/screenshot001.png)
![image](https://raw.githubusercontent.com/osuTitanic/titanic/main/.github/images/screenshot002.png)
![image](https://raw.githubusercontent.com/osuTitanic/titanic/main/.github/images/screenshot003.png)
![image](https://raw.githubusercontent.com/osuTitanic/titanic/main/.github/images/screenshot004.png)
![image](https://raw.githubusercontent.com/osuTitanic/titanic/main/.github/images/screenshot005.png)
![image](https://raw.githubusercontent.com/osuTitanic/titanic/main/.github/images/screenshot006.png)
![image](https://raw.githubusercontent.com/osuTitanic/titanic/main/.github/images/screenshot007.png)
![image](https://raw.githubusercontent.com/osuTitanic/titanic/main/.github/images/screenshot008.png)
