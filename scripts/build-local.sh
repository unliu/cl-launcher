#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUTPUT="${1:-"$ROOT_DIR/dist/local/cl-dev"}"

mkdir -p "$(dirname "$OUTPUT")" "$ROOT_DIR/.gocache"

(
  cd "$ROOT_DIR"
  GOFLAGS="${GOFLAGS:-} -modcacherw" GOCACHE="$ROOT_DIR/.gocache" go test ./...
  GOFLAGS="${GOFLAGS:-} -modcacherw" GOCACHE="$ROOT_DIR/.gocache" go build \
    -buildvcs=false \
    -ldflags "-s -w -X main.version=local" \
    -o "$OUTPUT" \
    .
)

echo "Built local test binary: $OUTPUT"
echo "Try it with: $OUTPUT <profile> [args...]"
