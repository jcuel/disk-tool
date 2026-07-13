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

    $disk = Invoke-RestMethod "$Base/api/disk?path=$([uri]::EscapeDataString($Root))"
    if (-not $disk.total) { throw "disk info failed" }
    Write-Host "OK GET /api/disk"

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

    $docker = Invoke-RestMethod "$Base/api/scans/$scanId/docker"
    if (-not $docker.usage) { throw "docker status missing usage" }
    Write-Host "OK GET /api/scans/{id}/docker"

    $dockerDry = Invoke-RestMethod -Method POST -Uri "$Base/api/scans/$scanId/docker/prune" -ContentType "application/json" -Body (@{
        dryRun = $true
        confirm = $false
        confirmPhrase = ""
    } | ConvertTo-Json)
    if (-not $dockerDry.dryRun) { throw "docker prune dry-run failed" }
    Write-Host "OK POST /api/scans/{id}/docker/prune dry-run"

    $ticket = Invoke-WebRequest "$Base/api/scans/$scanId/export?format=ticket" -UseBasicParsing
    if ($ticket.Content -notmatch "Disk usage report") { throw "ticket export failed" }
    Write-Host "OK export ticket"

    $smokeDir = Join-Path $Root ".smoke-cleanup-$PID"
    New-Item -ItemType Directory -Force -Path (Join-Path $smokeDir "nested") | Out-Null
    Set-Content -Path (Join-Path $smokeDir "nested\file.txt") -Value "test"
    $cleanupBody = @{
        paths = @((Join-Path $smokeDir "nested"))
        dryRun = $true
        confirm = $false
        confirmPhrase = ""
    } | ConvertTo-Json
    $cleanup = Invoke-RestMethod -Method POST -Uri "$Base/api/scans/$scanId/cleanup" -ContentType "application/json" -Body $cleanupBody
    if ($cleanup.results[0].status -ne "would_delete") { throw "cleanup dry-run failed" }
    Remove-Item -Recurse -Force $smokeDir -ErrorAction SilentlyContinue
    Write-Host "OK POST /api/scans/{id}/cleanup dry-run"

    $ui = Invoke-WebRequest "$Base/" -UseBasicParsing
    if ($ui.StatusCode -ne 200) { throw "UI status $($ui.StatusCode)" }
    Write-Host "OK UI /"

    Write-Host "smoke-api: all checks passed"
} finally {
    if ($server -and -not $server.HasExited) {
        Stop-Process -Id $server.Id -Force -ErrorAction SilentlyContinue
    }
}
