import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/finance/index.vue'

const mockRequest = vi.fn()
vi.stubGlobal('useApi', () => ({ request: mockRequest }))

const mountOptions = {
  global: {
    stubs: {
      AdminPageHeader: { props: ['title', 'subtitle'], template: '<div>{{ title }} {{ subtitle }}<slot /><slot name="actions" /></div>' },
      AdminDataCard: { props: ['flush'], template: '<div><slot /></div>' },
    },
  },
}

beforeEach(() => {
  mockRequest.mockReset()
  mockRequest.mockResolvedValue({ today_sales: 50000, today_orders: 9 })
})

describe('FinancePage', () => {
  it('调用统计接口加载财务概览', async () => {
    mount(Page, mountOptions)
    await flushPromises()
    expect(mockRequest).toHaveBeenCalledWith('/admin/analytics/summary')
  })

  it('渲染返回金额', async () => {
    const wrapper = mount(Page, mountOptions)
    await flushPromises()
    expect(wrapper.text()).toContain('50000')
  })
})
