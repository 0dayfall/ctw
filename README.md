# ctw - Command-Line Twitter Client

**A scriptable CLI for automating Twitter workflows using the Twitter v2 API.**

Perfect for building bots, monitoring social media, aggregating content, or integrating Twitter into your automation pipelines. Single binary, JSON output, works with standard Unix tools.

## Why ctw?

- **Built for Automation** - JSON output, exit codes, environment variables - integrates seamlessly with bash, cron, systemd
- **Complete API Coverage** - Tweets, DMs, users, media upload, streaming, likes, retweets, bookmarks, timelines
- **Single Binary** - No runtime dependencies, just download and run
- **Script-Friendly** - Pipeable output, works with `jq`, `grep`, `awk` and other Unix tools
- **Real-Time Processing** - Monitor tweets as they happen using Twitter's filtered stream API
- **Production-Ready** - Rate-limit handling, error reporting, comprehensive logging

## Quick Start

```bash
# Install (macOS)
brew tap 0dayfall/tap && brew install ctw

# Configure
export BEARER_TOKEN="your_twitter_bearer_token"

# Start automating
ctw init
ctw search recent --query "golang" | jq -r '.data[].text'
ctw users lookup --usernames "github,golang"
ctw tweets create --text "Hello from automation!"
```

## Common Use Cases

### Monitor Brand Mentions
```bash
# Real-time monitoring with keyword filtering
ctw watch --keyword "@YourBrand" --auto-setup | grep "urgent"

# JSON lines output for automation
ctw watch --keyword "@YourBrand" --auto-setup --json > brand_mentions.jsonl
```

More copy-paste examples live in `script/recipes/`.

### Collect Social Data
```bash
# Gather tweets for analysis
ctw search recent --query "climate change" --param "max_results=100" \
    | jq -r '.data[] | [.created_at, .text] | @csv' > data.csv
```

### Build a Twitter Bot
```bash
# Automated content curation
ctw search recent --query "golang tutorial" --param "max_results=10" \
    | jq -r '.data[0].id' \
    | xargs -I {} ctw retweets add --user-id YOUR_ID --tweet-id {}
```

### Schedule Posts
```bash
# Publish with media (cron job)
MEDIA_ID=$(ctw media upload --file chart.png | jq -r '.media_id_string')
ctw tweets create --text "Daily update" --media-ids "$MEDIA_ID"
```

### Get Alerts
```bash
# Stream to Slack/Discord webhook
ctw watch --keyword "BREAKING" --auto-setup \
    | while read -r line; do curl -X POST $WEBHOOK -d "$line"; done
```

## Complete Feature Set

**Tweets & Content**
- Create, delete, and lookup tweets
- Upload media (images, videos, GIFs) with chunked upload
- Search recent tweets with filtering
- Get tweet counts and analytics

**Streaming**
- Real-time filtered stream with keyword monitoring
- Rule management (add, list, delete)
- Watch command for easy keyword tracking

**User Operations**
- Lookup users by username or ID
- Follow/unfollow users
- Block/unblock users

**Engagement**
- Like/unlike tweets
- Retweet/unretweet
- Add/remove bookmarks
- List liked tweets, bookmarks, retweeters

**Timelines**
- User tweets timeline
- Mentions timeline
- Home timeline (reverse chronological)

**Direct Messages**
- Send DMs
- List conversations
- Delete messages

## Installation

**macOS (Homebrew)**

```bash
brew tap 0dayfall/tap && brew install ctw
```

**Ubuntu/Debian**

```bash
# Download from releases or build locally
make deb
sudo dpkg -i build/ctw_<version>_amd64.deb
```

**From Source**

```bash
git clone https://github.com/0dayfall/ctw.git
cd ctw
go build -o ctw ./cmd/ctw
sudo mv ctw /usr/local/bin/
```

**Pre-built Binaries**

Download from [releases](https://github.com/0dayfall/ctw/releases) for Linux, macOS, Windows (amd64/arm64).

See [INSTALL.md](INSTALL.md) for detailed installation instructions and [GORELEASER_SETUP.md](GORELEASER_SETUP.md) for building releases.

## Configuration

### Quick Setup (Recommended)

```bash
ctw init
```

This verifies your credentials and writes a config file with an env reference.

### Config File (TOML)

Default locations:
- macOS/Linux: `~/.config/ctw/config.toml`
- Windows: `%APPDATA%\\ctw\\config.toml`

Example:

```toml
[auth]
bearer_token = "env:BEARER_TOKEN"

[http]
user_agent = "ctw/0.2"
timeout = "15s"
retry = 3

[output]
pretty = false

[stream]
backoff_max = "2m"
```

See `config.example.toml` for a ready-to-copy template.

### Precedence

Flags > Env vars > Config file > Defaults

### Environment Variables

```bash
# Required: Twitter API bearer token
export BEARER_TOKEN="your_twitter_bearer_token_here"

# Optional: custom user agent
export USER_AGENT="my-bot/1.0"

# Optional: config overrides
export CTW_TIMEOUT="15s"
export CTW_RETRY="3"
export CTW_PRETTY="false"
export CTW_STREAM_BACKOFF_MAX="2m"
```

Get your bearer token from the Twitter Developer Portal.

## Automation Examples

### Real-Time Monitoring

```bash
# Monitor keywords and pipe to processing
ctw watch --keyword "golang" --auto-setup | grep "tutorial"

# Multi-keyword brand monitoring
ctw watch --keyword "@YourBrand" --keyword "YourProduct" --auto-setup --show-user

# Stream with complex rules
ctw stream rules add --value "bitcoin OR ethereum lang:en -is:retweet"
ctw stream
```

### Data Collection & Analysis

```bash
# Search and extract with jq
ctw search recent --query "AI ethics" --param "max_results=100" \
    | jq -r '.data[] | [.author_id, .text] | @csv'

# Get tweet counts over time
ctw counts recent --query "climate change" --granularity day

# Monitor user activity
ctw timelines user --user-id 123 --param "max_results=50"
```

### Content Publishing

```bash
# Simple tweet
ctw tweets create --text "Hello from automation"

# Tweet with media
MEDIA=$(ctw media upload --file chart.png --category tweet_image | jq -r '.media_id_string')
ctw tweets create --text "Daily metrics" --media-ids "$MEDIA"

# Scheduled via cron: 0 9 * * * /usr/local/bin/post_daily_update.sh
```

### User Management

```bash
# Lookup users
ctw users lookup --usernames "github,golang,docker"
ctw users lookup --ids "123,456,789"

# Manage relationships
ctw users follow --source-id YOUR_ID --target-id 123
ctw users block --source-id YOUR_ID --target-id 456
```

### Engagement Automation

```bash
# Like tweets matching criteria
ctw search recent --query "open source" --param "max_results=10" \
    | jq -r '.data[].id' \
    | xargs -I {} ctw likes add --user-id YOUR_ID --tweet-id {}

# Retweet quality content
ctw retweets add --user-id YOUR_ID --tweet-id 1234567890

# Manage bookmarks
ctw bookmarks add --user-id YOUR_ID --tweet-id 9876543210
ctw bookmarks list --user-id YOUR_ID
```

### Scripting Patterns

```bash
# Error handling with exit codes
if ctw tweets create --text "test"; then
    echo "Posted successfully"
else
    echo "Failed with code $?"
fi

# Loop through results
ctw search recent --query "golang" | jq -r '.data[].text' | while read -r tweet; do
    echo "Processing: $tweet"
done

# Combine multiple operations
USER_ID=$(ctw users lookup --usernames "twitter" | jq -r '.data[0].id')
ctw timelines user --user-id "$USER_ID"
```

## Command Reference

For detailed help on any command:

```bash
ctw --help
ctw stream --help
ctw search --help
# ... etc
```

**Available Commands:**
- `watch` - Monitor tweets with keyword filtering (easiest)
- `stream` - Manage filtered stream rules and connect
- `search` - Search recent or all tweets
- `counts` - Get tweet count aggregations
- `tweets` - Create, delete, and lookup tweets
- `users` - Lookup users and manage relationships
- `timelines` - Get user, mentions, and home timelines
- `likes` - Like, unlike, and list liked tweets
- `retweets` - Retweet, unretweet, and list retweeters
- `bookmarks` - Add, remove, and list bookmarks
- `dms` - Send, list, and delete direct messages
- `media` - Upload images, videos, and GIFs

## Documentation

- **[AUTOMATION.md](AUTOMATION.md)** - Comprehensive automation guide with real-world examples
- **[STREAMING.md](STREAMING.md)** - Detailed guide to Twitter's filtered stream API
- **[KEYWORD_SCREENING.md](KEYWORD_SCREENING.md)** - Quick reference for keyword filtering
- **[INSTALL.md](INSTALL.md)** - Installation instructions for all platforms
- **[script/sh/examples/](script/sh/examples/)** - Ready-to-use automation scripts

## How It Works

`ctw` is a thin CLI wrapper around the Twitter v2 API. It handles:

- **Authentication** - Bearer token management
- **HTTP** - Configurable client with rate-limit tracking
- **JSON** - Type-safe request/response handling
- **Errors** - Proper exit codes and error messages
- **Streaming** - Real-time filtered stream processing

All commands output JSON to stdout and logs/errors to stderr, making it perfect for Unix pipelines and automation scripts.

## Project Structure

```text
cmd/ctw/             # CLI commands (Cobra)
internal/client/     # HTTP client with auth and rate-limits
internal/tweet/      # Tweet services (publish, search, stream, likes, etc.)
internal/users/      # User services (lookup, follow, block)
internal/media/      # Media upload (chunked upload for large files)
internal/dm/         # Direct message services
script/sh/           # Shell script examples and testing utilities
```

## Development

**Run tests:**

```bash
go test ./...
```

Tests use `httptest.Server` mocks - no live API calls required.

**Build locally:**

```bash
make build        # Builds to bin/ctw
make install      # Installs to /usr/local/bin
make deb          # Creates .deb package
```

**Contributing:**

- Follow the service pattern: accept `context.Context`, return typed structs + `client.RateLimitSnapshot`
- Keep CLI commands thin - parse flags, call service, print JSON
- Add tests using `httptest.Server` helpers
- Update documentation when adding features

**Issues & Contributions:**

- Report bugs or request features via GitHub Issues
- PRs welcome -- small, focused changes are easiest to review

## API Rate Limits

Twitter API has rate limits. `ctw` tracks and reports rate-limit headers:

```bash
ctw search recent --query "test" 2>&1 | grep "Rate Limit"
```

Rate limits vary by endpoint. See [Twitter's documentation](https://developer.twitter.com/en/docs/twitter-api/rate-limits) for details.

## License

See [LICENSE](LICENSE) file.
