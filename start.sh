#!/bin/bash
set -eu
version=$1
GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -installsuffix cgo -o bin/gok8s .
# GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o bin/gok8s .
docker build -t gok8s:${version:-latest} .