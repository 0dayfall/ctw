#!/bin/bash
# Example: Watch tweets about golang in real-time
# 
# Prerequisites:
# 1. Set BEARER_TOKEN environment variable
# 2. Build ctw: go build -o bin/ctw ./cmd/ctw
#
# Usage: ./script/sh/examples/watch_golang.sh

if [ -z "$BEARER_TOKEN" ]; then
    echo "❌ Error: BEARER_TOKEN environment variable is not set"
    echo ""
    echo "Get your bearer token from https://developer.twitter.com/en/portal/dashboard"
    echo "Then run: export BEARER_TOKEN='your_token_here'"
    exit 1
fi

echo "🔴 Starting real-time tweet monitor for 'golang' keyword..."
echo "Press Ctrl+C to stop"
echo ""

# Watch tweets mentioning golang, showing usernames and metadata
./bin/ctw watch \
    --keyword "golang" \
    --auto-setup \
    --show-user \
    --show-meta
