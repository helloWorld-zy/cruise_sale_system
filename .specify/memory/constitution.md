<!--
Sync Impact Report:
- Version Change: 0.0.0 -> 1.0.0
- Added Principle: Uncompromising Quality (100% Test Coverage)
- Added Principle: Technology Standardization (Mandatory Tech Stack)
- Added Principle: API-First & Documentation
- Added Principle: Phased Delivery
- Governance: Initial ratification.
- Templates: No immediate template changes required.
-->

# CruiseBooking Constitution

## Core Principles

### I. Uncompromising Quality
Test coverage must be **100%** for all codebases (Backend, Web, Mini-program). This is non-negotiable. CI pipelines will reject any merge request that falls below this threshold. Testing includes Unit, Integration, and E2E layers as defined in the Test Strategy.

### II. Technology Standardization
Adherence to the specified technology stack is mandatory.
*   **Backend**: Go 1.26+, Gin, GORM, PostgreSQL 17, Redis 7.4.
*   **Web Frontends**: Nuxt 4.3.0+, Vue 3.5+, Vite, Tailwind CSS v4.
*   **Mini-program**: uni-app (Vue 3 mode).
*   **Infrastructure**: Kubernetes, Docker, MinIO.
Any deviation must be formally approved via constitution amendment.

### III. API-First & Documentation
All backend interfaces must follow RESTful standards and provide up-to-date Swagger/OpenAPI 3.1 documentation. API documentation is the source of truth for frontend-backend integration and must be auto-generated where possible.

### IV. Phased Delivery
Development follows a strict phased approach:
1.  **MVP**: Core Sales (Cruises, Cabins, Booking, Payment, Orders).
2.  **V1.0**: Complete Experience (Shore excursions, e-tickets).
3.  **V1.5**: Intelligent Features (AI, VR).
4.  **V2.0**: Ecosystem (Community, Distribution).
Features must not be implemented before their scheduled phase without explicit approval.

## Technical Standards

**Backend (Go)**:
*   Framework: Gin v1.11.0
*   ORM: GORM v2.x
*   Search: Meilisearch
*   Messaging: NATS JetStream

**Frontend (Nuxt/Vue)**:
*   State Management: Pinia v3
*   Language: TypeScript 5.9.x (Strict Mode)
*   UI Libs: Nuxt UI v3 (Web), uni-ui (Mini-program)

## Governance

This constitution governs the CruiseBooking project.
*   **Amendments**: Changes to the tech stack or principles require a constitution amendment and version bump.
*   **Compliance**: All PRs must be reviewed against these principles. Use the `check-prerequisites.ps1` script where applicable to ensure environment readiness.

**Version**: 1.0.0 | **Ratified**: 2026-02-10 | **Last Amended**: 2026-02-10