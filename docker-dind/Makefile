VERSION=dind
HUB=querycap
IMAGE=${HUB}/docker:${VERSION}

build:
	docker buildx build \
	--push \
	--platform linux/amd64,linux/arm64 \
	-t ${IMAGE} \
	.

sync:
	cd ../_sync/ && make -e IMAGE=${IMAGE} sync-image
