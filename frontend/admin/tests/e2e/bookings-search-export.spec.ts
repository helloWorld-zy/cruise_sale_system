import { expect, test } from '@playwright/test'

test('bookings page supports search, filters and export with current query', async ({ page }) => {
  let lastListQuery = ''
  let lastExportQuery = ''

  await page.addInitScript(() => {
    window.localStorage.setItem('admin_token', 'e2e-token')
  })

  await page.route('**/admin/bookings/export**', async (route) => {
    lastExportQuery = route.request().url()
    await route.fulfill({
      status: 200,
      contentType: 'text/csv; charset=utf-8',
      body: 'booking_no,phone,voyage_code,cruise_name\nBK-2026-001,13800000001,VOY-ALPHA,海洋量子号',
    })
  })

  await page.route('**/admin/bookings**', async (route) => {
    const requestUrl = new URL(route.request().url())
    if (requestUrl.pathname.endsWith('/export')) {
      await route.fallback()
      return
    }
    lastListQuery = requestUrl.search
    const keyword = requestUrl.searchParams.get('keyword') || ''
    const voyageCode = requestUrl.searchParams.get('voyage_code') || ''
    const list = [
      { id: 1, booking_no: 'BK-2026-001', status: 'paid', total_cents: 18800, user_id: 9, phone: '13800000001', voyage_code: 'VOY-ALPHA', cruise_name: '海洋量子号', created_at: '2026-03-07T08:00:00Z' },
      { id: 2, booking_no: 'BK-2026-002', status: 'pending_payment', total_cents: 9900, user_id: 10, phone: '13900000002', voyage_code: 'VOY-BETA', cruise_name: '海洋光谱号', created_at: '2026-03-07T09:00:00Z' },
    ].filter((item) => {
      const matchesKeyword = !keyword || [item.booking_no, item.phone, item.voyage_code, item.cruise_name, item.status, String(item.total_cents)].some((value) => value.includes(keyword))
      const matchesVoyageCode = !voyageCode || item.voyage_code.includes(voyageCode)
      return matchesKeyword && matchesVoyageCode
    })
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ code: 0, data: { list, total: list.length } }),
    })
  })

  await page.goto('/bookings')
  await expect(page.getByRole('heading', { name: '订单管理' })).toBeVisible()
  await expect(page.getByText('BK-2026-001')).toBeVisible()
  await expect(page.getByText('BK-2026-002')).toBeVisible()

  await page.locator('[data-test="booking-search-input"]').fill('BK-2026-001')
  await page.locator('[data-test="filter-voyage-code"]').fill('VOY-ALPHA')
  await page.locator('[data-test="booking-search-submit"]').click()

  await expect.poll(() => lastListQuery).toContain('keyword=BK-2026-001')
  await expect.poll(() => lastListQuery).toContain('voyage_code=VOY-ALPHA')
  await expect(page.getByText('BK-2026-001')).toBeVisible()
  await expect(page.getByText('BK-2026-002')).toHaveCount(0)

  await page.locator('[data-test="export"]').click()
  await expect.poll(() => lastExportQuery).toContain('/admin/bookings/export')
  await expect.poll(() => lastExportQuery).toContain('keyword=BK-2026-001')
  await expect.poll(() => lastExportQuery).toContain('voyage_code=VOY-ALPHA')
})