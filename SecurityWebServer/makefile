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

# ASCII symbols for statuses - NO COLORS (professional and works everywhere)
CHECK := [OK]
CROSS := [ERROR]
WARN := [WARN]
INFO := [INFO]
SEC := [SEC]
LOG := [LOG]
TEST := [TEST]
GO := [GO]
DONE := [DONE]
SCAN := [SCAN]

# ============================================================================
# MAIN TARGETS (Local Docker Management)
# ============================================================================

all: security-quick

build:
	@echo "--- Building Docker Image: $(FULL_IMAGE_NAME) ---"
	docker build --pull -t $(FULL_IMAGE_NAME) .
	@echo "--- Successfully built: $(FULL_IMAGE_NAME) ---"

run: stop build
	@echo "--- Running container on http://localhost:8080 ---"
	docker run -d --rm --name $(CONTAINER_NAME) -p 8080:8080 $(FULL_IMAGE_NAME)
	@echo "--- Container $(CONTAINER_NAME) started successfully ---"
	@echo.
	@echo "  #######################################################"
	@echo "  #          SECURE HALLOWEEN SERVER STARTED           #"
	@echo "  #######################################################"
	@echo.
	@echo "Waiting for server to start..."
	@timeout /t 3 /nobreak >nul
	@echo "$(CHECK) Server ready! Opening Halloween page..."
	@start http://localhost:8080/halloween 2>nul || echo "$(INFO) Manual: Open http://localhost:8080/halloween in browser"

test: run
	@echo.
	@echo "--- Testing endpoints ---"
	@curl -s http://localhost:8080/ | findstr "Status" || echo "$(CROSS) Health check failed"
	@curl -s http://localhost:8080/api/halloween | findstr "spookily_secure" && echo "$(CHECK) API test passed" || echo "$(CROSS) API test failed"
	@echo "--- All tests completed ---"

stop:
	@echo "--- Stopping container: $(CONTAINER_NAME) ---"
	@-docker stop $(CONTAINER_NAME) 2>nul && echo "$(CHECK) Container stopped" || echo "$(INFO) Container was not running"

clean: stop
	@echo "--- Cleaning up Docker resources ---"
	@-docker rmi $(FULL_IMAGE_NAME) 2>nul && echo "$(CHECK) Image deleted" || echo "$(INFO) Image not found"
	@-docker system prune -f --volumes 2>nul && echo "$(CHECK) System cleaned" || echo "$(CHECK) Cleanup completed"
	@echo "Docker cleanup completed"

logs:
	@echo "--- Showing logs for $(CONTAINER_NAME) ---"
	@docker logs -f $(CONTAINER_NAME)

# ============================================================================
# SECURITY SCANNING TARGETS
# ============================================================================

scan: build
	@echo "--- Scanning Docker Image for HIGH/CRITICAL vulnerabilities using Trivy ---"
	$(TRIVY_BIN) image --severity HIGH,CRITICAL $(FULL_IMAGE_NAME)
	@echo "--- Trivy scan completed ---"

vulnerability-scan: build
	@echo "$(SCAN) Scanning $(FULL_IMAGE_NAME) for vulnerabilities..."
	$(TRIVY_BIN) image --format table --scanners vuln,misconfig --severity MEDIUM,HIGH,CRITICAL $(FULL_IMAGE_NAME)
	@echo "$(CHECK) Vulnerability scan completed"

scan-docker: build
	@echo "--- Scanning with Docker Scout ---"
	docker scout quickview $(FULL_IMAGE_NAME)

# ============================================================================
# FALCO RUNTIME SECURITY TARGETS
# ============================================================================

falco-check:
	@echo "--- Checking Falco prerequisites ---"
	@where helm >nul 2>&1 || (echo "$(CROSS) Helm not found." && exit 1)
	@echo "$(CHECK) Helm is available"
	@kubectl version --client >nul 2>&1 || (echo "$(CROSS) ERROR: kubectl is required." && exit 1)
	@echo "$(CHECK) kubectl is available"
	@echo "--- Prerequisites check passed ---"

falco-install-clean: falco-check k8s-check
	@echo "--- Complete Falco Cleanup and Fresh Install ---"
	@echo "$(INFO) Removing previous Falco resources..."
	-helm uninstall falco -n $(FALCO_NAMESPACE) --wait 2>nul || echo "$(INFO) No helm release found"
	-kubectl delete namespace $(FALCO_NAMESPACE) --ignore-not-found=true --wait --timeout=60s
	@echo "$(INFO) Installing Falco (simple version)..."
	helm repo add falcosecurity https://falcosecurity.github.io/charts --force-update
	helm repo update
	helm upgrade --install falco falcosecurity/falco \
	   --namespace $(FALCO_NAMESPACE) \
	   --create-namespace \
	   --set falcosidekick.enabled=false \
	   --wait
	@echo "$(CHECK) Fresh Falco installation completed!"

falco-logs:
	@echo "--- Falco Security Events ---"
	kubectl logs -l app=falco -n $(FALCO_NAMESPACE) -c falco --tail=50 -f

falco-test:
	@echo "--- Testing Falco Detection ---"
	@echo "$(TEST) Generating test security event (cat /etc/shadow)..."
	kubectl run --rm -i --tty falco-test --image=alpine --restart=Never --namespace $(K8S_NAMESPACE) -- sh -c "cat /etc/shadow" || true
	@echo "$(INFO) Check logs with: make falco-logs"

runtime-security-test: falco-test
	@echo "--- Runtime Security Tests (Advanced) ---"
	@echo "$(TEST) Testing Falco detection capabilities against application namespace..."
	@echo "$(TEST) 1. Testing sensitive file access..."
	kubectl run --rm -i --tty test-file-access --image=alpine --restart=Never --namespace $(K8S_NAMESPACE) -- sh -c "ls -la /etc/passwd" || true
	@timeout /t 2 /nobreak >nul
	@echo "$(TEST) 2. Testing binary change..."
	kubectl run --rm -i --tty test-binary --image=alpine --restart=Never --namespace $(K8S_NAMESPACE) -- sh -c "touch /usr/bin/newfile" || true
	@echo.
	@echo "$(INFO) Check security events with: make falco-logs"

falco-clean:
	@echo "--- Removing Falco ---"
	helm uninstall falco -n $(FALCO_NAMESPACE) --wait 2>nul || echo "$(INFO) No Falco release found"
	kubectl delete namespace $(FALCO_NAMESPACE) --ignore-not-found=true --timeout=60s
	@echo "$(CHECK) Falco removed"

# ============================================================================
# KUBERNETES TARGETS (Deployment & Access)
# ============================================================================

k8s-check:
	@echo "--- Checking Kubernetes Access ---"
	kubectl config use-context $(KUBE_CONTEXT)
	kubectl cluster-info
	@echo "--- Kubernetes is ready! ---"

k8s-deploy: build k8s-check
	@echo "--- Deploying Go Secure Server to Kubernetes ---"
	kubectl create namespace $(K8S_NAMESPACE) --dry-run=client -o yaml | kubectl apply -f -
	kubectl apply -f k8s/ -n $(K8S_NAMESPACE)
	@echo "--- Waiting for deployment to be ready ---"
	kubectl wait --for=condition=ready pod -l app=go-secure-server -n $(K8S_NAMESPACE) --timeout=120s
	@echo "$(CHECK) Application Deployment completed! ---"

k8s-access:
	@echo "--- Starting Port Forward ---"
	@echo "$(INFO) Your Go Secure Server is now in K8s with Runtime Security!"
	@echo.
	@echo "$(INFO) AVAILABLE ENDPOINTS (http://localhost:8080/...):"
	@echo "$(INFO)   Health Check:    http://localhost:8080/"
	@echo "$(INFO)   Halloween Page:  http://localhost:8080/halloween"
	@echo "$(INFO)   Halloween API:   http://localhost:8080/api/halloween"
	@echo.
	@echo "$(INFO) Press Ctrl+C to stop port-forwarding"
	kubectl port-forward svc/go-secure-service 8080:80 -n $(K8S_NAMESPACE)

k8s-run: k8s-deploy k8s-access

k8s-clean: netpol-clean
	@echo "--- Cleaning Kubernetes Application Resources ---"
	kubectl delete -f k8s/ -n $(K8S_NAMESPACE) --ignore-not-found=true --wait
	kubectl delete namespace $(K8S_NAMESPACE) --ignore-not-found=true --timeout=60s
	@echo "$(CHECK) K8s application cleanup completed"

# ============================================================================
# NETWORK POLICY TARGETS (Zero Trust)
# ============================================================================

netpol-apply: k8s-check
	@echo "--- Applying Network Policy: $(NETPOL_FILE) ---"
	kubectl apply -f $(NETPOL_FILE) -n $(K8S_NAMESPACE)
	@echo "$(CHECK) NetworkPolicy applied successfully. (Zero Trust: Deny All, Allow 8080 and DNS)"

netpol-clean: k8s-check
	@echo "--- Removing Network Policy: $(NETPOL_FILE) ---"
	-kubectl delete -f $(NETPOL_FILE) -n $(K8S_NAMESPACE) --ignore-not-found=true
	@echo "$(CHECK) NetworkPolicy removed."

netpol-test: k8s-check
	@echo "--- Testing Network Policy Enforcement (Requires 'go-secure-server' deployment) ---"
	@echo "$(TEST) 1. Testing INGRESS (Ingress should be allowed on port 8080 only)"
	kubectl run test-client --rm -i --image=alpine/curl --restart=Never --namespace $(K8S_NAMESPACE) -- /bin/sh -c "curl -s --max-time 5 go-secure-service:8080/ || echo '$(CROSS) Ingress Test: FAILED (Unexpected)'"
	@echo "$(TEST) 2. Testing EGRESS (Egress to the external Internet should be blocked)"
	kubectl exec deployment/go-secure-server -n $(K8S_NAMESPACE) -- sh -c "echo 'Testing external connection...'; curl -s --max-time 3 https://google.com > /dev/null 2>&1 && echo '$(CROSS) Egress Test: FAILED (Unexpected - NetPol broken)' || echo '$(CHECK) Egress Test: PASSED (Blocked by NetPol)'"
	@echo "$(TEST) 3. Testing KubeDNS (DNS resolution must be allowed)"
	kubectl exec deployment/go-secure-server -n $(K8S_NAMESPACE) -- sh -c "nslookup google.com > /dev/null 2>&1 && echo '$(CHECK) DNS Test: PASSED (Allowed by NetPol)' || echo '$(CROSS) DNS Test: FAILED (NetPol blocking DNS)'"
	@echo "$(CHECK) Network Policy tests initiated. Review output manually."

# ============================================================================
# COMPREHENSIVE SECURITY PIPELINES
# ============================================================================

security-quick: build k8s-deploy
	@echo.
	@echo "$(DONE) QUICK SECURITY DEPLOYMENT COMPLETED!"
	@echo "  $(CHECK) Docker image built"
	@echo "  $(CHECK) K8s application deployed"
	@echo "  $(CHECK) Basic security context applied"
	@echo.
	@echo "  $(INFO) For full security stack run: make security-full"
	@echo "  $(GO) Quick access: make k8s-access"

security-full: scan falco-install-clean k8s-deploy netpol-apply
	@echo.
	@echo "$(DONE) FULL SECURITY STACK DEPLOYED!"
	@echo "  $(CHECK) Docker image built and scanned"
	@echo "  $(CHECK) Falco runtime security installed (clean install)"
	@echo "  $(CHECK) K8s application deployed with security context"
	@echo "  $(CHECK) Network Policy applied (Zero Trust)"
	@echo.
	@echo "$(SEC) SECURITY MONITORING ACTIVE!"
	@echo "  $(LOG) Run: make falco-logs    - View security events"
	@echo "  $(TEST) Run: make netpol-test   - Test Network Policy enforcement"
	@echo "  $(GO) Run: make k8s-access    - Access your application"

security-scan: vulnerability-scan k8s-security-audit
	@echo "--- Security Scanning Complete ---"

security-clean: k8s-clean falco-clean clean
	@echo "--- Complete Security Environment Cleanup ---"

# ============================================================================
# AUDIT & DEBUGGING TARGETS
# ============================================================================

k8s-security-audit:
	@echo "--- Kubernetes Security Audit ---"
	@echo "$(INFO) Checking K8s security configuration..."
	@echo.
	@echo "$(SEC) Pod Security Context (go-secure-server):"
	@-kubectl get pod -l app=go-secure-server -n $(K8S_NAMESPACE) -o jsonpath='{.items[*].spec.containers[*].securityContext}' 2>/dev/null && echo || echo "$(WARN) No security context found"
	@echo.
	@echo "$(INFO) Network Policy Check:"
	@-kubectl get netpol -n $(K8S_NAMESPACE) | findstr "go-server-default-deny" >nul && echo "$(CHECK) Network Policy is applied." || echo "$(CROSS) Network Policy is MISSING."
	@echo.
	@echo "$(CHECK) K8s security audit completed"

quick-debug:
	@echo "--- Quick Debug Information ---"
	@echo "$(INFO) Docker:"
	@-docker version --format 'Docker Version: {{.Server.Version}}' 2>/dev/null | findstr "Version" || echo "$(CROSS) Docker not available"
	@echo "$(INFO) Kubernetes:"
	@-kubectl version --short 2>/dev/null | findstr "Server" || echo "$(CROSS) Kubernetes not available"
	@echo "$(INFO) Falco:"
	@-kubectl get pods -n $(FALCO_NAMESPACE) 2>/dev/null | findstr "falco" && echo "$(CHECK) Falco: RUNNING" || echo "$(CROSS) Falco: NOT RUNNING"
	@echo "$(INFO) Application:"
	@-kubectl get pods -n $(K8S_NAMESPACE) 2>/dev/null | findstr "go-secure-server" && echo "$(CHECK) App: RUNNING" || echo "$(CROSS) App: NOT RUNNING"

security-report:
	@set REPORT_FILE=security-report-%DATE%_%TIME%.txt
	@echo "--- Generating Security Report ---"
	@echo SECURITY REPORT - %DATE% %TIME% > "%REPORT_FILE%"
	@echo ================================= >> "%REPORT_FILE%"
	@echo. >> "%REPORT_FILE%"
	@echo VULNERABILITY SCAN: >> "%REPORT_FILE%"
	$(TRIVY_BIN) image --severity HIGH,CRITICAL $(FULL_IMAGE_NAME) >> "%REPORT_FILE%" 2>&1 || echo Scan failed >> "%REPORT_FILE%"
	@echo. >> "%REPORT_FILE%"
	@echo KUBERNETES STATUS: >> "%REPORT_FILE%"
	@-kubectl get pods,svc -n $(K8S_NAMESPACE) >> "%REPORT_FILE%" 2>&1
	@echo. >> "%REPORT_FILE%"
	@echo FALCO STATUS: >> "%REPORT_FILE%"
	@-kubectl get pods -n $(FALCO_NAMESPACE) >> "%REPORT_FILE%" 2>&1
	@echo "$(INFO) Security report saved to: %REPORT_FILE%"

# ============================================================================
# INFORMATION & HELP TARGETS
# ============================================================================

info: help

help:
	@echo.
	@echo "GO K8S SECURITY LAB - COMPLETE MAKEFILE"
	@echo "==========================================="
	@echo.
	@echo "QUICK START & PIPELINES:"
	@echo "  make all              - Quick security deployment (build + k8s-deploy)"
	@echo "  make security-full    - Complete security stack (scan, falco, k8s, netpol)"
	@echo "  make security-clean   - Clean everything (k8s, falco, docker)"
	@echo.
	@echo "SECURITY COMMANDS:"
	@echo "  make scan             - Quick Trivy scan (HIGH/CRITICAL)"
	@echo "  make vulnerability-scan - Full vulnerability and misconfig scan"
	@echo "  make k8s-security-audit - Check K8s security context & NetPol status"
	@echo "  make falco-test       - Trigger a security event for Falco"
	@echo "  make netpol-test      - Test Network Policy ingress/egress"
	@echo "  make security-report  - Generate a file-based security report"
	@echo.
	@echo "K8S MANAGEMENT:"
	@echo "  make k8s-deploy       - Deploy app & NetPol to K8s"
	@echo "  make k8s-access       - Start port-forward to access app"
	@echo "  make falco-install-clean - Reinstall Falco securely"
	@echo "  make falco-logs       - Follow Falco security events"
	@echo.
	@echo "LOCAL DOCKER:"
	@echo "  make build            - Build Docker image"
	@echo "  make run              - Run Docker image locally"
	@echo "  make stop             - Stop local container"
	@echo.

.PHONY: all build run test stop clean logs status open-all deploy restart \
        scan vulnerability-scan scan-docker \
        falco-check falco-install-clean falco-logs falco-test runtime-security-test falco-clean \
        k8s-check k8s-deploy k8s-access k8s-run k8s-logs k8s-status k8s-clean \
        netpol-apply netpol-clean netpol-test \
        security-quick security-full security-scan security-clean info help \
        k8s-security-audit quick-debug security-report