#!/bin/bash

# Check if the correct number of arguments is provided
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <path_to_project> <branch_name>"
    exit 1
fi

# Assign arguments to variables
PROJECT_PATH=$1
BRANCH_NAME=$2

# Check if the project path exists
if [ ! -d "$PROJECT_PATH" ]; then
    echo "Error: Project path '$PROJECT_PATH' does not exist."
    exit 1
fi

# Change to the project directory
cd "$PROJECT_PATH" || exit

# Check if the directory is a git repository
if ! git rev-parse --is-inside-work-tree > /dev/null 2>&1; then
    echo "Error: '$PROJECT_PATH' is not a git repository."
    exit 1
fi

# Fetch the latest changes from the remote repository
git fetch origin
if [ $? -ne 0 ]; then
    echo "Error: Failed to fetch changes from the remote repository."
    exit 1
fi

# Check if the specified branch exists
if ! git show-ref --verify --quiet refs/heads/"$BRANCH_NAME"; then
    echo "Error: Branch '$BRANCH_NAME' does not exist."
    exit 1
fi

# Checkout the specified branch
git checkout "$BRANCH_NAME"
if [ $? -ne 0 ]; then
    echo "Error: Failed to checkout branch '$BRANCH_NAME'."
    exit 1
fi

# Stash any local changes
git stash
if [ $? -ne 0 ]; then
    echo "Error: Failed to stash local changes."
    exit 1
fi

# Pull the latest changes from the remote branch
git pull origin "$BRANCH_NAME"
if [ $? -ne 0 ]; then
    echo "Error: Failed to pull changes from branch '$BRANCH_NAME'."
    exit 1
fi

if [ -n "$(git stash list)" ]; then
    # Apply stashed changes
    git stash apply
    if [ $? -ne 0 ]; then
        echo "Error: Failed to apply stashed changes."
        exit 1
    fi
fi

# Update submodules, if any
git submodule update --recursive
if [ $? -ne 0 ]; then
    echo "Error: Failed to update submodules."
    exit 1
fi

# Check for merge conflicts
if [ -n "$(git status --porcelain)" ]; then
    echo "Warning: There are merge conflicts. Please resolve them manually."
else
    echo "Project updated successfully to branch '$BRANCH_NAME'."
fi