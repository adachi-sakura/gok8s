#!/bin/bash
version=$1
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -installsuffix cgo -o gok8s .
docker build -t gok8s:${version:-latest} .