# ============================================================================
# PROJECT CONFIGURATION (WINDOWS COMPATIBLE)
# ============================================================================
PROJECT_NAME := gok8ssecuritylab
GHCR_REGISTRY := ghcr.io
GITHUB_USERNAME := cooler-sai
IMAGE_REPO := go-secure-server
IMAGE_VERSION := latest
IMAGE_NAME := $(GHCR_REGISTRY)/$(GITHUB_USERNAME)/$(IMAGE_REPO):$(IMAGE_VERSION)

GIT_REPO_URL := https://github.com/$(GITHUB_USERNAME)/$(PROJECT_NAME)

CONTAINER_NAME := go-secure-server-test
K8S_NAMESPACE := go-security-lab
KUBE_CONTEXT := docker-desktop
FALCO_NAMESPACE := falco
TRIVY_BIN := trivy
NETPOL_FILE := k8s/network-policy.yaml

# Port configuration
HOST_PORT := 8081
CONTAINER_PORT := 8080

# ============================================================================
# TRIVY MANAGEMENT (SIMPLIFIED FOR WINDOWS)
# ============================================================================

trivy-update:
	@echo --- Updating Trivy to latest version ---
	@scoop update trivy
	@echo [OK] Trivy updated successfully

trivy-version:
	@echo --- Trivy Version Information ---
	@$(TRIVY_BIN) --version
	@echo.

trivy-check:
	@echo --- Checking Trivy ---
	@if exist $(TRIVY_BIN) (
		echo [OK] Trivy found at: $(TRIVY_BIN)
		@$(TRIVY_BIN) --version
	) else (
		echo [WARNING] Trivy not found. Install with: scoop install trivy
	)

# ============================================================================
# GHCR AUTHENTICATION
# ============================================================================

docker-login:
	@echo --- Authenticating with GHCR ---
	@echo NOTE: You need a GitHub Personal Access Token (PAT) with 'write:packages' scope
	@echo Create token at: https://github.com/settings/tokens
	@echo.
	@echo If you have a token, run manually:
	@echo   docker login ghcr.io -u $(GITHUB_USERNAME)
	@echo.
	@echo For automation (Windows):
	@echo   echo YOUR_PAT ^| docker login ghcr.io -u $(GITHUB_USERNAME) --password-stdin

# ============================================================================
# PREREQUISITE CHECKS
# ============================================================================

check-prerequisites:
	@echo --- Checking required tools ---
	@docker version >nul 2>&1 && echo [OK] Docker found || (echo [ERROR] Docker not found && exit 1)
	@kubectl version --client >nul 2>&1 && echo [OK] kubectl found || (echo [ERROR] kubectl not found && exit 1)
	@helm version >nul 2>&1 && echo [OK] Helm found || (echo [ERROR] Helm found && exit 1)
	@if exist $(TRIVY_BIN) (echo [OK] Trivy found) else (echo [WARNING] Trivy not found)
	@echo [OK] Prerequisites check completed

# ============================================================================
# GIT HELPER FUNCTIONS (WINDOWS COMPATIBLE)
# ============================================================================

# Get Git commit SHA
get-git-commit:
	@powershell -Command "git rev-parse --short HEAD"

# Get Git branch
get-git-branch:
	@powershell -Command "git rev-parse --abbrev-ref HEAD"

# Show Git info
show-git-info:
	@echo --- Git Information ---
	@powershell -Command "$$commit = git rev-parse --short HEAD; $$branch = git rev-parse --abbrev-ref HEAD; Write-Host 'Commit: $$commit'; Write-Host 'Branch: $$branch'"
	@echo.

# ============================================================================
# DOCKER BUILD (SIMPLE - ALWAYS WORKS)
# ============================================================================

# Standard build (backward compatibility)
build:
	@echo --- Building Docker image: $(IMAGE_NAME) ---
	docker build --pull -t $(IMAGE_NAME) .
	@echo [OK] Image successfully built: $(IMAGE_NAME)

# Test build - без сложностей
test-build:
	@echo --- Test build ---
	docker build -t ghcr.io/cooler-sai/go-secure-server:test .
	@echo [OK] Test build completed

# ============================================================================
# DOCKER BUILD WITH GIT METADATA (POWERSHELL VERSION - WORKS)
# ============================================================================

# PowerShell-based build (reliable on Windows)
build-with-git:
	@echo --- Building with Git metadata (PowerShell) ---
	@powershell -Command " \
		$$commit = git rev-parse --short HEAD; \
		$$branch = git rev-parse --abbrev-ref HEAD; \
		$$date = Get-Date -Format 'yyyy-MM-ddTHH:mm:ssZ'; \
		$$image = 'ghcr.io/cooler-sai/go-secure-server'; \
		Write-Host 'Git: $$commit ($$branch)' -ForegroundColor Cyan; \
		Write-Host 'Building: $$image:$$commit' -ForegroundColor Yellow; \
		docker build \
			--build-arg GIT_COMMIT=$$commit \
			--build-arg GIT_BRANCH=$$branch \
			--build-arg BUILD_DATE=$$date \
			--build-arg VERSION=latest \
			-t $$image`:$$commit \
			-t $$image`:latest \
			.; \
		if ($$LASTEXITCODE -eq 0) { \
			Write-Host '✅ Build successful!' -ForegroundColor Green; \
			docker images $$image`:*; \
		} else { \
			Write-Host '❌ Build failed!' -ForegroundColor Red; \
		} \
	"

# Simple build with commit tag (упрощенная версия)
build-simple:
	@echo --- Simple build with Git commit tag ---
	@powershell -Command " \
		$$commit = git rev-parse --short HEAD; \
		$$image = 'ghcr.io/cooler-sai/go-secure-server'; \
		Write-Host 'Commit: $$commit'; \
		Write-Host 'Building: $$image:$$commit'; \
		docker build -t $$image`:$$commit -t $$image`:latest .; \
	"

# ============================================================================
# DOCKER PUSH
# ============================================================================

# Standard push
push: build
	@echo --- Pushing Docker image to GHCR: $(IMAGE_NAME) ---
	@echo INFO: Ensure you are logged in to GHCR (make docker-login)
	docker push $(IMAGE_NAME)
	@echo [OK] Image successfully pushed: $(IMAGE_NAME)

# Push with Git metadata
push-with-git: build-with-git
	@echo --- Pushing Docker images with Git metadata ---
	@echo INFO: Ensure you are logged in to GHCR (make docker-login)
	@powershell -Command " \
		$$commit = git rev-parse --short HEAD; \
		$$image = 'ghcr.io/cooler-sai/go-secure-server'; \
		Write-Host 'Pushing: $$image:$$commit' -ForegroundColor Yellow; \
		docker push $$image`:$$commit; \
		docker push $$image`:latest; \
		Write-Host '✅ Push completed!' -ForegroundColor Green; \
	"

# ============================================================================
# DOCKER RUN & TEST
# ============================================================================

check-port:
	@echo --- Checking port $(HOST_PORT) availability ---
	@netstat -ano | findstr :$(HOST_PORT) >nul && \
		(echo [ERROR] Port $(HOST_PORT) is occupied! && \
		 echo Try: make run-on-port PORT=8082 && exit 1) || \
		echo [OK] Port $(HOST_PORT) is available

run: stop check-port build
	@echo --- Running container at http://localhost:$(HOST_PORT) ---
	docker run -d --rm --name $(CONTAINER_NAME) -p $(HOST_PORT):$(CONTAINER_PORT) $(IMAGE_NAME)
	@echo [OK] Container $(CONTAINER_NAME) successfully started
	@echo.
	@echo   #######################################################
	@echo   #         SECURE HALLOWEEN SERVER STARTED            #
	@echo   #######################################################
	@echo.
	@echo Waiting for server to start...
	@timeout /t 3 /nobreak >nul
	@echo Server ready! Opening page...
	@start http://localhost:$(HOST_PORT)/halloween 2>nul || echo Manually open http://localhost:$(HOST_PORT)/halloween in browser

run-on-port:
ifdef PORT
	@echo --- Running container at http://localhost:$(PORT) ---
	docker run -d --rm --name $(CONTAINER_NAME) -p $(PORT):$(CONTAINER_PORT) $(IMAGE_NAME)
	@echo [OK] Container $(CONTAINER_NAME) successfully started
	@echo Open: http://localhost:$(PORT)/halloween
	@timeout /t 3 /nobreak >nul
	@start http://localhost:$(PORT)/halloween 2>nul || echo Open manually in browser
else
	@echo [ERROR] Please specify PORT parameter
	@echo Example: make run-on-port PORT=9090
endif

test: run
	@echo.
	@echo --- Testing endpoints ---
	@curl -f http://localhost:$(HOST_PORT)/ >nul 2>&1 && echo Health check passed || echo Health check failed
	@curl -s http://localhost:$(HOST_PORT)/api/halloween | findstr "spookily_secure" >nul && echo API test passed || echo API test failed
	@echo [OK] All tests completed

stop:
	@echo --- Stopping container: $(CONTAINER_NAME) ---
	@docker stop $(CONTAINER_NAME) 2>nul && echo Container stopped || echo Container was not running

clean: stop
	@echo --- Cleaning Docker resources ---
	@docker rmi $(IMAGE_NAME) 2>nul && echo Image removed || echo Image not found
	@docker system prune -f --volumes 2>nul && echo System cleanup completed || echo Cleanup completed
	@echo [OK] Docker cleanup completed

logs:
	@echo --- Viewing logs for $(CONTAINER_NAME) ---
	@docker logs -f $(CONTAINER_NAME)

# Inspect image metadata
inspect:
	@echo --- Inspecting image: $(IMAGE_NAME) ---
	@docker inspect $(IMAGE_NAME) --format='{{range $$k, $$v := .Config.Labels}}{{printf "%-40s = %s\n" $$k $$v}}{{end}}'

# ============================================================================
# SECURITY SCANNING
# ============================================================================

check-trivy-version:
	@echo --- Trivy version check (run monthly) ---
	@echo [INFO] Run 'make trivy-update' to update Trivy
	@echo [INFO] Current version:
	@$(TRIVY_BIN) --version | findstr "Version"

trivy-scan:
	@echo --- Simple Trivy scan of $(IMAGE_NAME) ---
	@echo [INFO] "If image doesn't exist, run 'make build' first"
	@$(TRIVY_BIN) image --skip-version-check --severity HIGH,CRITICAL $(IMAGE_NAME)
	@echo [OK] Trivy scan completed

trivy-scan-custom:
ifdef IMAGE
	@echo --- Trivy scan of custom image: $(IMAGE) ---
	@$(TRIVY_BIN) image --skip-version-check --severity HIGH,CRITICAL $(IMAGE)
else
	@echo [ERROR] Please specify IMAGE parameter
	@echo Example: make trivy-scan-custom IMAGE=alpine:latest
endif

scan: build
	@echo --- Scanning Docker image for HIGH/CRITICAL vulnerabilities with Trivy ---
	$(TRIVY_BIN) image --skip-version-check --severity HIGH,CRITICAL $(IMAGE_NAME)
	@echo [OK] Trivy scan completed

vulnerability-scan: build
	@echo --- Comprehensive vulnerability assessment ---
	@echo [INFO] Scanning $(IMAGE_NAME) for vulnerabilities...
	$(TRIVY_BIN) image --skip-version-check --format table --scanners vuln,misconfig --severity MEDIUM,HIGH,CRITICAL $(IMAGE_NAME)
	@echo [OK] Vulnerability scan completed

scan-docker: build
	@echo --- Scanning with Docker Scout ---
	docker scout quickview $(IMAGE_NAME)

scan-strict: build
	@echo --- Strict vulnerability scan (Fails on CRITICAL/HIGH) ---
	@echo [INFO] Exits with code 1 if CRITICAL/HIGH vulnerabilities are found.
	$(TRIVY_BIN) image \
		--skip-version-check \
		--severity HIGH,CRITICAL \
		--exit-code 1 \
		--ignore-unfixed \
		$(IMAGE_NAME) && \
	echo "[PASS] No CRITICAL/HIGH vulnerabilities." || \
	(echo "[FAIL] Build blocked: CRITICAL/HIGH vulnerabilities found!" && exit 1)

scan-manifests:
	@echo --- Scanning Kubernetes YAML for misconfigurations ---
	@echo [INFO] Scanning directory: k8s/
	$(TRIVY_BIN) config --skip-version-check --severity HIGH,CRITICAL k8s/
	@echo [OK] Kubernetes manifest scan completed.

scan-cluster: k8s-check
	@echo --- Scanning the entire Kubernetes cluster ---
	@echo [INFO] This requires cluster-level permissions.
	$(TRIVY_BIN) k8s \
		--skip-version-check \
		--report summary \
		--namespace $(K8S_NAMESPACE) \
		cluster
	@echo [OK] Cluster scan completed.

security-pipeline: check-trivy-version scan-strict scan-manifests
	@echo --- SECURITY PIPELINE PASSED ---
	@echo [SUCCESS] All security checks passed.
	@echo "  ✓ Container image scan (strict)"
	@echo "  ✓ Kubernetes manifest scan"
	@echo.

# ============================================================================
# FALCO RUNTIME SECURITY
# ============================================================================

falco-install-clean: k8s-check
	@echo --- Complete Falco cleanup and fresh installation ---
	@echo [INFO] Removing previous Falco resources...
	@helm uninstall falco -n $(FALCO_NAMESPACE) --wait 2>nul || echo Helm release not found
	@kubectl delete namespace $(FALCO_NAMESPACE) --ignore-not-found=true --timeout=60s
	@echo [INFO] Installing Falco...
	@helm repo add falcosecurity https://falcosecurity.github.io/charts
	@helm repo update
	@helm upgrade --install falco falcosecurity/falco --namespace $(FALCO_NAMESPACE) --create-namespace --set falcosidekick.enabled=false
	@echo [OK] Falco installation completed!

falco-logs:
	@echo --- Falco security events ---
	@kubectl logs -l app=falco -n $(FALCO_NAMESPACE) -c falco --tail=50

falco-test:
	@echo --- Testing Falco detection ---
	@echo [TEST] Generating test security event (Reading /etc/shadow)...
	@kubectl run falco-test --rm -i --tty --image=alpine --restart=Never --namespace $(K8S_NAMESPACE) -- sh -c "cat /etc/shadow"

falco-clean:
	@echo --- Removing Falco ---
	@helm uninstall falco -n $(FALCO_NAMESPACE) --wait 2>nul || echo Falco release not found
	@kubectl delete namespace $(FALCO_NAMESPACE) --ignore-not-found=true --timeout=60s
	@echo [OK] Falco removed

# ============================================================================
# KUBERNETES DEPLOYMENT
# ============================================================================

k8s-check:
	@echo --- Checking Kubernetes access ---
	@kubectl config use-context $(KUBE_CONTEXT)
	@kubectl cluster-info
	@echo [OK] Kubernetes ready!

k8s-deploy: k8s-check
	@echo --- Deploying Go Secure Server to Kubernetes (using image from GHCR) ---
	@kubectl create namespace $(K8S_NAMESPACE) --dry-run=client -o yaml | kubectl apply -f -
	@kubectl apply -f k8s/ -n $(K8S_NAMESPACE)
	@echo --- Waiting for pods to start (30 seconds) ---
	@echo INFO: Checking pod status...
	@echo.
	@timeout /t 10 /nobreak >nul
	@kubectl get pods -n $(K8S_NAMESPACE)
	@echo.
	@echo Waiting additional 20 seconds for readiness...
	@timeout /t 20 /nobreak >nul
	@kubectl get pods -n $(K8S_NAMESPACE)
	@echo.
	@echo [OK] Application deployment submitted!
	@echo To check logs: kubectl logs -n $(K8S_NAMESPACE) -l app=go-secure-server

k8s-access:
	@echo --- Starting Port Forward ---
	@echo [INFO] Your Go Secure Server is now in K8s!
	@echo.
	@echo [ENDPOINTS] Available at http://localhost:8080/ :
	@echo    - Health Check:    http://localhost:8080/
	@echo    - Halloween Page:  http://localhost:8080/halloween
	@echo    - Halloween API:   http://localhost:8080/api/halloween
	@echo.
	@echo NOTE: Service uses port 80 internally, port-forward maps 8080->80
	@echo.
	@echo Press Ctrl+C to stop Port-Forwarding
	@echo.
	@kubectl port-forward svc/go-secure-service 8080:80 -n $(K8S_NAMESPACE)

k8s-run: k8s-deploy k8s-access

k8s-clean:
	@echo --- Cleaning Kubernetes application resources ---
	@kubectl delete -f k8s/ -n $(K8S_NAMESPACE) --ignore-not-found=true
	@kubectl delete namespace $(K8S_NAMESPACE) --ignore-not-found=true --timeout=60s
	@echo [OK] K8s application cleanup completed

# ============================================================================
# DEBUG & MONITORING
# ============================================================================

check-pods:
	@echo --- Current pod status ---
	@kubectl get pods -n $(K8S_NAMESPACE) -o wide

check-logs:
	@echo --- Recent logs ---
	@kubectl logs -n $(K8S_NAMESPACE) -l app=go-secure-server --tail=20

check-status:
	@echo --- Deployment status ---
	@kubectl get deployment -n $(K8S_NAMESPACE)
	@echo.
	@kubectl get svc -n $(K8S_NAMESPACE)
	@echo.
	@kubectl get pods -n $(K8S_NAMESPACE)

test-service:
	@echo --- Testing service connectivity (port 80) ---
	@echo Service go-secure-service uses port 80
	@echo Testing: http://go-secure-service:80/
	@echo.
	@kubectl run test-service-curl --rm -i --image=curlimages/curl --restart=Never --namespace $(K8S_NAMESPACE) -- curl -s http://go-secure-service:80/

test-service-simple:
	@echo --- Testing service (default port 80) ---
	@kubectl run curl-simple --rm -i --image=curlimages/curl --restart=Never --namespace $(K8S_NAMESPACE) -- curl -s http://go-secure-service/

validate-deployment:
	@echo --- Validating deployment ---
	@echo 1. Checking pods:
	@kubectl get pods -n $(K8S_NAMESPACE) -l app=go-secure-server
	@echo.
	@echo 2. Checking service:
	@kubectl get svc go-secure-service -n $(K8S_NAMESPACE)
	@echo.
	@echo 3. Testing service connectivity:
	@kubectl run validate-curl --rm -i --image=curlimages/curl --restart=Never --namespace $(K8S_NAMESPACE) -- sh -c "echo 'Testing service...' && curl -s http://go-secure-service/"
	@echo.
	@echo [OK] Deployment validated!

health-check:
	@echo --- Quick health check ---
	@echo Pods:
	@kubectl get pods -n $(K8S_NAMESPACE) -l app=go-secure-server -o custom-columns=NAME:.metadata.name,STATUS:.status.phase,IP:.status.podIP
	@echo.
	@echo To test service: make test-service
	@echo To access locally: make k8s-access

# ============================================================================
# NETWORK POLICY & SECURITY TESTING
# ============================================================================

netpol-apply: k8s-check
	@echo --- Applying Network Policy ---
	@kubectl apply -f $(NETPOL_FILE) -n $(K8S_NAMESPACE)
	@echo [OK] NetworkPolicy applied

netpol-clean: k8s-check
	@echo --- Removing Network Policy ---
	@kubectl delete -f $(NETPOL_FILE) -n $(K8S_NAMESPACE) --ignore-not-found=true
	@echo [OK] NetworkPolicy removed

netpol-test: k8s-check
	@echo --- Testing Network Policy ---
	@echo 1. Checking service exists...
	@kubectl get svc go-secure-service -n $(K8S_NAMESPACE)
	@echo.
	@echo 2. Test allowed port (80):
	@kubectl run netpol-test-allowed --rm -i --image=curlimages/curl --restart=Never --namespace $(K8S_NAMESPACE) -- curl -v http://go-secure-service:80/health
	@echo.
	@echo [OK] Network policy test completed

policy-test: k8s-check
	@echo --- Testing Security Policy with Simulated Attack ---
	@echo "[TEST 1] Deploying 'test-attacker' pod in namespace: $(K8S_NAMESPACE)"
	@kubectl run test-attacker \
		--image=alpine \
		--restart=Never \
		--namespace=$(K8S_NAMESPACE) \
		--command -- sh -c "sleep 3600" 2>nul || echo "Pod might already exist."
	@timeout /t 2 /nobreak >nul
	@echo.
	@echo "[TEST 2] Trying to access secure server from attacker pod..."
	@kubectl exec test-attacker -n $(K8S_NAMESPACE) -- \
		sh -c "echo 'Testing network policy...'; wget -qO- --timeout=3 http://go-secure-service:80/health || echo 'ACCESS BLOCKED (Expected with NetworkPolicy)'"
	@echo.
	@echo "[TEST 3] Cleanup..."
	@kubectl delete pod test-attacker -n $(K8S_NAMESPACE) --ignore-not-found=true
	@echo [OK] Security policy test completed.

# ============================================================================
# COMPREHENSIVE SECURITY PIPELINES
# ============================================================================

security-quick: build
	@echo --- Deploying to Kubernetes ---
	@kubectl create namespace $(K8S_NAMESPACE) --dry-run=client -o yaml | kubectl apply -f -
	@kubectl apply -f k8s/ -n $(K8S_NAMESPACE)
	@echo.
	@echo [SUCCESS] QUICK SECURITY DEPLOYMENT COMPLETED!
	@echo.
	@echo [STATUS] Checking pods:
	@timeout /t 5 /nobreak >nul
	@kubectl get pods -n $(K8S_NAMESPACE)
	@echo.
	@echo [NEXT] Access with: make k8s-access
	@echo [TEST] Validate deployment: make validate-deployment

security-quick-push: push
	@echo --- Deploying to Kubernetes (using GHCR image) ---
	@kubectl create namespace $(K8S_NAMESPACE) --dry-run=client -o yaml | kubectl apply -f -
	@kubectl apply -f k8s/ -n $(K8S_NAMESPACE)
	@echo.
	@echo [SUCCESS] DEPLOYMENT WITH PUSH COMPLETED!
	@echo.
	@echo [NEXT] Access with: make k8s-access

security-full: push security-pipeline
	@echo --- Deploying full security stack ---
	@kubectl create namespace $(K8S_NAMESPACE) --dry-run=client -o yaml | kubectl apply -f -
	@kubectl apply -f k8s/ -n $(K8S_NAMESPACE)
	@kubectl apply -f $(NETPOL_FILE) -n $(K8S_NAMESPACE)
	@echo.
	@echo [SUCCESS] FULL SECURITY STACK DEPLOYED!
	@echo [OK] Docker image built, pushed, and scanned
	@echo [OK] K8s application deployed
	@echo [OK] Network Policy applied
	@echo.
	@echo [NEXT] Access with: make k8s-access

security-full-with-git: push-with-git security-pipeline
	@echo --- Deploying full security stack with Git metadata ---
	@kubectl create namespace $(K8S_NAMESPACE) --dry-run=client -o yaml | kubectl apply -f -
	@kubectl apply -f k8s/ -n $(K8S_NAMESPACE)
	@kubectl apply -f $(NETPOL_FILE) -n $(K8S_NAMESPACE)
	@echo.
	@echo [SUCCESS] FULL SECURITY STACK WITH GIT METADATA DEPLOYED!
	@echo [OK] Docker image built with Git metadata
	@echo [OK] Image pushed with SHA tag
	@echo [OK] Security scan completed
	@echo [OK] K8s application deployed
	@echo [OK] Network Policy applied
	@echo.
	@echo [NEXT] Access with: make k8s-access

security-clean: k8s-clean falco-clean clean
	@echo --- Complete security environment cleanup ---

# ============================================================================
# CI/CD SIMULATION
# ============================================================================

ci-pipeline-local:
	@echo --- Simulating CI/CD Pipeline Locally ---
	@echo [STAGE 1] Prerequisites check
	make check-prerequisites
	@echo.
	@echo [STAGE 2] Security scanning
	make security-pipeline
	@echo.
	@echo [STAGE 3] Build with metadata
	make build-with-git
	@echo.
	@echo [STAGE 4] Inspect image
	make inspect
	@echo.
	@echo [STAGE 5] Deploy to Kubernetes
	make k8s-deploy
	@echo.
	@echo [SUCCESS] CI/CD pipeline simulation completed!

ci-build-and-push: build-with-git
	@echo --- CI/CD: Build and push with Git metadata ---
	@powershell -Command " \
		$$commit = git rev-parse --short HEAD; \
		$$image = 'ghcr.io/cooler-sai/go-secure-server'; \
		docker push $$image`:$$commit; \
		docker push $$image`:latest; \
		Write-Host '✅ CI/CD pipeline completed!' -ForegroundColor Green; \
	"

# ============================================================================
# INFORMATION & HELP
# ============================================================================

help:
	@echo.
	@echo "================================================================================"
	@echo "GO K8S SECURITY LAB - WINDOWS MAKEFILE (GHCR WITH GIT METADATA)"
	@echo "================================================================================"
	@echo.
	@echo "🚀 QUICK START:"
	@echo "   make build                      - Build image"
	@echo "   make build-with-git             - Build with Git metadata (RECOMMENDED)"
	@echo "   make test-build                 - Simple test build"
	@echo "   make push-with-git              - Push with Git metadata"
	@echo.
	@echo "🔧 GIT METADATA (WINDOWS):"
	@echo "   make get-git-commit             - Get current Git commit SHA"
	@echo "   make get-git-branch             - Get current Git branch"
	@echo "   make show-git-info              - Show Git information"
	@echo "   make inspect                    - View image labels and metadata"
	@echo.
	@echo "🐳 DOCKER & GHCR:"
	@echo "   make docker-login               - Show GHCR login instructions"
	@echo "   make build                      - Standard build"
	@echo "   make push                       - Standard push"
	@echo "   make run                        - Run locally on port 8081"
	@echo "   make run-on-port PORT=9090      - Run on custom port"
	@echo.
	@echo "🔐 SECURITY SCANNING:"
	@echo "   make security-pipeline          - Run full security scan pipeline"
	@echo "   make scan-strict                - Strict scan (fails on vulnerabilities)"
	@echo "   make scan-manifests             - Scan Kubernetes YAML files"
	@echo.
	@echo "☸️ KUBERNETES:"
	@echo "   make k8s-deploy                 - Deploy to Kubernetes"
	@echo "   make k8s-access                 - Access via port-forward"
	@echo "   make k8s-clean                  - Remove K8s resources"
	@echo.
	@echo "🔄 CI/CD SIMULATION:"
	@echo "   make ci-pipeline-local          - Simulate full CI/CD pipeline"
	@echo "   make ci-build-and-push          - Build and push for CI/CD"
	@echo.
	@echo "🧹 CLEANUP:"
	@echo "   make security-clean             - Complete cleanup"
	@echo "   make clean                      - Clean Docker resources"
	@echo.
	@echo "================================================================================"
	@echo "Image: $(IMAGE_NAME)"
	@echo "GitHub: $(GIT_REPO_URL)"
	@echo "GHCR: https://github.com/$(GITHUB_USERNAME)/$(PROJECT_NAME)/pkgs/container/$(IMAGE_REPO)"
	@echo "================================================================================"

# ============================================================================
# PHONY TARGETS
# ============================================================================
.PHONY: all build test-build build-with-git build-simple \
        push push-with-git \
        run run-on-port test stop clean logs inspect \
        check-prerequisites docker-login \
        get-git-commit get-git-branch show-git-info \
        trivy-scan trivy-scan-custom scan vulnerability-scan scan-docker scan-strict scan-manifests scan-cluster security-pipeline \
        falco-install-clean falco-logs falco-test falco-clean \
        k8s-check k8s-deploy k8s-access k8s-run k8s-clean \
        check-pods check-logs check-status test-service test-service-simple validate-deployment health-check \
        netpol-apply netpol-clean netpol-test policy-test check-port \
        security-quick security-quick-push security-full security-full-with-git security-clean \
        ci-pipeline-local ci-build-and-push help \
        trivy-update trivy-version trivy-check check-trivy-version