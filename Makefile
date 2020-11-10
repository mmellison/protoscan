# ---------------------------------
#   	Variables + Options

# The default set of binaries
binaries=bin/scanner

pkg_sources := $(shell find pkg/ -name '*.go')

# Defaults for the Container Image
IMAGE_NAME ?= seglberg/protoscan
IMAGE_TAG ?= local
IMAGE_MAKE_TARGET ?= all

# ---------------------------------
#   	TOP LEVEL TARGETS

all: fmt vet $(binaries);

%.dev: clean bin/%;

clean:
	rm -rf $(binaries)

# ---------------------------------
#   	Test / Lint Targets

.PHONY:
fmt:
	go fmt ./...

.PHONY:
vet:
	go vet ./...

.PHONY:
lint:
	golangci-lint run

.PHONY:
lint-fix:
	golangci-lint run --fix


# ---------------------------------
#		Build Targets

.SECONDEXPANSION:
bin/%: cmd/% $$(wildcard cmd/%/**/*) $(pkg_sources) go.mod go.sum
	CGO_ENABLED=0 go build -ldflags='${_LDFLAGS}' -o $@ ./$<

.PHONY:
compress_all:
	upx --exact --lzma --best bin/*

.PHONY:
compress_%:
	upx --exact --lzma --best bin/$*

# ---------------------------------
#		Container Targets

image-build-with-buildah:
	@echo "Building ${IMAGE_NAME}:${IMAGE_TAG} with Buildah"
	buildah bud -t ${IMAGE_NAME}:${IMAGE_TAG} --build-arg MAKE_TARGET=${IMAGE_MAKE_TARGET} .

image-build-with-docker:
	@echo "Building ${IMAGE_NAME}:${IMAGE_TAG} with Docker"
	docker build -t ${IMAGE_NAME}:${IMAGE_TAG} --build-arg MAKE_TARGET=${IMAGE_MAKE_TARGET} .

