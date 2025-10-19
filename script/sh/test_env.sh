#!/bin/bash
#
# Test Environment Configuration
#
# This script validates that your environment is properly configured
# for using the ctw CLI and Twitter API.
#
# Usage:
#   ./script/sh/test_env.sh

echo "🔍 Testing environment configuration..."
echo ""

# Check BEARER_TOKEN
if [ -z "$BEARER_TOKEN" ]; then
    echo "❌ BEARER_TOKEN: Not set"
    echo "   Fix: export BEARER_TOKEN='your_token_here'"
    TOKEN_OK=false
else
    # Mask the token for security
    MASKED_TOKEN="${BEARER_TOKEN:0:10}...${BEARER_TOKEN: -10}"
    echo "✅ BEARER_TOKEN: Set ($MASKED_TOKEN)"
    TOKEN_OK=true
fi

# Check for ctw binary
if [ -f "./bin/ctw" ]; then
    echo "✅ ctw binary: Found (./bin/ctw)"
    CTW_OK=true
elif command -v ctw &> /dev/null; then
    CTW_LOCATION=$(which ctw)
    echo "✅ ctw binary: Found ($CTW_LOCATION)"
    CTW_OK=true
else
    echo "⚠️  ctw binary: Not found"
    echo "   Build with: go build -o bin/ctw ./cmd/ctw"
    CTW_OK=false
fi

# Check for jq (optional but useful)
if command -v jq &> /dev/null; then
    echo "✅ jq: Installed"
else
    echo "⚠️  jq: Not installed (optional, for JSON formatting)"
    echo "   Install with: brew install jq"
fi

# Check for curl
if command -v curl &> /dev/null; then
    CURL_VERSION=$(curl --version | head -n1)
    echo "✅ curl: $CURL_VERSION"
else
    echo "❌ curl: Not found"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ "$TOKEN_OK" = true ] && [ "$CTW_OK" = true ]; then
    echo "🎉 Environment is ready!"
    echo ""
    echo "Try these commands:"
    echo "  ./bin/ctw stream rules list"
    echo "  ./bin/ctw watch --keyword 'golang' --auto-setup"
elif [ "$TOKEN_OK" = true ]; then
    echo "⚠️  Environment partially ready"
    echo ""
    echo "Build ctw first:"
    echo "  go build -o bin/ctw ./cmd/ctw"
else
    echo "❌ Environment not ready"
    echo ""
    echo "Set your bearer token:"
    echo "  export BEARER_TOKEN='your_token_here'"
    echo ""
    echo "Get a token from:"
    echo "  https://developer.twitter.com/en/portal/dashboard"
fi
