# Go Kubernetes Security Lab

This project is a hands-on lab demonstrating how to build and deploy a secure Go application on Kubernetes, following modern security best practices.

## Overview

The repository contains:
1.  A simple Go web server.
2.  A `Dockerfile` for containerizing the application using security-focused principles.
3.  Kubernetes manifests (`Deployment` and `Service`) to deploy the application securely.

The goal is to provide a practical, working example of securing a containerized workload in a Kubernetes environment.

## Security Best Practices Demonstrated

This lab highlights several key security concepts at different layers of the stack.

### 1. Secure Dockerfile

The `Dockerfile` is designed to create a minimal and secure container image.

-   **Multi-Stage Builds**: We use a `builder` stage to compile the Go application, and a final, separate stage for the production image. This ensures that build tools and source code are not included in the final image, reducing its size and attack surface.
-   **Minimal Base Image (`distroless`)**: The final image is based on `gcr.io/distroless/static-debian11`, which contains only the application and its runtime dependencies. It does not include a shell, package manager, or other utilities that could be exploited.
-   **Non-Root User**: The container runs the application as a `nonroot` user (`UID 65532`). This is a critical security measure to limit the potential damage if an attacker gains code execution inside the container.

### 2. Hardened Kubernetes Deployment

The `k8s/deployment.yaml` manifest includes a strict `securityContext` to lock down the pod's runtime privileges.

-   **Run as Non-Root**: The pod is forced to run as a non-root user (`runAsNonRoot: true`), with a specific user and group ID that matches the one in the `Dockerfile`.
-   **Read-Only Root Filesystem**: The container's root filesystem is mounted as read-only (`readOnlyRootFilesystem: true`). This prevents an attacker from modifying the application binaries or system files. A `/tmp` directory is added as a writable `emptyDir` volume for temporary application needs.
-   **Disable Privilege Escalation**: We prevent the process from gaining more privileges than its parent (`allowPrivilegeEscalation: false`).
-   **Drop All Capabilities**: All Linux capabilities are dropped (`drop: ["ALL"]`). This follows the principle of least privilege, ensuring the container has no special kernel permissions.

### 3. Secure Configuration and Secrets Management

-   **ConfigMaps for Configuration**: Non-sensitive configuration is injected into the application via environment variables using a `ConfigMap`.
-   **Secrets for Sensitive Data**: Sensitive data (like API keys or passwords) is mounted into the pod as a volume from a Kubernetes `Secret`. This is more secure than using environment variables for secrets, as it prevents accidental exposure through logs or debugging endpoints.

## How to Run the Lab

### Prerequisites

-   Go (1.18+)
-   Docker
-   A Kubernetes cluster (e.g., Minikube, Kind, Docker Desktop)
-   `kubectl` configured to talk to your cluster

### 1. Build and Load the Docker Image

Build the container image and load it into your local cluster's image registry.

```sh
# For Minikube
eval $(minikube docker-env)

# For Kind
# kind load docker-image gosec-lab:latest

# Build the image
docker build -t gosec-lab:latest .
```

### 2. Create Kubernetes Resources

First, create the `ConfigMap` and `Secret` that the application depends on.

```sh
# Create a ConfigMap for non-sensitive configuration
kubectl create configmap my-configmap --from-literal=CONFIG_VALUE="Hello from ConfigMap!"

# Create a Secret for sensitive data
kubectl create secret generic my-secret --from-literal=secret-key="s3cr3t-p@ssw0rd"
```

### 3. Deploy the Application

Apply the Kubernetes manifests to deploy the application.

```sh
kubectl apply -f k8s/
```

### 4. Access the Application

Forward a local port to the service to access it from your machine.

```sh
kubectl port-forward svc/gosec-lab 8080:80
```

Now you can test the different endpoints:

```sh
# Health check
curl http://localhost:8080/healthz

# Get config from ConfigMap
curl http://localhost:8080/config

# Get secret from mounted Secret volume
curl http://localhost:8080/secret
```

### 5. Clean Up

```sh
kubectl delete -f k8s/
kubectl delete secret my-secret
kubectl delete configmap my-configmap
```