VERSION=$(shell cat .version)
HUB=querycap
HUB2=registry.cn-hangzhou.aliyuncs.com/querycap

build:
	docker buildx build \
	--push \
	--platform linux/amd64,linux/arm64 \
	--build-arg=VERSION=${VERSION} \
	-t ${HUB}/docker-buildx:${VERSION} \
	-t ${HUB2}/docker-buildx:${VERSION} \
	.