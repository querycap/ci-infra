HUB=ghcr.io/querycap/ci-infra docker.io/querycap

gen:
	HUB="$(HUB)" go run ./cmd/imagetools

word-dot = $(word $2,$(subst ., ,$1))

# gn.gn
dockerx.%:
	$(MAKE) -C build/$(call word-dot,$*,1) dockerx HUB="$(HUB)" DOCKERX_NAME=$(call word-dot,$*,2)

imagetools.%:
	$(MAKE) -C build/$(call word-dot,$*,1) imagetools HUB="$(HUB)" DOCKERX_NAME=$(call word-dot,$*,2)