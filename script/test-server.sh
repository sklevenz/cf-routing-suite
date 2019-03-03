#!/usr/bin/env bash


echo "=== test server ==============================================="

# test local mongodb only if it is running
nc -z -w5 localhost 27017 >/dev/null 2>&1
if [ $? == 0 ]; then
  TAG_FLAG="-tags=mongodb"
fi

go test $TAG_FLAG ./server/... $@

