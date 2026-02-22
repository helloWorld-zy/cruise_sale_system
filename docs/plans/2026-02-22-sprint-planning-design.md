# Sprint Planning Design (TDD) — CruiseBooking

**Goal:** Define the structure and execution approach for writing full TDD-driven Sprint plans (Sprint 1–16) for the CruiseBooking platform.

**Scope:** Planning documentation only. No code changes.

---

## 1) Planning Approach

**Chosen option:** Plan A (per-Sprint files, per-stack sections).

Each Sprint plan is a standalone document organized into sections:
- Backend (Go + Gin + GORM)
- Admin Frontend (Nuxt 4)
- Web Frontend (Nuxt 4 SSR)
- Mini Program (uni-app)
- CI/CD + Infrastructure

This structure supports parallel execution and keeps dependencies explicit.

---

## 2) TDD Task Format (Mandatory)

Every task follows RED-GREEN-REFACTOR with exact file paths and runnable commands:

1. **Write the failing test** (real test code)
2. **Run test and confirm failure** (specific expected failure)
3. **Write minimal implementation** (only what makes the test pass)
4. **Run test and confirm pass**
5. **Commit** (small, focused commit)

Each step is a 2–5 minute action to ensure the plan is executable and easy to follow.

---

## 3) File Naming Conventions

Plans are saved as:
- `docs/plans/2026-02-22-sprint01.md`
- `docs/plans/2026-02-22-sprint02.md`
- ...
- `docs/plans/2026-02-22-sprint16.md`

The master roadmap is saved at:
- `docs/plans/2026-02-22-master-roadmap.md`

---

## 4) Content Requirements

Each Sprint plan MUST include:
- Goal, Architecture, Tech Stack
- Dependencies and integration points
- Precise API endpoints (if applicable)
- Database schema changes (if applicable)
- Test strategy and commands
- Exact file paths for every change
- Full code for test and minimal implementation steps

Coverage is always 100% for unit, integration, and E2E tests as defined in `功能列表.md`.

---

## 5) Execution Rules

- No production code without a failing test first.
- Use minimal code to pass the test.
- Refactor only when all tests are green.
- Commit after each task.
- Keep tasks small and reversible.

---

## 6) Sprint Breakdown

Sprint scope is defined by the agreed breakdown:

- Sprint 1: Infrastructure + Cruise Introduction module
- Sprint 2: Cabin product management
- Sprint 3: Booking flow + User system
- Sprint 4: Orders + Payments + Notifications + Analytics
- Sprint 5–8: Phase 2 features
- Sprint 9–12: Phase 3 features
- Sprint 13–16: Phase 4 features

This breakdown is also reflected in `docs/plans/2026-02-22-master-roadmap.md`.

---

## 7) Validation Checklist

Before considering each Sprint plan complete:
- All tasks follow the TDD step structure
- Every file path is explicit and correct
- All commands are runnable and include expected output
- Coverage commands are included for each subproject
- No missing dependencies or tools
