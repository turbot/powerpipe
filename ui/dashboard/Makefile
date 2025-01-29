package: zip

VERSION ?= 0.0.0

zip: copy
	mkdir -p assets
	cd build && zip -r -D ../assets/assets.zip *

copy: build
	mkdir -p ~/.powerpipe/dashboard/assets
	cd build && cp -r * ~/.powerpipe/dashboard/assets

build: setup-icons
	yarn build

setup-icons: install
	yarn setup-material-symbols

install: clean
	yarn install

clean:
	pwd
	[ ! -e assets/assets.zip ] || rm assets/assets.zip
	rm -rf build

run-local: setup-icons
	yarn start
