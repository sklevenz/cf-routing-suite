#!/usr/bin/env bash

echo "=== fet ================================================"

go vet ./cfrs-server/...
go vet ./cfrs-client/...

./script/test.sh

echo "=== build ==============================================="

rm -f ./cfrs-server/cfrs-server
go build -o ./cfrs-server/cfrs-server ./cfrs-server/cfrs-server.go
./cfrs-server/cfrs-server -v

rm -f ./cfrs-server/cfrs-client
go build -o ./cfrs-client/cfrs-client ./cfrs-client/cfrs-client.go
./cfrs-client/cfrs-client -v
