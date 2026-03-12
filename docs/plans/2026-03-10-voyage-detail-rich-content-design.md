# Voyage Detail Rich Content Design

**Date:** 2026-03-10

## Goal

增强 miniapp 航次详情页，让它完整呈现 admin 已维护的逐日餐食、住宿、抵离港时间、费用说明和预订须知，并在 admin 中补齐全局模板库与航次挂载能力。

## Confirmed Product Decisions

- 费用说明模板和预订须知模板采用全局模板库。
- 航次使用模板时采用混合模式：默认引用模板，也支持复制为航次专属内容后独立编辑。
- 行程中的餐食和住宿直接复用现有 itinerary 字段，不新建平行结构。
- “行程安排”和“可选舱型”之间新增抵离港时间表区域。

## Current State

- backend VoyageItinerary 已包含 `eta_time`、`etd_time`、`has_breakfast`、`has_lunch`、`has_dinner`、`has_accommodation`、`accommodation_text`。
- admin 航次新建与编辑页已经可以维护上述字段。
- miniapp 航次详情页已经从公开 `GET /api/v1/voyages/:id` 拉取基础详情、行程、舱型和设施，但当前只展示城市与摘要，没有把餐食、住宿和时间表做成完整信息区。
- 费用说明与预订须知目前既没有全局模板库，也没有航次侧挂载能力。

## Recommended Approach

采用“全局模板库 + 航次快照覆盖”的方案。

原因如下：

- 全局模板库能最大化复用对公司和邮轮无关的规则性文本。
- 航次快照覆盖能避免模板被修改后意外影响已发布航次。

## Data Model

### 1. 航次行程展示数据

不新增新的逐日服务表，继续以 `voyage_itineraries` 为事实来源。

miniapp 详情页将每个 itinerary row 解释为：

- 港口与摘要
- 抵港时间 / 离港时间
- 餐食服务标签：早餐 / 午餐 / 晚餐
- 住宿标签与住宿说明

### 2. 全局模板实体

新增一类通用模板实体，覆盖两种业务类型：

- `fee_note`：费用说明模板
- `booking_notice`：预订须知模板

建议字段：

- `id`
- `name`
- `kind`
- `status`
- `content_json`
- `created_at`
- `updated_at`

其中 `content_json` 采用结构化 JSON，而不是 HTML 字符串。

### 3. 航次模板挂载与覆盖

在航次上新增两组字段：

- `fee_note_template_id`
- `fee_note_mode`
- `fee_note_content_json`
- `booking_notice_template_id`
- `booking_notice_mode`
- `booking_notice_content_json`

`mode` 取值：

- `template`：直接引用模板内容
- `snapshot`：从模板复制后转为航次专属内容

读取规则：

- 若 mode=`template`，优先取模板当前内容
- 若 mode=`snapshot`，取航次自身快照内容

## API Design

### Admin API

新增模板管理接口：

- `GET /api/v1/admin/content-templates?kind=fee_note|booking_notice`
- `GET /api/v1/admin/content-templates/:id`
- `POST /api/v1/admin/content-templates`
- `PUT /api/v1/admin/content-templates/:id`
- `DELETE /api/v1/admin/content-templates/:id`

扩展 admin 航次创建和更新接口，让 payload 支持：

- 模板选择
- mode 选择
- 航次快照内容

### Public API

扩展 `GET /api/v1/voyages/:id` 返回值，直接包含：

- 完整 itinerary rows，包括餐食、住宿、抵离港时间
- resolved `fee_note`
- resolved `booking_notice`

这样 miniapp 航次详情页不需要为费用说明和预订须知再额外请求模板接口。

## Admin Experience

### 1. 新增“文案模板”模块

admin 导航新增“文案模板”入口，页面内用 tabs 或 segmented control 切换：

- 费用说明模板
- 预订须知模板

每类模板都支持：

- 列表
- 新建
- 编辑
- 启停用

### 2. 模板编辑器

费用说明模板结构：

- `included`: 费用包含条目数组
- `excluded`: 费用不包含条目数组
- 每条支持 `text` 和可选 `emphasis`

预订须知模板结构：

- `sections`: 章节数组
- 每章包含 `key`、`title`、`items[]`
- item 支持普通正文和强调文案

建议固定首版章节键：

- `booking_limit`
- `change_refund`
- `documents`
- `special_groups`
- `payment`
- `travel_notes`

### 3. 航次编辑页挂载区

在 admin 航次编辑页新增两个结构卡片：

- 费用说明
- 预订须知

每个卡片提供：

- 模板下拉选择
- 当前模式展示（引用模板 / 航次专属）
- “复制为航次专属内容”按钮
- 当 mode=`snapshot` 时显示结构化编辑器

## Miniapp Experience

### 1. 行程安排区升级

将当前“按天城市摘要列表”升级为每日行程卡片：

- 顶部显示 Day、港口、抵离港时间
- 中部显示摘要
- 底部显示餐食胶囊标签和住宿标签
- 若有住宿说明，用柔和背景做单独提示块

### 2. 抵离港时间表区

插入在“行程安排”和“可选舱型”之间。

- 自动生成 `天数 / 港口 / 抵港 / 离港` 表格
- 时间缺失时显示 `--`

### 3. 费用说明区

放在“船上设施”下方，采用移动端友好的双分区样式：

- `费用包含`
- `费用不包含`

保留图片参考中的粗体小标题和 bullet 层级，但收敛为适合 miniapp 的间距和字号。

### 4. 预订须知区

放在费用说明下方，采用：

- 顶部横向章节导航
- 下方章节内容卡片
- 橙色强调警示块

这样保留长文本结构，又避免手机端一次性展开过长内容导致阅读疲劳。

## Error Handling And Fallbacks

- 航次如果没有选择模板也没有专属内容，对应区块显示“暂无说明”。
- itinerary 某些时间为空时，时间表和卡片均显示 `--`，不阻塞渲染。
- 若模板被停用但已有航次引用：
  - admin 航次编辑页仍显示原已选择模板并给出“已停用”标记
  - miniapp 公开详情照常返回可解析内容

## Testing Strategy

backend：

- 模板 CRUD 与 kind 过滤
- 航次创建/更新携带模板引用与快照内容
- 公开航次详情返回 resolved 费用说明和预订须知

admin：

- 模板列表/编辑页渲染和提交
- 航次页选择模板、复制为专属内容、保存 payload 正确

miniapp：

- 每日行程卡片展示餐食、住宿、抵离港时间
- 时间表正确渲染
- 费用说明与预订须知按结构化数据渲染

## Non-Goals

- 首版不做模板版本历史与回滚。
- 首版不支持富文本 HTML 任意粘贴。