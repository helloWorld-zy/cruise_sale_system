# Sprint 9-16 Design

**Goal:** Provide consistent, TDD-first implementation plans for Sprints 9-16 covering backend, admin, web, miniapp, and CI with 100% coverage gates.

**Scope:** Social sharing, reviews/community, intelligent operations, pricing/i18n, performance, security, and production deployment/monitoring.

**Structure:** Each sprint plan mirrors `docs/plans/2026-02-22-sprint01.md`.
- Header with goal/architecture/tech stack
- Part A: Backend (domain, repo, service, handler)
- Part B: Admin (Nuxt 4 admin UI)
- Part C: Web (Nuxt 4 public UI)
- Part D: Miniapp (uni-app)
- Part E: CI (coverage checks + workflow updates)
- Final Verification (full coverage commands)

**TDD Pattern:**
1. Write failing test
2. Run test and confirm failure
3. Implement minimal code
4. Run test and confirm pass
5. Commit

**Frontend Conventions:**
- Vue 3 Composition API with `<script setup lang="ts">`
- Vitest + Vue Test Utils for unit tests
- Playwright for E2E tests

**Backend Conventions:**
- Go (Gin + GORM) layered architecture (handler/service/repository/domain)
- SQLite for repository unit tests
- go test with coverage flags for 100% coverage

**CI/CD Conventions:**
- GitHub Actions workflows for backend/admin/web/miniapp
- Coverage gates set to 100%
- Each sprint adds checks relevant to new modules (e.g., load tests, security scans, deploy manifests lint)
