#!/bin/sh
# This is a a script to install dependencies/packages, create user, and assign necessary permissions in the ubuntu 24 container.
# Used in release smoke tests.

set -e

# update apt and install required packages
apt-get update
apt-get install -y tar ca-certificates jq curl

# install steampipe latest
/bin/sh -c "$(curl -fsSL https://steampipe.io/install/steampipe.sh)"

# Extract the powerpipe binary
if [ "$(uname -m)" = "aarch64" ]; then
  tar -xzf /artifacts/linux-arm.tar.gz -C /usr/local/bin
else
  tar -xzf /artifacts/linux.tar.gz -C /usr/local/bin
fi

# Make the binaries executable
chmod +x /usr/local/bin/steampipe
chmod +x /usr/local/bin/powerpipe

# Create user, since steampipe cannot be run as root
useradd -m steampipe

# Make the scripts executable
chown steampipe:steampipe /scripts/smoke_test.sh
chmod +x /scripts/smoke_test.sh
