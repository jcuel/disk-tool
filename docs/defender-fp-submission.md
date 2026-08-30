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
4609b4c61d5ada14ca91c48d386062ce2e1e3ca706960086c88a143f67d5cc67  disk-tool_1.4.0_x64-setup.exe  (NSIS — submit this one)
0866488712adcb70a26988c0807d1001ed34a9474846799d4336af4acfddb327  disk-tool_1.4.0_x64_en-US.msi
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
| Release workflow run | https://github.com/jcuel/disk-tool/actions/runs/33291584015 |
| ClamAV scan (local, v1.5.0 artifacts) | Clean — 2 files scanned, 0 infected (see Summary below) |
| Security policy | https://github.com/jcuel/disk-tool/blob/dev/SECURITY.md#windows-defender--smartscreen-false-positives |

## Submission checklist

- [x] NSIS `.exe` built (`disk-tool_1.4.0_x64-setup.exe`, tag v1.5.0)
- [x] SHA-256 recorded above
- [x] ClamAV scan clean (NSIS + MSI)
- [ ] Submitted at WDSI portal (requires maintainer Microsoft account)
- [ ] Optional: VirusTotal upload

## After Microsoft clears the file

1. Publish GitHub Release (draft → latest) with the scanned artifacts.
2. Re-test download on a clean Windows VM with Defender enabled.
3. Track Authenticode signing for long-term SmartScreen reputation.
