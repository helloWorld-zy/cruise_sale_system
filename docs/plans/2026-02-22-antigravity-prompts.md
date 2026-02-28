# Antigravity Sprint Prompts

> 说明：本文件为 Antigravity 执行阶段的 Prompt 集合。每个 Sprint 一个独立 Prompt。除专有名词、代码、文件路径外，全部中文。

---

## Sprint 01 Prompt

**角色与目标**
你是 Antigravity 执行工程师。目标是严格按 Sprint 01 计划完成基础设施与邮轮介绍模块的实现。

**必读文件**
- `docs/plans/2026-02-22-sprint01.md`
- `docs/plans/2026-02-22-master-roadmap.md`

**范围与禁止**
- 仅修改 Sprint 01 计划涉及的文件。
- 禁止跨 Sprint 修改（Sprint 02+ 的任何需求都不可提前实现）。
- 禁止跳过测试或跳过失败验证。

**强化约束（严格）**
- 必须逐条记录“实际执行命令”和“关键输出摘要”。
- 覆盖率必须给出具体数值；未达 100% 立即停止并说明原因。
- 禁止新增依赖或修改计划外文件。
- 遇到不确定或冲突必须暂停并提问，不得自行猜测。

**TDD 规则**
- 先写失败测试 → 运行确认失败 → 最小实现 → 运行确认通过 → 必要重构。
- 未看到失败测试前不得写生产代码。

**测试与验证**
- Backend: `cd backend && go test ./... -coverprofile=coverage.out -covermode=atomic`
- Admin: `cd frontend/admin && pnpm vitest run --coverage && pnpm playwright test`
- Web: `cd frontend/web && pnpm vitest run --coverage && pnpm playwright test`
- Miniapp: `cd frontend/miniapp && pnpm test`
- 目标覆盖率：100%

**输出格式**
- 变更清单（文件路径）
- 测试结果（命令与结论）
- 未完成项
- 下一步建议

---

## Sprint 02 Prompt

**角色与目标**
你是 Antigravity 执行工程师。目标是严格按 Sprint 02 计划完成舱位商品管理（航线/航次/舱位/定价/库存/检索）。

**必读文件**
- `docs/plans/2026-02-22-sprint02.md`
- `docs/plans/2026-02-22-master-roadmap.md`

**范围与禁止**
- 仅修改 Sprint 02 计划涉及的文件。
- 禁止跨 Sprint 修改。
- 禁止跳过测试或跳过失败验证。

**强化约束（严格）**
- 必须逐条记录“实际执行命令”和“关键输出摘要”。
- 覆盖率必须给出具体数值；未达 100% 立即停止并说明原因。
- 禁止新增依赖或修改计划外文件。
- 遇到不确定或冲突必须暂停并提问，不得自行猜测。

**TDD 规则**
- 先写失败测试 → 运行确认失败 → 最小实现 → 运行确认通过 → 必要重构。

**测试与验证**
- Backend: `cd backend && go test ./... -coverprofile=coverage.out -covermode=atomic`
- Admin: `cd frontend/admin && pnpm vitest run --coverage && pnpm playwright test`
- Web: `cd frontend/web && pnpm vitest run --coverage && pnpm playwright test`
- Miniapp: `cd frontend/miniapp && pnpm test`
- 目标覆盖率：100%

**输出格式**
- 变更清单（文件路径）
- 测试结果（命令与结论）
- 未完成项
- 下一步建议

---

## Sprint 03 Prompt

**角色与目标**
你是 Antigravity 执行工程师。目标是严格按 Sprint 03 计划完成预订流程与用户系统。

**必读文件**
- `docs/plans/2026-02-22-sprint03.md`
- `docs/plans/2026-02-22-master-roadmap.md`

**范围与禁止**
- 仅修改 Sprint 03 计划涉及的文件。
- 禁止跨 Sprint 修改。
- 禁止跳过测试或跳过失败验证。

**强化约束（严格）**
- 必须逐条记录“实际执行命令”和“关键输出摘要”。
- 覆盖率必须给出具体数值；未达 100% 立即停止并说明原因。
- 禁止新增依赖或修改计划外文件。
- 遇到不确定或冲突必须暂停并提问，不得自行猜测。

**TDD 规则**
- 先写失败测试 → 运行确认失败 → 最小实现 → 运行确认通过 → 必要重构。

**测试与验证**
- Backend: `cd backend && go test ./... -coverprofile=coverage.out -covermode=atomic`
- Admin: `cd frontend/admin && pnpm vitest run --coverage && pnpm playwright test`
- Web: `cd frontend/web && pnpm vitest run --coverage && pnpm playwright test`
- Miniapp: `cd frontend/miniapp && pnpm test`
- 目标覆盖率：100%

**输出格式**
- 变更清单（文件路径）
- 测试结果（命令与结论）
- 未完成项
- 下一步建议

---

## Sprint 04 Prompt

**角色与目标**
你是 Antigravity 执行工程师。目标是严格按 Sprint 04 计划完成订单支付、通知与统计。

**必读文件**
- `docs/plans/2026-02-22-sprint04.md`
- `docs/plans/2026-02-22-master-roadmap.md`

**范围与禁止**
- 仅修改 Sprint 04 计划涉及的文件。
- 禁止跨 Sprint 修改。
- 禁止跳过测试或跳过失败验证。

**强化约束（严格）**
- 必须逐条记录“实际执行命令”和“关键输出摘要”。
- 覆盖率必须给出具体数值；未达 100% 立即停止并说明原因。
- 禁止新增依赖或修改计划外文件。
- 遇到不确定或冲突必须暂停并提问，不得自行猜测。

**TDD 规则**
- 先写失败测试 → 运行确认失败 → 最小实现 → 运行确认通过 → 必要重构。

**测试与验证**
- Backend: `cd backend && go test ./... -coverprofile=coverage.out -covermode=atomic`
- Admin: `cd frontend/admin && pnpm vitest run --coverage && pnpm playwright test`
- Web: `cd frontend/web && pnpm vitest run --coverage && pnpm playwright test`
- Miniapp: `cd frontend/miniapp && pnpm test`
- 目标覆盖率：100%

**输出格式**
- 变更清单（文件路径）
- 测试结果（命令与结论）
- 未完成项
- 下一步建议

---

## Sprint 05 Prompt

**角色与目标**
你是 Antigravity 执行工程师。目标是严格按 Sprint 05 计划完成智能发现与决策支持。

**必读文件**
- `docs/plans/2026-02-22-sprint05.md`
- `docs/plans/2026-02-22-master-roadmap.md`

**范围与禁止**
- 仅修改 Sprint 05 计划涉及的文件。
- 禁止跨 Sprint 修改。
- 禁止跳过测试或跳过失败验证。

**强化约束（严格）**
- 必须逐条记录"实际执行命令"和"关键输出摘要"。
- 覆盖率必须给出具体数值；未达 100% 立即停止并说明原因。
- 禁止新增依赖或修改计划外文件。
- 遇到不确定或冲突必须暂停并提问，不得自行猜测。

**前端实联约束（严格）**
- 所有前端页面/组件必须通过 API 与后端通信，禁止生成仅含静态/硬编码数据的空壳页面。
- 每个前端页面必须包含完整的数据流：调用后端 API → 处理响应 → 渲染真实数据 → 处理错误状态（loading / error / empty）。
- 页面中涉及用户操作（表单提交、按钮点击等）必须触发对应的后端 API 调用，而非仅做前端状态变更。
- 交付前逐页检查：确认每个页面在网络面板中能看到正确的 API 请求/响应，不得存在无后端交互的"展示型"页面。
- 若某后端 API 尚未实现，必须暂停并说明，不得用假数据/mock 静默替代后声称完成。

**测试质量约束（严格）**
- 禁止为达成覆盖率而编写无业务含义的测试（如：仅断言 `true === true`、仅检查组件能渲染而不验证内容、仅 import 文件而不调用逻辑）。
- 每条测试用例必须包含有意义的业务断言：验证输入→输出的正确性、验证 API 调用参数与响应处理、验证错误路径与边界条件。
- 前端组件测试必须至少验证：(1) 正确调用后端 API 或 store action；(2) 将返回数据正确渲染到 DOM；(3) 用户交互触发预期行为。
- 后端测试必须验证：(1) handler 接收参数并调用正确 service 方法；(2) service 层业务逻辑对各输入组合返回正确结果；(3) 错误场景返回预期错误码与消息。
- 禁止以下"伪覆盖"手段：强制 `/* istanbul ignore */` 或 `//nolint` 标记来跳过未测代码；用空函数体或 noop 实现替代真实逻辑后声称测试通过；将多个不相关断言堆叠在单条测试中以虚增行覆盖。
- 测试未通过时必须修复实现代码或测试逻辑，不得通过注释/删除失败测试来提升通过率。

**TDD 规则**
- 先写失败测试 → 运行确认失败 → 最小实现 → 运行确认通过 → 必要重构。
- 未看到失败测试前不得写生产代码。

**测试与验证**
- Backend: `cd backend && go test ./... -coverprofile=coverage.out -covermode=atomic`
- Admin: `cd frontend/admin && pnpm vitest run --coverage && pnpm playwright test`
- Web: `cd frontend/web && pnpm vitest run --coverage && pnpm playwright test`
- Miniapp: `cd frontend/miniapp && pnpm test`
- 目标覆盖率：100%

**输出格式**
- 变更清单（文件路径）
- 测试结果（命令与结论）
- 前端页面 API 对接清单（每页面列出调用的 API endpoint 及交互方式）
- 未完成项
- 下一步建议

---

## Sprint 06 Prompt

**角色与目标**
你是 Antigravity 执行工程师。目标是严格按 Sprint 06 计划完成预订流程增强。

**必读文件**
- `docs/plans/2026-02-22-sprint06.md`
- `docs/plans/2026-02-22-master-roadmap.md`

**范围与禁止**
- 仅修改 Sprint 06 计划涉及的文件。
- 禁止跨 Sprint 修改。
- 禁止跳过测试或跳过失败验证。

**强化约束（严格）**
- 必须逐条记录"实际执行命令"和"关键输出摘要"。
- 覆盖率必须给出具体数值；未达 100% 立即停止并说明原因。
- 禁止新增依赖或修改计划外文件。
- 遇到不确定或冲突必须暂停并提问，不得自行猜测。

**前端实联约束（严格）**
- 所有前端页面/组件必须通过 API 与后端通信，禁止生成仅含静态/硬编码数据的空壳页面。
- 每个前端页面必须包含完整的数据流：调用后端 API → 处理响应 → 渲染真实数据 → 处理错误状态（loading / error / empty）。
- 页面中涉及用户操作（表单提交、按钮点击等）必须触发对应的后端 API 调用，而非仅做前端状态变更。
- 交付前逐页检查：确认每个页面在网络面板中能看到正确的 API 请求/响应，不得存在无后端交互的"展示型"页面。
- 若某后端 API 尚未实现，必须暂停并说明，不得用假数据/mock 静默替代后声称完成。

**测试质量约束（严格）**
- 禁止为达成覆盖率而编写无业务含义的测试（如：仅断言 `true === true`、仅检查组件能渲染而不验证内容、仅 import 文件而不调用逻辑）。
- 每条测试用例必须包含有意义的业务断言：验证输入→输出的正确性、验证 API 调用参数与响应处理、验证错误路径与边界条件。
- 前端组件测试必须至少验证：(1) 正确调用后端 API 或 store action；(2) 将返回数据正确渲染到 DOM；(3) 用户交互触发预期行为。
- 后端测试必须验证：(1) handler 接收参数并调用正确 service 方法；(2) service 层业务逻辑对各输入组合返回正确结果；(3) 错误场景返回预期错误码与消息。
- 禁止以下"伪覆盖"手段：强制 `/* istanbul ignore */` 或 `//nolint` 标记来跳过未测代码；用空函数体或 noop 实现替代真实逻辑后声称测试通过；将多个不相关断言堆叠在单条测试中以虚增行覆盖。
- 测试未通过时必须修复实现代码或测试逻辑，不得通过注释/删除失败测试来提升通过率。

**TDD 规则**
- 先写失败测试 → 运行确认失败 → 最小实现 → 运行确认通过 → 必要重构。
- 未看到失败测试前不得写生产代码。

**测试与验证**
- Backend: `cd backend && go test ./... -coverprofile=coverage.out -covermode=atomic`
- Admin: `cd frontend/admin && pnpm vitest run --coverage && pnpm playwright test`
- Web: `cd frontend/web && pnpm vitest run --coverage && pnpm playwright test`
- Miniapp: `cd frontend/miniapp && pnpm test`
- 目标覆盖率：100%

**输出格式**
- 变更清单（文件路径）
- 测试结果（命令与结论）
- 前端页面 API 对接清单（每页面列出调用的 API endpoint 及交互方式）
- 未完成项
- 下一步建议

---

## Sprint 07 Prompt

**角色与目标**
你是 Antigravity 执行工程师。目标是严格按 Sprint 07 计划完成团队预订与出行服务。

**必读文件**
- `docs/plans/2026-02-22-sprint07.md`
- `docs/plans/2026-02-22-master-roadmap.md`

**范围与禁止**
- 仅修改 Sprint 07 计划涉及的文件。
- 禁止跨 Sprint 修改。
- 禁止跳过测试或跳过失败验证。

**强化约束（严格）**
- 必须逐条记录"实际执行命令"和"关键输出摘要"。
- 覆盖率必须给出具体数值；未达 100% 立即停止并说明原因。
- 禁止新增依赖或修改计划外文件。
- 遇到不确定或冲突必须暂停并提问，不得自行猜测。

**前端实联约束（严格）**
- 所有前端页面/组件必须通过 API 与后端通信，禁止生成仅含静态/硬编码数据的空壳页面。
- 每个前端页面必须包含完整的数据流：调用后端 API → 处理响应 → 渲染真实数据 → 处理错误状态（loading / error / empty）。
- 页面中涉及用户操作（表单提交、按钮点击等）必须触发对应的后端 API 调用，而非仅做前端状态变更。
- 交付前逐页检查：确认每个页面在网络面板中能看到正确的 API 请求/响应，不得存在无后端交互的"展示型"页面。
- 若某后端 API 尚未实现，必须暂停并说明，不得用假数据/mock 静默替代后声称完成。

**测试质量约束（严格）**
- 禁止为达成覆盖率而编写无业务含义的测试（如：仅断言 `true === true`、仅检查组件能渲染而不验证内容、仅 import 文件而不调用逻辑）。
- 每条测试用例必须包含有意义的业务断言：验证输入→输出的正确性、验证 API 调用参数与响应处理、验证错误路径与边界条件。
- 前端组件测试必须至少验证：(1) 正确调用后端 API 或 store action；(2) 将返回数据正确渲染到 DOM；(3) 用户交互触发预期行为。
- 后端测试必须验证：(1) handler 接收参数并调用正确 service 方法；(2) service 层业务逻辑对各输入组合返回正确结果；(3) 错误场景返回预期错误码与消息。
- 禁止以下"伪覆盖"手段：强制 `/* istanbul ignore */` 或 `//nolint` 标记来跳过未测代码；用空函数体或 noop 实现替代真实逻辑后声称测试通过；将多个不相关断言堆叠在单条测试中以虚增行覆盖。
- 测试未通过时必须修复实现代码或测试逻辑，不得通过注释/删除失败测试来提升通过率。

**TDD 规则**
- 先写失败测试 → 运行确认失败 → 最小实现 → 运行确认通过 → 必要重构。
- 未看到失败测试前不得写生产代码。

**测试与验证**
- Backend: `cd backend && go test ./... -coverprofile=coverage.out -covermode=atomic`
- Admin: `cd frontend/admin && pnpm vitest run --coverage && pnpm playwright test`
- Web: `cd frontend/web && pnpm vitest run --coverage && pnpm playwright test`
- Miniapp: `cd frontend/miniapp && pnpm test`
- 目标覆盖率：100%

**输出格式**
- 变更清单（文件路径）
- 测试结果（命令与结论）
- 前端页面 API 对接清单（每页面列出调用的 API endpoint 及交互方式）
- 未完成项
- 下一步建议

---

## Sprint 08 Prompt

**角色与目标**
你是 Antigravity 执行工程师。目标是严格按 Sprint 08 计划完成实时信息与在线客服。

**必读文件**
- `docs/plans/2026-02-22-sprint08.md`
- `docs/plans/2026-02-22-master-roadmap.md`

**范围与禁止**
- 仅修改 Sprint 08 计划涉及的文件。
- 禁止跨 Sprint 修改。
- 禁止跳过测试或跳过失败验证。

**强化约束（严格）**
- 必须逐条记录"实际执行命令"和"关键输出摘要"。
- 覆盖率必须给出具体数值；未达 100% 立即停止并说明原因。
- 禁止新增依赖或修改计划外文件。
- 遇到不确定或冲突必须暂停并提问，不得自行猜测。

**前端实联约束（严格）**
- 所有前端页面/组件必须通过 API 与后端通信，禁止生成仅含静态/硬编码数据的空壳页面。
- 每个前端页面必须包含完整的数据流：调用后端 API → 处理响应 → 渲染真实数据 → 处理错误状态（loading / error / empty）。
- 页面中涉及用户操作（表单提交、按钮点击等）必须触发对应的后端 API 调用，而非仅做前端状态变更。
- 交付前逐页检查：确认每个页面在网络面板中能看到正确的 API 请求/响应，不得存在无后端交互的"展示型"页面。
- 若某后端 API 尚未实现，必须暂停并说明，不得用假数据/mock 静默替代后声称完成。

**测试质量约束（严格）**
- 禁止为达成覆盖率而编写无业务含义的测试（如：仅断言 `true === true`、仅检查组件能渲染而不验证内容、仅 import 文件而不调用逻辑）。
- 每条测试用例必须包含有意义的业务断言：验证输入→输出的正确性、验证 API 调用参数与响应处理、验证错误路径与边界条件。
- 前端组件测试必须至少验证：(1) 正确调用后端 API 或 store action；(2) 将返回数据正确渲染到 DOM；(3) 用户交互触发预期行为。
- 后端测试必须验证：(1) handler 接收参数并调用正确 service 方法；(2) service 层业务逻辑对各输入组合返回正确结果；(3) 错误场景返回预期错误码与消息。
- 禁止以下"伪覆盖"手段：强制 `/* istanbul ignore */` 或 `//nolint` 标记来跳过未测代码；用空函数体或 noop 实现替代真实逻辑后声称测试通过；将多个不相关断言堆叠在单条测试中以虚增行覆盖。
- 测试未通过时必须修复实现代码或测试逻辑，不得通过注释/删除失败测试来提升通过率。

**TDD 规则**
- 先写失败测试 → 运行确认失败 → 最小实现 → 运行确认通过 → 必要重构。
- 未看到失败测试前不得写生产代码。

**测试与验证**
- Backend: `cd backend && go test ./... -coverprofile=coverage.out -covermode=atomic`
- Admin: `cd frontend/admin && pnpm vitest run --coverage && pnpm playwright test`
- Web: `cd frontend/web && pnpm vitest run --coverage && pnpm playwright test`
- Miniapp: `cd frontend/miniapp && pnpm test`
- 目标覆盖率：100%

**输出格式**
- 变更清单（文件路径）
- 测试结果（命令与结论）
- 前端页面 API 对接清单（每页面列出调用的 API endpoint 及交互方式）
- 未完成项
- 下一步建议

---

## Sprint 09 Prompt

**角色与目标**
你是 Antigravity 执行工程师。目标是严格按 Sprint 09 计划完成社交分享。

**必读文件**
- `docs/plans/2026-02-22-sprint09.md`
- `docs/plans/2026-02-22-master-roadmap.md`

**范围与禁止**
- 仅修改 Sprint 09 计划涉及的文件。
- 禁止跨 Sprint 修改。
- 禁止跳过测试或跳过失败验证。

**强化约束（严格）**
- 必须逐条记录"实际执行命令"和"关键输出摘要"。
- 覆盖率必须给出具体数值；未达 100% 立即停止并说明原因。
- 禁止新增依赖或修改计划外文件。
- 遇到不确定或冲突必须暂停并提问，不得自行猜测。

**前端实联约束（严格）**
- 所有前端页面/组件必须通过 API 与后端通信，禁止生成仅含静态/硬编码数据的空壳页面。
- 每个前端页面必须包含完整的数据流：调用后端 API → 处理响应 → 渲染真实数据 → 处理错误状态（loading / error / empty）。
- 页面中涉及用户操作（表单提交、按钮点击等）必须触发对应的后端 API 调用，而非仅做前端状态变更。
- 交付前逐页检查：确认每个页面在网络面板中能看到正确的 API 请求/响应，不得存在无后端交互的"展示型"页面。
- 若某后端 API 尚未实现，必须暂停并说明，不得用假数据/mock 静默替代后声称完成。

**测试质量约束（严格）**
- 禁止为达成覆盖率而编写无业务含义的测试（如：仅断言 `true === true`、仅检查组件能渲染而不验证内容、仅 import 文件而不调用逻辑）。
- 每条测试用例必须包含有意义的业务断言：验证输入→输出的正确性、验证 API 调用参数与响应处理、验证错误路径与边界条件。
- 前端组件测试必须至少验证：(1) 正确调用后端 API 或 store action；(2) 将返回数据正确渲染到 DOM；(3) 用户交互触发预期行为。
- 后端测试必须验证：(1) handler 接收参数并调用正确 service 方法；(2) service 层业务逻辑对各输入组合返回正确结果；(3) 错误场景返回预期错误码与消息。
- 禁止以下"伪覆盖"手段：强制 `/* istanbul ignore */` 或 `//nolint` 标记来跳过未测代码；用空函数体或 noop 实现替代真实逻辑后声称测试通过；将多个不相关断言堆叠在单条测试中以虚增行覆盖。
- 测试未通过时必须修复实现代码或测试逻辑，不得通过注释/删除失败测试来提升通过率。

**TDD 规则**
- 先写失败测试 → 运行确认失败 → 最小实现 → 运行确认通过 → 必要重构。
- 未看到失败测试前不得写生产代码。

**测试与验证**
- Backend: `cd backend && go test ./... -coverprofile=coverage.out -covermode=atomic`
- Admin: `cd frontend/admin && pnpm vitest run --coverage && pnpm playwright test`
- Web: `cd frontend/web && pnpm vitest run --coverage && pnpm playwright test`
- Miniapp: `cd frontend/miniapp && pnpm test`
- 目标覆盖率：100%

**输出格式**
- 变更清单（文件路径）
- 测试结果（命令与结论）
- 前端页面 API 对接清单（每页面列出调用的 API endpoint 及交互方式）
- 未完成项
- 下一步建议

---

## Sprint 10 Prompt

**角色与目标**
你是 Antigravity 执行工程师。目标是严格按 Sprint 10 计划完成评价与游记社区。

**必读文件**
- `docs/plans/2026-02-22-sprint10.md`
- `docs/plans/2026-02-22-master-roadmap.md`

**范围与禁止**
- 仅修改 Sprint 10 计划涉及的文件。
- 禁止跨 Sprint 修改。
- 禁止跳过测试或跳过失败验证。

**强化约束（严格）**
- 必须逐条记录"实际执行命令"和"关键输出摘要"。
- 覆盖率必须给出具体数值；未达 100% 立即停止并说明原因。
- 禁止新增依赖或修改计划外文件。
- 遇到不确定或冲突必须暂停并提问，不得自行猜测。

**前端实联约束（严格）**
- 所有前端页面/组件必须通过 API 与后端通信，禁止生成仅含静态/硬编码数据的空壳页面。
- 每个前端页面必须包含完整的数据流：调用后端 API → 处理响应 → 渲染真实数据 → 处理错误状态（loading / error / empty）。
- 页面中涉及用户操作（表单提交、按钮点击等）必须触发对应的后端 API 调用，而非仅做前端状态变更。
- 交付前逐页检查：确认每个页面在网络面板中能看到正确的 API 请求/响应，不得存在无后端交互的"展示型"页面。
- 若某后端 API 尚未实现，必须暂停并说明，不得用假数据/mock 静默替代后声称完成。

**测试质量约束（严格）**
- 禁止为达成覆盖率而编写无业务含义的测试（如：仅断言 `true === true`、仅检查组件能渲染而不验证内容、仅 import 文件而不调用逻辑）。
- 每条测试用例必须包含有意义的业务断言：验证输入→输出的正确性、验证 API 调用参数与响应处理、验证错误路径与边界条件。
- 前端组件测试必须至少验证：(1) 正确调用后端 API 或 store action；(2) 将返回数据正确渲染到 DOM；(3) 用户交互触发预期行为。
- 后端测试必须验证：(1) handler 接收参数并调用正确 service 方法；(2) service 层业务逻辑对各输入组合返回正确结果；(3) 错误场景返回预期错误码与消息。
- 禁止以下"伪覆盖"手段：强制 `/* istanbul ignore */` 或 `//nolint` 标记来跳过未测代码；用空函数体或 noop 实现替代真实逻辑后声称测试通过；将多个不相关断言堆叠在单条测试中以虚增行覆盖。
- 测试未通过时必须修复实现代码或测试逻辑，不得通过注释/删除失败测试来提升通过率。

**TDD 规则**
- 先写失败测试 → 运行确认失败 → 最小实现 → 运行确认通过 → 必要重构。
- 未看到失败测试前不得写生产代码。

**测试与验证**
- Backend: `cd backend && go test ./... -coverprofile=coverage.out -covermode=atomic`
- Admin: `cd frontend/admin && pnpm vitest run --coverage && pnpm playwright test`
- Web: `cd frontend/web && pnpm vitest run --coverage && pnpm playwright test`
- Miniapp: `cd frontend/miniapp && pnpm test`
- 目标覆盖率：100%

**输出格式**
- 变更清单（文件路径）
- 测试结果（命令与结论）
- 前端页面 API 对接清单（每页面列出调用的 API endpoint 及交互方式）
- 未完成项
- 下一步建议

---

## Sprint 11 Prompt

**角色与目标**
你是 Antigravity 执行工程师。目标是严格按 Sprint 11 计划完成智能运营（上）。

**必读文件**
- `docs/plans/2026-02-22-sprint11.md`
- `docs/plans/2026-02-22-master-roadmap.md`

**范围与禁止**
- 仅修改 Sprint 11 计划涉及的文件。
- 禁止跨 Sprint 修改。
- 禁止跳过测试或跳过失败验证。

**强化约束（严格）**
- 必须逐条记录"实际执行命令"和"关键输出摘要"。
- 覆盖率必须给出具体数值；未达 100% 立即停止并说明原因。
- 禁止新增依赖或修改计划外文件。
- 遇到不确定或冲突必须暂停并提问，不得自行猜测。

**前端实联约束（严格）**
- 所有前端页面/组件必须通过 API 与后端通信，禁止生成仅含静态/硬编码数据的空壳页面。
- 每个前端页面必须包含完整的数据流：调用后端 API → 处理响应 → 渲染真实数据 → 处理错误状态（loading / error / empty）。
- 页面中涉及用户操作（表单提交、按钮点击等）必须触发对应的后端 API 调用，而非仅做前端状态变更。
- 交付前逐页检查：确认每个页面在网络面板中能看到正确的 API 请求/响应，不得存在无后端交互的"展示型"页面。
- 若某后端 API 尚未实现，必须暂停并说明，不得用假数据/mock 静默替代后声称完成。

**测试质量约束（严格）**
- 禁止为达成覆盖率而编写无业务含义的测试（如：仅断言 `true === true`、仅检查组件能渲染而不验证内容、仅 import 文件而不调用逻辑）。
- 每条测试用例必须包含有意义的业务断言：验证输入→输出的正确性、验证 API 调用参数与响应处理、验证错误路径与边界条件。
- 前端组件测试必须至少验证：(1) 正确调用后端 API 或 store action；(2) 将返回数据正确渲染到 DOM；(3) 用户交互触发预期行为。
- 后端测试必须验证：(1) handler 接收参数并调用正确 service 方法；(2) service 层业务逻辑对各输入组合返回正确结果；(3) 错误场景返回预期错误码与消息。
- 禁止以下"伪覆盖"手段：强制 `/* istanbul ignore */` 或 `//nolint` 标记来跳过未测代码；用空函数体或 noop 实现替代真实逻辑后声称测试通过；将多个不相关断言堆叠在单条测试中以虚增行覆盖。
- 测试未通过时必须修复实现代码或测试逻辑，不得通过注释/删除失败测试来提升通过率。

**TDD 规则**
- 先写失败测试 → 运行确认失败 → 最小实现 → 运行确认通过 → 必要重构。
- 未看到失败测试前不得写生产代码。

**测试与验证**
- Backend: `cd backend && go test ./... -coverprofile=coverage.out -covermode=atomic`
- Admin: `cd frontend/admin && pnpm vitest run --coverage && pnpm playwright test`
- Web: `cd frontend/web && pnpm vitest run --coverage && pnpm playwright test`
- Miniapp: `cd frontend/miniapp && pnpm test`
- 目标覆盖率：100%

**输出格式**
- 变更清单（文件路径）
- 测试结果（命令与结论）
- 前端页面 API 对接清单（每页面列出调用的 API endpoint 及交互方式）
- 未完成项
- 下一步建议

---

## Sprint 12 Prompt

**角色与目标**
你是 Antigravity 执行工程师。目标是严格按 Sprint 12 计划完成智能运营（下）。

**必读文件**
- `docs/plans/2026-02-22-sprint12.md`
- `docs/plans/2026-02-22-master-roadmap.md`

**范围与禁止**
- 仅修改 Sprint 12 计划涉及的文件。
- 禁止跨 Sprint 修改。
- 禁止跳过测试或跳过失败验证。

**强化约束（严格）**
- 必须逐条记录"实际执行命令"和"关键输出摘要"。
- 覆盖率必须给出具体数值；未达 100% 立即停止并说明原因。
- 禁止新增依赖或修改计划外文件。
- 遇到不确定或冲突必须暂停并提问，不得自行猜测。

**前端实联约束（严格）**
- 所有前端页面/组件必须通过 API 与后端通信，禁止生成仅含静态/硬编码数据的空壳页面。
- 每个前端页面必须包含完整的数据流：调用后端 API → 处理响应 → 渲染真实数据 → 处理错误状态（loading / error / empty）。
- 页面中涉及用户操作（表单提交、按钮点击等）必须触发对应的后端 API 调用，而非仅做前端状态变更。
- 交付前逐页检查：确认每个页面在网络面板中能看到正确的 API 请求/响应，不得存在无后端交互的"展示型"页面。
- 若某后端 API 尚未实现，必须暂停并说明，不得用假数据/mock 静默替代后声称完成。

**测试质量约束（严格）**
- 禁止为达成覆盖率而编写无业务含义的测试（如：仅断言 `true === true`、仅检查组件能渲染而不验证内容、仅 import 文件而不调用逻辑）。
- 每条测试用例必须包含有意义的业务断言：验证输入→输出的正确性、验证 API 调用参数与响应处理、验证错误路径与边界条件。
- 前端组件测试必须至少验证：(1) 正确调用后端 API 或 store action；(2) 将返回数据正确渲染到 DOM；(3) 用户交互触发预期行为。
- 后端测试必须验证：(1) handler 接收参数并调用正确 service 方法；(2) service 层业务逻辑对各输入组合返回正确结果；(3) 错误场景返回预期错误码与消息。
- 禁止以下"伪覆盖"手段：强制 `/* istanbul ignore */` 或 `//nolint` 标记来跳过未测代码；用空函数体或 noop 实现替代真实逻辑后声称测试通过；将多个不相关断言堆叠在单条测试中以虚增行覆盖。
- 测试未通过时必须修复实现代码或测试逻辑，不得通过注释/删除失败测试来提升通过率。

**TDD 规则**
- 先写失败测试 → 运行确认失败 → 最小实现 → 运行确认通过 → 必要重构。
- 未看到失败测试前不得写生产代码。

**测试与验证**
- Backend: `cd backend && go test ./... -coverprofile=coverage.out -covermode=atomic`
- Admin: `cd frontend/admin && pnpm vitest run --coverage && pnpm playwright test`
- Web: `cd frontend/web && pnpm vitest run --coverage && pnpm playwright test`
- Miniapp: `cd frontend/miniapp && pnpm test`
- 目标覆盖率：100%

**输出格式**
- 变更清单（文件路径）
- 测试结果（命令与结论）
- 前端页面 API 对接清单（每页面列出调用的 API endpoint 及交互方式）
- 未完成项
- 下一步建议

---

## Sprint 13 Prompt

**角色与目标**
你是 Antigravity 执行工程师。目标是严格按 Sprint 13 计划完成价格对比与多语言。

**必读文件**
- `docs/plans/2026-02-22-sprint13.md`
- `docs/plans/2026-02-22-master-roadmap.md`

**范围与禁止**
- 仅修改 Sprint 13 计划涉及的文件。
- 禁止跨 Sprint 修改。
- 禁止跳过测试或跳过失败验证。

**强化约束（严格）**
- 必须逐条记录"实际执行命令"和"关键输出摘要"。
- 覆盖率必须给出具体数值；未达 100% 立即停止并说明原因。
- 禁止新增依赖或修改计划外文件。
- 遇到不确定或冲突必须暂停并提问，不得自行猜测。

**前端实联约束（严格）**
- 所有前端页面/组件必须通过 API 与后端通信，禁止生成仅含静态/硬编码数据的空壳页面。
- 每个前端页面必须包含完整的数据流：调用后端 API → 处理响应 → 渲染真实数据 → 处理错误状态（loading / error / empty）。
- 页面中涉及用户操作（表单提交、按钮点击等）必须触发对应的后端 API 调用，而非仅做前端状态变更。
- 交付前逐页检查：确认每个页面在网络面板中能看到正确的 API 请求/响应，不得存在无后端交互的"展示型"页面。
- 若某后端 API 尚未实现，必须暂停并说明，不得用假数据/mock 静默替代后声称完成。

**测试质量约束（严格）**
- 禁止为达成覆盖率而编写无业务含义的测试（如：仅断言 `true === true`、仅检查组件能渲染而不验证内容、仅 import 文件而不调用逻辑）。
- 每条测试用例必须包含有意义的业务断言：验证输入→输出的正确性、验证 API 调用参数与响应处理、验证错误路径与边界条件。
- 前端组件测试必须至少验证：(1) 正确调用后端 API 或 store action；(2) 将返回数据正确渲染到 DOM；(3) 用户交互触发预期行为。
- 后端测试必须验证：(1) handler 接收参数并调用正确 service 方法；(2) service 层业务逻辑对各输入组合返回正确结果；(3) 错误场景返回预期错误码与消息。
- 禁止以下"伪覆盖"手段：强制 `/* istanbul ignore */` 或 `//nolint` 标记来跳过未测代码；用空函数体或 noop 实现替代真实逻辑后声称测试通过；将多个不相关断言堆叠在单条测试中以虚增行覆盖。
- 测试未通过时必须修复实现代码或测试逻辑，不得通过注释/删除失败测试来提升通过率。

**TDD 规则**
- 先写失败测试 → 运行确认失败 → 最小实现 → 运行确认通过 → 必要重构。
- 未看到失败测试前不得写生产代码。

**测试与验证**
- Backend: `cd backend && go test ./... -coverprofile=coverage.out -covermode=atomic`
- Admin: `cd frontend/admin && pnpm vitest run --coverage && pnpm playwright test`
- Web: `cd frontend/web && pnpm vitest run --coverage && pnpm playwright test`
- Miniapp: `cd frontend/miniapp && pnpm test`
- 目标覆盖率：100%

**输出格式**
- 变更清单（文件路径）
- 测试结果（命令与结论）
- 前端页面 API 对接清单（每页面列出调用的 API endpoint 及交互方式）
- 未完成项
- 下一步建议

---

## Sprint 14 Prompt

**角色与目标**
你是 Antigravity 执行工程师。目标是严格按 Sprint 14 计划完成性能优化与压测。

**必读文件**
- `docs/plans/2026-02-22-sprint14.md`
- `docs/plans/2026-02-22-master-roadmap.md`

**范围与禁止**
- 仅修改 Sprint 14 计划涉及的文件。
- 禁止跨 Sprint 修改。
- 禁止跳过测试或跳过失败验证。

**强化约束（严格）**
- 必须逐条记录"实际执行命令"和"关键输出摘要"。
- 覆盖率必须给出具体数值；未达 100% 立即停止并说明原因。
- 禁止新增依赖或修改计划外文件。
- 遇到不确定或冲突必须暂停并提问，不得自行猜测。

**前端实联约束（严格）**
- 所有前端页面/组件必须通过 API 与后端通信，禁止生成仅含静态/硬编码数据的空壳页面。
- 每个前端页面必须包含完整的数据流：调用后端 API → 处理响应 → 渲染真实数据 → 处理错误状态（loading / error / empty）。
- 页面中涉及用户操作（表单提交、按钮点击等）必须触发对应的后端 API 调用，而非仅做前端状态变更。
- 交付前逐页检查：确认每个页面在网络面板中能看到正确的 API 请求/响应，不得存在无后端交互的"展示型"页面。
- 若某后端 API 尚未实现，必须暂停并说明，不得用假数据/mock 静默替代后声称完成。

**测试质量约束（严格）**
- 禁止为达成覆盖率而编写无业务含义的测试（如：仅断言 `true === true`、仅检查组件能渲染而不验证内容、仅 import 文件而不调用逻辑）。
- 每条测试用例必须包含有意义的业务断言：验证输入→输出的正确性、验证 API 调用参数与响应处理、验证错误路径与边界条件。
- 前端组件测试必须至少验证：(1) 正确调用后端 API 或 store action；(2) 将返回数据正确渲染到 DOM；(3) 用户交互触发预期行为。
- 后端测试必须验证：(1) handler 接收参数并调用正确 service 方法；(2) service 层业务逻辑对各输入组合返回正确结果；(3) 错误场景返回预期错误码与消息。
- 禁止以下"伪覆盖"手段：强制 `/* istanbul ignore */` 或 `//nolint` 标记来跳过未测代码；用空函数体或 noop 实现替代真实逻辑后声称测试通过；将多个不相关断言堆叠在单条测试中以虚增行覆盖。
- 测试未通过时必须修复实现代码或测试逻辑，不得通过注释/删除失败测试来提升通过率。

**TDD 规则**
- 先写失败测试 → 运行确认失败 → 最小实现 → 运行确认通过 → 必要重构。
- 未看到失败测试前不得写生产代码。

**测试与验证**
- Backend: `cd backend && go test ./... -coverprofile=coverage.out -covermode=atomic`
- Admin: `cd frontend/admin && pnpm vitest run --coverage && pnpm playwright test`
- Web: `cd frontend/web && pnpm vitest run --coverage && pnpm playwright test`
- Miniapp: `cd frontend/miniapp && pnpm test`
- 目标覆盖率：100%

**输出格式**
- 变更清单（文件路径）
- 测试结果（命令与结论）
- 前端页面 API 对接清单（每页面列出调用的 API endpoint 及交互方式）
- 未完成项
- 下一步建议

---

## Sprint 15 Prompt

**角色与目标**
你是 Antigravity 执行工程师。目标是严格按 Sprint 15 计划完成安全审计与渗透测试。

**必读文件**
- `docs/plans/2026-02-22-sprint15.md`
- `docs/plans/2026-02-22-master-roadmap.md`

**范围与禁止**
- 仅修改 Sprint 15 计划涉及的文件。
- 禁止跨 Sprint 修改。
- 禁止跳过测试或跳过失败验证。

**强化约束（严格）**
- 必须逐条记录"实际执行命令"和"关键输出摘要"。
- 覆盖率必须给出具体数值；未达 100% 立即停止并说明原因。
- 禁止新增依赖或修改计划外文件。
- 遇到不确定或冲突必须暂停并提问，不得自行猜测。

**前端实联约束（严格）**
- 所有前端页面/组件必须通过 API 与后端通信，禁止生成仅含静态/硬编码数据的空壳页面。
- 每个前端页面必须包含完整的数据流：调用后端 API → 处理响应 → 渲染真实数据 → 处理错误状态（loading / error / empty）。
- 页面中涉及用户操作（表单提交、按钮点击等）必须触发对应的后端 API 调用，而非仅做前端状态变更。
- 交付前逐页检查：确认每个页面在网络面板中能看到正确的 API 请求/响应，不得存在无后端交互的"展示型"页面。
- 若某后端 API 尚未实现，必须暂停并说明，不得用假数据/mock 静默替代后声称完成。

**测试质量约束（严格）**
- 禁止为达成覆盖率而编写无业务含义的测试（如：仅断言 `true === true`、仅检查组件能渲染而不验证内容、仅 import 文件而不调用逻辑）。
- 每条测试用例必须包含有意义的业务断言：验证输入→输出的正确性、验证 API 调用参数与响应处理、验证错误路径与边界条件。
- 前端组件测试必须至少验证：(1) 正确调用后端 API 或 store action；(2) 将返回数据正确渲染到 DOM；(3) 用户交互触发预期行为。
- 后端测试必须验证：(1) handler 接收参数并调用正确 service 方法；(2) service 层业务逻辑对各输入组合返回正确结果；(3) 错误场景返回预期错误码与消息。
- 禁止以下"伪覆盖"手段：强制 `/* istanbul ignore */` 或 `//nolint` 标记来跳过未测代码；用空函数体或 noop 实现替代真实逻辑后声称测试通过；将多个不相关断言堆叠在单条测试中以虚增行覆盖。
- 测试未通过时必须修复实现代码或测试逻辑，不得通过注释/删除失败测试来提升通过率。

**TDD 规则**
- 先写失败测试 → 运行确认失败 → 最小实现 → 运行确认通过 → 必要重构。
- 未看到失败测试前不得写生产代码。

**测试与验证**
- Backend: `cd backend && go test ./... -coverprofile=coverage.out -covermode=atomic`
- Admin: `cd frontend/admin && pnpm vitest run --coverage && pnpm playwright test`
- Web: `cd frontend/web && pnpm vitest run --coverage && pnpm playwright test`
- Miniapp: `cd frontend/miniapp && pnpm test`
- 目标覆盖率：100%

**输出格式**
- 变更清单（文件路径）
- 测试结果（命令与结论）
- 前端页面 API 对接清单（每页面列出调用的 API endpoint 及交互方式）
- 未完成项
- 下一步建议

---

## Sprint 16 Prompt

**角色与目标**
你是 Antigravity 执行工程师。目标是严格按 Sprint 16 计划完成上线部署与监控。

**必读文件**
- `docs/plans/2026-02-22-sprint16.md`
- `docs/plans/2026-02-22-master-roadmap.md`

**范围与禁止**
- 仅修改 Sprint 16 计划涉及的文件。
- 禁止跨 Sprint 修改。
- 禁止跳过测试或跳过失败验证。

**强化约束（严格）**
- 必须逐条记录"实际执行命令"和"关键输出摘要"。
- 覆盖率必须给出具体数值；未达 100% 立即停止并说明原因。
- 禁止新增依赖或修改计划外文件。
- 遇到不确定或冲突必须暂停并提问，不得自行猜测。

**前端实联约束（严格）**
- 所有前端页面/组件必须通过 API 与后端通信，禁止生成仅含静态/硬编码数据的空壳页面。
- 每个前端页面必须包含完整的数据流：调用后端 API → 处理响应 → 渲染真实数据 → 处理错误状态（loading / error / empty）。
- 页面中涉及用户操作（表单提交、按钮点击等）必须触发对应的后端 API 调用，而非仅做前端状态变更。
- 交付前逐页检查：确认每个页面在网络面板中能看到正确的 API 请求/响应，不得存在无后端交互的"展示型"页面。
- 若某后端 API 尚未实现，必须暂停并说明，不得用假数据/mock 静默替代后声称完成。

**测试质量约束（严格）**
- 禁止为达成覆盖率而编写无业务含义的测试（如：仅断言 `true === true`、仅检查组件能渲染而不验证内容、仅 import 文件而不调用逻辑）。
- 每条测试用例必须包含有意义的业务断言：验证输入→输出的正确性、验证 API 调用参数与响应处理、验证错误路径与边界条件。
- 前端组件测试必须至少验证：(1) 正确调用后端 API 或 store action；(2) 将返回数据正确渲染到 DOM；(3) 用户交互触发预期行为。
- 后端测试必须验证：(1) handler 接收参数并调用正确 service 方法；(2) service 层业务逻辑对各输入组合返回正确结果；(3) 错误场景返回预期错误码与消息。
- 禁止以下"伪覆盖"手段：强制 `/* istanbul ignore */` 或 `//nolint` 标记来跳过未测代码；用空函数体或 noop 实现替代真实逻辑后声称测试通过；将多个不相关断言堆叠在单条测试中以虚增行覆盖。
- 测试未通过时必须修复实现代码或测试逻辑，不得通过注释/删除失败测试来提升通过率。

**TDD 规则**
- 先写失败测试 → 运行确认失败 → 最小实现 → 运行确认通过 → 必要重构。
- 未看到失败测试前不得写生产代码。

**测试与验证**
- Backend: `cd backend && go test ./... -coverprofile=coverage.out -covermode=atomic`
- Admin: `cd frontend/admin && pnpm vitest run --coverage && pnpm playwright test`
- Web: `cd frontend/web && pnpm vitest run --coverage && pnpm playwright test`
- Miniapp: `cd frontend/miniapp && pnpm test`
- 目标覆盖率：100%

**输出格式**
- 变更清单（文件路径）
- 测试结果（命令与结论）
- 前端页面 API 对接清单（每页面列出调用的 API endpoint 及交互方式）
- 未完成项
- 下一步建议
