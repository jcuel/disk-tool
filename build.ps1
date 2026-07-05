# Build disk-tool on Windows
$ErrorActionPreference = "Stop"
Set-Location $PSScriptRoot

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
