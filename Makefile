OUTPUT_DIR?=/usr/local/bin

.PHONY: build

build:
	$(eval MAJOR := $(shell cat version.json | jq '.major'))
	$(eval MINOR := $(shell cat version.json | jq '.minor'))
	$(eval PATCH := $(shell cat version.json | jq '.patch'))
	$(eval TIMESTAMP := $(shell date +%Y%m%d%H%M%S))

	go build -o $(OUTPUT_DIR) -ldflags "-X main.version=$(MAJOR).$(MINOR).$(PATCH)-dev.$(TIMESTAMP)" .


dashboard_assets:
	$(MAKE) -C ui/dashboard

all:
	build ,dashboard_assets

