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


.PHONY: release
release:
	@if [ ! -f ".release-env" ]; then \
		echo ".release-env is required for release";\
		exit 1;\
	fi
	docker run \
		--rm \
		--platform=linux/arm64 \
		-e CGO_ENABLED=1 \
		--env-file .release-env \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/powerpipe \
		-v `pwd`/../pipe-fittings:/go/src/pipe-fittings \
		-w /go/src/powerpipe \
		ghcr.io/turbot/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		release --clean --skip=validate

.PHONY: test test-race

test:
	go test ./...

test-race:
	go test -race ./...

# Run the test generator to create the acceptance tests(update the paths in main function in tests/acceptance/test_generator/generate.go)
build-tests:
	go run tests/acceptance/test_generator/generate.go