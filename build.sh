#!/bin/bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o k8s-syn main.go
# This shell is executed before docker build.

cp k8s-sync ./manifest/docker
chmod 777 ./manifest/docker/manifest

