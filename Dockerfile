FROM debian:bullseye-slim
LABEL maintainer="Turbot Support <help@turbot.com>"

ARG TARGETVERSION
ARG TARGETARCH
# accept the GitHub token as a build argument (required till we make the repo public)
ARG GITHUB_TOKEN

# add a non-root 'steampipe' user
RUN adduser --system --disabled-login --ingroup 0 --gecos "powerpipe user" --shell /bin/false --uid 9193 powerpipe

# Copy the download script into the image
COPY download_release.sh /download_release.sh

# Install dependencies (jq and curl)
RUN apt-get update -y && apt-get install -y jq curl && rm -rf /var/lib/apt/lists/*

# Run the download script
RUN /download_release.sh turbot/powerpipe $TARGETVERSION powerpipe.linux.$TARGETARCH.tar.gz $GITHUB_TOKEN \
    && tar xzf powerpipe.linux.$TARGETARCH.tar.gz \
    && mv powerpipe /usr/local/bin/ \
    && rm -rf /tmp/* powerpipe.linux.$TARGETARCH.tar.gz /download_release.sh

# Change user to non-root
USER powerpipe:0

# Use a constant workspace directory that can be mounted to
WORKDIR /workspace

# expose dashboard service default port
EXPOSE 9194

COPY docker-entrypoint.sh /usr/local/bin
ENTRYPOINT [ "docker-entrypoint.sh" ]
CMD [ "powerpipe"]
