#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

PLATFORMS=${PLATFORMS:-"linux_amd64 linux_arm64"}

if [ -z ${IMAGE} ]; then
  echo "Please provide IMAGE."
  exit 1
fi

if [ -z ${CI_PIPELINE_ID} ]; then
  echo "Please provide CI_PIPELINE_ID."
  exit 1
fi

rm -rf ${HOME}/.docker/manifests/docker.io_${REGISTRY_PREFIX}_${CI_PROJECT_NAME}-${IMAGE}-${CI_COMMIT_REF_NAME_FIX}-${CI_PIPELINE_ID}
DES_REGISTRY=${REGISTRY_PREFIX}/${CI_PROJECT_NAME}
for platform in ${PLATFORMS}; do
  os=${platform%_*}
  arch=${platform#*_}

  docker manifest create --amend --insecure ${DES_REGISTRY}:${IMAGE}-${CI_COMMIT_REF_NAME_FIX}-${CI_PIPELINE_ID} \
    ${DES_REGISTRY}-${arch}:${IMAGE}-${CI_COMMIT_REF_NAME_FIX}-${CI_PIPELINE_ID}

  docker manifest annotate ${DES_REGISTRY}:${IMAGE}-${CI_COMMIT_REF_NAME_FIX}-${CI_PIPELINE_ID} \
		${DES_REGISTRY}-${arch}:${IMAGE}-${CI_COMMIT_REF_NAME_FIX}-${CI_PIPELINE_ID} \
		--os ${os} --arch ${arch}
done
docker manifest push --purge --insecure ${DES_REGISTRY}:${IMAGE}-${CI_COMMIT_REF_NAME_FIX}-${CI_PIPELINE_ID}
