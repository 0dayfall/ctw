# ctw

A Go 1.18 command-line toolkit for working with selected Twitter v2 REST endpoints.

## Features

- Configurable HTTP client (`internal/client`) with bearer-token auth, user-agent overrides, and rate-limit parsing.
- Service wrappers for filtered stream, recent search, recent counts, and user relationship endpoints under `internal/tweet` and `internal/users`.
- Cobra-powered CLI in `cmd/ctw` with commands: `stream`, `search recent`, `counts recent|all`, `users lookup|block|unblock|follow|unfollow`, and `tweets create|delete`.

## Getting Started

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
   ```

## CLI Examples

```bash
# Stream with selected fields
ctw stream --field tweet.fields=created_at --field expansions=author_id

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
internal/tweet   # Tweet-related services
internal/users   # User lookup & relationship services
script/sh        # Legacy curl helpers (optional)
```

## Contributing

- Follow the existing service pattern: accept `context.Context`, build query maps, decode into typed structs, return `client.RateLimitSnapshot` alongside responses.
- Keep the CLI thin—parse flags, call a service, print JSON, and surface rate-limit metadata via `printRateLimits`.

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
