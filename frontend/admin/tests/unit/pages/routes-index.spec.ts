import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/routes/index.vue'

const mockRequest = vi.fn()
vi.stubGlobal('useApi', () => ({ request: mockRequest }))

beforeEach(() => {
  mockRequest.mockReset()
  mockRequest.mockResolvedValue({ data: [{ id: 1, code: 'R-1', name: 'Route 1' }] })
})

describe('RoutesIndexPage', () => {
  it('调用 routes API', async () => {
    mount(Page)
    await flushPromises()
    expect(mockRequest).toHaveBeenCalledWith('/routes')
  })

  it('渲染返回航线数据', async () => {
    const wrapper = mount(Page)
    await flushPromises()
    expect(wrapper.text()).toContain('R-1')
  })
})
