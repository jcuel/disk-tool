param(
    [string]$Repo = "jcuel/disk-tool"
)

$ErrorActionPreference = "Stop"
$PatUrl = "https://github.com/settings/tokens/new?scopes=project,repo&description=disk-tool-GH_PROJECT_SYNC"

Write-Host "GitHub fine-grained PATs cannot access user-owned Projects (board #3)."
Write-Host "Use a classic PAT (ghp_...) with scopes: project, repo"
Write-Host ""
Write-Host "Opening token creation page..."
Start-Process $PatUrl

$Pat = Read-Host "Paste classic PAT (ghp_...)" -AsSecureString
$Plain = [Runtime.InteropServices.Marshal]::PtrToStringAuto(
    [Runtime.InteropServices.Marshal]::SecureStringToBSTR($Pat))
if ([string]::IsNullOrWhiteSpace($Plain)) { throw "No token provided." }
if (-not $Plain.StartsWith("ghp_")) { throw "Expected a classic PAT starting with ghp_." }

Write-Host "Setting GH_PROJECT_SYNC on $Repo..."
gh secret set GH_PROJECT_SYNC --repo $Repo --body $Plain

Write-Host "Verifying token can list project items..."
$env:GH_TOKEN = $Plain
gh project item-list 3 --owner jcuel --format json --limit 1 | Out-Null
Write-Host "Project API check passed."
Write-Host "Done. Next merge to master will auto-sync the board."
