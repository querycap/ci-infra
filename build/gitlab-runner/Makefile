VERSION = $(shell cat Dockerfile.version | grep "^FROM " | sed -e "s/FROM.*:alpine-v//g" )
TAG = $(VERSION)
HUB = querycap
NAME = gitlab-runner

build:
	$(MAKE) -f ../Makefile buildx TAGS=$(HUB)/$(NAME):$(TAG) BUILD_ARGS=VERSION=$(VERSION) DOCKERFILE=Dockerfile.$(NAME)