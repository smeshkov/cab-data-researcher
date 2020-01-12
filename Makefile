.PHONY: deps fmt clean build_darwin

APP_ENV?=local
VERSION?=$(shell git rev-list HEAD --max-count=1 --abbrev-commit)

deps:
	go get -u ./...

fmt:
    gofmt -w=true -s $(find . -type f -name '*.go' -not -path "./vendor/*")

clean: 
	rm -rf _dist/*
	
build:
	APP_ENV=$(APP_ENV) ./_bin/build.sh linux $(VERSION)

build_darwin:
	APP_ENV=$(APP_ENV) ./_bin/build.sh darwin $(VERSION)

test:
	./_bin/test.sh

run:
	go run cmd/app/app.go

up: build
	docker-compose up --build