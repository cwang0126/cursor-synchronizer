# Security Policy

## Supported Versions

`cursor-sync` is pre-1.0 software. Only the latest minor release line receives
security updates; older lines are not patched.

| Version | Supported          |
| ------- | ------------------ |
| 0.2.x   | :white_check_mark: |
| < 0.2   | :x:                |

Once a `1.0` is released, this table will be updated to cover the current
major line and the immediately preceding one.

## Reporting a Vulnerability

Please **do not** open a public GitHub issue for security problems.

Report vulnerabilities privately via GitHub's
["Report a vulnerability"](https://github.com/cwang0126/cursor-synchronizer/security/advisories/new)
flow on this repository. If that isn't available to you, email the maintainer
at the address listed on the GitHub profile
[@cwang0126](https://github.com/cwang0126) instead.

When reporting, please include:

- A description of the issue and its impact.
- The `cursor-sync` version (`cursor-sync --help` shows it) and your OS.
- Steps to reproduce, or a minimal proof-of-concept.

### What to expect

- **Acknowledgement:** within 7 days of your report.
- **Status updates:** at least every 14 days until the issue is resolved or
  closed.
- **Accepted reports:** we'll agree on a disclosure timeline (typically up to
  90 days), prepare a fix, cut a patched release, and publish a GitHub
  Security Advisory crediting you unless you prefer to remain anonymous.
- **Declined reports:** you'll get a written explanation of why the issue is
  out of scope or not considered a vulnerability.

### Scope

In scope:

- The `cursor-sync` CLI itself (clone / pull / list / config commands).
- The published `install.sh` and `install.ps1` installers.
- The release artifacts attached to GitHub Releases.

Out of scope:

- Vulnerabilities in upstream dependencies (`git`, Go standard library, etc.)
  — please report those to their maintainers. If the issue is in how
  `cursor-sync` *uses* a dependency, that is in scope.
- Content of third-party `.cursor/` repositories you choose to sync from.
