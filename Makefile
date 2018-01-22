TARGETS =  linux/arm64 darwin/amd64 linux/amd64
# windows/amd64 windows/386 linux/arm linux/386
GIT_COMMIT = $(shell git rev-parse HEAD)
BUILD_TIME = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ" | tr -d '\n')

usage:
	@echo ""
	@echo "Task                 : Description"
	@echo "-----------------    : -------------------"
	@echo "make setup           : Install all necessary dependencies"
	@echo "make build           : Generate production build for current OS"
	@echo "make bootstrap       : Install cross-compilation toolchain"
	@echo "make release         : Generate binaries for all supported OSes"
	@echo "make clean           : Remove all build files and reset assets"
	@echo "make upload          : Uploads binaries to S3, needs GOLANG_S3_KEY (${GOLANG_S3_KEY}) and GOLANG_S3_VERSION (${GOLANG_S3_VERSION})"
	@echo ""

assets:
	@echo Oh yeah, no assets here.

build: assets
	go build
	@echo "You can now execute ./downtimeserver"

release: assets
	mkdir -p ./bin
	@echo "Building binaries..."
	$(GOPATH)/bin/gox \
		-osarch="$(TARGETS)" \
		-output="./bin/oauth2_proxy_{{.OS}}_{{.Arch}}"

bootstrap:
	gox -build-toolchain

setup:
	go get github.com/mitchellh/gox
	go get

clean:
	rm -f ./downtimeserver
	rm -f ./downtimeserver_*
	rm -rf ./bin/*
	make assets

upload: release
	s3cmd --access_key="${GOLANG_S3_ID}" --secret_key="${GOLANG_S3_KEY}" --bucket-location="eu-west-1" --progress --stats --force --stop-on-error --acl-public put bin/oauth2_proxy_linux_arm64 s3://meceap-distro-box/oauth2_proxy/oauth2_proxy_linux_aarch64_v${GOLANG_S3_VERSION}
	s3cmd --access_key="${GOLANG_S3_ID}" --secret_key="${GOLANG_S3_KEY}" --bucket-location="eu-west-1" --progress --stats --force --stop-on-error --acl-public put bin/oauth2_proxy_linux_amd64 s3://meceap-distro-box/oauth2_proxy/oauth2_proxy_linux_x86_64_v${GOLANG_S3_VERSION}

