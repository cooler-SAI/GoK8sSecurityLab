PROJECT_NAME := gok8ssecuritylab
IMAGE_NAME := $(PROJECT_NAME)/go-secure-server
IMAGE_VERSION := latest
FULL_IMAGE_NAME := $(IMAGE_NAME):$(IMAGE_VERSION)

all: build

build:
	@echo "--- Building Docker Image: $(FULL_IMAGE_NAME) ---"
	docker build --pull -t $(FULL_IMAGE_NAME) .
	@echo "--- Successfully built: $(FULL_IMAGE_NAME) ---"

run: build
	@echo "--- Running container on http://localhost:8080 ---"
	docker run -d --rm --name go-secure-server-test -p 8080:8080 $(FULL_IMAGE_NAME)

stop:
	@echo "--- Stopping container: go-secure-server-test ---"
	-docker stop go-secure-server-test

clean: stop
	@echo "--- Deleting image: $(FULL_IMAGE_NAME) ---"
	-docker rmi $(FULL_IMAGE_NAME)

.PHONY: all build run stop clean