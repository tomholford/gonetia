#!/usr/bin/env sh
# Pull the latest urbit-wordlists submodule (source for go:embed).

set -e

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
REPO_ROOT="$SCRIPT_DIR/.."
cd "$REPO_ROOT"

git submodule foreach git pull origin master
