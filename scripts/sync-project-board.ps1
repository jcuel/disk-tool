param(
    [string]$Owner = "jcuel",
    [int]$ProjectNumber = 3,
    [string]$Repo = "jcuel/disk-tool"
)

$ErrorActionPreference = "Stop"
$ProjectId = "PVT_kwHOA0E_1c4Bcf_W"
$FieldId = "PVTSSF_lAHOA0E_1c4Bcf_WzhXIUKk"
$StatusDone = "0d582c4b"
$StatusReady = "7eda7d29"

$closed = @(gh issue list --repo $Repo --state closed --json number --jq '.[].number')
$data = gh project item-list $ProjectNumber --owner $Owner --format json --limit 100 | ConvertFrom-Json

foreach ($item in $data.items) {
    $num = $item.content.number
    if (-not $num) { continue }
    $issueState = gh issue view $num --repo $Repo --json state --jq '.state'
    $status = $item.status

    if ($issueState -eq "CLOSED" -and $status -ne "Done") {
        Write-Host "Issue #$num : $status -> Done"
        gh project item-edit --id $item.id --project-id $ProjectId --field-id $FieldId `
            --single-select-option-id $StatusDone | Out-Null
        continue
    }

    if ($issueState -eq "OPEN" -and $status -eq "Backlog") {
        $body = gh issue view $num --repo $Repo --json body --jq '.body'
        $change = $null
        if ($body -match '/propose\s+([a-z0-9-]+)') { $change = $Matches[1] }
        $fallback = @{ 4 = "duplicate-detection"; 5 = "age-based-cleanup" }
        if ($fallback.ContainsKey($num)) { $change = $fallback[$num] }
        if ($change -and (Test-Path "openspec/changes/$change/proposal.md")) {
            Write-Host "Issue #$num : Backlog -> Ready (proposal exists)"
            gh project item-edit --id $item.id --project-id $ProjectId --field-id $FieldId `
                --single-select-option-id $StatusReady | Out-Null
        }
    }
}

Write-Host "Project board sync complete."
