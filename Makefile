.PHONY: deps fmt clean build_darwin

deps:
	go get -u ./...

fmt:
    gofmt -w=true -s $(find . -type f -name '*.go' -not -path "./vendor/*")

clean: 
	rm -rf _dist/*
	
build:
	./_bin/build.sh linux

build_darwin:
	./_bin/build.sh darwin

test:
	./_bin/test.sh

run:
	go run cmd/app/app.go

up: build
	docker-compose up --build