NAME = docker
IMAGE = $(shell cat Dockerfile.$(NAME) | grep "^FROM " | sed -e "s/FROM //g" )
HUB ?= docker.io

echo:
	@echo $(IMAGE)

sync:
	docker buildx build \
    --push \
    --platform linux/arm64,linux/amd64 \
    --tag $(HUB)/$(IMAGE) \
    --file Dockerfile.$(NAME) .

