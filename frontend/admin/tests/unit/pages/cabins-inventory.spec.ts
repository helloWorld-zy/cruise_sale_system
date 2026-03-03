import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/cabins/inventory.vue'

const mockRequest = vi.fn()
vi.stubGlobal('useApi', () => ({ request: mockRequest }))
vi.stubGlobal('useRoute', () => ({ query: { skuId: '7' } }))

beforeEach(() => {
  mockRequest.mockReset()
  mockRequest.mockResolvedValue({ data: { total: 20, available: 9 } })
})

describe('CabinsInventoryPage', () => {
  it('挂载时调用库存查询 API', async () => {
    mount(Page)
    await flushPromises()
    expect(mockRequest).toHaveBeenCalledWith('/cabins/7/inventory')
  })

  it('点击调整触发调整接口', async () => {
    const wrapper = mount(Page)
    await flushPromises()
    const inputs = wrapper.findAll('input')
    await inputs[1].setValue('2')
    await wrapper.findAll('button')[1].trigger('click')
    await flushPromises()
    expect(mockRequest).toHaveBeenCalledWith('/cabins/7/inventory/adjust', expect.objectContaining({ method: 'POST' }))
  })

  it('保存预警阈值调用接口', async () => {
    const wrapper = mount(Page)
    await flushPromises()
    const inputs = wrapper.findAll('input')
    await inputs[0].setValue('3')
    await wrapper.findAll('button')[0].trigger('click')
    await flushPromises()
    expect(mockRequest).toHaveBeenCalledWith('/cabins/7/alert-threshold', expect.objectContaining({ method: 'PUT' }))
  })
})
