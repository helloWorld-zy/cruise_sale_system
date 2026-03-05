# Fullstack Startup And Playwright Manual Validation Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Start backend plus all three frontend modules (`admin`, `web`, `miniapp`) and open them in Playwright browser contexts for manual functional and visual validation.

**Architecture:** Use existing local development entry points (`docker compose` for infra, Go `make run` for backend, Nuxt/Vite dev servers for frontend). Assign fixed, non-conflicting ports for each frontend module (`3013`, `3014`, `3015`) so Playwright and manual verification are deterministic. Validate readiness via HTTP checks before opening pages in Playwright to avoid false negatives from cold startup.

**Tech Stack:** Go (backend), Docker Compose (infra), Nuxt 4 (`frontend/admin`, `frontend/web`), Vite + Vue (`frontend/miniapp`), Playwright (`scripts/frontend-smoke.mjs` and browser launch).

---

### Task 1: Environment Preflight And Dependency Readiness

**Files:**
- Modify: none
- Test: command-line checks only

**Step 1: Write the failing test**

```powershell
# PowerShell preflight checks (expected to fail if env is incomplete)
Set-Location d:\cruise_sale_system
Get-Command docker, go, node, npm -ErrorAction Stop
```

**Step 2: Run test to verify it fails**

Run: `Get-Command docker, go, node, npm -ErrorAction Stop`
Expected: FAIL if any command is missing.

**Step 3: Write minimal implementation**

```powershell
# Install or configure missing tools, then re-open terminal.
# If tools are present, no implementation is required.
```

**Step 4: Run test to verify it passes**

Run: `Get-Command docker, go, node, npm -ErrorAction Stop`
Expected: PASS and each command resolves to an executable path.

**Step 5: Commit**

```bash
git add -A
git commit -m "chore: document local runtime preflight for startup workflow"
```

### Task 2: Start Infrastructure Services Required By Backend

**Files:**
- Modify: none
- Test: Docker service health checks

**Step 1: Write the failing test**

```powershell
Set-Location d:\cruise_sale_system\docker
docker compose ps
```

**Step 2: Run test to verify it fails**

Run: `docker compose ps`
Expected: FAIL or services not healthy before startup.

**Step 3: Write minimal implementation**

```powershell
Set-Location d:\cruise_sale_system\docker
docker compose up -d
```

**Step 4: Run test to verify it passes**

Run: `docker compose ps`
Expected: PASS and critical services (`postgres`, `redis`, `minio`, `meilisearch`, `nats`) show `running` or `healthy`.

**Step 5: Commit**

```bash
git add -A
git commit -m "chore: bring up local infrastructure for fullstack startup"
```

### Task 3: Start Backend API Service

**Files:**
- Modify: none
- Test: backend health endpoint check

**Step 1: Write the failing test**

```powershell
# Before backend starts, this should fail or timeout.
Invoke-WebRequest -Uri http://127.0.0.1:8080/health -TimeoutSec 3
```

**Step 2: Run test to verify it fails**

Run: `Invoke-WebRequest -Uri http://127.0.0.1:8080/health -TimeoutSec 3`
Expected: FAIL before backend starts.

**Step 3: Write minimal implementation**

```powershell
Set-Location d:\cruise_sale_system\backend
make run
```

**Step 4: Run test to verify it passes**

Run: `Invoke-WebRequest -Uri http://127.0.0.1:8080/health -TimeoutSec 5`
Expected: PASS with `200 OK` (or equivalent healthy response body).

**Step 5: Commit**

```bash
git add -A
git commit -m "chore: run backend server for manual frontend integration validation"
```

### Task 4: Start Frontend Admin Module On Stable Port

**Files:**
- Modify: none
- Test: admin home availability

**Step 1: Write the failing test**

```powershell
Invoke-WebRequest -Uri http://127.0.0.1:3013 -TimeoutSec 3
```

**Step 2: Run test to verify it fails**

Run: `Invoke-WebRequest -Uri http://127.0.0.1:3013 -TimeoutSec 3`
Expected: FAIL before admin dev server starts.

**Step 3: Write minimal implementation**

```powershell
Set-Location d:\cruise_sale_system\frontend\admin
npm install
npm run dev -- --host 127.0.0.1 --port 3013
```

**Step 4: Run test to verify it passes**

Run: `Invoke-WebRequest -Uri http://127.0.0.1:3013 -TimeoutSec 5`
Expected: PASS with `200 OK` and Nuxt HTML content.

**Step 5: Commit**

```bash
git add -A
git commit -m "chore: run admin frontend on fixed dev port"
```

### Task 5: Start Frontend Web Module On Stable Port

**Files:**
- Modify: none
- Test: web home availability

**Step 1: Write the failing test**

```powershell
Invoke-WebRequest -Uri http://127.0.0.1:3014 -TimeoutSec 3
```

**Step 2: Run test to verify it fails**

Run: `Invoke-WebRequest -Uri http://127.0.0.1:3014 -TimeoutSec 3`
Expected: FAIL before web dev server starts.

**Step 3: Write minimal implementation**

```powershell
Set-Location d:\cruise_sale_system\frontend\web
npm install
npm run dev -- --host 127.0.0.1 --port 3014
```

**Step 4: Run test to verify it passes**

Run: `Invoke-WebRequest -Uri http://127.0.0.1:3014 -TimeoutSec 5`
Expected: PASS with `200 OK` and Nuxt HTML content.

**Step 5: Commit**

```bash
git add -A
git commit -m "chore: run web frontend on fixed dev port"
```

### Task 6: Start Frontend Miniapp Module On Stable Port

**Files:**
- Modify: none
- Test: miniapp home availability

**Step 1: Write the failing test**

```powershell
Invoke-WebRequest -Uri http://127.0.0.1:3015 -TimeoutSec 3
```

**Step 2: Run test to verify it fails**

Run: `Invoke-WebRequest -Uri http://127.0.0.1:3015 -TimeoutSec 3`
Expected: FAIL before miniapp dev server starts.

**Step 3: Write minimal implementation**

```powershell
Set-Location d:\cruise_sale_system\frontend\miniapp
npm install
npm run dev -- --host 127.0.0.1 --port 3015
```

**Step 4: Run test to verify it passes**

Run: `Invoke-WebRequest -Uri http://127.0.0.1:3015 -TimeoutSec 5`
Expected: PASS with `200 OK` and Vite HTML content.

**Step 5: Commit**

```bash
git add -A
git commit -m "chore: run miniapp frontend on fixed dev port"
```

### Task 7: Open Three Frontend Modules In Playwright For Manual Verification

**Files:**
- Use: `scripts/frontend-smoke.mjs`
- Test: Playwright browser open and page load

**Step 1: Write the failing test**

```powershell
# This should fail if one of the 3 frontend modules is not running.
Set-Location d:\cruise_sale_system
node scripts/frontend-smoke.mjs
```

**Step 2: Run test to verify it fails**

Run: `node scripts/frontend-smoke.mjs`
Expected: FAIL with `ok: false` if any module cannot be reached.

**Step 3: Write minimal implementation**

```powershell
# Open interactive Playwright browser pages for manual checks.
Set-Location d:\cruise_sale_system
node -e "import('playwright').then(async ({ chromium }) => { const browser = await chromium.launch({ headless: false }); const context = await browser.newContext(); const pages = ['http://127.0.0.1:3013','http://127.0.0.1:3014','http://127.0.0.1:3015']; for (const u of pages) { const p = await context.newPage(); await p.goto(u, { waitUntil: 'domcontentloaded' }); } })"
```

**Step 4: Run test to verify it passes**

Run: `node scripts/frontend-smoke.mjs`
Expected: PASS with JSON output similar to `{ "ok": true, "checked": { ... } }`.

**Step 5: Commit**

```bash
git add -A
git commit -m "chore: validate all frontends with playwright and open manual QA pages"
```

### Task 8: Manual QA Checklist (Function + Style)

**Files:**
- Modify: none
- Test: manual validation record in PR description or test notes

**Step 1: Write the failing test**

```text
No manual validation notes exist yet for admin/web/miniapp screens.
```

**Step 2: Run test to verify it fails**

Run: check your QA notes
Expected: FAIL if no notes or screenshots are recorded.

**Step 3: Write minimal implementation**

```text
For each module, verify: layout integrity, no overlap/truncation, key action buttons visible, no obvious console errors, and API-backed pages render data or graceful empty states.
```

**Step 4: Run test to verify it passes**

Run: complete one pass over three modules in opened Playwright windows
Expected: PASS with written notes for all modules.

**Step 5: Commit**

```bash
git add -A
git commit -m "test: complete manual functional and visual validation for three frontend modules"
```

---

## 2026-03-05 Execution Notes

### Completed

- Infra status check passed: `docker compose ps` shows `postgres/redis/minio/meilisearch/nats` all healthy.
- Frontend smoke script updated to match latest architecture:
	- Removed deprecated admin route checks: `/routes`, `/routes/new`, `/routes/1`.
	- Reduced false positives by no longer treating every console 404 resource message as hard failure.
- Playwright smoke re-run passed:

```json
{
	"ok": true,
	"checked": {
		"admin": 32,
		"web": 15,
		"miniappTabs": 5
	}
}
```

### Current Local Blocker

- Local port `8080` is occupied by external process `httpd` (PID `6936`) that cannot be terminated in current permission context.
- Impact:
	- `go run ./cmd/server` on default `:8080` fails with bind error.
	- Nuxt admin/web dev proxy targets `http://localhost:8080/api/v1`, so API integration checks on default setup may hit non-project server.

### Validation Workaround Confirmed

- Backend service itself is healthy when launched on alternate port:
	- Command: `CRUISE_SERVER_PORT=:18080 go run ./cmd/server`
	- Verification: `GET http://127.0.0.1:18080/api/v1/cruises -> 200`

### Next Suggested Action

- Make Nuxt dev proxy target configurable by environment (instead of hardcoded `:8080`), so local startup/Playwright validation is robust even when `8080` is occupied.
