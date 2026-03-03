import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/voyages/index.vue'

const mockRequest = vi.fn()
const confirmMock = vi.fn(() => true)
vi.stubGlobal('useApi', () => ({ request: mockRequest }))
vi.stubGlobal('confirm', confirmMock)

beforeEach(() => {
  mockRequest.mockReset()
  confirmMock.mockClear()
  confirmMock.mockReturnValue(true)
  mockRequest.mockResolvedValue({ data: [{ id: 1, code: 'V-1', name: 'Voyage', route_id: 10, departure_date: '2026-03-01' }] })
})

describe('VoyagesIndexPage', () => {
  it('调用 voyages API', async () => {
    mount(Page)
    await flushPromises()
    expect(mockRequest).toHaveBeenCalledWith('/voyages')
  })

  it('渲染航次行', async () => {
    const wrapper = mount(Page)
    await flushPromises()
    expect(wrapper.text()).toContain('V-1')
    expect(wrapper.text()).toContain('2026-03-01')
  })

  it('空数据时显示 No data', async () => {
    mockRequest.mockResolvedValueOnce({ data: [] })
    const wrapper = mount(Page)
    await flushPromises()
    expect(wrapper.text()).toContain('No data')
  })

  it('删除无效 id 时显示错误', async () => {
    const wrapper = mount(Page)
    await flushPromises()
    await (wrapper.vm as any).handleDelete('bad')
    await flushPromises()
    expect(wrapper.text()).toContain('无效记录 ID，无法删除')
  })

  it('删除确认取消时不调用 DELETE', async () => {
    confirmMock.mockReturnValueOnce(false)
    const wrapper = mount(Page)
    await flushPromises()
    await (wrapper.vm as any).handleDelete(1)
    await flushPromises()
    expect(mockRequest).not.toHaveBeenCalledWith('/voyages/1', { method: 'DELETE' })
  })

  it('删除成功后调用 DELETE', async () => {
    mockRequest.mockImplementation((url: string, options?: any) => {
      if (url === '/voyages/1' && options?.method === 'DELETE') return Promise.resolve({ data: { ok: true } })
      return Promise.resolve({ data: [{ id: 1, code: 'V-1', name: 'Voyage', route_id: 10, departure_date: '2026-03-01' }] })
    })
    const wrapper = mount(Page)
    await flushPromises()
    await (wrapper.vm as any).handleDelete(1)
    await flushPromises()
    expect(mockRequest).toHaveBeenCalledWith('/voyages/1', { method: 'DELETE' })
  })
})
