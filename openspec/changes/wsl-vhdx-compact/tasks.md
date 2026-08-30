# Tasks: wsl-vhdx-compact

**Change:** `wsl-vhdx-compact`

## Implementation Checklist

- [x] `internal/maintenance/wsl/wsl_windows.go` + `stub.go`
- [x] `GET /api/wsl/disks`, `POST /api/wsl/compact`
- [x] WSL compact UI flow in `web/src/main.ts`
- [x] Maintenance hub WSL card (Windows when disks found)
- [x] Delta spec + testing report

## Notes

Manual Windows QA requires Hyper-V or diskpart; elevation may be required on some systems.
