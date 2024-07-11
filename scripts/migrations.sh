#!/bin/bash

## Change your credentials here
POSTGRES_USER="bancho"
POSTGRES_PASSWORD="examplePassword"
POSTGRES_DB="bancho"
POSTGRES_PORT="5432"
##

DATABASE_URL="postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@127.0.0.1:$POSTGRES_PORT/$POSTGRES_DB?sslmode=disable"
MIGRATIONS_DIR="../migrations"

migrate -database "$DATABASE_URL" -path "$MIGRATIONS_DIR" $@