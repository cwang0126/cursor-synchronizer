# cursor-synchronizer

A small CLI that one-way syncs your project's `.cursor/` rules, skills, and
commands from a remote git repository so every developer on the team works with
the same Cursor agent context.

- **Single static binary** (Go, no runtime dependencies).
- **Simple model:** `.cursor-sync/config.yaml` records the remote;
`.cursor-sync/manifest.yaml` records which entries were pulled.
- **One-way pull only.** Files you add yourself under `.cursor/` are never
touched by `cursor-sync`.

## Install

### macOS / Linux

```bash
curl -fsSL https://raw.githubusercontent.com/cwang0126/cursor-synchronizer/master/install.sh | bash
```

### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/cwang0126/cursor-synchronizer/master/install.ps1 | iex
```

Both scripts download the appropriate prebuilt binary from the latest GitHub
Release, drop it into `/usr/local/bin` (or `~/.local/bin` as a fallback), and
make sure it's on your `PATH`.

### Build from source

Requires Go 1.22+.

```bash
git clone https://github.com/cwang0126/cursor-synchronizer.git
cd cursor-synchronizer
go build -ldflags="-s -w" -o cursor-sync .
./cursor-sync --help
```

To use `cursor-sync` from any project directory, install the freshly built
binary somewhere on your `PATH`:

```bash
# Preferred: user-local, no sudo required.
mkdir -p "$HOME/.local/bin"
install -m 0755 cursor-sync "$HOME/.local/bin/cursor-sync"

# If ~/.local/bin isn't on your PATH yet, add it (zsh shown):
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc
exec $SHELL -l

# Verify:
which cursor-sync
cursor-sync --help
```

Prefer a system-wide install? Use `sudo install -m 0755 cursor-sync /usr/local/bin/cursor-sync` instead.

## Usage

### `cursor-sync clone <repo-url> [directory]`

Shallow-clones the remote, lets you multi-select which `rules/`, `skills/`,
and `commands/` entries to import, then writes them under
`<directory>/.cursor/` (defaulting to the current folder). It also creates
`<directory>/.cursor-sync/config.yaml` and `manifest.yaml`.

```bash
# Into the current directory:
cursor-sync clone https://github.com/example/cursor-config.git

# Into a new subdirectory:
cursor-sync clone https://github.com/example/cursor-config.git my-project

# Skip the prompt and import everything:
cursor-sync clone --all https://github.com/example/cursor-config.git

# Use a non-default branch:
cursor-sync clone --branch dev https://github.com/example/cursor-config.git
```

The multi-select prompt uses arrow keys to navigate, space to toggle, and
Enter to confirm. The first option, `[Select All]`, picks everything.

### `cursor-sync pull [--yes]`

Re-syncs only the entries previously listed in `.cursor-sync/manifest.yaml`.
On per-file conflicts you'll be asked:

```
Overwrite rules/foo.mdc? [y/N/a/s]
```

- `y` overwrite this file
- `N` keep local (default)
- `a` overwrite all remaining conflicts
- `s` skip all remaining conflicts

`--yes` (`-y`) overwrites everything without prompting.

### `cursor-sync list`

Lists every entry under `.cursor/` and tags each one:

- `[remote]` — name appears in the manifest, pulled from the remote
- `[local]` — added by you, not tracked by `cursor-sync`

This is offline; no network call is made.

### `cursor-sync config`

```bash
cursor-sync config --show remote
cursor-sync config --set remote https://github.com/example/cursor-config.git
cursor-sync config --set branch dev
```

## How it works

```
your-project/
├── .cursor/                 # synced (rules/, skills/, commands/)
└── .cursor-sync/
    ├── config.yaml          # remote URL + branch
    └── manifest.yaml        # files pulled from the remote
```

`cursor-sync` shells out to your local `git` binary to do a `--depth 1` clone
into a temp directory, copies just the entries you selected, then deletes the
temp clone. Authentication uses your existing git credentials (SSH keys,
`gh auth`, `~/.gitconfig`, etc.); `cursor-sync` itself stores no secrets.

### Supported remote layouts

The remote repo's source layout can be any of the following. Detection runs
in this order and picks the first match:

```
<repo>/.cursor/{rules,skills,commands}/   # preferred
<repo>/cursor/{rules,skills,commands}/
<repo>/{rules,skills,commands}/           # groupings at the repo root
```

Regardless of which layout the remote uses, files always land locally under
`<project>/.cursor/` because that's what Cursor itself reads.

## License

MIT