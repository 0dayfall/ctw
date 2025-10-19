# Script Directory

This directory contains shell scripts for working with the Twitter API and `ctw` CLI.

## Structure

- **sh/** - Shell scripts for direct API calls and CLI usage
  - **examples/** - Ready-to-use example scripts for common use cases
  - Legacy curl-based scripts for testing

## Examples Directory

The `sh/examples/` folder contains practical scripts demonstrating `ctw` usage:

### Stream Monitoring Examples

**watch_golang.sh** - Watch tweets about Go programming in real-time
```bash
./script/sh/examples/watch_golang.sh
```

**monitor_crypto.sh** - Monitor cryptocurrency discussions with multiple keywords
```bash
./script/sh/examples/monitor_crypto.sh
```

**advanced_rules.sh** - Demonstrates complex stream rule management
```bash
./script/sh/examples/advanced_rules.sh
```

## Legacy Scripts

The following curl-based scripts provide low-level API access for testing:

### Stream Management

**connect_to_stream.sh** - Connect to the filtered stream with curl
```bash
./script/sh/connect_to_stream.sh
```
Direct API call to start streaming tweets. Requires active stream rules.

**twitter_stream.sh** - Add stream rules via curl
```bash
./script/sh/twitter_stream.sh
```
Adds sample rules (edit the script to customize). For production, use `ctw stream rules add`.

**validate_twitter_stream.sh** - List active stream rules
```bash
./script/sh/validate_twitter_stream.sh
```
Fetches and displays current rules with jq formatting.

**delete_rule.sh** - Delete a stream rule by ID
```bash
./script/sh/delete_rule.sh RULE_ID
```
Removes a specific rule. Get rule IDs from `validate_twitter_stream.sh`.

### User Management

**block.sh** - Block/unblock users (OAuth 1.0a required)
```bash
./script/sh/block.sh USER_ID [block|unblock]
```
Note: Requires OAuth 1.0a signature. Use `ctw users block/unblock` instead.

### Testing

**test_env.sh** - Validate environment configuration
```bash
./script/sh/test_env.sh
```
Checks for BEARER_TOKEN, ctw binary, and required tools.

## Prerequisites

All scripts require:
- `BEARER_TOKEN` environment variable set
- `ctw` binary built (for example scripts): `go build -o bin/ctw ./cmd/ctw`

```bash
export BEARER_TOKEN="your_twitter_bearer_token_here"
```

## Usage

Make scripts executable:
```bash
chmod +x script/sh/examples/*.sh
chmod +x script/sh/*.sh
```

Run any script:
```bash
./script/sh/examples/watch_golang.sh
```

## See Also

- [STREAMING.md](../../STREAMING.md) - Comprehensive streaming guide
- [KEYWORD_SCREENING.md](../../KEYWORD_SCREENING.md) - Quick reference for keyword screening
