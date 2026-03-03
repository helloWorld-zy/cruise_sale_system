import { expect, test } from '@playwright/test'
import fs from 'node:fs'
import path from 'node:path'

test.describe('Task22 key E2E', () => {
  test('empty login keeps user on login page', async ({ page }) => {
    await page.goto('/login')

    await page.getByRole('button', { name: '登录' }).click()
    await expect(page).toHaveURL(/\/login/)
  })

  test('authenticated token can open cruises and dashboard with auth header', async ({ page }) => {
    let seenCruisesAuth = ''
    let seenSummaryAuth = ''

    await page.route('**/admin/cruises**', async (route) => {
      seenCruisesAuth = route.request().headers()['authorization'] || ''
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ code: 0, data: { list: [], total: 0 } }),
      })
    })

    await page.route('**/admin/analytics/summary**', async (route) => {
      seenSummaryAuth = route.request().headers()['authorization'] || ''
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ code: 0, data: { today_sales: 8800, today_orders: 6, weekly_trend: [1, 2, 3, 4, 5, 6, 7] } }),
      })
    })

    await page.addInitScript(() => {
      window.localStorage.setItem('admin_token', 'e2e-token')
    })

    await page.goto('/cruises')
    await expect(page.getByRole('heading', { name: '邮轮管理' })).toBeVisible()
    await expect.poll(() => seenCruisesAuth).toBe('Bearer e2e-token')

    const evidenceDir = path.resolve(process.cwd(), '..', '..', 'docs', 'plans', 'evidence', 'task22')
    fs.mkdirSync(evidenceDir, { recursive: true })
    await page.screenshot({ path: path.join(evidenceDir, 'admin-cruises.png'), fullPage: true })

    await page.goto('/dashboard')
    await expect(page.getByText('Dashboard')).toBeVisible()
    await expect(page.getByText('Sales')).toBeVisible()
    await expect(page.getByText('Orders')).toBeVisible()
    await expect.poll(() => seenSummaryAuth).toBe('Bearer e2e-token')
    await page.screenshot({ path: path.join(evidenceDir, 'admin-dashboard.png'), fullPage: true })
  })
})
