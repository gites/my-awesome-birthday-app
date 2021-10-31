#!/bin/bash
set -e

DOCKER_REGISTRY="${DOCKER_REGISTRY:-docker.io/gites}"

export GIT_HASH=`git rev-parse --short HEAD`

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/my-awesome-birthday-app  -ldflags="-X github.com/gites/my-awesome-birthday-app/service.Version=${GIT_HASH}"

docker build . -t ${DOCKER_REGISTRY}/my-awesome-birthday-app:${GIT_HASH}
docker push ${DOCKER_REGISTRY}/my-awesome-birthday-app:${GIT_HASH}
