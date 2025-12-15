# Technical Documentation: URL Shortener Microservices

**Version:** 1.0
**Date:** December 15, 2025

## 1. Project Overview

This document outlines the architecture, setup, and operational procedures for the URL Shortener Microservices project. The system is a cloud-native application designed to shorten long URLs, track redirection analytics, and provide a modern user interface for management. It is built using Golang, following Clean Architecture principles, and is designed for deployment on Kubernetes.

### Core Features
- **URL Shortening**: Creates short, unique identifiers for long URLs.
- **Redirection**: Handles redirection from short URLs to their original destination.
- **Analytics**: Asynchronously tracks click data for each short URL.
- **Management UI**: A responsive frontend for creating links and viewing statistics.
- **Notifications**: An extensible service for sending notifications (e.g., to Slack).

## 2. System Architecture

The project follows a microservices architecture, promoting separation of concerns, independent scalability, and maintainability.

### 2.1. Application Services

-   **Link Service**: The core service responsible for the business logic of creating, listing, and deleting short links. It validates incoming URLs and generates unique short IDs.
-   **Redirect Service**: A high-performance service focused solely on handling incoming short URL requests, looking up the original URL, and issuing an HTTP 302 redirect. It also publishes a "click" event for analytics.
-   **Stats Service**: Consumes "click" events and aggregates statistics. It provides the data for the analytics dashboard on the frontend.
-   **Notification Service**: A consumer service that listens for specific events (e.g., new link creation) and sends notifications to external systems like Slack.
-   **Frontend**: A vanilla JavaScript single-page application that provides the user interface. It is served as a static asset and communicates with the backend via the API Gateway.

### 2.2. Infrastructure Components

-   **Nginx API Gateway**: The single entry point for all external traffic. It routes requests to the appropriate backend service, handles rate limiting, and manages CORS policies.
-   **PostgreSQL**: The primary relational database used for persistent storage of links and their associated statistics. It is deployed as a `StatefulSet` in Kubernetes for data stability.
-   **Redis**: An in-memory cache used by the `link-service` to speed up lookups and reduce database load.
-   **RabbitMQ**: A message broker used for asynchronous communication between services. For example, the `redirect-service` publishes a message to a queue, which the `stats-service` and `notification-service` consume. This decouples the click-tracking process from the user-facing redirect, ensuring low latency.

## 3. Local Development Environment

The entire stack can be run locally using Docker Compose, providing a development environment that closely mirrors production.

### 3.1. Prerequisites
-   Go (1.22+)
-   Docker & Docker Compose
-   A code editor (e.g., VS Code)

### 3.2. Initial Setup

1.  **Clone the Repository**:
    ```bash
    git clone <your-repository-url>
    cd microservices-url-shortener
    ```

2.  **Create Environment File**: The system uses an `.env` file for managing secrets and local configuration. Copy the example template to create your local version.
    ```bash
    cp .env.example .env
    ```
    The default values in this file are suitable for the local Docker Compose setup.

### 3.3. Running the Stack

-   **Start all services**: This command builds the Docker images for each microservice and starts all containers defined in `docker-compose.yml`.
    ```bash
    docker compose up --build
    ```

-   **Access Points**: Once running, the services are available at the following local addresses:
    -   **Frontend UI**: `http://localhost:3000`
    -   **API Gateway**: `http://localhost:8080`
    -   **RabbitMQ Management**: `http://localhost:15672` (user: `admin`, pass: `admin`)
    -   **PostgreSQL**: `localhost:5432`
    -   **Redis**: `localhost:6379`

### 3.4. Stopping the Stack
-   To stop all running containers, press `Ctrl+C` in the terminal where `docker compose` is running.
-   To stop and remove the containers, run:
    ```bash
    docker compose down
    ```

## 4. Kubernetes Deployment

The project is designed for production deployment on Kubernetes. The manifests are organized in the `k8s/base/` directory, following a modular, one-file-per-component structure.

### 4.1. Manifest Structure

The `k8s/base/` directory contains individual YAML files for each Kubernetes resource, such as:
-   `00-secrets.yml`: Manages all application secrets.
-   `01-postgres.yml`: Defines the database `StatefulSet`, `Service`, and `ConfigMap`.
-   `09-link-service.yml`: Defines the `Deployment` and `Service` for the link-service.
-   ...and so on for all other components.

This structure is a best practice that allows for granular control, easier debugging, and compatibility with GitOps tools.

### 4.2. Deployment Procedure

1.  **Prerequisites**:
    -   A running Kubernetes cluster.
    -   `kubectl` configured to connect to your cluster.

2.  **Apply Manifests**: Deploy the entire application stack by applying all manifests from the `k8s/base/` directory.
    ```bash
    kubectl apply -f k8s/base/
    ```
    This command will create all the necessary Deployments, Services, Secrets, and other resources in the `default` namespace.

3.  **Verify Deployment**: Check the status of the pods and services.
    ```bash
    kubectl get pods
    kubectl get services
    ```
    Look for the `nginx-gateway` service of type `LoadBalancer` to find the external IP address for accessing the application.


## 5. API Endpoints

All API endpoints are exposed through the Nginx gateway.

| Method | Endpoint                 | Service            | Description                               |
| :----- | :----------------------- | :----------------- | :---------------------------------------- |
| `PUT`  | `/api/generate`          | `link-service`     | Creates a new short URL.                  |
| `GET`  | `/api/links`             | `link-service`     | Retrieves a list of all created links.    |
| `DELETE`| `/api/delete/{id}`       | `link-service`     | Deletes a specific short URL by its ID.   |
| `GET`  | `/api/stats`             | `stats-service`    | Retrieves aggregated analytics data.      |
| `GET`  | `/{shortId}`             | `redirect-service` | Redirects to the original URL.            |
| `GET`  | `/health`                | `nginx-gateway`    | A simple health check for the gateway.    |