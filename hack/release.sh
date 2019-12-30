#!/bin/bash
export VERSION=0.3.1

# Tidy.
go mod tidy

# Tag & push branch for CI release.
git tag v$VERSION-preview -m "Preview release."
git push origin v$VERSION-preview

# Master.
git checkout master
git merge v$VERSION-preview
git push