# syntax = docker/dockerfile:experimental

FROM golang:1.15 as src

ARG VERSION=master
RUN git clone --depth 1 -b ${VERSION} https://github.com/docker/buildx.git /go/src/

WORKDIR /go/src/

RUN --mount=type=cache,id=gomod,target=/go/pkg/mod go mod download -x

FROM busybox

COPY --from=src /go/src /go/src