# Miniapp Cruise Wiki Home Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 为 `frontend/miniapp` 新增真实可用的 `邮轮百科` 首页，接通前台公司列表与邮轮筛选接口，并支持从百科卡片进入现有邮轮详情页。

**Architecture:** 后端沿用现有 `domain -> repository -> service -> handler -> router` 分层，为前台公开场景新增公司列表接口，并补齐前台邮轮列表的启用态与公司过滤约束。前端基于 Vue 3 `<script setup lang="ts">` 在 miniapp shell 内新增百科页面容器与轻量展示组件，复用现有详情页入口，并修正详情页中无效的 `/routes` 请求，确保新入口不会落到损坏页面。

**Tech Stack:** Go + Gin + GORM + Swagger docs + Vue 3 + TypeScript + Pinia + Vitest + Testing Library + uni-app style page structure

---

## Preconditions

- 参考技能：`@vue-best-practices` `@vue-development-guides` `@vue-testing-best-practices` `@ui-ux-pro-max` `@test-driven-development`
- 在独立 worktree 或隔离开发环境执行。
- 所有新增 Vue 页面与组件使用 `<script setup lang="ts">`。
- 保持 DRY / YAGNI，优先修复根因，不引入前台专用聚合接口。

### Task 1: Add Failing Backend Tests For Public Company Browsing

**Files:**
- Modify: `backend/internal/handler/all_handler_test.go`
- Modify: `backend/internal/repository/company_repo_test.go`
- Modify: `backend/internal/service/service_test.go`

**Step 1: Write the failing tests**

```go
func TestCompanyRepository_ListPublicEnabledOnly(t *testing.T) {
    // seed enabled + disabled companies
    // assert only enabled ones are returned for public listing
}

func TestCompanyService_ListPublic_ReturnsEnabledCompanies(t *testing.T) {
    // service-level contract for public company browsing
}

func TestAllHandler_PublicCompanies_ReturnsEnabledItems(t *testing.T) {
    // GET /companies
    // expect 200 and only enabled companies in payload
}
```

**Step 2: Run tests to verify they fail**

Run: `cd backend; go test ./internal/handler ./internal/repository ./internal/service -run "PublicCompanies|ListPublic" -v`
Expected: FAIL because the public company listing contract does not exist yet.

**Step 3: Write minimal implementation scaffolding**

```go
// company service
func (s *CompanyService) ListPublic(ctx context.Context, page, pageSize int) ([]domain.CruiseCompany, int64, error)

// company repository
func (r *CompanyRepo) ListPublic(ctx context.Context, page, pageSize int) ([]domain.CruiseCompany, int64, error)

// public route expectation
GET /api/v1/companies
```

**Step 4: Run tests to verify they pass**

Run: `cd backend; go test ./internal/handler ./internal/repository ./internal/service -run "PublicCompanies|ListPublic" -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add backend/internal/handler/all_handler_test.go backend/internal/repository/company_repo_test.go backend/internal/service/service_test.go
git commit -m "test(backend): cover public company browsing contract"
```

### Task 2: Implement Public Company API And Route

**Files:**
- Modify: `backend/internal/repository/company_repo.go`
- Modify: `backend/internal/domain/repo.go`
- Modify: `backend/internal/service/company_service.go`
- Modify: `backend/internal/handler/company_handler.go`
- Modify: `backend/internal/router/router.go`
- Modify: `backend/docs/swagger.yaml`
- Modify: `backend/docs/swagger.json`
- Modify: `backend/docs/docs.go`

**Step 1: Write the failing handler-level assertion if still missing**

```go
func TestAllHandler_PublicCompanies_SortsBySortOrderDesc(t *testing.T) {
    // expect public companies ordered for browsing
}
```

**Step 2: Run test to verify it fails**

Run: `cd backend; go test ./internal/handler -run "PublicCompanies" -v`
Expected: FAIL because route/ordering is not implemented yet.

**Step 3: Write minimal implementation**

```go
// router.go
api.GET("/companies", deps.Company.ListPublic)
```

```go
// company_handler.go
func (h *CompanyHandler) ListPublic(c *gin.Context) {
    page := queryInt(c, "page", 1)
    pageSize := queryInt(c, "page_size", 50)
    items, total, err := h.svc.ListPublic(c.Request.Context(), page, pageSize)
    if err != nil {
        response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
        return
    }
    response.Success(c, gin.H{"list": items, "total": total})
}
```

```go
// company_repo.go
q := db.WithContext(ctx).Model(&domain.CruiseCompany{}).Where("status = ?", 1)
q = q.Order("sort_order DESC").Order("id DESC")
```

**Step 4: Run tests to verify they pass**

Run: `cd backend; go test ./internal/handler ./internal/repository ./internal/service -run "PublicCompanies|ListPublic" -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add backend/internal/repository/company_repo.go backend/internal/domain/repo.go backend/internal/service/company_service.go backend/internal/handler/company_handler.go backend/internal/router/router.go backend/docs/swagger.yaml backend/docs/swagger.json backend/docs/docs.go
git commit -m "feat(backend): add public company list api"
```

### Task 3: Add Failing Backend Tests For Public Cruise Filtering

**Files:**
- Modify: `backend/internal/handler/all_handler_test.go`
- Modify: `backend/internal/repository/cruise_repo_test.go`
- Modify: `backend/internal/service/service_test.go`

**Step 1: Write the failing tests**

```go
func TestCruiseRepository_List_PublicOnlyEnabledCruises(t *testing.T) {
    // seed enabled + disabled cruises
    // assert public listing excludes disabled ones
}

func TestCruiseRepository_List_PublicFiltersByCompanyID(t *testing.T) {
    // seed multiple enabled companies/cruises
    // assert company_id narrows the list correctly
}

func TestAllHandler_PublicCruises_FilterByCompanyID(t *testing.T) {
    // GET /cruises?company_id=2
    // expect only company 2 enabled cruises
}
```

**Step 2: Run tests to verify they fail**

Run: `cd backend; go test ./internal/handler ./internal/repository ./internal/service -run "PublicCruises|FilterByCompanyID|OnlyEnabledCruises" -v`
Expected: FAIL if public listing currently returns disabled cruises or misses company filter guarantees.

**Step 3: Write minimal implementation**

```go
// public cruise path should enforce status=1 when status query is absent
// preserve company_id filtering in the existing service/repository flow
```

**Step 4: Run tests to verify they pass**

Run: `cd backend; go test ./internal/handler ./internal/repository ./internal/service -run "PublicCruises|FilterByCompanyID|OnlyEnabledCruises" -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add backend/internal/handler/all_handler_test.go backend/internal/repository/cruise_repo_test.go backend/internal/service/service_test.go
git commit -m "test(backend): cover public cruise status and company filtering"
```

### Task 4: Implement Public Cruise Filter Behavior

**Files:**
- Modify: `backend/internal/handler/cruise_handler.go`
- Modify: `backend/internal/service/cruise_service.go`
- Modify: `backend/internal/repository/cruise_repo.go`
- Modify: `backend/internal/router/router.go`

**Step 1: Write the failing public contract test if still needed**

```go
func TestAllHandler_PublicCruises_DefaultsToEnabledOnly(t *testing.T) {
    // GET /cruises
    // assert disabled cruises are excluded by default
}
```

**Step 2: Run test to verify it fails**

Run: `cd backend; go test ./internal/handler -run "PublicCruises" -v`
Expected: FAIL until the public-default status rule is explicit.

**Step 3: Write minimal implementation**

```go
// Keep admin list behavior unchanged.
// For public path usage, ensure status defaults to enabled-only semantics.
```

```go
// repository example
if status != nil {
    q = q.Where("status = ?", *status)
} else if publicOnly {
    q = q.Where("status = ?", 1)
}
```

**Step 4: Run tests to verify they pass**

Run: `cd backend; go test ./internal/handler ./internal/repository ./internal/service -run "PublicCruises|FilterByCompanyID|OnlyEnabledCruises" -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add backend/internal/handler/cruise_handler.go backend/internal/service/cruise_service.go backend/internal/repository/cruise_repo.go backend/internal/router/router.go
git commit -m "feat(backend): enforce public cruise filtering rules"
```

### Task 5: Add Failing Miniapp Tests For Wiki Page Behavior

**Files:**
- Create: `frontend/miniapp/tests/unit/pages/wiki/index.spec.ts`
- Modify: `frontend/miniapp/tests/unit/pages/cruise/detail.spec.ts`
- Modify: `frontend/miniapp/src/App.vue`

**Step 1: Write the failing tests**

```ts
it('renders all cruises by default on wiki home', async () => {
  // mock /companies and /cruises
  // expect 全部邮轮 selected
  // expect multiple cruise cards rendered
})

it('filters cruises after clicking a company', async () => {
  // click 皇家加勒比
  // expect second request carries company_id
  // expect filtered cards shown
})

it('navigates to cruise detail after clicking a cruise card', async () => {
  // click 海洋绿洲号
  // expect app shell switched to detail branch with correct cruise id
})
```

**Step 2: Run tests to verify they fail**

Run: `cd frontend/miniapp; npx vitest run tests/unit/pages/wiki/index.spec.ts tests/unit/pages/cruise/detail.spec.ts`
Expected: FAIL because wiki page and navigation wiring do not exist yet.

**Step 3: Write minimal implementation scaffolding**

```ts
// app shell state
const activeTab = ref<TabKey>('home')
const selectedCruiseId = ref<number | null>(null)
```

```vue
<WikiHomePage
  v-else-if="activeTab === 'wiki'"
  @open-cruise="(id) => { selectedCruiseId = id; activeTab = 'cruise-detail' }"
/>
```

**Step 4: Run tests to verify they pass**

Run: `cd frontend/miniapp; npx vitest run tests/unit/pages/wiki/index.spec.ts tests/unit/pages/cruise/detail.spec.ts`
Expected: PASS.

**Step 5: Commit**

```bash
git add frontend/miniapp/tests/unit/pages/wiki/index.spec.ts frontend/miniapp/tests/unit/pages/cruise/detail.spec.ts frontend/miniapp/src/App.vue
git commit -m "test(miniapp): cover wiki home flow and detail navigation"
```

### Task 6: Build Wiki Page UI With Vue And Mobile UX Best Practices

**Files:**
- Create: `frontend/miniapp/pages/wiki/index.vue`
- Create: `frontend/miniapp/components/wiki/WikiCompanySidebar.vue`
- Create: `frontend/miniapp/components/wiki/WikiCruiseList.vue`
- Create: `frontend/miniapp/components/wiki/WikiCruiseCard.vue`
- Modify: `frontend/miniapp/src/App.vue`
- Modify: `frontend/miniapp/src/style.css`

**Step 1: Write or extend a failing UI assertion**

```ts
it('shows company sidebar and cruise cards with reference-like structure', async () => {
  // expect 邮轮百科 title strip
  // expect 全部邮轮 in sidebar
  // expect card image, Chinese name, English name
})
```

**Step 2: Run test to verify it fails**

Run: `cd frontend/miniapp; npx vitest run tests/unit/pages/wiki/index.spec.ts`
Expected: FAIL because the real wiki UI is not implemented.

**Step 3: Write minimal implementation (`<script setup lang="ts">`)**

```vue
<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { request } from '../../src/utils/request'

const emit = defineEmits<{ (e: 'open-cruise', id: number): void }>()

const companies = ref<any[]>([])
const cruises = ref<any[]>([])
const selectedCompanyId = ref<number | 'all'>('all')
const loadingCompanies = ref(false)
const loadingCruises = ref(false)

const sidebarItems = computed(() => [
  { id: 'all', name: '全部邮轮' },
  ...companies.value,
])

async function loadCompanies() {}
async function loadCruises(companyId?: number) {}
</script>
```

**Step 4: Run tests to verify they pass**

Run: `cd frontend/miniapp; npx vitest run tests/unit/pages/wiki/index.spec.ts`
Expected: PASS.

**Step 5: Commit**

```bash
git add frontend/miniapp/pages/wiki/index.vue frontend/miniapp/components/wiki/WikiCompanySidebar.vue frontend/miniapp/components/wiki/WikiCruiseList.vue frontend/miniapp/components/wiki/WikiCruiseCard.vue frontend/miniapp/src/App.vue frontend/miniapp/src/style.css
git commit -m "feat(miniapp): add cruise wiki home page"
```

### Task 7: Fix Cruise Detail Page Invalid Data Dependency

**Files:**
- Modify: `frontend/miniapp/pages/cruise/detail.vue`
- Modify: `frontend/miniapp/tests/unit/pages/cruise/detail.spec.ts`

**Step 1: Write the failing test**

```ts
it('loads cruise detail without requesting unsupported routes endpoint', async () => {
  // mock detail/cabin/facility requests only
  // expect no request('/routes')
})
```

**Step 2: Run test to verify it fails**

Run: `cd frontend/miniapp; npx vitest run tests/unit/pages/cruise/detail.spec.ts`
Expected: FAIL because the page currently requests `/routes`.

**Step 3: Write minimal implementation**

```ts
const [detailRes, typeRes, facilityRes] = await Promise.all([
  request(`/cruises/${id}`),
  request(`/cabin-types?cruise_id=${id}&page=1&page_size=20`),
  request(`/facilities?cruise_id=${id}`),
])

routes.value = []
```

**Step 4: Run test to verify it passes**

Run: `cd frontend/miniapp; npx vitest run tests/unit/pages/cruise/detail.spec.ts`
Expected: PASS.

**Step 5: Commit**

```bash
git add frontend/miniapp/pages/cruise/detail.vue frontend/miniapp/tests/unit/pages/cruise/detail.spec.ts
git commit -m "fix(miniapp): remove unsupported routes dependency from cruise detail"
```

### Task 8: Run Targeted Verification

**Files:**
- Modify if needed: affected files from previous tasks only

**Step 1: Run backend targeted tests**

Run: `cd backend; go test ./internal/handler ./internal/repository ./internal/service -run "PublicCompanies|ListPublic|PublicCruises|FilterByCompanyID|OnlyEnabledCruises" -v`
Expected: PASS.

**Step 2: Run frontend targeted tests**

Run: `cd frontend/miniapp; npx vitest run tests/unit/pages/wiki/index.spec.ts tests/unit/pages/cruise/detail.spec.ts`
Expected: PASS.

**Step 3: Run miniapp preview verification**

Run: `cd frontend/miniapp; npm run dev -- --host 127.0.0.1 --port 3015`
Expected: dev server starts cleanly.

**Step 4: Validate UI and API behavior manually**

Run checks:
- open `http://127.0.0.1:3015`
- switch to `邮轮百科`
- verify default `全部邮轮`
- click one company and verify filtered cards
- click one cruise and verify detail loads without missing API requests

**Step 5: Commit**

```bash
git add <verified changed files>
git commit -m "feat(miniapp): deliver cruise wiki home browsing flow"
```