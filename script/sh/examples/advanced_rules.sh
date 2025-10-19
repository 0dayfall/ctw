#!/bin/bash
# Example: Advanced stream rule management
#
# This example shows how to manually manage stream rules
# for more complex filtering scenarios

if [ -z "$BEARER_TOKEN" ]; then
    echo "❌ Error: BEARER_TOKEN environment variable is not set"
    exit 1
fi

echo "📋 Step 1: List current rules"
./bin/ctw stream rules list
echo ""

echo "➕ Step 2: Add complex rules"
echo ""

# Rule 1: Track tech news with links in English
./bin/ctw stream rules add \
    --value "(Apple OR Google OR Microsoft) (announcement OR launch) has:links lang:en -is:retweet" \
    --tag "tech-news"

# Rule 2: Monitor AI discussions with images
./bin/ctw stream rules add \
    --value "(AI OR \"artificial intelligence\" OR \"machine learning\") has:images lang:en" \
    --tag "ai-visuals"

# Rule 3: Track a specific user's tweets
./bin/ctw stream rules add \
    --value "from:TwitterDev" \
    --tag "twitter-dev"

echo ""
echo "✅ Step 3: Verify rules were added"
./bin/ctw stream rules list
echo ""

echo "🔴 Step 4: Start streaming"
echo "Press Ctrl+C to stop"
echo ""

./bin/ctw stream \
    --field "tweet.fields=created_at,author_id,lang,possibly_sensitive" \
    --field "expansions=author_id" \
    --field "user.fields=username,name"
