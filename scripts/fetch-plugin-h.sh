#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
mkdir -p "$ROOT/include"
curl -fsSL "https://raw.githubusercontent.com/habi498/NPC-VCMP/master/plugin/plugin.h" \
  -o "$ROOT/include/plugin.h"
echo "Fetched plugin.h ($(wc -c < "$ROOT/include/plugin.h") bytes)"
