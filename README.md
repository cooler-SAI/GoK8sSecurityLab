# Go Kubernetes Security Lab

This project is a hands-on lab demonstrating how to build and deploy a secure Go application on Kubernetes, following modern security best practices.

## Overview

The repository contains:
* A simple Go web server.
* A `Dockerfile` for containerizing the application using security-focused principles.
* Kubernetes manifests (`Deployment` and `Service`) to deploy the application securely.

The goal is to provide a practical, working example of securing a containerized workload in a Kubernetes environment.

---

## Security Best Practices Demonstrated

### 1. Secure Dockerfile
* **Multi-Stage Builds:** Uses a builder stage to compile the Go application and a final stage for production to reduce attack surface.
* **Minimal Base Image (distroless):** Uses `gcr.io/distroless/static-debian11`, which contains no shell or package manager.
* **Non-Root User:** Runs as UID `65532` to limit potential damage from code execution exploits.

### 2. Hardened Kubernetes Deployment
The `k8s/deployment.yaml` uses a strict `securityContext`:
* **runAsNonRoot: true**: Forces the container to run without root privileges.
* **readOnlyRootFilesystem: true**: Prevents modification of system files; uses `emptyDir` for `/tmp`.
* **allowPrivilegeEscalation: false**: Prevents the process from gaining extra permissions.
* **capabilities: drop: ["ALL"]**: Removes all kernel-level privileges.

### 3. Secure Configuration
* **ConfigMaps:** Used for non-sensitive environment variables.
* **Secrets:** Used for sensitive data, mounted as volumes rather than environment variables to prevent accidental logging exposure.

---

## How to Run the Lab

### Prerequisites
* Go (1.18+)
* Docker & Kubectl
* A Kubernetes cluster (Minikube, Kind, or Docker Desktop)

### 1. Build and Load the Docker Image
```bash
# For Minikube
eval $(minikube docker-env)

# Build the image
docker build -t gosec-lab:latest .

# For Kind (uncomment if using Kind)
# kind load docker-image gosec-lab:latest

```

### 2. Create Kubernetes Resources
Create the configuration and security objects the application requires.

```bash
# Create a ConfigMap for non-sensitive configuration
kubectl create configmap my-configmap --from-literal=CONFIG_VALUE="Hello from ConfigMap!"

# Create a Secret for sensitive data
kubectl create secret generic my-secret --from-literal=secret-key="s3cr3t-p@ssw0rd"
```

### 3. Deploy the Application
Apply the manifests to the cluster.

```bash
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
```

### 4. Access and Test
Forward the service port to your local machine.

```bash
kubectl port-forward svc/gosec-lab 8080:80
```

#### Test the endpoints in a new terminal:
```bash
# Health check
curl http://localhost:8080/healthz

# Get config from ConfigMap
curl http://localhost:8080/config

# Get secret from mounted Secret volume
curl http://localhost:8080/secret
```

### 5. Clean Up
Remove all resources created during this lab.

```bash
kubectl delete -f k8s/
kubectl delete secret my-secret
kubectl delete configmap my-configmap
```