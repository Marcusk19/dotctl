#!/usr/bin/env bash
set -euo pipefail

REPO="Marcusk19/dotctl"
BINARY="dotctl"

# Detect OS
OS=$(uname -s)
case "$OS" in
  Linux)  OS="Linux" ;;
  Darwin) OS="Darwin" ;;
  *)
    echo "Unsupported OS: $OS"
    exit 1
    ;;
esac

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
  x86_64)          ARCH="x86_64" ;;
  aarch64 | arm64) ARCH="arm64" ;;
  i386 | i686)     ARCH="i386" ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

# Determine archive name and format
if [ "$OS" = "Windows" ]; then
  ARCHIVE="${BINARY}_${OS}_${ARCH}.zip"
else
  ARCHIVE="${BINARY}_${OS}_${ARCH}.tar.gz"
fi

DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${ARCHIVE}"

echo "Downloading dotctl for ${OS}/${ARCH}..."
echo "URL: ${DOWNLOAD_URL}"

# Download to a temp dir
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

curl -fsSL "$DOWNLOAD_URL" -o "$TMP_DIR/$ARCHIVE"

# Extract
cd "$TMP_DIR"
if [[ "$ARCHIVE" == *.tar.gz ]]; then
  tar xzf "$ARCHIVE"
else
  unzip -q "$ARCHIVE"
fi

# Determine install location
if [ -w "/usr/local/bin" ]; then
  INSTALL_DIR="/usr/local/bin"
else
  INSTALL_DIR="$HOME/.local/bin"
  mkdir -p "$INSTALL_DIR"
fi

# Install binary
mv "$TMP_DIR/$BINARY" "$INSTALL_DIR/$BINARY"
chmod +x "$INSTALL_DIR/$BINARY"

echo ""
echo "dotctl installed to $INSTALL_DIR/$BINARY"

# PATH hint if using ~/.local/bin
if [ "$INSTALL_DIR" = "$HOME/.local/bin" ]; then
  if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
    echo ""
    echo "NOTE: Add ~/.local/bin to your PATH if not already done:"
    echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
  fi
fi

# Run bootstrap if a repo URL was provided
if [ -n "${1:-}" ]; then
  echo ""
  echo "Running: dotctl bootstrap $1"
  "$INSTALL_DIR/$BINARY" bootstrap "$1"
fi
