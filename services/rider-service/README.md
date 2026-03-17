# Rider Service

Premium rider-delivery backend for large restaurant operations. The service is implemented as a clean modular monolith in Go so teams can ship quickly now and split into microservices later without rewriting core domain boundaries.

## What is in this scaffold

- REST API with JWT access and refresh token flow
- Rider profile, shift, assignment, delivery lifecycle, OTP, wallet, payout, rating, notification, support, and admin modules
- Structured logging, request IDs, RBAC, rate limiting, validation, and standardized error responses
- PostgreSQL and Mongo bootstrap hooks with an in-memory repository for local runnable development and tests
- Docker, docker-compose, `.env` template, SQL migrations, OpenAPI doc, and unit tests for core flows

## Project layout

- `cmd/api`: service entrypoint
- `config`: environment-driven config loading
- `internal/app`: dependency container and datastore bootstrap
- `internal/bootstrap`: HTTP server and router wiring
- `internal/domain`: roles, enums, permissions
- `internal/dto`: request DTOs
- `internal/handler`: REST handlers
- `internal/middleware`: auth, RBAC, rate limiting, request tracing, logging
- `internal/model`: domain models
- `internal/repository`: in-memory repository and seeded demo data
- `internal/service`: business logic
- `pkg`: reusable auth, validation, logging, pagination, response helpers
- `docs`: architecture, schema design, OpenAPI, error catalog
- `migrations`: PostgreSQL DDL and seed SQL

## Local run

```bash
cp .env.example .env
go test ./...
go run ./cmd/api
```

## Seed users for local testing

- Rider: `rider@rider.local` / `Rider@123`
- Dispatcher: `dispatch@rider.local` / `Dispatch@123`
- Super Admin: `admin@rider.local` / `Admin@123`

## Production direction

The current repository layer is intentionally in-memory so the service is runnable in this workspace without live databases. The schema, migration, API, and domain contracts are designed for PostgreSQL and MongoDB-backed repositories, and the code is structured so those repositories can replace the in-memory store without changing handler contracts.
