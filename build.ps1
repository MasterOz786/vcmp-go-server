# Build Safari plugin and deploy into plugins/.
param(
    [switch]$StartServer,
    [switch]$NoTest,
    [switch]$NoStop
)

$pluginRoot = Join-Path $PSScriptRoot "..\vcmp-go-plugin"
$buildScript = Join-Path $pluginRoot "build.ps1"

if (-not (Test-Path $buildScript)) {
    throw "vcmp-go-plugin not found at $pluginRoot"
}

$params = @{
    ServerRoot = $PSScriptRoot
    StartServer = $StartServer
}
if (-not $NoTest) { $params.Test = $true }
if (-not $NoStop) { $params.StopServer = $true }

& $buildScript @params
