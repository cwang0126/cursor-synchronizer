# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.3.1] - 2026-04-23 Windows Terminal arrow-key fix

### Fixed
- Arrow keys (↑/↓/←/→) no longer print as literal `[A` / `[B` / `[C` /
  `[D` in the `cursor-sync clone` multi-select prompt on Windows 11 /
  Windows Terminal. No effect on macOS or Linux.

## [v0.3.0] - 2026-04-22 Customizable remote source folder

### Added
- `--folder` / `-f` flag on `cursor-sync clone` and `cursor-sync pull` to
  point at any directory (relative to the repo root) that contains the
  `rules/`, `skills/`, and `commands/` groupings. Useful when the remote
  keeps them under a non-standard path like `configs/cursor/`. When
  omitted, the existing auto-detection (`.cursor/` → `cursor/` → repo
  root) is preserved.
- `folder` field in `.cursor-sync/config.yaml`. `clone --folder <path>`
  records the value; `pull` reads it as the default (and falls back to
  auto-detection when empty), writing back any new value passed via
  `--folder`.
- `cursor-sync config --show folder` / `--set folder <path>` for managing
  the new field.

## [v0.2.0] - 2026-04-22 Auto-detection feature on default branch name

### Added
- `cursor-sync clone` now auto-detects the default branch when `--branch` is
  omitted: it tries `master` first and falls back to `main`, returning a
  clear error if neither exists. The resolved branch is written to
  `.cursor-sync/config.yaml` so subsequent `cursor-sync pull` runs go
  straight to the right ref.

### Changed
- Replace ASCII art banner style.
- `--branch` flag default flipped from a hard-coded `master` to empty,
  with the help text updated to `(default: try master, then main)`.

## [v0.1.0] - 2026-04-21 Initial Release

Initial release of `cursor-sync` — a small Go CLI that one-way syncs a
project's `.cursor/` rules, skills, and commands from a remote git repository.

### Added
- `cursor-sync clone <repo-url> [directory]` — shallow-clones the remote,
  multi-selects which `rules/`, `skills/`, and `commands/` entries to import,
  and writes them under `<directory>/.cursor/`. Supports `--all` to import
  everything without prompting and `--branch` to pick a non-default branch.
- `cursor-sync pull [--yes]` — re-syncs only the entries listed in
  `.cursor-sync/manifest.yaml`. Per-file conflicts prompt `y/N/a/s`
  (overwrite / keep / overwrite-all / skip-all); `--yes` overwrites silently.
- `cursor-sync list` — offline listing of `.cursor/` entries tagged
  `[remote]` (tracked in the manifest) or `[local]` (added by the user).
- `cursor-sync config --show <key>` / `--set <key> <value>` — read and write
  `remote` and `branch` in `.cursor-sync/config.yaml`.
- YAML-backed state: `.cursor-sync/config.yaml` (remote URL + branch) and
  `.cursor-sync/manifest.yaml` (list of pulled entries).
- One-way pull model: files added locally under `.cursor/` are never touched.
- Uses the local `git` binary for `--depth 1` clones into a temp dir, so
  authentication piggybacks on your existing git credentials.
- `build.sh` — interactive cross-compiler that writes binaries to
  `bin/<os>_<arch>/cursor-sync[.exe]` for darwin/linux/windows × amd64/arm64.
  Requires Go 1.22+.
- `install.sh` (macOS/Linux) and `install.ps1` (Windows) — install the
  prebuilt binary from `bin/<os>_<arch>/` onto `PATH` after a local
  `git clone`, with no network access. Scripts error cleanly and point at
  `./build.sh` when no matching binary is committed for the host platform.
- Build-from-source instructions for Go 1.22+.

[0.3.1]: https://github.com/cwang0126/cursor-synchronizer/releases/tag/v0.3.1
[0.3.0]: https://github.com/cwang0126/cursor-synchronizer/releases/tag/v0.3.0
[0.2.0]: https://github.com/cwang0126/cursor-synchronizer/releases/tag/v0.2.0
[0.1.0]: https://github.com/cwang0126/cursor-synchronizer/releases/tag/v0.1.0
