#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Check if a version tag is provided
if [ -z "$1" ]; then
  echo "Usage: $0 <version>"
  exit 1
fi

VERSION=$1
MESSAGE=${2:-"Release $VERSION"}

# Switch to the master branch and pull the latest changes
git checkout master
git pull

# Create a new tag
git tag -a "$VERSION" -m "$MESSAGE"

# Push the tag to the remote repository
git push origin "$VERSION"

echo "Release $VERSION created and pushed successfully."