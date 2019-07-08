COMMIT=$(shell git rev-parse HEAD)
VERSION=$(shell git describe --tags --exact-match --always)
DATE=$(shell date +'%FT%TZ%z')

build:
	CGO_ENABLED=0 go build -v -a \
		-ldflags '-s -w -extldflags "-static" -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.Date=$(DATE)' \
		-o bin/trunks
