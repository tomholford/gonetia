#!/usr/bin/env sh

set -e

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
REPO_ROOT="$SCRIPT_DIR/.."
cd "$REPO_ROOT"

mkdir -p bin

VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo dev)

echo "Building gonetia ${VERSION}..."
go build -ldflags "-s -w -X main.Version=${VERSION}" -o bin/gonetia .
echo "Done → bin/gonetia"
