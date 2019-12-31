#!/bin/bash

set -e

green="\033[32m"
red="\033[31m"
reset="\033[0m"

if [ -z "$1" ]; then
  printf "${red}ERROR: Version missing${reset}\n"
  echo "Example usage: ./hack/release.sh 0.3.1"
  exit 1
fi

VERSION=$1-preview
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
sed -i'.bak' -e "s/\(VERSION=\)\(.*\)/\1$VERSION/" scripts/get-kuda-cli.sh
rm scripts/*.bak
git add scripts/get-kuda-cli.sh
git commit -m "update get script to version $VERSION"

# Tag & push branch for CI release.
git tag v$VERSION -m "Preview release."
git push origin v$VERSION

# Master.
# git checkout master
# git merge v$VERSION-preview
# git push