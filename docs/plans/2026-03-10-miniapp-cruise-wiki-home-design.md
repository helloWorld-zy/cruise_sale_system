# 2026-03-10 Miniapp Cruise Wiki Home Design

## 1. Background And Goal
- Target: Build a real `邮轮百科` home page for `frontend/miniapp` instead of the current placeholder.
- Primary UX: Follow the provided mobile reference layout with a left-side cruise company list and a right-side cruise card list.
- Default behavior: On first entry, the right side shows all enabled cruises.
- Filter behavior: After the user manually selects a company on the left, the right side shows only enabled cruises under that enabled company.
- Navigation behavior: Clicking a cruise card enters the existing cruise detail page.
- Data source requirement: Company data and cruise data must come from real backend public APIs.

## 2. Confirmed Constraints
- Frontend is a Vue-based uni-app style miniapp implementation and must follow Vue best practices.
- Use Vue 3 Composition API with `<script setup lang="ts">` for any new or refactored page/component code.
- UI implementation can additionally follow the provided `ui-ux-pro-max` skill for mobile interaction, spacing, contrast, touch area, and visual fidelity.
- Frontend should only show enabled companies and enabled cruises.
- Avoid introducing a full routing system for miniapp if existing app state switching already handles page changes.
- Keep changes minimal and aligned with current project structure.

## 3. Chosen Approach
- Preferred approach: Add one public backend company list API, reuse the existing public cruise list/detail APIs with targeted enhancement, and build a dedicated wiki page in miniapp.
- Why: This keeps responsibilities clear, minimizes backend coupling, preserves current app architecture, and still provides a complete real-data flow.

## 4. Backend API Design
### 4.1 New Public Company List API
- Route: `GET /api/v1/companies`
- Purpose: Return enabled cruise companies for miniapp/web public browsing.
- Query params:
  - `page` default `1`
  - `page_size` default `50`
  - optional `keyword` if later needed, but not required by the first miniapp page version
- Sort order:
  - `sort_order DESC`
  - `id DESC`
- Returned fields:
  - `id`
  - `name`
  - `english_name`
  - `logo_url`
  - `description`
  - `sort_order`
- Filter rule:
  - only `status = 1`

### 4.2 Public Cruise List API Adjustment
- Reuse: `GET /api/v1/cruises`
- Required behavior for public traffic:
  - default to enabled cruises only
  - support `company_id` filter for company-scoped browsing
- Initial wiki page request:
  - `GET /api/v1/cruises?page=1&page_size=30`
- Filtered request after company click:
  - `GET /api/v1/cruises?company_id=<id>&page=1&page_size=30`
- Returned items should remain compatible with current miniapp card and detail consumption.

### 4.3 Public Cruise Detail Continuity
- Reuse: `GET /api/v1/cruises/:id`
- No wiki-specific detail payload should be invented.
- The wiki page passes `cruiseId` and lets the detail page fetch its own source of truth.

## 5. Frontend Information Architecture
### 5.1 Page Replacement
- Replace the current `wiki` placeholder branch inside `frontend/miniapp/src/App.vue` with a real `WikiHomePage` component.

### 5.2 Component Boundaries
- `pages/wiki/index.vue`: page container, data orchestration, selected company state, loading/error handling, card click handling.
- Optional child components if the page grows too large:
  - `components/wiki/WikiCompanySidebar.vue`
  - `components/wiki/WikiCruiseCard.vue`
  - `components/wiki/WikiCruiseList.vue`
- Principle:
  - page container owns async state
  - child components stay mostly presentational with props/emits

### 5.3 App-Level Navigation
- Continue using current `activeTab` state switching in `App.vue`.
- Add a page-level selected cruise state at app shell level or reuse existing detail branch input pattern.
- Card click flow:
  - set selected cruise id
  - switch `activeTab` to `cabin-detail` or a proper cruise detail tab branch aligned with the current app shell
- If needed, normalize naming in `App.vue` so wiki detail navigation targets the real cruise detail page instead of cabin preview branches.

## 6. Page Interaction Design
### 6.1 Header And Top Region
- Keep the small-program visual structure shown in the reference:
  - top white native-style header area
  - blue title strip with `邮轮百科`
- Preserve bottom tab bar already used by the miniapp shell.

### 6.2 Main Two-Column Layout
- Left side:
  - fixed-width company list column
  - first local item is `全部邮轮`
  - enabled company items follow backend order
  - active item uses clear highlight, blue accent edge, and stronger text color
- Right side:
  - vertically scrollable cruise card list
  - cards prioritize cover image, Chinese name, and English name
  - cards should have comfortable radius, white background, strong image area, and clear tap affordance

### 6.3 Default And Filtered States
- First load:
  - selected company is local `全部邮轮`
  - right panel shows all enabled cruises
- After manual company click:
  - selected company id updates
  - right panel refreshes with `company_id`
  - only enabled cruises under that company remain visible

### 6.4 Loading, Empty, And Error States
- Company list loading:
  - left column shows lightweight skeleton or loading text
- Cruise list loading:
  - right panel shows loading placeholders without collapsing layout
- Empty state:
  - keep selected company visible and show a right-panel empty state message
- Error state:
  - do not blank the whole page
  - company and cruise areas should fail independently where practical

## 7. Vue And UI/UX Implementation Rules
- Use Vue 3 Composition API and `<script setup lang="ts">`.
- Keep derived display state in `computed` instead of mutating duplicated state.
- Keep API side effects in focused async functions triggered from lifecycle hooks or explicit user actions.
- Prefer smaller presentational child components if the page exceeds a single clear responsibility.
- Ensure touch targets remain at least mobile-friendly size.
- Preserve sufficient color contrast and visible selected/focus states.
- Use image `alt` text or equivalent descriptive attributes in web-preview-compatible markup.
- Avoid hover-only interaction assumptions; primary actions must work by tap.

## 8. Existing Detail Page Risk To Address
- Current `frontend/miniapp/pages/cruise/detail.vue` requests `/routes`, but there is no matching public backend route.
- This must be corrected during implementation; otherwise the new wiki entry path will navigate into a page that can still fail.
- Acceptable fixes:
  - remove the unsupported request and degrade gracefully
  - or replace it with a real public data source already supported by backend
- The implementation must prefer a real available backend contract and avoid console noise or visible errors in preview.

## 9. Testing Strategy
### 9.1 Backend
- Add handler/service/repository tests for the public company list endpoint.
- Add coverage ensuring public company list returns only enabled companies.
- Add coverage ensuring public cruise list respects enabled-only behavior and `company_id` filtering.

### 9.2 Frontend Unit Tests
- Initial render shows `全部邮轮` selected and all enabled cruises returned from API.
- Clicking a company triggers filtered cruise loading.
- Clicking a cruise card switches to detail with the correct cruise id.
- Error and empty states render without crashing the app shell.

### 9.3 Preview Verification
- Verify the `邮轮百科` tab no longer shows placeholder text.
- Verify visual structure matches the provided mobile reference closely.
- Verify there are no new missing API requests in browser preview.

## 10. Definition Of Done
- `邮轮百科` in miniapp is a real page, not a placeholder.
- Left column lists enabled companies from backend and includes local `全部邮轮` at the top.
- First entry shows all enabled cruises on the right.
- Selecting a company filters the right-side cruises to that company only.
- Clicking a cruise card enters the existing cruise detail page successfully.
- The implementation follows Vue best practices and mobile UI/UX best practices.
- Backend and frontend targeted tests pass.