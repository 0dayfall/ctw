# Twitter Keyword Screening - Quick Reference

## 🎯 What You Can Do

With `ctw`, you can **screen tweets in real-time** for specific keywords, hashtags, users, or complex search criteria. Perfect for:

- 📊 Brand monitoring and social listening
- 📰 Breaking news tracking
- 💰 Market intelligence and trend analysis
- 🔍 Research and data collection
- ⚡ Event monitoring
- 🎯 Customer support tracking

## 🚀 Fastest Way to Start

```bash
# 1. Set your Twitter API token
export BEARER_TOKEN="your_twitter_bearer_token_here"

# 2. Build the tool
go build -o bin/ctw ./cmd/ctw

# 3. Watch for keywords
./bin/ctw watch --keyword "golang" --auto-setup --show-user
```

That's it! You'll see tweets matching "golang" appear in real-time. Press Ctrl+C to stop.

## 📖 Three Methods to Screen Tweets

### Method 1: Watch (Recommended for Beginners)

**Best for:** Quick keyword monitoring with minimal setup

```bash
# Single keyword
ctw watch --keyword "bitcoin" --auto-setup

# Multiple keywords
ctw watch --keyword "AI" --keyword "GPT" --auto-setup --show-user

# Save to file
ctw watch --keyword "golang" --auto-setup > golang_tweets.log
```

**Pros:**
- ✅ Automatic rule management
- ✅ Human-readable output
- ✅ Real-time statistics
- ✅ Easy to use

### Method 2: Manual Rules (Advanced)

**Best for:** Complex filtering and fine-grained control

```bash
# Add rules
ctw stream rules add --value "golang" --tag "go"
ctw stream rules add --value "bitcoin OR ethereum" --tag "crypto"

# View rules
ctw stream rules list

# Connect to stream
ctw stream

# Delete rules when done
ctw stream rules delete --id "rule_id"
```

**Pros:**
- ✅ Persistent rules (survive restarts)
- ✅ Complex Boolean logic
- ✅ Twitter's full rule syntax
- ✅ Rate limit friendly

### Method 3: Programmatic (Integration)

**Best for:** Integrating with other tools

```bash
# Stream raw JSON
ctw stream | jq '.data[].text'

# Process with custom scripts
ctw stream | while read tweet; do
    echo "$tweet" | your_processing_script
done
```

## 🔥 Common Use Cases

### Brand Monitoring

```bash
# Monitor brand mentions
ctw watch --keyword "@YourBrand" --keyword "YourProduct" --auto-setup --show-user

# Track support requests
ctw stream rules add --value "(@YourBrand OR #YourProduct) (help OR support OR issue)"
```

### News Tracking

```bash
# Breaking news with media
ctw watch --keyword "breaking news has:images lang:en" --auto-setup

# Specific sources
ctw stream rules add --value "from:CNN OR from:BBCBreaking OR from:Reuters"
```

### Market Intelligence

```bash
# Crypto sentiment
ctw watch --keyword "bitcoin (bull OR bear OR crash)" --auto-setup

# Tech trends
ctw stream rules add --value "(AI OR GPT) (breakthrough OR innovation) lang:en -is:retweet"
```

### Event Coverage

```bash
# Conference hashtag
ctw watch --keyword "#DevConf2024" --auto-setup --show-user

# Live events with photos
ctw stream rules add --value "#SuperBowl has:images"
```

### Research & Data Collection

```bash
# Collect climate change discussions
ctw watch --keyword "climate change" --auto-setup > climate_data_$(date +%Y%m%d).jsonl

# Political sentiment
ctw stream rules add --value "election2024 lang:en (excited OR concerned OR worried)"
```

## 🎓 Rule Syntax Cheat Sheet

### Basic

| Syntax | Example | Matches |
|--------|---------|---------|
| keyword | `bitcoin` | Tweets containing "bitcoin" |
| "phrase" | `"climate change"` | Exact phrase |
| OR | `cat OR dog` | Either keyword |
| AND (space) | `cat images` | Both keywords |
| -exclude | `cat -dog` | "cat" but not "dog" |

### Advanced

| Operator | Example | Matches |
|----------|---------|---------|
| `from:` | `from:TwitterDev` | Tweets from specific user |
| `to:` | `to:Support` | Replies to user |
| `@mention` | `@elonmusk` | Mentions of user |
| `#hashtag` | `#bitcoin` | Tweets with hashtag |
| `has:images` | `cat has:images` | Tweets with images |
| `has:videos` | `news has:videos` | Tweets with videos |
| `has:links` | `article has:links` | Tweets with URLs |
| `lang:` | `bonjour lang:fr` | Specific language |
| `is:retweet` | `news is:retweet` | Only retweets |
| `-is:retweet` | `news -is:retweet` | Exclude retweets |

### Complex Examples

```bash
# Tech news with links, English only, no retweets
ctw stream rules add --value "(Apple OR Google) announcement has:links lang:en -is:retweet"

# Bitcoin discussions with sentiment
ctw stream rules add --value "bitcoin (bullish OR bearish OR moon OR crash)"

# Climate science with media
ctw stream rules add --value "(climate OR \"global warming\") (research OR study) (has:images OR has:videos) lang:en"
```

## 💡 Pro Tips

### 1. Test Rules First

```bash
# Dry-run before adding
ctw stream rules add --value "test" --dry-run
```

### 2. Use Tags for Organization

```bash
ctw stream rules add --value "golang" --tag "programming-go"
ctw stream rules add --value "python" --tag "programming-python"
```

### 3. Monitor Stream Health

```bash
# Check active rules
ctw stream rules list

# Watch tweet rate
ctw watch --keyword "test" --auto-setup | grep "Tweet #"
```

### 4. Save Output for Analysis

```bash
# Timestamped log files
ctw watch --keyword "AI" --auto-setup | tee ai_tweets_$(date +%Y%m%d).log

# JSON for processing
ctw stream > tweets.jsonl
```

### 5. Combine with Other Tools

```bash
# Real-time sentiment analysis
ctw stream | jq -r '.data[].text' | sentiment-analyzer

# Count tweets per minute
ctw watch --keyword "trending" --auto-setup | grep "Text:" | pv -l -i 60 > /dev/null
```

## 📊 Understanding Output

### Watch Command Output

```
─────────────────────────────────────────
🐦 Tweet #42
ID: 1234567890123456789
Time: 2025-10-19T15:30:00Z
Author: @username (Display Name)
Language: en

Text:
Just discovered this amazing #golang library! 🚀
Check it out: https://example.com
```

### Stream Command Output (JSON)

```json
{
  "data": [{
    "id": "1234567890",
    "text": "Example tweet text",
    "author_id": "987654321",
    "created_at": "2025-10-19T15:30:00.000Z",
    "lang": "en"
  }],
  "includes": {
    "users": [{
      "id": "987654321",
      "username": "example",
      "name": "Example User"
    }]
  }
}
```

## 🚨 Troubleshooting

### No tweets appearing?

1. Check rules: `ctw stream rules list`
2. Test with common keyword: `ctw watch --keyword "the" --auto-setup`
3. Verify token: `echo $BEARER_TOKEN`

### Too many tweets?

Make rules more specific:
```bash
# Instead of:
ctw watch --keyword "news" --auto-setup

# Try:
ctw stream rules add --value "news (breaking OR urgent) has:images lang:en -is:retweet"
```

### Stream disconnects?

Normal after 40s of no data. Broaden your rules or expect reconnections.

### Rate limited?

- Free tier: 50 requests per 15 minutes for rules
- Use `--dry-run` to test without consuming rate limit
- Avoid frequent rule changes

## 📚 More Resources

- **Full Guide**: See [STREAMING.md](STREAMING.md) for comprehensive documentation
- **Examples**: Check `script/sh/examples/` directory for ready-to-use scripts
- **Twitter Docs**: [developer.twitter.com/en/docs/twitter-api/tweets/filtered-stream](https://developer.twitter.com/en/docs/twitter-api/tweets/filtered-stream/introduction)

## 🎬 Example Scripts

We've included ready-to-use scripts in the `script/sh/examples/` directory:

```bash
# Watch Go programming tweets
./script/sh/examples/watch_golang.sh

# Monitor crypto discussions
./script/sh/examples/monitor_crypto.sh

# Advanced rule management
./script/sh/examples/advanced_rules.sh
```

## ⚡ Quick Commands Reference

```bash
# Basic keyword watch
ctw watch --keyword "KEYWORD" --auto-setup

# Multiple keywords with user info
ctw watch --keyword "K1" --keyword "K2" --auto-setup --show-user

# Add custom rule
ctw stream rules add --value "RULE" --tag "TAG"

# List all rules
ctw stream rules list

# Delete rules
ctw stream rules delete --id "ID1" --id "ID2"

# Raw stream connection
ctw stream

# Help for any command
ctw watch --help
ctw stream rules --help
```

---

**Ready to start?** Run: `ctw watch --keyword "your_keyword" --auto-setup`
