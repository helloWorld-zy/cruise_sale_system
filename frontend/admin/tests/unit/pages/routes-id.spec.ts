import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/routes/[id].vue'

const mockRequest = vi.fn()
const mockNavigateTo = vi.fn().mockResolvedValue(undefined)
const confirmMock = vi.fn(() => true)
const routeMock = { params: { id: '7' } }

vi.stubGlobal('useApi', () => ({ request: mockRequest }))
vi.stubGlobal('navigateTo', mockNavigateTo)
vi.stubGlobal('confirm', confirmMock)
vi.stubGlobal('useRoute', () => routeMock)

describe('Routes edit page', () => {
  beforeEach(() => {
    mockRequest.mockReset()
    mockNavigateTo.mockClear()
    confirmMock.mockClear()
    confirmMock.mockReturnValue(true)

    mockRequest.mockImplementation((url: string, options?: any) => {
      if (url === '/routes/7' && !options) return Promise.resolve({ data: { code: 'R7', name: '北欧航线' } })
      if (url === '/routes/7' && options?.method === 'PUT') return Promise.resolve({ data: { ok: true } })
      if (url === '/routes/7' && options?.method === 'DELETE') return Promise.resolve({ data: { ok: true } })
      return Promise.resolve({ data: {} })
    })
  })

  it('loads detail and saves', async () => {
    const wrapper = mount(Page)
    await flushPromises()

    expect(wrapper.text()).toContain('编辑航线 #7')
    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(mockRequest).toHaveBeenCalledWith('/routes/7', expect.objectContaining({ method: 'PUT' }))
    expect(mockNavigateTo).toHaveBeenCalledWith('/routes')
  })

  it('does not delete when confirm is cancelled', async () => {
    confirmMock.mockReturnValueOnce(false)
    const wrapper = mount(Page)
    await flushPromises()

    const deleteBtn = wrapper.findAll('button').find((b) => b.text().includes('删除'))
    await deleteBtn!.trigger('click')
    await flushPromises()

    expect(mockRequest).not.toHaveBeenCalledWith('/routes/7', { method: 'DELETE' })
  })
})
