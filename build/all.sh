#!/usr/bin/env bash
# Builds altping for every supported platform into one bin/<yyyy-mm-dd-HH-MM>-<version>/.
set -euo pipefail
DIR="$(dirname "$0")"
BUILD_STAMP="$(date +%Y-%m-%d-%H-%M)"
export BUILD_STAMP
for t in linux/amd64 linux/arm64 windows/amd64 darwin/arm64; do
  "$DIR/_build.sh" "${t%/*}" "${t#*/}"
done
# Notarize the signed dmg with Apple (a few minutes); macOS only, the dmg
# made on other systems is unsigned
if [ "$(uname)" = Darwin ]; then
  "$DIR/notarize.sh" "$DIR/../bin/${BUILD_STAMP}"-*/*.dmg
fi
# The README screenshot on model data, in a virtual X server; skipped where it cannot run
"$DIR/_screenshot.sh" || echo "Screenshot failed, the builds are done"
