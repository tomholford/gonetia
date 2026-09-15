#!/usr/bin/env sh

set -e

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
REPO_ROOT="$SCRIPT_DIR/.."
cd "$REPO_ROOT"

echo "Linting $(go list ./... | wc -w | tr -d ' ') packages"

if ! [ -x "$(command -v golangci-lint)" ]; then
  echo "golangci-lint is not installed. Please install it from https://golangci-lint.run/welcome/install/"
  exit 1
fi

golangci-lint run
