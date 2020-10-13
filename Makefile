HUB=ghcr.io/querycap/ci-infra docker.io/querycap

gen:
	HUB="$(HUB)" go run ./cmd/imagetools