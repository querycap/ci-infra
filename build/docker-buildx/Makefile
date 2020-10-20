VERSION = $(shell cat .version)
NAME ?= docker-buildx
TAG ?= $(VERSION)

prepare:
	@echo ::set-output name=image::$(NAME):$(TAG)
	@echo ::set-output name=build_args::VERSION=$(VERSION)

debug:
	docker buildx build --progress=plain --build-arg=VERSION=$(VERSION) --build-arg=GOPROXY=${GOPROXY} --push --platform=linux/arm64,linux/amd64 --tag $(NAME):$(TAG) -f Dockerfile.$(NAME) .
