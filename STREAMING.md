# Twitter Filtered Stream Guide

This guide shows you how to use `ctw` to screen tweets for keywords using Twitter's filtered stream API.

## Quick Start

### Option 1: Watch Mode (Easiest)

The `watch` command automatically sets up rules and monitors tweets in real-time:

```bash
# Watch for a single keyword
ctw watch --keyword "golang" --auto-setup

# Watch for multiple keywords
ctw watch --keyword "bitcoin" --keyword "crypto" --auto-setup

# Show detailed information including usernames
ctw watch --keyword "AI" --auto-setup --show-user --show-meta
```

Press `Ctrl+C` to stop watching.

When the stream stops, ctw prints a summary with runtime, tweet count, reconnects,
and the last disconnect reason so you can quickly see stream health.

### Option 2: Manual Rule Management

For more control, manually manage stream rules:

```bash
# Step 1: List existing rules
ctw stream rules list

# Step 2: Add rules for your keywords
ctw stream rules add --value "golang" --tag "golang-tweets"
ctw stream rules add --value "rust programming" --tag "rust-tweets"
ctw stream rules add --value "bitcoin OR ethereum" --tag "crypto"

# Step 3: Connect to the stream
ctw stream --field "tweet.fields=created_at,author_id,lang" \
            --field "expansions=author_id" \
            --field "user.fields=username,name"

# Step 4: Delete rules when done
ctw stream rules list  # Get rule IDs
ctw stream rules delete --id "1234567890" --id "9876543210"
```

## Twitter Rule Syntax

### Basic Keywords

```bash
# Single word
ctw stream rules add --value "bitcoin"

# Phrase (must include quotes in the value)
ctw stream rules add --value "\"climate change\""

# Multiple keywords (OR)
ctw stream rules add --value "cat OR dog OR pet"

# Multiple keywords (AND)
ctw stream rules add --value "cat has:images"
```

### Advanced Operators

```bash
# Tweets with media
ctw stream rules add --value "cats has:images"
ctw stream rules add --value "music has:videos"

# Tweets with links
ctw stream rules add --value "news has:links"

# Language filter
ctw stream rules add --value "bonjour lang:fr"

# From specific user
ctw stream rules add --value "from:TwitterDev"

# To specific user
ctw stream rules add --value "to:Twitter"

# Hashtags
ctw stream rules add --value "#bitcoin"

# Exclude terms
ctw stream rules add --value "cats -dogs"

# Retweets
ctw stream rules add --value "cats is:retweet"
ctw stream rules add --value "cats -is:retweet"  # Exclude retweets
```

### Complex Rules

```bash
# Combine multiple conditions
ctw stream rules add --value "bitcoin (price OR market) lang:en -is:retweet"

# Track mentions of multiple accounts
ctw stream rules add --value "(@elonmusk OR @BillGates) crypto"

# News about tech companies
ctw stream rules add --value "(Apple OR Google OR Microsoft) (announcement OR launch OR release) has:links lang:en"
```

## Real-World Examples

### 1. Monitor Your Brand

```bash
# Watch for mentions of your product
ctw watch --keyword "YourProduct" --keyword "@YourCompany" --auto-setup --show-user

# Track support requests
ctw stream rules add --value "(@YourCompany OR #YourProduct) (help OR support OR issue OR problem)"
```

### 2. Track Breaking News

```bash
# Breaking news with images
ctw watch --keyword "breaking news has:images lang:en" --auto-setup

# Major news sources
ctw stream rules add --value "from:CNN OR from:BBCBreaking OR from:Reuters"
```

### 3. Market Intelligence

```bash
# Cryptocurrency sentiment
ctw watch --keyword "bitcoin (bull OR bear OR crash OR moon)" --auto-setup

# Tech industry trends
ctw stream rules add --value "(AI OR \"machine learning\" OR GPT) (breakthrough OR innovation) lang:en"
```

### 4. Event Monitoring

```bash
# Conference hashtag
ctw watch --keyword "#DevConf2024" --auto-setup --show-user

# Sports event
ctw stream rules add --value "#SuperBowl has:images"
```

### 5. Research & Analysis

```bash
# Sentiment analysis data collection
ctw watch --keyword "climate change" --keyword "global warming" --auto-setup > climate_tweets.jsonl

# Location-based trending
ctw stream rules add --value "place_country:US (trending OR viral)"
```

## Output Formats

### JSON Format (Default)

The `stream` command outputs raw JSON:

```bash
ctw stream > tweets.json
```

### Human-Readable (Watch Command)

The `watch` command formats output for readability:

```
─────────────────────────────────────────
🐦 Tweet #1
ID: 1234567890
Time: 2025-10-19T12:34:56Z
Author: @username (Display Name)
Language: en

Text:
Check out this awesome #golang tutorial! 🚀
https://example.com/tutorial
```

### JSON Lines (Watch Command)

Use JSON output for automation pipelines:

```bash
ctw watch --keyword "golang" --auto-setup --json > watch.jsonl
```

For a compact, pretty JSON event:

```bash
ctw watch --keyword "golang" --auto-setup --json --pretty
```

## Best Practices

### 1. Rule Limits

- Free tier: Up to 25 rules
- Each rule can be up to 512 characters
- Combine related keywords in one rule when possible

### 2. Avoid Rate Limits

```bash
# Test rules with dry-run first
ctw stream rules add --value "test" --dry-run

# Start with narrow rules, then expand
ctw stream rules add --value "golang tutorial" --tag "specific"
# Later: broaden if needed
ctw stream rules add --value "golang" --tag "broad"
```

### 3. Save Stream Data

```bash
# Save to file with timestamps
ctw watch --keyword "bitcoin" --auto-setup | tee -a bitcoin_stream_$(date +%Y%m%d).log

# JSON Lines format for processing
ctw stream --field "tweet.fields=created_at,text,author_id" > stream_data.jsonl
```

### 4. Monitor Stream Health

```bash
# List current rules periodically
watch -n 60 'ctw stream rules list'

# Track tweet volume
ctw watch --keyword "trending" --auto-setup | grep "Tweet #" | wc -l

# Check reconnection behavior with summary output
ctw watch --keyword "the" --auto-setup
```

## Troubleshooting

### "No tweets appearing"

1. Check your rules are active:
   ```bash
   ctw stream rules list
   ```

2. Test with a common keyword:
   ```bash
   ctw watch --keyword "the" --auto-setup
   ```

3. Verify your bearer token:
   ```bash
   echo $BEARER_TOKEN
   ```

### "Too many rules" error

Delete unused rules:
```bash
# List all rules
ctw stream rules list

# Delete by ID
ctw stream rules delete --id "rule_id_1" --id "rule_id_2"
```

### "Connection reset"

The stream disconnects after 40 seconds of no data. This is normal - the client will reconnect automatically. Use broader rules if you experience frequent disconnects.

## Advanced: Processing Stream Data

### Parse with jq

```bash
# Extract just the tweet text
ctw stream | jq -r '.data[].text'

# Count tweets per minute
ctw stream | jq -r '.data[].created_at' | uniq -c

# Filter by language
ctw stream | jq 'select(.data[].lang == "en")'
```

### Pipe to Analysis Tools

```bash
# Sentiment analysis
ctw watch --keyword "AI" --auto-setup | grep "Text:" | your_sentiment_tool

# Store in database
ctw stream | while read line; do
  echo "$line" | psql -c "INSERT INTO tweets (data) VALUES ('$line')"
done
```

## Rule Management Workflow

```bash
# 1. Check current state
ctw stream rules list

# 2. Clear all rules (careful!)
ctw stream rules list | jq -r '.data[].id' | xargs -I {} ctw stream rules delete --id {}

# 3. Add new rule set
ctw stream rules add --value "keyword1" --tag "tag1"
ctw stream rules add --value "keyword2" --tag "tag2"

# 4. Verify
ctw stream rules list

# 5. Start watching
ctw watch --keyword "keyword1" --keyword "keyword2"
```

## API Reference

### Commands

- `ctw watch` - Watch tweets in real-time (recommended for beginners)
- `ctw stream` - Connect to filtered stream (raw JSON output)
- `ctw stream rules list` - Show active rules
- `ctw stream rules add` - Add a new rule
- `ctw stream rules delete` - Remove rules by ID

### Flags

**Watch command:**
- `--keyword` - Keyword to watch for (repeatable)
- `--auto-setup` - Automatically manage stream rules
- `--show-user` - Display author information
- `--show-meta` - Show additional metadata

**Stream command:**
- `--field` - Query parameter (e.g., `tweet.fields=created_at`)

**Rules commands:**
- `--value` - Rule value (required for add)
- `--tag` - Optional rule tag
- `--id` - Rule ID (required for delete)
- `--dry-run` - Test without applying changes

## Resources

- [Twitter Filtered Stream Docs](https://developer.twitter.com/en/docs/twitter-api/tweets/filtered-stream/introduction)
- [Twitter Rule Operators](https://developer.twitter.com/en/docs/twitter-api/tweets/filtered-stream/integrate/build-a-rule)
- [Rate Limits](https://developer.twitter.com/en/docs/twitter-api/rate-limits)
- Example scripts in `script/sh/examples/`
