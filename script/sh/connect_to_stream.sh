#!/bin/bash
#
# Connect to Twitter Filtered Stream (Direct API Call)
#
# This script demonstrates how to connect to the Twitter filtered stream
# API using curl. For production use, consider using the ctw CLI instead.
#
# Prerequisites:
#   - BEARER_TOKEN environment variable must be set
#
# Usage:
#   ./script/sh/connect_to_stream.sh

set -e

if [ -z "$BEARER_TOKEN" ]; then
    echo "❌ Error: BEARER_TOKEN environment variable is not set" >&2
    echo "" >&2
    echo "Get your bearer token from:" >&2
    echo "https://developer.twitter.com/en/portal/dashboard" >&2
    echo "" >&2
    echo "Then run: export BEARER_TOKEN='your_token_here'" >&2
    exit 1
fi

echo "🔴 Connecting to Twitter filtered stream..." >&2
echo "Press Ctrl+C to stop" >&2
echo "" >&2

curl -N -X GET \
    -H "Authorization: Bearer $BEARER_TOKEN" \
    -H "User-Agent: ctw-curl/1.0" \
    "https://api.twitter.com/2/tweets/search/stream?tweet.fields=created_at,author_id,lang&expansions=author_id&user.fields=username,name,created_at"
