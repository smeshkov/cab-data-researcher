#!/bin/sh

DIST_DIR="_dist"
BINARY="cabresearcher"
OS="$1"
VERSION="$2"

if [ -z "$OS" ]; then
    OS="darwin"
fi

if [ ! -z "$VERSION" ]; then
    VERSION="_$VERSION"
fi

if [ ! -d "$DIST_DIR" ]; then
  mkdir -p $DIST_DIR
fi

env GOOS=${OS} GOARCH=amd64 go build -i -v -o ${DIST_DIR}/${BINARY}_${OS}${VERSION} cmd/app/app.go