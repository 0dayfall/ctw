#!/bin/bash
# Example: Monitor crypto tweets with multiple keywords
#
# This example shows how to track multiple related keywords
# and save the output to a file for later analysis

if [ -z "$BEARER_TOKEN" ]; then
    echo "❌ Error: BEARER_TOKEN environment variable is not set"
    exit 1
fi

OUTPUT_FILE="crypto_tweets_$(date +%Y%m%d_%H%M%S).log"

echo "🔴 Monitoring crypto tweets..."
echo "📁 Saving to: $OUTPUT_FILE"
echo "Press Ctrl+C to stop"
echo ""

# Watch multiple crypto keywords and save to file
./bin/ctw watch \
    --keyword "bitcoin" \
    --keyword "ethereum" \
    --keyword "crypto" \
    --auto-setup \
    --show-user \
    | tee "$OUTPUT_FILE"
