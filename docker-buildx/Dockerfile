FROM golang:1.15 as builder

ARG VERSION
RUN git clone --depth 1 -b kubernetes-driver-enhancement https://github.com/morlay/buildx.git /go/src/github.com/docker/buildx
WORKDIR /go/src/github.com/docker/buildx

RUN CGO_ENABLED=0 go build -o /go/bin/docker-buildx ./cmd/buildx

FROM docker
COPY --from=builder /go/bin/docker-buildx /root/.docker/cli-plugins/docker-buildx
ENV DOCKER_CLI_EXPERIMENTAL=enabled
