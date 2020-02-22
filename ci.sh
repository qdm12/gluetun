#!/bin/bash

if [ "$TRAVIS_PULL_REQUEST" = "true" ]; then
  echo "Building pull request without pushing to Docker Hub"
  docker buildx build \
    --progress plain \
    --platform=linux/amd64,linux/386,linux/arm64,linux/arm/v7,linux/arm/v6,linux/ppc64le,linux/s390x \
    .
  exit $?
fi
echo $DOCKER_PASSWORD | docker login -u qmcgaw --password-stdin &> /dev/null
TAG="$TRAVIS_TAG"
if [ -z $TAG ]; then
  if [ "$TRAVIS_BRANCH" = "master" ]; then
    TAG=latest
  else
    TAG="$TRAVIS_BRANCH"
  fi
fi

echo "Building Docker images for \"$DOCKER_REPO:$TAG\""
docker buildx build \
    --progress plain \
    --platform=linux/amd64,linux/386,linux/arm64,linux/arm/v7,linux/arm/v6,linux/ppc64le,linux/s390x \
    --build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
    --build-arg VCS_REF=`git rev-parse --short HEAD` \
    --build-arg VERSION=$TAG \
    -t $DOCKER_REPO:$TAG \
    --push \
    .
