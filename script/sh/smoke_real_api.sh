#!/usr/bin/env bash
set -euo pipefail

if [[ -z "${BEARER_TOKEN:-}" ]]; then
  echo "BEARER_TOKEN not set; skipping smoke test."
  exit 0
fi

CTW_BIN="${CTW_BIN:-}"
if [[ -z "$CTW_BIN" ]]; then
  if command -v ctw >/dev/null 2>&1; then
    CTW_BIN="$(command -v ctw)"
  elif [[ -x "./bin/ctw" ]]; then
    CTW_BIN="./bin/ctw"
  elif [[ -x "./ctw" ]]; then
    CTW_BIN="./ctw"
  else
    echo "ctw binary not found (set CTW_BIN or add to PATH)." >&2
    exit 1
  fi
fi

temp_dir="$(mktemp -d)"
trap 'rm -rf "$temp_dir"' EXIT

echo "Running ctw smoke tests with $CTW_BIN"

echo "-> users lookup"
"$CTW_BIN" users lookup --usernames "twitter" > "$temp_dir/users.json"

echo "-> search recent"
"$CTW_BIN" search recent --query "golang" --param "max_results=5" > "$temp_dir/search.json"

if command -v jq >/dev/null 2>&1; then
  user_id="$(jq -r '.data[0].id // empty' "$temp_dir/users.json")"
  if [[ -n "$user_id" ]]; then
    echo "-> timelines user"
    "$CTW_BIN" timelines user --user-id "$user_id" --param "max_results=5" > "$temp_dir/timeline.json"
  else
    echo "No user id found in users lookup; skipping timelines test."
  fi
else
  echo "jq not available; skipping timelines test."
fi

echo "Smoke tests passed."
