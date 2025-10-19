#!/bin/bash
#
# Delete Twitter Filtered Stream Rule (Direct API Call)
#
# This script deletes a filtered stream rule by its ID.
#
# Prerequisites:
#   - BEARER_TOKEN environment variable must be set
#
# Usage:
#   ./script/sh/delete_rule.sh RULE_ID
#
# Example:
#   ./script/sh/delete_rule.sh 1234567890
#
# Get rule IDs with:
#   ctw stream rules list
#   OR
#   ./script/sh/validate_twitter_stream.sh

set -e

if [ -z "$BEARER_TOKEN" ]; then
    echo "❌ Error: BEARER_TOKEN environment variable is not set" >&2
    exit 1
fi

if [ -z "$1" ]; then
    echo "❌ Error: Rule ID is required" >&2
    echo "" >&2
    echo "Usage: $0 RULE_ID" >&2
    echo "" >&2
    echo "Example:" >&2
    echo "  $0 1234567890" >&2
    echo "" >&2
    echo "Get rule IDs with:" >&2
    echo "  ctw stream rules list" >&2
    echo "  OR" >&2
    echo "  ./script/sh/validate_twitter_stream.sh" >&2
    exit 1
fi

RULE_ID="$1"

echo "🗑️  Deleting rule ID: $RULE_ID..." >&2

curl -X POST \
    'https://api.twitter.com/2/tweets/search/stream/rules' \
    -H 'Content-Type: application/json' \
    -H "Authorization: Bearer $BEARER_TOKEN" \
    -H "User-Agent: ctw-curl/1.0" \
    -d "{\"delete\":{\"ids\":[\"$RULE_ID\"]}}"

echo "" >&2
echo "✅ Rule deleted successfully" >&2
