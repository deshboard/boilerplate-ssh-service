# Make variables
#
# This file contains project specific variables.

# Build variables
PACKAGE = $(shell go list .)
BINARY_NAME = $(shell echo ${PACKAGE} | cut -d '/' -f 3)

# Docker variables
DOCKER_IMAGE ?= $(shell echo ${PACKAGE} | cut -d '/' -f 2,3)
