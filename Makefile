.PHONY: default clean checks test build

export GO111MODULE=on

TAG_NAME := $(shell git tag -l --contains HEAD)
SHA := $(shell git rev-parse --short HEAD)
VERSION := $(if $(TAG_NAME),$(TAG_NAME),$(SHA))

BUILD_DATE := $(shell date -u '+%Y-%m-%d_%I:%M:%S%p')

default: clean checks test build

test: clean
	go test -v -cover ./...

clean:
	rm -rf dist/ cover.out

build: clean
	@echo Version: $(VERSION) $(BUILD_DATE)
	go build -v -ldflags '-X "github.com/ldez/traefik-certs-dumper/cmd.version=${VERSION}" -X "github.com/ldez/traefik-certs-dumper/cmd.commit=${SHA}" -X "github.com/ldez/traefik-certs-dumper/cmd.date=${BUILD_DATE}"' -o traefik-certs-dumper

checks:
	golangci-lint run

doc:
	go run . doc
