#!/usr/bin/env bash
# build.sh - cross-compile cursor-sync binaries into bin/<os>_<arch>/.
#
# Interactive: prompts for a target. Requires Go 1.22+ on PATH.

set -euo pipefail

cd "$(dirname "$0")"

if ! command -v go >/dev/null 2>&1; then
  echo "go not found on PATH. Install Go 1.22+ first." >&2
  exit 1
fi

targets=(
  "darwin/arm64"
  "darwin/amd64"
  "linux/arm64"
  "linux/amd64"
  "windows/arm64"
  "windows/amd64"
)

echo "Select build target:"
for i in "${!targets[@]}"; do
  printf "  %d) %s\n" "$((i + 1))" "${targets[$i]}"
done
printf "  %d) all\n" "$(( ${#targets[@]} + 1 ))"

read -rp "Enter choice [1-$(( ${#targets[@]} + 1 ))]: " choice

selected=()
if [[ "$choice" =~ ^[0-9]+$ ]] && [ "$choice" -ge 1 ] && [ "$choice" -le "${#targets[@]}" ]; then
  selected=("${targets[$((choice - 1))]}")
elif [ "$choice" = "$(( ${#targets[@]} + 1 ))" ]; then
  selected=("${targets[@]}")
else
  echo "invalid choice: $choice" >&2
  exit 1
fi

for target in "${selected[@]}"; do
  os="${target%/*}"
  arch="${target#*/}"
  out_dir="bin/${os}_${arch}"
  out_name="cursor-sync"
  [ "$os" = "windows" ] && out_name="cursor-sync.exe"
  mkdir -p "$out_dir"
  echo "Building $out_dir/$out_name..."
  GOOS="$os" GOARCH="$arch" CGO_ENABLED=0 \
    go build -ldflags="-s -w" -o "$out_dir/$out_name" .
done

echo "Done."
