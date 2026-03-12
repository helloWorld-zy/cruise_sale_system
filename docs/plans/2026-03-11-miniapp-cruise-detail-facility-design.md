# 2026-03-11 Miniapp Cruise Detail Facility Redesign Design

## 1. Background And Goal
- Target: Redesign the miniapp cruise detail page reached from `邮轮百科` cruise cards so it matches the provided small-program reference more closely.
- Primary UX: Show a top image carousel, a ship identity block, a dynamic basic-parameter section, a description section, and a facility section centered around eight fixed visual categories.
- Business rule: The `基本参数` block must only display ship specification items that currently have values from the admin-maintained cruise data.
- Content rule: The `邮轮介绍` block uses the cruise `description` field from operations configuration.
- Interaction rule: The `邮轮设施` block shows only eight category buttons by default; clicking one expands a new detail card below it, and clicking the same category again collapses that card.
- Admin rule: Facility category icon selection in admin must be updated to a new fixed icon family that corresponds to the eight frontend facility categories.

## 2. Confirmed Constraints
- The miniapp is implemented with Vue 3 and must follow Composition API with `<script setup lang="ts">`.
- The existing app shell keeps page navigation via local tab state in `frontend/miniapp/src/App.vue`; this redesign must preserve that navigation model.
- Existing cruise detail data comes from public endpoints already used by the miniapp page: cruise detail, cabin types, and facilities.
- The backend schema does not need to change for this task; the eight-category facility experience should be built from frontend mapping logic on top of current facility category and facility payloads.
- The facility category admin pages already use a visual icon picker; the redesign should reuse that interaction model rather than reverting to free-text icon input.

## 3. Chosen Approach
- Preferred approach: Keep the backend contract unchanged, redesign the cruise detail page UI in the miniapp, and add a frontend taxonomy layer that groups current facility categories into eight fixed display categories.
- Why: This delivers the requested experience with minimal backend risk, preserves existing operations data, and keeps the admin change limited to icon semantics and visual affordance.

## 4. Miniapp Information Architecture
### 4.1 Hero And Identity Area
- Top of the page becomes a full-width image carousel using cruise gallery images when present.
- The content immediately below the carousel contains:
  - Chinese ship name
  - English ship name
  - a compact company or brand pill when source data is available
- The overall visual rhythm follows the provided mobile screenshot: light gray page background, white cards, restrained shadows, and blue accent color.

### 4.2 Dynamic Basic Parameter Section
- Replace the current horizontal metric strip with a card-based parameter grid.
- Candidate parameters include:
  - `tonnage`
  - `build_year`
  - `passenger_capacity`
  - `deck_count`
  - `length`
  - `width`
  - `room_count`
  - `crew_count`
  - `refurbish_year`
- Render only parameters with non-empty, non-zero, or otherwise meaningful values.
- Each parameter tile shows a compact Chinese label and a value with unit.

### 4.3 Description Section
- Add a standalone `邮轮介绍` card that reads from `detail.description`.
- Use calm line length, slightly elevated line-height, and small section-header ornamentation to match polished miniapp patterns.
- If description is absent, omit the section instead of showing placeholder copy.

### 4.4 Facility Section
- Keep the section title as `邮轮设施`.
- First layer: eight circular or rounded category buttons with icon plus short label.
- Second layer: a conditional detail card inserted directly below the category buttons when a category is active.
- Detail-card interaction:
  - click inactive category -> expand its detail card
  - click active category -> collapse the detail card
  - click another category -> swap the detail card content
- Detail-card content shows only facilities grouped into the selected fixed category.
- Each facility item prioritizes name, then conditionally shows location, opening hours, target audience, charge state, charge tip, and short description.

## 5. Fixed Facility Taxonomy
### 5.1 The Eight Display Categories
- 免费餐厅
- 特色餐厅
- 酒吧
- 休闲娱乐
- 亲子童趣
- 舒享舱房
- 运动健身
- 其它

### 5.2 Mapping Strategy
- Implement the grouping in frontend constants or utilities.
- Matching order:
  - exact category-name mapping first
  - keyword fallback second
  - default to `其它`
- Representative keyword buckets:
  - 免费餐厅: 主餐厅, 自助餐厅, 免费餐饮, 免费咖啡厅
  - 特色餐厅: 特色餐厅, 收费餐厅, 牛排馆, 铁板烧, 日料, 火锅
  - 酒吧: 酒吧, 酒廊, 咖啡吧, 葡萄酒吧
  - 休闲娱乐: 剧院, 秀场, 赌场, KTV, 派对, 娱乐
  - 亲子童趣: 儿童, 青少年, 托管, 电玩, 亲子
  - 舒享舱房: 套房, 礼遇, 行政酒廊, 客房服务, 舱房体验
  - 运动健身: 泳池, SPA, 健身, 球场, 跑道, 运动
- The mapping is intentionally permissive so existing运营 data can continue working without schema or content migration.

## 6. Admin Icon Design Direction
- Replace the current generic icon set with eight fixed icon meanings aligned to the fixed facility taxonomy.
- The admin list page and edit page continue to show a visual picker, but the available options become the new eight icons only.
- Icon style should be consistent line icons suitable for both admin and possible future miniapp reuse.
- Existing stored icon values outside the new set should still render safely as fallback until users re-save them.

## 7. Visual Direction
- Overall style: premium but restrained small-program detail page.
- Palette:
  - background: soft gray
  - surface: white
  - accent: marine blue
  - text: dark slate
- Layout decisions:
  - large image hero
  - generous corner radius on cards
  - lighter shadow than the current page
  - clear section spacing without oversized decoration
- Facility detail card should feel more editorial and polished than a backend list: use small metadata pills, compact separators, and deliberate spacing.

## 8. Testing Strategy
### 8.1 Miniapp Unit Tests
- Add or update tests for:
  - top-level key sections render with the redesigned labels
  - only metrics with values are shown
  - description renders from `detail.description`
  - facility detail card is hidden by default
  - clicking one category expands its grouped facilities
  - clicking the same category again collapses the detail card
  - clicking a second category replaces the expanded content

### 8.2 Admin Unit Tests
- Add or update tests for:
  - list-page icon picker exposes only the fixed eight icons
  - edit-page icon picker exposes only the fixed eight icons
  - chosen icon value is submitted correctly
  - legacy unknown icon value still displays without crashing and can be replaced

### 8.3 Focused Verification
- Run only the targeted Vitest files for touched miniapp and admin pages during implementation.
- Avoid unrelated frontend suites unless a regression requires broader coverage.

## 9. Definition Of Done
- Miniapp cruise detail visually reflects the approved reference direction.
- Top of the page is a carousel hero.
- `基本参数` only shows fields with actual data.
- `邮轮介绍` displays the cruise description field.
- `邮轮设施` defaults to icon-only category buttons and expands a detail card on click.
- Clicking the active category collapses the facility detail card.
- Admin facility category pages use the new eight-icon system.
- Targeted miniapp and admin unit tests pass.