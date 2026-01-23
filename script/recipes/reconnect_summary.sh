#!/usr/bin/env bash
set -euo pipefail

# Highlight reconnect/backoff messages and final summary.
ctw watch --keyword "the" --auto-setup --json \
  2> >(grep -E "disconnected:|reconnecting in|Stream summary" >&2)
