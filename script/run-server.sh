#!/usr/bin/env bash

USAGE="Usage: run-server.sh [s|db] [8080]"

if [[ $# -eq 0 ]] ; then
    echo $USAGE
    exit 0
fi

case "$1" in
    "s") export MODE="simulator" ;;
    "db") export MODE="mongodb" ;;
    *) echo $USAGE; exit 0 ;;
esac

VERSION="snapshot-$(git log --format="%H" -n 1)"
if [[ $(git diff --stat) != '' ]]; then
  VERSION="$VERSION-dirty"
fi

if [[ $2 -eq 0 ]] ; then
  export PORT=8080
else
  export PORT=$2
fi

echo "=== port: $PORT";
echo "=== mode: $MODE";


pushd server
    go run -ldflags="-X main.version=$VERSION" ./
popd

