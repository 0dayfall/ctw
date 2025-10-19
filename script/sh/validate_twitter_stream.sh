#!/bin/bash
#
# List Twitter Filtered Stream Rules (Direct API Call)
#
# This script retrieves and displays all active filtered stream rules.
#
# Prerequisites:
#   - BEARER_TOKEN environment variable must be set
#
# Usage:
#   ./script/sh/validate_twitter_stream.sh
#
# For production use, consider: ctw stream rules list

set -e

if [ -z "$BEARER_TOKEN" ]; then
    echo "❌ Error: BEARER_TOKEN environment variable is not set" >&2
    exit 1
fi

echo "📋 Fetching active stream rules..." >&2
echo "" >&2

curl -s -X GET \
    'https://api.twitter.com/2/tweets/search/stream/rules' \
    -H "Authorization: Bearer $BEARER_TOKEN" \
    -H "User-Agent: ctw-curl/1.0" | jq '.'

echo "" >&2
echo "💡 Tip: Use 'ctw stream rules list' for better formatting" >&2
