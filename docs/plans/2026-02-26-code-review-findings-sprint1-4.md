# Cumulative Code Review Findings: Sprint 1 → Sprint 4

> **审查日期**: 2026-02-26
> **审查范围**: Sprint 1 至 Sprint 4 全部代码
> **审查方法**: 按 `2026-02-24-cumulative-code-review-sprint1-4.md` 指示执行累进审查

---

## 测试执行结果

### 后端 (Go)

```
config        100.0%
domain          0.0% (无可执行语句)
handler        99.2%
middleware    100.0%
database      100.0%
errcode       [no statements]
logger        100.0%
response      100.0%
search        100.0%
repository     80.0%  ← 未达标
router        100.0%
service        98.6%  ← 未达标
service/chk     0.0% (开发辅助工具，非交付物)
```

### 前端

| 应用 | 测试文件数 | 测试用例数 | 通过率 |
|------|-----------|-----------|--------|
| Admin | 13 | 24 | 100% |
| Web | 10 | 23 | 100% |
| MiniApp | 8 | 13 | 100% |

---

## Sprint 1 审查结论

**状态**: ❌ 不通过

### 通过项

- [x] Docker & Infra: `docker-compose.yml` 使用环境变量注入密码，符合最小权限原则
- [x] Auth & RBAC: JWT 中间件强制验证 HMAC 签名算法 (`alg == "HS256"`)；Casbin 中间件从 JWT Claims 提取角色
- [x] Domain 模型: Staff/Role/CruiseCompany/Cruise/CabinType/FacilityCategory/Facility/Image 定义完整
- [x] Repository 实现: 均使用 `WithContext(ctx)`
- [x] Config/Logger/Database: 配置加载、日志、数据库连接均实现且有测试
- [x] 迁移文件: `000001_init_schema.up.sql` 包含所有 Sprint 1 表，索引和外键完备
- [x] Swagger: `docs/` 目录包含 swagger.json/yaml + docs.go
- [x] 前端共享类型: `frontend/shared/types/` 定义了 api.ts、domain.ts、constants/index.ts
- [x] Admin/Web/MiniApp 布局和基础设施已就位
- [x] CI: `.github/workflows/ci.yml` 存在

### 问题清单

| # | 严重程度 | 文件 | 问题描述 | 修改建议 |
|---|---------|------|---------|---------|
| S1-01 | 阻塞 | `frontend/admin/app/pages/login.vue` | 登录提交使用 `setTimeout` 模拟，未调用任何后端 API (`// TODO: 实际环境替换为 API 调用`) | 调用 `useApi().request('/auth/login', { method: 'POST', body })` 实现实际认证 |
| S1-02 | 阻塞 | `frontend/admin/app/pages/cruises/index.vue` | 邮轮列表使用硬编码数据 `["示例邮轮"]`，未调用后端 API | 改为 `onMounted(() => useApi().request('/cruises'))` 获取真实数据 |
| S1-03 | 警告 | `backend/internal/service/auth_service.go` | `NewAuthService` 接收具体类型 `*repository.StaffRepository` 而非接口，违反 DI/DDD 原则 | 改为接收 `domain.StaffRepository` 接口 |
| S1-04 | 警告 | `frontend/web/app/stores/cruise.ts` | 使用 Pinia Options API 风格，与 admin 和 miniapp 的 Composition API 风格不一致 | 统一为 `defineStore('cruise', () => { ... })` Composition API 风格 |

### 被忽略的后续 Sprint 变更

- `domain/repository.go` 中包含 Sprint 2-4 新增的 Repository 接口 — 已按规则忽略
- `router.go` 中包含 Sprint 2-4 新增的路由注册 — 已按规则忽略
- `cmd/server/main.go` 中包含 Sprint 2-4 的完整 DI 链 — 已按规则忽略
- Sprint 2-4 的所有新增文件 — 已按规则忽略

---

## Sprint 2 审查结论

**状态**: ❌ 不通过

### 通过项

- [x] 库存一致性: `cabin_repo.go` 的 `AdjustInventoryAtomic()` 使用原子 UPDATE 防超卖
- [x] Domain 模型: `CabinPrice.PriceCents` 使用 `int64` 整数分
- [x] Repository 实现: route_repo/voyage_repo/cabin_repo 均使用 `WithContext(ctx)`
- [x] Pricing Service: `sameDay()` 使用 `.Date()` 方法比较日期
- [x] Inventory Service: 正确委托给 Repository 的原子操作
- [x] Admin CRUD Handler: Route/Voyage/Cabin handler 带参数验证、分页、MeiliSearch 索引
- [x] MeiliSearch: `meili.go` 适配器实现 `IndexDocuments`，有搜索重试队列

### 问题清单

| # | 严重程度 | 文件 | 问题描述 | 修改建议 |
|---|---------|------|---------|---------|
| S2-01 | 阻塞 | `frontend/admin/app/pages/cabins/index.vue` | **空壳页面**：仅包含 `<h1>Cabins</h1>` 和 TODO 注释，无任何 API 调用或数据展示 | 实现舱位列表，调用 `GET /api/v1/admin/cabins` |
| S2-02 | 阻塞 | `frontend/admin/app/pages/cabins/inventory.vue` | **空壳页面**：仅包含 `<h1>Inventory</h1>` 和 TODO 注释 | 实现库存管理页面，调用库存 API |
| S2-03 | 阻塞 | `frontend/admin/app/pages/routes/index.vue` | 使用硬编码数据 `[{ id: 1, code: 'R1', name: 'Route 1' }]`，未调用后端 API | 调用 `useApi().request('/routes')` 获取真实数据 |
| S2-04 | 阻塞 | `frontend/admin/app/pages/voyages/index.vue` | 使用硬编码数据，未调用后端 API | 调用 `useApi().request('/voyages')` 获取真实数据 |
| S2-05 | 阻塞 | `frontend/admin/app/pages/cabins/pricing.vue` | 使用硬编码价格数据，未调用后端定价 API | 调用 Pricing API 获取真实数据 |
| S2-06 | 阻塞 | `frontend/admin/app/pages/routes/new.vue` | 表单无提交逻辑，无 API 调用 | 实现表单提交调用 `POST /api/v1/admin/routes` |
| S2-07 | 阻塞 | `frontend/admin/app/pages/voyages/new.vue` | 表单无提交逻辑，无 API 调用 | 实现表单提交调用 `POST /api/v1/admin/voyages` |
| S2-08 | 阻塞 | `frontend/admin/tests/unit/pages/sprint2-pages.spec.ts` | **伪覆盖测试**：仅验证 DOM 渲染（`toContain('Cabins')`），无 API 调用断言、无数据流验证 | 为每个页面编写 API 调用、loading/error 状态、表单提交的完整测试 |
| S2-09 | 警告 | `frontend/miniapp/pages/cabin/detail.vue` | **空壳页面**：仅包含标题和 TODO 注释（小程序端 Part E: 舱位详情） | 实现舱房详情和价格日历 API 对接 |
| S2-10 | 警告 | 迁移编号 | Sprint 2 迁移文件编号为 `000003`（计划为 `000002`），因存在额外的 `000002_add_staff_roles` | 不影响功能但偏离计划文档 |

### 被忽略的后续 Sprint 变更

- Sprint 3/4 新增的所有文件 — 已按规则忽略
- `router.go` 中 Sprint 3/4 新增的路由 — 已按规则忽略

---

## Sprint 3 审查结论

**状态**: ❌ 不通过

### 通过项

- [x] 下单幂等与并发: Booking Service 在单一事务内完成 Hold → Price → Create，唯一索引防重复锁定
- [x] Auth Service: `UserAuthService` 限制尝试次数（MaxAttempts=5）、有效期（1min）、锁定时间（30min），Mutex 保护并发
- [x] Booking Transaction: `BookingRepo.InTx()` 确保单一事务，失败自动回滚
- [x] Domain 模型: User/Passenger/Booking/BookingPassenger 定义完整
- [x] Frontend Validation: Web booking/confirm.vue 有 `canSubmit` 防连点
- [x] Admin Bookings: bookings/index.vue 和 bookings/[id].vue 使用 `useApi()` 调用 API，正确处理 loading/error
- [x] Web LoginForm: 完整 SMS 登录流程
- [x] MiniApp Booking: booking/create.vue 调用 `/bookings` API

### 问题清单

| # | 严重程度 | 文件 | 问题描述 | 修改建议 |
|---|---------|------|---------|---------|
| S3-01 | 阻塞 | `frontend/miniapp/pages/login/login.vue` | **空壳页面**：仅显示 `<text>Login</text>` 和一个 PrimaryButton，无短信验证码输入、无 API 调用、无登录逻辑 | 实现完整的手机号+验证码登录流程，参考 web/LoginForm.vue |
| S3-02 | 阻塞 | `frontend/web/pages/booking/index.vue` | 预订第一步仅做本地表单跳转（`router.push`），未调用任何后端 API 验证航次/舱位有效性 | 跳转前调用 API 验证航次和舱位的有效性 |
| S3-03 | 警告 | `frontend/web/pages/booking/confirm.vue` | Token 从 `sessionStorage` 手动获取并注入 header，未使用 `useApi()` composable；与 admin 端 token 管理方式不一致 | 统一使用 `useApi()` composable 自动注入 Authorization header |
| S3-04 | 警告 | `backend/internal/repository/booking_repo.go` L37 | `FindBookingByID` 函数 0% 测试覆盖率 | 补充 `FindBookingByID` 的单元测试用例 |

### 被忽略的后续 Sprint 变更

- Sprint 4 新增的所有文件 — 已按规则忽略
- `router.go` 中 Sprint 4 新增的路由 — 已按规则忽略

---

## Sprint 4 审查结论

**状态**: ❌ 不通过

### 通过项

- [x] 支付回调安全: `PaymentCallbackService` 使用 `hmac.Equal()` 时序安全签名验证 + 幂等性检查
- [x] 退款逻辑: `RefundService` 验证 `amount > 0`、原支付状态 "paid"、累计退款 ≤ 原支付金额
- [x] 通知发件箱: `NotifyService` 采用 Outbox 模式，Enqueue 写入表，OutboxDispatcher 异步处理
- [x] Domain 模型: Payment/Refund 完整，`AmountCents` 使用 `int64`
- [x] Web Orders: `orders/[id].vue` 调用 API，完整处理 loading/error/empty 三态
- [x] Backend Tests: payment/refund/notify/analytics service 测试覆盖边界和正常场景

### 问题清单

| # | 严重程度 | 文件 | 问题描述 | 修改建议 |
|---|---------|------|---------|---------|
| S4-01 | 阻塞 | `frontend/admin/pages/dashboard/index.vue` | **空壳页面+硬编码数据**：`const summary = { sales: 1000, orders: 12 }`，未调用 analytics API，无 loading/error 处理 | 调用 `GET /api/v1/admin/analytics/dashboard` 获取真实数据 |
| S4-02 | 阻塞 | `frontend/admin/pages/finance/index.vue` | **空壳页面+硬编码数据**：表格中硬编码示例行，未调用任何财务 API | 调用财务 API 获取真实流水数据 |
| S4-03 | 阻塞 | `frontend/web/pages/pay/[id].vue` | **空壳页面**：仅导入 PayButton 组件，未获取订单信息、无支付金额展示、无支付 API 调用 | 调用订单 API 获取金额，点击时调用支付 API |
| S4-04 | 阻塞 | `frontend/web/components/PayButton.vue` | **空壳组件**：仅有 `<button>Pay Now</button>` 骨架，无任何支付逻辑或 API 调用 | 实现支付创建 API 调用、loading 状态、错误处理 |
| S4-05 | 阻塞 | `frontend/miniapp/pages/pay/pay.vue` | **空壳页面**：仅导入 PayButton，无支付逻辑 | 实现小程序支付流程 (wx.requestPayment) |
| S4-06 | 阻塞 | `frontend/admin/tests/unit/dashboard.page.spec.ts` | **伪覆盖测试**：仅 1 个测试 `toContain('Dashboard')`，无 API 调用断言、无数据渲染验证 | 编写完整测试：mock analytics API → 验证数据渲染 → 验证 loading/error |
| S4-07 | 阻塞 | `frontend/web/tests/unit/pay.page.spec.ts` | **伪覆盖测试**：仅 1 个测试 `toContain('Pay Now')`，无支付 API 调用验证 | 编写支付流程测试：点击 → API 调用 → 成功/失败场景 |
| S4-08 | 阻塞 | 后端覆盖率 | repository 层 80%：`analytics_repo.go` / `notification_repo.go` / `payment_repo.go` / `refund_repo.go` 全部 0%；handler 99.2%：`analytics_handler.go` 两个错误分支未覆盖；service 98.6%：`payment_service.go` 3 个分支 + `refund_service.go` 1 个分支未覆盖 | 为所有 0% 覆盖的 repo 编写集成测试；补充 handler/service 缺失分支 |
| S4-09 | 阻塞 | `backend/internal/service/chk/check.go` | 开发辅助工具（AST 扫描），0% 覆盖率，不属于 Sprint 4 交付物但纳入 coverage 统计 | 从 `internal/service/` 迁出或通过 build tag 排除出 coverage |
| S4-10 | 阻塞 | 迁移编号 | **迁移链冲突**：`000004_user_booking` (Sprint 3) 和 `000004_payment_notify` (Sprint 4) 共享 `000004` 编号前缀，golang-migrate 执行顺序不确定 | 将 Sprint 4 迁移重命名为 `000005_payment_notify`，Sprint 3 的 `000005` 重命名为 `000006` |

### 前端实联审查结果

| 页面 | API Endpoint | 是否实联 | loading/error/empty 处理 | 备注 |
|------|-------------|---------|------------------------|------|
| dashboard/index.vue | GET /api/v1/admin/analytics/dashboard | ❌ 否 | ❌ 否 | 硬编码 `{sales:1000, orders:12}` |
| finance/index.vue | GET /api/v1/admin/analytics/* | ❌ 否 | ❌ 否 | 硬编码示例表格行 |
| StatCard.vue | — | ✅ 组件级 (props) | — | 纯展示组件，依赖父页面传入数据 |
| orders/[id].vue | GET /api/v1/bookings/:id | ✅ 是 | ✅ 是 | 完整实现 |
| pay/[id].vue | POST /api/v1/payments | ❌ 否 | ❌ 否 | 仅导入 PayButton 壳组件 |
| PayButton.vue | — | ❌ 否 | ❌ 否 | 仅有按钮骨架 |
| miniapp pay/pay.vue | — | ❌ 否 | ❌ 否 | 仅导入 PayButton |

### 测试质量审查结果

| 测试文件 | 业务断言数 | 弱断言数 | 伪覆盖标记数 | 判定 |
|---------|-----------|---------|------------|------|
| dashboard.page.spec.ts | 0 | 1 | 0 | ❌ 伪覆盖 |
| pay.page.spec.ts | 0 | 1 | 0 | ❌ 伪覆盖 |
| bookings.list.spec.ts | 3 | 0 | 0 | ✅ 通过 |
| bookings.detail.spec.ts | 2 | 0 | 0 | ✅ 通过 |
| booking/confirm.spec.ts | 7 | 0 | 0 | ✅ 通过 |
| booking/create.spec.ts (miniapp) | 3 | 0 | 0 | ✅ 通过 |
| sprint2-pages.spec.ts | 0 | 3 | 0 | ❌ 伪覆盖 |
| payment_service_test.go | 多个 | 0 | 0 | ✅ 通过 |
| refund_service_test.go | 多个 | 0 | 0 | ✅ 通过 |
| analytics_handler_test.go | 多个 | 0 | 缺2个错误分支 | ⚠️ 不完整 |

---

## 全局审查项（跨 Sprint）

### 跨 Sprint 一致性检查

- [x] Domain 模型命名一致性：外键关联正确
- [x] Repository 接口一致性：统一遵循 `New*Repo(db *gorm.DB)` 构造模式
- [x] Router 注册完整性：Sprint 1-4 所有路由已注册
- [ ] ❌ 迁移链完整性：`000004` 重复编号（`000004_user_booking` + `000004_payment_notify`）
- [x] 共享类型同步：`frontend/shared/types/domain.ts` 覆盖全部后端结构
- [x] CI 管道覆盖

### 全局安全检查

- [x] 环境变量泄露：`config.yaml` 使用占位符，未提交真实密钥
- [x] SQL 注入：所有 Repository SQL 均使用 GORM 参数化查询
- [x] XSS 防护：Vue 模板使用 `{{ }}` 自动转义，无 `v-html` 使用
- [x] CORS 配置：`router.go` 使用特定域名，未用 `*` 通配
- [x] 支付安全：金额后端计算，回调 `hmac.Equal()` 验证签名

### 全局空壳页面汇总

| # | 页面 | Sprint | 问题类型 |
|---|------|--------|---------|
| 1 | admin/login.vue | S1 | setTimeout 模拟 |
| 2 | admin/cruises/index.vue | S1 | 硬编码数据 |
| 3 | admin/cabins/index.vue | S2 | 空壳 |
| 4 | admin/cabins/inventory.vue | S2 | 空壳 |
| 5 | admin/routes/index.vue | S2 | 硬编码数据 |
| 6 | admin/voyages/index.vue | S2 | 硬编码数据 |
| 7 | admin/cabins/pricing.vue | S2 | 硬编码数据 |
| 8 | admin/routes/new.vue | S2 | 表单无提交 |
| 9 | admin/voyages/new.vue | S2 | 表单无提交 |
| 10 | miniapp/login/login.vue | S3 | 空壳 |
| 11 | miniapp/cabin/detail.vue | S2 | 空壳 |
| 12 | admin/dashboard/index.vue | S4 | 硬编码数据 |
| 13 | admin/finance/index.vue | S4 | 硬编码数据 |
| 14 | web/pay/[id].vue | S4 | 空壳 |
| 15 | web/PayButton.vue | S4 | 空壳 |
| 16 | miniapp/pay/pay.vue | S4 | 空壳 |

### 后端未覆盖函数/文件明细

| 文件 | 覆盖率 | 未覆盖内容 |
|------|--------|-----------|
| `repository/analytics_repo.go` | 0% | TodaySales / WeeklyTrend / TodayOrderCount 全部 |
| `repository/notification_repo.go` | 0% | Enqueue / FindPending / MarkSent / MarkFailed / MarkProcessing 全部 |
| `repository/payment_repo.go` | 0% | Create / FindByOrderID / FindByID / UpdateStatus 全部 |
| `repository/refund_repo.go` | 0% | Create / FindByPaymentID / SumByPaymentID 全部 |
| `repository/booking_repo.go` L37 | 部分 | FindBookingByID 未覆盖 |
| `handler/analytics_handler.go` | 99.2% | WeeklyTrend 和 TodayOrderCount 错误分支 |
| `service/payment_service.go` | ~97% | 3 个错误分支未覆盖 |
| `service/refund_service.go` | ~97% | amountCents ≤ 0 错误分支未覆盖 |
| `domain/notification.go` | 0% | TableName() 方法 |
| `service/chk/check.go` | 0% | 开发辅助工具，非交付物 |

---

## 汇总报告

| Sprint | 状态 | 阻塞问题数 | 警告数 | 空壳页面数 | 伪覆盖测试数 |
|--------|------|-----------|--------|-----------|-------------|
| 1      | 不通过 | 2 | 2 | 2 | 0 |
| 2      | 不通过 | 8 | 2 | 7 | 1 |
| 3      | 不通过 | 2 | 2 | 1 | 0 |
| 4      | 不通过 | 10 | 0 | 5 | 2 |

### 总体结论

**❌ 需修改后重新 Review**

后端架构设计、安全机制（HMAC 签名验证、原子库存、事务预订、SMS 防爆破）、DDD 分层等核心部分质量良好。主要问题集中在：

1. **前端空壳页面泛滥**：16 个页面/组件缺少 API 集成
2. **后端测试覆盖率不达标 + 伪覆盖**：repository 层 80%（4 个 repo 文件 0%），前端 3 个伪覆盖测试文件

---

## 修改方案

### 优先级 1 - 迁移链修复（阻塞部署）

1. 将 `000004_payment_notify` 重命名为 `000005_payment_notify`
2. 将 `000005_cabin_hold_unique` 重命名为 `000006_cabin_hold_unique`
3. 更新迁移测试文件的编号引用

### 优先级 2 - 后端覆盖率提升至 100%

4. `analytics_repo.go` — 编写 SQLite 集成测试，覆盖 TodaySales / WeeklyTrend / TodayOrderCount
5. `notification_repo.go` — 编写 Enqueue / FindPending / MarkSent / MarkFailed 测试
6. `payment_repo.go` — 编写 Create / FindByOrderID / FindByID / UpdateStatus 测试
7. `refund_repo.go` — 编写 Create / FindByPaymentID / SumByPaymentID 测试
8. `booking_repo.go` L37 — 补充 `FindBookingByID` 测试
9. `analytics_handler.go` — 补充错误分支测试
10. `payment_service.go` — 补充 3 个未覆盖分支测试
11. `refund_service.go` — 补充 `amountCents <= 0` 错误分支测试
12. `service/chk/` — 迁出或排除出 coverage 统计

### 优先级 3 - Sprint 1/2 前端实联修复

13. `admin/login.vue` — 替换 `setTimeout` 为真实 API 调用
14. `admin/cruises/index.vue` — 调用 `GET /api/v1/admin/cruises` 替换硬编码
15. `admin/routes/index.vue` — 调用 `GET /api/v1/admin/routes` 替换硬编码
16. `admin/voyages/index.vue` — 调用 `GET /api/v1/admin/voyages` 替换硬编码
17. `admin/cabins/index.vue` — 实现完整舱位列表页
18. `admin/cabins/inventory.vue` — 实现完整库存管理页
19. `admin/cabins/pricing.vue` — 调用定价 API 替换硬编码
20. `admin/routes/new.vue` + `voyages/new.vue` — 实现表单提交逻辑

### 优先级 4 - Sprint 3/4 前端实联修复

21. `miniapp/login/login.vue` — 实现完整手机号+验证码登录
22. `miniapp/cabin/detail.vue` — 对接舱房详情 API
23. `admin/dashboard/index.vue` — 调用 analytics API 替换硬编码
24. `admin/finance/index.vue` — 调用财务 API 替换硬编码
25. `web/pay/[id].vue` + `PayButton.vue` — 实现支付流程 API 集成
26. `miniapp/pay/pay.vue` — 实现小程序支付流程

### 优先级 5 - 前端测试补全

27. 为修复后的每个页面编写含 API 调用断言、loading/error 状态验证的完整测试
28. 删除或重写 `sprint2-pages.spec.ts`、`dashboard.page.spec.ts`、`pay.page.spec.ts` 三个伪覆盖测试

### 优先级 6 - 代码质量改善（警告级）

29. `AuthService` 改为接收 `domain.StaffRepository` 接口
30. `web/stores/cruise.ts` 统一为 Composition API 风格
31. `web/booking/confirm.vue` 改用 `useApi()` composable 管理 token
