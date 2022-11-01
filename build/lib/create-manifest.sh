#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

PLATFORMS=${PLATFORMS:-"linux_amd64 linux_arm64"}

if [ -z ${IMAGE} ]; then
  echo "Please provide IMAGE."
  exit 1
fi

if [ -z ${VERSION} ]; then
  echo "Please provide VERSION."
  exit 1
fi

rm -rf ${HOME}/.docker/manifests/docker.io_${REGISTRY_PREFIX}_${CI_PROJECT_NAME}-${IMAGE}-${CI_COMMIT_REF_NAME_FIX}-${CI_PIPELINE_ID}
DES_REGISTRY=${REGISTRY_PREFIX}/${CI_PROJECT_NAME}
for platform in ${PLATFORMS}; do
  os=${platform%_*}
  arch=${platform#*_}
  variant=""
#  if [ ${arch} == "arm64" ]; then
#    variant="--variant unknown"
#  fi

  docker manifest create --amend ${DES_REGISTRY}:${IMAGE}-${CI_COMMIT_REF_NAME_FIX}-${CI_PIPELINE_ID} \
    ${DES_REGISTRY}-${arch}:${IMAGE}-${CI_COMMIT_REF_NAME_FIX}-${CI_PIPELINE_ID}

  docker manifest annotate ${DES_REGISTRY}:${IMAGE}-${CI_COMMIT_REF_NAME_FIX}-${CI_PIPELINE_ID} \
		${DES_REGISTRY}-${arch}:${IMAGE}-${CI_COMMIT_REF_NAME_FIX}-${CI_PIPELINE_ID} \
		--os ${os} --arch ${arch} ${variant}
done
docker manifest push --purge ${DES_REGISTRY}:${IMAGE}-${CI_COMMIT_REF_NAME_FIX}-${CI_PIPELINE_ID}
