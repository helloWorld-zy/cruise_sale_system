# Miniapp Cruise Detail Facility Redesign Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Redesign the miniapp cruise detail page to match the approved mobile reference, add dynamic basic specs and expandable fixed-category facility content, and align admin facility category icons to the same eight-category taxonomy.

**Architecture:** Keep the backend API contract unchanged and implement the behavior in the frontend layers. The miniapp cruise detail page will derive display sections and facility groups from existing cruise and facility payloads, while admin facility category pages will switch to a fixed eight-icon visual system that mirrors the miniapp taxonomy.

**Tech Stack:** Vue 3 with `<script setup lang="ts">`, Vitest, Vue Test Utils, Testing Library, existing miniapp/admin frontend structure.

---

### Task 1: Save planning artifacts and align the top-level plan

**Files:**
- Modify: `plan.md`
- Create: `docs/plans/2026-03-11-miniapp-cruise-detail-facility-design.md`
- Create: `docs/plans/2026-03-11-miniapp-cruise-detail-facility-implementation-plan.md`

**Step 1: Update the top-level plan**

Add unchecked items for the miniapp cruise detail redesign, facility grouping, admin icon update, and verification.

**Step 2: Save the design document**

Record the approved layout, interaction, and fixed eight-category mapping strategy.

**Step 3: Save the implementation plan**

Ensure later execution does not depend on hidden context.

### Task 2: Add failing miniapp tests for the redesigned cruise detail layout

**Files:**
- Modify: `frontend/miniapp/tests/unit/pages/cruise/detail.spec.ts`
- Modify: `frontend/miniapp/tests/unit/app.spec.ts`

**Step 1: Write the failing test for dynamic basic parameters and description**

Add assertions that the cruise detail page renders `基本参数` and `邮轮介绍`, shows populated metrics such as `总吨位` and `建造年份`, and omits empty metrics.

**Step 2: Write the failing test for facility toggle behavior**

Add assertions that:
- no facility detail card is visible on first render
- clicking one fixed category reveals grouped facilities
- clicking the same category again hides the grouped facilities
- clicking another fixed category swaps the visible group

**Step 3: Run the targeted miniapp tests to verify failure**

Run: `npx vitest run tests/unit/pages/cruise/detail.spec.ts tests/unit/app.spec.ts`

Expected: fail because the current page still renders the older layout and ungrouped facility list.

### Task 3: Add failing admin tests for the fixed eight-icon picker

**Files:**
- Modify: `frontend/admin/tests/unit/pages/facility-categories.spec.ts`
- Modify: `frontend/admin/tests/unit/pages/facility-categories-form.spec.ts`

**Step 1: Write the failing list-page icon test**

Assert the picker exposes the eight fixed icon values and no longer relies on the old generic set.

**Step 2: Write the failing edit-page icon test**

Assert the edit view exposes the same eight fixed icon choices and preserves fallback rendering for unknown stored icon values.

**Step 3: Run the targeted admin tests to verify failure**

Run: `npx vitest run tests/unit/pages/facility-categories.spec.ts tests/unit/pages/facility-categories-form.spec.ts`

Expected: fail because the current icon options are still the older generic set.

### Task 4: Implement the fixed facility taxonomy utilities for miniapp

**Files:**
- Create: `frontend/miniapp/src/constants/cruiseFacilityCategories.ts`
- Modify: `frontend/miniapp/pages/cruise/detail.vue`

**Step 1: Add the fixed category definitions**

Export the eight display categories with stable ids, labels, icon keys, and keyword match lists.

**Step 2: Add facility grouping helpers**

Implement helper functions that map current facility category names into the fixed eight buckets using exact-match and keyword fallback.

**Step 3: Keep the helper API small and predictable**

Return grouped facility collections that the page can consume without additional mutation-heavy logic.

### Task 5: Implement the miniapp cruise detail redesign

**Files:**
- Modify: `frontend/miniapp/pages/cruise/detail.vue`

**Step 1: Replace the old layout with the approved card structure**

Build the hero carousel, ship identity block, dynamic basic parameter grid, and description card.

**Step 2: Derive the metric list from available cruise detail fields**

Use computed state to show only fields with meaningful values.

**Step 3: Replace the existing facility tabs with fixed category buttons**

Render the eight fixed categories in a compact icon-first layout.

**Step 4: Implement expandable facility detail-card behavior**

Track the active category id, expand the selected category, collapse on second click, and switch content when another category is selected.

**Step 5: Render grouped facility items with compact metadata rows**

Display only fields with values and preserve a polished mobile visual hierarchy.

**Step 6: Run the targeted miniapp tests to verify green**

Run: `npx vitest run tests/unit/pages/cruise/detail.spec.ts tests/unit/app.spec.ts`

Expected: pass.

### Task 6: Replace admin facility category icon options with the fixed eight-icon system

**Files:**
- Modify: `frontend/admin/app/constants/facilityCategoryIcons.ts`
- Modify: `frontend/admin/app/components/facility-categories/FacilityCategoryIcon.vue`
- Modify: `frontend/admin/app/pages/facility-categories/index.vue`
- Modify: `frontend/admin/app/pages/facility-categories/[id].vue`

**Step 1: Replace the icon option constant set**

Define the eight fixed icon values and labels aligned to the approved taxonomy.

**Step 2: Add the new generated SVG icons**

Implement a consistent icon family for the eight values while keeping fallback rendering for unknown legacy icons.

**Step 3: Keep the picker behavior unchanged**

Reuse the current button-grid picker interaction so only the option set and visuals change.

**Step 4: Run the targeted admin tests to verify green**

Run: `npx vitest run tests/unit/pages/facility-categories.spec.ts tests/unit/pages/facility-categories-form.spec.ts`

Expected: pass.

### Task 7: Run focused verification across both frontends

**Files:**
- No code changes expected

**Step 1: Run all touched miniapp tests**

Run: `npx vitest run tests/unit/pages/cruise/detail.spec.ts tests/unit/app.spec.ts`

Expected: pass.

**Step 2: Run all touched admin tests**

Run: `npx vitest run tests/unit/pages/facility-categories.spec.ts tests/unit/pages/facility-categories-form.spec.ts`

Expected: pass.

**Step 3: Summarize residual risks**

Call out:
- facility grouping still depends on category-name keywords rather than backend-enforced taxonomy
- unknown legacy icon values remain supported only as fallback until users resave them
- previewing the exact small-program look in browser mode may still differ slightly from real WeChat rendering

Plan complete and saved to `docs/plans/2026-03-11-miniapp-cruise-detail-facility-implementation-plan.md`. Two execution options:

**1. Subagent-Driven (this session)** - I dispatch fresh subagent per task, review between tasks, fast iteration

**2. Parallel Session (separate)** - Open new session with executing-plans, batch execution with checkpoints

For this task, I will continue in this session and implement directly.