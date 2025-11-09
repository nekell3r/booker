# PowerShell script for setting up the booker project
# Alternative to 'make setup' for Windows PowerShell users

Write-Host "Starting full setup..." -ForegroundColor Cyan

# Function to retry a command
function Retry-Command {
    param(
        [scriptblock]$Command,
        [int]$MaxRetries = 3,
        [int]$DelaySeconds = 10
    )
    
    $retry = 0
    while ($retry -lt $MaxRetries) {
        Write-Host "Attempt $($retry + 1) of $MaxRetries..." -ForegroundColor Yellow
        try {
            & $Command
            if ($LASTEXITCODE -eq 0) {
                return $true
            }
        } catch {
            Write-Host "Error: $_" -ForegroundColor Red
        }
        
        $retry++
        if ($retry -lt $MaxRetries) {
            Write-Host "⚠️  Command failed, retrying in $DelaySeconds seconds..." -ForegroundColor Yellow
            Start-Sleep -Seconds $DelaySeconds
        }
    }
    
    Write-Host "❌ Failed after $MaxRetries attempts" -ForegroundColor Red
    return $false
}

# Step 1: Pull Docker images
Write-Host "`nStep 1: Pulling Docker images..." -ForegroundColor Cyan
$pullSuccess = Retry-Command -Command {
    docker-compose --profile infra-min --profile apps pull
}

if (-not $pullSuccess) {
    Write-Host "⚠️  Some images failed to pull. Continuing anyway..." -ForegroundColor Yellow
    Write-Host "💡 You can retry pulling images later with: docker-compose --profile infra-min --profile apps pull" -ForegroundColor Yellow
}

# Step 2: Start infrastructure and services
Write-Host "`nStep 2: Starting infrastructure and services..." -ForegroundColor Cyan
docker-compose --profile infra-min --profile apps up -d
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Failed to start services." -ForegroundColor Red
    Write-Host "This might be due to:" -ForegroundColor Yellow
    Write-Host "  - Network issues pulling images" -ForegroundColor Yellow
    Write-Host "  - Port conflicts" -ForegroundColor Yellow
    Write-Host "  - Docker daemon not running" -ForegroundColor Yellow
    exit 1
}

# Step 3: Wait for databases
Write-Host "`nStep 3: Waiting for databases to be ready..." -ForegroundColor Cyan
Start-Sleep -Seconds 15

# Step 4: Run migrations
Write-Host "`nStep 4: Running migrations..." -ForegroundColor Cyan
docker-compose --profile infra-min run --rm migrate
if ($LASTEXITCODE -ne 0) {
    Write-Host "Migration failed, retrying..." -ForegroundColor Yellow
    Start-Sleep -Seconds 5
    docker-compose --profile infra-min run --rm migrate
}

# Step 5: Seed data
Write-Host "`nStep 5: Seeding sample data..." -ForegroundColor Cyan
docker-compose --profile infra-min run --rm seed

# Success message
Write-Host "`n✅ Setup complete!" -ForegroundColor Green
Write-Host "📊 Admin Gateway: http://localhost:8080" -ForegroundColor Cyan
Write-Host "📈 Grafana: http://localhost:3000 (admin/admin)" -ForegroundColor Cyan
Write-Host "🔍 Jaeger: http://localhost:16686" -ForegroundColor Cyan
Write-Host "📉 Prometheus: http://localhost:9090" -ForegroundColor Cyan
Write-Host "📨 Kafka UI: http://localhost:8081" -ForegroundColor Cyan

