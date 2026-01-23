#!/usr/bin/env bash
set -euo pipefail

# Stream tweets as JSON Lines (good for ingestion).
ctw watch --keyword "golang" --auto-setup --json > watch.jsonl
