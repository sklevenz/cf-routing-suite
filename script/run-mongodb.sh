#!/usr/bin/env bash

SOURCE="${0%/*}"

mkdir -p $SOURCE/../gen/mongo/db

mongod --dbpath $SOURCE/../gen/mongo/db --logpath $SOURCE/../gen/mongo/db.log --port 27017

