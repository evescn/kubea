#!/bin/bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

docker build -t harbor.dayuan1997.com/devops/kubea-demo:v1.17 .
docker push harbor.dayuan1997.com/devops/kubea-demo:v1.17