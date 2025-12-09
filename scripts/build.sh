#!/usr/bin/env bash
set -euo pipefail

mkdir -p dist

targets=(
  "linux amd64"
  "linux arm64"
  "darwin amd64"
  "darwin arm64"
  "windows amd64"
)

for t in "${targets[@]}"; do
  read -r GOOS GOARCH <<<"$t"
  out="dist/opsdoctor-${GOOS}-${GOARCH}"
  if [[ "$GOOS" == "windows" ]]; then
    out="${out}.exe"
  fi
  echo "Building $out"
  GOOS=$GOOS GOARCH=$GOARCH go build -o "$out" ./...
done

echo "Built artifacts in dist/"
