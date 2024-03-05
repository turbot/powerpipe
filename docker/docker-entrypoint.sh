#!/usr/bin/env bash

# The `set -Eeo pipefail` option controls the behavior of the script on errors:
# - `set -E` ensures that any ERR traps are inherited by shell functions, command substitutions, and commands executed in subshells.
# - `set -e` exits the script immediately if any command within the script or a subshell returns a non-zero status (indicating failure).
# - `set -o pipefail` causes the script to fail if any of the commands in a pipeline fail, not just the last command.
set -Eeo pipefail

log_if_debug() {
    # Convert POWERPIPE_LOG_LEVEL to lowercase and compare
    if [ "${POWERPIPE_LOG_LEVEL,,}" = "debug" ]; then
        echo "$@"
    fi
}

log_if_debug "Running docker entrypoint script..."

# Function to check and change ownership of a mounted volume
check_and_change_ownership() {
    local mount_path=$1

    if [ "$mount_path" = "/workspace" ]; then
        log_if_debug "Skipping /workspace as it's the workspace directory"
        return
    fi

    log_if_debug "Checking the ownership of the volume mounted at $mount_path..."

    # Check if the volume is mounted
    if mount | grep -q "on $mount_path type"; then
        log_if_debug "Checking if the volume at $mount_path is owned by root:root..."

        # Check if the volume is empty and owned by root
        mount_ownership=$(stat -c "%u:%g" $mount_path)

        if [ "$mount_ownership" == "0:0" ]; then
            log_if_debug "Volume at $mount_path is owned by root:root."

            # Change the owner of the volume to USER_UID:USER_GID
            chown "$USER_UID:$USER_GID" $mount_path
            log_if_debug "Changed ownership of the volume at $mount_path to $USER_UID:$USER_GID."
        elif [ "$mount_ownership" != "$USER_UID:$USER_GID" ]; then
            echo "WARNING: Directory $mount_path has ownership of $mount_ownership and does not match the UID/GID $USER_UID:$USER_GID of the powerpipe user."
            echo "         Resolve by either overriding the environment variables USER_UID and USER_GID."
            echo "         Or by changing the ownership of the directory."
            echo "         Ownership $mount_ownership of $mount_path will not be modified."
        fi
    else
        log_if_debug "No volume is mounted at $mount_path. Skipping."
    fi
}

log_if_debug "Setting up default UID and GID if not provided..."

# Default UID and GID for powerpipe user if not provided
DEFAULT_UID=9033
DEFAULT_GID=0

log_if_debug "Checking and changing ownership of mounted volumes if necessary..."

# Find all unique devices associated with mounts within /etc or its subdirectories
readarray -t etc_devices < <(mount | grep ' on /etc' | awk '{print $1}' | sort -u)

# Convert array to a string for easy checking
ignore_devices=$(IFS="|"; echo "${etc_devices[*]}")

# Obtain mount points from the mount command and loop through them
while IFS= read -r line; do
    mount_device=$(echo "$line" | awk '{print $1}')
    mount_path=$(echo "$line" | awk '{print $3}')

    if [ -f "$mount_path" ]; then
        log_if_debug "Skipping $mount_path as it's a file"
        continue
    fi

    # This is mounted under a different partition mount point
    # We skip all mounts that are under the same device as /etc or its subdirectories.
    check_and_change_ownership "$mount_path"
done < <(mount | grep '^/dev')

# Check current ownership of /workspace
workspace_uid=$(stat -c '%u' /workspace)
workspace_gid=$(stat -c '%g' /workspace)

log_if_debug "Using USER_UID=$USER_UID and USER_GID=$USER_GID."
log_if_debug "Checking if the current user's UID/GID is different from the default UID/GID..."

# Check if the current UID/GID is different from the provided or default USER_UID/USER_GID
if [ "$(id -u powerpipe)" != "$USER_UID" ] || [ "$(id -g powerpipe)" != "$USER_GID" ]; then
    log_if_debug "Current UID/GID is different. Updating powerpipe user and group IDs..."

    # Create or modify the user's primary group if USER_GID is provided and it's not the default GID
    if [ "$USER_GID" != "$DEFAULT_GID" ]; then
        if ! getent group $USER_GID >/dev/null; then
            log_if_debug "Creating group powerpipegroup with GID $USER_GID."
            groupadd -g $USER_GID powerpipegroup
        fi
        log_if_debug "Modifying powerpipe's primary group to $USER_GID."
        usermod -g $USER_GID powerpipe
    fi

    # Modify the powerpipe user's UID if it's provided and not the default UID
    if [ "$USER_UID" != "$DEFAULT_UID" ]; then
        log_if_debug "Modifying powerpipe's UID to $USER_UID."
        usermod -u $USER_UID powerpipe
    fi

    # Compare current ownership with desired (USER_UID and USER_GID) and change if necessary
    if [ "$workspace_uid" -eq "$DEFAULT_UID" ] && [ "$workspace_gid" -eq "$DEFAULT_GID" ]; then
        log_if_debug "Ownership of /workspace is the default UID/GID. Changing..."
        chown "$USER_UID:$USER_GID" /workspace
        workspace_uid=$USER_UID
        workspace_gid=$USER_GID
    fi
else
    log_if_debug "Current UID/GID is the same as the provided or default UID/GID. Skipping user and group ID updates."
fi

cd /workspace
log_if_debug "Changed directory to /workspace."

log_if_debug "Evaluating the initial argument to determine if it's the 'powerpipe' command. If not, 'powerpipe' will be prepended to ensure the powerpipe CLI is executed."
# if first arg is anything other than `powerpipe`, assume we want to run powerpipe
# this is for when other commands are passed to the container
if [ "${1:0}" != 'powerpipe' ]; then
    set -- powerpipe "$@"
fi

log_if_debug "Final command configuration set. Proceeding to execute the 'powerpipe' CLI with the provided arguments."
# Now, execute the command provided to the docker run
exec gosu powerpipe "$@"
