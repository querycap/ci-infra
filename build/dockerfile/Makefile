VERSION = $(shell cat Dockerfile.version | grep "^FROM " | sed -e "s/FROM.*://g" )
NAME = docker
TAG = $(VERSION)

prepare:
	echo ::set-output name=image::$(NAME):$(TAG)
	echo ::set-output name=build_args::VERSION=$(VERSION)