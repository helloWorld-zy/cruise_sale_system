# Implementation Plan: Cruise Booking System (MVP + Full)

**Branch**: `001-cruise-booking-system` | **Date**: 2026-02-10 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `specs/001-cruise-booking-system/spec.md`

## Summary

Implementation of a comprehensive Cruise Booking System, starting with a strict MVP (Sales, Booking, Payment) and evolving into a full-service platform (Post-booking, Social, Intelligent Ops). The system comprises a Golang backend, Nuxt 4 Admin/Web frontends, and a Uni-app Mini-program.

## Technical Context

**Language/Version**: Go 1.26+ (Backend), TypeScript 5.9+ (Frontend)
**Primary Dependencies**: 
- Backend: Gin v1.11.0, GORM v2, NATS JetStream, Meilisearch
- Frontend: Nuxt 4.3.0, Vue 3.5, Pinia v3, Tailwind CSS v4, Uni-app
**Storage**: PostgreSQL 17.x, Redis 7.4.x, MinIO (S3 compatible)
**Testing**: 
- Backend: Testify, Gomock, Httptest (100% coverage mandatory)
- Frontend: Vitest, Vue Test Utils, Playwright (100% coverage mandatory)
**Target Platform**: Kubernetes (Containerized), WeChat Mini-program, Modern Web Browsers
**Project Type**: Full-stack System (Backend API + 3 Frontends)
**Performance Goals**: <2s page load, <15min inventory lock reliability, 1000+ concurrent users
**Constraints**: Strict 100% test coverage, Phased delivery (MVP -> V2.0)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- [x] **Uncompromising Quality**: Plan includes 100% test coverage strategy for all modules.
- [x] **Technology Standardization**: Stack matches Constitution exactly (Go/Gin, Nuxt/Vue, PG/Redis).
- [x] **API-First**: RESTful + Swagger/OpenAPI 3.1 is the integration standard.
- [x] **Phased Delivery**: Plan follows MVP -> V1.0 -> V1.5 -> V2.0 progression.

## Project Structure

### Documentation (this feature)

```text
specs/001-cruise-booking-system/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (OpenAPI)
└── checklists/          # Quality checks
```

### Source Code (repository root)

```text
backend/
├── cmd/server/          # Entry point
├── internal/
│   ├── api/             # Handlers & Routers (Gin)
│   ├── core/            # Business Logic (Services)
│   ├── data/            # Data Access (Repositories/GORM)
│   ├── middleware/      # Auth, Logging, CORS
│   └── model/           # Domain Entities
├── pkg/                 # Shared utilities
└── tests/               # Integration/E2E tests

admin/                   # Nuxt 4 Management Dashboard
├── components/
├── pages/
├── server/              # BFF if needed (or minimal)
└── stores/              # Pinia

web/                     # Nuxt 4 Customer Frontend
├── components/
├── pages/
└── stores/

mp/                      # Uni-app Mini-program
├── pages/
├── static/
└── stores/
```

**Structure Decision**: Multi-project repository (Monorepo-style) to manage all components in one place, facilitating code sharing (where applicable) and unified CI/CD.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| N/A | | |

## Implementation Strategy

### Phase 1: MVP (Sprint 1-4)
1. **Foundation**: Go setup, DB Schema (Cruises/Cabins), Nuxt/Uni-app init.
2. **Product**: Cruise/Cabin management (Admin), Browsing (Public).
3. **Booking**: Inventory logic, Locking, Order creation, Payment integration.
4. **Ops**: Basic Admin dashboard, User management.

### Phase 2: V1.0 - Full Experience (Sprint 5-8)
1. **Services**: Shore excursions, E-tickets.
2. **Interaction**: Notifications, Pre-trip checklist.
3. **Evaluation**: Reviews, basic loyalty.

### Phase 3: V1.5 - Intelligent (Sprint 9-12)
1. **AI**: Recommendation engine, Price prediction.
2. **Visuals**: VR integration, Interactive Deck plans.
3. **Fintech**: Installments, Multi-currency, OCR.

### Phase 4: V2.0 - Ecosystem (Sprint 13-16)
1. **Community**: UGC, Group buying.
2. **Advanced Ops**: Dynamic Pricing Engine, CRM, Marketing Automation.