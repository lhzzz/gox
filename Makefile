.PHONY: all
all: format build

# ==============================================================================
# Build options

ROOT_PACKAGE=.

# ==============================================================================
# Includes

include build/lib/common.mk
include build/lib/golang.mk
include build/lib/image.mk

# ==============================================================================
# Targets

.PHONY: format
format:
	go fmt ./internal/... ./cmd/...

## build: Build source code for host platform.
.PHONY: build
build:
	@$(MAKE) go.build

## build.multiarch: Build source code for multiple platforms. See option PLATFORMS.
.PHONY: build.multiarch
build.multiarch:
	@$(MAKE) go.build.multiarch

## image: Build docker images for host arch.
.PHONY: image
image:
	@$(MAKE) image.build

## image.multiarch: Build docker images for multiple platforms. See option PLATFORMS.
.PHONY: image.multiarch
image.multiarch:
	@$(MAKE) image.build.multiarch

## push: Build docker images for host arch and push images to registry.
.PHONY: push
push:
	@$(MAKE) image.push

## push.multiarch: Build docker images for multiple platforms and push images to registry.
.PHONY: push.multiarch
push.multiarch:
	@$(MAKE) image.push.multiarch

## manifest: Build docker images for host arch and push manifest list to registry.
.PHONY: manifest
manifest:
	@$(MAKE) image.manifest.push

## manifest.multiarch: Build docker images for multiple platforms and push manifest lists to registry.
.PHONY: manifest.multiarch
manifest.multiarch:
	@$(MAKE) image.manifest.push.multiarch

.PHONY: buildx.multiarch
buildx.multiarch:
	@$(MAKE) image.buildx.push.multiarch

## clean: Remove all files that are created by building.
.PHONY: clean
clean:
	@$(MAKE) go.clean

.PHONY: release.build
release.build:
	@$(MAKE) push.multiarch

## help: Show this help info.
.PHONY: help
help: Makefile
	@echo -e "\nUsage: make <TARGETS> <OPTIONS> ...\n\nTargets:"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo "$$USAGE_OPTIONS"
