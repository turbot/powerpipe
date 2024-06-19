#!/bin/sh -e

# This script is used in the mod.bats tests to update the version of the top level mod. This is done to simulate a new version of the top level 
# mod being available. This script is using jq to modify the version property in .mod.cache.json and also to rename the mod folders to reflect the new version.

# update the version of top level mod
mv .powerpipe/mods/github.com/pskrbasu/powerpipe-mod-1@v1.0.0 .powerpipe/mods/github.com/pskrbasu/powerpipe-mod-1@v0.1.0
JSON_FILE=".mod.cache.json"
jq '.local["github.com/pskrbasu/powerpipe-mod-1"].version = "0.1.0" | 
    . as $json | 
    .["github.com/pskrbasu/powerpipe-mod-1@v1.0.0"] as $oldKey | 
    del(.["github.com/pskrbasu/powerpipe-mod-1@v1.0.0"]) | 
    .["github.com/pskrbasu/powerpipe-mod-1@v0.1.0"] = $oldKey' "$JSON_FILE" > tmp.$$.json
mv tmp.$$.json "$JSON_FILE"