# syntax=docker/dockerfile:1.7

FROM --platform=$TARGETPLATFORM debian:trixie-slim AS builder

ARG TARGETARCH
ARG FFMPEG_REFRESH=default

WORKDIR /src

RUN <<'EOF'
set -eu

# Refresh debian package metadata when ci invalidates the cache
echo "ffmpeg refresh: $FFMPEG_REFRESH"
sed -i 's/^Types: deb$/Types: deb deb-src/' /etc/apt/sources.list.d/debian.sources
apt-get update

# Install build dependencies
set -- build-essential ca-certificates dpkg-dev libmp3lame-dev perl pkg-config xz-utils
if [ "$TARGETARCH" = "amd64" ]; then
    set -- "$@" nasm
fi
DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends "$@"

# Download and extract the Debian ffmpeg sources
apt-get source ffmpeg
cd ffmpeg-*

# Only enable the audio decoders required by deck
# TODO: Check if we catched all of them
decoders="flac,mp3,vorbis,\
pcm_alaw,pcm_mulaw,pcm_s8,pcm_u8,\
pcm_s16be,pcm_s16le,pcm_u16be,pcm_u16le,\
pcm_s24be,pcm_s24le,pcm_u24be,pcm_u24le,\
pcm_s32be,pcm_s32le,pcm_u32be,pcm_u32le,\
pcm_s64be,pcm_s64le,pcm_f32be,pcm_f32le,\
pcm_f64be,pcm_f64le"

./configure \
    --prefix=/usr \
    --disable-autodetect \
    --disable-doc \
    --disable-debug \
    --disable-network \
    --disable-everything \
    --disable-avdevice \
    --disable-postproc \
    --disable-swscale \
    --enable-small \
    --enable-static \
    --disable-shared \
    --pkg-config-flags=--static \
    --extra-ldflags=-static \
    --enable-ffmpeg \
    --enable-ffprobe \
    --enable-protocol=file,pipe \
    --enable-demuxer=flac,mp3,ogg,wav \
    --enable-muxer=mp3 \
    --enable-decoder="$decoders" \
    --enable-encoder=libmp3lame \
    --enable-parser=flac,mpegaudio,vorbis \
    --enable-filter=aformat,anull,aresample,atrim \
    --enable-libmp3lame

# Build and package the binaries
make -j"$(nproc)" ffmpeg ffprobe
strip --strip-unneeded ffmpeg ffprobe

mkdir -p /out/usr/bin /out/usr/share/ffmpeg
cp ffmpeg ffprobe /out/usr/bin/
dpkg-parsechangelog --show-field Version > /out/usr/share/ffmpeg/version

# TODO: perhaps add testing to this once the build is done
EOF

FROM scratch

LABEL org.opencontainers.image.source="https://github.com/osuTitanic/titanic"

COPY --from=builder /out/ /

ENTRYPOINT ["/usr/bin/ffmpeg"]
