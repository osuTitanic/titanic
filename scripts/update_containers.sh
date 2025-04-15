#!/bin/bash

# Check if the correct number of arguments is provided
if [ "$#" -lt 1 ]; then
    echo "Usage: $0 <container_name(s)>"
    exit 1
fi

# Assign arguments to variables
CONTAINER_NAMES=("$@")

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "Error: Docker is not installed."
    exit 1
fi

# Check if current directory or parent directory has a docker-compose file
if [ ! -f "docker-compose.yml" ] && [ ! -f "../docker-compose.yml" ]; then
    echo "Error: No docker-compose.yml file found in the current or parent directory."
    exit 1
fi

# Rebuild and restart the specified containers
docker compose build ${CONTAINER_NAMES[@]}
if [ $? -ne 0 ]; then
    echo "Error: Failed to build containers."
    exit 1
fi

docker compose up -d ${CONTAINER_NAMES[@]}
if [ $? -ne 0 ]; then
    echo "Error: Failed to start containers."
    exit 1
fi

echo "Containers updated and restarted successfully."