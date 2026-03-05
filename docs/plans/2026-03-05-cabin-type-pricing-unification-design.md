# 2026-03-05 舱型管理与价格管理重构设计

## 1. 背景与目标

当前后台存在独立的舱型管理与舱房管理，数据模型和交互路径存在重复与割裂。此次重构目标：

- 合并原有两模块能力，统一主菜单为`舱型管理`。
- 将四个大类（内舱房、海景房、阳台房、套房）作为可维护的分类字典（支持新增/改名/停用）。
- 实际售卖对象仍为具体舱型小类（例如标准内舱、豪华内舱）。
- 新增同级子路由`价格管理`，按`航次 + 具体舱型`管理库存、结算价、销售价。
- 价格支持历史版本管理，在线业务默认只读当前态；后续价格趋势读取历史序列。

## 2. 关键业务决策（已确认）

- 分类策略：预置四大类，后台可新增/改名/停用。
- 数据策略：强制重建，清空现有舱型相关数据，不保留旧数据。
- 菜单策略：新建统一菜单名`舱型管理`，原两个菜单下线。
- 舱型创建策略：多选邮轮时执行批量复制创建，每艘邮轮生成独立舱型记录。
- 价格粒度：`航次 + 具体舱型`，航次级固定价。
- 批量套价流程：先按公司/邮轮/出发区间筛航次，勾选航次后选舱型。
- 批量覆盖策略：覆盖当前价，同时保留旧版本为历史。
- 当前价判定：按`Asia/Shanghai`时区，取`effective_at <= now`最近一条。
- 媒体接口：使用专用舱型媒体接口，不复用通用上传。
- 媒体元数据第一版：类型（图片/平面图）、标题、排序、是否主图。
- 字段语义修正：去掉床型概念，仅保留`简介`字段。
- 在线读价例外：默认读当前态表，后续web/miniapp的价格走势组件可读取历史版本。

## 3. 信息架构与路由

- 主菜单：`舱型管理`
- 子路由：
  - `舱型管理/列表`
  - `舱型管理/新建`
  - `舱型管理/编辑`
  - `舱型管理/价格管理`（同级子路由）

## 4. 数据模型设计

### 4.1 舱型分类字典表 `cabin_type_categories`

字段建议：

- `id` BIGSERIAL PK
- `name` VARCHAR(64) NOT NULL
- `code` VARCHAR(32) NOT NULL UNIQUE
- `status` SMALLINT NOT NULL DEFAULT 1
- `sort_order` INT NOT NULL DEFAULT 0
- `created_at` TIMESTAMPTZ NOT NULL DEFAULT now()
- `updated_at` TIMESTAMPTZ NOT NULL DEFAULT now()
- `deleted_at` TIMESTAMPTZ NULL

初始化默认数据：内舱房、海景房、阳台房、套房。

### 4.2 实际舱型小类表 `cabin_types`

重定义为可售舱型小类，不再存储床型字段。

字段建议：

- `id` BIGSERIAL PK
- `category_id` BIGINT NOT NULL REFERENCES `cabin_type_categories(id)`
- `name` VARCHAR(100) NOT NULL
- `code` VARCHAR(50) NOT NULL
- `occupancy` INT NOT NULL
- `intro` TEXT NOT NULL DEFAULT ''
- `status` SMALLINT NOT NULL DEFAULT 1
- `sort_order` INT NOT NULL DEFAULT 0
- `created_at` TIMESTAMPTZ NOT NULL DEFAULT now()
- `updated_at` TIMESTAMPTZ NOT NULL DEFAULT now()
- `deleted_at` TIMESTAMPTZ NULL

建议唯一约束：`(code, deleted_at IS NULL)`或等价实现。

### 4.3 舱型-邮轮绑定表 `cabin_type_cruise_bindings`

字段建议：

- `id` BIGSERIAL PK
- `cabin_type_id` BIGINT NOT NULL REFERENCES `cabin_types(id)`
- `cruise_id` BIGINT NOT NULL REFERENCES `cruises(id)`
- `created_at` TIMESTAMPTZ NOT NULL DEFAULT now()

唯一约束：`(cabin_type_id, cruise_id)`。

### 4.4 舱型媒体表 `cabin_type_media`

字段建议：

- `id` BIGSERIAL PK
- `cabin_type_id` BIGINT NOT NULL REFERENCES `cabin_types(id)`
- `media_type` VARCHAR(20) NOT NULL -- image | floor_plan
- `url` TEXT NOT NULL
- `title` VARCHAR(120) NOT NULL DEFAULT ''
- `sort_order` INT NOT NULL DEFAULT 0
- `is_primary` BOOLEAN NOT NULL DEFAULT FALSE
- `created_at` TIMESTAMPTZ NOT NULL DEFAULT now()
- `updated_at` TIMESTAMPTZ NOT NULL DEFAULT now()
- `deleted_at` TIMESTAMPTZ NULL

约束建议：同`cabin_type_id + media_type`最多一个`is_primary=true`（部分唯一索引）。

### 4.5 航次舱型价格历史版本表 `voyage_cabin_type_price_versions`

字段建议：

- `id` BIGSERIAL PK
- `voyage_id` BIGINT NOT NULL REFERENCES `voyages(id)`
- `cabin_type_id` BIGINT NOT NULL REFERENCES `cabin_types(id)`
- `inventory_total` INT NOT NULL
- `settlement_price_cents` BIGINT NOT NULL
- `sale_price_cents` BIGINT NOT NULL
- `effective_at` TIMESTAMPTZ NOT NULL
- `created_by` BIGINT NULL
- `created_at` TIMESTAMPTZ NOT NULL DEFAULT now()

关键索引：`(voyage_id, cabin_type_id, effective_at DESC)`。

### 4.6 航次舱型当前态表 `voyage_cabin_type_current`

字段建议：

- `voyage_id` BIGINT NOT NULL REFERENCES `voyages(id)`
- `cabin_type_id` BIGINT NOT NULL REFERENCES `cabin_types(id)`
- `inventory_total` INT NOT NULL
- `settlement_price_cents` BIGINT NOT NULL
- `sale_price_cents` BIGINT NOT NULL
- `effective_at` TIMESTAMPTZ NOT NULL
- `version_id` BIGINT NOT NULL REFERENCES `voyage_cabin_type_price_versions(id)`
- `updated_at` TIMESTAMPTZ NOT NULL DEFAULT now()

主键或唯一键：`(voyage_id, cabin_type_id)`。

## 5. 强制重建迁移策略

单次迁移完成以下动作：

1. 清空旧舱型相关数据（按强制重建约定）。
2. 新建上述新表、索引、约束。
3. 写入四大类默认字典数据。
4. 下线旧舱型/舱房菜单依赖的数据读取路径。

说明：旧数据不做保留与映射。

## 6. 后端 API 设计

### 6.1 分类字典

- `GET /admin/cabin-type-categories`
- `POST /admin/cabin-type-categories`
- `PUT /admin/cabin-type-categories/:id`
- `PUT /admin/cabin-type-categories/:id/status`

### 6.2 舱型小类管理

- `GET /admin/cabin-types`（支持按分类、公司、邮轮、状态、关键词筛选）
- `POST /admin/cabin-types/batch-create`（多邮轮批量复制创建）
- `GET /admin/cabin-types/:id`
- `PUT /admin/cabin-types/:id`
- `DELETE /admin/cabin-types/:id`

### 6.3 舱型媒体专用接口

- `POST /admin/cabin-types/:id/media/upload`
- `GET /admin/cabin-types/:id/media`
- `PUT /admin/cabin-types/:id/media/:mediaId`
- `DELETE /admin/cabin-types/:id/media/:mediaId`

### 6.4 价格管理

- `GET /admin/price-management/voyages`（公司/邮轮/出发区间筛航次）
- `GET /admin/price-management/current`
- `POST /admin/price-management/batch-apply`（覆盖当前价 + 写历史版本）
- `GET /admin/price-management/history`（趋势数据源）

## 7. 当前价与历史价规则

- 当前价计算时区固定：`Asia/Shanghai`。
- 规则：`effective_at <= now`中最近的一条为当前有效版本。
- 批量覆盖会写入新版本，并刷新`voyage_cabin_type_current`。
- 在线业务默认读当前态；仅趋势类展示读历史版本。

## 8. 前端交互设计

### 8.1 舱型新建/编辑

- 所属邮轮：双列多选（左列公司，右列悬停公司后的邮轮列表）。
- 大类选择：从分类字典选择。
- 基础容量+最大容量合并为`容纳人数`。
- 删除床型字段，仅保留`简介`。
- 图片与平面图上传区接入专用媒体接口（非占位）。

### 8.2 价格管理子路由

- 第一步：按公司/邮轮/出发区间筛选航次。
- 第二步：勾选目标航次。
- 第三步：选择具体舱型。
- 第四步：设置库存、结算价、销售价、生效时间（北京时间）。
- 提交后展示成功/失败明细。

### 8.3 web/miniapp 消费逻辑

- 先按大类筛选，再展示具体舱型。
- 默认读取当前态价格与库存。
- 后续价格走势组件读取历史版本接口。

## 9. 错误处理与可观测性

- 批量接口返回结构化结果：`success_count`、`failed_count`、`failures[]`。
- 常见错误：分类停用冲突、主图冲突、生效时间非法、目标航次不可用。
- 关键操作写审计日志：批量创建舱型、批量套价。

## 10. 测试策略

- Migration 测试：新表、约束、默认字典、强制重建结果。
- 后端单元/集成：
  - 分类字典 CRUD 与停用
  - 舱型批量复制创建
  - 媒体接口与主图约束
  - 价格历史写入与当前态刷新
  - `Asia/Shanghai`生效判定
- 前端单测：
  - 双列多选创建舱型
  - 价格管理四步流程
  - 批量结果明细渲染
- E2E（建议）：
  - 舱型创建 -> 套价 -> web/miniapp读取当前价 -> 趋势读取历史

## 11. 非目标（本期不做）

- 自动迁移旧舱型数据映射。
- 多时区生效策略。
- 复杂版权/来源信息的媒体治理（可后续扩展）。
