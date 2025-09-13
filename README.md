# Shorty: A High-Performance URL Shortener

[![Go Report Card](https://goreportcard.com/badge/github.com/Kalmera74/Shorty)](https://goreportcard.com/report/github.com/Kalmera74/Shorty)
[![CI](https://github.com/Kalmera74/Shorty/actions/workflows/ci.yml/badge.svg)](https://github.com/Kalmera74/Shorty/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/docker-ready-blue.svg)](https://hub.docker.com/)

Shorty is an fully containerized URL shortener backend built with **Go**. It's designed to be fast, scalable, and maintainable. It features authentication & authorization, analytics, caching and message broking.

---

## Core Features

- **Secure User Management:** JWT-based authentication and authorization for multi-tenant usage.
- **Custom Short URLs:** Users can create personalized, branded short links.
- **Blazing-Fast Redirects:** Sub-millisecond response times for shortened URLs, with in-memory caching with Redis.
- **Asynchronous, High-Throughput Click Analytics:** Redirect analytics are processed asynchronously using RabbitMQ, ensuring the user-facing redirect path is never blocked.
- **Strict Input Validation:** All request payloads are validated with [`go-playground/validator`](https://github.com/go-playground/validator) to ensure data integrity and prevent malformed inputs.
- **Clean, Vertical Slice Architecture:** The codebase is organized by feature, making it highly modular, easy to test, and simple for teams to collaborate on.
- **API Rate Limiting:** Protects the API from abuse and ensures service stability.
- **Structured, Production-Ready Logging:** Uses [`zerolog`](https://github.com/rs/zerolog) for high-performance, structured (JSON) logging.
- **Developer-First Experience:** Auto-generated interactive API documentation via Swagger and a single-command setup with Docker or make.
- **Unit Tests**: Each key part of the project has unit tests to ensure they comply with business requirements throughout the development
- **Future-Proof**: The APIs are designed with an admin dashboard in mind, making it simple, fast, and efficient to set up a fully functional management interface.

---

## Quick Start

The entire project is containerized for a seamless, one-command setup.

### Prerequisites

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Running the Application

1.  Clone the repository:

    ```bash
    git clone https://github.com/Kalmera74/Shorty
    cd Shorty
    ```

2.  Bring the entire stack online:

    ```bash
    docker-compose up
    ```

    This command will build the Go application image and start containers for the API, the analytics worker, PostgreSQL, Redis, and RabbitMQ.

3.  Explore the API:
    Once the containers are running, the interactive Swagger API documentation is available at:
    **[http://localhost:8080/swagger](http://localhost:8080/swagger)**

---

## Architecture & Design Philosophy

This project was designed with performance, scalability, and maintainability in mind. Every technology and pattern was chosen to solve a specific problem.

### Tech Stack & Tooling

| Category                  | Technology                                                                                                  |
| :------------------------ | :---------------------------------------------------------------------------------------------------------- |
| **Core Backend**          | **Go**, **Fiber** (High-Performance Web Framework), **Gorm** (ORM), **Validator** (go-playground/validator) |
| **Data & Infrastructure** | **PostgreSQL** (Primary Datastore), **Redis** (Caching), **RabbitMQ** (Message Broker)                      |
| **DevOps & Tooling**      | **Docker & Docker Compose** (Containerization), **GitHub Actions** (CI/CD), **Makefile** (Build Automation) |
| **Observability & Docs**  | **Zerolog** (Structured Logging), **Swagger/OpenAPI** (API Documentation)                                   |

### System Flow Diagram

The system is split into two primary workflows: the read-optimized redirect path and the write-intensive analytics path.

1.  **URL Redirect (The Critical Path):**

    - A request for a short URL hits the API.
    - The service first checks **Redis** using a **cache-aside pattern**. A cache hit results in an immediate redirect, avoiding any database interaction.
    - On a cache miss, the service queries **PostgreSQL**, populates the Redis cache for future requests, and then redirects.
    - Crucially, a "click" event is published to **RabbitMQ** without waiting for a response. **The redirect is never blocked by analytics processing.**

2.  **Analytics Processing (The Decoupled Path):**
    - The `analytics-worker` is a separate, dedicated Go service.
    - It consumes messages from the RabbitMQ queue in the background.
    - This worker is solely responsible for parsing click data and persisting it to PostgreSQL. This decoupling means the two services can be scaled independently.

### Key Design Decisions

- **Clean Architecture with Vertical Slices:** The project structure under `internal/features` is deliberate. Each feature (e.g., `user`, `shortener`) is a self-contained module. This promotes **high cohesion and low coupling**, making the codebase easy to navigate, test, and extend, and allows development teams to work on features in parallel with minimal friction.

- **Validation First:** All request objects are validated with [`go-playground/validator`](https://github.com/go-playground/validator) to ensure strict type and format guarantees at the API boundary. This prevents invalid data from propagating into the system and reduces runtime errors.

- **Event-Driven & Asynchronous:** By using RabbitMQ, we decouple the critical, user-facing redirect logic from the non-critical, background analytics work. This drastically improves perceived performance and resilience. The redirect service remains lightweight and fast, even under heavy load.

- **Stateless for Scalability:** Using JWTs for authentication means the application is stateless. Any API instance can serve any authenticated request. This is crucial for **horizontal scalability**, allowing us to run multiple instances behind a load balancer without needing sticky sessions.

- **Role-Based Authorization (RBAC):** Shorty uses simple, string-based roles (e.g., `admin`, `user`) for authorization. This lightweight approach keeps the system easy to manage while still allowing future extension to fine-grained policies if needed.

- **API Rate Limiting:** Rate limiting is applied at the API gateway level to prevent abuse, protect from brute-force attacks, and maintain service stability under high load.

- **Structured Logging:** Logging is done with [`zerolog`](https://github.com/rs/zerolog), producing structured JSON logs. This enables seamless integration with monitoring/observability tools like ELK or Grafana Loki while remaining high-performance in production environments.

- **Configuration Management:** All configuration (database URLs, JWT secrets, etc.) is managed via environment variables and loaded with [`Godotenv`](https://github.com/joho/godotenv). This follows the twelve-factor app methodology, ensuring no secrets are hardcoded and the application is portable across environments.

---

## Project Structure

```
├── cmd/                # Entrypoints for our binaries
│   ├── api/            # Main API application
│   └── analytics-worker/ # Background worker for processing clicks
├── internal/           # Private application code, not for export
│   ├── apperrors/      # Custom application-specific errors
│   ├── features/       # Core business logic, organized by feature (Vertical Slices)
│   │   ├── shortener/  # All code related to URL shortening
│   │   |── user/       # All code related to user management
│   ├── middleware/     # Shared HTTP middleware (e.g., auth)
│   └── ...
├── pkg/                # Public library code, shareable with other projects
│   ├── auth/           # JWT generation and validation logic
│   ├── cache/          # Redis client wrapper
│   └── ...
├── docs/               # Swagger/OpenAPI documentation files
├── Dockerfile          # Defines the Go application container
└── docker-compose.yml  # Defines and orchestrates all services
```
