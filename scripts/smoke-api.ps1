# API smoke test — Windows
param(
    [string]$Bin = (Join-Path $PSScriptRoot "..\bin\disk-tool.exe"),
    [int]$Port = 18080
)

$ErrorActionPreference = "Stop"
$Root = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
$Base = "http://127.0.0.1:$Port"

if (-not (Test-Path $Bin)) {
    throw "Binary not found: $Bin"
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
    Write-Host "OK /api/roots"

    $body = @{ root = $Root } | ConvertTo-Json
    $scan = Invoke-RestMethod -Method POST -Uri "$Base/api/scans" -ContentType "application/json" -Body $body
    $scanId = $scan.scanId
    Write-Host "OK POST /api/scans ($scanId)"

    $status = ""
    for ($i = 0; $i -lt 60; $i++) {
        $job = Invoke-RestMethod "$Base/api/scans/$scanId"
        $status = $job.status
        if ($status -eq "completed") { break }
        Start-Sleep -Milliseconds 500
    }
    if ($status -ne "completed") { throw "scan not completed: $status" }
    Write-Host "OK GET /api/scans/{id} completed"

    if (-not $job.insights) { throw "missing insights" }
    if (-not $job.tree) { throw "missing tree" }
    Write-Host "OK insights + tree"

    if ($job.tree.children -and $job.tree.children.Count -gt 0) {
        $path = $job.tree.children[0].path
        $expand = @{ path = $path; depth = 5 } | ConvertTo-Json
        Invoke-RestMethod -Method POST -Uri "$Base/api/scans/$scanId/expand" -ContentType "application/json" -Body $expand | Out-Null
        Write-Host "OK POST /api/scans/{id}/expand"
    }

    $ticket = Invoke-WebRequest "$Base/api/scans/$scanId/export?format=ticket" -UseBasicParsing
    if ($ticket.Content -notmatch "Disk usage report") { throw "ticket export failed" }
    Write-Host "OK export ticket"

    $ui = Invoke-WebRequest "$Base/" -UseBasicParsing
    if ($ui.StatusCode -ne 200) { throw "UI status $($ui.StatusCode)" }
    Write-Host "OK UI /"

    Write-Host "smoke-api: all checks passed"
} finally {
    if ($server -and -not $server.HasExited) {
        Stop-Process -Id $server.Id -Force -ErrorAction SilentlyContinue
    }
}
