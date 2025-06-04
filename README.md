# Finplan Microservices Architecture

This repository contains a set of Go microservices orchestrated with Docker Compose.

## Microservices

- **API Gateway** - Exposes HTTP endpoints and routes requests to the underlying services.
- **Auth Service** - Manages user authentication and issues JWT tokens.
- **User Service** - Handles user profiles and related operations.
- **Transaction Service** - Stores and retrieves user transactions and categories.

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/)
- Go 1.20+ (only required for running tests without Docker)

Create environment files before starting the stack:

```bash
cp services/<service>/.env.example services/<service>/.env
```

Repeat the command for each service (`api-gateway`, `auth-service`, `user-service`, `transaction-service`). Edit the copied `.env` files if you need to override any defaults.

## Running with Docker Compose

Start all services along with PostgreSQL:

```bash
docker-compose up --build
```

Services will be available on the ports defined in `docker-compose.yml`.

## Running Tests

To run unit tests for all services locally:

```bash
go test ./...
```


