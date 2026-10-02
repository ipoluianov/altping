#!/usr/bin/env bash
# Builds altping for every supported platform into bin/.
set -euo pipefail
DIR="$(dirname "$0")"
for t in linux/amd64 linux/arm64 windows/amd64 windows/arm64 darwin/amd64 darwin/arm64; do
  "$DIR/_build.sh" "${t%/*}" "${t#*/}"
done
