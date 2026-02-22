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
