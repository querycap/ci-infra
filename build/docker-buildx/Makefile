VERSION = $(shell cat .version)
NAME = docker-buildx
TAG = $(VERSION)

prepare:
	echo ::set-output name=image::$(NAME):$(TAG)
	echo ::set-output name=build_args::VERSION=$(VERSION)