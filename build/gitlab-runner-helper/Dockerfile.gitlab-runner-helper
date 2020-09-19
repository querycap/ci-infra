ARG VERSION

FROM gitlab/gitlab-runner-helper:x86_64-${VERSION} as base-amd64
FROM gitlab/gitlab-runner-helper:arm64-${VERSION} as base-arm64

FROM base-${TARGETARCH}