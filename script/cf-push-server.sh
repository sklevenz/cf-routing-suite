#!/usr/bin/env bash

VERSION="snapshot-$(git log --format="%H" -n 1)"
USAGE="Usage: cf-push-server.sh [s|db]"

if [[ $# -eq 0 ]] ; then
    echo $USAGE
    exit 0
fi

case "$1" in
    "s") MODE="simulator" ;;
    "db") MODE="mongodb" ;;
    *) echo $USAGE; exit 0 ;;
esac

if [[ $(git diff --stat) != '' ]]; then
  VERSION="$VERSION-dirty"
fi

echo "=== mode: $MODE";
echo "=== version: $VERSION";

pushd cfrs-server
  mkdir -p gen
  sed -e "s/{{VERSION}}/$VERSION/" -e "s/{{MODE}}/$MODE/" manifest.yml.template > ./gen/manifest.yml
  cf push -f ./gen/manifest.yml
popd
