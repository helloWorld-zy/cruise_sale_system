# Admin Cruise Company Management And Cruise Company Selector Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 在 Admin 后台新增“邮轮公司管理”可用页面，并完成后端与 SQL 保障；同时把邮轮新增/编辑中的“公司ID”改为“所属公司”下拉选择。

**Architecture:** 后端沿用现有 `domain -> repository -> service -> handler -> router` 分层，在现有 `cruise_companies` 基础上增加数据质量 SQL 迁移与查询/更新修正。前端基于 Nuxt3 + Vue3 现有页面模式新增公司管理页面，并在邮轮页加载公司列表作为下拉选项。以 TDD 方式先补失败测试，再写最小实现，确保回归安全。

**Tech Stack:** Go + Gin + GORM + PostgreSQL migration SQL + Nuxt 3 + Vue 3 (`<script setup lang="ts">`) + Vitest + Vue Test Utils

---

## Preconditions

- 参考技能：`@vue-best-practices` `@vue-testing-best-practices` `@superpowers:tdd`
- 在独立 worktree 执行。
- 所有任务遵循 DRY / YAGNI，每个任务完成后立即小提交。

### Task 1: Add SQL Migration For Company Data Quality And Searchability

**Files:**
- Create: `backend/migrations/000015_company_admin_hardening.up.sql`
- Create: `backend/migrations/000015_company_admin_hardening.down.sql`
- Create: `backend/migrations/migrations_sprint43_test.go`
- Test: `backend/migrations/migrations_sprint43_test.go`

**Step 1: Write the failing test**

```go
package migrations

import (
    "os"
    "testing"
)

func TestSprint43MigrationFilesExist(t *testing.T) {
    files := []string{
        "000015_company_admin_hardening.up.sql",
        "000015_company_admin_hardening.down.sql",
    }
    for _, f := range files {
        if _, err := os.Stat(f); err != nil {
            t.Fatalf("expected %s", f)
        }
    }
}
```

**Step 2: Run test to verify it fails**

Run: `cd backend/migrations; go test -run TestSprint43MigrationFilesExist -v`
Expected: FAIL with `no such file or directory` for `000015_*`。

**Step 3: Write minimal implementation (SQL)**

```sql
-- 000015_company_admin_hardening.up.sql
-- 保障公司管理字段质量并提升中英文名称搜索性能。

ALTER TABLE cruise_companies
    ALTER COLUMN name SET NOT NULL;

UPDATE cruise_companies
SET english_name = NULL
WHERE english_name = '';

UPDATE cruise_companies
SET logo_url = NULL
WHERE logo_url = '';

CREATE INDEX IF NOT EXISTS idx_cruise_companies_name ON cruise_companies(name);
CREATE INDEX IF NOT EXISTS idx_cruise_companies_english_name ON cruise_companies(english_name);
```

```sql
-- 000015_company_admin_hardening.down.sql
DROP INDEX IF EXISTS idx_cruise_companies_english_name;
DROP INDEX IF EXISTS idx_cruise_companies_name;
```

**Step 4: Run test to verify it passes**

Run: `cd backend/migrations; go test -run TestSprint43MigrationFilesExist -v`
Expected: PASS。

**Step 5: Commit**

```bash
git add backend/migrations/000015_company_admin_hardening.up.sql backend/migrations/000015_company_admin_hardening.down.sql backend/migrations/migrations_sprint43_test.go
git commit -m "feat(db): add company admin hardening migration"
```

### Task 2: Backend Company Search And Cruise Update Fix

**Files:**
- Modify: `backend/internal/repository/company_repo.go`
- Modify: `backend/internal/service/cruise_service.go`
- Modify: `backend/internal/handler/cruise_handler.go`
- Test: `backend/internal/repository/company_repo_test.go`
- Test: `backend/internal/service/service_test.go`
- Test: `backend/internal/handler/cruise_handler_extended_test.go`

**Step 1: Write failing tests**

```go
func TestCompanyRepository_List_SearchesNameAndEnglishName(t *testing.T) {
    // seed: Name="皇家加勒比", EnglishName="Royal Caribbean"
    // assert keyword "Royal" can hit by english_name
}

func TestCruiseHandler_Update_ChangesCompanyID(t *testing.T) {
    // PUT /cruises/:id with company_id=2
    // assert response data has company_id=2
}
```

**Step 2: Run tests to verify they fail**

Run: `cd backend; go test ./internal/repository ./internal/handler ./internal/service -run "CompanyRepository_List_SearchesNameAndEnglishName|CruiseHandler_Update_ChangesCompanyID" -v`
Expected: FAIL，因为当前公司搜索仅查 `name`，且 `CruiseHandler.Update` 未回填 `existing.CompanyID`。

**Step 3: Write minimal implementation**

```go
// company_repo.go
if keyword != "" {
    like := "%" + keyword + "%"
    q = q.Where("name LIKE ? OR english_name LIKE ?", like, like)
}
```

```go
// cruise_handler.go (Update)
existing.CompanyID = req.CompanyID
```

```go
// cruise_service.go (optional guard for update path)
if _, err := s.companyRepo.GetByID(ctx, cruise.CompanyID); err != nil {
    return err
}
return s.cruiseRepo.Update(ctx, cruise)
```

**Step 4: Run tests to verify they pass**

Run: `cd backend; go test ./internal/repository ./internal/handler ./internal/service -run "CompanyRepository_List_SearchesNameAndEnglishName|CruiseHandler_Update_ChangesCompanyID" -v`
Expected: PASS。

**Step 5: Commit**

```bash
git add backend/internal/repository/company_repo.go backend/internal/service/cruise_service.go backend/internal/handler/cruise_handler.go backend/internal/repository/company_repo_test.go backend/internal/service/service_test.go backend/internal/handler/cruise_handler_extended_test.go
git commit -m "fix(backend): support company english search and cruise company update"
```

### Task 3: Add Backend Company Handler Validation Tests

**Files:**
- Create: `backend/internal/handler/company_handler_test.go`
- Modify: `backend/internal/handler/company_handler.go`

**Step 1: Write failing test**

```go
func TestCompanyHandler_Create_RequiresName(t *testing.T) {
    // POST /companies with empty name
    // expect 400 validation
}

func TestCompanyHandler_Create_AcceptsLogoAndDescription(t *testing.T) {
    // POST /companies with logo_url, english_name, description
    // expect 200 and payload fields returned
}
```

**Step 2: Run test to verify it fails**

Run: `cd backend; go test ./internal/handler -run "CompanyHandler_Create_RequiresName|CompanyHandler_Create_AcceptsLogoAndDescription" -v`
Expected: FAIL（当前缺少该测试文件或校验覆盖不足）。

**Step 3: Write minimal implementation**

```go
type CompanyRequest struct {
    Name        string `json:"name" binding:"required"`
    EnglishName string `json:"english_name"`
    Description string `json:"description"`
    LogoURL     string `json:"logo_url"`
    SortOrder   int    `json:"sort_order"`
}
```

```go
if err := c.ShouldBindJSON(&req); err != nil {
    response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
    return
}
```

**Step 4: Run test to verify it passes**

Run: `cd backend; go test ./internal/handler -run "CompanyHandler_Create_RequiresName|CompanyHandler_Create_AcceptsLogoAndDescription" -v`
Expected: PASS。

**Step 5: Commit**

```bash
git add backend/internal/handler/company_handler.go backend/internal/handler/company_handler_test.go
git commit -m "test(handler): cover company payload validation and fields"
```

### Task 4: Add Admin Company Management Pages

**Files:**
- Create: `frontend/admin/app/pages/companies/index.vue`
- Create: `frontend/admin/app/pages/companies/[id].vue`
- Modify: `frontend/admin/app/layouts/default.vue`
- Test: `frontend/admin/tests/unit/pages/companies-index.spec.ts`
- Test: `frontend/admin/tests/unit/pages/companies-id.spec.ts`
- Test: `frontend/admin/tests/unit/layouts/default.spec.ts`

**Step 1: Write failing tests**

```ts
it('renders company management page title and columns', async () => {
  // expect "邮轮公司管理"
  // expect columns: logo/name/english_name/description
})

it('saves company edit form', async () => {
  // PUT /companies/:id payload contains logo_url, name, english_name, description
})

it('default layout includes companies nav link', () => {
  // expect link to /companies
})
```

**Step 2: Run tests to verify they fail**

Run: `cd frontend/admin; npx vitest run tests/unit/pages/companies-index.spec.ts tests/unit/pages/companies-id.spec.ts tests/unit/layouts/default.spec.ts`
Expected: FAIL（页面与测试尚未存在）。

**Step 3: Write minimal implementation (`<script setup lang="ts">`)**

```vue
<!-- companies/index.vue -->
<script setup lang="ts">
import { onMounted, ref } from 'vue'
const { request } = useApi()
const rows = ref<any[]>([])

async function loadItems() {
  const res = await request('/companies')
  const payload = res?.data ?? res ?? {}
  rows.value = Array.isArray(payload) ? payload : payload?.list ?? []
}

onMounted(loadItems)
</script>
```

```vue
<!-- default.vue: add nav item -->
<NuxtLink class="admin-link" to="/companies">邮轮公司管理</NuxtLink>
```

**Step 4: Run tests to verify they pass**

Run: `cd frontend/admin; npx vitest run tests/unit/pages/companies-index.spec.ts tests/unit/pages/companies-id.spec.ts tests/unit/layouts/default.spec.ts`
Expected: PASS。

**Step 5: Commit**

```bash
git add frontend/admin/app/pages/companies/index.vue frontend/admin/app/pages/companies/[id].vue frontend/admin/app/layouts/default.vue frontend/admin/tests/unit/pages/companies-index.spec.ts frontend/admin/tests/unit/pages/companies-id.spec.ts frontend/admin/tests/unit/layouts/default.spec.ts
git commit -m "feat(admin): add cruise company management pages"
```

### Task 5: Replace Cruise Company ID Input With Company Selector

**Files:**
- Modify: `frontend/admin/app/pages/cruises/create.vue`
- Modify: `frontend/admin/app/pages/cruises/[id].vue`
- Modify: `frontend/admin/app/pages/cruises/index.vue`
- Test: `frontend/admin/tests/unit/pages/cruises-create.spec.ts`
- Test: `frontend/admin/tests/unit/pages/cruises-id.spec.ts`
- Test: `frontend/admin/tests/unit/pages/cruises-index.spec.ts`

**Step 1: Write failing tests**

```ts
it('create page shows 所属公司 select and loads options from /companies', async () => {
  // assert request('/companies') called
  // assert label "所属公司"
})

it('edit page updates selected company_id on save', async () => {
  // select company option then submit
  // assert PUT body.company_id equals selected option value
})

it('list page shows company name instead of raw company_id when company object exists', async () => {
  // item.company.name renders in table
})
```

**Step 2: Run tests to verify they fail**

Run: `cd frontend/admin; npx vitest run tests/unit/pages/cruises-create.spec.ts tests/unit/pages/cruises-id.spec.ts tests/unit/pages/cruises-index.spec.ts`
Expected: FAIL（当前是数字输入“公司 ID”）。

**Step 3: Write minimal implementation**

```vue
<script setup lang="ts">
const companies = ref<Array<{ id: number; name: string }>>([])

async function loadCompanies() {
  const res = await request('/companies', { query: { page: 1, page_size: 200 } })
  const payload = res?.data ?? res ?? {}
  companies.value = (Array.isArray(payload) ? payload : payload?.list ?? []).map((it: any) => ({
    id: Number(it.id),
    name: it.name || `公司#${it.id}`,
  }))
}
</script>

<label>所属公司
  <select v-model.number="form.company_id">
    <option :value="0" disabled>请选择所属公司</option>
    <option v-for="c in companies" :key="c.id" :value="c.id">{{ c.name }}</option>
  </select>
</label>
```

```vue
<!-- cruises/index.vue -->
<td class="p-3 text-slate-600">{{ item.company?.name || item.company_id || '-' }}</td>
```

**Step 4: Run tests to verify they pass**

Run: `cd frontend/admin; npx vitest run tests/unit/pages/cruises-create.spec.ts tests/unit/pages/cruises-id.spec.ts tests/unit/pages/cruises-index.spec.ts`
Expected: PASS。

**Step 5: Commit**

```bash
git add frontend/admin/app/pages/cruises/create.vue frontend/admin/app/pages/cruises/[id].vue frontend/admin/app/pages/cruises/index.vue frontend/admin/tests/unit/pages/cruises-create.spec.ts frontend/admin/tests/unit/pages/cruises-id.spec.ts frontend/admin/tests/unit/pages/cruises-index.spec.ts
git commit -m "feat(admin): use company selector in cruise create and edit"
```

### Task 6: End-To-End Validation And Swagger Consistency

**Files:**
- Modify: `backend/docs/swagger.yaml`
- Modify: `backend/docs/swagger.json`
- Optional Modify: `backend/internal/handler/company_handler.go` (注释一致性)

**Step 1: Write failing check (manual contract check)**

```bash
# 验证公司接口字段在 Swagger 中可见
rg "english_name|logo_url|description" backend/docs/swagger.yaml
```

**Step 2: Run check to verify it fails or is stale**

Run: `cd backend; rg "english_name|logo_url|description" docs/swagger.yaml`
Expected: 若缺失则需要刷新文档。

**Step 3: Write minimal implementation**

```bash
cd backend
swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal
```

**Step 4: Run verification tests**

Run: `cd backend; go test ./internal/handler ./internal/service ./internal/repository -v`
Expected: PASS。

Run: `cd frontend/admin; npx vitest run tests/unit/pages/companies-index.spec.ts tests/unit/pages/companies-id.spec.ts tests/unit/pages/cruises-create.spec.ts tests/unit/pages/cruises-id.spec.ts tests/unit/pages/cruises-index.spec.ts tests/unit/layouts/default.spec.ts`
Expected: PASS。

**Step 5: Commit**

```bash
git add backend/docs/swagger.yaml backend/docs/swagger.json
git commit -m "docs(api): refresh swagger for company management and cruise selector"
```

## Final Verification Checklist

1. 执行迁移并确认 `cruise_companies` 数据可正常读写：
`cd backend; make migrate-up`
2. 启动后端并验证接口：
`cd backend; go run ./cmd/server`
3. 启动 admin 前端并手测：
`cd frontend/admin; npm run dev`
4. 手工验收：
- 左侧菜单出现“邮轮公司管理”。
- 公司管理支持新增/编辑/删除，字段含 `logo_url`、中文名、英文名、介绍。
- 邮轮新增/编辑中不再出现“公司ID”数字输入，而是“所属公司”下拉。
- 下拉数据来自公司列表 API。

## Notes

- 若历史数据存在公司名为空，先做数据修复再加 `NOT NULL`。
- `cruise_handler.go` 的更新逻辑必须包含 `existing.CompanyID = req.CompanyID`，否则前端选择不会生效。
- 保持 Vue 页面使用 Composition API 和 `<script setup lang="ts">`。
