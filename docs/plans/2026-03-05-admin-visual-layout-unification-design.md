# 2026-03-05 Admin Visual And Layout Unification Design

## 1. Background And Goal
- Target: Refactor all `frontend/admin` pages to a highly faithful `vue-element-admin` style.
- Scope: Visual and layout only. Do not change business API contracts or behavior.
- Additional requirement: `/login` must become a standalone system login page with no admin shell elements.

## 2. Confirmed Constraints
- Full-page one-shot unification across all admin pages.
- High visual fidelity to `vue-element-admin` default style.
- Include interaction unification and page-template unification.
- DOM structure reordering is allowed, but business logic stays unchanged.

## 3. Chosen Approach
- Preferred approach: Token-driven CSS system + page template restructuring + reusable shell components.
- Why: Best balance of fidelity, consistency, and long-term maintainability for full-scope rollout.

## 4. Information Architecture And Shell
### 4.1 Global Shell
- `Sidebar + Topbar + Content` three-zone shell for all admin business pages.
- Sidebar: dark theme, active highlight, collapse support.
- Topbar: collapse trigger, breadcrumb, right-side action placeholders.
- Content: unified container width and paddings.

### 4.2 Route-Level Layout Split
- `layouts/default.vue`: admin shell only.
- `layouts/auth.vue`: standalone login shell.
- `/login`: uses auth layout and renders login content only.

## 5. Visual System (High Fidelity)
### 5.1 Design Tokens
- Primary: `#409EFF`
- Success: `#67C23A`
- Warning: `#E6A23C`
- Danger: `#F56C6C`
- Info: `#909399`
- Page background: `#f0f2f5`
- Card background: `#ffffff`
- Text primary: `#303133`
- Text secondary: `#606266`
- Border: `#DCDFE6`

### 5.2 Typography And Spacing
- Font stack: `"Helvetica Neue", Helvetica, "PingFang SC", "Microsoft YaHei", sans-serif`
- Title: `20/600`, section title: `16/600`, body/table: `14/400`, helper: `12/400`
- Spacing scale: `4/8/12/16/20/24/32`
- Radius: card `6px`, control `4px`

### 5.3 States And Motion
- Unified control states: default/hover/active/disabled/focus-visible.
- Transition duration: `120-180ms` with `ease-out`.
- Focus and contrast are preserved for readability and keyboard usage.

## 6. Unified Page Templates
### 6.1 List Page Template
- `PageHeader` (title + primary action)
- `FilterArea` (query controls)
- `TableCard` (loading/error/empty/data)
- `FooterBar` (pagination or batch actions)

### 6.2 Form Page Template
- `PageHeader`
- `FormCard` with section blocks and grid fields
- `ActionBar` with primary/secondary actions

### 6.3 Detail/Mixed Page Template
- `PageHeader`
- `InfoCard` + related list cards

## 7. Component Mapping Plan
- Reusable components: `AdminPageHeader`, `AdminFilterBar`, `AdminDataCard`, `AdminFormCard`, `AdminActionBar`, `AdminStatusTag`.
- Existing components (`AdminActionLink`, `AdminConfirmDialog`) keep APIs; only visual adaptation.
- Keep existing `data-test` selectors unchanged.

## 8. Target Page Coverage
### 8.1 List Pages
- `dashboard/index.vue`
- `cruises/index.vue`, `companies/index.vue`, `voyages/index.vue`
- `cabin-types/index.vue`, `cabins/index.vue`, `bookings/index.vue`
- `facilities/index.vue`, `facility-categories/index.vue`
- `staff/index.vue`, `finance/index.vue`, `notifications/templates.vue`

### 8.2 Form Pages
- `cruises/create.vue`, `cruises/[id].vue`
- `companies/[id].vue`
- `voyages/new.vue`, `voyages/[id].vue`
- `cabin-types/new.vue`, `cabin-types/[id].vue`
- `cabins/new.vue`, `cabins/[id].vue`
- `bookings/new.vue`
- `facilities/new.vue`, `facilities/[id].vue`
- `facility-categories/[id].vue`
- `settings/shop.vue`

### 8.3 Detail/Tool Pages
- `bookings/[id].vue`
- `cabins/inventory.vue`, `cabins/pricing.vue`, `cabin-types/pricing.vue`

## 9. Login Page Specific Design
- Standalone route-level page with auth-only layout.
- Centered login card and high-fidelity default style inspired by `vue-element-admin`.
- Inputs: username and password only.
- No sidebar, topbar, breadcrumb, or any other admin shell content.
- Error feedback appears inside login card.

## 10. Data Flow And Guarding
- Keep current auth token flow and `useApi` behavior unchanged.
- Route guard behavior:
  - unauthenticated to admin page -> redirect `/login`
  - authenticated to `/login` -> redirect `/dashboard`

## 11. Risk Control
- Main risk: full-scope DOM reordering can cause behavior drift.
- Mitigations:
  - do not modify request URLs/payloads or service logic
  - preserve `data-test` attributes
  - enforce table action no-wrap + horizontal overflow strategy

## 12. Verification Strategy
- Unit tests:
  - login layout isolation
  - shell rendering
  - representative list/form/detail pages
- Smoke tests:
  - unauthenticated redirect to `/login`
  - login success to `/dashboard`
  - core page navigation and key actions
- Manual checks:
  - responsive behavior on desktop and mobile widths
  - visual consistency across full menu routes

## 13. Definition Of Done
- All admin pages follow one visual language and one template system.
- `/login` is fully isolated as a standalone login page.
- Business behavior remains unchanged.
- Core tests and smoke checks pass.
