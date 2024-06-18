#!/bin/sh -e

# update the version of top level mod

mv .powerpipe/mods/github.com/pskrbasu/powerpipe-mod-1@v1.0.0 .powerpipe/mods/github.com/pskrbasu/powerpipe-mod-1@v0.1.0
JSON_FILE=".mod.cache.json"
jq '.local["github.com/pskrbasu/powerpipe-mod-1"].version = "0.1.0" | 
    . as $json | 
    .["github.com/pskrbasu/powerpipe-mod-1@v1.0.0"] as $oldKey | 
    del(.["github.com/pskrbasu/powerpipe-mod-1@v1.0.0"]) | 
    .["github.com/pskrbasu/powerpipe-mod-1@v0.1.0"] = $oldKey' "$JSON_FILE" > tmp.$$.json
mv tmp.$$.json "$JSON_FILE"