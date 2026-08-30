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
- Release binaries scanned by CI (ClamAV)

Out of scope:

- Issues requiring physical access to the machine running disk-tool
- Social engineering
- Denial-of-service against a single-user localhost instance (unless exploitable remotely)

## Windows Defender / SmartScreen false positives

**disk-tool release builds are not malware.** Windows Defender and SmartScreen sometimes flag our installers or `.exe` files because of how they are built and distributed, not because they contain a virus.

Common triggers for **unsigned** desktop software:

| Factor | disk-tool |
|--------|-----------|
| **No code signing** | Installers and binaries are not Authenticode-signed yet (no EV/OV cert). Unknown publisher = low reputation score. |
| **NSIS installer** | Tauri bundles Windows releases as NSIS `.exe` installers — a frequent heuristic target for consumer AV. |
| **Stripped Go sidecar** | Release builds use `-ldflags="-s -w"` (no debug symbols). Some engines treat stripped static binaries as suspicious. |
| **Embedded web UI** | The Go binary embeds the full Vite frontend — large, unusual PE layout vs. typical CLI tools. |
| **Sidecar spawn pattern** | The desktop shell launches a child process that binds `127.0.0.1` and serves HTTP — behavior similar to droppers/launchers in heuristic models. |

**What we do in CI**

- **Trivy** — dependency CVEs (filesystem + Docker image)
- **ClamAV** — malware scan on Linux + Windows cross-compiled binaries (every PR) and on full release artifacts (installers included) before publish
- **govulncheck** — Go stdlib/module vulnerabilities

A clean ClamAV result does not guarantee Defender will stay silent (different engines), but it confirms our build pipeline is not shipping known malware signatures.

**If Defender quarantines disk-tool**

1. Check the file hash against the [latest Release](https://github.com/jcuel/disk-tool/releases/latest) artifact you downloaded.
2. Review the CI ClamAV report attached to that release workflow run.
3. Report a false positive to Microsoft if you are satisfied: [Submit a file for analysis](https://www.microsoft.com/en-us/wdsi/filesubmission).
4. For production deployments, plan **Authenticode signing** for Windows installers (reduces SmartScreen warnings; tracked as future work).

Do not disable Defender permanently; use “Allow on device” only after verifying the download source is GitHub Releases for this repository.

### Submitting an NSIS false positive to Microsoft

When Defender flags the **Tauri NSIS installer** (`*_x64-setup.exe` from `desktop/src-tauri/target/release/bundle/nsis/`), submit it as a software developer at [Microsoft WDSI file submission](https://www.microsoft.com/en-us/wdsi/filesubmission).

**Before submitting**

1. Build or download the exact NSIS `.exe` (Release workflow artifact `desktop-windows`, or local `./build.ps1 -Desktop` then `npm run tauri build` in `desktop/`).
2. Confirm the release workflow **Security — ClamAV release artifacts** job passed and download `clamav-report-release-*` from that run.
3. Record SHA-256:

   ```powershell
   Get-FileHash -Algorithm SHA256 .\disk-tool_*_x64-setup.exe
   ```

   ```bash
   sha256sum disk-tool_*_x64-setup.exe
   ```

**Submission fields**

| Field | Value |
|-------|-------|
| Role | Software developer |
| File | NSIS setup `.exe` (not the CLI binary alone) |
| Product name | disk-tool |
| Description | Open-source disk usage analyzer; Tauri 2 desktop shell with embedded Go sidecar; unsigned build — heuristic false positive |
| Publisher | jcuel — https://github.com/jcuel/disk-tool |
| Download URL | https://github.com/jcuel/disk-tool/releases/latest (after publish) |

**Optional:** upload the same file to [VirusTotal](https://www.virustotal.com) for a public hash record.

Maintainers: after Microsoft clears the file, cut a tagged release so users fetch the updated reputation via GitHub Releases. Long-term mitigation is **Authenticode signing** for Windows installers.

## Safe harbor

We support good-faith security research on this repository. Do not access data outside your own systems or exfiltrate user data when testing.
