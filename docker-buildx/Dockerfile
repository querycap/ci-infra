FROM docker

ARG TARGETARCH
ARG VERSION

ADD https://github.com/docker/buildx/releases/download/v${VERSION}/buildx-v${VERSION}.linux-${TARGETARCH} /root/.docker/cli-plugins/docker-buildx
RUN chmod a+x /root/.docker/cli-plugins/docker-buildx
ENV DOCKER_CLI_EXPERIMENTAL=enabled
