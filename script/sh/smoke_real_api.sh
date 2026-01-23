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

run_ctw_to_file() {
  local out_file="$1"
  shift
  local err_file="$temp_dir/ctw_err.log"

  set +e
  "$CTW_BIN" "$@" > "$out_file" 2> "$err_file"
  local status=$?
  set -e

  if [[ $status -ne 0 ]]; then
    if grep -q "status 429" "$err_file"; then
      echo "Rate limited (429); skipping smoke tests."
      exit 0
    fi
    cat "$err_file" >&2
    return "$status"
  fi
}

echo "Running ctw smoke tests with $CTW_BIN"

echo "-> users lookup"
run_ctw_to_file "$temp_dir/users.json" users lookup --usernames "twitter"

echo "-> search recent"
run_ctw_to_file "$temp_dir/search.json" search recent --query "golang" --param "max_results=10"

if command -v jq >/dev/null 2>&1; then
  user_id="$(jq -r 'if type=="array" then (.[0].id // empty) else (.data[0].id // empty) end' "$temp_dir/users.json")"
  if [[ -n "$user_id" ]]; then
    echo "-> timelines user"
    run_ctw_to_file "$temp_dir/timeline.json" timelines user --user-id "$user_id" --param "max_results=10"
  else
    echo "No user id found in users lookup; skipping timelines test."
  fi
else
  echo "jq not available; skipping timelines test."
fi

echo "Smoke tests passed."
