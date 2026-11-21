# syntax=docker/dockerfile:1.7

# NOTE: This assumes you are running the build from the root, e.g.:
#       docker build -f ./services/stern/Dockerfile -t osutitanic/stern .

FROM --platform=$BUILDPLATFORM golang:1.27.1-alpine AS builder

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

RUN addgroup -g 10001 -S sailor && \
    adduser -u 10001 -S -D -G sailor -h /home/sailor sailor && \
    mkdir -p /out/data

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY internal ./internal
COPY services/stern ./services/stern

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOARM="" && \
    if [ "$TARGETARCH" = "arm" ] && [ "$TARGETVARIANT" = "v7" ]; then GOARM="7"; fi && \
    CGO_ENABLED=0 \
    GOOS="$TARGETOS" \
    GOARCH="$TARGETARCH" \
    GOARM="$GOARM" \
    go build -trimpath -ldflags="-s -w" -o /out/stern ./services/stern/cmd/web

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder --chown=10001:10001 /home/sailor /home/sailor

ENV HOME=/home/sailor

WORKDIR /app
COPY --from=builder --chown=10001:10001 /out/data /app/.data
COPY --from=builder /out/stern /app/stern

# NOTE: The assets directory is expected to be mounted as a volume
ENV STERN_STATIC_DIR=/app/static
VOLUME ["/app/static"]
VOLUME ["/app/.data"]

USER sailor:sailor

ENTRYPOINT ["/app/stern"]
