# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-04-21

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
- Prebuilt static binaries for macOS (amd64/arm64), Linux (amd64/arm64), and
  Windows (amd64), released via GitHub Actions on `v*` tags.
- `install.sh` (macOS/Linux) and `install.ps1` (Windows) to download the right
  prebuilt binary and put it on `PATH`.
- Build-from-source instructions for Go 1.22+.

[0.1.0]: https://github.com/cwang0126/cursor-synchronizer/releases/tag/v0.1.0
