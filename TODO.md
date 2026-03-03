# Sprint 4.2 Audit TODO Checklist

> 审计日期: 2026-03-02
> 用途: 按阻塞优先级逐项整改并勾选，作为复审核验清单。

## P0 Blockers (必须先完成)

- [ ] 前端页面实联未完成（存在空壳/缺页）
  - [ ] `frontend/admin/app/pages` 已覆盖 cruises/routes/voyages/cabins/bookings 等页面，但仍缺员工管理、店铺设置、通知模板等页面，需补齐并接后端真实 API。
  - [ ] `frontend/miniapp/pages/cabin/list.vue` + `frontend/miniapp/components/CabinCard.vue` 去除硬编码假数据，改为真实 API 渲染。
  - [ ] 所有页面统一具备 `loading` / `error` / `empty` 三态。

- [ ] 支付宝登录安全未实现
  - [ ] `backend/internal/service/user_auth_service.go` 的 `AlipayLogin` 增加支付宝回调签名验签。
  - [ ] 防伪造 `alipay_uid`（不得直接信任客户端传入 UID）。

- [ ] 账号绑定安全未实现
  - [ ] `BindAccount` 验证目标第三方账号未被其他用户绑定。
  - [ ] 绑定前增加身份确认（验证码/二次确认机制）。

- [ ] 订单超时关单不具备原子性与并发安全
  - [ ] `CloseExpiredOrders` 改为事务内执行（关单 + 释放库存 + 日志）。
  - [ ] 增加并发锁/幂等防重，避免并发重复关单。
  - [ ] 新增服务测试覆盖并发/失败回滚路径。

- [ ] 订单状态变更日志未落地
  - [ ] 每次状态变更强制写 `OrderStatusLog`（含 `operator_id` + `remark`）。
  - [ ] 状态更新与日志写入必须同一事务提交。
  - [ ] 建立统一状态机入口，禁止绕过 `CanTransitionTo` 直接改状态。

- [ ] Excel 导出仍为占位实现
  - [ ] `backend/internal/service/order_export_service.go` 实现真实导出逻辑，替换 `"Excel content placeholder"`。
  - [ ] 增加导出权限控制。
  - [ ] 增加导出数量上限。
  - [ ] 增加 CSV/公式注入防护（以 `=,+,-,@` 开头内容处理）。

- [ ] 通知模板渲染存在注入风险
  - [ ] `backend/internal/domain/notification_template.go` 改用 `text/template` 安全渲染。
  - [ ] 增加模板变量白名单与异常模板测试，防 SSTI 风险。

- [ ] 数据看板增强未完成（核心接口返回空数据）
  - [ ] `Trend` 实现真实 7/30 天趋势查询。
  - [ ] `CabinHotnessRanking` 实现真实排行查询。
  - [ ] `InventoryOverview` 实现真实库存概览。
  - [ ] `PageViewStats` 实现真实访问量统计。

## P1 High Priority

- [ ] 批量上架/下架/删除安全加固
  - [ ] 批量操作增加数量上限（防误操作/资源风险）。
  - [ ] 关键批量写操作加事务与失败回滚策略。
  - [ ] 补齐权限与审计日志校验用例。

- [ ] 库存预警能力与清单不一致
  - [ ] 增加 `AvailableWithAlert` 方法（返回可用库存 + 是否触发预警）。
  - [ ] 覆盖阈值边界测试（`==`、`<`、`threshold=0`）。

- [ ] 批量日期范围定价 API 缺失
  - [ ] 在 handler/router 暴露批量日期定价接口（调用 `BatchSetPrice`）。
  - [ ] 增加接口级测试覆盖日期区间与参数校验。

- [ ] 财务对账日报缺防重复生成机制
  - [ ] 以 `date` 做幂等控制（存在即拒绝或覆盖策略明确）。
  - [ ] 增加重复请求/并发生成测试。

- [ ] 员工角色分配未与 Casbin 同步
  - [ ] 角色变更同步更新 Casbin policy/grouping。
  - [ ] 权限变更写入操作日志（含操作者与变更内容）。

- [ ] ShopInfo 单行配置约束不足
  - [ ] 明确单行模式: 仅允许更新，不允许创建多条。
  - [ ] 数据库层与服务层双重约束（唯一键/固定主键策略）。

- [ ] 库存预警通知缺防重复机制
  - [ ] 增加通知去重键与时间窗口，避免重复轰炸。
  - [ ] 增加重复扫描测试。

## P2 Completeness / Test Quality

- [ ] 前端目录结构整理（admin）
  - [ ] 将 `frontend/admin/pages/dashboard` 与 `frontend/admin/pages/finance` 迁移到 `frontend/admin/app/pages`，统一页面目录结构。
  - [ ] 迁移后确认路由不变（`/dashboard`、`/finance`）且无重复路由冲突。
  - [ ] 迁移后执行 admin 端页面冒烟测试（登录、仪表盘、财务页可访问）。

- [ ] 阶梯退款边界测试补齐
  - [ ] 明确覆盖恰好 30 天、恰好 7 天、0 天边界。
  - [ ] 验证整数运算无精度丢失。

- [ ] 清理空实现/伪覆盖风险
  - [ ] 移除 `backend/internal/handler/booking_handler.go` 中测试用空实现 `UpdateStatus`，改为真实逻辑。
  - [ ] 继续巡检并禁止以空函数体替代核心业务实现。

- [ ] Sprint 4.2 Part A-E / 27 Tasks 逐项复核
  - [ ] 对照 `docs/plans/2026-02-22-sprint04.2.md` 建立任务到代码映射。
  - [ ] 每个 Task 提供“代码位置 + 测试用例 + 页面联调截图/录屏”证据。

## Re-Review Gate

- [ ] 全部 P0 项完成后方可申请复审。
- [ ] 复审前执行：后端单测、前端单测、关键 E2E、权限回归、导出安全回归。
- [ ] 任一项未通过，继续保持阻塞状态。
