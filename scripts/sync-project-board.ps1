$ErrorActionPreference = "Stop"
$repoRoot = Split-Path -Parent $PSScriptRoot
Push-Location $repoRoot
try {
    bash "$PSScriptRoot/sync-project-board.sh"
} finally {
    Pop-Location
}
