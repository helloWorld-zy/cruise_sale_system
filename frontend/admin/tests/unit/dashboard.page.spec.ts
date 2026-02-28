import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { nextTick } from 'vue'
import Page from '../../pages/dashboard/index.vue'

const mockRequest = vi.fn()
vi.stubGlobal('useApi', () => ({ request: mockRequest }))

beforeEach(() => {
  mockRequest.mockReset()
  mockRequest.mockResolvedValue({
    today_sales: 50000,
    weekly_trend: [1000, 2000, 3000, 4000, 5000, 6000, 7000],
    today_orders: 12,
  })
})

describe('Dashboard', () => {
  it('调用 analytics summary API', async () => {
    mount(Page)
    await flushPromises()
    expect(mockRequest).toHaveBeenCalledWith('/admin/analytics/summary')
  })

  it('渲染真实统计数据', async () => {
    const wrapper = mount(Page)
    await flushPromises()
    expect(wrapper.text()).toContain('50000')
    expect(wrapper.text()).toContain('12')
  })

  it('API 失败时显示错误', async () => {
    mockRequest.mockRejectedValueOnce(new Error('api error'))
    const wrapper = mount(Page)
    await flushPromises()
    expect(wrapper.text()).toContain('api error')
  })

  it('加载时显示 Loading', () => {
    mockRequest.mockImplementationOnce(() => new Promise(() => {}))
    const wrapper = mount(Page)
    return nextTick().then(() => {
      expect(wrapper.text()).toContain('Loading')
    })
  })
})
