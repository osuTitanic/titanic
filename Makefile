SHELL := /usr/bin/env bash

.DEFAULT_GOAL := help
.PHONY: help \
	up down stop restart build rebuild pull ps logs db-psql redis-cli

COMPOSE ?= docker compose
COMPOSE_FILE ?= docker-compose.yml

DC = $(COMPOSE) -f $(COMPOSE_FILE)

SERVICE ?=
CMD ?= /bin/bash

help: ## Show available commands
	@awk 'BEGIN {FS = ":.*##"}; /^[a-zA-Z0-9_-]+:.*##/ {printf "\033[36m%-18s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

up: ## Start the stack (detached)
	@$(DC) up -d $(SERVICE)

down: ## Stop and remove containers/networks
	@$(DC) down $(SERVICE)

stop: ## Stop running containers
	@$(DC) stop $(SERVICE)

restart: ## Restart containers
	@$(DC) restart $(SERVICE)

build: ## Build images
	@$(DC) build $(SERVICE)

pull: ## Pull upstream images
	@$(DC) pull

ps: ## Show container status
	@$(DC) ps

logs: ## Tail logs (all services)
	@$(DC) logs -f --tail=200 $(SERVICE)

rebuild: ## Rebuild images and restart
	@$(DC) build $(SERVICE)
	@$(DC) up -d $(SERVICE)

db-psql: ## Open psql in postgres container
	-@$(DC) exec db sh -c 'PGPASSWORD="$$POSTGRES_PASSWORD" psql -U "$$POSTGRES_USER" -d "$$POSTGRES_DB"'

redis-cli: ## Open redis-cli in cache container
	@$(DC) exec cache redis-cli