# syntax=docker/dockerfile:1.7

# NOTE: This assumes you are running the build from the root, e.g.:
#       docker build -f ./services/jobs/Dockerfile -t osutitanic/jobs .

FROM --platform=$BUILDPLATFORM golang:1.24.6-alpine AS builder

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOARM="" && \
    if [ "$TARGETARCH" = "arm" ] && [ "$TARGETVARIANT" = "v7" ]; then GOARM="7"; fi && \
    CGO_ENABLED=0 \
    GOOS="$TARGETOS" \
    GOARCH="$TARGETARCH" \
    GOARM="$GOARM" \
    go build -trimpath -ldflags="-s -w" -o /out/jobs ./services/jobs/cmd/cli

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app
COPY --from=builder /out/jobs /app/jobs
COPY ./services/jobs/schedule.example.json /app/schedule.json

VOLUME ["/app/.data"]
ENTRYPOINT ["/app/jobs"]
CMD ["-file", "/app/schedule.json"]
