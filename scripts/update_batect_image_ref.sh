#! /usr/bin/env bash

set -euo pipefail

SCRIPT_PATH="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
IMAGE_DIGEST=$(docker inspect "$IMAGE_TAG" --format '{{ index .RepoDigests 0 }}')
WORK_DIR=$(mktemp -d)
FILE="app/src/main/kotlin/batect/execution/CacheInitialisationImage.kt"

{
  echo "Cloning repo..."
  git clone https://$GITHUB_USERNAME:$GITHUB_TOKEN@github.com/batect/batect.git "$WORK_DIR"
  echo

  cd "$WORK_DIR"

  echo "Configuring Git..."
  git config user.name "batect-cache-init-image pipeline"
  git config user.email "batect-github-actions@batect.dev"
  echo

  echo "Setting image reference to '$IMAGE_DIGEST'..."
  mkdir -p "$(dirname "$FILE")"
  sed "s#REPLACE_WITH_IMAGE_TAG#$IMAGE_DIGEST#" "$SCRIPT_PATH/Template.kt" > "$FILE"
  echo

  echo "Preparing commit..."
  git add "$FILE"
  echo

  echo "Committing..."
  git commit -m "Update reference to cache init image."
  echo

  echo "Pushing..."
  git show
  #  git push
  echo
}

echo "Cleaning up..."
rm -rf "$WORK_DIR"
echo

echo "Done!"
