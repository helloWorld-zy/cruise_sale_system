# Task23 覆盖率缺口清单与最小补测计划

## 1. 基线快照（当前）
- Web: Statements 86.96%, Branches 69.82%, Functions 72.94%, Lines 91.02%
- Admin: Statements 44.89%, Branches 32.21%, Functions 45.76%, Lines 46.95%
- Miniapp: Statements 83.38%, Branches 59.31%, Functions 82.27%, Lines 85.27%

说明: 本文档仅给出最小补测路径，不新增业务功能代码。

## 2. 覆盖率缺口 Top（按收益排序）

### 2.1 Admin（最高优先）
当前总量最低，且多个目录仍在 low 区间。

- `frontend/admin/coverage/app/pages/cabin-types/index.html`
  - Statements 27.41%, Branches 20.87%, Functions 29.41%, Lines 28.57%
- `frontend/admin/coverage/app/pages/facilities/index.html`
  - Statements 26.11%, Branches 18.04%, Functions 26.66%, Lines 29.29%
- `frontend/admin/coverage/app/pages/facility-categories/index.html`
  - Statements 27.88%, Branches 24.16%, Functions 17.39%, Lines 28.42%
- `frontend/admin/coverage/app/pages/cruises/index.html`
  - Statements 35.11%, Branches 23.17%, Functions 33.33%, Lines 36.55%
- `frontend/admin/coverage/app/pages/bookings/index.html`
  - Statements 42.10%, Branches 38.04%, Functions 41.66%, Lines 44.18%
- `frontend/admin/coverage/app/components/index.html`
  - `PricingRow.vue` 0% / `RouteTable.vue` 0%

### 2.2 Miniapp（中优先）
整体可用，但分支覆盖偏低，且组件存在明显短板。

- `frontend/miniapp/coverage/components/index.html`
  - `PayButton.vue`: Statements 22.22%, Branches 6.66%, Functions 0%, Lines 23.52%
- `frontend/miniapp/coverage/pages/login/index.html`
  - Statements 73.58%, Branches 65.00%, Functions 87.5%, Lines 76.0%
- `frontend/miniapp/coverage/pages/pay/index.html`
  - Statements 77.27%, Branches 61.11%, Functions 50.0%, Lines 76.19%
- `frontend/miniapp/coverage/pages/cruise/index.html`
  - Branches 53.95%（语句高，分支仍偏低）

### 2.3 Web（低优先，已明显提升）
当前总覆盖已较高，优先补分支薄弱点。

- `frontend/web/app/layouts/default.vue`
  - Statements 60%, Branches 20%, Functions 0%, Lines 55.55%
- `frontend/web/app/pages/cruises/[id].vue`
  - Statements 73.56%, Branches 48.76%, Functions 47.61%, Lines 81.08%
- `frontend/web/app/pages/cruises/index.vue`
  - Statements 88.23%, Branches 57.14%
- `frontend/web/app/pages/booking/success.vue`
  - Statements 83.33%, Branches 81.25%, Functions 80%, Lines 86.66%

## 3. 最小补测计划（按批次）

### 批次 A（先拉升 Admin 总盘）
目标: 用最少新增 spec 拉升最多 low 目录。

1. 新增 `frontend/admin/tests/unit/pages/cruises-id.spec.ts`
- 覆盖编辑页加载成功/失败、保存成功/失败、删除确认分支。

2. 新增 `frontend/admin/tests/unit/pages/cabin-types-form.spec.ts`
- 覆盖 `new.vue` 与 `[id].vue` 的公共分支: 标签勾选、设施勾选、提交 body 序列化、删除分支。

3. 新增 `frontend/admin/tests/unit/pages/facilities-form.spec.ts`
- 覆盖收费开关（extra_charge）显隐逻辑、target_audience 勾选、保存/删除分支。

4. 新增 `frontend/admin/tests/unit/pages/facility-categories-form.spec.ts`
- 覆盖内联编辑与状态切换提交分支。

5. 新增 `frontend/admin/tests/unit/components/pricing-row.spec.ts`
- 覆盖 `PricingRow.vue` 当前 0% 的基础渲染与事件分支。

6. 新增 `frontend/admin/tests/unit/components/route-table.spec.ts`
- 覆盖 `RouteTable.vue` 当前 0% 的渲染、空态、交互分支。

### 批次 B（Miniapp 精准补短）
1. 增强 `frontend/miniapp/tests/unit/pages/pay.spec.ts`
- 补齐成功、失败、取消、重复点击防抖等分支。

2. 新增 `frontend/miniapp/tests/unit/components/pay-button.spec.ts`
- 独立覆盖 `PayButton.vue` 的 disabled/loading/emit 分支（当前函数 0%）。

3. 增强 `frontend/miniapp/tests/unit/pages/login.spec.ts`
- 补异常分支与边界输入。

### 批次 C（Web 收尾）
1. 新增 `frontend/web/tests/unit/layouts/default.spec.ts` 分支用例
- 覆盖导航/条件渲染分支（当前 branches 20%）。

2. 增强 `frontend/web/tests/unit/pages/cruises/id.spec.ts`
- 补 API 异常、空数据、回退展示、操作按钮分支。

3. 增强 `frontend/web/tests/unit/pages/cruises/index.spec.ts`
- 补搜索筛选、空态、错误态分支。

## 4. 执行与验收
- 每批次执行:
  - Admin: `cd frontend/admin && npx vitest run --coverage`
  - Miniapp: `cd frontend/miniapp && npx vitest run --coverage`
  - Web: `cd frontend/web && npx vitest run --coverage`
- 验收原则:
  - 批次内新增 spec 全绿。
  - 对应目标目录覆盖率出现可观提升（尤其 Branches）。
  - 不引入业务逻辑改动，仅测试与必要测试桩调整。

## 5. 建议目标（现实可达）
- 短期（1-2 轮）:
  - Admin Lines 提升至 60%+
  - Miniapp Branches 提升至 70%+
  - Web Branches 提升至 75%+
- 中期:
  - 三端整体继续爬升，但 100% 覆盖率目标建议按“核心路径优先”分阶段达成。
