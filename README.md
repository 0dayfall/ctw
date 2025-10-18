# ctw

Command-line toolkit for exploring Twitter v2 REST endpoints. The project is written in Go 1.18 and ships as a Cobra-based CLI bundled with reusable service packages under `internal/`.

## Prerequisites

- Go 1.18 or newer
- Twitter API bearer token (`BEARER_TOKEN` environment variable)

## Installation

```bash
go build -o bin/ctw ./cmd/ctw
```

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

## Development

- Services live under `internal/` and accept a shared `client.Client` for HTTP access.
- Unit tests use `httptest.Server` fixtures; run them with `go test ./...`.
- Scripts under `script/sh` contain raw curl examples that mirror the CLI behaviour.
