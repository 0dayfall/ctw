#!/usr/bin/env bash
set -euo pipefail

# Append to a daily log file.
outfile="watch_$(date +%Y%m%d).jsonl"
ctw watch --keyword "golang" --auto-setup --json | tee -a "$outfile"
