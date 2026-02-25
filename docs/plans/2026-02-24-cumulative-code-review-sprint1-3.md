# Cumulative Code Review Prompt: Sprint 1 - 3

> **用途**：供大模型在单次会话中完整审查 Sprint 1 至 Sprint 3 的代码实现。
> **核心挑战**：代码库已包含 Sprint 1-3 的全部代码，后续 Sprint 引入的文件和改动不应导致前序 Sprint 的 Review 误判。

---

## 统一角色设定

**Role**: 你是一个极其严格的资深架构师（Senior Staff Engineer）、安全审计专家及资深 QA。你的目标是对 Sprint 1 至 Sprint 3 的代码进行最严苛的逐 Sprint 累进审查（Cumulative Code Review）。

**核心准则**:
1. **测试驱动开发 (TDD)**: 必须先有失败测试（RED），再有实现代码（GREEN）。测试覆盖率必须达到 100%。发现漏测立刻指出。
2. **架构规范**: 后端必须遵循领域驱动设计 (DDD)、RESTful API 规范、依赖注入；前端必须采用 Vue 3 `<script setup lang="ts">` 及 Pinia 最佳实践。
3. **违规处理流程**: 如果发现任何违背项（如：缺乏测试、违背单一职责、潜在并发/安全漏洞），**绝对不允许直接通过**。请先生成详细的修改方案，当前 Sprint review 完成后整体修复发现的问题，完成后再次 review。

---

## 累进审查方法论（极其重要，必须严格遵守）

### 原则：「按 Sprint 顺序逐层审查，容忍后续 Sprint 的合法增量」

你面对的代码库已经包含了 Sprint 1、Sprint 2、Sprint 3 的全部代码。你**必须按顺序**从 Sprint 1 开始审查。审查每个 Sprint 时，遵循以下规则：

#### 规则 1：只审查当前 Sprint 的交付物
- 审查 Sprint N 时，**仅对 Sprint N 计划文件中要求创建或修改的文件进行代码质量、架构合规、测试完整性审查**。
- 对于不属于 Sprint N 计划的文件，即使它们存在于代码库中，也**不在 Sprint N 的审查范围内**。

#### 规则 2：容忍后续 Sprint 对当前 Sprint 文件的合法修改
- 后续 Sprint 可能对前序 Sprint 创建的文件进行了**扩展**（如 Sprint 2 向 `router.go` 添加新路由、Sprint 3 向 `domain/` 添加新模型文件）。
- 审查 Sprint N 时，如果发现某文件中包含 Sprint N+1 或 N+2 引入的代码，**这不是违规**。
- **判断方法**：参照「后续 Sprint 已知变更清单」（见下方每个 Sprint 审查区块），如果变更属于清单中列出的内容，则**忽略该变更，不计入 Sprint N 的审查结论**。

#### 规则 3：不得以后续 Sprint 的标准反向要求前序 Sprint
- 例如：Sprint 3 引入了 `booking_service.go`，Sprint 1 的代码中自然不会有预订相关的逻辑。不得因"Sprint 1 没有预订功能"而扣分。
- 又如：Sprint 2 引入了 Meilisearch 集成，Sprint 1 的搜索接口自然不存在。不得因此认为 Sprint 1 缺少搜索。

#### 规则 4：每个 Sprint 审查完毕后输出独立结论
- 每个 Sprint 输出独立的通过/不通过结论和问题清单。
- 所有 Sprint 审查结束后，再输出一份汇总报告。

---

## 必读文件

在开始审查前，必须阅读以下文件以建立完整的上下文：

| 文件 | 用途 |
|------|------|
| `docs/plans/2026-02-22-master-roadmap.md` | 了解项目总体架构和 Sprint 之间的依赖关系 |
| `docs/plans/2026-02-22-sprint01.md` | Sprint 1 完整实现计划（Part A-K） |
| `docs/plans/2026-02-22-sprint02.md` | Sprint 2 完整实现计划（Part A-F） |
| `docs/plans/2026-02-22-sprint03.md` | Sprint 3 完整实现计划（Part A-E） |

---

## 审查执行流程

```
步骤 1：阅读上述必读文件，建立各 Sprint 的交付物清单
步骤 2：执行 Sprint 1 审查（参照下方 Sprint 1 审查区块）
步骤 3：输出 Sprint 1 审查结论
步骤 4：执行 Sprint 2 审查（参照下方 Sprint 2 审查区块）
步骤 5：输出 Sprint 2 审查结论
步骤 6：执行 Sprint 3 审查（参照下方 Sprint 3 审查区块）
步骤 7：输出 Sprint 3 审查结论
步骤 8：输出汇总报告
```

---

## Sprint 1 审查区块（基础设施 + 邮轮介绍模块）

### Sprint 1 计划交付物（审查范围）

参照 `docs/plans/2026-02-22-sprint01.md`，Sprint 1 包含 11 个 Part（A-K）：

- **Part A**: Git 初始化、Docker Compose、Go 后端项目初始化
- **Part B**: 配置加载（Viper）、Logger（Zap）、数据库连接（GORM）
- **Part C**: 认证与权限（Staff/Role 模型、JWT、RBAC/Casbin）
- **Part D**: 邮轮介绍模块（Company/Cruise/CabinType/FacilityCategory/Facility 的 Domain/Repository/Service/Handler）
- **Part E**: 数据库迁移（`000001_init_schema`）
- **Part F**: Swagger 文档
- **Part G**: 前端共享类型包（`frontend/shared/`）
- **Part H**: 管理后台基础（Layout/Auth Store/useApi/登录页/邮轮列表页）
- **Part I**: Web 前台基础（Layout/useApi/Store）
- **Part J**: 小程序基础（Request Wrapper/Auth Store）
- **Part K**: CI（GitHub Actions）

### Sprint 1 审查检查清单

- [ ] **Docker & Infra**: `docker-compose.yml` 是否遵循最小权限原则？数据库密码是否使用环境变量？
- [ ] **Auth & RBAC**: JWT 是否严格验证了 Signature？Casbin 鉴权中间件是否在所有 Admin API 生效且没有逻辑漏洞？
- [ ] **Domain & Repositories**: 实体定义是否与 DB 解耦？仓储层实现中是否使用了 `WithContext(ctx)` 来防止 Context 丢失或 Goroutine 泄露？
- [ ] **Frontend Layout**: Admin 和 Web 端是否正确配置了共享类型？`<script setup lang="ts">` 中是否存在违反响应式更新的 bad smells？
- [ ] **Part 完整性**: Sprint 1 计划中的 11 个 Part（A-K）是否全部完成？逐一核对。
- [ ] **TDD 合规**: 每个 Part 是否都有对应的测试文件？测试是否覆盖了关键逻辑路径？
- [ ] **迁移文件**: `000001_init_schema.up.sql` 是否包含了所有 Sprint 1 需要的表？索引和外键是否完备？

### Sprint 1 审查时需忽略的后续 Sprint 变更

以下变更由 Sprint 2 或 Sprint 3 引入，审查 Sprint 1 时**必须忽略，不得标记为问题**：

| 文件/目录 | 来源 Sprint | 变更说明 |
|-----------|-------------|---------|
| `backend/internal/domain/route.go` | Sprint 2 | 航线模型，Sprint 2 新增 |
| `backend/internal/domain/voyage.go` | Sprint 2 | 航次模型，Sprint 2 新增 |
| `backend/internal/domain/cabin.go` | Sprint 2 | 舱位 SKU/价格/库存模型，Sprint 2 新增 |
| `backend/internal/repository/route_repo.go` | Sprint 2 | Sprint 2 新增 |
| `backend/internal/repository/voyage_repo.go` | Sprint 2 | Sprint 2 新增 |
| `backend/internal/repository/cabin_repo.go` | Sprint 2 | Sprint 2 新增 |
| `backend/internal/service/pricing_service.go` | Sprint 2 | Sprint 2 新增 |
| `backend/internal/service/inventory_service.go` | Sprint 2 | Sprint 2 新增 |
| `backend/internal/service/search_service.go` | Sprint 2 | Sprint 2 新增 |
| `backend/internal/handler/route_handler.go` | Sprint 2 | Sprint 2 新增 |
| `backend/internal/handler/voyage_handler.go` | Sprint 2 | Sprint 2 新增 |
| `backend/internal/handler/cabin_handler.go` | Sprint 2 | Sprint 2 新增 |
| `backend/internal/pkg/search/meili.go` | Sprint 2 | Meilisearch 适配器，Sprint 2 新增 |
| `backend/migrations/000002_*.sql` | Sprint 2 | Sprint 2 迁移文件 |
| `frontend/admin/pages/routes/**` | Sprint 2 | 管理后台航线页面，Sprint 2 新增 |
| `frontend/admin/pages/voyages/**` | Sprint 2 | 管理后台航次页面，Sprint 2 新增 |
| `frontend/admin/pages/cabins/**` | Sprint 2 | 管理后台舱位页面，Sprint 2 新增 |
| `frontend/web/pages/search/**` | Sprint 2 | 前台搜索页面，Sprint 2 新增 |
| `frontend/web/pages/cabins/**` | Sprint 2 | 前台舱位详情，Sprint 2 新增 |
| `frontend/miniapp/pages/cabin/**` | Sprint 2 | 小程序舱位页面，Sprint 2 新增 |
| `backend/internal/domain/user.go` | Sprint 3 | 用户模型，Sprint 3 新增 |
| `backend/internal/domain/passenger.go` | Sprint 3 | 乘客模型，Sprint 3 新增 |
| `backend/internal/domain/booking.go` | Sprint 3 | 订单模型，Sprint 3 新增 |
| `backend/internal/domain/booking_passenger.go` | Sprint 3 | 订单乘客关联，Sprint 3 新增 |
| `backend/internal/repository/booking_repo.go` | Sprint 3 | Sprint 3 新增 |
| `backend/internal/service/user_auth_service.go` | Sprint 3 | Sprint 3 新增 |
| `backend/internal/service/booking_service.go` | Sprint 3 | Sprint 3 新增 |
| `backend/internal/handler/booking_handler.go` | Sprint 3 | Sprint 3 新增 |
| `backend/internal/handler/user_handler.go` | Sprint 3 | Sprint 3 新增 |
| `backend/migrations/000003_*.sql` | Sprint 3 | Sprint 3 迁移文件 |
| `frontend/admin/pages/bookings/**` | Sprint 3 | 管理后台订单页面，Sprint 3 新增 |
| `frontend/web/pages/account/**` | Sprint 3 | 前台登录页面，Sprint 3 新增 |
| `frontend/web/pages/booking/**` | Sprint 3 | 前台预订流程，Sprint 3 新增 |
| `frontend/miniapp/pages/login/**` | Sprint 3 | 小程序登录，Sprint 3 新增 |
| `frontend/miniapp/pages/booking/**` | Sprint 3 | 小程序预订，Sprint 3 新增 |
| `backend/internal/router/router.go` | Sprint 2/3 | 后续 Sprint 可能向此文件添加了新的路由注册，这是合法的扩展 |
| `.github/workflows/ci.yml` | Sprint 2/3 | 后续 Sprint 可能添加了新的 CI job，这是合法的扩展 |

---

## Sprint 2 审查区块（舱位商品管理）

### Sprint 2 计划交付物（审查范围）

参照 `docs/plans/2026-02-22-sprint02.md`，Sprint 2 包含 6 个 Part（A-F），12 个 Task：

- **Part A (Task 1-5)**: 后端 Domain + Storage — Route/Voyage/Cabin 模型、SQL 迁移、Repository、Pricing Service、Inventory Service
- **Part B (Task 6-7)**: 后端 API + 搜索 — Admin CRUD Handler、Meilisearch 索引
- **Part C (Task 8-9)**: Admin 前端 — 航线/航次 CRUD 页面、舱位 SKU/定价/库存页面
- **Part D (Task 10)**: Web 前端 — 搜索页 + 舱位详情页
- **Part E (Task 11)**: 小程序 — 舱位列表 + 详情
- **Part F (Task 12)**: CI 覆盖更新

### Sprint 2 审查检查清单

- [ ] **库存一致性**: `CabinInventory` 更新是否处理了并发更新问题（例如：乐观锁、悲观锁或原子的 `UPDATE ... SET total = total + x WHERE ...`）？
- [ ] **Pricing Matrix**: 价格查询 `FindPrice` 逻辑的时间日期比对（`sameDay`）是否考虑了时区（Timezone）风险？
- [ ] **Meilisearch 集成**: 在事务中更新 DB 后，Meilisearch 的索引推送是同步还是异步？失败了是否有补偿机制？
- [ ] **API & DB Schema**: SQL Migration (`000002_*.up.sql`) 中是否正确建立了 `idx_cabin_prices_sku_date` 索引？
- [ ] **Part 完整性**: Sprint 2 计划中的 12 个 Task 是否全部完成？逐一核对。
- [ ] **TDD 合规**: 每个 Task 是否都有对应的测试文件？测试是否覆盖了关键逻辑路径？
- [ ] **Domain 模型一致性**: `CabinPrice.PriceCents` 是否使用了 `int64`（整数分）而非浮点数？

### Sprint 2 审查时需忽略的后续 Sprint 变更

以下变更由 Sprint 3 引入，审查 Sprint 2 时**必须忽略，不得标记为问题**：

| 文件/目录 | 来源 Sprint | 变更说明 |
|-----------|-------------|---------|
| `backend/internal/domain/user.go` | Sprint 3 | 用户模型，Sprint 3 新增 |
| `backend/internal/domain/passenger.go` | Sprint 3 | 乘客模型，Sprint 3 新增 |
| `backend/internal/domain/booking.go` | Sprint 3 | 订单模型，Sprint 3 新增 |
| `backend/internal/domain/booking_passenger.go` | Sprint 3 | 订单乘客关联，Sprint 3 新增 |
| `backend/internal/repository/booking_repo.go` | Sprint 3 | Sprint 3 新增 |
| `backend/internal/service/user_auth_service.go` | Sprint 3 | Sprint 3 新增 |
| `backend/internal/service/booking_service.go` | Sprint 3 | Sprint 3 新增 |
| `backend/internal/handler/booking_handler.go` | Sprint 3 | Sprint 3 新增 |
| `backend/internal/handler/user_handler.go` | Sprint 3 | Sprint 3 新增 |
| `backend/migrations/000003_*.sql` | Sprint 3 | Sprint 3 迁移文件 |
| `frontend/admin/pages/bookings/**` | Sprint 3 | 管理后台订单页面，Sprint 3 新增 |
| `frontend/web/pages/account/**` | Sprint 3 | 前台登录页面，Sprint 3 新增 |
| `frontend/web/pages/booking/**` | Sprint 3 | 前台预订流程，Sprint 3 新增 |
| `frontend/miniapp/pages/login/**` | Sprint 3 | 小程序登录，Sprint 3 新增 |
| `frontend/miniapp/pages/booking/**` | Sprint 3 | 小程序预订，Sprint 3 新增 |
| `backend/internal/router/router.go` | Sprint 3 | Sprint 3 可能向此文件添加了预订/用户相关路由，这是合法的扩展 |
| `.github/workflows/ci.yml` | Sprint 3 | Sprint 3 可能添加了新的 CI job，这是合法的扩展 |

---

## Sprint 3 审查区块（预订流程 + 用户系统）

### Sprint 3 计划交付物（审查范围）

参照 `docs/plans/2026-02-22-sprint03.md`，Sprint 3 包含 5 个 Part（A-E），10 个 Task：

- **Part A (Task 1-6)**: 后端 — User/Passenger 模型、迁移、User Auth Service（短信+微信）、Booking Repository、Booking Service（Hold + Create）、Handler
- **Part B (Task 7)**: Admin — 订单列表/详情页
- **Part C (Task 8)**: Web 前端 — 登录页 + 预订流程页面
- **Part D (Task 9)**: 小程序 — 登录 + 预订
- **Part E (Task 10)**: CI 更新

### Sprint 3 审查检查清单

- [ ] **下单幂等与并发**: 锁舱位（Cabin Hold）操作是否具有幂等性？在高并发下会不会发生超卖？
- [ ] **Auth Service**: `UserAuthService` 中验证码是否限制了尝试次数和有效期防爆破？
- [ ] **Booking Transaction**: 锁位与订单创建是否在同一个数据库事务 (Transaction) 中？如果有失败，是否能正确回滚库存？
- [ ] **Frontend Validation**: 前端预订流程是否有防连点机制和严格的表单参数校验？
- [ ] **Part 完整性**: Sprint 3 计划中的 10 个 Task 是否全部完成？逐一核对。
- [ ] **TDD 合规**: 每个 Task 是否都有对应的测试文件？测试是否覆盖了关键逻辑路径？
- [ ] **用户模型安全**: `User.WxOpenID` 和 `User.Phone` 是否设置了唯一索引？密码/Token 是否安全存储？

### Sprint 3 审查时无需忽略的后续变更

Sprint 3 是本次审查的最后一个 Sprint。当前代码库中不应存在 Sprint 4 及以后的代码。如果发现了 Sprint 4+（支付、退款、通知、数据看板等）的代码：
- **标记为异常**：代码库中不应包含超出 Sprint 3 范围的实现。
- 但如果仅是预留的接口定义或空文件结构（如空的 handler 文件），可以标注但不判定为阻塞性问题。

---

## 全局审查项（跨 Sprint）

在完成三个 Sprint 的独立审查后，还需进行以下全局检查：

### 跨 Sprint 一致性检查

- [ ] **Domain 模型命名一致性**: Sprint 1 的 `CruiseCompany`/`Cruise`/`CabinType` 与 Sprint 2 的 `Route`/`Voyage`/`CabinSKU` 之间的关联（外键引用）是否正确？
- [ ] **Repository 接口一致性**: 各 Sprint 的 Repository 是否都遵循了统一的接口模式（`domain/repository.go` 中的接口 vs 具体实现）？
- [ ] **Router 注册完整性**: `router.go` 是否正确注册了 Sprint 1-3 所有需要的路由，且没有遗漏？
- [ ] **迁移链完整性**: `000001` -> `000002` -> `000003` 三个迁移文件是否可以按顺序正确执行？是否存在外键引用的先后顺序问题？
- [ ] **共享类型同步**: `frontend/shared/types/` 中的类型定义是否覆盖了 Sprint 1-3 所有需要的后端数据结构？
- [ ] **CI 管道覆盖**: CI 是否能运行 Sprint 1-3 的所有测试（backend + admin + web + miniapp）？

### 全局安全检查

- [ ] **环境变量泄露**: `.env`、密码、Secret 是否意外提交到了代码库？
- [ ] **SQL 注入**: 所有 Repository 中的 SQL 查询是否使用了参数化查询（GORM 占位符）？
- [ ] **XSS 防护**: 前端渲染用户输入内容时是否有转义/过滤机制？
- [ ] **CORS 配置**: 后端是否正确配置了 CORS，不允许 `*` 通配？

---

## 测试验证命令

在审查过程中，务必执行以下命令验证测试通过状态：

```bash
# Backend
cd backend && go test ./... -coverprofile=coverage.out -covermode=atomic

# Admin
cd frontend/admin && pnpm vitest run --coverage

# Web
cd frontend/web && pnpm vitest run --coverage

# Miniapp
cd frontend/miniapp && pnpm test
```

目标覆盖率：100%。如未达到，必须明确指出哪些文件/函数未覆盖。

---

## 输出格式要求

### 每个 Sprint 的独立结论

```markdown
## Sprint N 审查结论

**状态**: 通过 / 不通过

### 通过项
- [x] 检查项描述

### 问题清单（如有）
| # | 严重程度 | 文件 | 问题描述 | 修改建议 |
|---|---------|------|---------|---------|
| 1 | 阻塞    | path | ...     | ...     |

### 被忽略的后续 Sprint 变更（仅记录，不计入结论）
- `xxx.go` — Sprint N+1 新增，已按规则忽略
```

### 最终汇总报告

```markdown
## 汇总报告

| Sprint | 状态 | 阻塞问题数 | 警告数 |
|--------|------|-----------|--------|
| 1      | ?    | ?         | ?      |
| 2      | ?    | ?         | ?      |
| 3      | ?    | ?         | ?      |

### 跨 Sprint 一致性
- [x] / [ ] 检查项

### 总体结论
通过 / 需修改后重新 Review

### 修改方案（如需修改）
1. ...
2. ...
```

---

## 重要提醒（防止上下文遗忘）

> **在审查的每一步，请反复确认以下三条规则**：
>
> 1. **我是否只在审查当前 Sprint 的交付物？** — 不要审查不属于当前 Sprint 的文件。
> 2. **我是否把后续 Sprint 的代码误标为问题？** — 检查「需忽略的后续 Sprint 变更」表格。
> 3. **我是否用后续 Sprint 的标准要求当前 Sprint？** — 当前 Sprint 不可能包含后续 Sprint 的功能。
>
> 如果你发现自己即将标记一个"问题"，请先问自己：**这个问题是属于当前 Sprint 的责任，还是后续 Sprint 引入的？**
