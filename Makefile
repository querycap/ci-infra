build:
	docker buildx build --platform linux/amd64,linux/arm64 -t querycap/docker-buildx:latest --push .