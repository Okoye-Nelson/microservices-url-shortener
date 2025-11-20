# URL Shortener Microservices — Production-Ready Architecture

A cloud-native URL shortener built with **Golang microservices**, **Clean Architecture**, **Docker**, and **Kubernetes**, following modern DevOps and platform engineering practices. This repository is part of an ongoing effort to learn, improve, and extend scalable distributed systems by studying real-world patterns and infrastructure workflows.

---

### 🎯 Skills Demonstrated

- Go microservices with domain boundaries

- Clean Architecture implementation

- Docker multi-arch builds

- Kubernetes deployment patterns

- GitOps workflow management

- Cloud-native observability practices 

---

## 📦 Features at a Glance

<p align="center">

<img src="https://img.shields.io/badge/Go-Microservices-00ADD8?style=flat-square" />
<img src="https://img.shields.io/badge/Clean_Architecture-Yes-9333EA?style=flat-square" />
<img src="https://img.shields.io/badge/PostgreSQL-Database-4169E1?style=flat-square" />
<img src="https://img.shields.io/badge/Docker_&_Compose-Ready-2496ED?style=flat-square" />
<img src="https://img.shields.io/badge/Multi-arch_Build-Yes-0EA5E9?style=flat-square" />
<img src="https://img.shields.io/badge/Kubernetes-Manifests-326CE5?style=flat-square" />
<img src="https://img.shields.io/badge/GitOps-Portainer-0EA5E9?style=flat-square" />
<img src="https://img.shields.io/badge/CI-Automation-0F766E?style=flat-square" />

</p>

- **Independent Go microservices:** Link, Redirect, Stats  
- **Clean Architecture + modular domain layers** for maintainable code  
- **PostgreSQL** as persistent data store  
- **Docker & Docker Compose** for full local stack  
- **Multi-arch Docker builds** + automated Docker Hub pipelines  
- **Kubernetes manifests** for production deployment  
- **Optional GitOps workflow** with Portainer for automated cluster sync  
- **Built-in scripts** (`push-to-dockerhub.sh`) for CI/CD style automation  

---

## 🛠 Prerequisites

Ensure the following are installed:

- **Go 1.22+**
- **Docker** & **Docker Compose**
- **kubectl**
- **PostgreSQL**
- A Kubernetes cluster (local or cloud)
- Optional: **Portainer** for GitOps workflows

---
### **Local Development with Docker Compose**

```bash
# Start all services with hot reload
docker-compose up --build

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
# Run database migrations
psql -h localhost -U postgres -d urlshortener -f scripts/init.sql

# Start individual services
cd services/link-service && go run main.go
cd services/redirect-service && go run main.go  
cd services/stats-service && go run main.go
```

## 🐳 **Docker Hub Build & Push Process**

The `push-to-dockerhub.sh` script automates building and pushing all service images to Docker Hub.

### **Script Configuration**

```bash
# Default configuration (edit in script)
DOCKER_HUB_USERNAME="piyushsachdeva"          # Your Docker Hub username
IMAGE_TAG="${IMAGE_TAG:-latest}"         # Configurable via environment
BUILD_PLATFORM="linux/amd64,linux/arm64" # Multi-architecture support

# Services built:
- url-shortener-link      (services/link-service/Dockerfile)
- url-shortener-redirect  (services/redirect-service/Dockerfile)  
- url-shortener-stats     (services/stats-service/Dockerfile)
- url-shortener-frontend  (frontend/Dockerfile)
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

### **Option 1: Standard Kubernetes** (`k8s/base/`)

```bash
# Deploy to any Kubernetes cluster
kubectl apply -f k8s/base/

# Check deployment status
kubectl get pods -n url-shortener
kubectl get services -n url-shortener
kubectl get ingress -n url-shortener
```

### **Option 2: Portainer GitOps** (`k8s/gitopsportainer/`)

Complete GitOps deployment with Portainer management interface. **For detailed Portainer setup instructions, see:**

📖 **[k8s/gitopsportainer/README-GITOPS.md](k8s/gitopsportainer/README-GITOPS.md)**

**Quick Portainer Deployment:**

```bash
# Navigate to GitOps directory
cd k8s/gitopsportainer/

# Option A: Automatic deployment (recommended)
./deploy.sh

```

**Access after Portainer deployment:**
```bash
# Via Ingress (recommended)
kubectl get ingress -n url-shortener

# Via LoadBalancer
kubectl get svc -n url-shortener | grep LoadBalancer

# Via Port Forward (development)
kubectl port-forward svc/frontend 8080:80 -n url-shortener
```

## 🔗 **Related Documentation**

- **Portainer GitOps Setup**: [k8s/gitopsportainer/README-GITOPS.md](k8s/gitopsportainer/README-GITOPS.md)
- **API Documentation**: Available at `/api/docs` when services are running
- **Database Schema**: See `scripts/init.sql` for complete schema

