#! /usr/bin/env bash

set -euo pipefail

IMAGE_TAG=$1
ARCHITECTURE=$2

echo "Building application..."
CGO_ENABLED=0 GOOS=linux GOARCH="$ARCHITECTURE" go build -o docker/app .
echo

echo "Building image..."
docker build -t "$IMAGE_TAG" --platform "linux/$ARCHITECTURE" docker
echo
