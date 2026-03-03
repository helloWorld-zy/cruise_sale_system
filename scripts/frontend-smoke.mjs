import { chromium } from '../frontend/admin/node_modules/playwright/index.mjs'

const adminBase = 'http://127.0.0.1:3013'
const webBase = 'http://127.0.0.1:3014'
const miniBase = 'http://127.0.0.1:3015'

const adminRoutes = [
  '/', '/login', '/admin/login', '/dashboard', '/finance', '/cruises', '/cruises/create', '/cruises/1',
  '/routes', '/routes/new', '/routes/1', '/voyages', '/voyages/new', '/voyages/1',
  '/cabin-types', '/cabin-types/new', '/cabin-types/1', '/cabins', '/cabins/new', '/cabins/1',
  '/cabins/inventory?skuId=1', '/cabins/pricing?skuId=1', '/cabins/alerts', '/bookings', '/bookings/new', '/bookings/1',
  '/staff', '/settings/shop', '/notifications/templates', '/facilities', '/facilities/new', '/facilities/1',
  '/facility-categories', '/facility-categories/1', '/admin/cruises'
]

const webRoutes = [
  '/', '/search', '/web/search', '/account', '/account/login', '/cruises', '/cruises/1', '/cabins', '/cabins/1',
  '/booking', '/booking/confirm?voyage_id=1&cabin_sku_id=1&guests=2', '/booking/success?order_id=1', '/orders', '/orders/1', '/pay/1'
]

const failures = []

function collectClientErrors(page, scope, route) {
  const bag = []
  page.on('pageerror', (err) => {
    bag.push(`pageerror: ${String(err?.message || err)}`)
  })
  page.on('console', (msg) => {
    if (msg.type() === 'error') {
      const t = msg.text() || ''
      if (t.includes('Failed to load resource: the server responded with a status of 404')) {
        bag.push(`console: ${t}`)
      }
      if (t.includes('Cannot read properties of undefined') || t.includes('ReferenceError') || t.includes('TypeError')) {
        bag.push(`console: ${t}`)
      }
    }
  })
  return () => {
    if (bag.length) {
      failures.push({ scope, route, issues: bag })
    }
  }
}

async function checkRoute(context, scope, base, route) {
  const page = await context.newPage()
  const done = collectClientErrors(page, scope, route)
  try {
    const res = await page.goto(`${base}${route}`, { waitUntil: 'domcontentloaded', timeout: 20000 })
    const status = res?.status() || 0
    await page.waitForTimeout(500)
    const bodyText = (await page.locator('body').innerText()).trim()

    if (status >= 400) {
      failures.push({ scope, route, issues: [`http status ${status}`] })
    }
    if (!bodyText || bodyText.length < 2) {
      failures.push({ scope, route, issues: ['blank or near-empty page body'] })
    }
    if (bodyText.includes('Page not found') || bodyText.includes('Cannot find any path matching')) {
      failures.push({ scope, route, issues: ['rendered page-not-found text'] })
    }
  } catch (err) {
    failures.push({ scope, route, issues: [`navigation failed: ${String(err?.message || err)}`] })
  } finally {
    done()
    await page.close()
  }
}

async function checkMiniapp(context) {
  const page = await context.newPage()
  const done = collectClientErrors(page, 'miniapp', '/')
  try {
    const res = await page.goto(miniBase, { waitUntil: 'domcontentloaded', timeout: 20000 })
    const status = res?.status() || 0
    if (status >= 400) {
      failures.push({ scope: 'miniapp', route: '/', issues: [`http status ${status}`] })
    }

    const expectedTabs = ['登录', '下单', '舱房列表', '舱房详情', '支付']
    for (const tab of expectedTabs) {
      const btn = page.locator('.preview-tabs').getByRole('button', { name: tab })
      if ((await btn.count()) === 0) {
        failures.push({ scope: 'miniapp', route: '/', issues: [`tab missing: ${tab}`] })
        continue
      }
      await btn.click()
      await page.waitForTimeout(400)
      const text = (await page.locator('.preview-stage').innerText()).trim()
      if (!text || text.length < 2) {
        failures.push({ scope: 'miniapp', route: '/', issues: [`tab blank content: ${tab}`] })
      }
    }
  } catch (err) {
    failures.push({ scope: 'miniapp', route: '/', issues: [`navigation failed: ${String(err?.message || err)}`] })
  } finally {
    done()
    await page.close()
  }
}

async function main() {
  const browser = await chromium.launch({ headless: true })
  const context = await browser.newContext()
  await context.addInitScript(() => {
    try {
      localStorage.setItem('admin_token', 'smoke-token')
      sessionStorage.setItem('auth_token', 'smoke-token')
    } catch {}
  })

  for (const route of adminRoutes) {
    await checkRoute(context, 'admin', adminBase, route)
  }

  for (const route of webRoutes) {
    await checkRoute(context, 'web', webBase, route)
  }

  await checkMiniapp(context)

  await browser.close()

  if (failures.length) {
    console.log(JSON.stringify({ ok: false, failures }, null, 2))
    process.exit(1)
  }

  console.log(JSON.stringify({ ok: true, checked: { admin: adminRoutes.length, web: webRoutes.length, miniappTabs: 5 } }, null, 2))
}

main().catch((err) => {
  console.error(err)
  process.exit(1)
})
