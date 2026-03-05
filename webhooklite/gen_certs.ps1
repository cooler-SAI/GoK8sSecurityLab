# gen_certs.ps1 - Generate TLS certificates for webhooklite

Write-Host "🎯 webhooklite - Certificate Generation" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan

# Create directory for certificates
if (!(Test-Path "certificates")) {
    New-Item -ItemType Directory -Path "certificates" -Force
    Write-Host "📁 Created certificates/ directory" -ForegroundColor Green
} else {
    Write-Host "📁 Directory certificates/ already exists" -ForegroundColor Yellow
}

# Create openssl.cnf
$opensslCnf = @"
[req]
distinguished_name = req_distinguished_name
x509_extensions = v3_req
prompt = no

[req_distinguished_name]
CN = webhooklite-service.default.svc

[v3_req]
keyUsage = keyEncipherment, dataEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[alt_names]
DNS.1 = webhooklite-service
DNS.2 = webhooklite-service.default
DNS.3 = webhooklite-service.default.svc
DNS.4 = webhooklite-service.default.svc.cluster.local
DNS.5 = localhost
IP.1 = 127.0.0.1
"@

$opensslCnf | Out-File -FilePath "certificates/openssl.cnf" -Encoding ascii
Write-Host "📄 Created openssl.cnf with SAN extensions" -ForegroundColor Green

Write-Host "`n🔐 Generating certificates using Docker..." -ForegroundColor Yellow

# Step 1: Generate private key
Write-Host "`n🔑 Step 1: Generating private key..." -ForegroundColor Yellow
docker run --rm -v "${PWD}:/work" -w /work alpine sh -c "apk add --no-cache openssl && openssl genrsa -out certificates/tls.key 2048"

if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Failed to generate private key" -ForegroundColor Red
    exit 1
}
Write-Host "✅ Private key created: certificates/tls.key" -ForegroundColor Green

# Step 2: Generate certificate
Write-Host "`n📝 Step 2: Generating certificate with SAN..." -ForegroundColor Yellow
docker run --rm -v "${PWD}:/work" -w /work alpine sh -c "apk add --no-cache openssl && openssl req -new -x509 -days 365 -key certificates/tls.key -out certificates/tls.crt -config certificates/openssl.cnf"

if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Failed to generate certificate" -ForegroundColor Red
    exit 1
}
Write-Host "✅ Certificate created: certificates/tls.crt" -ForegroundColor Green

# Step 3: Show certificate info
Write-Host "`n🔍 Certificate information:" -ForegroundColor Cyan
docker run --rm -v "${PWD}:/work" -w /work alpine sh -c "apk add --no-cache openssl > /dev/null && openssl x509 -in certificates/tls.crt -text -noout | grep -E 'Subject:|DNS:|IP Address:|Not Before|Not After'" | Select-String -Pattern "(Subject:|DNS:|IP Address:|Not Before|Not After|X509v3 Subject Alternative Name)" | ForEach-Object {
    Write-Host "   $_" -ForegroundColor Gray
}

# Step 4: Show Base64 for caBundle
Write-Host "`n🔑 Base64 for caBundle (paste into ValidatingWebhookConfiguration):" -ForegroundColor White
$base64Cert = [Convert]::ToBase64String([System.IO.File]::ReadAllBytes("certificates/tls.crt"))
Write-Host "caBundle: ""$base64Cert""" -ForegroundColor Gray

# Step 5: Create Kubernetes secret
Write-Host "`n📦 Creating Kubernetes secret..." -ForegroundColor Yellow
kubectl delete secret webhooklite-certs -n default --ignore-not-found
kubectl create secret tls webhooklite-certs --cert=certificates/tls.crt --key=certificates/tls.key -n default

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Secret created: webhooklite-certs" -ForegroundColor Green
} else {
    Write-Host "⚠️  Could not create secret. Make sure kubectl is configured." -ForegroundColor Yellow
}

Write-Host "`n⚠️  Important:" -ForegroundColor Red
Write-Host "   - Certificate valid for 365 days" -ForegroundColor Gray
Write-Host "   - Includes SAN for webhooklite service" -ForegroundColor Gray
Write-Host "   - Use default namespace for deployment" -ForegroundColor Gray

Write-Host "`n🏁 Certificate generation completed" -ForegroundColor Cyan
Write-Host "====================================" -ForegroundColor Cyan

# Quick commands
Write-Host "`n⚡ Quick commands:" -ForegroundColor Magenta
Write-Host "cd certificates" -ForegroundColor Gray
Write-Host "kubectl get secret webhooklite-certs -n default" -ForegroundColor Gray