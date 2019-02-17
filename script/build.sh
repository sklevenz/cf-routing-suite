#!/usr/bin/env bash

echo "=== fet ================================================"

go vet ./server/...
go vet ./client/...

./script/test.sh

echo "=== build ==============================================="

rm -f ./server/cfrs-server
go build -o ./server/cfrs-server ./server/cfrs-server.go
./server/cfrs-server -v

rm -f ./server/cfrs-client
go build -o ./client/cfrs-client ./client/cfrs-client.go
./client/cfrs-client -v
