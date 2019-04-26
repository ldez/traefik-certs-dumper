#!/bin/bash
set -ex

# base docker image tag
TAG="andig/traefik-certs-dumper"

# only linux for now
OS=linux

# target platforms in docker manifest notation
declare -a PLATFORMS=( "amd64" "arm.v6" )

# images from Dockerfile
IMAGE=$(grep "{RUNTIME_HASH}" < Dockerfile | sed "s/FROM //" | sed 's/\$.*//')

# manifest cache file
MANIFEST_FILE=/tmp/manifest.$IMAGE.json

# get platform image hash from docker manifest
function hash () {
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
    docker pull "$IMAGE"
    docker manifest inspect "$IMAGE" > "$MANIFEST_FILE"
fi

# main
for platform in "${PLATFORMS[@]}"; do 
    # split architecture.version
    IFS='.' read -r ARCHITECTURE VARIANT <<< "$platform"

    # add xargs to trim whitespace
    RUNTIME_HASH=$(hash "$ARCHITECTURE" "$VARIANT")

    # arm architectures flavors, strip "v" prefix
    GOARM=${VARIANT:1}

    # build for target runtime image and architecture
    docker build --build-arg RUNTIME_HASH=@${RUNTIME_HASH} --build-arg GOARCH=${ARCHITECTURE} --build-arg GOARM=${GOARM} -t "$TAG:latest-$platform" .
done

# push images
for platform in "${PLATFORMS[@]}"; do 
    docker push "$TAG:latest-$platform"
done

# create manifest
TAG_LIST=$(printf "$TAG:latest-%s " "${PLATFORMS[@]}")
# shellcheck disable=SC2086
docker manifest create --amend "$TAG:latest" $TAG_LIST

for platform in "${PLATFORMS[@]}"; do 
    # split architecture.version
    IFS='.' read -r ARCHITECTURE VARIANT <<< "$platform"
    
    # docker and go architectures don't match
    if [ "arm" == "$ARCHITECTURE" ] && [ -n "$VARIANT" ]; then
        VARIANT="$ARCHITECTURE$VARIANT"
    fi

    docker manifest annotate "$TAG:latest" "$TAG:latest-$platform" --os "$OS" --arch "$ARCHITECTURE" --variant "$VARIANT"
done

docker manifest push "$TAG:latest"
