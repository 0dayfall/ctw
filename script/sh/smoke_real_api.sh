#!/usr/bin/env bash
set -euo pipefail

if [[ -z "${BEARER_TOKEN:-}" ]]; then
  echo "BEARER_TOKEN not set; skipping smoke test."
  exit 0
fi

cleanup_paths=()
trap 'rm -rf "${cleanup_paths[@]}"' EXIT

resolve_ctw_bin() {
  local candidates=()
  local candidate=""

  if [[ -n "${CTW_BIN:-}" ]]; then
    candidates+=("$CTW_BIN")
  fi
  if command -v ctw >/dev/null 2>&1; then
    candidates+=("$(command -v ctw)")
  fi
  if [[ -x "./bin/ctw" ]]; then
    candidates+=("./bin/ctw")
  fi
  if [[ -x "./ctw" ]]; then
    candidates+=("./ctw")
  fi

  for candidate in "${candidates[@]}"; do
    set +e
    "$candidate" --version >/dev/null 2>&1
    status=$?
    set -e
    if [[ $status -eq 0 ]]; then
      echo "$candidate"
      return 0
    fi
  done

  if command -v go >/dev/null 2>&1; then
    local build_dir
    build_dir="$(mktemp -d)"
    cleanup_paths+=("$build_dir")
    go build -o "$build_dir/ctw" ./cmd/ctw
    echo "$build_dir/ctw"
    return 0
  fi

  return 1
}

CTW_BIN="$(resolve_ctw_bin)" || {
  echo "ctw binary not found or not runnable; install Go or set CTW_BIN." >&2
  exit 1
}

temp_dir="$(mktemp -d)"
cleanup_paths+=("$temp_dir")

echo "Running ctw smoke tests with $CTW_BIN"

echo "-> users lookup"
"$CTW_BIN" users lookup --usernames "twitter" > "$temp_dir/users.json"

echo "-> search recent"
"$CTW_BIN" search recent --query "golang" --param "max_results=10" > "$temp_dir/search.json"

if command -v jq >/dev/null 2>&1; then
  user_id="$(jq -r 'if type=="array" then (.[0].id // empty) else (.data[0].id // empty) end' "$temp_dir/users.json")"
  if [[ -n "$user_id" ]]; then
    echo "-> timelines user"
    "$CTW_BIN" timelines user --user-id "$user_id" --param "max_results=10" > "$temp_dir/timeline.json"
  else
    echo "No user id found in users lookup; skipping timelines test."
  fi
else
  echo "jq not available; skipping timelines test."
fi

echo "Smoke tests passed."
