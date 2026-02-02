# INSTRUCTION: If local OpenSSL is not found, use Docker
if (!(Test-Path "certs")) {
    New-Item -ItemType Directory -Path "certs" -Force
    Write-Host "📁 Created certs/ directory" -ForegroundColor Cyan
}

Write-Host "🔐 Attempting to generate certificates using Docker container..." -ForegroundColor Yellow

try {
    # Run alpine container, mount current folder and generate keys
    docker run --rm -v "${PWD}:/work" -w /work alpine sh -c "
        echo '📦 Installing OpenSSL in Alpine container...' && \
        apk add --no-cache openssl > /dev/null 2>&1 && \
        echo '🔑 Generating private key...' && \
        openssl genrsa -out certs/tls.key 2048 && \
        echo '📄 Generating self-signed certificate...' && \
        openssl req -new -x509 -sha256 -key certs/tls.key -out certs/tls.crt -days 365 -subj '/CN=sentinel-service.default.svc' && \
        echo '✅ Certificate generation complete!'"

    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Success! Certificates created in certs/ folder (using Docker)" -ForegroundColor Green
        Write-Host "Files created:" -ForegroundColor White
        Write-Host "  • certs/tls.key (private key)" -ForegroundColor Gray
        Write-Host "  • certs/tls.crt (certificate)" -ForegroundColor Gray
        Write-Host "" -ForegroundColor White
        Write-Host "📋 Next step:" -ForegroundColor Cyan
        Write-Host "kubectl create secret tls sentinel-certs --cert=certs/tls.crt --key=certs/tls.key" -ForegroundColor Yellow
    } else {
        Write-Host "❌ Error: Certificate generation failed." -ForegroundColor Red
        Write-Host "Please ensure Docker Desktop is running." -ForegroundColor Red
    }
} catch {
    Write-Host "❌ Error: $($_.Exception.Message)" -ForegroundColor Red
}