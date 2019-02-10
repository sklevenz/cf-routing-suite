#!/usr/bin/env bash


./script/test.sh

rm -f ./cfrs-server/cfrs-server
go build -o ./cfrs-server/cfrs-server ./cfrs-server/cfrs-server.go
./cfrs-server/cfrs-server -version

rm -f ./cfrs-server/cfrs-client
go build -o ./cfrs-client/cfrs-client ./cfrs-client/cfrs-client.go
./cfrs-client/cfrs-client -version
