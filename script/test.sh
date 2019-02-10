#!/usr/bin/env bash


echo "=== fet ================================================"

go vet ./cfrs-server/...
go vet ./cfrs-client/...

echo "=== test ==============================================="

export GOCACHE=off

# test local mongodb
nc -z -w5 localhost 27017
if [ $? == 0 ]; then
  TAG_FLAG=-tags=mongodb
  echo "=== mongodb listening at localhost:27017"
else
  echo "=== simulation only mode"
fi

go test  $TAG_FLAG ./cfrs-server/...
go test ./cfrs-client/...

