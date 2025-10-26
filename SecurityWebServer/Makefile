PROJECT_NAME := gok8ssecuritylab
IMAGE_NAME := $(PROJECT_NAME)/go-secure-server
IMAGE_VERSION := latest
FULL_IMAGE_NAME := $(IMAGE_NAME):$(IMAGE_VERSION)
CONTAINER_NAME := go-secure-server-test

all: build

build:
	@echo "--- Building Docker Image: $(FULL_IMAGE_NAME) ---"
	docker build --pull -t $(FULL_IMAGE_NAME) .
	@echo "--- Successfully built: $(FULL_IMAGE_NAME) ---"

run: stop build
	@echo "--- Running container on http://localhost:8080 ---"
	docker run -d --rm --name $(CONTAINER_NAME) -p 8080:8080 $(FULL_IMAGE_NAME)
	@echo "--- Container $(CONTAINER_NAME) started successfully ---"

stop:
	@echo "--- Stopping container: $(CONTAINER_NAME) ---"
	-docker stop $(CONTAINER_NAME) 2>nul && echo "--- Container stopped ---" || echo "--- Container was not running ---"

clean: stop
	@echo "--- Deleting image: $(FULL_IMAGE_NAME) ---"
	-docker rmi $(FULL_IMAGE_NAME) 2>nul && echo "--- Image deleted ---" || echo "--- Image not found ---"

.PHONY: all build run stop clean