# API + Cypress E2E — Windows
param(
    [string]$Bin = (Join-Path $PSScriptRoot "..\bin\disk-tool.exe"),
    [int]$Port = 18081
)

$ErrorActionPreference = "Stop"
$Root = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
$Fixture = Join-Path $Root "testdata\e2e-root"
$Base = "http://127.0.0.1:$Port"

Push-Location (Join-Path $Root "web")
try {
    & "C:\Program Files\nodejs\npm.cmd" ci 2>$null
    if ($LASTEXITCODE -ne 0) {
        node node_modules\vite\bin\vite.js build 2>$null
    } else {
        & "C:\Program Files\nodejs\npm.cmd" install --no-save cypress@14.3.2
        & "C:\Program Files\nodejs\npm.cmd" run build
    }
} finally {
    Pop-Location
}

Remove-Item -Recurse -Force (Join-Path $Root "cmd\disk-tool\static\*") -ErrorAction SilentlyContinue
Copy-Item -Recurse -Force (Join-Path $Root "web\dist\*") (Join-Path $Root "cmd\disk-tool\static\")

if (-not (Test-Path $Bin)) {
    Push-Location $Root
    go build -o $Bin ./cmd/disk-tool
    Pop-Location
}

$server = Start-Process -FilePath $Bin -ArgumentList "serve", "--port", $Port, "--no-open" -PassThru -WindowStyle Hidden
try {
    $ready = $false
    for ($i = 0; $i -lt 30; $i++) {
        try {
            Invoke-RestMethod "$Base/api/roots" | Out-Null
            $ready = $true
            break
        } catch { Start-Sleep -Milliseconds 200 }
    }
    if (-not $ready) { throw "server not ready" }

    Push-Location (Join-Path $Root "web")
    $env:CYPRESS_BASE_URL = $Base
    npx cypress@14.3.2 run --env "scanRoot=$Fixture"
    Pop-Location
    Write-Host "e2e-run: all Cypress specs passed"
} finally {
    if ($server -and -not $server.HasExited) {
        Stop-Process -Id $server.Id -Force -ErrorAction SilentlyContinue
    }
}
