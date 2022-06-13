#!/bin/bash

if [ -e ./k8s-sync/manifest/docker/k8s-sync ]
then
  echo "clear old k8s-sync"
  rm -rf ./k8s-sync/manifest/docker/k8s-sync
fi
if [ -e ./k8s-sync ]
then
  rm -rf ./k8s-sync
fi

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 /home/go/go/bin/go build -o k8s-sync main.go
# This shell is executed before docker build.


cp ./k8s-sync ./manifest/docker
#rm -rf ./k8s-sync
chmod 777 ./manifest/docker/k8s-sync

