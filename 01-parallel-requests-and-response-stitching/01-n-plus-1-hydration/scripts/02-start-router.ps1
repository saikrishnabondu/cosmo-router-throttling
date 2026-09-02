<#
Starts the STOCK Cosmo router (no throttling, no custom modules) against this
demo's clean config. Frees port 3003 first in case another router is bound there.
#>
$RouterConfigDir = Split-Path -Parent $MyInvocation.MyCommand.Path | Join-Path -ChildPath "..\router-config" | Resolve-Path
Push-Location $RouterConfigDir

# Free port 3003 (e.g. if the throttling router is running there)
$conn = Get-NetTCPConnection -LocalPort 3003 -State Listen -ErrorAction SilentlyContinue | Select-Object -First 1
if ($conn) {
    Write-Host "Freeing port 3003 (stopping PID $($conn.OwningProcess))..." -ForegroundColor DarkGray
    Stop-Process -Id $conn.OwningProcess -Force -ErrorAction SilentlyContinue
    Start-Sleep -Seconds 1
}

Write-Host "Starting CLEAN Cosmo router (no throttling) on http://localhost:3003 ..." -ForegroundColor Cyan
& "C:\cosmo\router\router.exe" -config router.config.yaml
Pop-Location
