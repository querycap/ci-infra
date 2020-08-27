VERSION=$(shell cat .version)
HUB=querycap
IMAGE=${HUB}/gitlab-runner:${VERSION}

build:
	docker buildx build \
	--push \
	--platform linux/amd64,linux/arm64 \
	--build-arg=VERSION=${VERSION} \
	-t ${IMAGE} \
	.

sync:
	cd ../_sync/ && make -e IMAGE=${IMAGE} sync-image
