#!/usr/bin/env bash
# install.sh - download the latest cursor-sync binary for macOS or Linux.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/cwang0126/cursor-synchronizer/main/install.sh | bash
#
# Env vars:
#   REPO       owner/name of the GitHub repo (default: cwang0126/cursor-synchronizer)
#   VERSION    release tag to install (default: latest)
#   INSTALL_DIR install location (default: /usr/local/bin, falls back to ~/.local/bin)

set -euo pipefail

REPO="${REPO:-cwang0126/cursor-synchronizer}"
VERSION="${VERSION:-latest}"
INSTALL_DIR="${INSTALL_DIR:-}"

uname_s="$(uname -s)"
case "$uname_s" in
  Darwin) goos="darwin" ;;
  Linux)  goos="linux" ;;
  *) echo "unsupported OS: $uname_s" >&2; exit 1 ;;
esac

uname_m="$(uname -m)"
case "$uname_m" in
  x86_64|amd64) goarch="amd64" ;;
  arm64|aarch64) goarch="arm64" ;;
  *) echo "unsupported architecture: $uname_m" >&2; exit 1 ;;
esac

if [ "$VERSION" = "latest" ]; then
  url_prefix="https://github.com/${REPO}/releases/latest/download"
  # Resolve the actual tag for the file name.
  VERSION="$(curl -fsSL -o /dev/null -w '%{url_effective}' "https://github.com/${REPO}/releases/latest" | sed 's|.*/tag/||')"
else
  url_prefix="https://github.com/${REPO}/releases/download/${VERSION}"
fi

asset="cursor-sync_${VERSION}_${goos}_${goarch}.tar.gz"
url="${url_prefix}/${asset}"

echo "Downloading ${asset}..."
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT
curl -fsSL "$url" -o "$tmp/$asset"
tar -xzf "$tmp/$asset" -C "$tmp"

bin_src="$tmp/cursor-sync_${VERSION}_${goos}_${goarch}/cursor-sync"
chmod +x "$bin_src"

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
  mv "$bin_src" "$dest"
else
  echo "sudo required to write to $INSTALL_DIR"
  sudo mv "$bin_src" "$dest"
fi

echo
echo "Installed cursor-sync ${VERSION} to $dest"
case ":$PATH:" in
  *":$INSTALL_DIR:"*) ;;
  *) echo "Note: add $INSTALL_DIR to your PATH." ;;
esac
echo "Run: cursor-sync --help"
