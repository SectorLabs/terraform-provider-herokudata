#!/usr/bin/env bash

OSs=("darwin" "linux")
ARCHs=("amd64")

OUTPUT_DIR=$1
if [[ -z "${OUTPUT_DIR}" ]]; then
    OUTPUT_DIR=$(dirname $0)
fi

FILE_NAME="terraform-provider-herokudata"
VERSION=$(git tag --points-at HEAD)
if [[ -n "${VERSION}" ]]; then
    FILE_NAME="${FILE_NAME}_${VERSION}"
fi

for GOOS in "${OSs[@]}"; do
  for GOARCH in "${ARCHs[@]}"; do
    echo "Building $GOOS-$GOARCH"
    export GOOS=$GOOS
    export GOARCH=$GOARCH
    go build -o "${FILE_NAME}"
    gzip -c "${FILE_NAME}" > "${OUTPUT_DIR}/${FILE_NAME}_${GOOS}_${GOARCH}.tar.gz"
  done
done
