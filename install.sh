#!/usr/bin/env bash
# install.sh — build and install the helios-cli binary.
#
# Usage:
#   ./install.sh                # installs to $HOME/.local/bin (default)
#   PREFIX=/usr/local ./install.sh
#
# Environment variables:
#   PREFIX     install root prefix; the binary is placed under $PREFIX/bin (default: $HOME/.local)
#   BIN_NAME   binary name (default: helios-cli)

set -euo pipefail

PREFIX="${PREFIX:-$HOME/.local}"
BIN_NAME="${BIN_NAME:-helios-cli}"
INSTALL_DIR="$PREFIX/bin"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if ! command -v go >/dev/null 2>&1; then
  echo "error: 'go' is required but not found in PATH" >&2
  exit 1
fi

mkdir -p "$INSTALL_DIR"

echo "==> building $BIN_NAME"
( cd "$SCRIPT_DIR" && go build -o "$INSTALL_DIR/$BIN_NAME" . )

echo "==> installed to $INSTALL_DIR/$BIN_NAME"

case ":$PATH:" in
  *":$INSTALL_DIR:"*) ;;
  *)
    echo "note: $INSTALL_DIR is not in your PATH." >&2
    echo "      add it via: export PATH=\"$INSTALL_DIR:\$PATH\"" >&2
    ;;
esac
