#!/usr/bin/env bash
# Install atlassian-cli binaries for the current platform.
# Downloads a single archive (jira + confluence) from GitHub Releases.
#
# Usage:
#   ./install.sh              # install latest version
#   ./install.sh v1.2.3       # install specific version

set -euo pipefail

REPO="PigZyj2333/atlassian-cli"
VERSION="${1:-latest}"
SKILL_DIR="$(cd "$(dirname "$0")/.." && pwd)"

# --- Detect platform ---
case "$(uname -s)" in
  Linux*)               _OS="linux" ;;
  Darwin*)              _OS="darwin" ;;
  CYGWIN*|MINGW*|MSYS*) _OS="windows" ;;
  *)  echo "Error: Unsupported OS: $(uname -s)" >&2; exit 1 ;;
esac

case "$(uname -m)" in
  x86_64|amd64)  _ARCH="amd64" ;;
  aarch64|arm64) _ARCH="arm64" ;;
  *)  echo "Error: Unsupported architecture: $(uname -m)" >&2; exit 1 ;;
esac

_EXT="tar.gz"
[ "$_OS" = "windows" ] && _EXT="zip"

_BIN_EXT=""
[ "$_OS" = "windows" ] && _BIN_EXT=".exe"

PLATFORM="${_OS}_${_ARCH}"
BIN_DIR="$SKILL_DIR/bin/$PLATFORM"

echo "Platform: ${_OS}/${_ARCH}"

# --- Resolve version ---
if [ "$VERSION" = "latest" ]; then
  echo "Fetching latest release tag..."
  VERSION=$(curl -sI "https://github.com/$REPO/releases/latest" \
    | grep -i "^location:" | sed 's|.*/tag/||' | tr -d '\r\n')
  if [ -z "$VERSION" ]; then
    echo "Error: Could not determine latest version. Specify manually: ./install.sh v1.0.0" >&2
    exit 1
  fi
fi

echo "Version: $VERSION"

# --- Check if already installed ---
VERSION_FILE="$BIN_DIR/.version"
if [ -f "$VERSION_FILE" ] && [ "$(cat "$VERSION_FILE")" = "$VERSION" ]; then
  echo "Already installed: $VERSION"
  echo "  jira:       $BIN_DIR/jira${_BIN_EXT}"
  echo "  confluence: $BIN_DIR/confluence${_BIN_EXT}"
  exit 0
fi

# --- Download single archive ---
mkdir -p "$BIN_DIR"
TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

CLEAN_VERSION="${VERSION#v}"
ARCHIVE="atlassian-cli_${CLEAN_VERSION}_${_OS}_${_ARCH}.${_EXT}"
URL="https://github.com/$REPO/releases/download/${VERSION}/${ARCHIVE}"

echo "Downloading atlassian-cli ($PLATFORM)..."
echo "  $URL"

if ! curl -fsSL -o "$TMPDIR/$ARCHIVE" "$URL"; then
  echo "Error: Failed to download $URL" >&2
  echo "Check that version '$VERSION' exists at https://github.com/$REPO/releases" >&2
  exit 1
fi

# --- Extract ---
if [ "$_EXT" = "zip" ]; then
  unzip -qo "$TMPDIR/$ARCHIVE" -d "$TMPDIR/extracted"
else
  mkdir -p "$TMPDIR/extracted"
  tar -xzf "$TMPDIR/$ARCHIVE" -C "$TMPDIR/extracted"
fi

# --- Copy binaries ---
for TOOL in jira confluence; do
  BINARY=$(find "$TMPDIR/extracted" -name "${TOOL}${_BIN_EXT}" -type f | head -1)
  if [ -z "$BINARY" ]; then
    echo "Error: '${TOOL}${_BIN_EXT}' not found in archive" >&2
    exit 1
  fi
  cp "$BINARY" "$BIN_DIR/${TOOL}${_BIN_EXT}"
  chmod +x "$BIN_DIR/${TOOL}${_BIN_EXT}" 2>/dev/null || true
done

echo "$VERSION" > "$VERSION_FILE"

echo ""
echo "Installed successfully:"
echo "  jira:       $BIN_DIR/jira${_BIN_EXT}"
echo "  confluence: $BIN_DIR/confluence${_BIN_EXT}"
