# Sprint 4.2 Audit Evidence

## P0

### Baseline Snapshot (2026-03-03)

- Backend baseline command:
  - `cd backend && go test ./... -count=1`
- Backend baseline result:
  - `FAIL` (`github.com/cruisebooking/backend/internal/config`)
  - Failure 1: `TestLoad` expected `8080`, got `:8080`
  - Failure 2: `TestLoadPanicsOnUnmarshalError` did not panic as expected

- Admin baseline command:
  - `cd frontend/admin && pnpm vitest run`
- Admin baseline result:
  - `PASS` (`37 files`, `103 tests`)
  - Residual warning risk: frequent `NuxtLink` unresolved warnings in unit tests

- Web baseline command:
  - `cd frontend/web && pnpm vitest run`
- Web baseline result:
  - `PASS` (`17 files`, `57 tests`)
  - Residual warning risk: `--localstorage-file` warning in node runtime

- Miniapp baseline command:
  - `cd frontend/miniapp && pnpm vitest run`
- Miniapp baseline result:
  - `PASS` (`13 files`, `27 tests`)
  - Residual warning risk: unresolved `swiper/swiper-item/scroll-view` component warnings in unit tests

### Task 1 Completion Evidence

- Plan checklist updated: `plan.md` Task 1 marked done.
- TODO evidence entry updated: link added to this file for follow-up remediation rounds.

### Task 2: Admin 缺失页接入真实 API + 三态

- 代码改动（后端 API 支撑）：
  - `backend/internal/repository/staff_repo.go`：补 `List` 方法。
  - `backend/internal/repository/shop_info_repo.go`：新增店铺信息单行读写仓储。
  - `backend/internal/repository/notification_template_repo.go`：新增通知模板 CRUD 仓储。
  - `backend/internal/service/notification_template_service.go`：新增通知模板业务服务。
  - `backend/internal/handler/notification_template_handler.go`：新增通知模板处理器。
  - `backend/internal/handler/staff_handler.go`：统一 `context.Context` 服务接口。
  - `backend/internal/handler/shop_info_handler.go`：统一 `context.Context` 服务接口。
  - `backend/internal/router/router.go`：新增 `/admin/staffs`、`/admin/shop-info`、`/admin/notification-templates` 路由。
  - `backend/cmd/server/main.go`：注入 staff/shop/template 新依赖。

- 代码改动（前端页面）：
  - `frontend/admin/app/pages/staff/index.vue`
  - `frontend/admin/app/pages/settings/shop.vue`
  - `frontend/admin/app/pages/notifications/templates.vue`
  - 三页均实现 `loading / error / empty` 三态，且走 `useApi().request` 真实后端路径。

- 测试新增：
  - `frontend/admin/tests/unit/pages/staff.index.spec.ts`
  - `frontend/admin/tests/unit/pages/settings.shop.spec.ts`
  - `frontend/admin/tests/unit/pages/notifications.templates.spec.ts`

- 验证命令与结果：
  - `cd backend && go test ./internal/router ./cmd/server ./internal/handler ./internal/repository ./internal/service -count=1` -> PASS
  - `cd frontend/admin && pnpm vitest run tests/unit/pages/staff.index.spec.ts tests/unit/pages/settings.shop.spec.ts tests/unit/pages/notifications.templates.spec.ts` -> PASS (`3 files`, `9 tests`)

### Task 3: Miniapp 舱位列表去硬编码 + 三态

- 代码改动：
  - `frontend/miniapp/pages/cabin/list.vue`：接入 `request('/cabins')` 真实数据请求，新增 `loading / error / empty / data` 分支渲染。
  - `frontend/miniapp/components/CabinCard.vue`：删除硬编码文案，改为 `props.item` 驱动渲染（名称、描述、价格）。

- 测试补齐：
  - `frontend/miniapp/tests/cabin-list.spec.ts`：新增 4 个用例覆盖 `loading`、`empty`、`error`、`data`。

- 验证命令与结果：
  - `cd frontend/miniapp && pnpm vitest run tests/cabin-list.spec.ts` -> PASS (`1 file`, `4 tests`)

### Task 4: 支付宝登录签名验签与 UID 防伪造

- 代码改动：
  - `backend/internal/service/user_auth_service.go`
  - `AlipayLogin` 从“直接返回客户端传入 UID”改为“验签后返回 provider UID”，新增：
    - 参数完整性校验
    - HMAC-SHA256 签名校验
    - 客户端 UID 与验签 UID 不一致拒绝（防伪造）
  - 新增错误语义：`ErrAlipayPayloadInvalid`、`ErrAlipaySignatureInvalid`、`ErrAlipayUIDMismatch`。

- 测试补齐：
  - `backend/internal/service/user_auth_service_test.go`
  - 新增/更新用例：
    - `TestUserAuthAlipayLogin`
    - `TestUserAuthAlipayLoginRejectsInvalidSignature`
    - `TestUserAuthAlipayLoginRejectsForgedClientUID`

- 验证命令与结果：
  - `cd backend && go test ./internal/service -run TestUserAuth -count=1` -> PASS

### Task 5: 账号绑定唯一性 + 绑定前身份确认

- 代码改动：
  - `backend/internal/service/user_auth_service.go`
  - `BindAccount` 新增：
    - 绑定前身份确认窗口校验（未确认直接拒绝）
    - 第三方账号唯一归属校验（同一 `provider:identifier` 禁止跨用户重复绑定）
    - 绑定成功后消费确认窗口（避免二次滥用）
  - 新增 `AuthorizeBinding` 用于“验证码/二次确认后授权绑定”。
  - 新增错误语义：`ErrBindingConfirmationRequired`、`ErrThirdPartyAlreadyBound`、`ErrBindPayloadInvalid`。

- 测试补齐：
  - `backend/internal/service/user_auth_service_test.go`
  - 新增/更新用例：
    - `TestUserAuthBindAccount`
    - `TestUserAuthBindAccountRequiresConfirmation`
    - `TestUserAuthBindAccountRejectsDuplicateIdentifier`

- 验证命令与结果：
  - `cd backend && go test ./internal/service -run TestUserAuth -count=1` -> PASS

### Task 6: 超时关单事务化 + 并发幂等

- 代码改动：
  - `backend/internal/service/order_timeout_service.go`
  - 关闭流程改为串行临界区执行，避免并发任务重复处理同一批过期订单。
  - 对每个订单新增失败回滚语义：
    - 先更新状态为 `cancelled`
    - 库存释放失败时，状态回滚到原状态（当前为 `pending_payment`）

- 测试新增：
  - `backend/internal/service/order_timeout_service_test.go`
  - 用例覆盖：
    - `TestCloseExpiredOrdersRollbackOnInventoryReleaseFailure`
    - `TestCloseExpiredOrdersConcurrentIdempotent`

- 验证命令与结果：
  - `cd backend && go test ./internal/service -run TestCloseExpiredOrders -count=1` -> PASS
  - `cd backend && go test ./internal/service -count=1` -> PASS

### Task 7: 订单状态机统一入口 + 状态日志同事务

- 代码改动：
  - `backend/internal/repository/booking_repo.go`
  - 新增 `TransitionStatus(ctx, id, status, operatorID, remark)`：
    - 事务内读取当前状态并执行 `CanTransitionTo` 校验
    - 状态更新与 `OrderStatusLog` 写入同事务提交
    - 非法流转返回 `ErrInvalidOrderStatusTransition`
  - `UpdateStatus` 改为统一转发至 `TransitionStatus`，避免绕过状态机。

- 入口收口：
  - `backend/internal/handler/booking_handler.go`
  - `AdminUpdate` 改走 `TransitionStatus`，并携带 `operator_id`（来自 `ContextKeyStaffID`）与 `remark`（空值时自动补默认文案）。
  - `BookingHandler.UpdateStatus` 去除空实现，改为转发到统一入口。

- 测试补齐：
  - `backend/internal/repository/booking_repo_test.go`
  - 新增用例：
    - `TestBookingRepoUpdateStatusWritesLog`
    - `TestBookingRepoUpdateStatusRejectsInvalidTransition`
  - `backend/internal/handler/p06_coverage_test.go`、`backend/internal/handler/delete_conflict_handler_test.go`：更新 mock 以适配统一入口签名。

- 验证命令与结果：
  - `cd backend && go test ./internal/repository -run TestBookingRepo -count=1` -> PASS
  - `cd backend && go test ./internal/handler -run TestBookingHandler_AdminUpdate -count=1` -> PASS
  - `cd backend && go test ./internal/... -count=1` -> 部分失败（仅 `internal/config` 既有基线失败，非本任务回归）

### Task 8: 订单导出真实实现 + 权限 + 上限 + CSV 注入防护

- 代码改动：
  - `backend/internal/service/order_export_service.go`
  - 导出内容由占位文本改为真实 CSV 数据（含表头与订单行）。
  - 新增导出权限校验（无权限直接拒绝）。
  - 新增导出上限（`5000` 行）控制，超限拒绝。
  - 新增 CSV 注入防护：单元格以 `=,+,-,@` 开头时自动前缀 `'`。

- 测试新增：
  - `backend/internal/service/order_export_service_test.go`
  - 用例覆盖：
    - `TestOrderExportServiceDeniedWithoutPermission`
    - `TestOrderExportServiceRejectsOverLimit`
    - `TestOrderExportServiceSanitizesCSVInjection`

- 验证命令与结果：
  - `cd backend && go test ./internal/service -run TestOrderExportService -count=1` -> PASS
  - `cd backend && go test ./internal/service -count=1` -> PASS

### Task 9: 通知模板 text/template + 变量白名单

- 代码改动：
  - `backend/internal/domain/notification_template.go`
  - `Render` 从字符串替换改为 `text/template` 渲染（`missingkey=error`）。
  - 新增模板变量白名单校验，禁止未授权变量进入模板执行。
  - 新增错误语义：
    - `ErrNotificationTemplateInvalid`
    - `ErrNotificationTemplateVarNotAllowed`

- 测试补齐：
  - `backend/internal/service/notify_template_test.go`
  - 新增/更新用例：
    - `TestNotifyTemplateRender`
    - `TestNotifyTemplateRender_Multiple`
    - `TestNotifyTemplateRenderRejectsUnknownVariable`
    - `TestNotifyTemplateRenderRejectsInvalidSyntax`

- 验证命令与结果：
  - `cd backend && go test ./internal/service -run TestNotifyTemplateRender -count=1` -> PASS
  - `cd backend && go test ./internal/domain -count=1` -> PASS
  - `cd backend && go test ./internal/service -count=1` -> PASS

### Task 10: Analytics 真实查询 + 索引 + 性能说明

- 代码改动：
  - `backend/internal/repository/analytics_repo.go`
  - 由空实现改为真实 SQL 查询：
    - `Trend(days)`：按天聚合销售额与订单数（SQLite 递归 CTE / PostgreSQL generate_series）
    - `CabinHotnessRanking(limit)`：按舱位聚合销量与热度
    - `InventoryOverview()`：统计总舱位、低库存、售罄数
    - `PageViewStats()`：基于业务表聚合页面访问指标

- 索引与性能改动：
  - `backend/migrations/000012_analytics_indexes.up.sql`
  - `backend/migrations/000012_analytics_indexes.down.sql`
  - 新增仪表盘相关索引，覆盖 `bookings.created_at/status/cabin_sku_id`、`payments.status/created_at`、`cabin_inventories.alert_threshold`。

- 测试补齐：
  - `backend/internal/repository/analytics_repo_test.go`
  - 新增用例：
    - `TestAnalyticsRepository_Trend`
    - `TestAnalyticsRepository_CabinHotnessRanking`
    - `TestAnalyticsRepository_InventoryOverviewAndPageViewStats`

- 验证命令与结果：
  - `cd backend && go test ./internal/repository -run TestAnalyticsRepository -count=1` -> PASS
  - `cd backend && go test ./internal/repository -count=1` -> PASS

### Task 11: 批量操作安全加固（已完成）

- 批量数量上限 + 审计日志
  - `backend/internal/handler/bulk_limits.go`：统一批量上限常量。
  - `backend/internal/handler/cruise_handler.go`：`BatchUpdateStatus` 增加上限校验与审计日志。
  - `backend/internal/handler/cabin_handler.go`：`BatchUpdateStatus` 增加上限校验与审计日志。
  - `backend/internal/handler/cruise_handler_extended_test.go`：新增 `TestCruiseHandler_BatchUpdateStatusRejectsOversize`。
  - `backend/internal/handler/cabin_handler_extended_test.go`：新增 `TestCabinHandler_BatchStatusRejectsOversize`。

- 关键批量写事务化与失败回滚
  - `backend/internal/repository/cruise_repo.go`：新增事务化 `BatchUpdateStatus`，并在 `RowsAffected != len(ids)` 时整体回滚。
  - `backend/internal/repository/cabin_repo.go`：`BatchUpdateStatus` 改为事务化，并在目标数量不匹配时回滚。
  - `backend/internal/service/cruise_service.go`：优先走仓储批量事务接口（若实现）。
  - `backend/internal/repository/cruise_repo_test.go`：新增 `TestCruiseRepository_BatchUpdateStatusRollbackOnPartialMatch`。
  - `backend/internal/repository/cabin_repo_extended_test.go`：新增 `TestCabinRepository_BatchUpdateStatusRollbackOnPartialMatch`。

- 权限与审计日志校验用例
  - `backend/internal/handler/cruise_handler_extended_test.go`：新增 `TestCruiseHandler_BatchUpdateStatusWritesAuditLog`。
  - `backend/internal/handler/cabin_handler_extended_test.go`：新增 `TestCabinHandler_BatchStatusWritesAuditLog`。
  - `backend/internal/router/router_test.go`：新增 `TestSetup_ProtectedBatchEndpointsRequireAuthAndRole`，覆盖 `401`（无 JWT）与 `403`（无角色）场景。

- 验证命令与结果：
  - `cd backend && go test ./internal/handler -run "TestCruiseHandler_BatchUpdateStatus|TestCabinHandler_BatchStatus" -count=1` -> PASS
  - `cd backend && go test ./internal/repository -run "TestCruiseRepository_BatchUpdateStatusRollbackOnPartialMatch|TestCabinRepository_BatchUpdateStatusRollbackOnPartialMatch" -count=1` -> PASS
  - `cd backend && go test ./internal/handler -run "TestCruiseHandler_BatchUpdateStatusWritesAuditLog|TestCabinHandler_BatchStatusWritesAuditLog" -count=1` -> PASS
  - `cd backend && go test ./internal/router -run TestSetup_ProtectedBatchEndpointsRequireAuthAndRole -count=1` -> PASS
  - `cd backend && go test ./internal/repository -count=1` -> PASS
  - `cd backend && go test ./internal/handler -count=1` -> PASS
  - `cd backend && go test ./internal/service -count=1` -> PASS

### Task 12: AvailableWithAlert 与阈值边界测试

- 代码改动：
  - `backend/internal/service/inventory_alert_service.go`
  - 新增 `AvailableWithAlert(ctx, skuID)`，返回 `(available, isAlert)`。
  - 预警判定规则明确为：`threshold > 0 && available <= threshold`。

- 测试补齐：
  - `backend/internal/service/inventory_alert_service_test.go`
  - 新增 `TestInventoryAlertService_AvailableWithAlertBoundaries`，覆盖：
    - `available == threshold`
    - `available < threshold`
    - `threshold = 0`

- 验证命令与结果：
  - `cd backend && go test ./internal/service -run TestInventoryAlertService -count=1` -> PASS
  - `cd backend && go test ./internal/service -count=1` -> PASS

### Task 13: 批量日期范围定价 API 暴露与接口测试

- 代码改动：
  - `backend/internal/repository/cabin_repo.go`
  - 新增 `BatchSetPrice(...)` 仓储实现，事务内按日期区间批量写入价格。
  - `backend/internal/service/cabin_admin_service.go`
  - 新增 `BatchSetPrice(...)` 服务入口并转发仓储。
  - `backend/internal/handler/cabin_handler.go`
  - 新增 `POST /api/v1/admin/cabins/:id/prices/batch` 处理器：
    - 参数校验（起止日期、区间顺序、occupancy）
    - 调用 `BatchSetPrice`
  - `backend/internal/router/router.go`
  - 注册新路由 `:id/prices/batch`。

- 测试补齐：
  - `backend/internal/handler/cabin_handler_extended_test.go`：新增 `TestCabinHandler_BatchSetPrice`（成功 + 日期区间非法）。
  - `backend/internal/handler/all_handler_test.go`、`backend/internal/service/cabin_admin_service_test.go`、`backend/internal/service/p06_coverage_test.go`、`backend/internal/handler/delete_conflict_handler_test.go`：补齐接口扩展后的 mock 方法。

- 验证命令与结果：
  - `cd backend && go test ./internal/handler -run TestCabinHandler_BatchSetPrice -count=1` -> PASS
  - `cd backend && go test ./internal/handler -count=1` -> PASS
  - `cd backend && go test ./internal/service -count=1` -> PASS
  - `cd backend && go test ./internal/router -count=1` -> PASS

## P1

### Task 14: 对账日报幂等与并发测试

- 代码改动：
  - `backend/internal/service/reconciliation_service.go`
  - 新增按 `date(YYYY-MM-DD)` 的幂等控制策略：同日期重复生成直接拒绝。
  - 新增并发安全保护：并发同日期请求最终仅允许一个成功，其余返回重复生成错误。
  - 新增错误语义：`ErrReconciliationReportAlreadyGenerated`。

- 测试补齐：
  - `backend/internal/service/reconciliation_service_test.go`
  - 新增：
    - `TestReconciliationService_GenerateDailyReportRejectsDuplicateDate`
    - `TestReconciliationService_GenerateDailyReportConcurrentIdempotent`

- 验证命令与结果：
  - `cd backend && go test ./internal/service -run TestReconciliationService_GenerateDailyReport -count=1` -> PASS
  - `cd backend && go test ./internal/service -count=1` -> PASS

### Task 15: 员工角色变更同步 Casbin + 审计日志

- 代码改动：
  - `backend/internal/service/staff_service.go`
  - `AssignRole` 扩展为 `AssignRole(ctx, id, role, operatorID)`，角色更新后触发：
    - Casbin 角色同步
    - 审计日志写入（含 `operatorID`、目标员工、旧角色、新角色）
  - `backend/internal/service/staff_role_sync.go`
  - 新增 Casbin 同步实现：更新 grouping policy 并持久化策略。
  - `backend/internal/service/staff_audit_logger.go`
  - 新增操作日志写入器，落表 `operation_logs`。
  - `backend/internal/repository/operation_log_repo.go`
  - 新增操作日志仓储 `Create`。
  - `backend/internal/handler/staff_handler.go`
  - `AssignRole` 从 JWT 上下文提取操作者 `staffID` 并透传到服务层。
  - `backend/cmd/server/main.go`
  - 注入 `CasbinStaffRoleSync` 与 `StaffOperationLogger` 到 `StaffService`。

- 测试补齐：
  - `backend/internal/service/staff_service_test.go`
  - 新增：
    - `TestStaffServiceAssignRoleSyncsCasbinAndWritesAudit`
    - `TestStaffServiceAssignRoleReturnsErrorWhenCasbinSyncFails`

- 验证命令与结果：
  - `cd backend && go test ./internal/service -run TestStaffServiceAssignRole -count=1` -> PASS
  - `cd backend && go test ./internal/handler -count=1` -> PASS
  - `cd backend && go test ./internal/router -count=1` -> PASS
  - `cd backend && go test ./cmd/server -count=1` -> PASS

### Task 16: ShopInfo 单行约束（DB + Service）

- 代码改动：
  - `backend/internal/service/shop_info_service.go`
  - 新增服务层单行约束：
    - 非空写入时仅允许 `ID=0/1`
    - 强制归一化为固定主键 `ID=1`
  - 新增错误语义：`ErrInvalidShopInfoSingletonID`。
  - `backend/migrations/000013_shop_info_singleton.up.sql`
  - 新增数据库约束：`CHECK (id = 1)`，从 DB 层防止多行主键写入。
  - `backend/migrations/000013_shop_info_singleton.down.sql`
  - 提供回滚删除约束脚本。

- 测试补齐：
  - `backend/internal/service/shop_info_service_test.go`
  - 新增/更新：
    - `TestShopInfoService_Update` 断言更新后固定 `ID=1`
    - `TestShopInfoService_UpdateRejectsNonSingletonID`

- 验证命令与结果：
  - `cd backend && go test ./internal/service -run TestShopInfoService_Update -count=1` -> PASS
  - `cd backend && go test ./internal/service -run "TestShopInfoService|TestStaffServiceAssignRole|TestReconciliationService_GenerateDailyReport" -count=1` -> PASS
  - `cd backend && go test ./cmd/server ./internal/handler ./internal/router -count=1` -> PASS

### Task 17: 库存预警通知去重窗口

- 代码改动：
  - `backend/internal/service/inventory_alert_service.go`
  - 扩展库存预警服务支持通知发送依赖与去重窗口：
    - 新增通知接口 `inventoryAlertNotifier`
    - 新增构造器 `NewInventoryAlertServiceWithNotify`
    - 新增 `ScanAndNotify(...)`：对低库存项发送通知，并基于去重键（`sku:{id}`）与时间窗口抑制重复通知

- 测试补齐：
  - `backend/internal/service/inventory_alert_service_test.go`
  - 新增：
    - `TestInventoryAlertService_ScanAndNotifyDeduplicatesWithinWindow`

- 验证命令与结果：
  - `cd backend && go test ./internal/service -run TestInventoryAlertService_ScanAndNotifyDeduplicatesWithinWindow -count=1` -> PASS
  - `cd backend && go test ./internal/service -run TestInventoryAlertService -count=1` -> PASS

## P2

### Task 18: 退款边界测试补齐（30/7/0 天）

- 测试补齐：
  - `backend/internal/service/refund_service_tiered_test.go`
  - 新增：
    - `TestTieredRefund_BoundaryExactDays30_7_0`
    - `TestTieredRefund_IntegerArithmeticNoPrecisionLoss`
  - 明确覆盖：恰好 `30` 天、恰好 `7` 天、`0` 天，以及非整除金额（`9999 * 80%`）的整数计算结果。

- 验证命令与结果：
  - `cd backend && go test ./internal/service -run TestTieredRefund -count=1` -> PASS

### Task 19: 清理空实现与伪覆盖

- 代码改动：
  - `backend/internal/handler/task19_no_empty_impl_guard_test.go`
  - 新增守卫测试 `TestTask19_NoEmptyFunctionBodiesInProduction`：
    - 自动扫描 `backend/internal` 下生产代码（排除 `*_test.go`）
    - 发现空函数体（`{}`）即失败，防止以空实现替代核心逻辑

- 巡检结论：
  - `UpdateStatus` 空实现已在前序任务中移除，当前 `BookingHandler.UpdateStatus` 已转发至真实状态机入口。
  - 本轮扫描未发现新的生产代码空函数体风险。

- 验证命令与结果：
  - `cd backend && go test ./internal/handler -run TestTask19_NoEmptyFunctionBodiesInProduction -count=1` -> PASS
  - `cd backend && go test ./internal/handler ./internal/service ./internal/router ./cmd/server -count=1` -> PASS

### Task 20: Sprint 4.2 27 Tasks 对照复核证据（已完成）

- 产出文档：
  - `docs/plans/2026-03-03-sprint42-27tasks-mapping.md`
  - 对照 `docs/plans/2026-02-22-sprint04.2.md` 全量 27 个任务建立映射矩阵：
    - `Task -> Status(complete/partial/gap) -> Code Evidence -> Test Evidence -> UI Evidence`

- 当前结论（持续补齐后）：
  - `complete`: 27
  - `partial`: 0
  - `gap`: 0

- 关键补齐项（本轮新增）：
  - 已补充 UI 证据：
    - `docs/plans/evidence/task20/task8-cruises.png`
    - `docs/plans/evidence/task20/task8-cabins.png`
    - `docs/plans/evidence/task20/task9-pricing.png`
    - `docs/plans/evidence/task20/task21-bookings-export.png`
    - `docs/plans/evidence/task20/task22-staff.png`
    - `docs/plans/evidence/task20/task23-shop.png`
    - `docs/plans/evidence/task20/task24-templates.png`
    - `docs/plans/evidence/task20/task26-dashboard.png`
  - 已补代码与测试：
    - Task 10：`frontend/miniapp/pages/cabin/list.vue`、`frontend/miniapp/components/CabinCard.vue`、`frontend/shared/components/InventoryBadge.vue`，并新增/更新 `tests/cabin-list.spec.ts`、`tests/cabin-detail.spec.ts`、`frontend/web/tests/unit/pages/cabins/index.spec.ts`。
    - Task 15：新增 `frontend/web/app/pages/orders/index.vue`、`frontend/miniapp/pages/orders/list.vue` 与对应单测。
    - Task 16：新增 `backend/migrations/000010_sprint42_user_passenger_extend.*.sql`，并由 `backend/migrations/migrations_sprint42_test.go` 覆盖存在性。
    - Task 21：增强 `frontend/admin/app/pages/bookings/index.vue`（筛选/tab/导出/操作），并补 `frontend/admin/tests/unit/bookings.list.spec.ts`。

- 本轮验证命令与结果：
  - `cd backend && go test ./migrations -run TestSprint42MigrationFilesExist -count=1` -> PASS
  - `cd frontend/admin && pnpm vitest run tests/unit/pages/cabins-index.spec.ts tests/unit/pages/cabins-pricing.spec.ts tests/unit/pages/notifications.templates.spec.ts` -> PASS
  - `cd frontend/web && pnpm vitest run tests/unit/pages/orders/index.spec.ts tests/unit/pages/cabins/index.spec.ts` -> PASS
  - `cd frontend/miniapp && pnpm vitest run tests/orders-list.spec.ts tests/cabin-list.spec.ts tests/cabin-detail.spec.ts` -> PASS
  - `cd frontend/admin && pnpm vitest run tests/unit/bookings.list.spec.ts` -> PASS
  - `cd frontend/admin && pnpm test:e2e` -> PASS（`3 passed`）

- 收口说明：
  - Task 8/9/24/27 已从 `partial` 收口为 `complete`，Task 20 对照复核结论同步更新为“27 项全 complete”。

### Task 21: Admin 目录迁移与路由冒烟

- 代码改动：
  - 新增：
    - `frontend/admin/app/pages/dashboard/index.vue`
    - `frontend/admin/app/pages/finance/index.vue`
  - 删除：
    - `frontend/admin/pages/dashboard/index.vue`
    - `frontend/admin/pages/finance/index.vue`
  - 测试更新：
    - `frontend/admin/tests/unit/dashboard.page.spec.ts`（改为引用 `app/pages/dashboard/index.vue`）
    - `frontend/admin/tests/unit/pages/finance-index.spec.ts`（改为引用 `app/pages/finance/index.vue`）
    - `frontend/admin/tests/unit/pages/task21-route-smoke.spec.ts`（新增）

- 路由与冲突验证：
  - 通过 `task21-route-smoke.spec.ts` 验证 `app/pages/login.vue`、`app/pages/dashboard/index.vue`、`app/pages/finance/index.vue` 页面存在。
  - 验证旧目录 `pages/dashboard/index.vue`、`pages/finance/index.vue` 已删除，避免重复路由冲突。

- 验证命令与结果：
  - `cd frontend/admin && pnpm vitest run tests/unit/dashboard.page.spec.ts tests/unit/pages/finance-index.spec.ts tests/unit/pages/task21-route-smoke.spec.ts` -> PASS
  - `cd frontend/admin && pnpm vitest run tests/unit/pages/login.spec.ts tests/unit/routes.list.spec.ts` -> PASS

### Task 22: 全量回归 + TODO 收口（已执行）

- 后端全量回归：
  - `cd backend && go test ./internal/config -count=1` -> PASS
  - `cd backend && go test ./... -count=1` -> PASS

- 关键专项回归：
  - `cd backend && go test ./internal/router -run TestSetup_ProtectedBatchEndpointsRequireAuthAndRole -count=1` -> PASS
  - `cd backend && go test ./internal/service -run "TestOrderExportServiceDeniedWithoutPermission|TestOrderExportServiceRejectsOverLimit|TestOrderExportServiceSanitizesCSVInjection" -count=1` -> PASS

- 前端全量回归：
  - `cd frontend/admin && pnpm vitest run` -> PASS（`41 files`, `114 tests`）
  - `cd frontend/web && pnpm vitest run` -> PASS（`17 files`, `57 tests`）
  - `cd frontend/miniapp && pnpm vitest run` -> PASS（`13 files`, `30 tests`）

- E2E 检查：
  - 新增 Playwright 配置与关键 E2E：
    - `frontend/admin/playwright.config.ts`
    - `frontend/admin/tests/e2e/task22-auth-and-dashboard.spec.ts`
    - `frontend/admin/tests/e2e/task20-ui-evidence.spec.ts`
  - 执行：`cd frontend/admin && pnpm test:e2e` -> PASS（`3 passed`）
  - 产出截图：
    - `docs/plans/evidence/task22/admin-cruises.png`
    - `docs/plans/evidence/task22/admin-dashboard.png`
    - `docs/plans/evidence/task20/task21-bookings-export.png`

- 前端三态收口（loading/error/empty）：
  - 本轮补齐页面：
    - `frontend/admin/app/pages/dashboard/index.vue`
    - `frontend/admin/app/pages/cabins/pricing.vue`
    - `frontend/admin/app/pages/bookings/[id].vue`
    - `frontend/admin/app/pages/cabins/[id].vue`
    - `frontend/admin/app/pages/cruises/[id].vue`
    - `frontend/admin/app/pages/routes/[id].vue`
    - `frontend/admin/app/pages/voyages/[id].vue`
    - `frontend/admin/app/pages/facilities/[id].vue`
    - `frontend/admin/app/pages/facility-categories/[id].vue`
    - `frontend/admin/app/pages/cabin-types/[id].vue`
    - `frontend/admin/app/pages/cabin-types/new.vue`
    - `frontend/admin/app/pages/facilities/new.vue`
  - 验证命令与结果：
    - `cd frontend/admin && pnpm vitest run tests/unit/pages/cabins-id.spec.ts tests/unit/pages/cruises-id.spec.ts tests/unit/pages/routes-id.spec.ts tests/unit/pages/voyages-id.spec.ts tests/unit/pages/cabins-pricing.spec.ts tests/unit/pages/facilities-form.spec.ts tests/unit/pages/facility-categories-form.spec.ts tests/unit/pages/cabin-types-form.spec.ts tests/unit/pages/facilities-new.spec.ts` -> PASS（`9 files`, `28 tests`）
    - `frontend/admin` 数据加载页三态扫描：PASS（无缺失项）。

- 收口结论：
  - 已完成“后端单测/前端单测/关键 E2E/权限回归/导出安全回归”执行与结果归档。
  - `TODO.md` 中 P0 与 27-task 复核条目已同步收口。

## 回归测试汇总

- Backend:
  - `go test ./internal/config -count=1` -> PASS
  - `go test ./... -count=1` -> PASS
  - `go test ./internal/router -run TestSetup_ProtectedBatchEndpointsRequireAuthAndRole -count=1` -> PASS
  - `go test ./internal/service -run "TestOrderExportServiceDeniedWithoutPermission|TestOrderExportServiceRejectsOverLimit|TestOrderExportServiceSanitizesCSVInjection" -count=1` -> PASS

- Frontend:
  - `frontend/admin`: PASS（`41 files`, `114 tests`）
  - `frontend/web`: PASS（`17 files`, `57 tests`）
  - `frontend/miniapp`: PASS（`13 files`, `30 tests`）
  - `frontend/admin e2e`: PASS（`3 tests`）

- Final gate status:
  - Ready for review（整改项与证据已闭环）

## 联调记录

- Task 20/27 映射已收口为 `complete=27, partial=0, gap=0`。
