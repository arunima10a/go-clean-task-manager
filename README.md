# Go Clean Task Manager 

A production-grade distributed Task Management system built using **Clean Architecture (Onion Pattern)**. This project demonstrates high-performance Go engineering, multi-protocol support, and resilient infrastructure.

##  Architectural Highlights
- **Hexagonal/Clean Architecture:** Strict separation of Domain logic from Infrastructure.
- **Transport Agnostic:** Simultaneously serves **REST (Gin)** and **gRPC (Protobuf)** via a unified UseCase layer.
- **Messaging:** Event-driven background processing using **RabbitMQ**.
- **Performance:** Multi-level storage with **PostgreSQL** persistence and **Redis** caching.
- **Observability:** Full instrumentation with **Prometheus** metrics and structured **JSON logging (Zerolog)**.

##  Features
- **Security:** JWT-based Authentication with custom Middleware (REST) and Interceptors (gRPC).
- **Concurrency:** Graceful shutdown coordination using `sync.WaitGroup` to protect in-flight background workers.
- **Resilience:** Automatic connection retries for Postgres and RabbitMQ; Atomic database transactions via Context.
- **Quality:** 100% linting compliance using `golangci-lint` and automated CI via GitHub Actions.

##  Getting Started
1. **Prerequisites:** Docker & Docker Compose.
2. **Launch:** `make docker-up`
3. **API Docs:** `http://localhost:8080/swagger/index.html`
4. **Monitoring:** `http://localhost:9090` (Prometheus)