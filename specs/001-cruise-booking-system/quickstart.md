# Quickstart: Cruise Booking System

## Prerequisites

*   **Go**: 1.26 or later
*   **Node.js**: v20.x or later
*   **Docker & Docker Compose**: For running DB, Redis, MinIO, Meilisearch

## Infrastructure Setup

Start the dependencies using Docker Compose:

```bash
docker-compose up -d postgres redis minio meilisearch nats
```

## Backend (Go)

1.  Navigate to backend directory:
    ```bash
    cd backend
    ```
2.  Install dependencies:
    ```bash
    go mod tidy
    ```
3.  Run migrations (auto-migrate enabled in dev):
    ```bash
    go run cmd/migrate/main.go
    ```
4.  Start server:
    ```bash
    go run cmd/server/main.go
    ```
    API will be available at `http://localhost:8080`.

## Frontend - Admin (Nuxt 4)

1.  Navigate to admin directory:
    ```bash
    cd admin
    ```
2.  Install dependencies:
    ```bash
    npm install
    ```
3.  Start dev server:
    ```bash
    npm run dev
    ```
    Admin panel at `http://localhost:3000`.

## Frontend - Web (Nuxt 4)

1.  Navigate to web directory:
    ```bash
    cd web
    ```
2.  Install & Start:
    ```bash
    npm install && npm run dev
    ```
    Public site at `http://localhost:3001`.

## Environment Variables

Copy `.env.example` to `.env` in `backend/`:

```env
DB_DSN="host=localhost user=postgres password=postgres dbname=cruise_booking port=5432 sslmode=disable"
REDIS_ADDR="localhost:6379"
MINIO_ENDPOINT="localhost:9000"
JWT_SECRET="change_me_in_prod"
```
