#!/usr/bin/env bash
# Dev helper: rebuild and restart server when Go sources change (run from server/).
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
WATCH="${1:-$ROOT/../vcmp-go-plugin/plugin $ROOT/../vcmp-go-plugin/vcmp $ROOT/safari}"

if ! command -v fswatch >/dev/null 2>&1; then
  echo "Install fswatch: brew install fswatch" >&2
  exit 1
fi

echo "[dev-watch] watching: $WATCH"
fswatch -o $WATCH | while read -r _; do
  echo "[dev-watch] change detected — rebuilding..."
  if command -v pwsh >/dev/null 2>&1; then
    pwsh -NoProfile -ExecutionPolicy Bypass -File "$ROOT/build.ps1" -StartServer -NoTest || true
  elif command -v powershell.exe >/dev/null 2>&1; then
    powershell.exe -NoProfile -ExecutionPolicy Bypass -File "$ROOT/build.ps1" -StartServer -NoTest || true
  else
    echo "[dev-watch] PowerShell required to build the Windows plugin" >&2
  fi
done
