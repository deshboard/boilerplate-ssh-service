# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

include .make/variables.mk

# Build variables
VERSION ?= $(shell git rev-parse --abbrev-ref HEAD)
COMMIT_HASH = $(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_DATE = $(shell date +%FT%T%z)
LDFLAGS = -ldflags "-w -X main.Version=${VERSION} -X main.CommitHash=${COMMIT_HASH} -X main.BuildDate=${BUILD_DATE}"

# Docker variables
DOCKER_TAG ?= ${VERSION}
DOCKER_LATEST ?= false

# Dev variables
GO_SOURCE_FILES = $(shell find . -type f -name "*.go" -not -name "bindata.go" -not -path "./vendor/*")
GO_PACKAGES = $(shell go list ./app/...)

.PHONY: setup
setup:: dep .env .env.test ## Setup the project for development

.PHONY: dep
dep: ## Install dependencies
	@glide install

.env: ## Create local env file
	cp .env.dist .env

.env.test: ## Create local env file for running tests
	cp .env.dist .env.test

.PHONY: clean
clean:: ## Clean the working area
	rm -rf ${BUILD_DIR}/ vendor/ .env .env.test

.PHONY: run
run: TAGS += dev
run: build .env ## Build and execute a binary
	${BUILD_DIR}/${BINARY_NAME} ${ARGS}

.PHONY: watch
watch: ## Watch for file changes and run the built binary
	reflex -s -t 3s -d none -r '\.go$$' -- $(MAKE) ARGS="${ARGS}" run

.PHONY: build
build: ## Build a binary
	CGO_ENABLED=0 go build -tags '${TAGS}' ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME} ${PACKAGE}/cmd

.PHONY: build-docker
build-docker:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-docker ${PACKAGE}/cmd

.PHONY: docker
docker: build-docker ## Build a Docker image
	docker build --build-arg BUILD_DIR=${BUILD_DIR} --build-arg BINARY_NAME=${BINARY_NAME}-docker -t ${DOCKER_IMAGE}:${DOCKER_TAG} .
ifeq (${DOCKER_LATEST}, true)
	docker tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${DOCKER_IMAGE}:latest
endif

.PHONY: check
check:: test cs ## Run tests and linters

PASS=$(shell printf "\033[32mPASS\033[0m")
FAIL=$(shell printf "\033[31mFAIL\033[0m")
COLORIZE=sed ''/PASS/s//${PASS}/'' | sed ''/FAIL/s//${FAIL}/''

.PHONY: test
test: .env.test ## Run unit tests
	@go test -tags '${TAGS}' ${ARGS} ${GO_PACKAGES} | ${COLORIZE}

.PHONY: watch-test
watch-test: ## Watch for file changes and run tests
	reflex -t 2s -d none -r '\.go$$' -- $(MAKE) ARGS="${ARGS}" test

.PHONY: cs
cs: ## Check that all source files follow the Go coding style
	@gofmt -l ${GO_SOURCE_FILES} | read something && echo "Code differs from gofmt's style" 1>&2 && exit 1 || true

.PHONY: csfix
csfix: ## Fix Go coding style violations
	@gofmt -l -w -s ${GO_SOURCE_FILES}

.PHONY: envcheck
envcheck:: ## Check environment for all the necessary requirements
	$(call executable_check,Go,go)
	$(call executable_check,Glide,glide)
	$(call executable_check,Docker,docker)
	$(call executable_check,Docker Compose,docker-compose)
	$(call executable_check,Reflex,reflex)

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

include .make/service.mk
-include .make/custom.mk
