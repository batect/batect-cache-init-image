#! /usr/bin/env bash

set -euo pipefail

IMAGE_TAG=${1:-batect-cache-init-image}

./scripts/build.sh "$IMAGE_TAG"
echo

echo "Running tests..."
go test .

