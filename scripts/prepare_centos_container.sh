#!/bin/sh
# This is a a script to install dependencies/packages, create user, and assign necessary permissions in the centos 9 container.
# Used in release smoke tests. 

# update yum and install required packages
yum install -y epel-release tar ca-certificates jq curl

export PATH=$PATH:/usr/local/bin:/usr/bin:/bin

# install steampipe latest
/bin/sh -c "$(curl -fsSL https://steampipe.io/install/steampipe.sh)"

# Extract the powerpipe binary
tar -xzf  /artifacts/linux.tar.gz -C /usr/local/bin

# Make the binaries executable
chmod +x /usr/local/bin/steampipe
chmod +x /usr/local/bin/powerpipe

# Create user, since steampipe cannot be run as root
useradd -m steampipe

# Verify jq installation
jq --version

# Make the scripts executable
chown steampipe:steampipe /scripts/smoke_test.sh
chmod +x /scripts/smoke_test.sh