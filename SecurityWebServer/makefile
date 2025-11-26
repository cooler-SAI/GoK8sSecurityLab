PROJECT_NAME := gok8ssecuritylab
IMAGE_NAME := $(PROJECT_NAME)/go-secure-server
IMAGE_VERSION := latest
FULL_IMAGE_NAME := $(IMAGE_NAME):$(IMAGE_VERSION)
CONTAINER_NAME := go-secure-server-test

all: build

build:
	@echo --- Building Docker Image: $(FULL_IMAGE_NAME) ---
	docker build --pull -t $(FULL_IMAGE_NAME) .
	@echo --- Successfully built: $(FULL_IMAGE_NAME) ---

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

# Docker Scout as an alternative
scan-docker: build
	@echo --- Scanning with Docker Scout ---
	docker scout quickview $(FULL_IMAGE_NAME)

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
	@echo   [X] X-Content-Type-Options: nosniff
	@echo   [X] X-Frame-Options: DENY
	@echo   [X] X-XSS-Protection enabled
	@echo   [X] Secure Headers
	@echo.
	@echo   *** QUICK TEST ***
	@echo   Try refreshing http://localhost:8080/halloween quickly
	@echo   to see rate limiting in action!
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

info:
	@echo.
	@echo PROJECT INFORMATION:
	@echo    Image: $(FULL_IMAGE_NAME)
	@echo    Container: $(CONTAINER_NAME)
	@echo    Port: 8080
	@echo.
	@echo AVAILABLE COMMANDS:
	@echo    make build      - Build Docker image
	@echo    make run        - Build and run container
	@echo    make test       - Run tests
	@echo    make stop       - Stop container
	@echo    make clean      - Stop and remove everything
	@echo    make info       - Show this info
	@echo    make open-all   - Open all endpoints in browser
	@echo    make logs       - Show container logs
	@echo    make status     - Check container status
	@echo    make scan       - Security scan with Trivy
	@echo    make scan-fast  - Fast vulnerability scan
	@echo    make scan-full  - Full security scan
	@echo    make scan-docker - Scan with Docker Scout

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
# KUBERNETES TARGETS
# ============================================================================

K8S_NAMESPACE := go-security-lab
KUBE_CONTEXT := docker-desktop

# K8s check
k8s-check:
	@echo --- Checking Kubernetes Access ---
	kubectl config use-context $(KUBE_CONTEXT)
	kubectl cluster-info
	@echo --- Kubernetes is ready! ---

# Deploy to K8s
k8s-deploy: build k8s-check
	@echo --- Deploying to Docker Desktop Kubernetes ---
	kubectl create namespace $(K8S_NAMESPACE) --dry-run=client -o yaml | kubectl apply -f -
	kubectl apply -f k8s/ -n $(K8S_NAMESPACE)
	@echo --- Deployment completed! ---

# Port-forward for access
k8s-access:
	@echo --- Starting Port Forward ---
	@echo "🎃 Your Go Secure Server is now in K8s!"
	@echo ""
	@echo "📡 AVAILABLE ENDPOINTS:"
	@echo "   🔹 Health Check:    http://localhost:8080/"
	@echo "   🔹 Halloween Page:  http://localhost:8080/halloween"
	@echo "   🔹 Halloween API:   http://localhost:8080/api/halloween"
	@echo "   🔹 Server Info:     http://localhost:8080/info"
	@echo ""
	@echo "⚡ SECURITY FEATURES (now in K8s!):"
	@echo "   ✅ Rate Limiting (10 requests/second)"
	@echo "   ✅ Security Headers"
	@echo "   ✅ K8s Security Context"
	@echo "   ✅ Non-root user execution"
	@echo ""
	@echo "Press Ctrl+C to stop port-forwarding"
	kubectl port-forward svc/go-secure-service 8080:80 -n $(K8S_NAMESPACE)

# Full deploy with port-forward
k8s-run: k8s-deploy
	@echo --- Waiting for pods to be ready ---
	kubectl wait --for=condition=ready pod -l app=go-secure-server -n $(K8S_NAMESPACE) --timeout=120s
	@echo --- Starting automatic port-forward ---
	$(MAKE) k8s-access

# K8s logs
k8s-logs:
	@echo --- K8s Container Logs ---
	kubectl logs -l app=go-secure-server -n $(K8S_NAMESPACE) --tail=50 -f

# K8s state
k8s-status:
	@echo --- Kubernetes Status ---
	@echo "📊 Namespace: $(K8S_NAMESPACE)"
	kubectl get pods,svc,deploy -n $(K8S_NAMESPACE)

# K8s clean
k8s-clean:
	@echo --- Cleaning Kubernetes Resources ---
	kubectl delete -f k8s/ -n $(K8S_NAMESPACE) --ignore-not-found=true
	kubectl delete namespace $(K8S_NAMESPACE) --ignore-not-found=true
	@echo --- K8s cleanup completed ---

# Renew deployment (after make build)
k8s-update: build
	@echo --- Updating K8s Deployment ---
	kubectl rollout restart deployment/go-secure-server -n $(K8S_NAMESPACE)
	@echo --- Update initiated. Check status with: make k8s-status ---

# Security scan в K8s
k8s-scan: build
	@echo --- Comprehensive Security Scan ---
	@echo "🔍 Docker Image Scan:"
	@"C:\Users\coole\scoop\shims\trivy.exe" image --severity HIGH,CRITICAL $(FULL_IMAGE_NAME)
	@echo ""
	@echo "🛡️ K8s Security Check:"
	kubectl apply -f k8s/ -n $(K8S_NAMESPACE) --dry-run=server --validate=true

# Full security pipeline
security-pipeline: scan k8s-scan k8s-deploy
	@echo "🎉 SECURITY PIPELINE COMPLETED!"
	@echo "✅ Docker image built and scanned"
	@echo "✅ K8s manifests validated"
	@echo "✅ Application deployed to Kubernetes"
	@echo ""
	@echo "Next: make k8s-run to access your secured app!"

# K8s info
k8s-info:
	@echo.
	@echo "🚀 KUBERNETES DEPLOYMENT INFO:"
	@echo "   Project:    $(PROJECT_NAME)"
	@echo "   Image:      $(FULL_IMAGE_NAME)"
	@echo "   Namespace:  $(K8S_NAMESPACE)"
	@echo "   Context:    $(KUBE_CONTEXT)"
	@echo.
	@echo "📋 AVAILABLE K8s COMMANDS:"
	@echo "   make k8s-run        - Full deploy + auto port-forward"
	@echo "   make k8s-deploy     - Deploy to K8s only"
	@echo "   make k8s-access     - Start port-forward"
	@echo "   make k8s-status     - Check K8s status"
	@echo "   make k8s-logs       - Show K8s logs"
	@echo "   make k8s-update     - Update deployment"
	@echo "   make k8s-clean      - Clean K8s resources"
	@echo "   make k8s-scan       - Security scan in K8s"
	@echo "   make security-pipeline - Full security workflow"

# Update info target for enable K8s without
info:
	@echo.
	@echo "PROJECT INFORMATION:"
	@echo "   Image:      $(FULL_IMAGE_NAME)"
	@echo "   Container:  $(CONTAINER_NAME)"
	@echo "   Port:       8080"
	@echo "   K8s Namespace: $(K8S_NAMESPACE)"
	@echo.
	@echo "🦀 DOCKER COMMANDS:"
	@echo "   make build      - Build Docker image"
	@echo "   make run        - Build and run container"
	@echo "   make test       - Run tests"
	@echo "   make stop       - Stop container"
	@echo "   make clean      - Stop and remove everything"
	@echo "   make logs       - Show container logs"
	@echo "   make status     - Check container status"
	@echo "   make scan       - Security scan with Trivy"
	@echo "   make scan-fast  - Fast vulnerability scan"
	@echo "   make scan-full  - Full security scan"
	@echo "   make scan-docker - Scan with Docker Scout"
	@echo "   make open-all   - Open all endpoints in browser"
	@echo "   make deploy     - Run and show quick access"
	@echo "   make restart    - Stop and run again"
	@echo.
	@echo "🚀 KUBERNETES COMMANDS:"
	@echo "   make k8s-run    - Full K8s deploy + access"
	@echo "   make k8s-deploy - Deploy to K8s"
	@echo "   make k8s-status - Check K8s status"
	@echo "   make k8s-logs   - Show K8s logs"
	@echo "   make k8s-clean  - Clean K8s resources"
	@echo "   make k8s-info   - Show K8s help"

.PHONY: k8s-check k8s-deploy k8s-access k8s-run k8s-logs k8s-status k8s-clean k8s-update k8s-scan security-pipeline k8s-info