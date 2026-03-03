# Sprint 4.2 Audit Remediation Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 按 `TODO.md` 审计清单完成 `P0 -> P1 -> P2` 全量整改，补齐测试与证据，达到可复审状态。

**Architecture:** 采用后端优先的安全/一致性加固策略，先落地交易链路（认证、支付、关单、状态机、导出、通知模板、看板查询）并加事务/幂等，再补前端真实 API 联调与三态兜底。所有改动均按 TDD（RED-GREEN-REFACTOR）推进，并在每批次同步更新 `TODO.md` 和证据文档，避免末尾集中返工。关键路径（下单、支付、关单、导出、权限）以单测 + 接口测试 + 前端冒烟联调组成回归闭环。

**Tech Stack:** Go 1.26, Gin, GORM, PostgreSQL 17, Redis 7.4, Casbin, Nuxt 4, Vue 3, Pinia, Vitest, Playwright, uni-app.

---

## 工作前准备

- Skill refs: `@vue-best-practices` `@vue-development-guides` `@vue-testing-best-practices`
- 证据输出文件：`docs/plans/2026-03-02-sprint42-audit-evidence.md`
- 每个任务完成后同步更新：`TODO.md`

### Task 1: 基线检查与证据骨架

**Files:**
- Create: `docs/plans/2026-03-02-sprint42-audit-evidence.md`
- Modify: `TODO.md`

**Step 1: 写失败前置检查脚本（命令清单）**

```bash
cd backend && go test ./... -count=1
cd ../frontend/admin && pnpm vitest run
cd ../web && pnpm vitest run
cd ../miniapp && pnpm vitest run
```

**Step 2: 运行并记录当前失败点**

Run: 上述命令
Expected: 至少存在与 TODO 对应的失败/缺失项

**Step 3: 建立证据文档骨架**

```markdown
# Sprint 4.2 Audit Evidence

## P0
## P1
## P2
## 回归测试汇总
## 联调记录
```

**Step 4: 保存基线结果**

Run: 将失败摘要写入证据文档
Expected: 可追溯整改前状态

**Step 5: Commit**

```bash
git add docs/plans/2026-03-02-sprint42-audit-evidence.md TODO.md
git commit -m "chore(sprint4.2): add audit evidence scaffold and baseline checkpoints"
```

---

## P0 Blockers

### Task 2: Admin 缺失页接入真实 API + 三态

**Files:**
- Create: `frontend/admin/app/pages/staff/index.vue`
- Create: `frontend/admin/app/pages/settings/shop.vue`
- Create: `frontend/admin/app/pages/notifications/templates.vue`
- Modify: `frontend/admin/app/composables/useApi.ts`
- Modify: `frontend/admin/app/types/*.ts`
- Test: `frontend/admin/tests/pages/staff.index.spec.ts`
- Test: `frontend/admin/tests/pages/settings.shop.spec.ts`
- Test: `frontend/admin/tests/pages/notifications.templates.spec.ts`

**Step 1: 写失败测试（页面需触发真实 API 与三态）**

```ts
it('shows loading, error, empty, and data states', async () => {
  // mock API: pending -> error -> [] -> [data]
})
```

**Step 2: 运行失败测试**

Run: `cd frontend/admin && pnpm vitest run tests/pages/staff.index.spec.ts`
Expected: FAIL（页面不存在或无三态）

**Step 3: 最小实现页面 + API 调用**

```vue
<script setup lang="ts">
const { data, pending, error, refresh } = await useFetch('/api/admin/staff')
</script>
```

**Step 4: 运行通过测试**

Run: `cd frontend/admin && pnpm vitest run tests/pages/staff.index.spec.ts tests/pages/settings.shop.spec.ts tests/pages/notifications.templates.spec.ts`
Expected: PASS

**Step 5: Commit**

```bash
git add frontend/admin/app/pages frontend/admin/app/composables frontend/admin/tests/pages
git commit -m "feat(admin): add staff/shop/template pages with real api and loading-error-empty states"
```

### Task 3: Miniapp 舱位列表去硬编码 + 三态

**Files:**
- Modify: `frontend/miniapp/pages/cabin/list.vue`
- Modify: `frontend/miniapp/components/CabinCard.vue`
- Modify: `frontend/miniapp/services/cabin.ts`
- Test: `frontend/miniapp/tests/pages/cabin.list.spec.ts`
- Test: `frontend/miniapp/tests/components/CabinCard.spec.ts`

**Step 1: 写失败测试（禁止硬编码）**

```ts
expect(screen.queryByText('mock-cabin-hardcode')).toBeNull()
```

**Step 2: 运行失败测试**

Run: `cd frontend/miniapp && pnpm vitest run tests/pages/cabin.list.spec.ts`
Expected: FAIL

**Step 3: 改为真实接口渲染并补三态**

```ts
const { list, loading, error } = useCabinList(params)
```

**Step 4: 运行通过测试**

Run: `cd frontend/miniapp && pnpm vitest run tests/pages/cabin.list.spec.ts tests/components/CabinCard.spec.ts`
Expected: PASS

**Step 5: Commit**

```bash
git add frontend/miniapp/pages/cabin/list.vue frontend/miniapp/components/CabinCard.vue frontend/miniapp/services/cabin.ts frontend/miniapp/tests
git commit -m "feat(miniapp): replace hardcoded cabin data with real api and full state handling"
```

### Task 4: 支付宝登录签名验签与 UID 防伪造

**Files:**
- Modify: `backend/internal/service/user_auth_service.go`
- Create: `backend/internal/service/alipay_verify.go`
- Test: `backend/internal/service/user_auth_service_test.go`
- Test: `backend/internal/service/alipay_verify_test.go`

**Step 1: 写失败测试（签名错误/重放/UID伪造）**

```go
func TestAlipayLoginRejectInvalidSignature(t *testing.T) {}
func TestAlipayLoginRejectMismatchedUID(t *testing.T) {}
```

**Step 2: 运行失败测试**

Run: `cd backend && go test ./internal/service -run "TestAlipayLogin" -count=1`
Expected: FAIL

**Step 3: 最小实现验签与防伪造**

```go
if !verifier.Verify(payload, signature) { return "", ErrInvalidSignature }
if payload.AlipayUID != verifiedUID { return "", ErrUIDSpoofing }
```

**Step 4: 运行通过测试**

Run: `cd backend && go test ./internal/service -run "TestAlipayLogin|TestAlipayVerify" -count=1`
Expected: PASS

**Step 5: Commit**

```bash
git add backend/internal/service/user_auth_service.go backend/internal/service/alipay_verify.go backend/internal/service/*_test.go
git commit -m "fix(auth): verify alipay signature and prevent uid spoofing"
```

### Task 5: 账号绑定唯一性 + 绑定前身份确认

**Files:**
- Modify: `backend/internal/service/user_auth_service.go`
- Modify: `backend/internal/repository/user_repo.go`
- Test: `backend/internal/service/user_auth_service_test.go`

**Step 1: 写失败测试（已绑定冲突、验证码缺失）**

```go
func TestBindAccountRejectAlreadyBound(t *testing.T) {}
func TestBindAccountRequireIdentityChallenge(t *testing.T) {}
```

**Step 2: 运行失败测试**

Run: `cd backend && go test ./internal/service -run "TestBindAccount" -count=1`
Expected: FAIL

**Step 3: 实现唯一绑定校验与确认流程**

```go
if boundUserID != 0 && boundUserID != userID { return ErrIdentifierAlreadyBound }
if !challengeStore.Verify(userID, challengeToken) { return ErrIdentityChallengeRequired }
```

**Step 4: 运行通过测试**

Run: `cd backend && go test ./internal/service -run "TestBindAccount" -count=1`
Expected: PASS

**Step 5: Commit**

```bash
git add backend/internal/service/user_auth_service.go backend/internal/repository/user_repo.go backend/internal/service/user_auth_service_test.go
git commit -m "fix(auth): enforce unique account binding with identity challenge"
```

### Task 6: 超时关单事务化 + 并发幂等

**Files:**
- Modify: `backend/internal/service/order_timeout_service.go`
- Modify: `backend/internal/repository/order_repo.go`
- Modify: `backend/internal/repository/cabin_inventory_repo.go`
- Test: `backend/internal/service/order_timeout_service_test.go`
- Test: `backend/internal/service/order_timeout_service_concurrency_test.go`

**Step 1: 写失败测试（并发重复关单/回滚）**

```go
func TestCloseExpiredOrders_IdempotentUnderConcurrency(t *testing.T) {}
func TestCloseExpiredOrders_RollbackOnInventoryReleaseFailure(t *testing.T) {}
```

**Step 2: 运行失败测试**

Run: `cd backend && go test ./internal/service -run "TestCloseExpiredOrders" -count=1`
Expected: FAIL

**Step 3: 实现事务 + 行锁 + 幂等条件更新**

```go
// tx: lock order row FOR UPDATE; update where status=pending_payment; release inventory; insert log
```

**Step 4: 运行通过测试**

Run: `cd backend && go test ./internal/service -run "TestCloseExpiredOrders" -race -count=1`
Expected: PASS

**Step 5: Commit**

```bash
git add backend/internal/service/order_timeout_service.go backend/internal/repository/order_repo.go backend/internal/repository/cabin_inventory_repo.go backend/internal/service/order_timeout_service*_test.go
git commit -m "fix(order): make expired-order closing transactional idempotent and concurrency-safe"
```

### Task 7: 订单状态机统一入口 + 状态日志同事务

**Files:**
- Modify: `backend/internal/service/order_service.go`
- Modify: `backend/internal/domain/order.go`
- Modify: `backend/internal/repository/order_repo.go`
- Create: `backend/internal/repository/order_status_log_repo.go`
- Test: `backend/internal/service/order_service_status_test.go`

**Step 1: 写失败测试（绕过状态机、日志缺失）**

```go
func TestChangeStatusRejectInvalidTransition(t *testing.T) {}
func TestChangeStatusWritesLogInSameTx(t *testing.T) {}
```

**Step 2: 运行失败测试**

Run: `cd backend && go test ./internal/service -run "TestChangeStatus" -count=1`
Expected: FAIL

**Step 3: 实现统一 `ChangeStatus` 入口**

```go
if !order.CanTransitionTo(to) { return ErrInvalidTransition }
// tx: update order + insert order_status_logs(operator_id, remark)
```

**Step 4: 运行通过测试**

Run: `cd backend && go test ./internal/service -run "TestChangeStatus" -count=1`
Expected: PASS

**Step 5: Commit**

```bash
git add backend/internal/service/order_service.go backend/internal/domain/order.go backend/internal/repository/order*.go backend/internal/service/order_service_status_test.go
git commit -m "fix(order): enforce state-machine-only transitions with transactional status logs"
```

### Task 8: 订单导出真实实现 + 权限 + 上限 + CSV 注入防护

**Files:**
- Modify: `backend/internal/service/order_export_service.go`
- Modify: `backend/internal/handler/order_handler.go`
- Modify: `backend/internal/middleware/rbac.go`
- Test: `backend/internal/service/order_export_service_test.go`
- Test: `backend/internal/handler/order_handler_export_test.go`

**Step 1: 写失败测试（占位内容/无权限/超限/注入）**

```go
func TestExportRejectWithoutPermission(t *testing.T) {}
func TestExportRejectWhenRowsExceedLimit(t *testing.T) {}
func TestExportSanitizeCSVFormulaInjection(t *testing.T) {}
```

**Step 2: 运行失败测试**

Run: `cd backend && go test ./internal/service ./internal/handler -run "TestExport" -count=1`
Expected: FAIL

**Step 3: 最小实现真实导出**

```go
if !authz.Can(ctxUser, "order:export") { return ErrForbidden }
if total > maxExportRows { return ErrExportLimitExceeded }
func sanitizeCSV(v string) string { /* prefix ' for = + - @ */ }
```

**Step 4: 运行通过测试**

Run: `cd backend && go test ./internal/service ./internal/handler -run "TestExport" -count=1`
Expected: PASS

**Step 5: Commit**

```bash
git add backend/internal/service/order_export_service.go backend/internal/handler/order_handler.go backend/internal/middleware/rbac.go backend/internal/service/order_export_service_test.go backend/internal/handler/order_handler_export_test.go
git commit -m "fix(export): implement secure order export with permission limits and csv injection guard"
```

### Task 9: 通知模板 `text/template` + 变量白名单

**Files:**
- Modify: `backend/internal/domain/notification_template.go`
- Modify: `backend/internal/service/notify_service.go`
- Test: `backend/internal/service/notify_template_test.go`

**Step 1: 写失败测试（非法变量/模板函数注入）**

```go
func TestRenderRejectUnknownVariables(t *testing.T) {}
func TestRenderRejectDangerousTemplateExpr(t *testing.T) {}
```

**Step 2: 运行失败测试**

Run: `cd backend && go test ./internal/service -run "TestRender" -count=1`
Expected: FAIL

**Step 3: 实现安全渲染**

```go
tpl, err := template.New("notify").Option("missingkey=error").Parse(safeTemplate)
if !allVarsInWhitelist(parsedVars, whitelist) { return "", ErrInvalidTemplateVariables }
```

**Step 4: 运行通过测试**

Run: `cd backend && go test ./internal/service -run "TestRender" -count=1`
Expected: PASS

**Step 5: Commit**

```bash
git add backend/internal/domain/notification_template.go backend/internal/service/notify_service.go backend/internal/service/notify_template_test.go
git commit -m "fix(notify): use text template with variable whitelist and ssti-safe validation"
```

### Task 10: Analytics 真实查询 + 索引 + 性能说明

**Files:**
- Modify: `backend/internal/service/analytics_service.go`
- Modify: `backend/internal/repository/analytics_repo.go`
- Modify: `backend/internal/handler/analytics_handler.go`
- Create: `backend/migrations/000012_sprint42_analytics_indexes.up.sql`
- Create: `backend/migrations/000012_sprint42_analytics_indexes.down.sql`
- Test: `backend/internal/service/analytics_service_test.go`
- Test: `backend/internal/repository/analytics_repo_test.go`
- Modify: `docs/plans/2026-03-02-sprint42-audit-evidence.md`

**Step 1: 写失败测试（接口返回空、天数不正确）**

```go
func TestTrendReturns7And30Days(t *testing.T) {}
func TestCabinHotnessRankingNotEmpty(t *testing.T) {}
func TestInventoryOverviewAndPageViewStats(t *testing.T) {}
```

**Step 2: 运行失败测试**

Run: `cd backend && go test ./internal/service ./internal/repository -run "TestTrend|TestCabinHotness|TestInventoryOverview|TestPageView" -count=1`
Expected: FAIL

**Step 3: 实现真实聚合查询并补索引**

```sql
CREATE INDEX IF NOT EXISTS idx_orders_paid_at_status ON orders(status, paid_at);
CREATE INDEX IF NOT EXISTS idx_order_items_cabin_sku_id ON order_items(cabin_sku_id);
CREATE INDEX IF NOT EXISTS idx_page_views_created_at ON page_views(created_at);
```

**Step 4: 运行通过测试并记录耗时**

Run: `cd backend && go test ./internal/service ./internal/repository -run "TestTrend|TestCabinHotness|TestInventoryOverview|TestPageView" -count=1 -v`
Expected: PASS

**Step 5: Commit**

```bash
git add backend/internal/service/analytics_service.go backend/internal/repository/analytics_repo.go backend/internal/handler/analytics_handler.go backend/migrations/000012_sprint42_analytics_indexes.* backend/internal/service/analytics_service_test.go backend/internal/repository/analytics_repo_test.go docs/plans/2026-03-02-sprint42-audit-evidence.md
git commit -m "feat(analytics): implement real trend ranking inventory and pageview queries with indexes"
```

---

## P1 High Priority

### Task 11: 批量操作安全加固（限流 + 事务 + 权限审计）

**Files:**
- Modify: `backend/internal/handler/cabin_handler.go`
- Modify: `backend/internal/service/cabin_service.go`
- Modify: `backend/internal/middleware/operation_log.go`
- Test: `backend/internal/handler/cabin_handler_batch_test.go`
- Test: `backend/internal/service/cabin_service_batch_test.go`

**Step 1: 写失败测试（超上限/权限不足/回滚）**

```go
func TestBatchUpdateRejectTooManyIDs(t *testing.T) {}
func TestBatchUpdateRollbackOnPartialFailure(t *testing.T) {}
```

**Step 2: Run**

Run: `cd backend && go test ./internal/handler ./internal/service -run "TestBatchUpdate" -count=1`
Expected: FAIL

**Step 3: 实现**

```go
if len(req.IDs) > maxBatchIDs { return ErrBatchLimitExceeded }
// tx for batch write + op log
```

**Step 4: Re-run**

Run: `cd backend && go test ./internal/handler ./internal/service -run "TestBatchUpdate" -count=1`
Expected: PASS

**Step 5: Commit**

```bash
git add backend/internal/handler/cabin_handler.go backend/internal/service/cabin_service.go backend/internal/middleware/operation_log.go backend/internal/handler/cabin_handler_batch_test.go backend/internal/service/cabin_service_batch_test.go
git commit -m "fix(cabin): harden batch operations with limits tx rollback and audit checks"
```

### Task 12: `AvailableWithAlert` 与阈值边界测试

**Files:**
- Modify: `backend/internal/service/inventory_service.go`
- Test: `backend/internal/service/inventory_service_test.go`

**Step 1: 写失败边界测试**

```go
func TestAvailableWithAlert_EqualThreshold(t *testing.T) {}
func TestAvailableWithAlert_LessThanThreshold(t *testing.T) {}
func TestAvailableWithAlert_ThresholdZero(t *testing.T) {}
```

**Step 2: Run**

Run: `cd backend && go test ./internal/service -run "TestAvailableWithAlert" -count=1`
Expected: FAIL

**Step 3: 实现最小逻辑**

```go
below := inv.AlertThreshold > 0 && avail <= inv.AlertThreshold
```

**Step 4: Re-run**

Run: `cd backend && go test ./internal/service -run "TestAvailableWithAlert" -count=1`
Expected: PASS

**Step 5: Commit**

```bash
git add backend/internal/service/inventory_service.go backend/internal/service/inventory_service_test.go
git commit -m "test(inventory): add alert threshold boundary coverage"
```

### Task 13: 批量日期范围定价 API 暴露与接口测试

**Files:**
- Modify: `backend/internal/handler/pricing_handler.go`
- Modify: `backend/internal/router/router.go`
- Test: `backend/internal/handler/pricing_handler_test.go`

**Step 1: 写失败接口测试**

```go
func TestBatchSetPriceAPI(t *testing.T) {}
func TestBatchSetPriceAPIRejectInvalidDateRange(t *testing.T) {}
```

**Step 2: Run**

Run: `cd backend && go test ./internal/handler -run "TestBatchSetPriceAPI" -count=1`
Expected: FAIL

**Step 3: 实现 handler + router**

```go
admin.POST("/pricing/batch-set", pricingHandler.BatchSetPrice)
```

**Step 4: Re-run**

Run: `cd backend && go test ./internal/handler -run "TestBatchSetPriceAPI" -count=1`
Expected: PASS

**Step 5: Commit**

```bash
git add backend/internal/handler/pricing_handler.go backend/internal/router/router.go backend/internal/handler/pricing_handler_test.go
git commit -m "feat(pricing): expose batch date-range pricing api with validation"
```

### Task 14: 对账日报幂等与并发测试

**Files:**
- Modify: `backend/internal/service/reconciliation_service.go`
- Modify: `backend/internal/repository/reconciliation_repo.go`
- Test: `backend/internal/service/reconciliation_service_test.go`

**Step 1: 写失败测试（同日重复/并发）**

```go
func TestGenerateDailyReport_IdempotentByDate(t *testing.T) {}
func TestGenerateDailyReport_ConcurrentCalls(t *testing.T) {}
```

**Step 2: Run**

Run: `cd backend && go test ./internal/service -run "TestGenerateDailyReport" -race -count=1`
Expected: FAIL

**Step 3: 实现唯一键冲突处理与锁**

```go
// ON CONFLICT(report_date) DO NOTHING/UPDATE per policy
```

**Step 4: Re-run**

Run: `cd backend && go test ./internal/service -run "TestGenerateDailyReport" -race -count=1`
Expected: PASS

**Step 5: Commit**

```bash
git add backend/internal/service/reconciliation_service.go backend/internal/repository/reconciliation_repo.go backend/internal/service/reconciliation_service_test.go
git commit -m "fix(reconciliation): enforce date-level idempotency and concurrency safety"
```

### Task 15: 员工角色变更同步 Casbin + 审计日志

**Files:**
- Modify: `backend/internal/service/staff_service.go`
- Modify: `backend/internal/service/rbac_service.go`
- Modify: `backend/internal/domain/operation_log.go`
- Test: `backend/internal/service/staff_service_test.go`

**Step 1: 写失败测试（角色改了但 Casbin 未变）**

```go
func TestAssignRoleSyncsCasbinGrouping(t *testing.T) {}
func TestAssignRoleWritesOperationLog(t *testing.T) {}
```

**Step 2: Run**

Run: `cd backend && go test ./internal/service -run "TestAssignRole" -count=1`
Expected: FAIL

**Step 3: 实现**

```go
err := enforcer.AddGroupingPolicy(user, role)
// write operation log in same tx boundary
```

**Step 4: Re-run**

Run: `cd backend && go test ./internal/service -run "TestAssignRole" -count=1`
Expected: PASS

**Step 5: Commit**

```bash
git add backend/internal/service/staff_service.go backend/internal/service/rbac_service.go backend/internal/domain/operation_log.go backend/internal/service/staff_service_test.go
git commit -m "fix(staff): sync role assignment to casbin and persist audit logs"
```

### Task 16: ShopInfo 单行约束（DB + Service）

**Files:**
- Create: `backend/migrations/000013_shop_info_singleton.up.sql`
- Create: `backend/migrations/000013_shop_info_singleton.down.sql`
- Modify: `backend/internal/service/shop_info_service.go`
- Test: `backend/internal/service/shop_info_service_test.go`

**Step 1: 写失败测试（创建多条）**

```go
func TestShopInfoRejectMultipleRows(t *testing.T) {}
```

**Step 2: Run**

Run: `cd backend && go test ./internal/service -run "TestShopInfo" -count=1`
Expected: FAIL

**Step 3: 实现固定主键/唯一约束 + 仅更新模式**

```sql
CREATE UNIQUE INDEX IF NOT EXISTS idx_shop_info_singleton ON shop_info((1));
```

**Step 4: Re-run**

Run: `cd backend && go test ./internal/service -run "TestShopInfo" -count=1`
Expected: PASS

**Step 5: Commit**

```bash
git add backend/migrations/000013_shop_info_singleton.* backend/internal/service/shop_info_service.go backend/internal/service/shop_info_service_test.go
git commit -m "fix(shop): enforce singleton row at db and service layers"
```

### Task 17: 库存预警通知去重窗口

**Files:**
- Modify: `backend/internal/service/inventory_alert_service.go`
- Modify: `backend/internal/repository/inventory_alert_repo.go`
- Test: `backend/internal/service/inventory_alert_service_test.go`

**Step 1: 写失败测试（重复扫描轰炸）**

```go
func TestInventoryAlertDeduplicateWithinWindow(t *testing.T) {}
```

**Step 2: Run**

Run: `cd backend && go test ./internal/service -run "TestInventoryAlert" -count=1`
Expected: FAIL

**Step 3: 实现去重键 + 时间窗口**

```go
dedupeKey := fmt.Sprintf("inventory-alert:%d:%s", skuID, dateBucket)
```

**Step 4: Re-run**

Run: `cd backend && go test ./internal/service -run "TestInventoryAlert" -count=1`
Expected: PASS

**Step 5: Commit**

```bash
git add backend/internal/service/inventory_alert_service.go backend/internal/repository/inventory_alert_repo.go backend/internal/service/inventory_alert_service_test.go
git commit -m "fix(alert): dedupe inventory alerts with time-window keys"
```

---

## P2 Completeness / Test Quality

### Task 18: 退款边界测试补齐（30/7/0 天）

**Files:**
- Modify: `backend/internal/service/refund_service_tiered_test.go`

**Step 1: 写失败边界测试**

```go
func TestRefundBoundaryAt30Days(t *testing.T) {}
func TestRefundBoundaryAt7Days(t *testing.T) {}
func TestRefundBoundaryAt0Days(t *testing.T) {}
```

**Step 2: Run**

Run: `cd backend && go test ./internal/service -run "TestRefundBoundary" -count=1`
Expected: FAIL

**Step 3: 修正区间判断**

```go
// inclusive/exclusive boundaries documented and asserted
```

**Step 4: Re-run**

Run: `cd backend && go test ./internal/service -run "TestRefundBoundary" -count=1`
Expected: PASS

**Step 5: Commit**

```bash
git add backend/internal/service/refund_service_tiered_test.go backend/internal/service/refund_service.go
git commit -m "test(refund): cover exact 30 7 0 day boundaries"
```

### Task 19: 清理空实现与伪覆盖

**Files:**
- Modify: `backend/internal/handler/booking_handler.go`
- Test: `backend/internal/handler/booking_handler_test.go`
- Modify: `docs/plans/2026-03-02-sprint42-audit-evidence.md`

**Step 1: 写失败测试（UpdateStatus 不得为空）**

```go
func TestBookingUpdateStatusExecutesRealLogic(t *testing.T) {}
```

**Step 2: Run**

Run: `cd backend && go test ./internal/handler -run "TestBookingUpdateStatus" -count=1`
Expected: FAIL

**Step 3: 实现真实逻辑并补异常路径**

```go
// validate input -> service.ChangeStatus -> response
```

**Step 4: Re-run**

Run: `cd backend && go test ./internal/handler -run "TestBookingUpdateStatus" -count=1`
Expected: PASS

**Step 5: Commit**

```bash
git add backend/internal/handler/booking_handler.go backend/internal/handler/booking_handler_test.go docs/plans/2026-03-02-sprint42-audit-evidence.md
git commit -m "fix(booking): replace placeholder update-status with real implementation"
```

### Task 20: Sprint 4.2 27 Tasks 对照复核证据

**Files:**
- Modify: `docs/plans/2026-03-02-sprint42-audit-evidence.md`
- Modify: `TODO.md`

**Step 1: 建立 task->code/test mapping 表**

```markdown
| Task | Code Path | Test Case | API/UI Evidence |
```

**Step 2: 填充映射并标注缺口**

Run: 手工梳理 + grep 校验
Expected: 每个 Task 有证据落点

**Step 3: 补齐缺口并更新文档**

```markdown
- 关键截图/录屏路径（如 docs/evidence/screenshots/...）
```

**Step 4: 复查一致性**

Run: 对照 `docs/plans/2026-02-22-sprint04.2.md`
Expected: 无遗漏

**Step 5: Commit**

```bash
git add docs/plans/2026-03-02-sprint42-audit-evidence.md TODO.md
git commit -m "docs(audit): add full sprint4.2 task-to-evidence mapping"
```

### Task 21: Admin 目录迁移与路由冒烟

**Files:**
- Move: `frontend/admin/pages/dashboard/**` -> `frontend/admin/app/pages/dashboard/**`
- Move: `frontend/admin/pages/finance/**` -> `frontend/admin/app/pages/finance/**`
- Modify: `frontend/admin/nuxt.config.ts`
- Test: `frontend/admin/tests/smoke/routes.spec.ts`

**Step 1: 写失败冒烟测试（/dashboard /finance）**

```ts
it('routes /dashboard and /finance are reachable', async () => {})
```

**Step 2: Run**

Run: `cd frontend/admin && pnpm vitest run tests/smoke/routes.spec.ts`
Expected: FAIL

**Step 3: 迁移目录并修正路由配置**

```ts
// ensure no duplicate route names
```

**Step 4: Re-run**

Run: `cd frontend/admin && pnpm vitest run tests/smoke/routes.spec.ts`
Expected: PASS

**Step 5: Commit**

```bash
git add frontend/admin/app/pages frontend/admin/pages frontend/admin/nuxt.config.ts frontend/admin/tests/smoke/routes.spec.ts
git commit -m "refactor(admin): migrate dashboard and finance pages into app/pages with route smoke coverage"
```

---

## 收尾与门禁验证

### Task 22: 全量回归 + TODO 收口

**Files:**
- Modify: `TODO.md`
- Modify: `docs/plans/2026-03-02-sprint42-audit-evidence.md`

**Step 1: 后端全量回归**

Run: `cd backend && go test ./... -coverprofile=coverage.out -covermode=atomic`
Expected: PASS

**Step 2: 前端全量回归**

Run: 
- `cd frontend/admin && pnpm vitest run`
- `cd frontend/web && pnpm vitest run`
- `cd frontend/miniapp && pnpm vitest run`
Expected: PASS

**Step 3: 关键链路联调回归**

Run: 依次验证 下单 -> 支付回调 -> 超时关单 -> 导出 -> 权限拒绝
Expected: 全通过并有证据

**Step 4: 更新 TODO 勾选 + 证据链接**

```markdown
- [x] 项目名（证据: docs/plans/2026-03-02-sprint42-audit-evidence.md#...）
```

**Step 5: Commit**

```bash
git add TODO.md docs/plans/2026-03-02-sprint42-audit-evidence.md backend/coverage.out
git commit -m "chore(audit): close sprint4.2 todo with full regression evidence"
```
