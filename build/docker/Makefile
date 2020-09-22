VERSION = $(shell cat Dockerfile.version | grep "^FROM " | sed -e "s/FROM.*://g" )
TAG = $(VERSION)
HUB = querycap
NAME = docker
TARGET_PLATFORM = linux/arm64,linux/amd64

build:
	docker buildx build \
		--push \
		--platform $(TARGET_PLATFORM) \
		--build-arg VERSION=$(VERSION) \
		--tag $(HUB)/$(NAME):$(TAG) \
		--file Dockerfile.$(NAME) .

version:
	@echo "VERSION: $(VERSION)"
