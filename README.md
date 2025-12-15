# URL Shortener Microservices — Production-Ready Architecture

A cloud-native URL shortener built with **Golang microservices**, **Clean Architecture**, **Docker**, and **Kubernetes**, following modern DevOps and platform engineering practices. This repository is part of an ongoing effort to learn, improve, and extend scalable distributed systems by studying real-world patterns and infrastructure workflows.

---

### 🎯 Skills Demonstrated

- Go microservices with domain boundaries

- Clean Architecture implementation

- Docker multi-arch builds

- Kubernetes modular deployment patterns

- Cloud-native observability practices 

---

## 📦 Features at a Glance

<p align="center">

<img src="https://img.shields.io/badge/Go-Microservices-00ADD8?style=flat-square" />
<img src="https://img.shields.io/badge/Clean_Architecture-Yes-9333EA?style=flat-square" />
<img src="https://img.shields.io/badge/PostgreSQL-Database-4169E1?style=flat-square" />
<img src="https://img.shields.io/badge/Docker_&_Compose-Ready-2496ED?style=flat-square" />
<img src="https://img.shields.io/badge/Multi-arch_Build-Yes-0EA5E9?style=flat-square" />
<img src="https://img.shields.io/badge/CI-Automation-0F766E?style=flat-square" />

</p>

- **Independent Go microservices:** Link, Redirect, Stats  
- **Clean Architecture + modular domain layers** for maintainable code  
- **PostgreSQL** as persistent data store  
- **Docker & Docker Compose** for full local stack  
- **Multi-arch Docker builds** + automated Docker Hub pipelines  
- **Built-in scripts** (`push-to-dockerhub.sh`) for CI/CD style automation  

---

## 🛠 Prerequisites

Ensure the following are installed:

- **Go 1.22+**
- **Docker** & **Docker Compose**
- **kubectl**
- A Kubernetes cluster (local or cloud)

---
### **Local Environment Setup**

Before running the services, you need to create a local environment file.

```bash
cp .env.example .env
```
---
### **Local Development with Docker Compose**

```bash
# Start all services with hot reload
docker compose up --build

# Access points:
# Frontend: http://localhost:3000
# API Gateway: http://localhost:8080  
# Link service :8001 
# Redirect :8002
# Stats service :8003
# PostgreSQL: localhost:5432
```

### **Manual Service Development**

```bash
# The database is initialized automatically by Docker Compose using scripts/init.sql
# To run manually: psql -h localhost -U postgres -d urlshortener -f scripts/init.sql

# Start individual services
cd services/link-service && go run main.go
cd services/redirect-service && go run main.go  
cd services/stats-service && go run main.go
```

## 🐳 **Docker Hub Build & Push Process**

The `push-to-dockerhub.sh` script automates building and pushing all service images to Docker Hub.

### **Script Configuration**

```bash
# Default configuration (can be overridden by environment variables)
DOCKER_HUB_USERNAME="<your-dockerhub-username>" # Your Docker Hub username
IMAGE_TAG="${IMAGE_TAG:-latest}"         # Configurable via environment
BUILD_PLATFORM="linux/amd64,linux/arm64" # Multi-architecture support

# Services built:
- link-service      (services/link-service/Dockerfile)
- redirect-service  (services/redirect-service/Dockerfile)  
- stats-service     (services/stats-service/Dockerfile)
- frontend          (frontend/Dockerfile)
```

### **Complete Build & Push Workflow**

```bash
# 1. Build and push all images (recommended)
./push-to-dockerhub.sh deploy

# 2. Using custom tag
IMAGE_TAG=v1.2.0 ./push-to-dockerhub.sh deploy

# 3. Update Kubernetes manifests with new image tags
./push-to-dockerhub.sh update-k8s

# 4. Verify images on Docker Hub
./push-to-dockerhub.sh verify
```

### **Individual Operations**

```bash
# Build images locally only
./push-to-dockerhub.sh build

# Push existing images to Docker Hub
./push-to-dockerhub.sh login  # First-time setup
./push-to-dockerhub.sh push

# List local images
./push-to-dockerhub.sh list

# Clean up local images
./push-to-dockerhub.sh cleanup

# Get Docker Hub repository info
./push-to-dockerhub.sh info
```

### **Environment Variables**

```bash
# Custom image tag
export IMAGE_TAG="v2.1.0"

# Custom Docker Hub username (overrides script default)
export DOCKER_HUB_USERNAME="myusername"

# Custom build platform
export BUILD_PLATFORM="linux/amd64"
```

## ⚓ **Kubernetes Deployment Options**

### **Standard Kubernetes Deployment**

All Kubernetes resources are defined in a single manifest file. For a production environment, it is highly recommended to split this into individual resource files (e.g., in a `k8s/base/` directory).
```bash
# Deploy to any Kubernetes cluster
kubectl apply -f k8s/base/

# Check deployment status
kubectl get pods
kubectl get services
```

## 🔗 **Related Documentation**

- **API Documentation**: Available at `/api/docs` when services are running
- **Database Schema**: See `scripts/init.sql` for complete schema
