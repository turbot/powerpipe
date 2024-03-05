#!/bin/bash
set -e

# GitHub repository details and token as arguments
REPO=$1
TAG_NAME=$2
FILE_NAME=$3
GITHUB_TOKEN=$4

# Get the release info including assets
RELEASE_INFO=$(curl -s -H "Authorization: token $GITHUB_TOKEN" "https://api.github.com/repos/$REPO/releases/tags/$TAG_NAME")

# Extract the asset download URL
ASSET_URL=$(echo $RELEASE_INFO | jq -r ".assets[] | select(.name == \"$FILE_NAME\").url")

# Download the asset
curl -L -H "Authorization: token $GITHUB_TOKEN" -H "Accept: application/octet-stream" "$ASSET_URL" -o $FILE_NAME