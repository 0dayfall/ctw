# Copilot Instructions for ctw

## Architecture & Scope
- **Purpose**: Go 1.18 toolkit plus Cobra CLI for Twitter v2 endpoints. Runtime packages sit under `internal/`; the shipping binary lives in `cmd/ctw`.
- **HTTP client**: `internal/client` exposes a configurable `Client` that handles base URL resolution, auth headers, rate-limit parsing, and JSON error decoding. Always prefer this over direct `http.Client` use.
- **Shared DTOs**: `internal/data/json.go` carries common tweet entity structs. Keep endpoint-specific envelopes near their services (`internal/tweet/...`, `internal/users/...`).

## Client & Service Pattern
- Construct services with `client.New` and dependency-inject them (`filteredstream.NewService`, `recentsearch.NewService`, `recentcount.NewService`, `users/lookup.NewService`). Tests rely on this to swap in `httptest.Server` URLs.
- Each service returns a typed response plus `client.RateLimitSnapshot`. Propagate rate-limit data when adding new flows so the CLI can surface it.
- Use `client.CheckResponse` and `client.SafeClose` instead of manual error logging; they centralize API error formatting and body cleanup.

## Feature Packages
- **Filtered stream**: `internal/tweet/filteredstream` contains rule management and stream DTOs. Tests use `newTestService` helpers—mirror that for new methods.
- **Recent search & counts**: `internal/tweet/recentsearch` and `internal/tweet/recentcount` expose `SearchRecent`/`GetRecentCount` methods that accept query param maps; reuse this pattern for additional pagination or expansions.
- **Users domain**: `internal/users/lookup` now wraps lookup/block/follow endpoints in a single service. Legacy `block/` and `follow/` helpers are deprecated—avoid reviving them.
- **Tweet publish**: `internal/tweet/publish` implements create/delete operations. Request structs (`CreateTweetRequest`) stay minimal—extend them cautiously when adding optional API fields.
- **Sampled stream**: `internal/tweet/sampledstream` remains a stub; fix the malformed URL if you plan to ship sampled stream support.

## CLI Commands (`cmd/ctw`)
- Root command wires `--bearer-token`, `--base-url`, and `--user-agent` flags into `client.Config`.
- Subcommands: `stream`, `stream rules add|list`, `search recent`, `counts recent|all`, `users` (lookup, block, unblock, follow, unfollow), and `tweets` (create, delete). Follow existing usage patterns when adding new commands.
- Helper utilities (`helpers.go`) contain `printJSON`, `printRateLimits`, and `parseKeyValuePairs`; reuse them instead of duplicating formatting logic.

## Testing & Tooling
- Unit tests live alongside services and rely on `newTestService` helpers to spin up `httptest.Server` instances. Stick to that approach for deterministic testing.
- Integration-style tests have been removed; `go test ./...` should pass without hitting the real Twitter API. Provide skips if you reintroduce live calls.
- Cobra adds indirect dependencies (`spf13/pflag`, `mousetrap`). Run `go test ./...` after dependency updates to ensure `go.sum` stays in sync.

## Developer Workflows
- Export `BEARER_TOKEN` (and optional `USER_AGENT`) before running the CLI. Scripts under `script/sh` remain handy for quick curl smoke tests.
- When adding endpoints, create a new service method that accepts a `context.Context`, builds query maps, decodes into typed structs, and returns rate-limit metadata.
- Keep CLI commands thin: parse flags, build param maps, call the relevant service, print JSON, then log rate-limit info to stderr.
