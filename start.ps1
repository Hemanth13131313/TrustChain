Write-Host "Starting TRUSTCHAIN Infrastructure..." -ForegroundColor Cyan
docker-compose up -d

Write-Host "Waiting 5 seconds for PostgreSQL to initialize..." -ForegroundColor Yellow
Start-Sleep -Seconds 5

Write-Host "Starting API Gateway in the background..." -ForegroundColor Cyan
Start-Process -NoNewWindow -FilePath "go" -ArgumentList "run ./services/api-gateway/cmd/server/main.go"

Write-Host "Starting Frontend Dashboard..." -ForegroundColor Cyan
Set-Location web/dashboard
npm run dev
