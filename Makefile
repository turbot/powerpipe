OUTPUT_DIR?=/usr/local/bin
PACKAGE_NAME          := github.com/turbot/powerpipe
GOLANG_CROSS_VERSION  ?= v1.21.5

.PHONY: build
build:
	$(eval MAJOR := $(shell cat internal/version/version.json | jq '.major'))
	$(eval MINOR := $(shell cat internal/version/version.json | jq '.minor'))
	$(eval PATCH := $(shell cat internal/version/version.json | jq '.patch'))
	$(eval TIMESTAMP := $(shell date +%Y%m%d%H%M%S))

	go build -o $(OUTPUT_DIR) -ldflags "-X main.version=$(MAJOR).$(MINOR).$(PATCH)-dev.$(TIMESTAMP)" .


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
		-e CGO_ENABLED=1 \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/powerpipe \
		-v `pwd`/../pipe-fittings:/go/src/pipe-fittings \
		-w /go/src/powerpipe \
		ghcr.io/goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		--clean --skip-validate --skip-publish --snapshot --rm-dist


.PHONY: release
release:
	@if [ ! -f ".release-env" ]; then \
		echo ".release-env is required for release";\
		exit 1;\
	fi
	docker run \
		--rm \
		-e CGO_ENABLED=1 \
		--env-file .release-env \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/powerpipe \
		-v `pwd`/../pipe-fittings:/go/src/pipe-fittings \
		-w /go/src/powerpipe \
		ghcr.io/goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		release --clean --skip-validate
