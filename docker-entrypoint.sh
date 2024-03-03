#!/usr/bin/env bash
set -Eeo pipefail

# if first arg is anything other than `powerpipe`, assume we want to run powerpipe
# this is for when other commands are passed to the container
if [ "${1:0}" != 'powerpipe' ]; then
    set -- powerpipe "$@"
fi

exec "$@"