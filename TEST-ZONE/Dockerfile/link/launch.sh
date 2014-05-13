#!/bin/bash

cd postgres ; docker build -t postgres .
cd node ; docker build -t node .

timeout 20 docker run -d  -p 5432:5432 postgres &
docker run -rm -p 8080:8080 node

