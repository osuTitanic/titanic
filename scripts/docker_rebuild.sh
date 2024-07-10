#!/bin/bash

echo "Rebuilding docker containers..."
docker compose build $@
docker compose up -d $@