# Runs after /reload server shuts down the process: rebuild plugin, deploy, start server64.
param(
    [int]$WaitSeconds = 2
)

$ErrorActionPreference = "Stop"
$Root = Split-Path $PSScriptRoot -Parent
$Build = Join-Path $Root "build.ps1"

Start-Sleep -Seconds $WaitSeconds

if (-not (Test-Path $Build)) {
    Write-Error "build.ps1 not found at $Build"
    exit 1
}

Write-Host "[hotreload] rebuilding plugin and starting server..."
& $Build -StartServer -NoTest
