FROM docker

RUN apk add --update --no-cache curl make git bash && \
    mkdir -p ~/.docker/cli-plugins && cd ~/.docker/cli-plugins && \
    curl -sL --output docker-buildx https://github.com/docker/buildx/releases/download/v0.4.1/buildx-v0.4.1.linux-amd64 && \
    chmod a+x docker-buildx

ENV DOCKER_CLI_EXPERIMENTAL=enabled
