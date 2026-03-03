# Sprint 4.2 Audit TODO Checklist

> 审计日期: 2026-03-02
> 用途: 按阻塞优先级逐项整改并勾选，作为复审核验清单。
> 执行记录: 基线测试与证据骨架已完成（2026-03-03），见 `docs/plans/2026-03-02-sprint42-audit-evidence.md`。

## P0 Blockers (必须先完成)

- [x] 前端页面实联未完成（存在空壳/缺页）
  - [x] `frontend/admin/app/pages` 已覆盖 cruises/routes/voyages/cabins/bookings 等页面，但仍缺员工管理、店铺设置、通知模板等页面，需补齐并接后端真实 API。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-2-admin-缺失页接入真实-api--三态`）
  - [x] `frontend/miniapp/pages/cabin/list.vue` + `frontend/miniapp/components/CabinCard.vue` 去除硬编码假数据，改为真实 API 渲染。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-3-miniapp-舱位列表去硬编码--三态`）
  - [x] 所有页面统一具备 `loading` / `error` / `empty` 三态。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-22-全量回归--todo-收口已执行`）

- [x] 支付宝登录安全未实现
  - [x] `backend/internal/service/user_auth_service.go` 的 `AlipayLogin` 增加支付宝回调签名验签。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-4-支付宝登录签名验签与-uid-防伪造`）
  - [x] 防伪造 `alipay_uid`（不得直接信任客户端传入 UID）。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-4-支付宝登录签名验签与-uid-防伪造`）

- [x] 账号绑定安全未实现
  - [x] `BindAccount` 验证目标第三方账号未被其他用户绑定。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-5-账号绑定唯一性--绑定前身份确认`）
  - [x] 绑定前增加身份确认（验证码/二次确认机制）。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-5-账号绑定唯一性--绑定前身份确认`）

- [x] 订单超时关单不具备原子性与并发安全
  - [x] `CloseExpiredOrders` 改为事务内执行（关单 + 释放库存 + 日志）。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-6-超时关单事务化--并发幂等`）
  - [x] 增加并发锁/幂等防重，避免并发重复关单。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-6-超时关单事务化--并发幂等`）
  - [x] 新增服务测试覆盖并发/失败回滚路径。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-6-超时关单事务化--并发幂等`）

- [x] 订单状态变更日志未落地
  - [x] 每次状态变更强制写 `OrderStatusLog`（含 `operator_id` + `remark`）。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-7-订单状态机统一入口--状态日志同事务`）
  - [x] 状态更新与日志写入必须同一事务提交。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-7-订单状态机统一入口--状态日志同事务`）
  - [x] 建立统一状态机入口，禁止绕过 `CanTransitionTo` 直接改状态。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-7-订单状态机统一入口--状态日志同事务`）

- [x] Excel 导出仍为占位实现
  - [x] `backend/internal/service/order_export_service.go` 实现真实导出逻辑，替换 `"Excel content placeholder"`。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-8-订单导出真实实现--权限--上限--csv-注入防护`）
  - [x] 增加导出权限控制。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-8-订单导出真实实现--权限--上限--csv-注入防护`）
  - [x] 增加导出数量上限。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-8-订单导出真实实现--权限--上限--csv-注入防护`）
  - [x] 增加 CSV/公式注入防护（以 `=,+,-,@` 开头内容处理）。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-8-订单导出真实实现--权限--上限--csv-注入防护`）

- [x] 通知模板渲染存在注入风险
  - [x] `backend/internal/domain/notification_template.go` 改用 `text/template` 安全渲染。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-9-通知模板-texttemplate--变量白名单`）
  - [x] 增加模板变量白名单与异常模板测试，防 SSTI 风险。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-9-通知模板-texttemplate--变量白名单`）

- [x] 数据看板增强未完成（核心接口返回空数据）
  - [x] `Trend` 实现真实 7/30 天趋势查询。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-10-analytics-真实查询--索引--性能说明`）
  - [x] `CabinHotnessRanking` 实现真实排行查询。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-10-analytics-真实查询--索引--性能说明`）
  - [x] `InventoryOverview` 实现真实库存概览。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-10-analytics-真实查询--索引--性能说明`）
  - [x] `PageViewStats` 实现真实访问量统计。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-10-analytics-真实查询--索引--性能说明`）

## P1 High Priority

- [x] 批量上架/下架/删除安全加固
  - [x] 批量操作增加数量上限（防误操作/资源风险）。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-11-批量操作安全加固已完成`）
  - [x] 关键批量写操作加事务与失败回滚策略。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-11-批量操作安全加固已完成`）
  - [x] 补齐权限与审计日志校验用例。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-11-批量操作安全加固已完成`）

- [x] 库存预警能力与清单不一致
  - [x] 增加 `AvailableWithAlert` 方法（返回可用库存 + 是否触发预警）。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-12-availablewithalert-与阈值边界测试`）
  - [x] 覆盖阈值边界测试（`==`、`<`、`threshold=0`）。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-12-availablewithalert-与阈值边界测试`）

- [x] 批量日期范围定价 API 缺失
  - [x] 在 handler/router 暴露批量日期定价接口（调用 `BatchSetPrice`）。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-13-批量日期范围定价-api-暴露与接口测试`）
  - [x] 增加接口级测试覆盖日期区间与参数校验。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-13-批量日期范围定价-api-暴露与接口测试`）

- [x] 财务对账日报缺防重复生成机制
  - [x] 以 `date` 做幂等控制（存在即拒绝或覆盖策略明确）。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-14-对账日报幂等与并发测试`）
  - [x] 增加重复请求/并发生成测试。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-14-对账日报幂等与并发测试`）

- [x] 员工角色分配未与 Casbin 同步
  - [x] 角色变更同步更新 Casbin policy/grouping。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-15-员工角色变更同步-casbin--审计日志`）
  - [x] 权限变更写入操作日志（含操作者与变更内容）。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-15-员工角色变更同步-casbin--审计日志`）

- [x] ShopInfo 单行配置约束不足
  - [x] 明确单行模式: 仅允许更新，不允许创建多条。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-16-shopinfo-单行约束db--service`）
  - [x] 数据库层与服务层双重约束（唯一键/固定主键策略）。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-16-shopinfo-单行约束db--service`）

- [x] 库存预警通知缺防重复机制
  - [x] 增加通知去重键与时间窗口，避免重复轰炸。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-17-库存预警通知去重窗口`）
  - [x] 增加重复扫描测试。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-17-库存预警通知去重窗口`）

## P2 Completeness / Test Quality

- [x] 前端目录结构整理（admin）
  - [x] 将 `frontend/admin/pages/dashboard` 与 `frontend/admin/pages/finance` 迁移到 `frontend/admin/app/pages`，统一页面目录结构。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-21-admin-目录迁移与路由冒烟`）
  - [x] 迁移后确认路由不变（`/dashboard`、`/finance`）且无重复路由冲突。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-21-admin-目录迁移与路由冒烟`）
  - [x] 迁移后执行 admin 端页面冒烟测试（登录、仪表盘、财务页可访问）。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-21-admin-目录迁移与路由冒烟`）

- [x] 阶梯退款边界测试补齐
  - [x] 明确覆盖恰好 30 天、恰好 7 天、0 天边界。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-18-退款边界测试补齐3070-天`）
  - [x] 验证整数运算无精度丢失。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-18-退款边界测试补齐3070-天`）

- [x] 清理空实现/伪覆盖风险
  - [x] 移除 `backend/internal/handler/booking_handler.go` 中测试用空实现 `UpdateStatus`，改为真实逻辑。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-19-清理空实现与伪覆盖`）
  - [x] 继续巡检并禁止以空函数体替代核心业务实现。（证据：`docs/plans/2026-03-02-sprint42-audit-evidence.md#task-19-清理空实现与伪覆盖`）

- [x] Sprint 4.2 Part A-E / 27 Tasks 逐项复核
  - [x] 对照 `docs/plans/2026-02-22-sprint04.2.md` 建立任务到代码映射。（证据：`docs/plans/2026-03-03-sprint42-27tasks-mapping.md`）
  - [x] 每个 Task 提供“代码位置 + 测试用例 + 页面联调截图/录屏”证据。（证据：`docs/plans/2026-03-03-sprint42-27tasks-mapping.md`，`complete=27`）

## Re-Review Gate

- [x] 全部 P0 项完成后方可申请复审。
- [x] 复审前执行：后端单测、前端单测、关键 E2E、权限回归、导出安全回归。
  - [x] 后端单测已执行并通过（`go test ./...`），见证据 Task 22。
  - [x] 前端单测已执行（admin/web/miniapp 全量通过），见证据 Task 22。
  - [x] 关键 E2E：`frontend/admin` Playwright 已新增并通过（`pnpm test:e2e`，`3 passed`）。
  - [x] 权限回归已执行并通过（`TestSetup_ProtectedBatchEndpointsRequireAuthAndRole`）。
  - [x] 导出安全回归已执行并通过（`TestOrderExportServiceDeniedWithoutPermission` 等）。
- [ ] 任一项未通过，继续保持阻塞状态。
