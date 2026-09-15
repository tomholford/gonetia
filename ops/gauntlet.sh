#!/usr/bin/env sh

set -e

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
REPO_ROOT="$SCRIPT_DIR/.."
cd "$REPO_ROOT"

echo "Running gonetia gauntlet: format → lint → build → test"
echo

echo "=== Step 1/4: Format ==="
./ops/fmt.sh

echo
echo "=== Step 2/4: Lint ==="
./ops/lint.sh

echo
echo "=== Step 3/4: Build ==="
./ops/build.sh

echo
echo "=== Step 4/4: Test ==="
./ops/test.sh

echo
echo "Gauntlet completed successfully!"
