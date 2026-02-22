# Antigravity Sprint Prompt Design

**Goal:** 生成一个单一文件，包含 Sprint 1–16 的可执行 Prompt，供 Antigravity 在执行阶段使用。

**Scope:** 仅提示语设计，不包含实现代码。

---

## 1) 方案选择

采用 **方案 A：单文件 + 每 Sprint 一段完整 Prompt**。

原因：每个 Sprint Prompt 独立完整，减少模型遗漏约束的风险，最利于执行。

---

## 2) 文件结构

- 文件路径：`docs/plans/2026-02-22-antigravity-prompts.md`
- 结构：按 Sprint 1–16 顺序，每个 Sprint 一个完整 Prompt 段落。

---

## 3) Prompt 模板内容

每个 Sprint Prompt 包含以下部分：

1. **角色与目标**：明确 Antigravity 是执行工程师，目标是按该 Sprint 计划完成任务。
2. **必读文件**：必须先阅读对应 `docs/plans/2026-02-22-sprintXX.md`，必要时参考 `docs/plans/2026-02-22-master-roadmap.md`。
3. **范围与禁止**：仅允许修改当前 Sprint 涉及的文件；禁止跨 Sprint 改动；禁止跳过测试或跳过失败验证。
4. **TDD 规则**：先写失败测试→确认失败→最小实现→确认通过→必要重构；不得先写生产代码。
5. **测试与验证**：列出该 Sprint 关键测试命令与覆盖率要求（100%）。
6. **输出格式**：要求输出变更清单、测试结果、未完成项、下一步建议。
7. **中文要求**：除专有名词/代码/路径外，Prompt 全部中文。

---

## 4) 约束与最佳实践

- 统一要求 Antigravity 遵守仓库规范、技术栈与测试策略。
- 强制执行 RED/GREEN/REFACTOR。
- 所有测试命令必须真实执行并记录结果。
- 代码与文件路径必须与 Sprint 计划一致。
