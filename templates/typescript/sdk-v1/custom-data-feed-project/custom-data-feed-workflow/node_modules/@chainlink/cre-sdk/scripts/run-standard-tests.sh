#!/bin/zsh

# Run standard tests script
# Usage: ./scripts/run-standard-tests.sh

set -e

# Create dist test workflow folder
mkdir -p ./dist/workflows/standard_tests

# Build javy wasm
if [ ! -f ../cre-sdk-javy-plugin/dist/javy_chainlink_sdk.wasm ]; then
    echo "Error: javy_chainlink_sdk.wasm not found"
    exit 1
fi

# Copy standard test files to dist
echo "Copying test files..." && cp -r ./src/standard_tests/* ./dist/workflows/standard_tests/

# Run standard tests
make standard_tests
