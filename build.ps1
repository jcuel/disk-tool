# Build disk-tool on Windows
param(
    [switch]$Desktop
)

$ErrorActionPreference = "Stop"
Set-Location $PSScriptRoot

function Build-DiskTool {
    Push-Location web
    npm ci
    npm run build
    Pop-Location

    Remove-Item -Recurse -Force cmd\disk-tool\static\* -ErrorAction SilentlyContinue
    Copy-Item -Recurse -Force web\dist\* cmd\disk-tool\static\

    New-Item -ItemType Directory -Force -Path bin | Out-Null
    go test ./...
    go build -o bin\disk-tool.exe .\cmd\disk-tool
    Write-Host "Built bin\disk-tool.exe"
}

function Install-DesktopSidecar {
    param([string]$Source = "bin\disk-tool.exe")
    $target = (rustc -vV | Select-String '^host: (.+)$').Matches.Groups[1].Value
    $destDir = "desktop\src-tauri\binaries"
    New-Item -ItemType Directory -Force -Path $destDir | Out-Null
    $dest = Join-Path $destDir "disk-tool-$target.exe"
    Copy-Item -Force $Source $dest
    Write-Host "Installed sidecar: $dest"
}

function Build-Desktop {
    Build-DiskTool
    Install-DesktopSidecar
    if (-not (Test-Path desktop\src-tauri\icons\icon.ico)) {
        Push-Location desktop
        npm ci
        npx tauri icon src-tauri/icons/icon.png
        Pop-Location
    }
    Push-Location desktop
    npm ci
    npm run build
    Pop-Location
    Write-Host "Desktop build complete"
}

if ($Desktop) {
    Build-Desktop
} else {
    Build-DiskTool
}
