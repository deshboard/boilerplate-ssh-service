# Make variables
#
# This file contains project specific variables.

# Build variables
PACKAGE = $(shell go list .)
BINARY_NAME = $(shell echo ${PACKAGE} | cut -d '/' -f 3)
BUILD_DIR = build

# Docker variables
DOCKER_IMAGE ?= $(shell echo ${PACKAGE} | cut -d '/' -f 2,3)
