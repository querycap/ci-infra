VERSION ?= $(shell cat Dockerfile.version | grep "^FROM " | sed -e "s/FROM.*:v//g")

DOCKERX_TAGS ?= $(VERSION)
DOCKERX_NAME ?= gitlab-runner-helper
DOCKERX_BUILD_ARGS ?= VERSION=$(VERSION)

include ../Makefile