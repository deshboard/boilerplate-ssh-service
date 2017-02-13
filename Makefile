GLIDE:=$(shell if which glide > /dev/null 2>&1; then echo "glide"; fi)
GO_SOURCE_FILES=$(shell find . -type f -name "*.go" -not -name "bindata.go" -not -path "./vendor/*")

# Setup environment
setup: install
	mkdir -p var/

# Install dependencies locally, optionally using go get
install:
ifdef GLIDE
	@$(GLIDE) install
else ifeq ($(FORCE), true)
	@go get
else
	@echo "Glide is necessary for installing project dependencies: http://glide.sh/ Run this command with FORCE=true to fall back to go get" 1>&2 && exit 1
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
ifeq ($(INTEGRATION), true)
	@docker-compose -f docker-compose.test.yml run --rm test -v -tags=integration
else
	@docker-compose -f docker-compose.test.yml run --rm test -v
endif

# Check that all source files follow the Coding Style
cs:
	@gofmt -l $(GO_SOURCE_FILES) | read something && echo "Code differs from gofmt's style" 1>&2 && exit 1 || true

# Fix Coding Standard violations
csfix:
	@gofmt -l -w -s $(GO_SOURCE_FILES)

.PHONY: setup install start stop clean test cs csfix
