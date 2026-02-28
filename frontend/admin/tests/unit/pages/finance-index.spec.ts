import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../pages/finance/index.vue'

const mockRequest = vi.fn()
vi.stubGlobal('useApi', () => ({ request: mockRequest }))

beforeEach(() => {
  mockRequest.mockReset()
  mockRequest.mockResolvedValue({ today_sales: 50000, today_orders: 9 })
})

describe('FinancePage', () => {
  it('调用统计接口加载财务概览', async () => {
    mount(Page)
    await flushPromises()
    expect(mockRequest).toHaveBeenCalledWith('/admin/analytics/summary')
  })

  it('渲染返回金额', async () => {
    const wrapper = mount(Page)
    await flushPromises()
    expect(wrapper.text()).toContain('50000')
  })
})
