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

.PHONY: all scan scan-fast scan-full scan-docker build run stop clean test info logs status deploy open-all restart