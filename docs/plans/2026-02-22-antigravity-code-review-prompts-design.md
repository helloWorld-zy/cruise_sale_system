# Antigravity Code Review Prompts Design

## 概述
本设计文档旨在为 `cruise_sale_system` 项目的 16 个 Sprint 生成用于 Antigravity（智能 Review 代理）的 Code Review Prompts。审查语言为中文，且标准极为严格。

## 核心设计要点
1. **严苛的 Persona 设定**：设定 Antigravity 为极其严格的资深架构师（Senior Staff Engineer）、安全审计专家及资深 QA，以最高标准审视代码。
2. **强制校验 TDD 与 100% 覆盖率**：审查必须确认所有新功能均有对应的测试代码（RED-GREEN-REFACTOR），并达到 100% 测试覆盖。
3. **架构与规范强制要求**：
   - 后端：遵循 Go 1.26 最佳实践、领域驱动设计（DDD）分层架构、依赖注入及 RESTful 规范。
   - 前端：强制 Vue 3 `<script setup lang="ts">`，响应式与状态管理（Pinia）的最佳实践。
4. **针对各 Sprint 的专项风险审查**：例如，Sprint 2 重点关注超卖防范机制与 Meilisearch 索引事务；Sprint 4 重点关注支付回调幂等与事务安全。
5. **处理违背项的流程**：若发现任何代码违规或架构缺陷，**不直接 Reject，而是先生成详细的修改方案，待当前 Review 流程完成后，要求开发侧整体修复所有发现的问题，并在修改完成后发起新的一轮 Review**。

## 产出格式
所有的 Prompts（每个 Sprint 一个）将统一存放在单一的 Markdown 文件中（例如 `docs/antigravity-prompts.md`），结构如下：

```markdown
# Antigravity Code Review Prompts (Sprint 1 - 16)

## Sprint 1 Code Review Prompt
...

## Sprint 2 Code Review Prompt
...
```