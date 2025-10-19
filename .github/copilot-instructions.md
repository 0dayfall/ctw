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
- **Tweet publish**: `internal/tweet/publish` implements create/delete operations. Request structs (`CreateTweetRequest`) stay minimal—extend them cautiously when adding optional API fields. Supports media attachments via `MediaIDs` array.
- **Tweet lookup**: `internal/tweet/lookup` fetches single or multiple tweets by ID with support for expansions and field parameters.
- **Timelines**: `internal/tweet/timelines` provides user tweets, mentions, and reverse chronological home timeline endpoints.
- **Likes**: `internal/tweet/likes` handles like/unlike operations and listing liked tweets.
- **Retweets**: `internal/tweet/retweets` supports retweet/unretweet actions and listing retweeters.
- **Bookmarks**: `internal/tweet/bookmarks` manages bookmark add/remove/list operations.
- **Direct messages**: `internal/dm` wraps DM send/list/delete flows. Follow the existing pattern (service + CLI) and surface rate-limit headers.
- **Media upload**: `internal/media` implements chunked upload to `upload.twitter.com` (different base URL). Uses three-phase protocol: INIT (get media_id), APPEND (send 5MB chunks via multipart/form-data), FINALIZE (complete upload). Supports async video processing with polling loop (max 60 attempts, 5s intervals). Media categories: `tweet_image`, `tweet_video`, `tweet_gif`, `dm_image`, `dm_video`, `dm_gif`. Service struct holds `uploadBaseURL` field for testability.
- **Sampled stream**: `internal/tweet/sampledstream` remains a stub; fix the malformed URL if you plan to ship sampled stream support.

## CLI Commands (`cmd/ctw`)
- Root command wires `--bearer-token`, `--base-url`, and `--user-agent` flags into `client.Config`.
- Subcommands: `stream`, `search recent`, `counts recent|all`, `users` (lookup, block, unblock, follow, unfollow), `tweets` (create, delete, get), `dms` (send, list, delete), `likes` (add, remove, list), `retweets` (add, remove, list), `bookmarks` (add, remove, list), `timelines` (user, mentions, home), and `media upload`.
- Tweet creation accepts `--media-ids` flag (comma-separated) to attach uploaded media. Validation allows text OR media (not both empty).
- Media upload command accepts `--file` (required) and `--category` (optional) flags. Outputs `media_id_string` for use in tweets/DMs.
- Helper utilities (`helpers.go`) contain `printJSON`, `printRateLimits`, and `parseKeyValuePairs`; reuse them instead of duplicating formatting logic.

## Testing & Tooling
- Unit tests live alongside services and rely on `newTestService` helpers to spin up `httptest.Server` instances. Stick to that approach for deterministic testing.
- Integration-style tests have been removed; `go test ./...` should pass without hitting the real Twitter API. Provide skips if you reintroduce live calls.
- Cobra adds indirect dependencies (`spf13/pflag`, `mousetrap`). Run `go test ./...` after dependency updates to ensure `go.sum` stays in sync.

## Developer Workflows
- Export `BEARER_TOKEN` (and optional `USER_AGENT`) before running the CLI. Scripts under `script/sh` remain handy for quick curl smoke tests.
- When adding endpoints, create a new service method that accepts a `context.Context`, builds query maps, decodes into typed structs, and returns rate-limit metadata.
- Keep CLI commands thin: parse flags, build param maps, call the relevant service, print JSON, then log rate-limit info to stderr.
