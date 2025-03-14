#!/bin/bash

# This script bumps the version in the VERSION file and updates chadburn.go
# Usage: ./scripts/bump-version.sh [major|minor|patch]

set -e

# Create scripts directory if it doesn't exist
mkdir -p $(dirname "$0")

# Check if the version type is provided
if [ $# -ne 1 ]; then
    echo "Usage: $0 [major|minor|patch]"
    exit 1
fi

TYPE=$1

# Check if the version type is valid
if [ "$TYPE" != "major" ] && [ "$TYPE" != "minor" ] && [ "$TYPE" != "patch" ]; then
    echo "Invalid version type. Use 'major', 'minor', or 'patch'."
    exit 1
fi

# Read the current version
if [ -f VERSION ]; then
    CURRENT_VERSION=$(cat VERSION)
else
    CURRENT_VERSION="0.0.0"
    echo $CURRENT_VERSION > VERSION
fi

# Split the version into components
IFS='.' read -r -a VERSION_PARTS <<< "$CURRENT_VERSION"
MAJOR=${VERSION_PARTS[0]}
MINOR=${VERSION_PARTS[1]}
PATCH=${VERSION_PARTS[2]}

# Bump the version according to the type
if [ "$TYPE" == "major" ]; then
    MAJOR=$((MAJOR + 1))
    MINOR=0
    PATCH=0
elif [ "$TYPE" == "minor" ]; then
    MINOR=$((MINOR + 1))
    PATCH=0
elif [ "$TYPE" == "patch" ]; then
    PATCH=$((PATCH + 1))
fi

# Create the new version
NEW_VERSION="${MAJOR}.${MINOR}.${PATCH}"

# Update the VERSION file
echo $NEW_VERSION > VERSION

# Update chadburn.go
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
sed -i "s/var version string/var version = \"${NEW_VERSION}\"/" chadburn.go
sed -i "s/var build string/var build = \"${BUILD_DATE}\"/" chadburn.go

echo "Version bumped from $CURRENT_VERSION to $NEW_VERSION"
echo "Build date set to $BUILD_DATE"
echo "Don't forget to commit the changes!" 