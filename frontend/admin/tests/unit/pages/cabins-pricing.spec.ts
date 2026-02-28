import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/cabins/pricing.vue'

const mockRequest = vi.fn()
vi.stubGlobal('useApi', () => ({ request: mockRequest }))
vi.stubGlobal('useRoute', () => ({ query: { skuId: '9' } }))

beforeEach(() => {
  mockRequest.mockReset()
  mockRequest.mockResolvedValue({ data: [{ date: '2026-05-01', occupancy: 2, price_cents: 19900 }] })
})

describe('CabinsPricingPage', () => {
  it('挂载时加载价格列表', async () => {
    mount(Page)
    await flushPromises()
    expect(mockRequest).toHaveBeenCalledWith('/cabins/9/prices')
  })

  it('渲染价格数据', async () => {
    const wrapper = mount(Page)
    await flushPromises()
    expect(wrapper.text()).toContain('19900')
  })
})
