# Services

Titanic is split across several services written in Python and Go. The Python services are included as submodules, while newer services and rewrites live directly in this repository. Eventually, all services will live in this repository, rewritten in Go.

- [Bancho](bancho): Bancho game / IRC server (Python)
- [Deck](deck): Score server (Python)
- [Stern](stern): Website frontend (Go)
- [Keel](keel): API server (Python)
- [Bot](bot): Discord bot (Python)
- [Jobs](jobs): Background task runner (Go)
- [Caddy](caddy): Optional reverse proxy configuration

The [Docker configuration](../docker-compose.yml) runs the application services together with PostgreSQL, Redis, and database migrations.
