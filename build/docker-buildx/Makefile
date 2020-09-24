VERSION = $(shell cat .version)
TAG = $(VERSION)
HUB = querycap
NAME = docker-buildx
IMAGE = $(HUB)/$(NAME):$(TAG)

image:
	@echo $(IMAGE)

build:
	$(MAKE) -f ../Makefile buildx TAGS=$(IMAGE) BUILD_ARGS=VERSION=$(VERSION) DOCKERFILE=Dockerfile.$(NAME)
