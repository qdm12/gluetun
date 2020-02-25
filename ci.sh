#!/bin/bash

ARCHS=linux/amd64,linux/386,linux/arm64,linux/arm/v7,linux/arm/v6,linux/ppc64le,linux/s390x

echo -e "\n\nPull request: $TRAVIS_PULL_REQUEST\nRelease tag: $TRAVIS_TAG\nBranch: $TRAVIS_BRANCH\n\nTarget arch: $ARCHS\n\n"

if [  "$TRAVIS_PULL_REQUEST" != "false" ]; then
  echo -e "\n\nBuilding pull request without pushing to Docker Hub\n\n"
  docker buildx build \
    --progress plain \
    --platform="$ARCHS" \
    .
  exit $?
fi

echo $DOCKER_PASSWORD | docker login -u qmcgaw --password-stdin 2>&1

TAG="$TRAVIS_TAG"
if [ -z "$TAG" ]; then
  TAG=latest
  if [ "$TRAVIS_BRANCH" != "master" ]; then
    TAG="$TRAVIS_BRANCH"
  fi
fi

echo -e "\n\nBuilding Docker images for \"$DOCKER_REPO:$TAG\"\n\n"
docker buildx build \
    --progress plain \
    --platform="$ARCHS" \
    --build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
    --build-arg VCS_REF=`git rev-parse --short HEAD` \
    --build-arg VERSION=$TAG \
    -t $DOCKER_REPO:$TAG \
    --push \
    .
