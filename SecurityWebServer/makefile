# ============================================================================
# PROJECT CONFIGURATION
# ============================================================================
PROJECT_NAME := gok8ssecuritylab
GHCR_REGISTRY := ghcr.io
GITHUB_USERNAME := cooler-sai
IMAGE_REPO := go-secure-server
IMAGE_VERSION := latest
IMAGE_NAME := $(GHCR_REGISTRY)/$(GITHUB_USERNAME)/$(IMAGE_REPO):$(IMAGE_VERSION)

CONTAINER_NAME := go-secure-server-test
K8S_NAMESPACE := go-security-lab
KUBE_CONTEXT := docker-desktop
FALCO_NAMESPACE := falco
TRIVY_BIN := trivy
NETPOL_FILE := k8s/network-policy.yaml

# ============================================================================
# TRIVY MANAGEMENT (SIMPLIFIED FOR WINDOWS)
# ============================================================================

# Simple Trivy update target
trivy-update:
	@echo --- Updating Trivy to latest version ---
	@scoop update trivy
	@echo [OK] Trivy updated successfully

# Show Trivy version
trivy-version:
	@echo --- Trivy Version Information ---
	@$(TRIVY_BIN) --version
	@echo.

# Check if Trivy needs update (simplified)
trivy-check:
	@echo --- Checking Trivy ---
	@if exist $(TRIVY_BIN) (
		echo [OK] Trivy found at: $(TRIVY_BIN)
		@$(TRIVY_BIN) --version
	) else (
		echo [WARNING] Trivy not found. Install with: scoop install trivy
	)

# ============================================================================
# GHCR Login
# ============================================================================

# echo <TOKEN> | docker login ghcr.io -u <GITHUB_NAME> --password-stdin

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
# DOCKER AUTHENTICATION
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
	@echo   echo YOUR_PAT | docker login ghcr.io -u $(GITHUB_USERNAME) --password-stdin

# ============================================================================
# MAIN TARGETS (Docker Management)
# ============================================================================

all: security-quick

build:
	@echo --- Building Docker image: $(IMAGE_NAME) ---
	docker build --pull -t $(IMAGE_NAME) .
	@echo [OK] Image successfully built: $(IMAGE_NAME)

push: build
	@echo --- Pushing Docker image: $(IMAGE_NAME) to GHCR ---
	@echo INFO: Ensure you are logged in to GHCR (make docker-login)
	docker push $(IMAGE_NAME)
	@echo [OK] Image successfully pushed: $(IMAGE_NAME)

run: stop build
	@echo --- Running container at http://localhost:8080 ---
	docker run -d --rm --name $(CONTAINER_NAME) -p 8080:8080 $(IMAGE_NAME)
	@echo [OK] Container $(CONTAINER_NAME) successfully started
	@echo.
	@echo   #######################################################
	@echo   #         SECURE HALLOWEEN SERVER STARTED            #
	@echo   #######################################################
	@echo.
	@echo Waiting for server to start...
	@timeout /t 3 /nobreak >nul
	@echo Server ready! Opening page...
	@start http://localhost:8080/halloween 2>nul || echo Manually open http://localhost:8080/halloween in browser

test: run
	@echo.
	@echo --- Testing endpoints ---
	@curl -f http://localhost:8080/ >nul 2>&1 && echo Health check passed || echo Health check failed
	@curl -s http://localhost:8080/api/halloween | findstr "spookily_secure" >nul && echo API test passed || echo API test failed
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

# ============================================================================
# SECURITY SCANNING TARGETS (ENHANCED)
# ============================================================================

# Version check helper
check-trivy-version:
	@echo --- Trivy version check (run monthly) ---
	@echo [INFO] Run 'make trivy-update' to update Trivy
	@echo [INFO] Current version:
	@$(TRIVY_BIN) --version | findstr "Version"

# SIMPLE TRIVY SCAN (JUST SCAN, NO BUILD)
trivy-scan:
	@echo --- Simple Trivy scan of $(IMAGE_NAME) ---
	@echo [INFO] If image doesn''t exist, run 'make build' first
	@$(TRIVY_BIN) image --skip-version-check --severity HIGH,CRITICAL $(IMAGE_NAME)
	@echo [OK] Trivy scan completed

# TRIVY SCAN WITH PARAMETER
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

# STRICT IMAGE SCAN (CI/CD FAILURE)
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

# SCAN KUBERNETES MANIFESTS (IaC SECURITY)
scan-manifests:
	@echo --- Scanning Kubernetes YAML for misconfigurations ---
	@echo [INFO] Scanning directory: k8s/
	$(TRIVY_BIN) config --skip-version-check --severity HIGH,CRITICAL k8s/
	@echo [OK] Kubernetes manifest scan completed.

# SCAN RUNNING KUBERNETES CLUSTER
scan-cluster: k8s-check
	@echo --- Scanning the entire Kubernetes cluster ---
	@echo [INFO] This requires cluster-level permissions.
	$(TRIVY_BIN) k8s \
		--skip-version-check \
		--report summary \
		--namespace $(K8S_NAMESPACE) \
		cluster
	@echo [OK] Cluster scan completed.

# COMPREHENSIVE SCAN PIPELINE (ALL CHECKS)
security-pipeline: check-trivy-version scan-strict scan-manifests
	@echo --- SECURITY PIPELINE PASSED ---
	@echo [SUCCESS] All security checks passed.
	@echo "  ✓ Container image scan (strict)"
	@echo "  ✓ Kubernetes manifest scan"
	@echo.

# ============================================================================
# FALCO RUNTIME SECURITY TARGETS
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
# KUBERNETES TARGETS (Deployment & Access)
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
# DEBUG & MONITORING TARGETS
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

# Test service on port 80 (service uses port 80, not 8080)
test-service:
	@echo --- Testing service connectivity (port 80) ---
	@echo Service go-secure-service uses port 80
	@echo Testing: http://go-secure-service:80/
	@echo.
	@kubectl run test-service-curl --rm -i --image=curlimages/curl --restart=Never --namespace $(K8S_NAMESPACE) -- curl -s http://go-secure-service:80/

# Test service without port (defaults to 80)
test-service-simple:
	@echo --- Testing service (default port 80) ---
	@kubectl run curl-simple --rm -i --image=curlimages/curl --restart=Never --namespace $(K8S_NAMESPACE) -- curl -s http://go-secure-service/

# Validate everything works
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

# Quick health check
health-check:
	@echo --- Quick health check ---
	@echo Pods:
	@kubectl get pods -n $(K8S_NAMESPACE) -l app=go-secure-server -o custom-columns=NAME:.metadata.name,STATUS:.status.phase,IP:.status.podIP
	@echo.
	@echo To test service: make test-service
	@echo To access locally: make k8s-access

# ============================================================================
# NETWORK POLICY & SECURITY TESTING TARGETS
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

# POLICY TEST: SIMULATE MALICIOUS POD
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

security-clean: k8s-clean falco-clean clean
	@echo --- Complete security environment cleanup ---

# ============================================================================
# INFORMATION & HELP
# ============================================================================

help:
	@echo.
	@echo "================================================================================"
	@echo "GO K8S SECURITY LAB - WINDOWS MAKEFILE (GHCR)"
	@echo "================================================================================"
	@echo.
	@echo "QUICK START:"
	@echo "   make security-quick        - Build and deploy locally"
	@echo "   make security-quick-push   - Build, push to GHCR, deploy"
	@echo "   make security-full         - Full stack (push, scan, deploy, netpol)"
	@echo "   make k8s-access            - Access app via port-forward (8080->80)"
	@echo.
	@echo "DOCKER & REGISTRY (GHCR):"
	@echo "   make build                 - Build image locally"
	@echo "   make docker-login          - Show GHCR login instructions"
	@echo "   make push                  - Push image to GHCR (requires login)"
	@echo "   make run                   - Run locally"
	@echo.
	@echo "SECURITY SCANNING (ENHANCED):"
	@echo "   make trivy-scan            - Simple scan of $(IMAGE_NAME) (no build)"
	@echo "   make trivy-scan-custom     - Scan custom image (IMAGE=alpine:latest)"
	@echo "   make scan                  - Vulnerability scan (High/Critical)"
	@echo "   make scan-strict           - Strict scan (fails on CRITICAL/HIGH)"
	@echo "   make scan-manifests        - Scan k8s/ YAML files for misconfigs"
	@echo "   make scan-cluster          - Scan the running K8s cluster"
	@echo "   make security-pipeline     - Run full security scan pipeline"
	@echo "   make vulnerability-scan    - Comprehensive vulnerability assessment"
	@echo "   make scan-docker           - Scan with Docker Scout"
	@echo.
	@echo "TRIVY MANAGEMENT:"
	@echo "   make trivy-update          - Update Trivy to latest version via Scoop"
	@echo "   make trivy-version         - Show current Trivy version"
	@echo "   make trivy-check           - Check if Trivy is installed"
	@echo.
	@echo "RUNTIME SECURITY (FALCO):"
	@echo "   make falco-install-clean   - Install Falco"
	@echo "   make falco-logs            - View Falco events"
	@echo "   make falco-test            - Run test to trigger Falco"
	@echo.
	@echo "KUBERNETES DEPLOYMENT:"
	@echo "   make k8s-deploy            - Deploy to Kubernetes"
	@echo "   make k8s-clean             - Remove K8s application resources"
	@echo.
	@echo "DEBUG & TESTING:"
	@echo "   make check-pods            - Check pod status"
	@echo "   make check-logs            - View recent logs"
	@echo "   make check-status          - Check deployment status"
	@echo "   make test-service          - Test service connectivity (port 80)"
	@echo "   make test-service-simple   - Test service (default port)"
	@echo "   make validate-deployment   - Full deployment validation"
	@echo "   make health-check          - Quick health check"
	@echo "   make policy-test           - Simulate attack to test NetworkPolicy"
	@echo "   make netpol-test           - Test network policy"
	@echo.
	@echo "================================================================================"

.PHONY: all build push run test stop clean logs \
        check-prerequisites docker-login \
        trivy-scan trivy-scan-custom \
        scan vulnerability-scan scan-docker scan-strict scan-manifests scan-cluster security-pipeline \
        falco-install-clean falco-logs falco-test falco-clean \
        k8s-check k8s-deploy k8s-access k8s-run k8s-clean \
        check-pods check-logs check-status test-service test-service-simple validate-deployment health-check \
        netpol-apply netpol-clean netpol-test policy-test \
        security-quick security-quick-push security-full security-clean help \
        trivy-update trivy-version trivy-check check-trivy-version