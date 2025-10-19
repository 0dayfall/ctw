#!/bin/bash
#
# Block/Unblock Twitter Users (Direct API Call)
#
# This script demonstrates blocking/unblocking users via the Twitter API.
# Note: This uses OAuth 1.0a which requires additional setup.
#
# Prerequisites:
#   - OAUTH_SIGNATURE environment variable (OAuth 1.0a signature)
#   - User ID to block/unblock
#
# Usage:
#   ./script/sh/block.sh USER_ID [action]
#
# Example:
#   ./script/sh/block.sh 2244994945 block
#   ./script/sh/block.sh 2244994945 unblock
#
# For production use, consider: ctw users block/unblock

set -e

if [ -z "$OAUTH_SIGNATURE" ]; then
    echo "❌ Error: OAUTH_SIGNATURE environment variable is not set" >&2
    echo "" >&2
    echo "This script requires OAuth 1.0a authentication." >&2
    echo "For easier usage, use the ctw CLI instead:" >&2
    echo "" >&2
    echo "  ctw users block --source-id YOUR_ID --target-id $1" >&2
    echo "  ctw users unblock --source-id YOUR_ID --target-id $1" >&2
    exit 1
fi

if [ -z "$1" ]; then
    echo "❌ Error: User ID is required" >&2
    echo "" >&2
    echo "Usage: $0 USER_ID [action]" >&2
    echo "" >&2
    echo "Example:" >&2
    echo "  $0 2244994945 block" >&2
    echo "  $0 2244994945 unblock" >&2
    exit 1
fi

USER_ID="$1"
ACTION="${2:-block}"

if [ "$ACTION" = "block" ]; then
    echo "🚫 Blocking user $USER_ID..." >&2
    METHOD="POST"
elif [ "$ACTION" = "unblock" ]; then
    echo "✅ Unblocking user $USER_ID..." >&2
    METHOD="DELETE"
else
    echo "❌ Error: Invalid action '$ACTION'" >&2
    echo "Valid actions: block, unblock" >&2
    exit 1
fi

curl -X $METHOD \
    "https://api.twitter.com/2/users/$USER_ID/blocking" \
    -H "Authorization: OAuth $OAUTH_SIGNATURE" \
    -H "User-Agent: ctw-curl/1.0"

echo "" >&2
echo "✅ Action completed" >&2
