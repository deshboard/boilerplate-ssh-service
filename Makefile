GLIDE:=$(shell if which glide > /dev/null 2>&1; then echo "glide"; fi)
BUILD_SERVICES=service test

# Setup environment
setup: build
	mkdir -p var/

# Build the service and test containers
build:
ifeq ($(FORCE), true)
	@docker-compose build --force-rm $(BUILD_SERVICES)
else
	@docker-compose build $(BUILD_SERVICES)
endif

# Start the environment
start:
	@docker-compose up -d

# Stop the environment
stop:
ifeq ($(FORCE), true)
	@docker-compose kill
else
	@docker-compose stop
endif

# Clean environment
clean: stop
	@rm -rf vendor/ var/
	@docker-compose rm --force
	@go clean

# Run test suite
test:
	@docker-compose run --rm test

# Install dependencies locally, optionally using go get
install:
ifdef GLIDE
	@$(GLIDE) install
else ifeq ($(FORCE), true)
	@go get
else
	@echo "Glide is necessary for installing project dependencies: http://glide.sh/ Run this command with FORCE=true to fall back to go get" 1>&2 && exit 1
endif

# Check that all source files follow the Coding Style
cs:
	@gofmt -l . | read something && echo "Code differs from gofmt's style" 1>&2 && exit 1 || true

# Fix Coding Standard violations
csfix:
	@gofmt -l -w -s .

.PHONY: setup build start stop clean test install cs csfix
