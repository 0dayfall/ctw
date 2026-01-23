# Scripts for ctw Automation

Shell scripts demonstrating how to automate Twitter workflows with `ctw`.

## Purpose

These scripts show you how to:
- Validate your environment setup
- Use `ctw` in automation pipelines
- Integrate Twitter data with other tools
- Build real-time monitoring systems

## Structure

- **sh/examples/** - Production-ready automation scripts
- **sh/** - Testing utilities and low-level API examples

## Quick Start

1. **Validate your setup:**
   ```bash
   ./script/sh/test_env.sh
   ```

2. **Run an automation example:**
   ```bash
   ./script/sh/examples/watch_golang.sh
   ```

3. **Integrate with your pipeline:**
   ```bash
   ctw search recent --query "AI" | jq -r '.data[].text' | your_script.sh
   ```

## Automation Examples

Located in `sh/examples/` - ready to use in production:

**watch_golang.sh**  
Monitor Go programming discussions in real-time. Shows how to use `ctw watch` with multiple keywords.

**monitor_crypto.sh**  
Track cryptocurrency tweets with filtering. Demonstrates real-time data collection.

**advanced_rules.sh**  
Complex Boolean filtering with Twitter's rule syntax. Shows advanced stream management.

## Testing Utilities

Located in `sh/` - for environment validation and API exploration:

**test_env.sh**  
Validates your setup: checks `BEARER_TOKEN`, finds `ctw` binary, verifies dependencies (`jq`, `curl`).

**validate_twitter_stream.sh**  
Tests API connectivity and lists active stream rules.

**smoke_real_api.sh**  
Opt-in smoke test for real API calls (requires `BEARER_TOKEN`). Intended for release checks.

**connect_to_stream.sh**, **twitter_stream.sh**, **delete_rule.sh**, **block.sh**  
Low-level curl examples showing raw API usage. For production, use the `ctw` CLI instead.

## Usage Patterns

**Pattern 1: Real-Time Monitoring**

```bash
# Monitor and process tweets as they arrive
ctw watch --keyword "urgent" --auto-setup | while read -r line; do
    # Your processing logic here
    echo "$line" | send_to_slack.sh
done
```

**Pattern 2: Data Collection**

```bash
# Collect tweets for analysis
ctw search recent --query "data science" --param "max_results=100" \
    | jq -r '.data[] | [.id, .text, .created_at] | @csv' \
    > tweets.csv
```

**Pattern 3: Scheduled Automation**

```bash
# Add to crontab for scheduled runs
0 */6 * * * /usr/local/bin/ctw search recent --query "myapp" >> /var/log/twitter_mentions.json
```

**Pattern 4: Conditional Actions**

```bash
# Take action based on content
ctw search recent --query "urgent support" | jq -r '.data[].id' | while read -r id; do
    ctw likes add --user-id YOUR_ID --tweet-id "$id"
    # Notify support team
done
```

## Prerequisites

```bash
# Required
export BEARER_TOKEN="your_twitter_bearer_token"

# Optional but recommended
export PATH="$PATH:/path/to/ctw/bin"
```

Dependencies:
- `ctw` binary (built or installed)
- `jq` for JSON processing (recommended)
- `curl` for low-level API scripts

Check with: `./script/sh/test_env.sh`

## Documentation

- **[AUTOMATION.md](../AUTOMATION.md)** - Comprehensive automation guide with real-world examples
- **[STREAMING.md](../STREAMING.md)** - Twitter filtered stream API details  
- **[KEYWORD_SCREENING.md](../KEYWORD_SCREENING.md)** - Keyword filtering syntax reference
- **`ctw --help`** - Built-in CLI documentation
