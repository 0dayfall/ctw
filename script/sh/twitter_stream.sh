#!/bin/bash
#
# Add Twitter Filtered Stream Rules (Direct API Call)
#
# This script adds sample rules to the Twitter filtered stream.
# Edit the rules in the JSON payload below to match your needs.
#
# Prerequisites:
#   - BEARER_TOKEN environment variable must be set
#
# Usage:
#   ./script/sh/twitter_stream.sh
#
# For production use, consider: ctw stream rules add --value "keyword" --tag "tag"

set -e

if [ -z "$BEARER_TOKEN" ]; then
    echo "❌ Error: BEARER_TOKEN environment variable is not set" >&2
    exit 1
fi

echo "📝 Adding filtered stream rules..." >&2

# Edit these rules to match your use case
RULES_JSON='{
  "add": [
    {
      "value": "golang OR rust",
      "tag": "programming-languages"
    },
    {
      "value": "bitcoin OR ethereum",
      "tag": "crypto"
    },
    {
      "value": "AI OR \"machine learning\"",
      "tag": "artificial-intelligence"
    }
  ]
}'

curl -X POST \
    'https://api.twitter.com/2/tweets/search/stream/rules' \
    -H 'Content-Type: application/json' \
    -H "Authorization: Bearer $BEARER_TOKEN" \
    -H "User-Agent: ctw-curl/1.0" \
    -d "$RULES_JSON"

echo "" >&2
echo "✅ Rules added successfully" >&2
echo "" >&2
echo "💡 Tip: Use 'ctw stream rules list' to verify" >&2
