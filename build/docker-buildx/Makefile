VERSION = $(shell cat .version)
NAME ?= docker-buildx
TAG ?= $(VERSION)

prepare: gomodcache
	@echo ::set-output name=image::$(NAME):$(TAG)
	@echo ::set-output name=build_args::VERSION=$(VERSION)

gomodcache:
	docker buildx build --push --build-arg VERSION=$(VERSION) -t localhost:5000/src:$(VERSION) -f ./src.Dockerfile .