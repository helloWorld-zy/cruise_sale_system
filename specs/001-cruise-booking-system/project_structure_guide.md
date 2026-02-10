# Project Structure Guide: Cruise Booking System

**Feature Branch**: `001-cruise-booking-system`
**Date**: 2026-02-10

This document outlines the file structure established in Phase 1-11 for the Cruise Booking System. Use this as a reference for locating modules and understanding the architectural layout.

## Root Directory

*   `.github/workflows/ci.yml`: CI/CD pipeline configuration (GitHub Actions) enforcing 100% test coverage.
*   `docker-compose.yml`: Local development infrastructure (Postgres, Redis, MinIO, Meilisearch, NATS).
*   `.gitignore`: Global ignore rules.

## Backend (Go 1.26 + Gin)
Path: `/backend`

*   `cmd/server/`: Entry point for the application (`main.go`).
*   `internal/`: Private application code.
    *   `api/`: HTTP Handlers and Router configuration (Gin).
        *   `v1/`: API Version 1.
            *   `admin/`: Admin Handlers.
                *   `cruise.go`: Admin Cruise CRUD.
                *   `voyage.go`: Admin Voyage/Inventory.
                *   `order.go`: Admin Order List/Detail/Upload.
            *   `cruise.go`: Cruise Listing/Detail handlers.
            *   `order.go`: Booking & User Order handlers.
            *   `price_trend.go`: Price Trend Analysis.
            *   `recommendation.go`: AI Recommendations.
            *   `excursion.go`: Shore Excursions.
    *   `core/`: Business logic and Service layer.
        *   `cruise_service.go`: Logic for searching cruises.
        *   `inventory_service.go`: Logic for inventory and locking.
        *   `order_service.go`: Logic for creating/cancelling orders and notices.
        *   `payment_service.go`: Logic for payments.
        *   `voyage_service.go`: Logic for voyages.
        *   `report_service.go`: Financial reporting logic.
        *   `recommendation.go`: Recommendation logic.
        *   `trend_service.go`: Price trend logic.
        *   `excursion_service.go`: Shore excursion logic.
        *   `pricing_engine.go`: Dynamic pricing logic.
        *   `marketing.go`: Marketing automation.
    *   `data/`: Data Access Layer.
        *   `db.go`: GORM connection.
        *   `redis.go`: Redis connection.
        *   `minio.go`: MinIO connection.
        *   `cruise_repo.go`: DB operations for Cruises.
        *   `inventory_repo.go`: DB operations for Inventory.
        *   `order_repo.go`: DB operations for Orders.
        *   `voyage_repo.go`: DB operations for Voyages.
        *   `excursion_repo.go`: DB operations for Excursions.
    *   `middleware/`: HTTP Middleware.
        *   `auth.go`: JWT Authentication middleware.
        *   `casbin.go`: RBAC Authorization middleware.
    *   `model/`: Domain entities and DTOs.
        *   `user.go`: User and Staff models.
        *   `cruise.go`: Cruise and CabinType models.
        *   `facility.go`: Facility model.
        *   `inventory.go`: Voyage, Cabin, Inventory models.
        *   `order.go`: Order, OrderItem, Passenger models.
        *   `excursion.go`: Shore Excursion models.
        *   `review.go`: User Reviews.
*   `pkg/`: Public/Shared libraries.
        *   `logger/`: Centralized Zap logger.
        *   `response/`: Standard API response structure.
        *   `storage/`: File upload utilities.
        *   `image/`: Image generation (Poster).
*   `tests/`: Integration and E2E tests.
    *   `unit/`: Unit tests.
        *   `middleware_test.go`: Auth/Middleware tests.
        *   `cruise_service_test.go`: Logic tests for CruiseService.
        *   `order_state_test.go`: Logic tests for Order state.
        *   `cancellation_test.go`: Logic tests for Order cancellation.
    *   `integration/`: Integration tests.
        *   `inventory_test.go`: Locking tests.
        *   `rbac_test.go`: Permission tests.
*   `go.mod`: Go module definition.

## Frontend - Admin (Nuxt 4)
Path: `/admin`

*   `components/`: Vue components.
*   `pages/`: Nuxt pages (File-based routing).
    *   `cruises/`: Cruise Management.
        *   `index.vue`: List/Add.
    *   `inventory/`: Inventory Management.
        *   `index.vue`: Dashboard.
    *   `orders/`: Order Management.
        *   `index.vue`: List/Detail.
    *   `marketing/`: Marketing & Pricing.
        *   `pricing.vue`: Pricing Rules.
*   `stores/`: Pinia state management stores.
*   `nuxt.config.ts`: Nuxt configuration.

## Frontend - Web (Nuxt 4)
Path: `/web`

*   `components/`: Vue components.
    *   `Recommendations.vue`: AI Recommendation list.
    *   `Reviews.vue`: User Reviews list.
*   `pages/`: Nuxt pages (File-based routing).
    *   `cruises/`: Cruise feature pages.
        *   `index.vue`: List/Search page.
        *   `[id].vue`: Detail page.
    *   `booking/`: Booking feature pages.
        *   `create.vue`: Booking form.
        *   `pay.vue`: Payment page.
    *   `user/`: User feature pages.
        *   `orders.vue`: My Orders page.
        *   `order_detail.vue`: Order Detail & Notice Download.
*   `stores/`: Pinia state management stores.
*   `nuxt.config.ts`: Nuxt configuration.

## Mini-program (Uni-app)
Path: `/mp`

*   `pages/`: Vue pages for Mini-program.
    *   `cruises/`: Cruise feature pages.
        *   `index.vue`: List page.
        *   `detail.vue`: Detail page.
    *   `user/`: User feature pages.
        *   `orders.vue`: My Orders page.
*   `static/`: Static assets.
*   `stores/`: Pinia state management stores.
*   `package.json`: Dependencies.

## Specifications & Docs
Path: `/specs/001-cruise-booking-system/`

*   `plan.md`: Implementation Plan.
*   `spec.md`: Feature Specification.
*   `tasks.md`: Actionable Task List.
*   `data-model.md`: Database Schema.
*   `contracts/`: API Specifications (OpenAPI).
*   `research.md`: Technical Decisions.
