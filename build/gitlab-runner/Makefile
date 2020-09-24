VERSION = $(shell cat Dockerfile.version | grep "^FROM " | sed -e "s/FROM.*:alpine-v//g" )
TAG = $(VERSION)
HUB = querycap
NAME = gitlab-runner

IMAGE = $(HUB)/$(NAME):$(TAG)

image:
	@echo $(IMAGE)

build:
	$(MAKE) -f ../Makefile buildx TAGS=$(IMAGE) BUILD_ARGS=VERSION=$(VERSION) DOCKERFILE=Dockerfile.$(NAME)