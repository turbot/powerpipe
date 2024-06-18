#!/bin/sh -e

# update the commit hash of top level mod

JSON_FILE=".mod.cache.json"
jq '.local["github.com/pskrbasu/powerpipe-mod-1"].commit = "43746d1cd00bb9ecdccc9c6babe0a93bef4ee446"' "$JSON_FILE" > tmp.$$.json && mv tmp.$$.json "$JSON_FILE"