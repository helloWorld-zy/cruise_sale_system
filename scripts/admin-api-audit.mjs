import fs from 'node:fs'
import path from 'node:path'
import { chromium } from '../frontend/admin/node_modules/playwright/index.mjs'

const adminBase = 'http://127.0.0.1:3013'
const backendBase = 'http://localhost:8080'

// Focus on admin management modules to audit real API calls and statuses.
const adminRoutes = [
  '/dashboard',
  '/finance',
  '/cruises',
  '/routes',
  '/voyages',
  '/cabin-types',
  '/cabins',
  '/cabins/alerts',
  '/cabins/inventory?skuId=1',
  '/cabins/pricing?skuId=1',
  '/bookings',
  '/staff',
  '/settings/shop',
  '/notifications/templates',
  '/facilities',
  '/facility-categories',
]

const outputPathArgIndex = process.argv.indexOf('--out')
const usernameArgIndex = process.argv.indexOf('--username')
const passwordArgIndex = process.argv.indexOf('--password')
const outputPath =
  outputPathArgIndex > -1 && process.argv[outputPathArgIndex + 1]
    ? process.argv[outputPathArgIndex + 1]
    : path.resolve(process.cwd(), 'docs', 'plans', 'evidence', 'admin-api-audit-latest.json')
const username =
  usernameArgIndex > -1 && process.argv[usernameArgIndex + 1]
    ? process.argv[usernameArgIndex + 1]
    : 'admin'
const password =
  passwordArgIndex > -1 && process.argv[passwordArgIndex + 1]
    ? process.argv[passwordArgIndex + 1]
    : 'admin123'

async function getAdminToken() {
  const res = await fetch(`${backendBase}/api/v1/admin/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
  })

  const payload = await res.json().catch(() => null)
  if (!res.ok) {
    throw new Error(`login failed: HTTP ${res.status}`)
  }

  const token = payload?.data?.token || payload?.token || ''
  if (!token) {
    throw new Error('login response has no token')
  }
  return token
}

function toApiKey(method, url) {
  try {
    const parsed = new URL(url)
    return `${method.toUpperCase()} ${parsed.pathname}${parsed.search}`
  } catch {
    return `${method.toUpperCase()} ${url}`
  }
}

async function auditRoute(context, route) {
  const page = await context.newPage()
  const apis = new Map()
  const requestFailures = []
  const pageErrors = []

  page.on('response', async (response) => {
    const url = response.url()
    if (!url.includes('/api/v1/')) return

    const req = response.request()
    const key = toApiKey(req.method(), url)
    const entry = apis.get(key) || {
      method: req.method().toUpperCase(),
      url,
      statuses: [],
      count: 0,
    }

    entry.count += 1
    entry.statuses.push(response.status())
    apis.set(key, entry)
  })

  page.on('requestfailed', (req) => {
    const url = req.url()
    if (!url.includes('/api/v1/')) return
    requestFailures.push({
      method: req.method().toUpperCase(),
      url,
      errorText: req.failure()?.errorText || 'request failed',
    })
  })

  page.on('pageerror', (err) => {
    pageErrors.push(String(err?.message || err))
  })

  let navigationStatus = 0
  let navigationError = ''
  try {
    const res = await page.goto(`${adminBase}${route}`, { waitUntil: 'domcontentloaded', timeout: 20000 })
    navigationStatus = res?.status() || 0
    await page.waitForTimeout(1200)
  } catch (err) {
    navigationError = String(err?.message || err)
  }

  await page.close()

  const apiCalls = [...apis.values()].map((item) => {
    const maxStatus = item.statuses.length ? Math.max(...item.statuses) : 0
    return {
      method: item.method,
      url: item.url,
      count: item.count,
      statuses: item.statuses,
      hasHttpError: maxStatus >= 400,
    }
  })

  return {
    route,
    navigationStatus,
    navigationError,
    apiCalls,
    requestFailures,
    pageErrors,
    ok:
      !navigationError &&
      navigationStatus < 400 &&
      requestFailures.length === 0 &&
      apiCalls.every((x) => !x.hasHttpError),
  }
}

async function main() {
  let adminToken = ''
  let loginError = ''
  try {
    adminToken = await getAdminToken()
  } catch (err) {
    loginError = String(err?.message || err)
  }

  const browser = await chromium.launch({ headless: true })
  const context = await browser.newContext()

  await context.addInitScript((token) => {
    try {
      if (token) {
        localStorage.setItem('admin_token', token)
      }
    } catch {}
  }, adminToken)

  const startedAt = new Date().toISOString()
  const routeResults = []

  for (const route of adminRoutes) {
    const result = await auditRoute(context, route)
    routeResults.push(result)
  }

  await browser.close()

  const totalApiCalls = routeResults.reduce((sum, r) => sum + r.apiCalls.reduce((x, c) => x + c.count, 0), 0)
  const failedRoutes = routeResults.filter((r) => !r.ok)

  const report = {
    ok: failedRoutes.length === 0,
    auth: {
      username,
      loginOk: Boolean(adminToken),
      loginError,
    },
    startedAt,
    finishedAt: new Date().toISOString(),
    base: adminBase,
    checkedRoutes: adminRoutes.length,
    totalApiCalls,
    failedRouteCount: failedRoutes.length,
    failedRoutes: failedRoutes.map((r) => r.route),
    routeResults,
  }

  fs.mkdirSync(path.dirname(outputPath), { recursive: true })
  fs.writeFileSync(outputPath, JSON.stringify(report, null, 2), 'utf8')

  console.log(JSON.stringify({
    ok: report.ok,
    checkedRoutes: report.checkedRoutes,
    totalApiCalls: report.totalApiCalls,
    failedRouteCount: report.failedRouteCount,
    outputPath,
  }, null, 2))

  if (!report.ok) {
    process.exit(1)
  }
}

main().catch((err) => {
  console.error(err)
  process.exit(1)
})
