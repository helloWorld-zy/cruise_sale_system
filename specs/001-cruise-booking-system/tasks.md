# Tasks: Cruise Booking System (MVP + Full)

**Input**: Design documents from `specs/001-cruise-booking-system/`
**Prerequisites**: plan.md, spec.md, data-model.md, contracts/openapi.yaml
**Feature Branch**: `001-cruise-booking-system`

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and environment setup.

- [ ] T001 Create project structure (backend, admin, web, mp) per plan.md
- [ ] T002 [P] Initialize Backend (Go 1.26 + Gin) with `go mod init`
- [ ] T003 [P] Initialize Admin (Nuxt 4) with `npx nuxi init admin`
- [ ] T004 [P] Initialize Web (Nuxt 4) with `npx nuxi init web`
- [ ] T005 [P] Initialize Mini-program (Uni-app)
- [ ] T006 Configure Docker Compose (Postgres 17, Redis 7, MinIO, Meilisearch, NATS)
- [ ] T007 Configure CI/CD pipeline for 100% coverage check

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure, database schema, and authentication.

- [ ] T008 Setup GORM connection and migration framework in `backend/internal/data/db.go`
- [ ] T009 Implement User & Staff models in `backend/internal/model/user.go`
- [ ] T010 Implement JWT Authentication middleware in `backend/internal/middleware/auth.go`
- [ ] T011 Implement Casbin RBAC middleware in `backend/internal/middleware/casbin.go`
- [ ] T012 Implement centralized Error Handling & Logger in `backend/pkg/logger/`
- [ ] T013 Create Base Response structure in `backend/pkg/response/response.go`
- [ ] T014 [P] Test Auth & Middleware (Unit Tests) in `backend/tests/unit/middleware_test.go`

## Phase 3: User Story 1 - Guest Cruise Browsing (Priority: P1)

**Goal**: Users can browse, filter, and view cruise details.
**Independent Test**: Public user visits home, filters by date, views cruise details.

### Implementation
- [ ] T015 [P] [US1] Create Cruise & CabinType models in `backend/internal/model/cruise.go`
- [ ] T016 [P] [US1] Create Facility model in `backend/internal/model/facility.go`
- [ ] T017 [US1] Implement CruiseRepository (List/Get) in `backend/internal/data/cruise_repo.go`
- [ ] T018 [US1] Implement CruiseService (Filter logic) in `backend/internal/core/cruise_service.go`
- [ ] T019 [US1] Implement API Handlers (GET /cruises, /cruises/:id) in `backend/internal/api/v1/cruise.go`
- [ ] T020 [P] [US1] Implement Cruise List page in Web `web/pages/cruises/index.vue`
- [ ] T021 [P] [US1] Implement Cruise Detail page in Web `web/pages/cruises/[id].vue`
- [ ] T022 [P] [US1] Implement Cruise List/Detail in MP `mp/pages/cruises/`

### Tests
- [ ] T023 [P] [US1] Unit Test CruiseService in `backend/tests/unit/cruise_service_test.go`
- [ ] T024 [P] [US1] E2E Test Browser Flow in `web/tests/e2e/browse.spec.ts`

## Phase 4: User Story 2 - Booking & Checkout (Priority: P1)

**Goal**: Users can select cabins, book, and pay.
**Independent Test**: User selects cabin -> Locks inventory -> Pays -> Order Confirmed.

### Implementation
- [ ] T025 [P] [US2] Create Voyage, Cabin, Inventory models in `backend/internal/model/inventory.go`
- [ ] T026 [P] [US2] Create Order, OrderItem, Passenger models in `backend/internal/model/order.go`
- [ ] T027 [US2] Implement InventoryService (Lock/Unlock/Check) in `backend/internal/core/inventory_service.go`
- [ ] T028 [US2] Implement OrderService (Create/Update Status) in `backend/internal/core/order_service.go`
- [ ] T029 [US2] Implement PaymentService (Mock/WeChat) in `backend/internal/core/payment_service.go`
- [ ] T030 [US2] Implement Booking API (POST /orders) in `backend/internal/api/v1/order.go`
- [ ] T031 [P] [US2] Implement Booking Form (Passenger Input) in Web `web/pages/booking/create.vue`
- [ ] T032 [P] [US2] Implement Payment Page in Web `web/pages/booking/pay.vue`

### Tests
- [ ] T033 [P] [US2] Integration Test Locking Mechanism in `backend/tests/integration/inventory_test.go`
- [ ] T034 [P] [US2] Unit Test Order State Machine in `backend/tests/unit/order_state_test.go`

## Phase 5: User Story 3 - User Order Management (Priority: P1)

**Goal**: Users can view and manage their orders.

- [ ] T035 [US3] Implement My Orders API (GET /orders/mine) in `backend/internal/api/v1/user_order.go`
- [ ] T036 [US3] Implement Cancel Order logic in `backend/internal/core/order_service.go`
- [ ] T037 [P] [US3] Implement My Orders Page in Web `web/pages/user/orders.vue`
- [ ] T038 [P] [US3] Implement My Orders Page in MP `mp/pages/user/orders.vue`
- [ ] T039 [P] [US3] Test Cancellation Policy Logic in `backend/tests/unit/cancellation_test.go`

## Phase 6: User Story 4 - Admin Content & Inventory (Priority: P1)

**Goal**: Admins can manage cruises and inventory.

- [ ] T040 [US4] Implement Admin Cruise CRUD API in `backend/internal/api/v1/admin/cruise.go`
- [ ] T041 [US4] Implement Admin Inventory/Voyage API in `backend/internal/api/v1/admin/voyage.go`
- [ ] T042 [P] [US4] Implement Admin Cruise Management UI in `admin/pages/cruises/index.vue`
- [ ] T043 [P] [US4] Implement Admin Inventory Dashboard in `admin/pages/inventory/index.vue`
- [ ] T044 [P] [US4] Test Admin Permissions in `backend/tests/integration/rbac_test.go`

## Phase 7: User Story 5 - Admin Order & Finance (Priority: P2)

**Goal**: Ops can view orders and finance reports.

- [ ] T045 [US5] Implement Admin Order List/Detail API in `backend/internal/api/v1/admin/order.go`
- [ ] T046 [US5] Implement Financial Report Aggregation in `backend/internal/core/report_service.go`
- [ ] T047 [P] [US5] Implement Admin Order UI in `admin/pages/orders/index.vue`

## Phase 8: User Story 6 - Intelligent Discovery (Priority: P2)

**Goal**: Recommendations and Price Trends.

- [ ] T048 [US6] Implement Recommendation Engine Stub in `backend/internal/core/recommendation.go`
- [ ] T049 [US6] Implement Price Trend API in `backend/internal/api/v1/price_trend.go`
- [ ] T050 [P] [US6] Add Recommendation Component to Web Home `web/components/Recommendations.vue`

## Phase 9: User Story 7 - Post-Booking Experience (Priority: P2)

**Goal**: E-tickets and Shore Excursions.

- [ ] T051 [US7] Implement Admin Upload API for Departure Notice in `backend/internal/api/v1/admin/order.go`
- [ ] T052 [US7] Implement Shore Excursion Models & API in `backend/internal/api/v1/excursion.go`
- [ ] T053 [P] [US7] Implement Departure Notice Pop-up & Download in Web `web/pages/user/order_detail.vue`

## Phase 10: User Story 9 - Intelligent Ops (Priority: P2)

**Goal**: Dynamic Pricing and Marketing.

- [ ] T054 [US9] Implement Dynamic Pricing Rules Engine in `backend/internal/core/pricing_engine.go`
- [ ] T055 [US9] Implement Marketing Trigger Service in `backend/internal/core/marketing.go`
- [ ] T056 [P] [US9] Implement Pricing Rules UI in `admin/pages/marketing/pricing.vue`

## Phase 11: User Story 8 - Social & Community (Priority: P3)

**Goal**: Reviews and Sharing.

- [ ] T057 [US8] Implement Review Model & API in `backend/internal/api/v1/review.go`
- [ ] T058 [US8] Implement Poster Generation Service in `backend/pkg/image/poster.go`
- [ ] T059 [P] [US8] Add Review Section to Cruise Detail `web/components/Reviews.vue`

## Phase 12: Polish & Cross-Cutting

- [ ] T060 Update Swagger/OpenAPI Documentation
- [ ] T061 Perform Load Testing (k6) for Inventory Locking
- [ ] T062 Verify 100% Test Coverage Report
- [ ] T063 Final Localization Check (Chinese Default)

## Dependencies & Execution Order

- **Phase 1 & 2** (Setup/Foundational) BLOCKS all User Stories.
- **US1** (Browsing) is prerequisite for **US2** (Booking).
- **US2** (Booking) is prerequisite for **US3** (User Orders), **US5** (Admin Orders), **US7** (Post-Booking).
- **US4** (Admin Content) can run parallel with **US1** (Browsing).
- **US6** (Discovery) and **US8** (Social) are largely independent enhancements.
- **US9** (Ops) depends on **US2** (Booking Data).

## Parallel Example: User Story 2 (Booking)

```bash
# Backend Team
Task T027: Implement InventoryService
Task T028: Implement OrderService
Task T030: Implement Booking API

# Frontend Team
Task T031: Implement Booking Form
Task T032: Implement Payment Page
```
