# The image that tests/acceptance/run-docker.sh runs the acceptance tests in: the tools of the
# CI acceptance workflow (.github/workflows/11-test-acceptance.yaml), and powerpipe built from
# the checkout.
#
# The base image digest is also in .claude/agent-fleet/inbox-gates.yml (check `acceptance`,
# `images`). Change both together.

FROM golang:1.26.10-bookworm@sha256:5de7c99eb78ca5e515a4c340235d0955ec5964a5071c492c15503d7c267da473 AS tools

RUN apt-get update \
 && apt-get install -y --no-install-recommends jq openssl \
 && rm -rf /var/lib/apt/lists/*

# Steampipe, checked against the SHA256 of steampipe_linux_<arch>.tar.gz of STEAMPIPE_VERSION
ARG STEAMPIPE_VERSION=v2.4.9
ARG STEAMPIPE_SHA256_AMD64=3a59bbed43019cb17b9b9adbc8b559772c6b49ac22e6579bd63c9cae5e4ba413
ARG STEAMPIPE_SHA256_ARM64=fcaea67ce59cc5f50bbf35d8739443f61e4c6c423add8747ff41fb32fd75e050
RUN set -eu; \
    arch="$(dpkg --print-architecture)"; \
    case "$arch" in \
      amd64) sha256="$STEAMPIPE_SHA256_AMD64" ;; \
      arm64) sha256="$STEAMPIPE_SHA256_ARM64" ;; \
      *) echo "No Steampipe SHA256 for $arch" >&2; exit 1 ;; \
    esac; \
    file="steampipe_linux_${arch}.tar.gz"; \
    cd /tmp; \
    curl -fsSLO "https://github.com/turbot/steampipe/releases/download/${STEAMPIPE_VERSION}/$file"; \
    echo "$sha256  $file" | sha256sum -c -; \
    tar -xzf "$file" -C /usr/local/bin steampipe; \
    rm "$file"

ARG JD_VERSION=v1.9.2
RUN GOBIN=/usr/local/bin go install "github.com/josephburnett/jd@${JD_VERSION}"

# Steampipe does not run as root
RUN useradd --create-home --uid 1000 tester
USER tester
ENV GOCACHE=/home/tester/.cache/go-build \
    GOMODCACHE=/home/tester/go/pkg/mod

# install the Steampipe database in the image, so that every test file starts from the same one
RUN steampipe service start && steampipe service stop

FROM tools AS test

RUN mkdir -p /home/tester/powerpipe /home/tester/bin
WORKDIR /home/tester/powerpipe
COPY --chown=tester:tester . .

# with CGO, as the release build: the DuckDB and SQLite backends need it
RUN --mount=type=cache,target=/home/tester/.cache/go-build,uid=1000,gid=1000 \
    --mount=type=cache,target=/home/tester/go/pkg/mod,uid=1000,gid=1000 \
    CGO_ENABLED=1 go build -buildvcs=false -o /home/tester/bin/powerpipe .

# LANG as on the CI runners
ENV PATH=/home/tester/bin:$PATH \
    POWERPIPE_UPDATE_CHECK=false \
    LANG=C.UTF-8
