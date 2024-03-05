FROM debian:bullseye-slim

LABEL org.opencontainers.image.ref.name powerpipe \
      org.opencontainers.image.version ${TARGETVERSION} \
      org.opencontainers.image.url "https://powerpipe.io" \
      org.opencontainers.image.authors "Turbot HQ, Inc" \
      org.opencontainers.image.source "https://github.com/turbot/powerpipe" \
      org.opencontainers.image.description "Powerpipe container image ${TARGETVERSION}"

# Set environment variables
# USER_UID: Override the powerpipe UID (default: 9193)
# USER_GID: Override the powerpipe GID (default: 0)
# POWERPIPE_LOG_LEVEL: Set logging level (default: "off")
# POWERPIPE_UPDATE_CHECK: Disable auto-update checks (default: false)
# POWERPIPE_TELEMETRY: Disable telemetry (default: none)
ENV USER_UID=9193 \
    USER_GID=0 \
    POWERPIPE_LOG_LEVEL="off" \
    POWERPIPE_UPDATE_CHECK=false \
    POWERPIPE_TELEMETRY=none

# Declare build arguments
# TARGETVERSION: Specify the version to target
# TARGETARCH: Specify the architecture to target
ARG TARGETVERSION \
    TARGETARCH

# NOTE: Temporary workaround until repository is public    
ARG GH_ACCESS_TOKEN

# This RUN command configures Powerpipe within the container through the following steps:
# 1. Updates package lists with `apt-get update` to ensure installation of the latest versions of packages and dependencies.
# 2. Adds a system user 'powerpipe' with no login capabilities and assigns it to the existing group defined by USER_GID, facilitating secure Powerpipe operations.
# 3. Conditionally installs necessary tools and downloads Powerpipe based on the GH_ACCESS_TOKEN environment variable:
#    - If GH_ACCESS_TOKEN is set (indicating private or rate-limited GitHub access):
#        a. Installs `curl`, `jq`, and `gosu` via apt-get for release information fetching, JSON parsing, and specific user command execution, respectively.
#        b. Fetches Powerpipe release information from GitHub using `curl` with GH_ACCESS_TOKEN for authentication.
#        c. Extracts the download URL for the Powerpipe binary that matches the TARGETARCH architecture using `jq`.
#        d. Downloads the Powerpipe binary (tar.gz format) to /tmp, authenticated by GH_ACCESS_TOKEN.
#    - If GH_ACCESS_TOKEN is not set (public access):
#        a. Installs `wget` and `gosu` via apt-get for file downloading and specific user command execution.
#        b. Downloads the Powerpipe binary (tar.gz format) for the specified TARGETVERSION and TARGETARCH from GitHub using `wget`, without authentication.
# 4. Creates a directory at /opt/powerpipe for storing the Powerpipe binary.
# 5. Extracts the Powerpipe binary from the tar.gz file in /tmp to /opt/powerpipe.
# 6. Moves the Powerpipe executable to /usr/local/bin, adding it to the system PATH for easy execution.
# 7. Cleans up by running `apt-get clean` and removing unnecessary files, including the downloaded tar.gz file and apt lists, to minimize the image size.
# NOTE: Conditional handling of GH_ACCESS_TOKEN will be removed once the Powerpipe repository becomes public.
RUN apt-get update && \
    group_name=$(getent group ${USER_GID} | cut -d: -f1) && \
    adduser --system --disabled-login --ingroup $group_name --gecos "powerpipe user" --shell /bin/false --uid $USER_UID powerpipe && \
    if [ -n "$GH_ACCESS_TOKEN" ]; then \
        apt-get install -y curl jq gosu && \
        bash -c "\
        release_info=\$(curl -s -H 'Authorization: token $GH_ACCESS_TOKEN' 'https://api.github.com/repos/turbot/powerpipe/releases/tags/$TARGETVERSION') && \
        asset_url=\$(echo \$release_info | jq -r '.assets[] | select(.name == \"powerpipe.linux.$TARGETARCH.tar.gz\").url') && \
        curl -L -H 'Authorization: token $GH_ACCESS_TOKEN' -H 'Accept: application/octet-stream' \"\$asset_url\" -o /tmp/powerpipe.linux.$TARGETARCH.tar.gz"; \
    else \
        apt-get install -y wget gosu && \
        wget -nv https://github.com/turbot/powerpipe/releases/download/${TARGETVERSION}/powerpipe.linux.${TARGETARCH}.tar.gz -O /tmp/powerpipe.linux.${TARGETARCH}.tar.gz; \
    fi && \
    mkdir -p /opt/powerpipe && \
    tar xzf /tmp/powerpipe.linux.${TARGETARCH}.tar.gz -C /opt/powerpipe && \
    mv /opt/powerpipe/powerpipe /usr/local/bin/powerpipe && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/powerpipe.linux.${TARGETARCH}.tar.gz

# Creates the workspace directory with the correct ownership
USER powerpipe
WORKDIR /workspace
USER root

# This is the default port for the dashboard service
EXPOSE 9033

# Copy the entrypoint script into the image
COPY docker-entrypoint.sh /usr/local/bin

# Define the entrypoint and default command
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["powerpipe"]
