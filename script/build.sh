#!/usr/bin/env bash

echo "=== fet ================================================"

go vet ./server/...
go vet ./client/...

echo "=== build ==============================================="

rm -f ./server/cfrs-server
go build -o ./server/cfrs-server ./server
./server/cfrs-server -v

rm -f ./server/cfrs-client
go build -o ./client/cfrs-client ./client
./client/cfrs-client -v
