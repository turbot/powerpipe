#!/bin/sh
# TODO(everyone): Keep this script simple and easily auditable.

# Function to check if a command exists
command_exists() {
  type "$1" &> /dev/null
}

set -e

if ! command -v tar >/dev/null; then
	echo "Error: 'tar' is required to install Powerpipe." 1>&2
	exit 1
fi

if ! command -v gzip >/dev/null; then
	echo "Error: 'gzip' is required to install Powerpipe." 1>&2
	exit 1
fi

if ! command -v install >/dev/null; then
	echo "Error: 'install' is required to install Powerpipe." 1>&2
	exit 1
fi

if command -v powerpipe >/dev/null; then
	# Check if port 9033 is in use
	if command_exists ss; then
		# Use ss if available
		if ss -lnt | grep ':9033' >/dev/null; then
			echo "Error: Port 9033 is already in use. Please stop the service using this port before running installation." 1>&2
			exit 1
		fi
	elif command_exists netstat; then
		# Fallback to netstat if ss is not available
		if netstat -an | grep ':9033' >/dev/null; then
			echo "Error: Port 9033 is already in use. Please stop the service using this port before running installation." 1>&2
			exit 1
		fi
	else
		echo "Error: Neither ss nor netstat command is available." 1>&2
		exit 1
	fi
fi

if [ "$OS" = "Windows_NT" ]; then
	echo "Error: Windows is not supported yet." 1>&2
	exit 1
else
	case $(uname -sm) in
	"Darwin x86_64") target="darwin.amd64.tar.gz" ;;
	"Darwin arm64") target="darwin.arm64.tar.gz" ;;
	"Linux x86_64") target="linux.amd64.tar.gz" ;;
	"Linux aarch64") target="linux.arm64.tar.gz" ;;
	*) echo "Error: '$(uname -sm)' is not supported yet." 1>&2;exit 1 ;;
	esac
fi

if [ $# -eq 0 ]; then
	powerpipe_uri="https://github.com/turbot/powerpipe/releases/latest/download/powerpipe.${target}"
else
	powerpipe_uri="https://github.com/turbot/powerpipe/releases/download/${1}/powerpipe.${target}"
fi

bin_dir="/usr/local/bin"
exe="$bin_dir/powerpipe"

test -z "$tmp_dir" && tmp_dir="$(mktemp -d)"
mkdir -p "${tmp_dir}"
tmp_dir="${tmp_dir%/}"

echo "Created temporary directory at $tmp_dir. Changing to $tmp_dir"
cd "$tmp_dir"

# set a trap for a clean exit - even in failures
trap 'rm -rf $tmp_dir' EXIT

case $(uname -s) in
	"Darwin") zip_location="$tmp_dir/powerpipe.tar.gz" ;;
	"Linux") zip_location="$tmp_dir/powerpipe.tar.gz" ;;
	*) echo "Error: powerpipe is not supported on '$(uname -s)' yet." 1>&2;exit 1 ;;
esac

echo "Downloading from $powerpipe_uri"
if command -v wget >/dev/null; then
	# because --show-progress was introduced in 1.16.
	wget --help | grep -q '\--showprogress' && _FORCE_PROGRESS_BAR="--no-verbose --show-progress" || _FORCE_PROGRESS_BAR=""
	# prefer an IPv4 connection, since github.com does not handle IPv6 connections properly.
	if ! wget --prefer-family=IPv4 --progress=bar:force:noscroll $_FORCE_PROGRESS_BAR -O "$zip_location" "$powerpipe_uri"; then
        echo "Could not find version $1"
        exit 1
    fi
elif command -v curl >/dev/null; then
	# curl uses HappyEyeball for connections, therefore, no preference is required
    if ! curl --fail --location --progress-bar --output "$zip_location" "$powerpipe_uri"; then
        echo "Could not find version $1"
        exit 1
    fi
else
    echo "Unable to find wget or curl. Cannot download."
    exit 1
fi

echo "Deflating downloaded archive"
tar -xf "$zip_location" -C "$tmp_dir"

echo "Installing"
install -d "$bin_dir"
install "$tmp_dir/powerpipe" "$bin_dir"

echo "Applying necessary permissions"
chmod +x $exe

echo "Removing downloaded archive"
rm "$zip_location"

echo "Powerpipe was installed successfully to $exe"

if ! command -v $bin_dir/powerpipe >/dev/null; then
	echo "Powerpipe was installed, but could not be executed. Are you sure '$bin_dir/powerpipe' has the necessary permissions?"
	exit 1
fi
