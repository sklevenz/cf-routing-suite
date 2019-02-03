#!/usr/bin/env bash

set -e

SOURCE="${0%/*}"
VERSION=$1

if [ "$VERSION" == "" ]; then
    echo Usage: delete-tag VERSION
    exit 0
fi

git push --delete origin "$VERSION"
git tag -d "$VERSION"
