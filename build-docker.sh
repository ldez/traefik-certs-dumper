#!/usr/bin/env bash

set -o errexit
set -o pipefail

VERSION=v666

# safe guard
#if [ -n "$TRAVIS_TAG" ] && [ -n "$VERSION" ]; then
#  echo "Deploying..."
#else
#  echo "Skipping deploy"
#  exit 0
#fi

# base docker image name
IMAGE_NAME="ldez/traefik-certs-dumper"

# only linux for now
OS=linux

# target platforms in docker manifest notation
declare -a PLATFORMS=( "amd64" "arm.v6" "arm.v7" "arm64")

# images from Dockerfile
FROM_IMAGE=$(grep "{RUNTIME_HASH}" < Dockerfile | sed "s/FROM //" | sed 's/\$.*//')

# manifest cache file
MANIFEST_FILE=/tmp/tcd-manifest.${FROM_IMAGE}.json

# get platform image hash from docker manifest
function platformHash () {
    local ARCHITECTURE VARIANT HASH
    read -r ARCHITECTURE VARIANT <<< "$@"

    if [ -z "$VARIANT" ]; then
        HASH=$(jq -r ".manifests[] | select(.platform.architecture == \"$ARCHITECTURE\") | .digest" < "$MANIFEST_FILE")
    else
        HASH=$(jq -r ".manifests[] | select(.platform.architecture == \"$ARCHITECTURE\" and .platform.variant == \"$VARIANT\") | .digest" < "$MANIFEST_FILE")
    fi

    echo "$HASH"
}

# get manifest
if [ ! -f "$MANIFEST_FILE" ]; then
    docker pull "$FROM_IMAGE"
    echo "docker manifest inspect $FROM_IMAGE"
    DOCKER_CLI_EXPERIMENTAL=enabled docker manifest inspect "$FROM_IMAGE" > "$MANIFEST_FILE"
fi

# create and push images
for platform in "${PLATFORMS[@]}"; do
    # split architecture.version
    IFS='.' read -r ARCHITECTURE VARIANT <<< "$platform"

    # add xargs to trim whitespace
    RUNTIME_HASH=$(platformHash "$ARCHITECTURE" "$VARIANT")

    # arm architectures flavors, strip "v" prefix
    GOARM=${VARIANT:1}

    # build for target runtime image and architecture
#    echo "docker build --build-arg=RUNTIME_HASH=@${RUNTIME_HASH} --build-arg=GOARCH=${ARCHITECTURE} --build-arg=GOARM=${GOARM} -t $IMAGE_NAME:${VERSION}-$platform" .
    docker build --build-arg="RUNTIME_HASH=@${RUNTIME_HASH}" --build-arg="GOARCH=${ARCHITECTURE}" --build-arg="GOARM=${GOARM}" -t "$IMAGE_NAME:${VERSION}-$platform" .

    # push images
    echo "docker push $IMAGE_NAME:${VERSION}-$platform"
#    docker push "$IMAGE_NAME:${VERSION}-$platform"
done

# create manifest
TAG_LIST=$(printf "$IMAGE_NAME:${VERSION}-%s " "${PLATFORMS[@]}")
# shellcheck disable=SC2086
echo "docker manifest create --amend $IMAGE_NAME:${VERSION} $TAG_LIST"
#DOCKER_CLI_EXPERIMENTAL=enabled docker manifest create --amend "$IMAGE_NAME:${VERSION}" $TAG_LIST

for platform in "${PLATFORMS[@]}"; do
    # split architecture.version
    IFS='.' read -r ARCHITECTURE VARIANT <<< "$platform"

    echo "docker manifest annotate $IMAGE_NAME:${VERSION} $IMAGE_NAME:${VERSION}-$platform --os $OS --arch $ARCHITECTURE --variant $VARIANT"
#    DOCKER_CLI_EXPERIMENTAL=enabled docker manifest annotate "$IMAGE_NAME:${VERSION}" "$IMAGE_NAME:${VERSION}-$platform" --os "$OS" --arch "$ARCHITECTURE" --variant "$VARIANT"
done

# push manifest
echo "docker manifest push $IMAGE_NAME:${VERSION}"
#DOCKER_CLI_EXPERIMENTAL=enabled docker manifest push "$IMAGE_NAME:${VERSION}"
