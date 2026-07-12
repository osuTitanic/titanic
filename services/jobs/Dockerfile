# syntax=docker/dockerfile:1.7

# NOTE: This assumes you are running the build from the root, e.g.:
#       docker build -f ./services/jobs/Dockerfile -t osutitanic/jobs .

FROM --platform=$BUILDPLATFORM golang:1.26-trixie AS builder

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        gcc \
        gcc-aarch64-linux-gnu \
        gcc-arm-linux-gnueabihf && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

# This shit is a hot mess to look at, but I don't care because it works

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    # Figure out which cross-compiler to use
    case "$TARGETARCH/$TARGETVARIANT" in \
        amd64/) \
            CC="gcc"; \
            GOARM=""; \
            NATIVE_ARCH="linux-x64" \
            ;; \
        arm64/) \
            CC="aarch64-linux-gnu-gcc"; \
            GOARM=""; \
            NATIVE_ARCH="linux-arm64" \
            ;; \
        arm/v7) \
            CC="arm-linux-gnueabihf-gcc"; \
            GOARM="7"; \
            NATIVE_ARCH="linux-arm" \
            ;; \
        *) \
            echo "Unsupported target: $TARGETARCH/$TARGETVARIANT" >&2; \
            exit 1 \
            ;; \
    esac && \
    # Build the binary with those settings
    CGO_ENABLED=1 \
    CC="$CC" \
    GOOS="$TARGETOS" \
    GOARCH="$TARGETARCH" \
    GOARM="$GOARM" \
    CGO_LDFLAGS="-lm" \
    go build \
        -tags=native \
        -trimpath \
        -ldflags="-s -w" \
        -o /out/jobs \
        ./services/jobs/cmd/cli && \
    # Build done -> copy osu-native .so file so that we can use it later
    mkdir -p /out/lib && \
    NATIVE_LIBRARY="$(find /go/pkg/mod/github.com/7mochi/osu-native-go@*/native/bin/"$NATIVE_ARCH" \
        -maxdepth 1 \
        -name 'libosu.Native.so' \
        -print \
        -quit)" && \
    if [ -z "$NATIVE_LIBRARY" ]; then \
        echo "libosu.Native.so not found for $NATIVE_ARCH" >&2; \
        exit 1; \
    fi && \
    cp "$NATIVE_LIBRARY" /out/lib/libosu.Native.so

FROM gcr.io/distroless/base-debian13:nonroot

WORKDIR /app
COPY --from=builder /out/jobs /app/jobs
COPY ./services/jobs/schedule.example.json /app/schedule.json

COPY --from=builder /out/lib/libosu.Native.so /app/lib/libosu.Native.so
ENV LD_LIBRARY_PATH=/app/lib

VOLUME ["/app/.data"]
ENTRYPOINT ["/app/jobs"]
CMD ["-file", "/app/schedule.json"]
