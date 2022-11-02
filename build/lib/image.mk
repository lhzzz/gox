DOCKER := docker

ifeq ($(origin CI_PIPELINE_ID), undefined)
CI_PIPELINE_ID := local
endif

ifeq ($(origin CI_PROJECT_NAMESPACE), undefined)
CI_PROJECT_NAMESPACE := singer
endif

ifeq ($(origin CI_PROJECT_NAME), undefined)
CI_PROJECT_NAME := go-server
endif

ifeq ($(origin CI_COMMIT_REF_NAME), undefined)
CI_COMMIT_REF_NAME := $(shell git symbolic-ref --short HEAD)
endif

CI_COMMIT_REF_NAME_FIX=$(subst /,-,$(CI_COMMIT_REF_NAME))

REGISTRY_PREFIX ?= registry.com/$(CI_PROJECT_NAMESPACE)

EXTRA_ARGS ?=
_DOCKER_BUILD_EXTRA_ARGS :=

ifneq ($(EXTRA_ARGS), )
_DOCKER_BUILD_EXTRA_ARGS += $(EXTRA_ARGS)
endif

# Determine image files by looking into build/docker/*/Dockerfile
IMAGES_DIR ?= $(wildcard ${ROOT_DIR}/build/docker/*)
# Determine images names by stripping out the dir names
IMAGES ?= $(filter-out tools,$(foreach image,${IMAGES_DIR},$(notdir ${image})))

ifeq (${IMAGES},)
  $(error Could not determine IMAGES, set ROOT_DIR or run in source dir)
endif

.PHONY: image.daemon.verify
image.daemon.verify:
	$(eval PASS := $(shell $(DOCKER) version | grep -q -E 'Experimental: {1,5}true' && echo 1 || echo 0))
	@if [ $(PASS) -ne 1 ]; then \
		echo "Experimental features of Docker daemon is not enabled. Please add \"experimental\": true in '/etc/docker/daemon.json' and then restart Docker daemon."; \
		exit 1; \
	fi

.PHONY: image.build
image.build: $(addprefix image.build., $(addprefix $(IMAGE_PLAT)., $(IMAGES)))

.PHONY: image.build.multiarch
image.build.multiarch: $(foreach p,$(PLATFORMS),$(addprefix image.build., $(addprefix $(p)., $(IMAGES))))

.PHONY: image.build.%
# image.build.%: go.build.%
image.build.%:
	$(eval COMMAND := $(word 2,$(subst ., ,$*)))
	$(eval PLATFORM := $(word 1,$(subst ., ,$*)))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	$(eval IMAGE := $(COMMAND))
	$(eval IMAGE_PLAT := $(subst _,/,$(PLATFORM)))
	@echo "===========> Building docker image $(IMAGE) $(CI_PIPELINE_ID) for $(IMAGE_PLAT)"
	@mkdir -p $(TMP_DIR)/$(IMAGE)
	@cat $(ROOT_DIR)/build/docker/$(IMAGE)/Dockerfile > $(TMP_DIR)/$(IMAGE)/Dockerfile
	@cp $(OUTPUT_DIR)/$(IMAGE_PLAT)/$(IMAGE) $(TMP_DIR)/$(IMAGE)/
	@DST_DIR=$(TMP_DIR)/$(IMAGE) $(ROOT_DIR)/build/docker/$(IMAGE)/build.sh 2>/dev/null || true
	$(eval BUILD_SUFFIX := $(_DOCKER_BUILD_EXTRA_ARGS) --pull -t $(REGISTRY_PREFIX)/$(CI_PROJECT_NAME)-$(ARCH):$(IMAGE)-$(CI_COMMIT_REF_NAME_FIX)-$(CI_PIPELINE_ID) $(TMP_DIR)/$(IMAGE))
	@if [ $(shell $(GO) env GOARCH) != $(ARCH) ] ; then \
		$(MAKE) image.daemon.verify ;\
		$(DOCKER) build --platform $(IMAGE_PLAT) $(BUILD_SUFFIX) ; \
	else \
		$(DOCKER) build $(BUILD_SUFFIX) ; \
	fi
	@rm -rf $(TMP_DIR)/$(IMAGE)

.PHONY: image.push
image.push: $(addprefix image.push., $(addprefix $(IMAGE_PLAT)., $(IMAGES)))

.PHONY: image.push.multiarch
image.push.multiarch: $(foreach p,$(PLATFORMS),$(addprefix image.push., $(addprefix $(p)., $(IMAGES)))) 

.PHONY: image.push.%
image.push.%: image.build.%
	@echo "===========> Pushing image $(IMAGE)-$(CI_COMMIT_REF_NAME_FIX)-$(CI_PIPELINE_ID) to $(REGISTRY_PREFIX)/$(CI_PROJECT_NAME)-$(ARCH)"
	$(DOCKER) push $(REGISTRY_PREFIX)/$(CI_PROJECT_NAME)-$(ARCH):$(IMAGE)-$(CI_COMMIT_REF_NAME_FIX)-$(CI_PIPELINE_ID)

.PHONY: image.manifest.push
image.manifest.push: export DOCKER_CLI_EXPERIMENTAL := enabled
image.manifest.push: $(addprefix image.manifest.push., $(addprefix $(IMAGE_PLAT)., $(IMAGES)))

.PHONY: image.manifest.push.%
image.manifest.push.%: image.push.% image.manifest.remove.%
	@echo "===========> Pushing manifest $(IMAGE) $(CI_PIPELINE_ID) to $(REGISTRY_PREFIX) and then remove the local manifest list"
	@$(DOCKER) manifest create $(REGISTRY_PREFIX)/$(CI_PROJECT_NAME):$(IMAGE)-$(CI_COMMIT_REF_NAME_FIX)-$(CI_PIPELINE_ID) \
		$(REGISTRY_PREFIX)/$(CI_PROJECT_NAME)-$(ARCH):$(IMAGE)-$(CI_COMMIT_REF_NAME_FIX)-$(CI_PIPELINE_ID)
	@$(DOCKER) manifest annotate $(REGISTRY_PREFIX)/$(CI_PROJECT_NAME):$(IMAGE)-$(CI_COMMIT_REF_NAME_FIX)-$(CI_PIPELINE_ID) \
		$(REGISTRY_PREFIX)/$(CI_PROJECT_NAME)-$(ARCH):$(IMAGE)-$(CI_COMMIT_REF_NAME_FIX)-$(CI_PIPELINE_ID) \
		--os $(OS) --arch ${ARCH}
	@$(DOCKER) manifest push --purge $(REGISTRY_PREFIX)/$(CI_PROJECT_NAME):$(IMAGE)-$(CI_COMMIT_REF_NAME_FIX)-$(CI_PIPELINE_ID)

# Docker cli has a bug: https://github.com/docker/cli/issues/954
# If you find your manifests were not updated,
# Please manually delete them in $HOME/.docker/manifests/
# and re-run.
.PHONY: image.manifest.remove.%
image.manifest.remove.%:
	@rm -rf ${HOME}/.docker/manifests/docker.io_$(REGISTRY_PREFIX)_$(CI_PROJECT_NAME)-$(IMAGE)-$(CI_COMMIT_REF_NAME_FIX)-$(CI_PIPELINE_ID)

.PHONY: image.manifest.push.multiarch
image.manifest.push.multiarch: image.push.multiarch $(addprefix image.manifest.push.multiarch., $(IMAGES))

.PHONY: image.manifest.push.multiarch.%
image.manifest.push.multiarch.%:
	@echo "===========> Pushing manifest $* $(CI_PIPELINE_ID) to $(REGISTRY_PREFIX) and then remove the local manifest list"
	REGISTRY_PREFIX=$(REGISTRY_PREFIX) PLATFROMS="$(PLATFORMS)" IMAGE=$* CI_PIPELINE_ID=$(CI_PIPELINE_ID) CI_COMMIT_REF_NAME_FIX=$(CI_COMMIT_REF_NAME_FIX) CI_PROJECT_NAME=$(CI_PROJECT_NAME)  DOCKER_CLI_EXPERIMENTAL=enabled \
	  $(ROOT_DIR)/build/lib/create-manifest.sh -



.PHONY: image.buildx.push.multiarch
image.buildx.push.multiarch: $(addprefix image.buildx.push.multiarch., $(IMAGES))

.PHONY: image.buildx.push.multiarch.%
image.buildx.push.multiarch.%:
	@echo "===========> Pushing manifest $* $(CI_PIPELINE_ID) to $(REGISTRY_PREFIX) and then remove the local manifest list"
	$(MAKE) image.daemon.verify
	$(eval IMAGE := $*)
	@echo "===========> Building docker image $(IMAGE) $(CI_PIPELINE_ID) for $(PLATFORMS)"
	@mkdir -p $(TMP_DIR)/$(IMAGE)
	@cat $(ROOT_DIR)/build/docker/$(IMAGE)/Dockerfile | sed 's#FROM#FROM --platform=$$TARGETPLATFORM#g' | sed 's#$(IMAGE)#$${TARGETARCH}/$(IMAGE)#' | sed '1a ARG TARGETARCH'  > $(TMP_DIR)/$(IMAGE)/Dockerfile
	@$(foreach var, $(PLATFORMS), $(eval ARCH = $(word 2,$(subst _, ,$(var)))) $(eval IMAGE_PLATFORM = $(subst _,/,$(var))) mkdir -p $(TMP_DIR)/$(IMAGE)/$(ARCH); cp $(OUTPUT_DIR)/$(IMAGE_PLATFORM)/$(IMAGE) $(TMP_DIR)/$(IMAGE)/$(ARCH);)
	$(eval BUILDX_PLATFORMS := linux/arm64,linux/amd64)
	$(eval BUILDX_SUFFIX := --platform $(BUILDX_PLATFORMS) -t $(REGISTRY_PREFIX)/$(CI_PROJECT_NAME):$(IMAGE)-$(CI_COMMIT_REF_NAME_FIX)-$(CI_PIPELINE_ID) $(TMP_DIR)/$(IMAGE))
	$(DOCKER) buildx build $(BUILDX_SUFFIX) --push
	@rm -rf $(TMP_DIR)/$(IMAGE)
