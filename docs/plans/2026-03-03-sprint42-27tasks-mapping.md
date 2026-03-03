# Sprint 4.2 27 Tasks Mapping Matrix

> Source: `docs/plans/2026-02-22-sprint04.2.md`
> Generated: 2026-03-03
> Legend: `complete` = code+test evidence available; `partial` = partially matched or evidence incomplete; `gap` = not implemented or missing key evidence.

| Task | Requirement | Status | Code Evidence | Test Evidence | UI Evidence |
|---|---|---|---|---|---|
| 1 | Route 港口/停靠字段扩展 | complete | `backend/internal/domain/route.go` | `backend/internal/domain/route_test.go` | N/A (backend) |
| 2 | CabinSKU 完整属性 | complete | `backend/internal/domain/cabin.go` | `backend/internal/domain/cabin_test.go` | N/A (backend) |
| 3 | CabinPrice 多类型定价字段 | complete | `backend/internal/domain/cabin.go` | `backend/internal/domain/cabin_test.go` | N/A (backend) |
| 4 | CabinInventory 预警阈值 + 可用量预警判定 | complete | `backend/internal/domain/cabin.go`, `backend/internal/service/inventory_alert_service.go` | `backend/internal/service/inventory_alert_service_test.go` | N/A (backend) |
| 5 | PricingService 多价格类型查询 + 日期区间批量设价 | complete | `backend/internal/service/pricing_service.go`, `backend/internal/service/cabin_admin_service.go` | `backend/internal/service/pricing_service_test.go`, `backend/internal/handler/cabin_handler_extended_test.go` | N/A (backend) |
| 6 | 舱位批量上/下架 + 三级分类 API | complete | `backend/internal/handler/cabin_handler.go`, `backend/internal/router/router.go` | `backend/internal/handler/cabin_handler_extended_test.go` | N/A (backend) |
| 7 | 舱位扩展字段 migration | complete | `backend/migrations/000009_sprint42_cabin_extend.up.sql`, `backend/migrations/000009_sprint42_cabin_extend.down.sql` | `backend/migrations/migrations_sprint42_test.go` | N/A (backend) |
| 8 | Admin 舱位管理页增强 | complete | `frontend/admin/app/pages/cabins/index.vue`, `frontend/admin/app/pages/cabins/new.vue`, `frontend/admin/app/pages/cabins/[id].vue` | `frontend/admin/tests/unit/pages/cabins-index.spec.ts`, `frontend/admin/tests/unit/pages/cabins-new.spec.ts`, `frontend/admin/tests/unit/pages/cabins-id.spec.ts` | `docs/plans/evidence/task20/task8-cabins.png`, `docs/plans/evidence/task20/task8-cruises.png` |
| 9 | Admin 定价日历 + 批量设价页面 | complete | `frontend/admin/app/pages/cabins/pricing.vue` | `frontend/admin/tests/unit/pages/cabins-pricing.spec.ts` | `docs/plans/evidence/task20/task9-pricing.png` |
| 10 | 前台舱位浏览/详情增强 (web+miniapp) | complete | `frontend/web/app/pages/cabins/index.vue`, `frontend/miniapp/pages/cabin/list.vue`, `frontend/miniapp/components/CabinCard.vue`, `frontend/shared/components/InventoryBadge.vue` | `frontend/web/tests/unit/pages/cabins/index.spec.ts`, `frontend/miniapp/tests/cabin-list.spec.ts`, `frontend/miniapp/tests/cabin-detail.spec.ts` | `docs/plans/evidence/task20/task8-cabins.png` |
| 11 | Passenger 完整字段 | complete | `backend/internal/domain/passenger.go` | `backend/internal/domain/passenger_test.go` | N/A (backend) |
| 12 | User 扩展 (支付宝/邮箱) | complete | `backend/internal/domain/user.go` | `backend/internal/domain/user_test.go` | N/A (backend) |
| 13 | UserAuth 支付宝登录 + 账号绑定 | complete | `backend/internal/service/user_auth_service.go` | `backend/internal/service/user_auth_service_test.go` | N/A (backend) |
| 14 | 常用乘客服务 | complete | `backend/internal/service/passenger_service.go` | `backend/internal/service/passenger_service_test.go` | N/A (backend) |
| 15 | 前台我的订单页面 (web+miniapp) | complete | `frontend/web/app/pages/orders/index.vue`, `frontend/miniapp/pages/orders/list.vue`, `frontend/shared/types/order.ts` | `frontend/web/tests/unit/pages/orders/index.spec.ts`, `frontend/miniapp/tests/orders-list.spec.ts` | N/A (frontend unit evidence) |
| 16 | 用户/乘客扩展 migration | complete | `backend/migrations/000010_sprint42_user_passenger_extend.up.sql`, `backend/migrations/000010_sprint42_user_passenger_extend.down.sql` | `backend/migrations/migrations_sprint42_test.go` | N/A (backend) |
| 17 | 订单状态生命周期 + 状态日志 | complete | `backend/internal/repository/booking_repo.go`, `backend/internal/handler/booking_handler.go` | `backend/internal/repository/booking_repo_test.go` | N/A (backend) |
| 18 | 支付金额校验 + 超时关单 | complete | `backend/internal/service/payment_service.go`, `backend/internal/service/order_timeout_service.go` | `backend/internal/service/payment_service_test.go`, `backend/internal/service/order_timeout_service_test.go` | N/A (backend) |
| 19 | 阶梯退款规则 | complete | `backend/internal/service/refund_service.go` | `backend/internal/service/refund_service_tiered_test.go` | N/A (backend) |
| 20 | 财务对账报表 | complete | `backend/internal/service/reconciliation_service.go` | `backend/internal/service/reconciliation_service_test.go` | N/A (backend) |
| 21 | 后台订单列表增强 + Excel 导出 | complete | `frontend/admin/app/pages/bookings/index.vue`, `backend/internal/service/order_export_service.go`, `backend/internal/handler/booking_handler.go` | `frontend/admin/tests/unit/bookings.list.spec.ts`, `backend/internal/service/order_export_service_test.go` | `docs/plans/evidence/task20/task21-bookings-export.png` |
| 22 | 员工账号管理 + 角色分配 | complete | `backend/internal/handler/staff_handler.go`, `backend/internal/service/staff_service.go`, `backend/internal/router/router.go` | `backend/internal/service/staff_service_test.go` | `docs/plans/evidence/task20/task22-staff.png` |
| 23 | 店铺/品牌信息管理 | complete | `backend/internal/handler/shop_info_handler.go`, `backend/internal/service/shop_info_service.go` | `backend/internal/service/shop_info_service_test.go`, `frontend/admin/tests/unit/pages/settings.shop.spec.ts` | `docs/plans/evidence/task20/task23-shop.png` |
| 24 | 通知模板配置 + 多渠道通知 | complete | `backend/internal/handler/notification_template_handler.go`, `backend/internal/service/notify_service.go`, `frontend/admin/app/pages/notifications/templates.vue` | `backend/internal/service/notify_template_test.go`, `backend/internal/service/notify_service_test.go`, `frontend/admin/tests/unit/pages/notifications.templates.spec.ts` | `docs/plans/evidence/task20/task24-templates.png` |
| 25 | 库存预警通知 | complete | `backend/internal/service/inventory_alert_service.go` | `backend/internal/service/inventory_alert_service_test.go` | N/A (backend) |
| 26 | 数据看板增强 | complete | `backend/internal/repository/analytics_repo.go`, `backend/internal/service/analytics_service.go` | `backend/internal/repository/analytics_repo_test.go` | `docs/plans/evidence/task20/task26-dashboard.png`, `docs/plans/evidence/task22/admin-dashboard.png` |
| 27 | Sprint 4.2 扩展 migration 汇总 | complete | `backend/migrations/000009_sprint42_cabin_extend.*.sql`, `backend/migrations/000010_sprint42_user_passenger_extend.*.sql`, `backend/migrations/000011_sprint42_order_notify_extend.*.sql`, `backend/migrations/000012_analytics_indexes.*.sql`, `backend/migrations/000013_shop_info_singleton.*.sql` | `backend/migrations/migrations_sprint42_test.go` | N/A (backend) |

## Summary

- complete: 27
- partial: 0
- gap: 0

## Gap List (for next tasks)

当前无 `gap` 或 `partial` 项；27 个任务均已具备代码、测试与（适用时）UI 证据。
