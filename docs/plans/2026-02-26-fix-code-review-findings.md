# Code Review Findings 修复 Prompt

> **用途**：供大模型在单次或多次会话中，按优先级依次修复 `2026-02-26-code-review-findings-sprint1-4.md` 中报告的全部问题。
> **前提**：执行本 Prompt 前，大模型必须先阅读 `docs/plans/2026-02-26-code-review-findings-sprint1-4.md` 理解全部问题上下文。

---

## 角色设定

**Role**: 你是一名极其严格的资深全栈工程师（Staff Engineer），擅长 Go（Gin + GORM + DDD）和 Vue 3（Nuxt 4 + `<script setup lang="ts">` + Pinia Composition API）。你的任务是逐一修复 Code Review 报告中发现的全部问题，确保每个修复都：

1. **遵循 TDD**：先阅读已有测试（如有），理解预期行为，然后修改代码，最后补充/重写测试并运行验证通过。
2. **符合现有架构**：后端遵循 DDD 分层（domain → repository → service → handler → router）；前端使用项目已建立的 `useApi()` composable（admin/web）或 `request()` 工具（miniapp）。
3. **代码风格一致**：参照同层已有文件的命名规范、注释风格、错误处理模式。
4. **100% 覆盖率目标**：修复完成后运行测试确认覆盖率提升至 100%。

---

## 执行流程（必须严格按顺序）

```
Phase 1: 迁移链修复（阻塞部署，最高优先级）
Phase 2: 后端测试覆盖率提升至 100%
Phase 3: 后端代码质量修复（DI 接口化）
Phase 4: 前端空壳页面实联修复（Sprint 1/2）
Phase 5: 前端空壳页面实联修复（Sprint 3/4）
Phase 6: 前端伪覆盖测试重写 + 新增页面测试
Phase 7: 前端代码质量修复（Pinia 风格统一、Token 管理统一）
Phase 8: 全量测试验证
```

---

## Phase 1: 迁移链修复

### 问题描述（S4-10）

当前迁移目录 `backend/migrations/` 中存在编号冲突：
- `000004_user_booking.up.sql`（Sprint 3）
- `000004_payment_notify.up.sql`（Sprint 4）

两个文件共享 `000004` 前缀，golang-migrate 无法确定执行顺序。后续的 `000005_cabin_hold_unique.up.sql` 依赖 `000004_user_booking` 中的 `cabin_holds` 表。

### 修复指令

1. 将 Sprint 3 的迁移保持 `000004` 不变（因为它创建了被后续迁移依赖的表）：
   - `000004_user_booking.up.sql` → 保持不变
   - `000004_user_booking.down.sql` → 保持不变

2. 将 Sprint 4 的迁移重命名为 `000005`：
   - `000004_payment_notify.up.sql` → 重命名为 `000005_payment_notify.up.sql`
   - `000004_payment_notify.down.sql` → 重命名为 `000005_payment_notify.down.sql`

3. 将原 `000005_cabin_hold_unique` 重命名为 `000006`：
   - `000005_cabin_hold_unique.up.sql` → 重命名为 `000006_cabin_hold_unique.up.sql`
   - `000005_cabin_hold_unique.down.sql` → 重命名为 `000006_cabin_hold_unique.down.sql`

4. 检查并更新以下测试文件中对迁移编号的引用：
   - `backend/migrations/migrations_sprint4_test.go` — 搜索 `000004_payment` 或 `000005_cabin` 字符串并替换
   - `backend/migrations/migrations_sprint3_test.go` — 确认 `000004_user_booking` 引用不变
   - `backend/migrations/migrations_sprint2_test.go` — 确认无影响

5. **验证**：正则搜索整个 `backend/` 目录中所有对旧编号的引用（`000004_payment` / `000005_cabin`），确保全部更新。

---

## Phase 2: 后端测试覆盖率提升至 100%

### 2.1 Repository 层集成测试（0% → 100%）

以下 4 个 Repository 文件当前 0% 覆盖率，需要编写 SQLite 内存数据库的集成测试。

**参考模式**（来自 `repository/company_repo_test.go`）：
```go
func TestXxxRepository_Method(t *testing.T) {
    db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
    if err != nil {
        t.Fatalf("db error: %v", err)
    }
    if err := db.AutoMigrate(&domain.XxxModel{}); err != nil {
        t.Fatalf("migrate error: %v", err)
    }
    repo := NewXxxRepository(db)
    // ... test logic with real DB operations
}
```

#### 2.1.1 `analytics_repo.go` — 创建 `analytics_repo_test.go`

需要测试的方法：
- `TodaySales(ctx)` — 注意：原始 SQL 使用 PostgreSQL 的 `CURRENT_DATE` 和 `generate_series`，SQLite 不支持。**处理方式**：
  - 方案 A（推荐）：重构 `AnalyticsRepository` 的 SQL 查询使其兼容 SQLite（使用 `date('now')` 替代 `CURRENT_DATE`），或者
  - 方案 B：在测试中使用 GORM 的 `db.Exec()` 手动创建简化的 `payments` 和 `bookings` 表，然后测试 `TodaySales` 和 `TodayOrderCount` 的 SQLite 兼容版本。
  - 方案 C：如果查询无法兼容 SQLite，则为该 repo 编写一个 mock 测试，测试 Repository 的构造函数和方法签名至少被调用。但**必须**标注 `// TODO: 需要 PostgreSQL 集成测试环境验证 SQL 语法`。
- `WeeklyTrend(ctx)` — 同上，`generate_series` 是 PG 特有函数
- `TodayOrderCount(ctx)` — 相对简单，`COUNT(*)` + `WHERE created_at >= date('now')` 在 SQLite 中可用

**测试用例清单**：
- `TestAnalyticsRepository_TodaySales_NoPayments` — 无数据时返回 0
- `TestAnalyticsRepository_TodaySales_WithPaidPayments` — 有已支付记录时返回正确总额
- `TestAnalyticsRepository_TodayOrderCount_NoBookings` — 无数据时返回 0
- `TestAnalyticsRepository_TodayOrderCount_WithBookings` — 有预订时返回正确数量
- `TestAnalyticsRepository_WeeklyTrend` — 返回 7 天数据（如果 SQLite 不支持则用 mock）

#### 2.1.2 `notification_repo.go` — 创建 `notification_repo_test.go`

需要 AutoMigrate 的 Domain 模型：`domain.Notification`

**测试用例清单**：
- `TestNotificationRepository_CreateOutbox` — 创建一条 pending 通知
- `TestNotificationRepository_ListPending` — 返回 pending 状态的通知，按 created_at 升序
- `TestNotificationRepository_ListPending_Limit` — 验证 limit 参数生效
- `TestNotificationRepository_MarkSent` — 标记为 sent 后，ListPending 不再返回
- `TestNotificationRepository_MarkFailed` — 标记为 failed 后，ListPending 不再返回

#### 2.1.3 `payment_repo.go` — 创建 `payment_repo_test.go`

需要 AutoMigrate 的 Domain 模型：`domain.Payment`

**测试用例清单**：
- `TestPaymentRepository_Create` — 创建支付记录并验证 ID 自增
- `TestPaymentRepository_FindByTradeNo_Found` — 根据 trade_no 查找
- `TestPaymentRepository_FindByTradeNo_NotFound` — 找不到时返回错误
- `TestPaymentRepository_FindByID_Found` — 根据 ID 查找
- `TestPaymentRepository_FindByID_NotFound` — 找不到时返回错误
- `TestPaymentRepository_UpdateStatus` — 更新状态后再次查询验证

#### 2.1.4 `refund_repo.go` — 创建 `refund_repo_test.go`

需要 AutoMigrate 的 Domain 模型：`domain.Refund`（注意：`Refund` 的 `payment_id` 外键引用 `payments` 表，需要同时 Migrate `domain.Payment`）

**测试用例清单**：
- `TestRefundRepository_Create` — 创建退款记录
- `TestRefundRepository_SumByPaymentID_NoRefunds` — 无退款时返回 0
- `TestRefundRepository_SumByPaymentID_WithRefunds` — 有退款时返回正确总额
- `TestRefundRepository_SumByPaymentID_ExcludesCancelled` — 取消状态的退款不计入总额

### 2.2 `booking_repo.go` L37 — 补充 `FindBookingByID` 测试（S3-04）

在已有的 `booking_repo_test.go` 或 `booking_repo_blackbox_test.go` 中补充：

**注意**：先检查 `booking_repo.go` 中是否真的有 `FindBookingByID` 方法。根据当前文件内容（只有 `Create`、`InTx`、`UpdateStatus`），如果该方法不存在，则需要确认是否在其他测试文件中已测试了所有方法。如果所有现有方法已被测试，则 S3-04 可能是误报——需要确认后再决定。

### 2.3 Handler 错误分支补充（S4-08）

#### `analytics_handler_test.go` — 补充 `WeeklyTrend` 和 `TodayOrderCount` 的独立错误测试

当前测试中 `fakeAnalyticsSvc` 的 `err` 字段是共享的——当 `err != nil` 时，`TodaySales` 先返回错误，后续的 `WeeklyTrend` 和 `TodayOrderCount` 不会被调用。需要创建可以让不同方法独立返回错误的 fake：

```go
type fakeAnalyticsSvcSelective struct {
    sales       int64
    trend       []int64
    orders      int64
    salesErr    error
    trendErr    error
    ordersErr   error
}

func (f *fakeAnalyticsSvcSelective) TodaySales(_ context.Context) (int64, error) {
    return f.sales, f.salesErr
}
func (f *fakeAnalyticsSvcSelective) WeeklyTrend(_ context.Context) ([]int64, error) {
    return f.trend, f.trendErr
}
func (f *fakeAnalyticsSvcSelective) TodayOrderCount(_ context.Context) (int64, error) {
    return f.orders, f.ordersErr
}
```

新增测试：
- `TestAnalyticsHandler_Summary_WeeklyTrendError` — `salesErr = nil, trendErr = errors.New("trend fail")` → 验证返回 500
- `TestAnalyticsHandler_Summary_OrderCountError` — `salesErr = nil, trendErr = nil, ordersErr = errors.New("orders fail")` → 验证返回 500

### 2.4 Service 层错误分支补充（S4-08）

#### `payment_service.go` — 3 个未覆盖分支

在 `payment_service_test.go` 中补充：

1. `PaymentService.Create` 中 `s.payRepo.Create()` 失败的分支 → mock payRepo.Create 返回错误
2. `PaymentCallbackServiceImpl.HandleCallback` 中 `v.ExtractTradeNo(body)` 失败的分支 → mock Verifier.Verify 成功但 ExtractTradeNo 返回错误
3. `PaymentCallbackServiceImpl.HandleCallback` 中 `s.payRepo.FindByTradeNo()` 失败的分支 → mock FindByTradeNo 返回错误

**检查方式**：运行 `go test ./internal/service/ -coverprofile=svc_cov.out -covermode=atomic && go tool cover -func=svc_cov.out | grep payment_service` 确认 100%。

#### `refund_service.go` — `amountCents <= 0` 分支

在 `refund_service_test.go` 中补充：
- `TestRefundService_Create_ZeroAmount` — `amountCents = 0` → 返回 "refund amount must be positive"
- `TestRefundService_Create_NegativeAmount` — `amountCents = -100` → 返回 "refund amount must be positive"

### 2.5 `service/chk/check.go` 排除出 Coverage（S4-09）

该文件是 `package main` 的开发辅助工具，不属于 Sprint 交付物。处理方式（任选其一）：

- **方案 A（推荐）**：将 `service/chk/` 目录移到 `backend/tools/chk/` 或 `backend/cmd/chk/`，使其不在 `internal/service/` 下被 `go test ./internal/...` 扫描到。
- **方案 B**：在 `check.go` 文件顶部添加 `//go:build ignore` 构建标签。
- **方案 C**：为 `check.go` 编写一个简单的测试来覆盖它（但这是一个 `main` 包，测试意义不大）。

### 2.6 `domain/notification.go` — `TableName()` 方法覆盖

检查 `domain/notification.go` 中是否有 `TableName()` 方法，如有，在 `domain/notification_test.go` 中添加一行测试：
```go
func TestNotification_TableName(t *testing.T) {
    n := Notification{}
    assert.Equal(t, "notifications", n.TableName())
}
```

### Phase 2 验证

```bash
cd backend && go test ./... -coverprofile=coverage.out -covermode=atomic && go tool cover -func=coverage.out
```

确认所有包覆盖率 ≥ 99%（`domain` 包因无可执行语句除外）。

---

## Phase 3: 后端代码质量修复

### 3.1 `AuthService` 接口化 DI（S1-03）

**当前代码**（`service/auth_service.go`）：
```go
type AuthService struct {
    staffRepo   *repository.StaffRepository  // ← 具体类型
    ...
}
func NewAuthService(staffRepo *repository.StaffRepository, ...) *AuthService
```

**修复步骤**：

1. 在 `domain/repository.go` 中添加 `StaffRepository` 接口（如果尚未存在）：
```go
type StaffRepository interface {
    Create(ctx context.Context, staff *Staff) error
    GetByUsername(ctx context.Context, username string) (*Staff, error)
    GetByID(ctx context.Context, id int64) (*Staff, error)
    Update(ctx context.Context, staff *Staff) error
    Delete(ctx context.Context, id int64) error
}
```

2. 修改 `service/auth_service.go`：
```go
import (
    // 移除 "github.com/cruisebooking/backend/internal/repository"
    "github.com/cruisebooking/backend/internal/domain"
)

type AuthService struct {
    staffRepo   domain.StaffRepository  // ← 改为接口
    ...
}

func NewAuthService(staffRepo domain.StaffRepository, ...) *AuthService
```

3. 确认 `repository.StaffRepository` struct 已隐式实现 `domain.StaffRepository` 接口（Go 鸭子类型）。

4. 更新 `cmd/server/main.go` 中的 DI 代码（如需要，因为 `*repository.StaffRepository` 会自动满足 `domain.StaffRepository` 接口）。

5. 更新 `service/auth_service_test.go`：如果测试中直接构造了 `repository.StaffRepository`，改为使用 mock/fake 实现 `domain.StaffRepository` 接口。

6. **验证**：`go build ./...` 确认编译通过，`go test ./internal/service/...` 确认测试通过。

---

## Phase 4: 前端空壳页面实联修复（Sprint 1/2）

### 核心原则

每个页面修复必须包含以下要素：
1. 使用 `useApi()` 发起真实 API 请求（admin/web）或 `request()` 发起请求（miniapp）
2. 三态处理：`loading`（加载中）、`error`（错误信息）、`empty`（空数据提示）
3. 保持 `<script setup lang="ts">` 风格

**参照模板**（来自 `admin/app/pages/bookings/index.vue`）：

```vue
<script setup lang="ts">
import { ref, onMounted } from 'vue'
declare const useApi: any

const { request } = useApi()
const items = ref<any[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

onMounted(async () => {
  loading.value = true
  try {
    const res = await request('/xxx')
    items.value = res?.data ?? []
  } catch (e: any) {
    error.value = e?.message ?? 'failed to load'
  } finally {
    loading.value = false
  }
})
</script>
<template>
  <div class="page">
    <h1>Title</h1>
    <p v-if="loading">Loading...</p>
    <p v-else-if="error" class="error">{{ error }}</p>
    <p v-else-if="items.length === 0">No data</p>
    <div v-else><!-- data rendering --></div>
  </div>
</template>
```

### 4.1 `admin/app/pages/login.vue`（S1-01）

**当前问题**：使用 `setTimeout` 模拟登录。

**修复**：
- 将 `setTimeout` 逻辑替换为调用 `useApi().request('/admin/auth/login', { method: 'POST', body: { username, password } })`
- 成功时将返回的 token 存入 `useAuthStore().setToken(res.token)`
- 失败时显示错误信息
- 提交时添加 loading 状态禁用按钮防连点

### 4.2 `admin/app/pages/cruises/index.vue`（S1-02）

**当前问题**：硬编码 `["示例邮轮"]`。

**修复**：
- `onMounted` 中调用 `request('/cruises')` 获取邮轮列表
- 渲染真实数据（id, name, company_id 等）
- 添加 loading/error/empty 三态

### 4.3 `admin/app/pages/routes/index.vue`（S2-03）

**当前问题**：硬编码 `[{ id: 1, code: 'R1', name: 'Route 1' }]`。

**修复**：
- `onMounted` 中调用 `request('/routes')` 获取航线列表
- 使用 `RouteTable` 组件渲染真实数据
- 添加 loading/error/empty 三态

### 4.4 `admin/app/pages/voyages/index.vue`（S2-04）

**当前问题**：硬编码数据。

**修复**：
- `onMounted` 中调用 `request('/voyages')` 获取航次列表
- 渲染真实数据（表格显示 id, code, name, route_id, departure_date 等）
- 添加 loading/error/empty 三态

### 4.5 `admin/app/pages/cabins/index.vue`（S2-01）

**当前问题**：仅有 `<h1>Cabins</h1>` 空壳。

**修复**：
- `onMounted` 中调用 `request('/cabins')` 获取舱位 SKU 列表
- 渲染舱位列表（表格显示 id, voyage_id, cabin_type_id, total, available 等）
- 添加 loading/error/empty 三态
- 可选：添加按航次筛选的 select

### 4.6 `admin/app/pages/cabins/inventory.vue`（S2-02）

**当前问题**：仅有 `<h1>Inventory</h1>` 空壳。

**修复**：
- 使用 route query 参数获取 `cabinSkuId`（如 `/cabins/inventory?skuId=1`）
- 调用 `request('/cabins/${skuId}/inventory')` 获取库存信息
- 显示当前库存数量（total, available）
- 提供调整表单：输入 delta 数值，调用 `request('/cabins/${skuId}/inventory/adjust', { method: 'POST', body: { delta } })`
- 添加 loading/error 三态

### 4.7 `admin/app/pages/cabins/pricing.vue`（S2-05）

**当前问题**：硬编码价格行。

**修复**：
- 使用 route query 参数获取 `cabinSkuId`
- 调用 `request('/cabins/${skuId}/prices')` 获取价格列表
- 使用 `PricingRow` 组件渲染真实数据
- 提供新增/编辑价格的表单，调用 `request('/cabins/${skuId}/prices', { method: 'POST', body: { date, price_cents } })`
- 添加 loading/error/empty 三态

### 4.8 `admin/app/pages/routes/new.vue`（S2-06）

**当前问题**：表单无提交逻辑。

**修复**：
- 添加 `handleSubmit` 函数，调用 `request('/routes', { method: 'POST', body: { code, name } })`
- 成功后跳转到 `/routes` 列表页（`navigateTo('/routes')`）
- 失败时显示错误信息
- 提交时 loading 状态禁用按钮防连点

### 4.9 `admin/app/pages/voyages/new.vue`（S2-07）

**当前问题**：表单无提交逻辑。

**修复**：
- 添加 `handleSubmit` 函数，调用 `request('/voyages', { method: 'POST', body: { code, name, route_id, departure_date, arrival_date } })`
- 成功后跳转到 `/voyages`
- 失败时显示错误信息
- 提交时 loading 状态禁用按钮防连点

### 4.10 `miniapp/pages/cabin/detail.vue`（S2-09）

**当前问题**：仅有标题和 TODO 注释。

**修复**：
- 从页面参数获取 `cabinSkuId`（uni-app: `onLoad(options)` 或 defineProps）
- 调用 `request('/cabins/${cabinSkuId}')` 获取舱位详情
- 渲染舱房信息：名称、描述、价格、库存
- 调用 `request('/cabins/${cabinSkuId}/prices')` 获取价格日历
- 添加 loading/error 三态

---

## Phase 5: 前端空壳页面实联修复（Sprint 3/4）

### 5.1 `miniapp/pages/login/login.vue`（S3-01）

**当前问题**：仅显示 `<text>Login</text>` 和一个按钮，无任何登录逻辑。

**参考实现**（来自 `web/components/LoginForm.vue`）：完整的 SMS 登录流程，包含手机号验证、发送验证码、倒计时、登录 API 调用。

**修复**：
- 添加手机号输入框 + 验证码输入框
- 「发送验证码」按钮调用 `request('/users/sms-code', { method: 'POST', data: { phone } })`
- 添加 60 秒倒计时防重发
- 「登录」按钮调用 `request('/users/login', { method: 'POST', data: { phone, code } })`
- 成功后将 token 存入 auth store（`useAuthStore().setToken(res.token)`）并跳转
- 失败时显示错误信息
- 手机号格式验证（11 位数字）

### 5.2 `web/pages/booking/index.vue`（S3-02）

**当前问题**：预订第一步仅做本地表单跳转，未验证航次/舱位有效性。

**修复**：
- 在 `router.push` 前调用 API 验证航次和舱位的有效性：
  ```ts
  const res = await request(`/voyages/${voyageId}`)
  // 验证航次存在且未过期
  const cabinRes = await request(`/cabins/${cabinSkuId}`)
  // 验证舱位存在且有库存
  ```
- 验证失败时显示错误信息，不跳转
- 添加 loading 状态

### 5.3 `admin/pages/dashboard/index.vue`（S4-01）

**当前问题**：硬编码 `const summary = { sales: 1000, orders: 12 }`。

**修复**：
- `onMounted` 中调用 `request('/admin/analytics/summary')` 获取真实数据
- 后端返回格式：`{ today_sales: number, weekly_trend: number[], today_orders: number }`
- 使用 `StatCard` 组件渲染各指标
- 添加 loading/error 三态

**注意**：此页面位于 `frontend/admin/pages/dashboard/index.vue`（非 `app/pages/`），确认 Nuxt 路由配置是否扫描了此目录。如果 Nuxt 只扫描 `app/pages/`，则可能需要将文件移到 `app/pages/dashboard/index.vue`。先检查 `nuxt.config.ts` 中的 pages 配置。无论在哪个目录，页面内部逻辑修复方式相同。

### 5.4 `admin/pages/finance/index.vue`（S4-02）

**当前问题**：硬编码示例表格行。

**修复**：
- `onMounted` 中调用 `request('/admin/analytics/summary')` 或专用财务 API（检查 router.go 中是否有财务专用路由）
- 如果没有专用财务 API，可以暂时复用 analytics summary 数据，但需添加 `// TODO: 替换为专用财务 API`
- 渲染支付流水列表（如果有对应的 API endpoint）
- 添加 loading/error/empty 三态

### 5.5 `web/pages/pay/[id].vue`（S4-03）

**当前问题**：仅导入 PayButton 组件，未获取订单信息。

**修复**：
- 使用 `useRoute()` 获取订单 ID
- `onMounted` 中调用 `request('/bookings/${id}')` 获取订单信息（金额、状态等）
- 在页面中显示订单金额、订单号等信息
- 将 `amountCents` 和 `bookingId` 作为 props 传递给 `PayButton`
- 添加 loading/error 三态

### 5.6 `web/components/PayButton.vue`（S4-04）

**当前问题**：仅有 `<button>Pay Now</button>` 骨架。

**修复**：
```vue
<script setup lang="ts">
import { ref } from 'vue'
declare const useApi: any

const props = defineProps<{
  bookingId: number
  amountCents: number
}>()

const emit = defineEmits<{
  (e: 'paid', payUrl: string): void
  (e: 'error', msg: string): void
}>()

const { request } = useApi()
const loading = ref(false)

const handlePay = async () => {
  if (loading.value) return  // 防连点
  loading.value = true
  try {
    const res = await request('/payments', {
      method: 'POST',
      body: {
        order_id: props.bookingId,
        amount_cents: props.amountCents,
        provider: 'wechat'
      }
    })
    emit('paid', res.pay_url)
  } catch (e: any) {
    emit('error', e?.message ?? 'payment failed')
  } finally {
    loading.value = false
  }
}
</script>
<template>
  <button @click="handlePay" :disabled="loading">
    {{ loading ? 'Processing...' : 'Pay Now' }}
  </button>
</template>
```

### 5.7 `miniapp/pages/pay/pay.vue`（S4-05）

**当前问题**：仅导入 PayButton，无支付逻辑。

**修复**：
- 从页面参数获取 booking ID
- 调用 `request('/bookings/${id}')` 获取订单信息
- 显示订单金额
- 点击支付按钮时调用 `request('/payments', { method: 'POST', data: { order_id, amount_cents, provider: 'wechat' } })`
- 成功后使用小程序支付能力（`uni.requestPayment` 或跳转支付 URL）
- 添加 loading/error 三态

### 5.8 `miniapp/components/PayButton.vue`

**当前问题**：仅有 `<text>Pay Now</text>`。

**修复**：与 web 端 PayButton 类似，但使用 miniapp 的 `request()` 工具和 `uni.requestPayment`。

---

## Phase 6: 前端伪覆盖测试重写 + 新增页面测试

### 核心原则

每个页面测试必须包含以下断言类型：
1. **API 调用断言**：验证组件挂载后调用了正确的 API endpoint
2. **数据渲染断言**：验证 API 返回的数据正确渲染到 DOM
3. **loading/error 状态断言**：验证 loading 和 error 文本正确显示
4. **用户交互断言**（如有表单/按钮）：验证点击/提交触发正确的 API 调用

**参照模板**（来自 `admin/tests/unit/bookings.list.spec.ts`）：
```ts
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../app/pages/xxx/index.vue'

const mockRequest = vi.fn().mockResolvedValue({ data: [...] })
vi.stubGlobal('useApi', () => ({ request: mockRequest }))
beforeEach(() => mockRequest.mockClear())

describe('XxxPage', () => {
    it('调用 API 获取数据', async () => {
        mount(Page)
        await flushPromises()
        expect(mockRequest).toHaveBeenCalledWith('/xxx')
    })

    it('渲染返回的数据', async () => {
        const wrapper = mount(Page)
        await flushPromises()
        expect(wrapper.text()).toContain('expectedData')
    })

    it('显示错误信息', async () => {
        mockRequest.mockRejectedValueOnce(new Error('fail'))
        const wrapper = mount(Page)
        await flushPromises()
        expect(wrapper.text()).toContain('fail')
    })
})
```

### 6.1 删除/重写 `sprint2-pages.spec.ts`（S2-08）

**当前问题**：仅验证文本渲染和 input 值，无 API 调用断言。

**修复**：删除此文件，为每个 Sprint 2 页面编写独立的测试文件：
- `admin/tests/unit/pages/cabins-index.spec.ts` — 测试 cabins 列表的 API 调用 + 数据渲染 + error 处理
- `admin/tests/unit/pages/cabins-inventory.spec.ts` — 测试库存 API 调用 + 调整表单提交
- `admin/tests/unit/pages/cabins-pricing.spec.ts`（如果已有 `cabins.pricing.spec.ts`，则检查并增强）
- `admin/tests/unit/pages/routes-index.spec.ts`（如果已有 `routes.list.spec.ts`，则检查并增强）
- `admin/tests/unit/pages/routes-new.spec.ts` — 测试表单提交的 API 调用
- `admin/tests/unit/pages/voyages-index.spec.ts` — 测试 voyages 列表的 API 调用
- `admin/tests/unit/pages/voyages-new.spec.ts` — 测试表单提交的 API 调用

### 6.2 重写 `dashboard.page.spec.ts`（S4-06）

**当前问题**：仅 1 个测试 `toContain('Dashboard')`。

**重写**：
```ts
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../pages/dashboard/index.vue'

const mockRequest = vi.fn().mockResolvedValue({
    today_sales: 50000,
    weekly_trend: [1000, 2000, 3000, 4000, 5000, 6000, 7000],
    today_orders: 12
})
vi.stubGlobal('useApi', () => ({ request: mockRequest }))
beforeEach(() => mockRequest.mockClear())

describe('Dashboard', () => {
    it('调用 analytics summary API', async () => {
        mount(Page)
        await flushPromises()
        expect(mockRequest).toHaveBeenCalledWith('/admin/analytics/summary')
    })

    it('渲染真实统计数据', async () => {
        const wrapper = mount(Page)
        await flushPromises()
        // 根据实际渲染方式调整断言
        expect(wrapper.text()).toContain('50000')  // 或格式化后的金额
        expect(wrapper.text()).toContain('12')
    })

    it('API 失败时显示错误', async () => {
        mockRequest.mockRejectedValueOnce(new Error('api error'))
        const wrapper = mount(Page)
        await flushPromises()
        expect(wrapper.text()).toContain('api error')
    })

    it('加载时显示 Loading', () => {
        const wrapper = mount(Page)
        expect(wrapper.text()).toContain('Loading')
    })
})
```

### 6.3 重写 `pay.page.spec.ts`（S4-07）

**当前问题**：仅 1 个测试 `toContain('Pay Now')`。

**重写**：
```ts
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../pages/pay/[id].vue'

const mockRequest = vi.fn()
vi.stubGlobal('useApi', () => ({ request: mockRequest }))
vi.stubGlobal('useRoute', () => ({ params: { id: '42' } }))
beforeEach(() => {
    mockRequest.mockClear()
    mockRequest.mockResolvedValue({ data: { id: 42, total_cents: 19900, status: 'created' } })
})

describe('Pay Page', () => {
    it('加载订单信息', async () => {
        mount(Page)
        await flushPromises()
        expect(mockRequest).toHaveBeenCalledWith(expect.stringContaining('/bookings/42'))
    })

    it('显示订单金额', async () => {
        const wrapper = mount(Page)
        await flushPromises()
        expect(wrapper.text()).toContain('19900')  // 或格式化
    })

    it('点击支付按钮触发支付 API', async () => {
        mockRequest
            .mockResolvedValueOnce({ data: { id: 42, total_cents: 19900 } })  // 订单信息
            .mockResolvedValueOnce({ pay_url: 'https://pay.example.com' })     // 支付创建
        const wrapper = mount(Page)
        await flushPromises()
        await wrapper.find('button').trigger('click')
        await flushPromises()
        expect(mockRequest).toHaveBeenCalledTimes(2)
    })

    it('支付失败时显示错误', async () => {
        mockRequest
            .mockResolvedValueOnce({ data: { id: 42, total_cents: 19900 } })
            .mockRejectedValueOnce(new Error('payment failed'))
        const wrapper = mount(Page)
        await flushPromises()
        await wrapper.find('button').trigger('click')
        await flushPromises()
        expect(wrapper.text()).toContain('payment failed')
    })
})
```

### 6.4 新增 Sprint 1 页面测试

- `admin/tests/unit/pages/login.spec.ts` — 如果已存在则增强，确保测试真实 API 调用（非 setTimeout）
- `admin/tests/unit/pages/cruises-index.spec.ts` — 如果已存在则增强，确保测试真实 API 调用

### 6.5 新增 Sprint 3 页面测试

- `miniapp/tests/unit/pages/login.spec.ts` — 测试完整 SMS 登录流程
- `web/tests/unit/pages/booking/index.spec.ts` — 如果已存在则增强，确保验证 API 调用

### 6.6 新增 Sprint 4 页面测试

- `admin/tests/unit/pages/finance-index.spec.ts` — 测试财务页面 API 调用
- `miniapp/tests/unit/pages/pay.spec.ts` — 如果已存在则增强，确保测试支付流程
- `web/tests/unit/components/pay-button.spec.ts` — 测试 PayButton 组件的支付 API 调用

---

## Phase 7: 前端代码质量修复

### 7.1 `web/app/stores/cruise.ts` 统一为 Composition API（S1-04）

**当前**（Options API）：
```ts
export const useCruiseStore = defineStore('cruise', {
    state: () => ({ list: [] as any[], detail: null as any }),
    actions: {
        setList(list: any[]) { this.list = list },
        setDetail(detail: any) { this.detail = detail },
    },
})
```

**修复**（Setup Store / Composition API）：
```ts
import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useCruiseStore = defineStore('cruise', () => {
    const list = ref<any[]>([])
    const detail = ref<any>(null)

    function setList(newList: any[]) { list.value = newList }
    function setDetail(newDetail: any) { detail.value = newDetail }

    return { list, detail, setList, setDetail }
})
```

**验证**：运行 `web/tests/unit/stores/cruise.spec.ts` 确认测试仍通过。

### 7.2 `web/pages/booking/confirm.vue` Token 管理统一（S3-03）

**当前问题**：手动从 `sessionStorage.getItem('auth_token')` 获取 token 并注入 header。

**修复**：
- 将手动的 `fetch` + `sessionStorage` 替换为使用 `useApi()` composable
- 如果 web 端的 `useApi` 不带 token（当前实现确实不带），需要增强 `web/app/composables/useApi.ts` 使其支持可选的 auth token 注入（类似 admin 端，从 store 或 sessionStorage 读取）
- 或者创建一个 `useAuthApi()` composable 专门用于需要认证的请求
- 确保 `confirm.vue` 中所有 API 调用都通过 composable 发起

---

## Phase 8: 全量测试验证

所有修复完成后，按以下顺序运行全量测试：

```bash
# 1. 后端
cd backend && go test ./... -coverprofile=coverage.out -covermode=atomic
go tool cover -func=coverage.out | grep -v "100.0%"
# 预期：除 domain（无语句）和已排除的 chk 包外，全部 100%

# 2. Admin 前端
cd frontend/admin && pnpm vitest run --coverage
# 预期：所有测试通过，无伪覆盖测试

# 3. Web 前端
cd frontend/web && pnpm vitest run --coverage
# 预期：所有测试通过

# 4. MiniApp 前端
cd frontend/miniapp && pnpm test
# 预期：所有测试通过
```

### 验收标准

- [ ] 后端覆盖率：repository ≥ 99%，service ≥ 99%，handler ≥ 99%
- [ ] 迁移链：`000001` → `000002` → `000003` → `000004` → `000005` → `000006` 顺序正确，无重复编号
- [ ] 前端零空壳页面：所有 16 个空壳页面全部实联
- [ ] 前端零伪覆盖测试：3 个伪覆盖测试文件全部重写
- [ ] 所有测试通过，无跳过、无注释、无 `istanbul ignore`
- [ ] `AuthService` 使用 `domain.StaffRepository` 接口而非具体类型
- [ ] `web/stores/cruise.ts` 使用 Composition API 风格
- [ ] `web/booking/confirm.vue` 使用 composable 管理 token

---

## 附录：修复项与 Review 报告编号对照表

| 修复编号 | Review 编号 | Phase | 简述 |
|---------|-----------|-------|------|
| 1.1-1.5 | S4-10 | 1 | 迁移编号冲突修复 |
| 2.1.1 | S4-08 | 2 | analytics_repo 测试 |
| 2.1.2 | S4-08 | 2 | notification_repo 测试 |
| 2.1.3 | S4-08 | 2 | payment_repo 测试 |
| 2.1.4 | S4-08 | 2 | refund_repo 测试 |
| 2.2 | S3-04 | 2 | FindBookingByID 测试 |
| 2.3 | S4-08 | 2 | analytics_handler 错误分支 |
| 2.4.1 | S4-08 | 2 | payment_service 错误分支 |
| 2.4.2 | S4-08 | 2 | refund_service 错误分支 |
| 2.5 | S4-09 | 2 | chk 工具排除 |
| 2.6 | S4-08 | 2 | notification TableName |
| 3.1 | S1-03 | 3 | AuthService DI 接口化 |
| 4.1 | S1-01 | 4 | admin login 实联 |
| 4.2 | S1-02 | 4 | admin cruises 实联 |
| 4.3 | S2-03 | 4 | admin routes 实联 |
| 4.4 | S2-04 | 4 | admin voyages 实联 |
| 4.5 | S2-01 | 4 | admin cabins 实联 |
| 4.6 | S2-02 | 4 | admin inventory 实联 |
| 4.7 | S2-05 | 4 | admin pricing 实联 |
| 4.8 | S2-06 | 4 | admin routes/new 实联 |
| 4.9 | S2-07 | 4 | admin voyages/new 实联 |
| 4.10 | S2-09 | 4 | miniapp cabin detail 实联 |
| 5.1 | S3-01 | 5 | miniapp login 实联 |
| 5.2 | S3-02 | 5 | web booking 验证 |
| 5.3 | S4-01 | 5 | admin dashboard 实联 |
| 5.4 | S4-02 | 5 | admin finance 实联 |
| 5.5 | S4-03 | 5 | web pay page 实联 |
| 5.6 | S4-04 | 5 | web PayButton 实联 |
| 5.7 | S4-05 | 5 | miniapp pay 实联 |
| 5.8 | — | 5 | miniapp PayButton 实联 |
| 6.1 | S2-08 | 6 | sprint2-pages 测试重写 |
| 6.2 | S4-06 | 6 | dashboard 测试重写 |
| 6.3 | S4-07 | 6 | pay.page 测试重写 |
| 6.4-6.6 | — | 6 | 新增页面测试 |
| 7.1 | S1-04 | 7 | cruise store Composition API |
| 7.2 | S3-03 | 7 | confirm.vue token 管理 |

---

## 重要提醒

> 1. **每修复一个文件后立即运行相关测试**，不要等全部修复完再测试。
> 2. **不要引入新的空壳页面或伪覆盖测试**。每个新增的测试文件必须包含有意义的业务断言。
> 3. **保持 Git 提交粒度**：要求每个 Phase 一个 commit，commit message 格式：`fix(phase-N): 修复描述`。
> 4. **如果某个修复需要修改多个文件**（如 Phase 3 的 DI 接口化），确保所有相关文件都更新后再运行测试。
> 5. **前端测试中的 mock 数据结构必须与后端 API 实际返回结构一致**。参照 `backend/internal/handler/` 中对应 handler 的 `response.Success(c, gin.H{...})` 调用来确定返回结构。
