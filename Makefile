.PHONY: tidy lint vet staticcheck check build test

VERSION:=$(shell git describe --tags)

tidy:
	go mod tidy -go=1.22.0

lint:
	revive -config revive/config.toml -formatter friendly ./...

vet:
	go vet ./...

staticcheck:
	staticcheck ./...

check: tidy lint vet staticcheck

build:
	go build -ldflags "-X main.version=${VERSION}" -o ov

test:
	go test -v ./... -race -covermode=atomic -shuffle=on
