# Security Policy

## Supported versions

| Version | Supported |
|---------|-----------|
| `master` (latest release integration) | Yes |
| `dev` | Best-effort (pre-release) |
| Older tags | No |

## Reporting a vulnerability

**Please do not open a public GitHub issue for security vulnerabilities.**

Report privately using one of:

1. **[GitHub Security Advisories](https://github.com/jcuel/disk-tool/security/advisories/new)** (preferred) — confidential, works for collaborators and external reporters.
2. **Maintainer contact** — if Advisories are unavailable, email or DM the repo owner via GitHub profile contact methods.

Include:

- Description of the issue and impact
- Steps to reproduce
- Affected versions or commits
- Suggested fix (if any)

## Response expectations

- **Acknowledgment:** within 7 days
- **Triage:** severity assessment and planned fix timeline
- **Disclosure:** coordinated after a fix is available (credit given unless you prefer anonymity)

## Scope

In scope:

- `disk-tool` server/API (localhost binding, path validation, scan/delete/cleanup flows)
- Docker image and CI-supplied artifacts
- Dependency vulnerabilities surfaced by CI (Trivy)

Out of scope:

- Issues requiring physical access to the machine running disk-tool
- Social engineering
- Denial-of-service against a single-user localhost instance (unless exploitable remotely)

## Safe harbor

We support good-faith security research on this repository. Do not access data outside your own systems or exfiltrate user data when testing.
