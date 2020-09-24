VERSION = $(shell cat Dockerfile.version | grep "^FROM " | sed -e "s/FROM.*://g" )
TAG = $(VERSION)
HUB = querycap
NAME = buildkit
IMAGE = $(HUB)/$(NAME):$(TAG)

image:
	@echo $(IMAGE)

build:
	$(MAKE) -f ../Makefile buildx TAGS=$(IMAGE) BUILD_ARGS=VERSION=$(VERSION) DOCKERFILE=Dockerfile.$(NAME)
