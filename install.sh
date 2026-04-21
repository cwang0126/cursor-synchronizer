#!/usr/bin/env bash
# install.sh - install the prebuilt cursor-sync binary from bin/<os>_<arch>/.
#
# Usage:
#   ./install.sh
#
# Env vars:
#   INSTALL_DIR   install location (default: /usr/local/bin, falls back to ~/.local/bin)

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
INSTALL_DIR="${INSTALL_DIR:-}"

case "$(uname -s)" in
  Darwin) os=darwin ;;
  Linux)  os=linux ;;
  *) echo "unsupported OS: $(uname -s)" >&2; exit 1 ;;
esac

case "$(uname -m)" in
  x86_64|amd64)  arch=amd64 ;;
  arm64|aarch64) arch=arm64 ;;
  *) echo "unsupported architecture: $(uname -m)" >&2; exit 1 ;;
esac

src="$SCRIPT_DIR/bin/${os}_${arch}/cursor-sync"
if [ ! -f "$src" ]; then
  echo "No prebuilt binary for ${os}_${arch} at $src." >&2
  echo "Run ./build.sh to build it, then re-run ./install.sh." >&2
  exit 1
fi

if [ -z "$INSTALL_DIR" ]; then
  if [ -w "/usr/local/bin" ] || [ "$(id -u)" -eq 0 ]; then
    INSTALL_DIR="/usr/local/bin"
  else
    INSTALL_DIR="$HOME/.local/bin"
  fi
fi
mkdir -p "$INSTALL_DIR"

dest="$INSTALL_DIR/cursor-sync"
if [ -w "$INSTALL_DIR" ]; then
  install -m 0755 "$src" "$dest"
else
  echo "sudo required to write to $INSTALL_DIR"
  sudo install -m 0755 "$src" "$dest"
fi

echo
echo "Installed cursor-sync to $dest"
case ":$PATH:" in
  *":$INSTALL_DIR:"*) ;;
  *) echo "Note: add $INSTALL_DIR to your PATH." ;;
esac
echo "Run: cursor-sync --help"
