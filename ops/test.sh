#!/usr/bin/env sh

set -e

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
REPO_ROOT="$SCRIPT_DIR/.."
cd "$REPO_ROOT"

echo "Running tests for gonetia"
go test -count=1 ./...
echo "Tests completed"
