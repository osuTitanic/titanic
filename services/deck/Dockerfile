# syntax=docker/dockerfile:1.7

# NOTE: This assumes you are running the build from the root, e.g.:
#       docker build -f ./services/deck/Dockerfile -t osutitanic/deck .

FROM --platform=$BUILDPLATFORM golang:1.27.1-trixie AS builder

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        gcc \
        crossbuild-essential-arm64 \
        crossbuild-essential-armhf && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY internal ./internal
COPY services/deck ./services/deck

# This shit is a hot mess to look at, but I don't care because it works

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    # Figure out which cross-compiler to use
    case "$TARGETARCH/$TARGETVARIANT" in \
        amd64/) \
            CC="gcc"; \
            GOARM=""; \
            NATIVE_ARCH="linux-x64"; \
            STRIP="strip" \
            ;; \
        arm64/) \
            CC="aarch64-linux-gnu-gcc"; \
            GOARM=""; \
            NATIVE_ARCH="linux-arm64"; \
            STRIP="aarch64-linux-gnu-strip" \
            ;; \
        arm/v7) \
            CC="arm-linux-gnueabihf-gcc"; \
            GOARM="7"; \
            NATIVE_ARCH="linux-arm"; \
            STRIP="arm-linux-gnueabihf-strip" \
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
        -o /out/deck \
        ./services/deck/cmd/api && \
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
    cp "$NATIVE_LIBRARY" /out/lib/libosu.Native.so && \
    "$STRIP" --strip-unneeded /out/lib/libosu.Native.so

FROM debian:trixie-slim AS runtime-deps

# Collect FFmpeg, FFprobe & their shared libraries
RUN <<'EOF'
set -eu

apt-get update
apt-get install -y --no-install-recommends ffmpeg

mkdir -p /out/bin /out/lib

# Copy the FFmpeg executables
cp /usr/bin/ffmpeg /usr/bin/ffprobe /out/bin/

# Copy shared libraries required by ffmpeg / ffprobe
# ldd resolves dependencies to their absolute paths
ldd /usr/bin/ffmpeg /usr/bin/ffprobe \
    | awk '
        $2 == "=>" && $3 ~ /^\// { print $3 }
        $1 ~ /^\// && $1 !~ /:$/ { print $1 }
    ' \
    | sort -u \
    | while IFS= read -r lib; do
        cp -L "$lib" /out/lib/
    done

# Remove cached package metadata to reduce image size.
rm -rf /var/lib/apt/lists/*
EOF

FROM gcr.io/distroless/base-nossl-debian13:latest AS runtime

WORKDIR /app
COPY --from=builder /out/deck /app/deck
COPY --from=builder /out/lib/libosu.Native.so /app/lib/libosu.Native.so
COPY --from=runtime-deps /out/bin/ /usr/bin/
COPY --from=runtime-deps /out/lib/ /app/lib/

ENV LD_LIBRARY_PATH=/app/lib

VOLUME ["/app/.data"]
ENTRYPOINT ["/app/deck"]
