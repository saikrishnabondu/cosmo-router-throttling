<#
Starts Docker Desktop (if not already running) and brings up Redis + NATS.
The router needs NATS (port 4222) to build the composed schema, and Redis is
used by the event providers. Run this FIRST.
#>
$dockerRunning = $null
try { docker ps *> $null; $dockerRunning = $LASTEXITCODE -eq 0 } catch { $dockerRunning = $false }

if (-not $dockerRunning) {
    Write-Host "Starting Docker Desktop... (this can take a minute)" -ForegroundColor Cyan
    Start-Process "C:\Program Files\Docker\Docker\Docker Desktop.exe"
    $waited = 0
    while ($waited -lt 180) {
        Start-Sleep -Seconds 5; $waited += 5
        docker ps *> $null
        if ($LASTEXITCODE -eq 0) { Write-Host "Docker is ready." -ForegroundColor Green; break }
        Write-Host "  waiting for Docker... ($waited s)" -ForegroundColor DarkGray
    }
}

Push-Location "C:\cosmo"
docker compose up -d redis nats
Pop-Location
Write-Host "Redis: redis://localhost:6379   |   NATS: nats://localhost:4222" -ForegroundColor Green
