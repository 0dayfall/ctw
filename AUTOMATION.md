# Automation Guide for ctw

This guide shows how to use `ctw` in scripts for Twitter automation workflows.

## Why ctw is Perfect for Scripting

- **CLI-first design** - Every feature accessible via command line
- **JSON output** - Easy to parse with `jq`
- **Exit codes** - Standard Unix conventions for error handling
- **Environment variables** - Configure once, use everywhere
- **Pipeable** - Works with Unix pipes and redirects
- **No dependencies** - Single binary, no runtime required

## Quick Start

### 1. Setup for Automation

```bash
# Add to your ~/.bashrc or ~/.zshrc
export BEARER_TOKEN="your_twitter_bearer_token"
export PATH="$PATH:/path/to/ctw/bin"

# Or use a config file
cat > ~/.ctw_config << EOF
export BEARER_TOKEN="your_token_here"
export USER_AGENT="my-automation-bot/1.0"
EOF

# Source it in your scripts
source ~/.ctw_config
```

### 2. Test Your Setup

```bash
ctw users lookup --usernames twitter
```

## Automation Patterns

### Pattern 1: Monitor and React

Watch for keywords and trigger actions:

```bash
#!/bin/bash
# Monitor brand mentions and log them

ctw watch --keyword "@YourBrand" --auto-setup | while read -r line; do
    if echo "$line" | grep -q "Text:"; then
        # Extract the tweet text
        tweet=$(echo "$line" | sed 's/Text://')
        
        # Log it
        echo "[$(date)] New mention: $tweet" >> mentions.log
        
        # Send alert (example)
        # curl -X POST https://your-webhook.com/alert -d "tweet=$tweet"
    fi
done
```

### Pattern 2: Scheduled Data Collection

Use cron for periodic data gathering:

```bash
#!/bin/bash
# Collect tweets about a topic every hour
# Add to crontab: 0 * * * * /path/to/collect_tweets.sh

KEYWORD="artificial intelligence"
OUTPUT_DIR="/data/tweets"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

ctw search recent \
    --query "$KEYWORD" \
    --param "max_results=100" \
    > "${OUTPUT_DIR}/ai_tweets_${TIMESTAMP}.json"

# Process results
jq -r '.data[].text' "${OUTPUT_DIR}/ai_tweets_${TIMESTAMP}.json" \
    >> "${OUTPUT_DIR}/all_tweets.txt"
```

### Pattern 3: Automated Publishing

Schedule tweets or publish based on conditions:

```bash
#!/bin/bash
# Publish daily summary

# Generate content
SUMMARY=$(./generate_daily_summary.sh)

# Upload image if exists
if [ -f "daily_chart.png" ]; then
    MEDIA_ID=$(ctw media upload --file daily_chart.png --category tweet_image \
        | jq -r '.media_id_string')
    
    # Tweet with media
    ctw tweets create \
        --text "$SUMMARY" \
        --media-ids "$MEDIA_ID"
else
    # Tweet without media
    ctw tweets create --text "$SUMMARY"
fi
```

### Pattern 4: Bulk Operations

Process multiple items efficiently:

```bash
#!/bin/bash
# Bulk follow users from a list

USER_ID="your_user_id"

while IFS= read -r target_username; do
    # Look up user
    TARGET_ID=$(ctw users lookup --usernames "$target_username" \
        | jq -r '.data[0].id')
    
    if [ -n "$TARGET_ID" ]; then
        echo "Following $target_username (ID: $TARGET_ID)"
        ctw users follow --source-id "$USER_ID" --target-id "$TARGET_ID"
        
        # Rate limit friendly delay
        sleep 2
    fi
done < users_to_follow.txt
```

### Pattern 5: Data Pipeline

Build data processing pipelines:

```bash
#!/bin/bash
# Extract, transform, and load tweet data

# Extract: Get tweets
ctw search recent --query "golang" --param "max_results=100" |
    # Transform: Extract relevant fields
    jq -r '.data[] | [.id, .author_id, .text, .created_at] | @csv' |
    # Load: Import to database
    psql -d mydb -c "COPY tweets(id, author_id, text, created_at) FROM STDIN CSV"
```

## Real-World Automation Examples

### Example 1: Social Media Dashboard

```bash
#!/bin/bash
# Update dashboard with latest metrics

DASHBOARD_DIR="/var/www/dashboard"

# Get account stats
ctw users lookup --usernames "YourAccount" > "${DASHBOARD_DIR}/account.json"

# Get recent mentions
ctw search recent --query "@YourAccount" --param "max_results=50" \
    > "${DASHBOARD_DIR}/mentions.json"

# Get engagement metrics
ctw timelines user --user-id "YOUR_ID" --param "max_results=20" \
    > "${DASHBOARD_DIR}/recent_tweets.json"

# Generate HTML report
./generate_report.py "${DASHBOARD_DIR}"/*.json > "${DASHBOARD_DIR}/index.html"
```

### Example 2: Customer Support Bot

```bash
#!/bin/bash
# Monitor support mentions and auto-respond

SUPPORT_USER_ID="your_support_account_id"

ctw watch --keyword "@YourSupport" --auto-setup | while read -r line; do
    # Extract tweet ID from the stream
    if echo "$line" | grep -q "^ID:"; then
        TWEET_ID=$(echo "$line" | awk '{print $2}')
        
        # Get full tweet details
        TWEET_DATA=$(ctw tweets get --id "$TWEET_ID")
        
        # Check if it's a question
        if echo "$TWEET_DATA" | jq -r '.data.text' | grep -qi "how\|help\|support"; then
            # Auto-reply
            AUTHOR_ID=$(echo "$TWEET_DATA" | jq -r '.data.author_id')
            ctw tweets create \
                --text "Thanks for reaching out! Our team will respond soon. In the meantime, check out our FAQ: https://example.com/faq"
            
            # Log for follow-up
            echo "$TWEET_ID" >> support_queue.txt
        fi
    fi
done
```

### Example 3: Content Aggregator

```bash
#!/bin/bash
# Aggregate tweets from multiple sources

SOURCES=("techcrunch" "verge" "wired" "arstechnica")
OUTPUT_FILE="tech_news_$(date +%Y%m%d).json"

echo '{"aggregated_tweets": []}' > "$OUTPUT_FILE"

for source in "${SOURCES[@]}"; do
    echo "Fetching from @$source..."
    
    # Get user ID
    USER_ID=$(ctw users lookup --usernames "$source" | jq -r '.data[0].id')
    
    # Get their recent tweets
    ctw timelines user --user-id "$USER_ID" --param "max_results=10" |
        jq '.data[]' >> temp_tweets.json
done

# Combine and deduplicate
jq -s 'add | unique_by(.id)' temp_tweets.json > "$OUTPUT_FILE"
rm temp_tweets.json

echo "Aggregated $(jq '.[]  | length' "$OUTPUT_FILE") unique tweets"
```

### Example 4: Sentiment Analysis Pipeline

```bash
#!/bin/bash
# Real-time sentiment analysis

KEYWORD="YourProduct"
SENTIMENT_API="http://localhost:5000/analyze"

ctw watch --keyword "$KEYWORD" --auto-setup | while read -r line; do
    if echo "$line" | grep -q "^Text:"; then
        TWEET_TEXT=$(echo "$line" | sed 's/Text://' | xargs)
        
        # Analyze sentiment
        SENTIMENT=$(curl -s -X POST "$SENTIMENT_API" \
            -H "Content-Type: application/json" \
            -d "{\"text\": \"$TWEET_TEXT\"}" |
            jq -r '.sentiment')
        
        # Log with sentiment
        echo "[$(date)] [$SENTIMENT] $TWEET_TEXT" >> sentiment_log.txt
        
        # Alert on negative sentiment
        if [ "$SENTIMENT" = "negative" ]; then
            ./send_alert.sh "Negative mention detected: $TWEET_TEXT"
        fi
    fi
done
```

### Example 5: Automated Content Curation

```bash
#!/bin/bash
# Find and retweet quality content

USER_ID="your_user_id"
KEYWORDS=("golang tutorial" "golang tips" "golang best practices")

for keyword in "${KEYWORDS[@]}"; do
    echo "Searching for: $keyword"
    
    # Find tweets
    ctw search recent \
        --query "$keyword lang:en -is:retweet has:links" \
        --param "max_results=10" |
        jq -r '.data[] | select(.public_metrics.retweet_count > 100) | .id' |
        head -n 3 |
        while read -r tweet_id; do
            echo "Retweeting $tweet_id"
            ctw retweets add --user-id "$USER_ID" --tweet-id "$tweet_id"
            sleep 5
        done
done
```

## Advanced Techniques

### Error Handling

```bash
#!/bin/bash
# Robust error handling

set -euo pipefail  # Exit on error, undefined vars, pipe failures

function handle_error() {
    echo "❌ Error on line $1" >&2
    # Send alert
    curl -X POST https://alerts.example.com/error \
        -d "script=automation.sh&line=$1"
    exit 1
}

trap 'handle_error $LINENO' ERR

# Your automation code here
ctw tweets create --text "Automated tweet" || {
    echo "Failed to post tweet" >&2
    exit 1
}
```

### Rate Limit Handling

```bash
#!/bin/bash
# Handle rate limits gracefully

function api_call_with_retry() {
    local max_retries=3
    local retry_delay=60
    local attempt=1
    
    while [ $attempt -le $max_retries ]; do
        if "$@"; then
            return 0
        else
            echo "Attempt $attempt failed. Waiting ${retry_delay}s..." >&2
            sleep $retry_delay
            retry_delay=$((retry_delay * 2))  # Exponential backoff
            attempt=$((attempt + 1))
        fi
    done
    
    return 1
}

# Usage
api_call_with_retry ctw search recent --query "test"
```

### Parallel Processing

```bash
#!/bin/bash
# Process multiple queries in parallel

QUERIES=("golang" "python" "rust" "javascript")
MAX_PARALLEL=4

export BEARER_TOKEN  # Make available to subshells

process_query() {
    local query="$1"
    echo "Processing: $query"
    ctw search recent --query "$query" --param "max_results=100" \
        > "results_${query}.json"
}

export -f process_query

# Run in parallel
printf '%s\n' "${QUERIES[@]}" | xargs -P "$MAX_PARALLEL" -I {} bash -c 'process_query "$@"' _ {}
```

### Logging and Monitoring

```bash
#!/bin/bash
# Comprehensive logging

LOG_DIR="/var/log/ctw-automation"
mkdir -p "$LOG_DIR"

# Redirect all output
exec 1> >(tee -a "${LOG_DIR}/automation_$(date +%Y%m%d).log")
exec 2> >(tee -a "${LOG_DIR}/errors_$(date +%Y%m%d).log" >&2)

echo "[$(date)] Starting automation"

# Track metrics
START_TIME=$(date +%s)
TWEETS_PROCESSED=0

# Your automation code
# ...

# Summary
END_TIME=$(date +%s)
DURATION=$((END_TIME - START_TIME))
echo "[$(date)] Completed. Processed $TWEETS_PROCESSED tweets in ${DURATION}s"
```

## Integration Examples

### Slack Integration

```bash
#!/bin/bash
# Post tweet notifications to Slack

SLACK_WEBHOOK="https://hooks.slack.com/services/YOUR/WEBHOOK"

ctw watch --keyword "#YourHashtag" --auto-setup | while read -r line; do
    if echo "$line" | grep -q "^Text:"; then
        TWEET=$(echo "$line" | sed 's/Text://')
        
        curl -X POST "$SLACK_WEBHOOK" \
            -H 'Content-Type: application/json' \
            -d "{\"text\": \"New tweet: $TWEET\"}"
    fi
done
```

### Database Integration

```bash
#!/bin/bash
# Store tweets in PostgreSQL

DB_CONN="postgresql://user:pass@localhost/twitter_db"

ctw search recent --query "data science" --param "max_results=100" |
    jq -r '.data[] | [.id, .author_id, .text, .created_at] | @csv' |
    while IFS=, read -r id author_id text created_at; do
        psql "$DB_CONN" -c \
            "INSERT INTO tweets (id, author_id, text, created_at) 
             VALUES ('$id', '$author_id', '$text', '$created_at')
             ON CONFLICT (id) DO NOTHING"
    done
```

### Webhook Server

```bash
#!/bin/bash
# Simple webhook receiver

while true; do
    # Listen for webhooks
    REQUEST=$(nc -l 8080)
    
    # Extract action from webhook
    ACTION=$(echo "$REQUEST" | grep -oP '(?<=action=)[^&]+')
    
    case "$ACTION" in
        "post_tweet")
            TEXT=$(echo "$REQUEST" | grep -oP '(?<=text=)[^&]+' | sed 's/+/ /g')
            ctw tweets create --text "$TEXT"
            ;;
        "check_mentions")
            ctw search recent --query "@YourAccount" --param "max_results=10"
            ;;
    esac
done
```

## Cron Examples

```bash
# Monitor mentions every 5 minutes
*/5 * * * * /opt/ctw/scripts/check_mentions.sh

# Daily content curation at 9 AM
0 9 * * * /opt/ctw/scripts/curate_content.sh

# Hourly data backup
0 * * * * /opt/ctw/scripts/backup_tweets.sh

# Weekly analytics report every Monday at 8 AM
0 8 * * 1 /opt/ctw/scripts/generate_weekly_report.sh

# Real-time monitoring (runs continuously)
@reboot /opt/ctw/scripts/start_monitoring.sh
```

## Systemd Service Example

```ini
# /etc/systemd/system/ctw-monitor.service

[Unit]
Description=CTW Twitter Monitor
After=network.target

[Service]
Type=simple
User=twitter-bot
Environment="BEARER_TOKEN=your_token_here"
ExecStart=/usr/local/bin/ctw watch --keyword "YourBrand" --auto-setup
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl enable ctw-monitor
sudo systemctl start ctw-monitor
sudo journalctl -u ctw-monitor -f  # Watch logs
```

## Docker Integration

```dockerfile
FROM golang:1.18-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o ctw ./cmd/ctw

FROM alpine:latest
RUN apk --no-cache add ca-certificates jq
COPY --from=builder /app/ctw /usr/local/bin/
ENV BEARER_TOKEN=""
CMD ["ctw", "watch", "--keyword", "automation", "--auto-setup"]
```

## Best Practices

### 1. Use Environment Variables

```bash
# Don't hardcode tokens
export BEARER_TOKEN="..."

# Use config files for complex setups
source /etc/ctw/config
```

### 2. Implement Retry Logic

```bash
for i in {1..3}; do
    ctw tweets create --text "Retry test" && break
    sleep $((i * 30))
done
```

### 3. Log Everything

```bash
ctw search recent --query "test" 2>&1 | tee -a ctw.log
```

### 4. Use Locks for Cron Jobs

```bash
#!/bin/bash
LOCKFILE=/tmp/ctw-script.lock

if [ -e "$LOCKFILE" ]; then
    echo "Already running"
    exit 1
fi

trap "rm -f $LOCKFILE" EXIT
touch "$LOCKFILE"

# Your automation code here
```

### 5. Monitor Script Health

```bash
#!/bin/bash
# Heartbeat script

while true; do
    if ! pgrep -f "ctw watch" > /dev/null; then
        echo "Monitor died, restarting..."
        /opt/ctw/scripts/start_monitor.sh &
    fi
    sleep 60
done
```

## Troubleshooting

### Check Exit Codes

```bash
ctw tweets create --text "test"
if [ $? -eq 0 ]; then
    echo "Success"
else
    echo "Failed with code $?"
fi
```

### Debug Mode

```bash
# Enable verbose curl output
export CTW_DEBUG=1

# Or use set -x for bash debugging
set -x
ctw search recent --query "debug"
set +x
```

### Test Before Deploying

```bash
# Dry-run mode (if supported)
ctw stream rules add --value "test" --dry-run

# Test with small limits
ctw search recent --query "test" --param "max_results=1"
```

## Resources

- **Examples**: `script/sh/examples/` directory
- **CLI Reference**: `ctw --help` and `ctw COMMAND --help`
- **API Limits**: https://developer.twitter.com/en/docs/twitter-api/rate-limits
- **Error Codes**: Check exit codes in your scripts (0 = success, >0 = error)

## Summary

`ctw` is designed for automation:
- ✅ Single binary deployment
- ✅ JSON output for easy parsing
- ✅ Standard exit codes
- ✅ Environment variable configuration
- ✅ Stream processing capabilities
- ✅ Rate-limit aware
- ✅ Scriptable and pipeable

Start automating your Twitter workflows today! 🚀
