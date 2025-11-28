PROJECT_NAME := gok8ssecuritylab
IMAGE_NAME := $(PROJECT_NAME)/go-secure-server
IMAGE_VERSION := latest
FULL_IMAGE_NAME := $(IMAGE_NAME):$(IMAGE_VERSION)
CONTAINER_NAME := go-secure-server-test
K8S_NAMESPACE := go-security-lab
KUBE_CONTEXT := docker-desktop
FALCO_NAMESPACE := falco

# ============================================================================
# MAIN TARGETS
# ============================================================================

all: security-quick

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
# FALCO RUNTIME SECURITY TARGETS (ОБНОВЛЕННЫЕ)
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

helm-cleanup:
	@echo --- Cleaning up stuck Helm operations ---
	@-helm rollback falco -n $(FALCO_NAMESPACE) 0 --cleanup-on-fail 2>nul || echo "No rollback needed"
	@-helm delete falco -n $(FALCO_NAMESPACE) --wait 2>nul || echo "No release to delete"
	@-kubectl delete secrets -n $(FALCO_NAMESPACE) -l "owner=helm" --ignore-not-found=true
	@echo "✅ Helm cleanup completed"

falco-clean-complete:
	@echo --- Complete Falco Cleanup ---
	@echo "🧹 Removing all Falco resources..."
	-helm uninstall falco -n $(FALCO_NAMESPACE) --wait 2>nul || echo "No helm release found"
	-kubectl delete namespace $(FALCO_NAMESPACE) --ignore-not-found=true --wait --timeout=60s
	-kubectl delete role,rolebinding,serviceaccount,configmap,secret,daemonset,service,deployment -l app=falco -A --ignore-not-found=true --wait
	-kubectl delete role,rolebinding,serviceaccount,configmap,secret,daemonset,service,deployment -l app.kubernetes.io/name=falco -A --ignore-not-found=true --wait
	@echo "✅ Complete cleanup finished"
	@echo "Wait 30 seconds for resources to be fully removed..."
	@timeout /t 30 /nobreak >nul

falco-install-simple: falco-check k8s-check
	@echo --- Simple Falco Installation ---
	helm repo add falcosecurity https://falcosecurity.github.io/charts --force-update
	helm repo update
	@echo "📦 Installing Falco (simple version)..."
	helm upgrade --install falco falcosecurity/falco \
		--namespace $(FALCO_NAMESPACE) \
		--create-namespace \
		--wait
	@echo "✅ Falco installed successfully"

falco-install-no-wait: falco-check k8s-check helm-cleanup
	@echo --- Installing Falco without wait ---
	helm repo add falcosecurity https://falcosecurity.github.io/charts --force-update
	helm repo update
	@echo "📦 Installing Falco (without wait flag)..."
	helm upgrade --install falco falcosecurity/falco \
		--namespace $(FALCO_NAMESPACE) \
		--create-namespace \
		--set falcosidekick.enabled=true \
		--set falcosidekick.webui.enabled=true \
		--set tty=true
	@echo "⌛ Waiting separately for pods..."
	@timeout /t 30 /nobreak >nul
	@kubectl wait --for=condition=ready pod -l app=falco -n $(FALCO_NAMESPACE) --timeout=120s
	@echo "✅ Falco installation completed!"

falco-install: falco-install-no-wait
	@echo "🎉 Falco installation completed!"

falco-install-clean: falco-clean-complete falco-install-simple
	@echo "🎉 Fresh Falco installation completed!"

falco-logs:
	@echo --- Falco Security Events ---
	kubectl logs -l app=falco -n $(FALCO_NAMESPACE) -c falco --tail=50 -f

falco-test:
	@echo --- Testing Falco Detection ---
	@echo "🧪 Generating test security event..."
	kubectl run --rm -i --tty falco-test --image=alpine --restart=Never -- sh -c "cat /etc/shadow" || true
	@echo "📋 Check logs with: make falco-logs"

falco-status:
	@echo --- Falco Status ---
	kubectl get pods,svc,daemonset -n $(FALCO_NAMESPACE)

falco-clean:
	@echo --- Removing Falco ---
	helm uninstall falco -n $(FALCO_NAMESPACE) --wait
	kubectl delete namespace $(FALCO_NAMESPACE) --ignore-not-found=true
	@echo "✅ Falco removed"

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

security-quick: build k8s-deploy
	@echo.
	@echo "🎉 QUICK SECURITY DEPLOYMENT COMPLETED!"
	@echo "✅ Docker image built"
	@echo "✅ K8s application deployed"
	@echo "✅ Basic security context applied"
	@echo.
	@echo "🔍 For full security stack run: make security-full"
	@echo "🚀 Quick access: make k8s-access"

security-pipeline: scan falco-install-simple k8s-deploy
	@echo.
	@echo "🎉 SECURITY PIPELINE COMPLETED!"
	@echo "✅ Docker image built and scanned"
	@echo "✅ Falco runtime security installed"
	@echo "✅ K8s manifests validated and deployed"
	@echo "✅ Application running in secure environment"
	@echo.
	@echo "🔍 NEXT STEPS:"
	@echo "  make k8s-access     - Access your secured application"
	@echo "  make falco-logs     - Monitor security events"
	@echo "  make falco-test     - Test Falco detection"

security-full: scan falco-install-clean k8s-deploy
	@echo.
	@echo "🎉 FULL SECURITY STACK DEPLOYED!"
	@echo "✅ Docker image built and scanned"
	@echo "✅ Falco runtime security installed (clean install)"
	@echo "✅ K8s application deployed with security context"
	@echo.
	@echo "🛡️  SECURITY MONITORING ACTIVE!"
	@echo "🔍 Run: make falco-logs    - View security events"
	@echo "🚀 Run: make k8s-access    - Access your application"

security-scan: scan k8s-scan
	@echo --- Security Scanning Complete ---

security-clean: k8s-clean falco-clean clean
	@echo --- Complete Security Environment Cleanup ---

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
	@echo "   make all              - Quick security deployment"
	@echo "   make security-full    - Complete security stack"
	@echo "   make run              - Quick Docker run"
	@echo.
	@echo "🔒 SECURITY COMMANDS:"
	@echo "   make scan             - Trivy vulnerability scan"
	@echo "   make falco-install    - Install runtime security"
	@echo "   make falco-install-simple - Simple Falco install"
	@echo "   make falco-test       - Test security detection"
	@echo "   make falco-logs       - Security monitoring"
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
        falco-check install-helm install-helm-choco helm-cleanup \
        falco-clean-complete falco-install-simple falco-install-no-wait \
        falco-install falco-install-clean falco-logs falco-test \
        falco-status falco-clean \
        k8s-check k8s-deploy k8s-access k8s-run k8s-logs k8s-status \
        k8s-clean k8s-update k8s-scan \
        security-quick security-pipeline security-full security-scan \
        security-clean security-status info help