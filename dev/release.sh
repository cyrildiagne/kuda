#!/bin/bash

set -e

green="\033[32m"
red="\033[31m"
reset="\033[0m"

# Ensure version was given as arg.
if [ -z "$1" ]; then
  printf "${red}ERROR:${reset} Version missing\n"
  echo "Example usage: ./scripts/release.sh 0.3.1"
  exit 1
fi

VERSION=$1-preview

# Check if current branch name matches with given version
BRANCH_NAME=$(git rev-parse --abbrev-ref HEAD)

if [[ "$VERSION" != "$BRANCH_NAME"* ]]; then
  printf "${red}ERROR:${reset} Git branch version mismatch.\n"
  echo "Version: $VERSION should start with Git branch name: $BRANCH_NAME"
  exit 1
fi

printf "${green}Releasing v$VERSION...${reset}\n"

# Tidy.
go mod tidy

# Check git state.
if [[ $(git diff --stat) != '' ]]; then
  printf "${red}ERROR: Git state is dirty. Commit changes and try again:${reset}\n"
  git status --short
  exit 1
fi

# Update scripts/get-kuda-cli.sh
sed -i'.bak' -e "s/\(VERSION=\)\(.*\)/\1$VERSION/" scripts/get-cli.sh
rm scripts/*.bak
git add scripts/get-cli.sh

if git diff --exit-code; then
  echo "nothing to commit"
else
  git commit -m "update get script to version $VERSION"
fi

# Tag & push branch for CI release.
git tag v$VERSION -m "Preview release."
git push origin v$VERSION

# Master.
# git checkout master
# git merge v$VERSION-preview
# git push