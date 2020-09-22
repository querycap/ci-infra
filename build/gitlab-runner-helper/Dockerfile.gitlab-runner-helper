ARG VERSION

FROM gitlab/gitlab-runner-helper:x86_64-v${VERSION} as base-amd64
FROM gitlab/gitlab-runner-helper:arm64-v${VERSION} as base-arm64

FROM base-${TARGETARCH}