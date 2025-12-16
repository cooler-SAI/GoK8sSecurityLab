# ============================================================================
# PROJECT CONFIGURATION
# ============================================================================
PROJECT_NAME := gok8ssecuritylab
# !!! IMPORTANT: Replace 'YOUR_GITHUB_USERNAME' with your actual GitHub username !!!
GHCR_REGISTRY := ghcr.io
GITHUB_USERNAME := cooler-sai
IMAGE_REPO := go-secure-server
IMAGE_VERSION := latest
# Full image name for GHCR: ghcr.io/<username>/go-secure-server:latest
IMAGE_NAME := $(GHCR_REGISTRY)/$(GITHUB_USERNAME)/$(IMAGE_REPO):$(IMAGE_VERSION)

CONTAINER_NAME := go-secure-server-test
K8S_NAMESPACE := go-security-lab
KUBE_CONTEXT := docker-desktop
FALCO_NAMESPACE := falco
TRIVY_BIN := "C:\Users\coole\scoop\shims\trivy.exe"
NETPOL_FILE := k8s/network-policy.yaml

# ============================================================================
# PREREQUISITE CHECKS
# ============================================================================

check-prerequisites:
	@echo --- Checking required tools ---
	@docker version >nul 2>&1 && echo [OK] Docker found || (echo [ERROR] Docker not found && exit 1)
	@kubectl version --client >nul 2>&1 && echo [OK] kubectl found || (echo [ERROR] kubectl not found && exit 1)
	@helm version >nul 2>&1 && echo [OK] Helm found || (echo [ERROR] Helm not found && exit 1)
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
# SECURITY SCANNING TARGETS
# ============================================================================

scan: build
	@echo --- Scanning Docker image for HIGH/CRITICAL vulnerabilities with Trivy ---
	$(TRIVY_BIN) image --severity HIGH,CRITICAL $(IMAGE_NAME)
	@echo [OK] Trivy scan completed

vulnerability-scan: build
	@echo --- Comprehensive vulnerability assessment ---
	@echo [INFO] Scanning $(IMAGE_NAME) for vulnerabilities...
	$(TRIVY_BIN) image --format table --scanners vuln,misconfig --severity MEDIUM,HIGH,CRITICAL $(IMAGE_NAME)
	@echo [OK] Vulnerability scan completed

scan-docker: build
	@echo --- Scanning with Docker Scout ---
	docker scout quickview $(IMAGE_NAME)

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
	@echo Press Ctrl+C to stop Port-Forwarding
	@echo.
	@kubectl port-forward svc/go-secure-service 8080:8080 -n $(K8S_NAMESPACE)

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

# Get first pod name
get-pod-name:
	@kubectl get pods -n $(K8S_NAMESPACE) -l app=go-secure-server -o jsonpath="{.items[0].metadata.name}"

# Auto test inside cluster (takes first pod)
test-cluster-auto:
	@echo --- Testing application inside cluster (auto-select pod) ---
	@for /f "tokens=*" %%i in ('kubectl get pods -n $(K8S_NAMESPACE) -l app^=go-secure-server -o jsonpath^="{.items[0].metadata.name}" 2^>nul') do @set POD_NAME=%%i
	@if "%POD_NAME%"=="" (
		echo [ERROR] No pods found with label app=go-secure-server!
		kubectl get pods -n $(K8S_NAMESPACE) --show-labels
		exit 1
	)
	@echo Testing pod: %POD_NAME%
	@echo 1. Health check:
	@kubectl exec -n $(K8S_NAMESPACE) %POD_NAME% -- curl -s http://localhost:8080/
	@echo.
	@echo 2. Halloween API:
	@kubectl exec -n $(K8S_NAMESPACE) %POD_NAME% -- curl -s http://localhost:8080/api/halloween
	@echo.
	@echo [OK] Cluster test completed for pod: %POD_NAME%

# Specific pod tests (no search required)
test-pod-qbfds:
	@echo --- Testing specific pod: go-secure-server-78545bbf69-qbfds ---
	@kubectl exec -n $(K8S_NAMESPACE) go-secure-server-78545bbf69-qbfds -- curl -s http://localhost:8080/
	@echo [OK] Test completed

test-both-pods:
	@echo --- Testing both pods ---
	@echo Pod 1: go-secure-server-78545bbf69-qbfds
	@kubectl exec -n $(K8S_NAMESPACE) go-secure-server-78545bbf69-qbfds -- curl -s http://localhost:8080/
	@echo.
	@echo Pod 2: go-secure-server-78545bbf69-xgvfp
	@kubectl exec -n $(K8S_NAMESPACE) go-secure-server-78545bbf69-xgvfp -- curl -s http://localhost:8080/
	@echo.
	@echo [OK] Both pods responding

# Detailed application test inside cluster
test-cluster-detailed:
	@echo --- Detailed application test inside cluster ---
	@for /f "tokens=*" %%i in ('kubectl get pods -n $(K8S_NAMESPACE) -l app^=go-secure-server -o jsonpath^="{.items[0].metadata.name}"') do @set POD_NAME=%%i
	@if not defined POD_NAME (echo [ERROR] No pods found! && exit 1)
	@echo Testing pod: %POD_NAME%
	@echo.
	@echo 1. Health endpoint:
	@kubectl exec -n $(K8S_NAMESPACE) %POD_NAME% -- curl -s http://localhost:8080/
	@echo.
	@echo 2. Halloween API:
	@kubectl exec -n $(K8S_NAMESPACE) %POD_NAME% -- curl -s http://localhost:8080/api/halloween
	@echo.
	@echo 3. Security headers check:
	@kubectl exec -n $(K8S_NAMESPACE) %POD_NAME% -- curl -I -s http://localhost:8080/ | findstr "X-Content-Type-Options\|X-Frame-Options\|Content-Security-Policy"
	@echo.
	@echo 4. Service connectivity from another pod:
	@echo Run in separate terminal: make test-service-connectivity
	@echo.
	@echo [OK] Detailed cluster test completed

# Test service accessibility from another pod
test-service-connectivity:
	@echo --- Testing service connectivity from another pod ---
	@echo Creating test pod to connect to go-secure-service...
	@kubectl run test-connectivity --rm -i --image=curlimages/curl --restart=Never --namespace $(K8S_NAMESPACE) -- curl -s http://go-secure-service:8080/ | head -2

# Validate entire deployment
validate-deployment: check-pods test-both-pods
	@echo.
	@echo [SUCCESS] Deployment validated!
	@echo [OK] Pods are running
	@echo [OK] Application is responding
	@echo [NEXT] Access locally: make k8s-access

# ============================================================================
# NETWORK POLICY TARGETS
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
	@echo 2. Run manual test (run this in a separate terminal):
	@echo make test-service-connectivity
	@echo.
	@echo [OK] Manual testing instructions displayed

netpol-test-step1:
	@echo Testing Port 80 (should fail)...
	@kubectl run test-80 --rm -i --image=alpine --restart=Never --namespace $(K8S_NAMESPACE) --command -- sh -c "wget -qO- --timeout=5 http://go-secure-service:80/"

netpol-test-step2:
	@echo Testing Port 8080 (should work)...
	@kubectl run test-8080 --rm -i --image=alpine --restart=Never --namespace $(K8S_NAMESPACE) --command -- sh -c "wget -qO- --timeout=5 http://go-secure-service:8080/health"

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
	@echo [DEBUG] View logs: make check-logs
	@echo [TEST] Test inside cluster: make test-pod-qbfds

security-quick-push: push
	@echo --- Deploying to Kubernetes (using GHCR image) ---
	@kubectl create namespace $(K8S_NAMESPACE) --dry-run=client -o yaml | kubectl apply -f -
	@kubectl apply -f k8s/ -n $(K8S_NAMESPACE)
	@echo.
	@echo [SUCCESS] DEPLOYMENT WITH PUSH COMPLETED!
	@echo.
	@echo [NEXT] Access with: make k8s-access

security-full: push scan
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
	@echo "   make k8s-access            - Access app via port-forward"
	@echo.
	@echo "DOCKER & REGISTRY (GHCR):"
	@echo "   make build                 - Build image locally"
	@echo "   make docker-login          - Show GHCR login instructions"
	@echo "   make push                  - Push image to GHCR (requires login)"
	@echo "   make run                   - Run locally"
	@echo.
	@echo "SECURITY (SCAN & FALCO):"
	@echo "   make scan                  - Vulnerability scan (High/Critical)"
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
	@echo "   make test-pod-qbfds        - Test specific pod (quick)"
	@echo "   make test-both-pods        - Test both pods"
	@echo "   make test-cluster-auto     - Test app inside cluster (auto)"
	@echo "   make test-cluster-detailed - Detailed cluster test"
	@echo "   make test-service-connectivity - Test service access"
	@echo "   make validate-deployment   - Full deployment validation"
	@echo.
	@echo "NETWORK POLICY:"
	@echo "   make netpol-apply          - Apply network policy"
	@echo "   make netpol-test           - Test network policy"
	@echo.
	@echo "================================================================================"

.PHONY: all build push run test stop clean logs \
        check-prerequisites docker-login \
        scan vulnerability-scan scan-docker \
        falco-install-clean falco-logs falco-test falco-clean \
        k8s-check k8s-deploy k8s-access k8s-run k8s-clean \
        check-pods check-logs check-status get-pod-name \
        test-cluster-auto test-pod-qbfds test-both-pods test-cluster-detailed test-service-connectivity validate-deployment \
        netpol-apply netpol-clean netpol-test netpol-test-step1 netpol-test-step2 \
        security-quick security-quick-push security-full security-clean help