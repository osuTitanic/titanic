SHELL := /usr/bin/env bash

.DEFAULT_GOAL := help
.PHONY: help \
	up down stop restart build rebuild pull ps logs db-psql redis-cli

COMPOSE ?= docker compose
COMPOSE_FILE ?= docker-compose.yml
COMPOSE_FILE_OVERRIDE ?= docker-compose.caddy.yml

DC = $(COMPOSE) -f $(COMPOSE_FILE) -f $(COMPOSE_FILE_OVERRIDE)

SERVICE ?=
CMD ?= /bin/bash

help: ## Show available commands
	@awk 'BEGIN {FS = ":.*##"}; /^[a-zA-Z0-9_-]+:.*##/ {printf "\033[36m%-18s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

up: ## Start the stack (detached)
	@$(DC) up -d $(SERVICE)

down: ## Stop and remove containers/networks
	@$(DC) down

stop: ## Stop running containers
	@$(DC) stop $(SERVICE)

restart: ## Restart containers
	@$(DC) restart $(SERVICE)

build: ## Build images
	@$(DC) build $(SERVICE)

rebuild: ## Rebuild images and restart
	@$(DC) build $(SERVICE)
	@$(DC) up -d $(SERVICE)

pull: ## Pull upstream images
	@$(DC) pull $(SERVICE)

ps: ## Show container status
	@$(DC) ps

config: ## Render the compose config
	@$(DC) config

logs: ## Tail logs
	@$(DC) logs -f --tail=$(TAIL) $(SERVICE)

logs-tail: ## Show recent logs
	@$(DC) logs --tail=$(TAIL) $(SERVICE)

shell: ## Open a shell/command in SERVICE, e.g. make shell SERVICE=stern
	@test -n "$(SERVICE)" || (echo "SERVICE is required, e.g. make shell SERVICE=stern"; exit 1)
	@$(DC) exec $(SERVICE) $(CMD)

db-psql: ## Open psql in the postgres container
	@$(DC) exec db sh -c 'PGPASSWORD="$$POSTGRES_PASSWORD" psql -U "$$POSTGRES_USER" -d "$${POSTGRES_DATABASE:-$$POSTGRES_USER}"'

redis-cli: ## Open redis-cli in the cache container
	@$(DC) exec cache redis-cli

migrate-up: ## Run database migrations
	@$(DC) up migrations

update: ## Pull repo changes and update submodules
	@git pull
	@git submodule update --init --recursive

submodules: ## Initialize or update submodules
	@git submodule update --init --recursive

fmt: ## Format Go packages
	@go fmt ./...

vet: ## Run go vet
	@go vet ./...

test: ## Run Go tests, override with TEST=./internal/...
	@go test $(TEST)

test-integration: ## Run Go integration tests
	@go test -tags=integration -count=1 -p 1 $(TEST)

check: fmt vet test ## Format, vet, and test Go code

tidy: ## Tidy Go modules
	@go mod tidy

clean: ## Remove local build/test cache
	@go clean -cache -testcache
