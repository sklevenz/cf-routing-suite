#!/usr/bin/env bash

pushd cfrs-server
    go run ./cfrs-server.go $1
popd