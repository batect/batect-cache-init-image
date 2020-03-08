#! /usr/bin/env bash

set -euo pipefail

IMAGE_TAG=${1:-batect-cache-init-image}

echo "Building application..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o docker/app .
echo

echo "Building image..."
docker build -t $IMAGE_TAG docker
