# Microsoft WDSI false positive — disk-tool NSIS installer

Prepared for submission at [Microsoft WDSI file submission](https://www.microsoft.com/en-us/wdsi/filesubmission).

## File to submit

| Field | Value |
|-------|-------|
| **File** | Tauri NSIS setup `.exe` from `desktop-windows` release artifact |
| **Path in bundle** | `*_x64-setup.exe` (under `bundle/nsis/`) |
| **Product name** | disk-tool |
| **Product version** | 1.5.0 |
| **Publisher** | jcuel — https://github.com/jcuel/disk-tool |
| **Role** | Software developer |
| **Download URL** | https://github.com/jcuel/disk-tool/releases/latest |

## SHA-256

```
<!-- FILL AFTER RELEASE BUILD -->
```

PowerShell:

```powershell
Get-FileHash -Algorithm SHA256 .\disk-tool_*_x64-setup.exe
```

## Description (paste into WDSI form)

```
disk-tool is an open-source cross-platform disk usage analyzer (https://github.com/jcuel/disk-tool).
The Windows desktop app is built with Tauri 2 and bundles an embedded Go sidecar server.
The NSIS installer is unsigned (no Authenticode certificate yet). Windows Defender flagged
this as a false positive — heuristic match on unsigned NSIS + stripped Go binary + sidecar
process spawn. CI ClamAV scan on the same artifact reported clean. No malware intended or present.
```

## Evidence links

| Item | URL |
|------|-----|
| Release workflow run | <!-- FILL: e.g. https://github.com/jcuel/disk-tool/actions/runs/NNNN --> |
| ClamAV report artifact | `clamav-report-release-*` from the same run |
| Security policy | https://github.com/jcuel/disk-tool/blob/dev/SECURITY.md#windows-defender--smartscreen-false-positives |

## Submission checklist

- [ ] NSIS `.exe` uploaded (not CLI `disk-tool.exe` alone)
- [ ] SHA-256 recorded above
- [ ] ClamAV release report attached or linked
- [ ] Submitted as **Software developer**
- [ ] Optional: same file uploaded to VirusTotal for public hash record

## After Microsoft clears the file

1. Publish GitHub Release (draft → latest) with the scanned artifacts.
2. Re-test download on a clean Windows VM with Defender enabled.
3. Track Authenticode signing for long-term SmartScreen reputation.
