mkfile := $(abspath $(lastword $(MAKEFILE_LIST)))
dir := $(dir $(mkfile))

BUILD_DIR := build
DOCKER_REPO := digitalocean/artifactory-docker-resource
BINARY := check
EXTENSION :=

export LOG_TRUNCATE=true
export LOG_DIRECTORY=$(dir)

.PHONY: test
test:
	@go test --cover github.com/digitalocean/artifactory-docker-resource/...

.PHONY: gofmt
gofmt:
	@gofmt -w .

.PHONY: docker
docker:
	@docker build -t $(DOCKER_REPO):dev .

.PHONY: build
build: test
	make go-build BINARY=check --no-print-directory
	make go-build BINARY=get --no-print-directory
	make go-build BINARY=put --no-print-directory

.PHONY: go-build
go-build:
	@CGO_ENABLED=0 GOOS=${OS} GOOS=${ARCH} go build -o $(BUILD_DIR)/$(BINARY)$(EXTENSION) -ldflags="-s -w" -v cmd/$(BINARY)/main.go

.PHONY: docker-push
docker-push:
	@docker push $(DOCKER_REPO):dev
