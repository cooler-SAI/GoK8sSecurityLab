# ============================================================================
# PROJECT CONFIGURATION
# ============================================================================
PROJECT_NAME := gok8ssecuritylab
GHCR_REGISTRY := ghcr.io
GITHUB_USERNAME := cooler-sai
IMAGE_REPO := go-secure-server
IMAGE_VERSION := latest
IMAGE_NAME := $(GHCR_REGISTRY)/$(GITHUB_USERNAME)/$(IMAGE_REPO):$(IMAGE_VERSION)

# Git metadata for traceability
GIT_COMMIT := $(shell git rev-parse --short HEAD)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
GIT_REPO_URL := https://github.com/$(GITHUB_USERNAME)/$(PROJECT_NAME)

CONTAINER_NAME := go-secure-server-test
K8S_NAMESPACE := go-security-lab
KUBE_CONTEXT := docker-desktop
FALCO_NAMESPACE := falco
TRIVY_BIN := trivy
NETPOL_FILE := k8s/network-policy.yaml

# Port configuration
HOST_PORT := 8081  # Changed from 8080 to avoid conflict with MiniTool ShadowMaker
CONTAINER_PORT := 8080

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
# GHCR AUTHENTICATION
# ============================================================================

# Note: Use PAT token for authentication
# Command template: echo <TOKEN> | docker login ghcr.io -u <GITHUB_NAME> --password-stdin

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
# DOCKER BUILD WITH GIT METADATA
# ============================================================================

# Build with Git metadata for traceability
build-with-metadata:
	@echo --- Building Docker image with Git metadata: $(IMAGE_NAME):$(GIT_COMMIT) ---
	docker build \
		--pull \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		--build-arg GIT_BRANCH=$(GIT_BRANCH) \
		--build-arg VERSION=latest \
		--label "org.opencontainers.image.revision=$(GIT_COMMIT)" \
		--label "org.opencontainers.image.source=$(GIT_REPO_URL)" \
		--label "org.opencontainers.image.created=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')" \
		-t $(IMAGE_NAME):$(GIT_COMMIT) \
		-t $(IMAGE_NAME):latest \
		.
	@echo [OK] Image successfully built: $(IMAGE_NAME):$(GIT_COMMIT) and :latest

# Standard build (backward compatibility)
build:
	@echo --- Building Docker image: $(IMAGE_NAME) ---
	docker build --pull -t $(IMAGE_NAME) .
	@echo [OK] Image successfully built: $(IMAGE_NAME)

# Build multiple tags including Git SHA
build-all: build build-with-metadata

# ============================================================================
# DOCKER PUSH WITH MULTIPLE TAGS
# ============================================================================

# Push with Git SHA tag for traceability
push-with-metadata: build-with-metadata
	@echo --- Pushing Docker images to GHCR ---
	@echo INFO: Ensure you are logged in to GHCR (make docker-login)
	docker push $(IMAGE_NAME):$(GIT_COMMIT)
	docker push $(IMAGE_NAME):latest
	@echo [OK] Images successfully pushed:
	@echo "  - $(IMAGE_NAME):$(GIT_COMMIT)"
	@echo "  - $(IMAGE_NAME):latest"

# Standard push (backward compatibility)
push: build
	@echo --- Pushing Docker image to GHCR: $(IMAGE_NAME) ---
	@echo INFO: Ensure you are logged in to GHCR (make docker-login)
	docker push $(IMAGE_NAME)
	@echo [OK] Image successfully pushed: $(IMAGE_NAME)

# ============================================================================
# DOCKER RUN & TEST
# ============================================================================

# Check port availability
check-port:
	@echo --- Checking port $(HOST_PORT) availability ---
	@netstat -ano | findstr :$(HOST_PORT) >nul && \
		(echo [ERROR] Port $(HOST_PORT) is occupied! && \
		 echo Try: make run-on-port PORT=8082 && exit 1) || \
		echo [OK] Port $(HOST_PORT) is available

# Run on default port (8081)
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

# Run on custom port
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

# Test endpoints
test: run
	@echo.
	@echo --- Testing endpoints ---
	@curl -f http://localhost:$(HOST_PORT)/ >nul 2>&1 && echo Health check passed || echo Health check failed
	@curl -s http://localhost:$(HOST_PORT)/api/halloween | findstr "spookily_secure" >nul && echo API test passed || echo API test failed
	@echo [OK] All tests completed

# Stop container
stop:
	@echo --- Stopping container: $(CONTAINER_NAME) ---
	@docker stop $(CONTAINER_NAME) 2>nul && echo Container stopped || echo Container was not running

# Clean Docker resources
clean: stop
	@echo --- Cleaning Docker resources ---
	@docker rmi $(IMAGE_NAME) 2>nul && echo Image removed || echo Image not found
	@docker system prune -f --volumes 2>nul && echo System cleanup completed || echo Cleanup completed
	@echo [OK] Docker cleanup completed

# View container logs
logs:
	@echo --- Viewing logs for $(CONTAINER_NAME) ---
	@docker logs -f $(CONTAINER_NAME)

# Inspect image metadata
inspect-metadata:
	@echo --- Image metadata for $(IMAGE_NAME):latest ---
	@docker inspect $(IMAGE_NAME):latest --format='{{range $k, $v := .Config.Labels}}{{printf "%s = %s\n" $k $v}}{{end}}'

# ============================================================================
# SECURITY SCANNING TARGETS
# ============================================================================

# Version check helper
check-trivy-version:
	@echo --- Trivy version check (run monthly) ---
	@echo [INFO] Run 'make trivy-update' to update Trivy
	@echo [INFO] Current version:
	@$(TRIVY_BIN) --version | findstr "Version"

# Simple Trivy scan (no build)
trivy-scan:
	@echo --- Simple Trivy scan of $(IMAGE_NAME) ---
	@echo [INFO] "If image doesn't exist, run 'make build' first"
	@$(TRIVY_BIN) image --skip-version-check --severity HIGH,CRITICAL $(IMAGE_NAME)
	@echo [OK] Trivy scan completed

# Trivy scan with parameter
trivy-scan-custom:
ifdef IMAGE
	@echo --- Trivy scan of custom image: $(IMAGE) ---
	@$(TRIVY_BIN) image --skip-version-check --severity HIGH,CRITICAL $(IMAGE)
else
	@echo [ERROR] Please specify IMAGE parameter
	@echo Example: make trivy-scan-custom IMAGE=alpine:latest
endif

# Standard vulnerability scan
scan: build
	@echo --- Scanning Docker image for HIGH/CRITICAL vulnerabilities with Trivy ---
	$(TRIVY_BIN) image --skip-version-check --severity HIGH,CRITICAL $(IMAGE_NAME)
	@echo [OK] Trivy scan completed

# Comprehensive vulnerability assessment
vulnerability-scan: build
	@echo --- Comprehensive vulnerability assessment ---
	@echo [INFO] Scanning $(IMAGE_NAME) for vulnerabilities...
	$(TRIVY_BIN) image --skip-version-check --format table --scanners vuln,misconfig --severity MEDIUM,HIGH,CRITICAL $(IMAGE_NAME)
	@echo [OK] Vulnerability scan completed

# Scan with Docker Scout
scan-docker: build
	@echo --- Scanning with Docker Scout ---
	docker scout quickview $(IMAGE_NAME)

# Strict scan (fails on CRITICAL/HIGH)
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

# Scan Kubernetes manifests
scan-manifests:
	@echo --- Scanning Kubernetes YAML for misconfigurations ---
	@echo [INFO] Scanning directory: k8s/
	$(TRIVY_BIN) config --skip-version-check --severity HIGH,CRITICAL k8s/
	@echo [OK] Kubernetes manifest scan completed.

# Scan running Kubernetes cluster
scan-cluster: k8s-check
	@echo --- Scanning the entire Kubernetes cluster ---
	@echo [INFO] This requires cluster-level permissions.
	$(TRIVY_BIN) k8s \
		--skip-version-check \
		--report summary \
		--namespace $(K8S_NAMESPACE) \
		cluster
	@echo [OK] Cluster scan completed.

# Comprehensive security pipeline
security-pipeline: check-trivy-version scan-strict scan-manifests
	@echo --- SECURITY PIPELINE PASSED ---
	@echo [SUCCESS] All security checks passed.
	@echo "  ✓ Container image scan (strict)"
	@echo "  ✓ Kubernetes manifest scan"
	@echo.

# ============================================================================
# FALCO RUNTIME SECURITY
# ============================================================================

# Clean Falco installation
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

# View Falco logs
falco-logs:
	@echo --- Falco security events ---
	@kubectl logs -l app=falco -n $(FALCO_NAMESPACE) -c falco --tail=50

# Test Falco detection
falco-test:
	@echo --- Testing Falco detection ---
	@echo [TEST] Generating test security event (Reading /etc/shadow)...
	@kubectl run falco-test --rm -i --tty --image=alpine --restart=Never --namespace $(K8S_NAMESPACE) -- sh -c "cat /etc/shadow"

# Clean Falco
falco-clean:
	@echo --- Removing Falco ---
	@helm uninstall falco -n $(FALCO_NAMESPACE) --wait 2>nul || echo Falco release not found
	@kubectl delete namespace $(FALCO_NAMESPACE) --ignore-not-found=true --timeout=60s
	@echo [OK] Falco removed

# ============================================================================
# KUBERNETES DEPLOYMENT
# ============================================================================

# Check Kubernetes access
k8s-check:
	@echo --- Checking Kubernetes access ---
	@kubectl config use-context $(KUBE_CONTEXT)
	@kubectl cluster-info
	@echo [OK] Kubernetes ready!

# Deploy to Kubernetes
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

# Access via port-forward
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

# Deploy and access
k8s-run: k8s-deploy k8s-access

# Clean Kubernetes resources
k8s-clean:
	@echo --- Cleaning Kubernetes application resources ---
	@kubectl delete -f k8s/ -n $(K8S_NAMESPACE) --ignore-not-found=true
	@kubectl delete namespace $(K8S_NAMESPACE) --ignore-not-found=true --timeout=60s
	@echo [OK] K8s application cleanup completed

# ============================================================================
# DEBUG & MONITORING
# ============================================================================

# Check pod status
check-pods:
	@echo --- Current pod status ---
	@kubectl get pods -n $(K8S_NAMESPACE) -o wide

# View recent logs
check-logs:
	@echo --- Recent logs ---
	@kubectl logs -n $(K8S_NAMESPACE) -l app=go-secure-server --tail=20

# Check deployment status
check-status:
	@echo --- Deployment status ---
	@kubectl get deployment -n $(K8S_NAMESPACE)
	@echo.
	@kubectl get svc -n $(K8S_NAMESPACE)
	@echo.
	@kubectl get pods -n $(K8S_NAMESPACE)

# Test service connectivity (port 80)
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

# Validate deployment
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
# NETWORK POLICY & SECURITY TESTING
# ============================================================================

# Apply network policy
netpol-apply: k8s-check
	@echo --- Applying Network Policy ---
	@kubectl apply -f $(NETPOL_FILE) -n $(K8S_NAMESPACE)
	@echo [OK] NetworkPolicy applied

# Remove network policy
netpol-clean: k8s-check
	@echo --- Removing Network Policy ---
	@kubectl delete -f $(NETPOL_FILE) -n $(K8S_NAMESPACE) --ignore-not-found=true
	@echo [OK] NetworkPolicy removed

# Test network policy
netpol-test: k8s-check
	@echo --- Testing Network Policy ---
	@echo 1. Checking service exists...
	@kubectl get svc go-secure-service -n $(K8S_NAMESPACE)
	@echo.
	@echo 2. Test allowed port (80):
	@kubectl run netpol-test-allowed --rm -i --image=curlimages/curl --restart=Never --namespace $(K8S_NAMESPACE) -- curl -v http://go-secure-service:80/health
	@echo.
	@echo [OK] Network policy test completed

# Policy test: simulate malicious pod
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

# Quick security deployment
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

# Quick deployment with push
security-quick-push: push
	@echo --- Deploying to Kubernetes (using GHCR image) ---
	@kubectl create namespace $(K8S_NAMESPACE) --dry-run=client -o yaml | kubectl apply -f -
	@kubectl apply -f k8s/ -n $(K8S_NAMESPACE)
	@echo.
	@echo [SUCCESS] DEPLOYMENT WITH PUSH COMPLETED!
	@echo.
	@echo [NEXT] Access with: make k8s-access

# Full security stack
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

# Full security stack with Git metadata
security-full-with-metadata: push-with-metadata security-pipeline
	@echo --- Deploying full security stack with Git metadata ---
	@kubectl create namespace $(K8S_NAMESPACE) --dry-run=client -o yaml | kubectl apply -f -
	@kubectl apply -f k8s/ -n $(K8S_NAMESPACE)
	@kubectl apply -f $(NETPOL_FILE) -n $(K8S_NAMESPACE)
	@echo.
	@echo [SUCCESS] FULL SECURITY STACK WITH GIT METADATA DEPLOYED!
	@echo [OK] Docker image built with Git metadata
	@echo [OK] Image pushed with SHA tag: $(GIT_COMMIT)
	@echo [OK] Security scan completed
	@echo [OK] K8s application deployed
	@echo [OK] Network Policy applied
	@echo.
	@echo [NEXT] Access with: make k8s-access

# Complete security environment cleanup
security-clean: k8s-clean falco-clean clean
	@echo --- Complete security environment cleanup ---

# ============================================================================
# CI/CD SIMULATION
# ============================================================================

# Simulate CI/CD pipeline locally
ci-pipeline-local:
	@echo --- Simulating CI/CD Pipeline Locally ---
	@echo [STAGE 1] Prerequisites check
	make check-prerequisites
	@echo.
	@echo [STAGE 2] Security scanning
	make security-pipeline
	@echo.
	@echo [STAGE 3] Build with metadata
	make build-with-metadata
	@echo.
	@echo [STAGE 4] Image metadata inspection
	make inspect-metadata
	@echo.
	@echo [STAGE 5] Deploy to Kubernetes
	make k8s-deploy
	@echo.
	@echo [SUCCESS] CI/CD pipeline simulation completed!

# GitHub Actions compatible commands
ci-build-and-push:
	@echo --- CI/CD: Build and push with Git metadata ---
	make build-with-metadata
	docker push $(IMAGE_NAME):$(GIT_COMMIT)
	docker push $(IMAGE_NAME):latest
	@echo [SUCCESS] CI/CD build and push completed

# ============================================================================
# INFORMATION & HELP
# ============================================================================

help:
	@echo.
	@echo "================================================================================"
	@echo "GO K8S SECURITY LAB - WINDOWS MAKEFILE (GHCR WITH GIT METADATA)"
	@echo "================================================================================"
	@echo.
	@echo "GIT METADATA & TRACEABILITY:"
	@echo "   make build-with-metadata    - Build with Git commit SHA as tag"
	@echo "   make push-with-metadata     - Push both :SHA and :latest tags"
	@echo "   make inspect-metadata       - View image labels and metadata"
	@echo "   make ci-build-and-push      - CI/CD compatible build and push"
	@echo.
	@echo "QUICK START:"
	@echo "   make security-quick         - Build and deploy locally (port 8081)"
	@echo "   make run                    - Run locally on port 8081"
	@echo "   make run-on-port PORT=9090  - Run on custom port"
	@echo "   make k8s-access             - Access K8s app via port-forward"
	@echo.
	@echo "DOCKER & GHCR:"
	@echo "   make docker-login           - Show GHCR login instructions"
	@echo "   make build                  - Standard build"
	@echo "   make push                   - Standard push"
	@echo "   make build-all              - Build both standard and metadata versions"
	@echo.
	@echo "SECURITY SCANNING:"
	@echo "   make security-pipeline      - Run full security scan pipeline"
	@echo "   make scan-strict            - Strict scan (fails on vulnerabilities)"
	@echo "   make scan-manifests         - Scan Kubernetes YAML files"
	@echo "   make trivy-scan             - Simple vulnerability scan"
	@echo.
	@echo "KUBERNETES:"
	@echo "   make k8s-deploy             - Deploy to Kubernetes"
	@echo "   make k8s-clean              - Remove K8s resources"
	@echo "   make validate-deployment    - Validate K8s deployment"
	@echo "   make policy-test            - Test network policies"
	@echo.
	@echo "FALCO RUNTIME SECURITY:"
	@echo "   make falco-install-clean    - Install Falco runtime security"
	@echo "   make falco-test             - Test Falco detection"
	@echo.
	@echo "CI/CD SIMULATION:"
	@echo "   make ci-pipeline-local      - Simulate full CI/CD pipeline locally"
	@echo "   make security-full-with-metadata - Full stack with Git traceability"
	@echo.
	@echo "CLEANUP:"
	@echo "   make security-clean         - Complete cleanup"
	@echo "   make clean                  - Clean Docker resources"
	@echo.
	@echo "================================================================================"
	@echo "Git Commit: $(GIT_COMMIT) | Branch: $(GIT_BRANCH)"
	@echo "Image: $(IMAGE_NAME)"
	@echo "GHCR: https://github.com/$(GITHUB_USERNAME)/$(PROJECT_NAME)/pkgs/container/$(IMAGE_REPO)"
	@echo "================================================================================"

.PHONY: all build build-with-metadata build-all push push-with-metadata run run-on-port test stop clean logs \
        check-prerequisites docker-login inspect-metadata \
        trivy-scan trivy-scan-custom scan vulnerability-scan scan-docker scan-strict scan-manifests scan-cluster security-pipeline \
        falco-install-clean falco-logs falco-test falco-clean \
        k8s-check k8s-deploy k8s-access k8s-run k8s-clean \
        check-pods check-logs check-status test-service test-service-simple validate-deployment health-check \
        netpol-apply netpol-clean netpol-test policy-test check-port \
        security-quick security-quick-push security-full security-full-with-metadata security-clean \
        ci-pipeline-local ci-build-and-push help \
        trivy-update trivy-version trivy-check check-trivy-version