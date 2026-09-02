<#
Starts the two backend services (subgraphs) this demo uses:
  employees -> http://localhost:4101/graphql
  products  -> http://localhost:4004/graphql

These are shared demo backend services (from C:\cosmo\demo). If they are already
running (e.g. from another demo), you will see a harmless "address already in use"
bind error in one window — just ignore it.

Requires Docker infra (Redis + NATS) to be up so the schema builds. If it isn't,
start it first with the infra script from the cosmo demo, or:
  cd C:\cosmo ; docker compose up -d redis nats
#>
$env:OTEL_HTTP_ENDPOINT = "localhost:0"
Start-Process powershell -ArgumentList "-NoExit","-Command", `
  "cd 'C:\cosmo\demo\cmd\employees'; `$env:PORT=4101; `$env:OTEL_HTTP_ENDPOINT='localhost:0'; go run main.go"
Start-Process powershell -ArgumentList "-NoExit","-Command", `
  "cd 'C:\cosmo\demo\cmd\products'; `$env:OTEL_HTTP_ENDPOINT='localhost:0'; go run main.go"
Write-Host "employees -> http://localhost:4101/graphql"
Write-Host "products  -> http://localhost:4004/graphql"
Write-Host "(first run downloads Go modules; wait for each window to print the playground URL)"
