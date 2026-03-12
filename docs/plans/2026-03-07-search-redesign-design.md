# NCL-Style UI & Search Redesign Design Doc

## 1. 目标 (Goal)
修复 Web 前端 `search/index.vue` 的样式以完全匹配 `www.ncl.com`（NCL 领航风格），同时在不修改后端和 Admin 代码的前提下实现高级搜索栏（支持**目的地(Where)**、**出发地(Departure)**、**日期(Dates)** 的联动及条件搜索）。

---

## 2. 界面设计 (UI Design)

按照 NCL 提供的参考图，页面将被重构为以下层次：
1. **Header (导航栏)**: 绝对定位/浮动在 Hero Banner 顶部（透明背景，滚动后白底），左侧为 Logotype，中间是导航菜单（Explore, Ships, Deals 等），最右侧为 User/Search 图标。
2. **Hero Banner**: 全屏宽度背景大图（使用高品质邮轮背景），叠加居中的营销推广文案（例如 "THIS WEEKEND ONLY! FREE + FREE..."）和倒计时组件（可选）。
3. **Floating Search Bar (悬浮搜索栏)**: 叠加在 Hero Banner 下方中间。包含三个主要的选择框：
   - **Where (目的地 / Arrival Port)**
   - **Departure (出发地 / Departure Port)**
   - **Dates (日期 / Depart Date)**
   - 右侧是显眼的粉/品红色大按钮 **Find Cruises**。
4. **Cruise Cards**: 下方以水平或网格方式平铺推荐的各条航线/航次的信息卡片（带有图像、价格、标签如“Last Minute Alaska Cruises”等）。

---

## 3. 搜索架构设计 (Search Architecture Approaches)

由于我们**不能修改任何后端或 Admin 代码**，后端的 `/cabins` 或 `/voyages` 可能并不直接支持类似 `?arrival_port=xxx&depart_date=xxx` 这样复杂的复合字段查询。我们必须利用现有 API 完成高级搜索。

### Approach A: 纯前端数据聚合 + 全量过滤 (推荐 / Recommended)
**工作流**：
1. 页面加载时，并行请求所有启用的航线 `/routes`（获取 `DeparturePort`, `ArrivalPort/Name`）以及所有的航次 `/voyages`（获取 `DepartDate`, `ReturnDate`, 关联的 `RouteID` / `CruiseID`）。
2. 在前端搜索栏，提取所有唯一的“目的地(Where)”、“出发地(Departure)”和“年月(Dates)”作为下拉框的选项。
3. 用户在表单选择条件后，点击 **Find Cruises**，前端逻辑对组合数据进行 `filter`，筛选出符合条件的 `voyage`（航次）的 ID 集合。
4. 将匹配的这些 `Voyage ID`（或 CruiseID）传入 `/cabin-types` 的迭代请求或显示所有符合条件的航次（而非直接显示舱位，NCL 风格一般先展示 Cruises，点击后再看 Cabins）。

**优势**：UI 响应极快，下拉框的联动状态（如选了目的地之后，自动过滤可用的出发地）容易实现；不需要任何新的后端接口。
**劣势**：如果数据量上千，首屏会一次性下载较大 JSON，但对于一般的初创邮轮系统，前期的 Route 和 Voyage 数量前端完全 hold 得住。

### Approach B: 基于现存 keyword/航线类型 的懒加载搜索
**工作流**：
利用现有的 `search/index.vue` 的 `keyword` 和下拉选项，通过把目的地/出发地/日期的选择组合为一个模糊的 `keyword` 字符串，或者分别匹配相关的下发拉框 `CruiseID` 传给后端。如果后端的 query 接口支持部分筛选，可以直接用 query 请求。
**优势**：不依赖全量加载。
**劣势**：无法准确实现精确下拉框（不知道系统里有哪些不重复的出发地和目的地可供选择）。

---

## 4. 请您确认 / Please Confirm

为了继续实施，请告诉我：
1. **是否同意使用 Approach A（前端进行所有航线和航次数据的加载与级联筛选，再匹配舱位/航次展示）以符合“不改后端代码”的限制？**
2. 搜索出的结果列表展示，是**先展示一个个航次 (Voyage / Cruise 卡片)**，点击进去后再展示舱型和预订，还是**直接平铺所有符合条件航次的舱型 (Cabin Types)**？（NCL 原版搜索出来的是 Cruises）。
