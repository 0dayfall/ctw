# ctw

A Go 1.18 command-line toolkit for working with Twitter v2 REST endpoints, with powerful **real-time tweet streaming** and keyword monitoring.

## ✨ Key Features

- **🔴 Real-time Tweet Streaming** - Watch tweets matching keywords as they happen
- **🔍 Filtered Stream** - Screen tweets for specific keywords, hashtags, users, or complex rules
- **📊 Tweet Analytics** - Track timelines, likes, retweets, bookmarks
- **📤 Publishing** - Create tweets, upload media, send DMs
- **👥 User Management** - Lookup users, manage follows/blocks
- **🔧 Flexible** - Configurable HTTP client with bearer-token auth and rate-limit tracking

## Quick Start: Tweet Streaming

Watch tweets for keywords in real-time:

```bash
# Set your Twitter API bearer token
export BEARER_TOKEN="your_twitter_bearer_token"

# Watch for keywords (auto-setup stream rules)
ctw watch --keyword "golang" --auto-setup --show-user

# Monitor multiple keywords
ctw watch --keyword "bitcoin" --keyword "crypto" --auto-setup

# Track your brand mentions
ctw watch --keyword "@YourCompany" --keyword "YourProduct" --auto-setup
```

**📚 Documentation:**
- [KEYWORD_SCREENING.md](KEYWORD_SCREENING.md) - Quick reference for keyword screening
- [STREAMING.md](STREAMING.md) - Comprehensive streaming guide with examples
- `script/sh/examples/` directory - Ready-to-use shell scripts

## Features

- Configurable HTTP client (`internal/client`) with bearer-token auth, user-agent overrides, and rate-limit parsing.
- Service wrappers for tweets (timelines, likes, retweets, bookmarks, lookup, publish), users (lookup, relationships), direct messages, filtered stream, and search/counts.
- Cobra-powered CLI in `cmd/ctw` with commands: `stream`, `search`, `counts`, `users`, `tweets`, `dms`, `likes`, `retweets`, `bookmarks`, and `timelines`.

## Installation

### macOS (Homebrew)

```bash
brew tap 0dayfall/tap
brew install ctw
```

### Ubuntu/Debian (.deb package)

Build and install the `.deb` package:

```bash
make deb
sudo dpkg -i build/ctw_<version>_amd64.deb
```

### Using GoReleaser (All Platforms)

GoReleaser builds for multiple platforms and handles releases automatically:

```bash
# Install GoReleaser
brew install goreleaser/tap/goreleaser

# Create a release (requires git tag)
git tag -a v0.1.0 -m "Release v0.1.0"
goreleaser release --clean
```

For detailed installation instructions including manual builds, see [INSTALL.md](INSTALL.md).

### From Source

1. **Install dependencies / tidy modules**

   ```bash
   go mod tidy
   ```

2. **Set credentials**

   ```bash
   export BEARER_TOKEN="<your twitter bearer token>"
   # Optional: override defaults
   export USER_AGENT="my-client/1.0"
   ```

3. **Build or run the CLI**

   ```bash
   go run ./cmd/ctw --help
   go build -o bin/ctw ./cmd/ctw
   
   # Or install to /usr/local/bin
   make build
   sudo make install
   ```

## CLI Examples

### Real-Time Streaming

```bash
# Watch tweets for keywords (easiest method)
ctw watch --keyword "AI" --keyword "machine learning" --auto-setup --show-user

# Manual stream rule management
ctw stream rules add --value "golang" --tag "go-lang"
ctw stream rules add --value "bitcoin OR ethereum" --tag "crypto"
ctw stream rules list
ctw stream  # Start streaming

# Advanced filtering
ctw stream rules add --value "cats has:images lang:en -is:retweet"
ctw watch --keyword "breaking news" --auto-setup

# Delete rules
ctw stream rules delete --id "rule_id_here"
```

### Tweets & Timelines

```bash
# Stream with selected fields (old method - use watch instead)
ctw stream --field tweet.fields=created_at --field expansions=author_id

# Search and analyze
ctw search recent --query "golang" --param max_results=20

# User timelines
ctw timelines user --user-id 123 --param max_results=10

# Add a filtered stream rule (dry run)
ctw stream rules add --value "cats has:images" --tag "cats" --dry-run

# Recent search with pagination token
ctw search recent --query "golang" --param max_results=20 --next-token <token>

# Tweet counts (recent)
ctw counts recent --query "from:TwitterDev" --granularity hour

# Lookup multiple usernames
ctw users lookup --usernames alice,bob --param "user.fields=created_at"

# Follow a user
ctw users follow --source-id 123 --target-id 456

# Publish and delete tweets
ctw tweets create --text "automation ready"
ctw tweets delete --id 1234567890

# Direct messages
ctw dms send --user-id 987654321 --text "hey there"
ctw dms list --param pagination_token=abc123
ctw dms delete --id event-123

# Likes
ctw likes add --user-id 123 --tweet-id 456
ctw likes remove --user-id 123 --tweet-id 456
ctw likes list --user-id 123

# Retweets
ctw retweets add --user-id 123 --tweet-id 789
ctw retweets remove --user-id 123 --tweet-id 789
ctw retweets list --tweet-id 789

# Bookmarks
ctw bookmarks add --user-id 123 --tweet-id 999
ctw bookmarks remove --user-id 123 --tweet-id 999
ctw bookmarks list --user-id 123

# Timelines
ctw timelines user --user-id 123 --param max_results=10
ctw timelines mentions --user-id 123
ctw timelines home --user-id 123

# Tweet lookup
ctw tweets get --id 1234567890
ctw tweets get --ids "123,456,789"

# Media upload
ctw media upload --file path/to/image.jpg --category tweet_image
ctw tweets create --text "check this out" --media-ids "1234567890"
```

## Testing

Unit tests are colocated with their services and rely on `httptest.Server` helpers for determinism. Run the full suite with:

```bash
go test ./...
```

No tests hit live Twitter endpoints. If you introduce integration tests, guard them with `t.Skip` unless credentials are configured.

## Project Structure

```text
cmd/ctw          # Cobra CLI entrypoint and commands
internal/client  # Shared HTTP client abstraction
internal/data    # Shared DTOs
internal/tweet   # Tweet-related services (lookup, publish, timelines, likes, retweets, bookmarks)
internal/users   # User lookup & relationship services
internal/media   # Media upload service (chunked upload to upload.twitter.com)
internal/dm      # Direct messages service
script/sh        # Shell scripts for testing and examples
```

## Contributing

- Follow the existing service pattern: accept `context.Context`, build query maps, decode into typed structs, return `client.RateLimitSnapshot` alongside responses.
- Keep the CLI thin—parse flags, call a service, print JSON, and surface rate-limit metadata via `printRateLimits`.

Command-line toolkit for exploring Twitter v2 REST endpoints. The project is written in Go 1.18 and ships as a Cobra-based CLI bundled with reusable service packages under `internal/`.

## Prerequisites

- Go 1.18 or newer
- Twitter API bearer token (`BEARER_TOKEN` environment variable)

## Usage

All commands read credentials from the `BEARER_TOKEN` environment variable. You can override authentication or base URL per invocation via flags (`--bearer-token`, `--base-url`, `--user-agent`).

### Filtered Stream

```bash
# List existing rules
ctw stream rules list

# Add a dry-run rule
ctw stream rules add --value "cats has:images" --tag "cat-images" --dry-run

# Start streaming with selected fields
ctw stream --field "tweet.fields=created_at" --field "expansions=author_id"
```

### Recent Search and Counts

```bash
# Fetch recent tweets
ctw search recent --query "golang lang:en" --param "max_results=10"

# Aggregate tweet counts
ctw counts recent --query "golang" --granularity hour
ctw counts all --query "from:TwitterDev" --granularity day
```

### Users

```bash
# Lookup users
ctw users lookup --username TwitterDev
ctw users lookup --ids "2244994945,6253282"

# Mutate relationships
ctw users follow --source-id 1 --target-id 2
ctw users block --source-id 1 --target-id 3
```

### Direct Messages

```bash
# Send a DM to a user
ctw dms send --user-id 2244994945 --text "Hello from ctw"

# List DM events with pagination
ctw dms list --param "pagination_token=some-token"

# Delete a DM event
ctw dms delete --id event-123
```

### Tweets

```bash
# Create and delete tweets
ctw tweets create --text "Hello Twitter"
ctw tweets delete --id 1234567890

# Fetch tweets by ID
ctw tweets get --id 1234567890
ctw tweets get --ids "123,456,789" --param "tweet.fields=created_at"
```

### Likes

```bash
# Like and unlike tweets
ctw likes add --user-id 123 --tweet-id 456
ctw likes remove --user-id 123 --tweet-id 456

# List liked tweets
ctw likes list --user-id 123 --param max_results=20
```

### Retweets

```bash
# Retweet and unretweet
ctw retweets add --user-id 123 --tweet-id 789
ctw retweets remove --user-id 123 --tweet-id 789

# List retweeters
ctw retweets list --tweet-id 789
```

### Bookmarks

```bash
# Add and remove bookmarks
ctw bookmarks add --user-id 123 --tweet-id 999
ctw bookmarks remove --user-id 123 --tweet-id 999

# List bookmarks
ctw bookmarks list --user-id 123
```

### Timelines

```bash
# Get user's tweets
ctw timelines user --user-id 123 --param max_results=10

# Get mentions
ctw timelines mentions --user-id 456

# Get home timeline
ctw timelines home --user-id 789
```

## Development

- Services live under `internal/` and accept a shared `client.Client` for HTTP access.
- Unit tests use `httptest.Server` fixtures; run them with `go test ./...`.
- Scripts under `script/sh` contain raw curl examples that mirror the CLI behaviour.
