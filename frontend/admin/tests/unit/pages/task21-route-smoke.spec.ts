import { describe, it, expect, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import fs from 'node:fs'
import path from 'node:path'

import DashboardPage from '../../../app/pages/dashboard/index.vue'
import FinancePage from '../../../app/pages/finance/index.vue'

const mockRequest = vi.fn().mockResolvedValue({
  today_sales: 12345,
  weekly_trend: [1, 2, 3, 4, 5, 6, 7],
  today_orders: 8,
})

vi.stubGlobal('useApi', () => ({ request: mockRequest }))

const mountOptions = {
  global: {
    stubs: {
      AdminPageHeader: { props: ['title', 'subtitle'], template: '<div>{{ title }} {{ subtitle }}<slot /><slot name="actions" /></div>' },
      AdminDataCard: { props: ['flush'], template: '<div><slot /></div>' },
      AdminStatusTag: { props: ['text'], template: '<span>{{ text }}</span>' },
    },
  },
}

describe('Task21 route smoke', () => {
  it('login/dashboard/finance route pages exist and render', async () => {
    const root = process.cwd()
    const loginPage = path.join(root, 'app', 'pages', 'login.vue')
    const dashboardPage = path.join(root, 'app', 'pages', 'dashboard', 'index.vue')
    const financePage = path.join(root, 'app', 'pages', 'finance', 'index.vue')

    expect(fs.existsSync(loginPage)).toBe(true)
    expect(fs.existsSync(dashboardPage)).toBe(true)
    expect(fs.existsSync(financePage)).toBe(true)

    const dashboard = mount(DashboardPage, mountOptions)
    await flushPromises()
    expect(dashboard.text()).toContain('Dashboard')

    const finance = mount(FinancePage, mountOptions)
    await flushPromises()
    expect(finance.text()).toContain('Finance')
  })

  it('legacy pages directory no longer has dashboard/finance entries', () => {
    const root = process.cwd()
    const oldDashboard = path.join(root, 'pages', 'dashboard', 'index.vue')
    const oldFinance = path.join(root, 'pages', 'finance', 'index.vue')

    expect(fs.existsSync(oldDashboard)).toBe(false)
    expect(fs.existsSync(oldFinance)).toBe(false)
  })
})
