# scripts/show-security-status.ps1
param(
    [string]$Action = "deploy"
)

function Show-SecurityDeploy {
    Write-Host "`n🎉 FULL SECURITY STACK DEPLOYED!" -ForegroundColor Green
    Write-Host "✅ Docker image built and scanned" -ForegroundColor Green
    Write-Host "✅ Falco runtime security installed (clean install)" -ForegroundColor Green
    Write-Host "✅ K8s application deployed with security context" -ForegroundColor Green
    Write-Host "✅ Network Policy applied (Zero Trust)" -ForegroundColor Green
    Write-Host ""
    Write-Host "🛡️  SECURITY MONITORING ACTIVE!" -ForegroundColor Cyan
    Write-Host "🔍 Run: make falco-logs    - View security events" -ForegroundColor Gray
    Write-Host "🧪 Run: make netpol-test   - Test Network Policy enforcement" -ForegroundColor Gray
    Write-Host "🚀 Run: make k8s-access    - Access your application" -ForegroundColor Gray
}

function Show-SecurityTest {
    Write-Host "`n🧪 SECURITY TESTING MODE" -ForegroundColor Yellow
    Write-Host "⚡ Running security scans..." -ForegroundColor Yellow
    Write-Host "📊 Results will be shown below" -ForegroundColor Gray
}

function Show-SecurityClean {
    Write-Host "`n🧹 CLEANING SECURITY STACK" -ForegroundColor Magenta
    Write-Host "🗑️  Removing all security components..." -ForegroundColor Magenta
}

switch ($Action) {
    "deploy" { Show-SecurityDeploy }
    "test"   { Show-SecurityTest }
    "clean"  { Show-SecurityClean }
    default  { Show-SecurityDeploy }
}