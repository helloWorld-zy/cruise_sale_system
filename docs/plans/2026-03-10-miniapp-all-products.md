# Miniapp All Products Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a miniapp all-products page with company expand-collapse navigation, cruise-level filtering, and voyage cards backed by a new public voyage listing API.

**Architecture:** Add a public voyage listing route on the backend that returns front-end-safe voyage list data, then replace the current miniapp all-products tab entry with a dedicated page that loads companies, cruises, and voyages together. The page will keep company expand-collapse state locally and filter the right-side voyage cards by selected cruise.

**Tech Stack:** Go, Gin, Vue 3, Vite, Vitest, Vue Testing Library

---

### Task 1: Add backend red tests for public voyage listing

**Files:**
- Modify: `backend/internal/handler/*voyage*test*.go`
- Modify: `backend/internal/router/router_test.go`

**Step 1: Write the failing test**

Add handler and router tests that assert:

- `GET /api/v1/voyages?page=1&page_size=20` returns public voyage data
- `GET /api/v1/voyages?cruise_id=11&page=1&page_size=20` filters by cruise
- route exists without admin auth

**Step 2: Run test to verify it fails**

Run: `go test ./internal/handler ./internal/router`

Expected: FAIL because public voyage list endpoint is not wired or not implemented.

**Step 3: Write minimal implementation**

Implement the smallest backend changes needed to expose a public voyage listing route and response DTO.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/handler ./internal/router`

Expected: PASS.

### Task 2: Implement backend public voyage list endpoint

**Files:**
- Modify: `backend/internal/router/router.go`
- Modify: `backend/internal/handler/voyage_handler.go`
- Modify: backend service or repository files that currently back voyage listing

**Step 1: Add public list query support**

Support `cruise_id`, `page`, and `page_size` query parameters for public listing.

**Step 2: Return front-end-safe voyage card data**

Ensure each item includes enough fields for miniapp rendering: voyage ID, cruise ID, cruise name, image, brief info or title text, dates, and any available price fields.

**Step 3: Wire public route**

Register `GET /api/v1/voyages` under the public API group.

**Step 4: Re-run focused backend tests**

Run: `go test ./internal/handler ./internal/router ./internal/service ./internal/repository`

Expected: PASS.

### Task 3: Add miniapp red tests for all-products interaction

**Files:**
- Create or modify: `frontend/miniapp/tests/unit/pages/products/*.spec.ts`
- Modify: `frontend/miniapp/tests/unit/app.spec.ts`

**Step 1: Write the failing tests**

Cover these behaviors:

- all-products tab renders the new page instead of cabin list preview
- page loads companies, cruises, voyages on mount
- clicking company toggles the nested cruise list only
- clicking cruise filters right-side voyage cards
- clicking a voyage card emits open-voyage or routes into voyage detail flow

**Step 2: Run test to verify it fails**

Run: `npx vitest run frontend/miniapp/tests/unit/app.spec.ts frontend/miniapp/tests/unit/pages/products`

Expected: FAIL because the new page and interactions do not exist yet.

**Step 3: Write minimal implementation**

Build the page, tree navigation, and voyage cards with only the behavior required by the tests.

**Step 4: Run test to verify it passes**

Run: `npx vitest run frontend/miniapp/tests/unit/app.spec.ts frontend/miniapp/tests/unit/pages/products`

Expected: PASS.

### Task 4: Replace tab wiring and add voyage detail entry flow

**Files:**
- Modify: `frontend/miniapp/src/App.vue`
- Create or modify: `frontend/miniapp/pages/...`

**Step 1: Switch all-products tab**

Point the `全部商品` tab at the new page component.

**Step 2: Add voyage selection flow**

When a voyage card is clicked, switch into voyage detail entry flow. If no dedicated voyage detail page exists, route through the closest existing flow with explicit prop naming and a TODO-safe fallback.

**Step 3: Verify with tests**

Run the focused miniapp tests again.

### Task 5: Final verification

**Files:**
- No new files

**Step 1: Run backend verification**

Run: `go test ./internal/handler ./internal/router ./internal/service ./internal/repository`

**Step 2: Run miniapp verification**

Run: `npx vitest run`

**Step 3: Optional preview smoke if environment is available**

Run preview or existing smoke for miniapp all-products view.

Expected: targeted tests pass and no regression is introduced in touched areas.