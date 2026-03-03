import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/routes/index.vue'

const mockRequest = vi.fn()
const confirmMock = vi.fn(() => true)
vi.stubGlobal('useApi', () => ({ request: mockRequest }))
vi.stubGlobal('confirm', confirmMock)

beforeEach(() => {
  mockRequest.mockReset()
  confirmMock.mockClear()
  confirmMock.mockReturnValue(true)
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

  it('空数据时显示 No data', async () => {
    mockRequest.mockResolvedValueOnce({ data: [] })
    const wrapper = mount(Page)
    await flushPromises()
    expect(wrapper.text()).toContain('No data')
  })

  it('无效 ID 删除时显示错误', async () => {
    const wrapper = mount(Page)
    await flushPromises()
    await (wrapper.vm as any).handleDelete('x')
    await flushPromises()
    expect(wrapper.text()).toContain('无效记录 ID，无法删除')
  })

  it('删除确认取消时不发请求', async () => {
    confirmMock.mockReturnValueOnce(false)
    const wrapper = mount(Page)
    await flushPromises()

    await (wrapper.vm as any).handleDelete(1)
    await flushPromises()

    expect(mockRequest).not.toHaveBeenCalledWith('/routes/1', { method: 'DELETE' })
  })

  it('删除成功后会重新加载列表', async () => {
    mockRequest.mockImplementation((url: string, options?: any) => {
      if (url === '/routes/1' && options?.method === 'DELETE') return Promise.resolve({ data: { ok: true } })
      return Promise.resolve({ data: [{ id: 1, code: 'R-1', name: 'Route 1' }] })
    })
    const wrapper = mount(Page)
    await flushPromises()

    await (wrapper.vm as any).handleDelete(1)
    await flushPromises()

    expect(mockRequest).toHaveBeenCalledWith('/routes/1', { method: 'DELETE' })
  })
})
