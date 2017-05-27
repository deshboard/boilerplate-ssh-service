# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

include etc/variables.mk

# Build variables
VERSION ?= $(shell git rev-parse --abbrev-ref HEAD)
COMMIT_HASH = $(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_DATE = $(shell date +%FT%T%z)
LDFLAGS = -ldflags "-w -X ${PACKAGE}/app.Version=${VERSION} -X ${PACKAGE}/app.CommitHash=${COMMIT_HASH} -X ${PACKAGE}/app.BuildDate=${BUILD_DATE}"

# Docker variables
DOCKER_TAG ?= ${VERSION}
DOCKER_LATEST ?= false

# Dev variables
GO_SOURCE_FILES = $(shell find . -type f -name "*.go" -not -name "bindata.go" -not -path "./vendor/*")
GO_PACKAGES = $(shell go list ./... | grep -v /vendor/)
GODOTENV = $(shell if which godotenv > /dev/null 2>&1; then echo "godotenv"; fi)

.PHONY: setup install clean run watch build build-docker docker check test watch-test fmt csfix envcheck

setup:: install .env .env.test ## Setup the project for development

install: ## Install dependencies
	@glide install

.env: ## Create local env file
	cp .env.example .env

.env.test: ## Create local env file for running tests
	cp .env.example .env.test

clean:: ## Clean the working area
	rm -rf build/ vendor/ .env .env.test

run: build ## Build and execute a binary
	${GODOTENV} build/${BINARY_NAME} ${ARGS}

watch: ## Watch for file changes and run the built binary
	reflex -s -t 3s -d none -r '\.go$$' -- $(MAKE) ARGS="${ARGS}" run

build: ## Build a binary
	CGO_ENABLED=0 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}

build-docker:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-docker

docker: build-docker ## Build a Docker image
	docker build --build-arg BUILD_DIR=${BUILD_DIR} --build-arg BINARY_NAME=${BINARY_NAME}-docker -t ${DOCKER_IMAGE}:${DOCKER_TAG} .
ifeq (${DOCKER_LATEST}, true)
	docker tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${DOCKER_IMAGE}:latest
endif

check:: test cs ## Run tests and linters

test: .env.test ## Run unit tests
ifdef GODOTENV
	@${GODOTENV} -f .env.test go test ${ARGS} ${GO_PACKAGES}
else
	@go test ${ARGS} ${GO_PACKAGES}
endif

watch-test: ## Watch for file changes and run tests
	reflex -t 2s -d none -r '\.go$$' -- $(MAKE) ARGS="${ARGS}" test

cs: ## Check that all source files follow the Go coding style
	@gofmt -l ${GO_SOURCE_FILES} | read something && echo "Code differs from gofmt's style" 1>&2 && exit 1 || true

csfix: ## Fix Go coding style violations
	@gofmt -l -w -s ${GO_SOURCE_FILES}

envcheck:: ## Check environment for all the necessary requirements
	$(call executable_check,Go,go)
	$(call executable_check,Glide,glide)
	$(call executable_check,Docker,docker)
	$(call executable_check,Reflex,reflex)
	$(call executable_check,Godotenv,godotenv)

define executable_check
    @printf "\033[36m%-30s\033[0m %s\n" "$(1)" `if which $(2) > /dev/null 2>&1; then echo "\033[0;32m✓\033[0m"; else echo "\033[0;31m✗\033[0m"; fi`
endef

.PHONY: help
.DEFAULT_GOAL := help
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# Variable outputting/exporting rules
var-%: ; @echo $($*)
varexport-%: ; @echo $*=$($*)

include etc/service.mk
-include custom.mk
