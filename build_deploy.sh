#!/bin/bash

set -exv

# Initialize git submodules (pg_repack)
git submodule update --init --recursive

IMAGE="quay.io/cloudservices/vuln4shift-backend-app"
IMAGE_TAG=$(git rev-parse --short=7 HEAD)

if [[ -z "$QUAY_USER" || -z "$QUAY_TOKEN" ]]; then
    echo "QUAY_USER and QUAY_TOKEN must be set"
    exit 1
fi

if [[ -z "$RH_REGISTRY_USER" || -z "$RH_REGISTRY_TOKEN" ]]; then
    echo "RH_REGISTRY_USER and RH_REGISTRY_TOKEN  must be set"
    exit 1
fi

AUTH_CONF_DIR="$(pwd)/.podman"
mkdir -p $AUTH_CONF_DIR
export REGISTRY_AUTH_FILE="$AUTH_CONF_DIR/auth.json"

# Last tag if present
VERSION="$(git describe --tags --abbrev=0 || echo dev)"

podman login -u="$QUAY_USER" -p="$QUAY_TOKEN" quay.io
podman login -u="$RH_REGISTRY_USER" -p="$RH_REGISTRY_TOKEN" registry.redhat.io
podman build --build-arg="VERSION=$VERSION" --pull=true -f Dockerfile -t "${IMAGE}:${IMAGE_TAG}" .
podman push "${IMAGE}:${IMAGE_TAG}"
podman tag "${IMAGE}:${IMAGE_TAG}" "${IMAGE}:latest"
podman push "${IMAGE}:latest"
