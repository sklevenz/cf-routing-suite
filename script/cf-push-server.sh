#!/usr/bin/env bash

VERSION=$(git log --format="%H" -n 1)

if [[ $(git diff --stat) != '' ]]; then
  VERSION="$VERSION-dirty"
fi


pushd cfrs-server
  mkdir -p gen
  sed -e "s/{{VERSION}}/$VERSION/" manifest.yml.template > ./gen/manifest.yml

  cf push -f ./gen/manifest.yml

popd
