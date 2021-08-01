#! /usr/bin/env bash

set -euo pipefail

IMAGE_MANIFEST_TAG=${1:-batect-cache-init-image:local-testing}

./scripts/build_image.sh "$IMAGE_MANIFEST_TAG" "amd64"
echo

echo "Running tests..."
IMAGE_MANIFEST_TAG=$IMAGE_MANIFEST_TAG go test .

