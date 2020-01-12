#!/bin/sh

DIST_DIR="_dist"
BINARY="cabresearcher"
OS="$1"
VERSION="$2"

if [ -z "$OS" ]; then
    OS="darwin"
fi

if [ ! -d "$DIST_DIR" ]; then
  mkdir -p $DIST_DIR
fi

if [ -z "$APP_ENV" ]; then
  APP_ENV="local"
fi

env GOOS=${OS} GOARCH=amd64 go build -ldflags "-X main.env=$APP_ENV -X main.version=$VERSION" -i -v -o ${DIST_DIR}/${BINARY}_${OS} cmd/app/app.go