VERSION = $(shell cat .version)
TAG = $(VERSION)
HUB = querycap
NAME = docker-buildx
TARGET_PLATFORM = linux/arm64,linux/amd64

build:
	$(MAKE) -f ../Makefile buildx TAGS=$(HUB)/$(NAME):$(TAG) BUILD_ARGS=VERSION=$(VERSION) DOCKERFILE=Dockerfile.$(NAME)
