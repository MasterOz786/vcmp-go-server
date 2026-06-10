# Build Safari plugin from sibling vcmp-go-plugin and deploy into plugins/.
param(
    [switch]$StartServer,
    [switch]$NoTest,
    [switch]$NoStop
)

$pluginRoot = Join-Path $PSScriptRoot "..\vcmp-go-plugin"
$buildScript = Join-Path $pluginRoot "build-safari.ps1"

if (-not (Test-Path $buildScript)) {
    throw "vcmp-go-plugin not found at $pluginRoot"
}

$params = @{
    StartServer = $StartServer
}
if ($NoTest) { $params.NoTest = $true }
if ($NoStop) { $params.NoStop = $true }

& $buildScript @params
