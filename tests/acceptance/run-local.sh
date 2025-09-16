#!/bin/bash

MY_PATH="`dirname \"$0\"`"              # relative
MY_PATH="`( cd \"$MY_PATH\" && pwd )`"  # absolutized and normalized

# TODO PSKR review all exports and remove unused ones in powerpipe
export POWERPIPE_INSTALL_DIR=$(mktemp -d)
export TZ=UTC
export WD=$(mktemp -d)

trap "cd -;code=$?;rm -rf $POWERPIPE_INSTALL_DIR; exit $code" EXIT

cd "$WD"
echo "Working directory: $WD"
# setup a powerpipe installation
echo "Install directory: $POWERPIPE_INSTALL_DIR"

# Temporarily disable 'exit on error' since we want to run the check command and not exit if it fails
set +e
powerpipe check > /dev/null 2>&1
check_status=$?
set -e

echo "Installation complete at $POWERPIPE_INSTALL_DIR"

echo "Starting steampipe service..."
steampipe service stop --force > /dev/null 2>&1 || true
steampipe service start > /dev/null 2>&1

if [ $# -eq 0 ]; then
  # Run all test files
  "$MY_PATH/run.sh"
else
  "$MY_PATH/run.sh" "${1}"
fi

echo "Stopping steampipe service..."
steampipe service stop
