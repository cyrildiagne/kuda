#!/bin/bash

if [ -z "$1" ]; then
  echo "ERROR: Version missing"
  echo "Example usage: ./hack/release.sh 0.3.1"
  exit 1
fi

VERSION=$1
echo "Releasing v$VERSION..."

# Tidy.
go mod tidy

# Tag & push branch for CI release.
git tag v$VERSION-preview -m "Preview release."
git push origin v$VERSION-preview

# Master.
# git checkout master
# git merge v$VERSION-preview
# git push