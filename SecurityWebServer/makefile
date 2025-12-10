# ============================================================================
# PROJECT CONFIGURATION
# ============================================================================
PROJECT_NAME := gok8ssecuritylab
IMAGE_NAME := $(PROJECT_NAME)/go-secure-server
IMAGE_VERSION := latest
FULL_IMAGE_NAME := $(IMAGE_NAME):$(IMAGE_VERSION)
CONTAINER_NAME := go-secure-server-test
K8S_NAMESPACE := go-security-lab
KUBE_CONTEXT := docker-desktop
FALCO_NAMESPACE := falco
TRIVY_BIN := "C:\Users\coole\scoop\shims\trivy.exe"
NETPOL_FILE := k8s/network-policy.yaml

# ============================================================================
# PREREQUISITE CHECKS - SIMPLIFIED FOR WINDOWS
# ============================================================================

check-prerequisites:
	@echo --- Checking Prerequisites ---
	@docker version >nul 2>&1 && echo [OK] Docker found || (echo [ERROR] Docker not found && exit 1)
	@kubectl version --client >nul 2>&1 && echo [OK] kubectl found || (echo [ERROR] kubectl not found && exit 1)
	@helm version >nul 2>&1 && echo [OK] Helm found || (echo [ERROR] Helm not found && exit 1)
	@if exist $(TRIVY_BIN) (echo [OK] Trivy found) else (echo [WARN] Trivy not found)
	@echo [OK] All prerequisites checked

# ============================================================================
# MAIN TARGETS (Local Docker Management)
# ============================================================================

all: security-quick

build:
	@echo --- Building Docker Image: $(FULL_IMAGE_NAME) ---
	docker build --pull -t $(FULL_IMAGE_NAME) .
	@echo [OK] Successfully built: $(FULL_IMAGE_NAME)

run: stop build
	@echo --- Running container on http://localhost:8080 ---
	docker run -d --rm --name $(CONTAINER_NAME) -p 8080:8080 $(FULL_IMAGE_NAME)
	@echo [OK] Container $(CONTAINER_NAME) started successfully
	@echo.
	@echo   #######################################################
	@echo   #          SECURE HALLOWEEN SERVER STARTED           #
	@echo   #######################################################
	@echo.
	@echo Waiting for server to start...
	@timeout /t 3 /nobreak >nul
	@echo Server ready! Opening Halloween page...
	@start http://localhost:8080/halloween 2>nul || echo Manual: Open http://localhost:8080/halloween in browser

test: run
	@echo.
	@echo --- Testing endpoints ---
	@curl -s http://localhost:8080/ | findstr "Status" || echo Health check failed
	@curl -s http://localhost:8080/api/halloween | findstr "spookily_secure" && echo API test passed || echo API test failed
	@echo [OK] All tests completed

stop:
	@echo --- Stopping container: $(CONTAINER_NAME) ---
	@docker stop $(CONTAINER_NAME) 2>nul && echo Container stopped || echo Container was not running

clean: stop
	@echo --- Cleaning up Docker resources ---
	@docker rmi $(FULL_IMAGE_NAME) 2>nul && echo Image deleted || echo Image not found
	@docker system prune -f --volumes 2>nul && echo System cleaned || echo Cleanup completed
	@echo [OK] Docker cleanup completed

logs:
	@echo --- Showing logs for $(CONTAINER_NAME) ---
	@docker logs -f $(CONTAINER_NAME)

# ============================================================================
# SECURITY SCANNING TARGETS
# ============================================================================

scan: build
	@echo --- Scanning Docker Image for HIGH/CRITICAL vulnerabilities using Trivy ---
	$(TRIVY_BIN) image --severity HIGH,CRITICAL $(FULL_IMAGE_NAME)
	@echo [OK] Trivy scan completed

vulnerability-scan: build
	@echo --- Comprehensive Vulnerability Assessment ---
	@echo [INFO] Scanning $(FULL_IMAGE_NAME) for vulnerabilities...
	$(TRIVY_BIN) image --format table --scanners vuln,misconfig --severity MEDIUM,HIGH,CRITICAL $(FULL_IMAGE_NAME)
	@echo [OK] Vulnerability scan completed

scan-docker: build
	@echo --- Scanning with Docker Scout ---
	docker scout quickview $(FULL_IMAGE_NAME)

# ============================================================================
# FALCO RUNTIME SECURITY TARGETS
# ============================================================================

falco-install-clean: k8s-check
	@echo --- Complete Falco Cleanup and Fresh Install ---
	@echo [INFO] Removing previous Falco resources...
	@helm uninstall falco -n $(FALCO_NAMESPACE) --wait 2>nul || echo No helm release found
	@kubectl delete namespace $(FALCO_NAMESPACE) --ignore-not-found=true --timeout=60s
	@echo [INFO] Installing Falco...
	@helm repo add falcosecurity https://falcosecurity.github.io/charts
	@helm repo update
	@helm upgrade --install falco falcosecurity/falco --namespace $(FALCO_NAMESPACE) --create-namespace --set falcosidekick.enabled=false
	@echo [OK] Falco installation completed!

falco-logs:
	@echo --- Falco Security Events ---
	@kubectl logs -l app=falco -n $(FALCO_NAMESPACE) -c falco --tail=50

falco-test:
	@echo --- Testing Falco Detection ---
	@echo [TEST] Generating test security event...
	@kubectl run falco-test --rm -i --tty --image=alpine --restart=Never --namespace $(K8S_NAMESPACE) -- sh -c "cat /etc/shadow"

falco-clean:
	@echo --- Removing Falco ---
	@helm uninstall falco -n $(FALCO_NAMESPACE) --wait 2>nul || echo No Falco release found
	@kubectl delete namespace $(FALCO_NAMESPACE) --ignore-not-found=true --timeout=60s
	@echo [OK] Falco removed

# ============================================================================
# KUBERNETES TARGETS (Deployment & Access)
# ============================================================================

k8s-check:
	@echo --- Checking Kubernetes Access ---
	@kubectl config use-context $(KUBE_CONTEXT)
	@kubectl cluster-info
	@echo [OK] Kubernetes is ready!

k8s-deploy: build k8s-check
	@echo --- Deploying Go Secure Server to Kubernetes ---
	@kubectl create namespace $(K8S_NAMESPACE) --dry-run=client -o yaml | kubectl apply -f -
	@kubectl apply -f k8s/ -n $(K8S_NAMESPACE)
	@echo --- Waiting for deployment to be ready ---
	@kubectl wait --for=condition=ready pod -l app=go-secure-server -n $(K8S_NAMESPACE) --timeout=120s
	@echo [OK] Application Deployment completed!

k8s-access:
	@echo --- Starting Port Forward ---
	@echo [INFO] Your Go Secure Server is now in K8s!
	@echo.
	@echo [ENDPOINTS] Available at http://localhost:8080/ :
	@echo    - Health Check:    http://localhost:8080/
	@echo    - Halloween Page:  http://localhost:8080/halloween
	@echo    - Halloween API:   http://localhost:8080/api/halloween
	@echo.
	@echo Press Ctrl+C to stop port-forwarding
	@kubectl port-forward svc/go-secure-service 8080:8080 -n $(K8S_NAMESPACE)

k8s-run: k8s-deploy k8s-access

k8s-clean:
	@echo --- Cleaning Kubernetes Application Resources ---
	@kubectl delete -f k8s/ -n $(K8S_NAMESPACE) --ignore-not-found=true
	@kubectl delete namespace $(K8S_NAMESPACE) --ignore-not-found=true --timeout=60s
	@echo [OK] K8s application cleanup completed

# ============================================================================
# NETWORK POLICY TARGETS - SIMPLIFIED FOR WINDOWS
# ============================================================================

netpol-apply: k8s-check
	@echo --- Applying Network Policy ---
	@kubectl apply -f $(NETPOL_FILE) -n $(K8S_NAMESPACE)
	@echo [OK] NetworkPolicy applied

netpol-clean: k8s-check
	@echo --- Removing Network Policy ---
	@kubectl delete -f $(NETPOL_FILE) -n $(K8S_NAMESPACE) --ignore-not-found=true
	@echo [OK] NetworkPolicy removed

# SIMPLE NETPOL TEST - ONE COMMAND AT A TIME
netpol-test: k8s-check
	@echo --- Testing Network Policy ---
	@echo 1. Testing service exists...
	@kubectl get svc go-secure-service -n $(K8S_NAMESPACE)
	@echo.
	@echo 2. Running manual test (run this in separate terminal):
	@echo kubectl run test-pod --rm -it --image=curlimages/curl --restart=Never --namespace $(K8S_NAMESPACE) -- curl -v http://go-secure-service:8080/health
	@echo.
	@echo [OK] Manual test instructions displayed

# ALTERNATIVE: STEP-BY-STEP TEST
netpol-test-step1:
	@echo Testing Port 80 (should fail)...
	@kubectl run test-80 --rm -i --image=alpine --restart=Never --namespace $(K8S_NAMESPACE) --command -- sh -c "wget -qO- --timeout=5 http://go-secure-service:80/"

netpol-test-step2:
	@echo Testing Port 8080 (should work)...
	@kubectl run test-8080 --rm -i --image=alpine --restart=Never --namespace $(K8S_NAMESPACE) --command -- sh -c "wget -qO- --timeout=5 http://go-secure-service:8080/health"

# ============================================================================
# DEBUGGING UTILITIES
# ============================================================================

check-service:
	@echo --- Checking Service ---
	@kubectl get svc go-secure-service -n $(K8S_NAMESPACE) -o yaml

check-pods:
	@echo --- Checking Pods ---
	@kubectl get pods -n $(K8S_NAMESPACE) -o wide

check-netpol:
	@echo --- Checking Network Policies ---
	@kubectl get netpol -n $(K8S_NAMESPACE) -o wide

manual-test:
	@echo --- Manual Test Command ---
	@echo Run this command to test connectivity:
	@echo kubectl run test-curl --rm -it --image=curlimages/curl --restart=Never --namespace $(K8S_NAMESPACE) -- curl -v http://go-secure-service:8080/health

# ============================================================================
# COMPREHENSIVE SECURITY PIPELINES - SIMPLIFIED
# ============================================================================

security-quick: build k8s-deploy
	@echo.
	@echo [SUCCESS] QUICK SECURITY DEPLOYMENT COMPLETED!
	@echo [OK] Docker image built
	@echo [OK] K8s application deployed
	@echo.
	@echo [NEXT] Access with: make k8s-access

security-full: build scan k8s-deploy netpol-apply
	@echo.
	@echo [SUCCESS] FULL SECURITY STACK DEPLOYED!
	@echo [OK] Docker image built and scanned
	@echo [OK] K8s application deployed
	@echo [OK] Network Policy applied
	@echo.
	@echo [NEXT] Access with: make k8s-access

security-clean: k8s-clean falco-clean clean
	@echo --- Complete Security Environment Cleanup ---

# ============================================================================
# INFORMATION & HELP
# ============================================================================

help:
	@echo.
	@echo "================================================================================"
	@echo "GO K8S SECURITY LAB - WINDOWS MAKEFILE"
	@echo "================================================================================"
	@echo.
	@echo "QUICK START:"
	@echo "   make security-quick   - Build and deploy"
	@echo "   make security-full    - Full security stack"
	@echo "   make k8s-access       - Access application"
	@echo.
	@echo "UTILITIES:"
	@echo "   make check-service    - Check service"
	@echo "   make check-pods       - Check pods"
	@echo "   make check-netpol     - Check network policies"
	@echo "   make manual-test      - Get test command"
	@echo.
	@echo "DOCKER:"
	@echo "   make build            - Build image"
	@echo "   make run              - Run locally"
	@echo "   make stop             - Stop container"
	@echo.
	@echo "SECURITY:"
	@echo "   make scan             - Scan for vulnerabilities"
	@echo "   make falco-install-clean - Install Falco"
	@echo "   make falco-logs       - View security events"
	@echo.
	@echo "================================================================================"

.PHONY: all build run test stop clean logs \
        check-prerequisites \
        scan vulnerability-scan scan-docker \
        falco-install-clean falco-logs falco-test falco-clean \
        k8s-check k8s-deploy k8s-access k8s-run k8s-clean \
        netpol-apply netpol-clean netpol-test netpol-test-step1 netpol-test-step2 \
        check-service check-pods check-netpol manual-test \
        security-quick security-full security-clean help