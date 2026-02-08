# gen_certs.ps1 - Generate TLS certificates for Kubernetes Webhook
# INSTRUCTION: If local OpenSSL is not found, use Docker with SAN support

Write-Host "🎯 Sentinel Webhook - Certificate Generation" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan

# Create directory for certificates
if (!(Test-Path "certs")) {
    New-Item -ItemType Directory -Path "certs" -Force
    Write-Host "📁 Created certs/ directory" -ForegroundColor Green
} else {
    Write-Host "📁 Directory certs/ already exists" -ForegroundColor Yellow
}

Write-Host "`n🔐 Generating certificates using Docker (with SAN for K8s)..." -ForegroundColor Yellow

# Create temporary configuration file for OpenSSL with Subject Alternative Names (SAN)
$extFile = @"
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = sentinel-service
DNS.2 = sentinel-service.default
DNS.3 = sentinel-service.default.svc
DNS.4 = sentinel-service.default.svc.cluster.local
DNS.5 = sentinel-service.sentinel-system
DNS.6 = sentinel-service.sentinel-system.svc
DNS.7 = sentinel-service.sentinel-system.svc.cluster.local
DNS.8 = localhost
IP.1 = 127.0.0.1
"@

$extFile | Out-File -FilePath "certs/extfile.cnf" -Encoding ascii
Write-Host "📄 Created configuration file with SAN extensions" -ForegroundColor Green

try {
    # Run Alpine container, install openssl and generate certificate with extensions
    Write-Host "🐳 Starting Docker container for certificate generation..." -ForegroundColor Yellow

    # Create a script file for Docker to execute
    $dockerScript = @"
apk add --no-cache openssl > /dev/null 2>&1

echo '🔑 Generating private key (2048 bit)...'
openssl genrsa -out certs/tls.key 2048

echo '📝 Creating CSR (Certificate Signing Request)...'
openssl req -new -key certs/tls.key -out certs/tls.csr -subj '/CN=sentinel-service.default.svc/C=US/ST=CA/L=SanFrancisco/O=Sentinel/OU=Security'

echo '🏷️  Creating certificate with SAN extensions...'
openssl x509 -req -days 365 -in certs/tls.csr -signkey certs/tls.key -out certs/tls.crt -extfile certs/extfile.cnf

echo '🧹 Cleaning temporary files...'
rm certs/tls.csr certs/extfile.cnf

echo '✅ Certificates successfully created!'
echo ''
echo '📊 Certificate information:'
openssl x509 -in certs/tls.crt -text -noout | grep -E '(Subject:|DNS:|IP Address:|Not Before|Not After)' | head -10
"@

    # Save Docker script to file
    $dockerScript | Out-File -FilePath "docker_script.sh" -Encoding ascii

    # Execute Docker command
    docker run --rm -v "${PWD}:/work" -w /work alpine sh docker_script.sh

    # Clean up script file
    Remove-Item docker_script.sh -ErrorAction SilentlyContinue

    if ($LASTEXITCODE -eq 0) {
        Write-Host "`n✅ Success! Certificates with SAN created in certs/ folder" -ForegroundColor Green
        Write-Host "📋 Created files:" -ForegroundColor White
        Write-Host "   • certs/tls.key - private key" -ForegroundColor Gray
        Write-Host "   • certs/tls.crt - certificate with SAN" -ForegroundColor Gray

        # Show certificate information
        Write-Host "`n🔍 Certificate information:" -ForegroundColor Cyan
        if (Test-Path "certs/tls.crt") {
            $certInfo = & openssl x509 -in certs/tls.crt -text -noout 2>$null
            if ($LASTEXITCODE -eq 0) {
                $certInfoLines = $certInfo -split "`n"
                $certInfoLines | Select-String -Pattern "(Subject:|DNS:|IP Address:|Not Before|Not After|X509v3 Subject Alternative Name)" | ForEach-Object {
                    Write-Host "   $_" -ForegroundColor Gray
                }
            }
        }

        Write-Host "`n📋 Next steps:" -ForegroundColor Cyan

        # Option 1: Create Secret via kubectl
        Write-Host "1. Create Kubernetes Secret:" -ForegroundColor Yellow
        Write-Host "   kubectl create secret tls sentinel-certs --cert=certs/tls.crt --key=certs/tls.key --namespace=sentinel-system" -ForegroundColor Gray

        # Option 2: Create YAML file for Secret
        Write-Host "`n2. Or create secret.yaml file:" -ForegroundColor Yellow
        Write-Host "   kubectl create secret tls sentinel-certs --cert=certs/tls.crt --key=certs/tls.key --namespace=sentinel-system --dry-run=client -o yaml > secret.yaml" -ForegroundColor Gray

        # Show Base64 for caBundle
        Write-Host "`n🔑 Base64 for caBundle (paste into ValidatingWebhookConfiguration):" -ForegroundColor White
        $base64Cert = [Convert]::ToBase64String([System.IO.File]::ReadAllBytes("certs/tls.crt"))
        Write-Host "caBundle: ""$base64Cert""" -ForegroundColor Gray

        Write-Host "`n⚠️  Important:" -ForegroundColor Red
        Write-Host "   - Certificate valid for 365 days" -ForegroundColor Gray
        Write-Host "   - Includes SAN for all typical Kubernetes service DNS names" -ForegroundColor Gray
        Write-Host "   - Use sentinel-system namespace for deployment" -ForegroundColor Gray

        # Verify certificate
        Write-Host "`n🧪 Certificate verification:" -ForegroundColor Cyan
        if (Test-Path "certs/tls.crt") {
            try {
                $verifyResult = & openssl verify -CAfile certs/tls.crt certs/tls.crt 2>&1
                if ($verifyResult -match "OK") {
                    Write-Host "   ✅ Certificate is valid" -ForegroundColor Green
                } else {
                    Write-Host "   ⚠️  Warning: $verifyResult" -ForegroundColor Yellow
                }
            } catch {
                Write-Host "   ⚠️  Could not verify certificate" -ForegroundColor Yellow
            }
        }
    } else {
        Write-Host "`n❌ Error: Docker command failed" -ForegroundColor Red
        Write-Host "   Make sure Docker Desktop is running" -ForegroundColor Gray
        Write-Host "   Check permissions for certs/ folder" -ForegroundColor Gray
    }
}
catch {
    Write-Host "`n❌ Error: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "   Check Docker installation" -ForegroundColor Gray
    Write-Host "   Check permissions" -ForegroundColor Gray
}
finally {
    # Cleanup
    Remove-Item docker_script.sh -ErrorAction SilentlyContinue
}

Write-Host "`n🏁 Certificate generation completed" -ForegroundColor Cyan
Write-Host "====================================" -ForegroundColor Cyan

# Quick commands for convenience
Write-Host "`n⚡ Quick commands to copy:" -ForegroundColor Magenta
Write-Host "cd certs" -ForegroundColor Gray
Write-Host "kubectl create ns sentinel-system" -ForegroundColor Gray
Write-Host "kubectl create secret tls sentinel-certs --cert=tls.crt --key=tls.key --namespace=sentinel-system" -ForegroundColor Gray