import { expect, test } from '@playwright/test'
import fs from 'node:fs'
import path from 'node:path'

test('capture Task20 UI evidence pages', async ({ page }) => {
  await page.addInitScript(() => {
    window.localStorage.setItem('admin_token', 'e2e-token')
  })

  await page.route('**/admin/analytics/summary**', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ code: 0, data: { today_sales: 12000, today_orders: 9, weekly_trend: [1, 2, 3, 4, 5, 6, 7] } }),
    })
  })

  await page.route('**/admin/shop-info**', async (route) => {
    if (route.request().method() === 'PUT') {
      await route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify({ code: 0, data: { ok: true } }) })
      return
    }
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ code: 0, data: { name: 'Cruise Shop', contact_phone: '13800138000', contact_email: 'ops@cruise.test', icp_number: 'ICP-TEST' } }),
    })
  })

  await page.route('**/admin/notification-templates**', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ code: 0, data: [{ id: 1, event_type: 'order_paid', channel: 'sms', template: '订单{{.OrderNo}}已支付' }] }),
    })
  })

  await page.route('**/admin/staffs**', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ code: 0, data: [{ id: 101, real_name: '审核员A', email: 'audit@cruise.test', role: 'operator', status: 1 }] }),
    })
  })

  await page.route('**/admin/cabins**', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ code: 0, data: { list: [{ id: 1, code: 'C-1', status: 1 }], total: 1 } }),
    })
  })

  await page.route('**/admin/cabins/**', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ code: 0, data: { list: [{ id: 1, date: '2026-03-03', occupancy: 2, price_cents: 19900 }], total: 1 } }),
    })
  })

  await page.route('**/admin/cruises**', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ code: 0, data: { list: [{ id: 11, name: 'E2E Cruise', code: 'E2E-11', company_id: 1, tonnage: 10000, passenger_capacity: 1200, status: 1 }], total: 1 } }),
    })
  })

  await page.route('**/admin/bookings**', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        code: 0,
        data: {
          list: [
            { id: 701, status: 'pending_payment', total_cents: 19900, voyage_id: 9, user_id: 1, created_at: '2026-03-03T10:00:00Z' },
            { id: 702, status: 'paid', total_cents: 25900, voyage_id: 9, user_id: 1, created_at: '2026-03-03T11:00:00Z' },
          ],
          total: 2,
        },
      }),
    })
  })

  const evidenceDir = path.resolve(process.cwd(), '..', '..', 'docs', 'plans', 'evidence', 'task20')
  fs.mkdirSync(evidenceDir, { recursive: true })

  const pages: Array<{ route: string; heading: string; file: string }> = [
    { route: '/cruises', heading: '邮轮管理', file: 'task8-cruises.png' },
    { route: '/cabins', heading: '舱位商品管理', file: 'task8-cabins.png' },
    { route: '/cabins/pricing?skuId=1', heading: '价格矩阵', file: 'task9-pricing.png' },
    { route: '/finance', heading: 'Finance', file: 'task21-finance.png' },
    { route: '/bookings', heading: 'Bookings', file: 'task21-bookings-export.png' },
    { route: '/staff', heading: '员工管理', file: 'task22-staff.png' },
    { route: '/settings/shop', heading: '店铺设置', file: 'task23-shop.png' },
    { route: '/notifications/templates', heading: '通知模板', file: 'task24-templates.png' },
    { route: '/dashboard', heading: 'Dashboard', file: 'task26-dashboard.png' },
  ]

  for (const item of pages) {
    await page.goto(item.route)
    await expect(page.getByRole('heading', { name: item.heading })).toBeVisible()
    if (item.route === '/bookings') {
      await expect(page.locator('[data-test="export"]')).toBeVisible()
    }
    await page.screenshot({ path: path.join(evidenceDir, item.file), fullPage: true })
  }
})
