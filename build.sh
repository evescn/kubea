#!/bin/bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

docker build -t docker.io/evescn/kubea:v1.$1 .
docker push docker.io/evescn/kubea:v1.$1
