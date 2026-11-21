OUTPUT_DIR?=/usr/local/bin
PACKAGE_NAME          := github.com/turbot/powerpipe
GOLANG_CROSS_VERSION  ?= gcc13-osxcross-20251006102018

.PHONY: build
build:
	$(eval TIMESTAMP := $(shell date +%Y%m%d%H%M%S))
	$(eval GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD | sed 's/[\/_]/-/g' | sed 's/[^a-zA-Z0-9.-]//g'))

	go build -o $(OUTPUT_DIR) -ldflags "-X main.version=0.0.0-dev-$(GIT_BRANCH).$(TIMESTAMP)" .

dashboard_assets:
	$(MAKE) -C ui/dashboard

# Builds the dashboard UI and writes the archive the binary embeds (internal/dashboardassets/assets.tar.gz).
# The release build (Shipyard) runs it before goreleaser, with VERSION set to the release's version (1.6.1),
# which becomes the UI's version. Without VERSION the UI has none, and skips its version check.
# CI=false: on a CI runner, the UI build would otherwise fail on lint warnings.
.PHONY: dashboard_assets_archive
dashboard_assets_archive:
	cd ui/dashboard && corepack yarn install --immutable && corepack yarn setup-material-symbols
	cd ui/dashboard && rm -rf build && CI=false REACT_APP_VERSION=$(VERSION) corepack yarn build
	if [ -n "$(VERSION)" ]; then printf '{"version":"%s"}\n' "$(VERSION)" > ui/dashboard/build/version.json; fi
	tar -czf internal/dashboardassets/assets.tar.gz -C ui/dashboard/build .

# The release's tests (Shipyard's results job runs them): the Go unit tests and the dashboard UI's unit tests.
# Run dashboard_assets_archive first, which installs the UI's packages.
.PHONY: release_test
release_test:
	go test ./...
	cd ui/dashboard && CI=true corepack yarn test --watchAll=false

.PHONY: all
all:
	$(MAKE) build
	$(MAKE) dashboard_assets

.PHONY: release-dry-run
release-dry-run:
	@docker run \
		--rm \
		--platform=linux/arm64 \
		-e CGO_ENABLED=1 \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/powerpipe \
		-v `pwd`/../pipe-fittings:/go/src/pipe-fittings \
		-w /go/src/powerpipe \
		ghcr.io/turbot/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		--clean --skip=validate --skip=publish --snapshot


# Run the test generator to create the acceptance tests(update the paths in main function in tests/acceptance/test_generator/generate.go)
build-tests:
	go run tests/acceptance/test_generator/generate.go