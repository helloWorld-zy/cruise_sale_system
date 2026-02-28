import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/cabins/index.vue'

// 模拟 API 请求。
const mockRequest = vi.fn()
vi.stubGlobal('useApi', () => ({ request: mockRequest }))

beforeEach(() => {
  mockRequest.mockReset()
  mockRequest.mockResolvedValue({ data: [{ id: 1, voyage_id: 2, cabin_type_id: 3, total: 10, available: 8 }] })
})

describe('CabinsIndexPage', () => {
  it('调用 API 获取舱位列表', async () => {
    mount(Page)
    await flushPromises()
    expect(mockRequest).toHaveBeenCalledWith('/cabins')
  })

  it('渲染返回数据', async () => {
    const wrapper = mount(Page)
    await flushPromises()
    expect(wrapper.text()).toContain('8')
  })

  it('API 失败时显示错误', async () => {
    mockRequest.mockRejectedValueOnce(new Error('load fail'))
    const wrapper = mount(Page)
    await flushPromises()
    expect(wrapper.text()).toContain('load fail')
  })
})
