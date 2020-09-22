VERSION = $(shell cat .version)
TAG = $(VERSION)
HUB = querycap
NAME = gitlab-runner-helper
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

ORIGIN = gitlab/gitlab-runner-helper

upgrade:
	echo ::set-output name=prev-version::$(VERSION)
	$(MAKE) upgrade-version VERSION=$(shell curl -sL "https://hub.docker.com/v2/repositories/$(ORIGIN)/tags?page=1&name=-v&ordering=last_updated" | jq -r ".results[].name" | grep arm64- | head -1 | sed -e "s/^arm64-//" )

upgrade-version: version
	echo "$(VERSION)" > .version
	echo ::set-output name=version::$(VERSION)