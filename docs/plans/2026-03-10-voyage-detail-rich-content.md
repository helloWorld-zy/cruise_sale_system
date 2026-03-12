# Voyage Detail Rich Content Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add global fee-note and booking-notice template management, voyage-level template snapshot overrides, and upgrade miniapp voyage detail to display itinerary meals, accommodation, route overview, schedule table, fee notes, and booking notices.

**Architecture:** Extend the voyage aggregate with template reference and resolved content fields, add an admin content-template module plus voyage form controls, and render the richer voyage detail experience from the public voyage detail API. The implementation stays TDD-first and keeps the miniapp consuming a single enriched voyage detail response rather than stitching multiple template requests on the client.

**Tech Stack:** Go, Gin, GORM, Vue 3, Nuxt 4 admin, Vite miniapp, Vitest

---

### Task 1: Add backend red tests for template domain and voyage detail enrichment

**Files:**
- Modify: `backend/internal/handler/route_handler_test.go`
- Modify: `backend/internal/router/router_test.go`
- Modify: repository or service tests under `backend/internal/repository` and `backend/internal/service`

**Step 1: Write the failing tests**

Cover these cases:

- admin content-template list/create/update routes exist
- templates can be filtered by `kind`
- public `GET /api/v1/voyages/:id` returns resolved `fee_note` and `booking_notice`
- voyage create/update accepts template IDs, mode, and snapshot content

**Step 2: Run test to verify it fails**

Run: `go test ./internal/handler ./internal/router ./internal/service ./internal/repository`

Expected: FAIL because the template domain, API surface, and enriched voyage payload do not exist yet.

**Step 3: Write minimal implementation**

Add only the types and route wiring needed to make the tests compile and fail for the intended reason.

**Step 4: Run test to verify the failure is correct**

Run the same command and confirm the failure is due to missing implementation, not broken tests.

### Task 2: Add migration and backend domain support for global templates and voyage overrides

**Files:**
- Create: `backend/migrations/0000xx_content_templates.up.sql`
- Create: `backend/migrations/0000xx_content_templates.down.sql`
- Modify: `backend/internal/domain/voyage.go`
- Create: `backend/internal/domain/content_template.go`
- Modify: `backend/internal/domain/repository.go`

**Step 1: Add migration**

Create tables/columns for:

- `content_templates`
- voyage-side template reference and snapshot fields for fee note and booking notice

**Step 2: Add domain models**

Define `ContentTemplate` and voyage fields for template references, mode, and snapshot content.

**Step 3: Run focused backend tests**

Run: `go test ./internal/repository ./internal/service`

Expected: still failing, but now against repository/service behavior rather than missing schema types.

### Task 3: Implement backend repository/service/handler/router for templates and resolved voyage detail

**Files:**
- Create: `backend/internal/repository/content_template_repo.go`
- Create or modify: `backend/internal/service/content_template_service.go`
- Create: `backend/internal/handler/content_template_handler.go`
- Modify: `backend/internal/repository/voyage_repo.go`
- Modify: `backend/internal/handler/voyage_handler.go`
- Modify: `backend/internal/router/router.go`

**Step 1: Implement template CRUD**

Support admin list/detail/create/update/delete with `kind` filtering.

**Step 2: Persist voyage template selection and snapshot mode**

Extend voyage upsert payload handling and persistence.

**Step 3: Resolve voyage detail content**

When public voyage detail is fetched, resolve `fee_note` and `booking_notice` according to voyage mode.

**Step 4: Run focused backend verification**

Run: `go test ./internal/handler ./internal/router ./internal/service ./internal/repository`

Expected: PASS for the newly added backend behavior.

### Task 4: Add admin red tests for template management and voyage template selection

**Files:**
- Create: `frontend/admin/tests/unit/pages/content-templates.spec.ts`
- Modify: `frontend/admin/tests/unit/pages/voyages-new.spec.ts`
- Modify: `frontend/admin/tests/unit/pages/voyages-id.spec.ts`
- Modify if needed: `frontend/admin/tests/unit/pages/task21-route-smoke.spec.ts`

**Step 1: Write the failing tests**

Cover these behaviors:

- admin nav exposes the content-template module
- template list switches kinds and loads rows
- template form submits structured fee-note and booking-notice payloads
- voyage create/edit page can select a template and switch to snapshot mode

**Step 2: Run test to verify it fails**

Run: `Set-Location d:\cruise_sale_system\frontend\admin; npx vitest run tests/unit/pages/content-templates.spec.ts tests/unit/pages/voyages-new.spec.ts tests/unit/pages/voyages-id.spec.ts`

Expected: FAIL because the new admin pages and voyage form controls do not exist yet.

### Task 5: Implement admin template pages and voyage form extensions

**Files:**
- Modify: `frontend/admin/app/layouts/default.vue`
- Create: `frontend/admin/app/pages/content-templates/index.vue`
- Create: `frontend/admin/app/pages/content-templates/[id].vue`
- Create if needed: `frontend/admin/app/components/content-template/*`
- Modify: `frontend/admin/app/pages/voyages/new.vue`
- Modify: `frontend/admin/app/pages/voyages/[id].vue`

**Step 1: Add template management UI**

Build the admin module for fee-note and booking-notice templates.

**Step 2: Extend voyage forms**

Add template selector, mode status, snapshot editor, and cleaner itinerary service chips for meals/accommodation.

**Step 3: Run focused admin tests**

Run the same Vitest command as Task 4.

Expected: PASS.

### Task 6: Add miniapp red tests for enriched voyage detail rendering

**Files:**
- Modify: `frontend/miniapp/tests/unit/pages/voyage/detail.spec.ts`
- Modify if needed: `frontend/miniapp/tests/unit/app.spec.ts`

**Step 1: Write the failing tests**

Cover these behaviors:

- daily itinerary cards show meals, accommodation, and arrival/departure times
- route overview SVG and schedule table render from itinerary data
- fee-note section renders `费用包含` and `费用不包含`
- booking-notice section renders structured tabs/sections with emphasis blocks

**Step 2: Run test to verify it fails**

Run: `Set-Location d:\cruise_sale_system\frontend\miniapp; npx vitest run tests/unit/pages/voyage/detail.spec.ts tests/unit/app.spec.ts`

Expected: FAIL because the miniapp detail page does not yet render the new sections.

### Task 7: Implement miniapp voyage detail redesign

**Files:**
- Modify: `frontend/miniapp/pages/voyage/detail.vue`
- Create if needed: `frontend/miniapp/components/voyage/*`

**Step 1: Upgrade itinerary cards**

Render meals, accommodation, and time chips from itinerary rows.

**Step 2: Add route overview + schedule table block**

Generate a client-side SVG route schematic and paired schedule table.

**Step 3: Add fee-note and booking-notice sections**

Render structured content with a mobile-first layout.

**Step 4: Run focused miniapp tests**

Run: `Set-Location d:\cruise_sale_system\frontend\miniapp; npx vitest run tests/unit/pages/voyage/detail.spec.ts tests/unit/app.spec.ts`

Expected: PASS.

### Task 8: Final verification

**Files:**
- No new files

**Step 1: Run backend verification**

Run: `Set-Location d:\cruise_sale_system\backend; go test ./internal/handler ./internal/router ./internal/service ./internal/repository`

**Step 2: Run admin verification**

Run: `Set-Location d:\cruise_sale_system\frontend\admin; npx vitest run tests/unit/pages/content-templates.spec.ts tests/unit/pages/voyages-new.spec.ts tests/unit/pages/voyages-id.spec.ts`

**Step 3: Run miniapp verification**

Run: `Set-Location d:\cruise_sale_system\frontend\miniapp; npx vitest run tests/unit/pages/voyage/detail.spec.ts tests/unit/app.spec.ts`

**Step 4: Run broader local acceptance when practical**

Run: `node scripts/local-acceptance.mjs --backend-base http://127.0.0.1:8080`

Expected: touched backend, admin, and miniapp flows remain green.