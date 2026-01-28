# Go Kubernetes Security Lab

This project is a hands-on lab demonstrating how to build and deploy a secure Go application on Kubernetes, following modern security best practices.

## Overview

This repository is a comprehensive security lab that provides a practical, hands-on experience for securing a containerized workload in a Kubernetes environment. It includes multiple services that work together to demonstrate a variety of security concepts, from application-level security to infrastructure-level enforcement.

## Project Components

This lab is composed of three main services:

*   **`websecure`**: A Go web server that demonstrates application-level security best practices. It includes features like JWT-based authentication, rate limiting, security headers, and protection against common vulnerabilities like XSS. It also includes a vulnerable endpoint for demonstrating the importance of secure coding practices.

*   **`emuserver`**: A chaos engineering tool designed to simulate an unstable or unreliable service. It can be configured to introduce random delays and errors into its responses, allowing you to test the resilience and fault tolerance of your applications.

*   **`sentinel`**: A Kubernetes admission webhook that enforces security policies on pods before they are deployed to the cluster. It acts as a gatekeeper, preventing pods that violate predefined security rules (such as running as a privileged user) from being scheduled.

## Security Features

This lab demonstrates a wide range of security features, including:

*   **Application-Level Security**:
    *   **JWT Authentication**: Securely manage user sessions and protect endpoints with JSON Web Tokens.
    *   **Rate Limiting**: Protect your services from abuse and denial-of-service attacks.
    *   **Security Headers**: Harden your application against common web vulnerabilities like clickjacking and cross-site scripting (XSS).
    *   **Role-Based Access Control (RBAC)**: Enforce different levels of access for different users.

*   **Infrastructure-Level Security**:
    *   **Kubernetes Admission Webhooks**: Enforce custom security policies on your cluster.
    *   **Hardened Dockerfiles**: Build minimal, secure container images using multi-stage builds and non-root users.
    *   **Secure Kubernetes Deployments**: Configure your deployments with a strict security context to limit the blast radius of a potential compromise.
    *   **Network Policies**: Isolate your services and control the flow of traffic between them.

## How to Run the Lab

### Prerequisites

*   Go (1.18+)
*   Docker & Kubectl
*   A Kubernetes cluster (Minikube, Kind, or Docker Desktop)

### 1. Build and Deploy the Services

Each service can be built and deployed independently. Refer to the `README.md` file in each service's directory for detailed instructions.

### 2. Explore the Security Features

Once the services are deployed, you can explore the various security features they demonstrate. For example, you can:

*   Attempt to deploy a privileged pod and see it get blocked by `sentinel`.
*   Use `emuserver` to test the resilience of `websecure`.
*   Explore the different authentication and authorization features of `websecure`.

## Contributing info

Contributions are welcome! If you have any ideas for new features or improvements, please open an issue or submit a pull request.