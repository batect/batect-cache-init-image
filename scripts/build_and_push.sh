#! /usr/bin/env bash

set -euo pipefail

IMAGE_MANIFEST_TAG=$1
ARCHITECTURES=(amd64 arm64 arm)

function main() {
  for architecture in "${ARCHITECTURES[@]}"; do
    buildForArchitecture "$architecture"
  done

  for architecture in "${ARCHITECTURES[@]}"; do
    pushImage "$architecture"
  done

  createManifest
  pushManifest

  echo
  echoHeader "Done."
}

function buildForArchitecture() {
  local architecture=$1

  echoHeader "Building for $architecture..."
  ./scripts/build_image.sh "$IMAGE_MANIFEST_TAG-$architecture" "$architecture"
}

function pushImage() {
  local architecture=$1

  echoHeader "Pushing image for $architecture..."
  docker push "$IMAGE_MANIFEST_TAG-$architecture"
  echo
}

function createManifest() {
  echoHeader "Creating manifest..."
  docker manifest create "$IMAGE_MANIFEST_TAG" "${ARCHITECTURES[@]/#/$IMAGE_MANIFEST_TAG-}"
  echo
}

function pushManifest() {
  echoHeader "Pushing manifest..."
  docker manifest push "$IMAGE_MANIFEST_TAG"
  echo
}

function echoHeader() {
  local text=$1

    echo "-------------------------------------------"
    echo "$text"
    echo "-------------------------------------------"
}

main
