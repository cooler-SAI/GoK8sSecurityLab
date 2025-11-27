PROJECT_NAME := gok8ssecuritylab
IMAGE_NAME := $(PROJECT_NAME)/go-secure-server
IMAGE_VERSION := latest
FULL_IMAGE_NAME := $(IMAGE_NAME):$(IMAGE_VERSION)
CONTAINER_NAME := go-secure-server-test
K8S_NAMESPACE := go-security-lab
KUBE_CONTEXT := docker-desktop
FALCO_NAMESPACE := falco
FALCO_VERSION := 3.5.4

# ============================================================================
# MAIN TARGETS
# ============================================================================

all: security-full

build:
	@echo --- Building Docker Image: $(FULL_IMAGE_NAME) ---
	docker build --pull -t $(FULL_IMAGE_NAME) .
	@echo --- Successfully built: $(FULL_IMAGE_NAME) ---

run: stop build
	@echo --- Running container on http://localhost:8080 ---
	docker run -d --rm --name $(CONTAINER_NAME) -p 8080:8080 $(FULL_IMAGE_NAME)
	@echo --- Container $(CONTAINER_NAME) started successfully ---
	@echo.
	@echo   #######################################################
	@echo   #          SECURE HALLOWEEN SERVER STARTED           #
	@echo   #######################################################
	@echo.
	@echo   *** AVAILABLE ENDPOINTS ***
	@echo   [1] Health Check:    http://localhost:8080/
	@echo   [2] Halloween Page:  http://localhost:8080/halloween
	@echo   [3] Halloween API:   http://localhost:8080/api/halloween
	@echo   [4] Server Info:     http://localhost:8080/info
	@echo.
	@echo   *** SECURITY FEATURES ***
	@echo   [X] Rate Limiting (10 requests/second)
	@echo   [X] Security Headers
	@echo   [X] Non-root user execution
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
	@echo --- All tests completed ---

stop:
	@echo --- Stopping container: $(CONTAINER_NAME) ---
	@-docker stop $(CONTAINER_NAME) 2>nul && echo Container stopped || echo Container was not running

clean: stop
	@echo --- Cleaning up resources ---
	@-docker rmi $(FULL_IMAGE_NAME) 2>nul && echo Image deleted || echo Image not found
	@-docker system prune -f 2>nul && echo System cleaned || echo Cleanup completed
	@echo Cleanup completed

logs:
	@echo --- Showing logs for $(CONTAINER_NAME) ---
	@docker logs -f $(CONTAINER_NAME)

status:
	@docker ps -a | findstr "$(CONTAINER_NAME)" >nul && echo Container is running || echo Container is not running

open-all:
	@echo Opening all endpoints in browser...
	@start http://localhost:8080/ 2>nul
	@timeout /t 1 /nobreak >nul
	@start http://localhost:8080/halloween 2>nul
	@timeout /t 1 /nobreak >nul
	@start http://localhost:8080/api/halloween 2>nul
	@echo All endpoints opened in browser tabs

deploy: run
	@echo.
	@echo *** QUICK ACCESS ***
	@echo Health:    http://localhost:8080/
	@echo Halloween: http://localhost:8080/halloween
	@echo API:       http://localhost:8080/api/halloween
	@echo.
	@echo Use 'make open-all' to open all links in browser

restart: stop run

# ============================================================================
# SECURITY SCANNING TARGETS
# ============================================================================

scan: build
	@echo --- Scanning Docker Image for HIGH/CRITICAL vulnerabilities using Trivy ---
	@"C:\Users\coole\scoop\shims\trivy.exe" image --severity HIGH,CRITICAL $(FULL_IMAGE_NAME)
	@echo --- Trivy scan completed ---

scan-fast: build
	@echo --- Fast scan (vulnerabilities only, no secrets) ---
	@"C:\Users\coole\scoop\shims\trivy.exe" image --scanners vuln --severity HIGH,CRITICAL $(FULL_IMAGE_NAME)

scan-full: build
	@echo --- Full security scan with secrets detection ---
	@"C:\Users\coole\scoop\shims\trivy.exe" image --severity HIGH,CRITICAL $(FULL_IMAGE_NAME)

scan-docker: build
	@echo --- Scanning with Docker Scout ---
	docker scout quickview $(FULL_IMAGE_NAME)

# ============================================================================
# FALCO RUNTIME SECURITY TARGETS (IMPROVED)
# ============================================================================

falco-check:
	@echo --- Checking Falco prerequisites ---
	@where helm >nul 2>&1 || (echo "❌ Helm not found. Installing Helm..." && $(MAKE) install-helm)
	@echo ✓ Helm is available
	@kubectl version --client >nul 2>&1 || (echo "ERROR: kubectl is required." && exit 1)
	@echo ✓ kubectl is available
	@echo --- Prerequisites check passed ---

install-helm:
	@echo --- Installing Helm ---
	@scoop install helm || (echo "Failed to install via scoop. Please install Helm manually: https://helm.sh/docs/intro/install/" && exit 1)
	@echo ✓ Helm installed successfully

install-helm-choco:
	@echo --- Installing Helm via Chocolatey ---
	choco install kubernetes-helm -y
	@echo ✓ Helm installed via Chocolatey

falco-install: falco-check k8s-check
	@echo --- Installing Falco Runtime Security ---
	@echo "🔒 Deploying Falco for real-time security monitoring"
	helm repo add falcosecurity https://falcosecurity.github.io/charts --force-update
	helm repo update
	@echo "📦 Installing Falco version $(FALCO_VERSION)..."
	helm upgrade --install falco falcosecurity/falco \
		--namespace $(FALCO_NAMESPACE) \
		--create-namespace \
		--version $(FALCO_VERSION) \
		--set falcosidekick.enabled=true \
		--set falcosidekick.webui.enabled=true \
		--set tty=true \
		--wait
	@echo --- Waiting for Falco pods to be ready ---
	@kubectl wait --for=condition=ready pod -l app=falco -n $(FALCO_NAMESPACE) --timeout=180s
	@echo.
	@echo "🎉 FALCO SUCCESSFULLY INSTALLED!"
	@echo "   📊 Web UI: http://localhost:2802"
	@echo "   📋 Logs:   make falco-logs"
	@echo "   🧪 Test:   make falco-test"

# Альтернативная установка без Helm (если проблемы с Helm)
falco-install-direct: k8s-check
	@echo --- Installing Falco using direct manifests ---
	kubectl create namespace $(FALCO_NAMESPACE) --dry-run=client -o yaml | kubectl apply -f -
	kubectl apply -f https://raw.githubusercontent.com/falcosecurity/charts/master/falco/templates/namespace.yaml
	@echo "⚠️  For direct installation, see: https://falco.org/docs/getting-started/installation/"
	@echo "📖 Or use: kubectl apply -f https://falco.org/helm/install"
# ============================================================================
# KUBERNETES TARGETS
# ============================================================================

k8s-check:
	@echo --- Checking Kubernetes Access ---
	kubectl config use-context $(KUBE_CONTEXT)
	kubectl cluster-info
	@echo --- Kubernetes is ready! ---

k8s-deploy: build k8s-check
	@echo --- Deploying to Docker Desktop Kubernetes ---
	kubectl create namespace $(K8S_NAMESPACE) --dry-run=client -o yaml | kubectl apply -f -
	kubectl apply -f k8s/ -n $(K8S_NAMESPACE)
	@echo --- Waiting for deployment to be ready ---
	kubectl wait --for=condition=ready pod -l app=go-secure-server -n $(K8S_NAMESPACE) --timeout=120s
	@echo --- Deployment completed! ---

k8s-access:
	@echo --- Starting Port Forward ---
	@echo "🎃 Your Go Secure Server is now in K8s with Runtime Security!"
	@echo ""
	@echo "📡 AVAILABLE ENDPOINTS:"
	@echo "   🔹 Health Check:    http://localhost:8080/"
	@echo "   🔹 Halloween Page:  http://localhost:8080/halloween"
	@echo "   🔹 Halloween API:   http://localhost:8080/api/halloween"
	@echo "   🔹 Server Info:     http://localhost:8080/info"
	@echo ""
	@echo "⚡ SECURITY FEATURES:"
	@echo "   ✅ Static Analysis (Trivy)"
	@echo "   ✅ Runtime Security (Falco)"
	@echo "   ✅ K8s Security Context"
	@echo "   ✅ Non-root user execution"
	@echo ""
	@echo "Press Ctrl+C to stop port-forwarding"
	kubectl port-forward svc/go-secure-service 8080:80 -n $(K8S_NAMESPACE)

k8s-run: k8s-deploy
	@echo --- Full Kubernetes Deployment Complete ---
	@echo --- Starting automatic port-forward ---
	$(MAKE) k8s-access

k8s-logs:
	@echo --- K8s Application Logs ---
	kubectl logs -l app=go-secure-server -n $(K8S_NAMESPACE) --tail=50 -f

k8s-status:
	@echo --- Kubernetes Cluster Status ---
	@echo "📊 Namespace: $(K8S_NAMESPACE)"
	kubectl get pods,svc,deploy -n $(K8S_NAMESPACE)
	@echo.
	@echo "🔒 Security Context:"
	kubectl get pod -l app=go-secure-server -n $(K8S_NAMESPACE) -o jsonpath='{.items[*].spec.securityContext}' | echo

k8s-clean:
	@echo --- Cleaning Kubernetes Resources ---
	kubectl delete -f k8s/ -n $(K8S_NAMESPACE) --ignore-not-found=true --wait
	kubectl delete namespace $(K8S_NAMESPACE) --ignore-not-found=true --timeout=60s
	@echo --- K8s cleanup completed ---

k8s-update: build
	@echo --- Updating K8s Deployment ---
	kubectl rollout restart deployment/go-secure-server -n $(K8S_NAMESPACE)
	kubectl wait --for=condition=ready pod -l app=go-secure-server -n $(K8S_NAMESPACE) --timeout=120s
	@echo --- Update completed ---

k8s-scan: build
	@echo --- Comprehensive Kubernetes Security Scan ---
	@echo "🔍 Docker Image Scan:"
	@"C:\Users\coole\scoop\shims\trivy.exe" image --severity HIGH,CRITICAL $(FULL_IMAGE_NAME)
	@echo.
	@echo "🛡️ K8s Manifest Validation:"
	kubectl apply -f k8s/ -n $(K8S_NAMESPACE) --dry-run=server --validate=true
	@echo.
	@echo "✅ Kubernetes security scan completed"

# ============================================================================
# COMPREHENSIVE SECURITY PIPELINES
# ============================================================================

security-pipeline: scan falco-install k8s-deploy
	@echo.
	@echo "🎉 SECURITY PIPELINE COMPLETED!"
	@echo "✅ Docker image built and scanned"
	@echo "✅ Falco runtime security installed"
	@echo "✅ K8s manifests validated and deployed"
	@echo "✅ Application running in secure environment"
	@echo.
	@echo "🔍 NEXT STEPS:"
	@echo "  make k8s-access     - Access your secured application"
	@echo "  make falco-ui       - Monitor security events"
	@echo "  make falco-test     - Test Falco detection"

security-full: scan falco-install k8s-deploy falco-ui
	@echo --- Complete Security Stack Activated ---


falco-ui:
security-scan: scan k8s-scan
	@echo --- Security Scanning Complete ---

security-clean: k8s-clean falco-clean clean
	@echo --- Complete Security Environment Cleanup ---


falco-clean:
security-status:
	@echo --- Comprehensive Security Status ---
	@echo.
	@echo "🔍 STATIC ANALYSIS:"
	@-docker images | findstr "$(IMAGE_NAME)" >nul && echo "✅ Docker image built" || echo "❌ Docker image not built"
	@echo.
	@echo "🔒 RUNTIME SECURITY:"
	@-kubectl get pods -n $(FALCO_NAMESPACE) 2>nul | findstr "falco" >nul && echo "✅ Falco installed" || echo "❌ Falco not installed"
	@echo.
	@echo "🚀 KUBERNETES:"
	@-kubectl get pods -n $(K8S_NAMESPACE) 2>nul | findstr "go-secure-server" >nul && echo "✅ Application deployed" || echo "❌ Application not deployed"
	@echo.
	@echo "📊 SECURITY MONITORING:"
	@echo "   make falco-logs    - View security events"
	@echo "   make falco-ui      - Web dashboard"
	@echo "   make k8s-logs      - Application logs"

# ============================================================================
# INFORMATION TARGETS
# ============================================================================

info:
	@echo.
	@echo "🦀 GO K8S SECURITY LAB - COMPLETE MAKEFILE"
	@echo "==========================================="
	@echo.
	@echo "📦 PROJECT INFO:"
	@echo "   Image:          $(FULL_IMAGE_NAME)"
	@echo "   Container:      $(CONTAINER_NAME)"
	@echo "   K8s Namespace:  $(K8S_NAMESPACE)"
	@echo "   Falco Namespace: $(FALCO_NAMESPACE)"
	@echo.
	@echo "🚀 QUICK START:"
	@echo "   make all              - Full security deployment"
	@echo "   make security-full    - Complete security stack"
	@echo "   make run              - Quick Docker run"
	@echo.
	@echo "🔒 SECURITY COMMANDS:"
	@echo "   make scan             - Trivy vulnerability scan"
	@echo "   make falco-install    - Install runtime security"
	@echo "   make falco-test       - Test security detection"
	@echo "   make falco-ui         - Security dashboard"
	@echo "   make security-status  - Overall security status"
	@echo.
	@echo "🛠️  DEVELOPMENT:"
	@echo "   make build            - Build Docker image"
	@echo "   make test             - Run tests"
	@echo "   make k8s-deploy       - Deploy to Kubernetes"
	@echo "   make k8s-update       - Update deployment"
	@echo.
	@echo "📊 MONITORING:"
	@echo "   make falco-logs       - Security events"
	@echo "   make k8s-logs         - Application logs"
	@echo "   make k8s-status       - K8s status"
	@echo.
	@echo "🧹 CLEANUP:"
	@echo "   make clean            - Clean Docker resources"
	@echo "   make k8s-clean        - Clean K8s resources"
	@echo "   make falco-clean      - Remove Falco"
	@echo "   make security-clean   - Clean everything"

help: info

# ============================================================================
# UTILITY TARGETS
# ============================================================================

.PHONY: all build run test stop clean logs status open-all deploy restart \
        scan scan-fast scan-full scan-docker \
        falco-check falco-install falco-status falco-logs falco-events \
        falco-test falco-ui falco-monitor falco-clean falco-rules \
        k8s-check k8s-deploy k8s-access k8s-run k8s-logs k8s-status \
        k8s-clean k8s-update k8s-scan \
        security-pipeline security-full security-scan security-clean \
        security-status info help