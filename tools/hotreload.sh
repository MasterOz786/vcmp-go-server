#!/usr/bin/env bash
# Runs after /reload server shuts down: rebuild plugin, deploy, start server64.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
sleep 2

if command -v pwsh >/dev/null 2>&1; then
  exec pwsh -NoProfile -ExecutionPolicy Bypass -File "$ROOT/build.ps1" -StartServer -NoTest
fi
if command -v powershell.exe >/dev/null 2>&1; then
  exec powershell.exe -NoProfile -ExecutionPolicy Bypass -File "$ROOT/build.ps1" -StartServer -NoTest
fi

echo "[hotreload] build.ps1 requires PowerShell (pwsh or powershell.exe)" >&2
exit 1
