#!/bin/bash

docker stop $(docker ps -a -q)
sleep 2
docker rm $(docker ps -a -q)
