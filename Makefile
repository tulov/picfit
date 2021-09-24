PROJECT_NAME ?= image_server
ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
VERSION=`git describe --tags`
PROJECT_NAMESPACE ?= tulovalex
REGISTRY_IMAGE ?= $(PROJECT_NAMESPACE)/$(PROJECT_NAME)

branch = $(shell git rev-parse --abbrev-ref HEAD)
commitMessage = $(shell git log -1 --pretty=%B)
commit = $(shell git log --pretty=format:'%h' -n 1)
now = $(shell date "+%Y-%m-%d %T UTC%z")
compiler = $(shell go version)

BIN_DIR = $(ROOT_DIR)/bin
PICFIT_CONFIG_PATH ?= `pwd`/config.json
BIN = $(BIN_DIR)/picfit
SSL_DIR = $(ROOT_DIR)/ssl
APP_DIR = /go/src/github.com/thoas/picfit

export GO111MODULE=on

test: unit

vendorize:
	find vendor/ -type f -not -path "*/.git*" -exec git add {} \;

run-server:
	@PICFIT_CONFIG_PATH=$(PICFIT_CONFIG_PATH) $(BIN)

serve:
	@modd

unit:
	go test -mod=vendor -v -cover ./...

all: picfit
	@(mkdir -p $(BIN_DIR))

build:
	@(echo "-> Compiling picfit binary")
	@(mkdir -p $(BIN_DIR))
	go build -mod=vendor -ldflags "\
		-X github.com/thoas/picfit/constants.Branch=$(branch) \
		-X github.com/thoas/picfit/constants.Revision=$(commit) \
		-X 'github.com/thoas/picfit/constants.BuildTime=$(now)' \
		-X 'github.com/thoas/picfit/constants.LatestCommitMessage=$(commitMessage)' \
		-X 'github.com/thoas/picfit/constants.Compiler=$(compiler)'" -o $(BIN_DIR)/picfit ./cmd/picfit/main.go
	@(echo "-> picfit binary created")

format:
	@(go fmt ./...)
	@(go vet ./...)

build-static:
	@(echo "-> Creating statically linked binary...")
	mkdir -p $(BIN_DIR)
	go build -mod=vendor -ldflags "\
		-X github.com/thoas/picfit/constants.Branch=$(branch) \
		-X github.com/thoas/picfit/constants.Revision=$(commit) \
		-X 'github.com/thoas/picfit/constants.BuildTime=$(now)' \
		-X 'github.com/thoas/picfit/constants.LatestCommitMessage=$(commitMessage)' \
		-X 'github.com/thoas/picfit/constants.Compiler=$(compiler)'" -a -installsuffix cgo -o $(BIN_DIR)/picfit ./cmd/picfit/main.go

docker-build-static: build-static


.PNONY: all test format

docker-build:
	@(echo "-> Preparing builder...")
	@(docker build -t picfit-builder -f Dockerfile.build .)
	@(mkdir -p $(BIN_DIR))
	@(docker run --rm -v $(BIN_DIR):$(APP_DIR)/bin picfit-builder)

dockerize:
	docker build -t $(PROJECT_NAME):$(VERSION) .

upload: dockerize
	docker tag $(PROJECT_NAME):$(VERSION) $(REGISTRY_IMAGE):$(VERSION)
	docker tag $(PROJECT_NAME):$(VERSION) $(REGISTRY_IMAGE):latest
	docker push $(REGISTRY_IMAGE):$(VERSION)
	docker push $(REGISTRY_IMAGE):latest
