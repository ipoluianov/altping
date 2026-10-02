#!/usr/bin/env bash
# Builds altping for one platform into bin/.
# Usage: build/_build.sh <goos> <goarch>
set -euo pipefail

GOOS_T="$1"
GOARCH_T="$2"
APP=altping
NFPM_VERSION=v2.47.0
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

VERSION="$(git describe --tags --always --dirty 2>/dev/null || echo dev)"
LDFLAGS="-s -w -X github.com/ipoluianov/altping/app.Version=${VERSION}"
EXT=""
if [ "$GOOS_T" = "windows" ]; then
  EXT=".exe"
  LDFLAGS="${LDFLAGS} -H=windowsgui"
fi

OUT="bin/${APP}-${GOOS_T}-${GOARCH_T}${EXT}"
mkdir -p bin
echo "Building ${OUT} (${VERSION})"
CGO_ENABLED=0 GOOS="$GOOS_T" GOARCH="$GOARCH_T" \
  go build -trimpath -ldflags="${LDFLAGS}" -o "$OUT" .

if [ "$GOOS_T" = "darwin" ]; then
  "$ROOT/build/_dmg.sh" "$OUT" "$VERSION"
fi

if [ "$GOOS_T" = "linux" ]; then
  # nfpm is pure Go, so packages can be built on any OS
  mkdir -p bin/.pkg
  cp "$OUT" bin/.pkg/altping
  for fmt in deb rpm; do
    PKG_ARCH="$GOARCH_T" PKG_VERSION="$VERSION" \
      go run "github.com/goreleaser/nfpm/v2/cmd/nfpm@${NFPM_VERSION}" \
      pkg --config build/nfpm.yaml --packager "$fmt" --target "${OUT}.${fmt}"
  done
  rm -rf bin/.pkg
fi
