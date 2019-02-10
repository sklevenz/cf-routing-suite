#!/usr/bin/env bash


echo "=== test ==============================================="

# test local mongodb only if it is running
nc -z -w5 localhost 27017
if [ $? == 0 ]; then
  TAG_FLAG=-tags=mongodb
  echo "=== mongodb listening at localhost:27017"
else
  echo "=== simulation only mode"
fi

go test $TAG_FLAG ./cfrs-server/...
go test ./cfrs-client/...

