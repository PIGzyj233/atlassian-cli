#!/usr/bin/env bash
# Resolve the correct atlassian-cli binary for the current OS/arch.
# After sourcing, use $JIRA and $CONFLUENCE as command prefixes.
#
# Usage:
#   source "<skill-dir>/scripts/resolve-bin.sh"
#   $JIRA issue get PROJ-123
#   $CONFLUENCE page get 12345678

set -euo pipefail

SKILL_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# --- Detect platform ---
case "$(uname -s)" in
  Linux*)            _OS="linux" ;;
  Darwin*)           _OS="darwin" ;;
  CYGWIN*|MINGW*|MSYS*) _OS="windows" ;;
  *)  echo "Error: Unsupported OS: $(uname -s)" >&2; return 1 2>/dev/null || exit 1 ;;
esac

case "$(uname -m)" in
  x86_64|amd64)   _ARCH="amd64" ;;
  aarch64|arm64)   _ARCH="arm64" ;;
  *)  echo "Error: Unsupported arch: $(uname -m)" >&2; return 1 2>/dev/null || exit 1 ;;
esac

_BIN_EXT=""
[ "$_OS" = "windows" ] && _BIN_EXT=".exe"

BIN_DIR="$SKILL_DIR/bin/${_OS}_${_ARCH}"
JIRA="$BIN_DIR/jira${_BIN_EXT}"
CONFLUENCE="$BIN_DIR/confluence${_BIN_EXT}"

# Auto-install if missing
if [ ! -f "$JIRA" ] || [ ! -f "$CONFLUENCE" ]; then
  echo "Binaries not found. Running install script..." >&2
  bash "$SKILL_DIR/scripts/install.sh"
fi

export JIRA CONFLUENCE
