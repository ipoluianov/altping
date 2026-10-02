#!/usr/bin/env bash
# Builds altping for one platform into bin/.
# Usage: build/_build.sh <goos> <goarch>
set -euo pipefail

GOOS_T="$1"
GOARCH_T="$2"
APP=altping
REPO=ipoluianov/altping
NFPM_VERSION=v2.47.0
WINRES_VERSION=v0.3.3
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

VERSION="$(git describe --tags --always --dirty 2>/dev/null || echo dev)"
# Numeric part for the exe version info: v1.2.3-4-gabc -> 1.2.3
NUM_VERSION="$(printf '%s' "$VERSION" | sed -E 's/^v//; s/-.*//')"
[[ "$NUM_VERSION" =~ ^[0-9]+(\.[0-9]+)*$ ]] || NUM_VERSION="0.0.0"
LDFLAGS="-s -w -X github.com/ipoluianov/altping/app.Version=${VERSION}"
EXT=""
if [ "$GOOS_T" = "windows" ]; then
  EXT=".exe"
  LDFLAGS="${LDFLAGS} -H=windowsgui"
fi

OUT="bin/${APP}-${GOOS_T}-${GOARCH_T}${EXT}"
mkdir -p bin
echo "Building ${OUT} (${VERSION})"
if [ "$GOOS_T" = "windows" ]; then
  # Icon and version info shown by Explorer; go build links the .syso in
  trap 'rm -f "$ROOT"/rsrc_windows_*.syso' EXIT
  go run "github.com/tc-hib/go-winres@${WINRES_VERSION}" simply     --arch "$GOARCH_T" --out rsrc --manifest none --icon icon.png     --product-name AltPing --file-description AltPing --original-filename altping.exe     --copyright "Ivan Poluianov" --file-version "$NUM_VERSION" --product-version "$VERSION"
fi
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

  # Archive for scripts/linux-x64-install.sh: the binary and the menu icon
  cp icon.svg bin/.pkg/altping.svg
  chmod 755 bin/.pkg/altping
  tar -czf "bin/${APP}-${VERSION}-linux-${GOARCH_T}.tar.gz" -C bin/.pkg altping altping.svg
  rm -rf bin/.pkg

  # The installer downloads that archive from the release of this tag
  if [ "$GOARCH_T" = "amd64" ]; then
    tr -d '\r' < scripts/linux-x64-install.sh | sed \
      -e "s|__APP__|${APP}|g" \
      -e "s|__DISPLAY_NAME__|AltPing|g" \
      -e "s|__TAG__|${VERSION}|g" \
      -e "s|__REPO__|${REPO}|g" > bin/linux-x64-install.sh
    chmod 755 bin/linux-x64-install.sh
    [[ "$VERSION" =~ ^v[0-9.]+$ ]] ||
      echo "Warning: ${VERSION} is not a clean tag, linux-x64-install.sh points to a release that may not exist" >&2
  fi
fi
