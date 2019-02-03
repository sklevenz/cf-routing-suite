#!/usr/bin/env bash

pushd cfrs-server
    go run ./main.go $1
popd