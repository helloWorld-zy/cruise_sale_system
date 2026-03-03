import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/cabins/[id].vue'

const mockRequest = vi.fn()
const mockNavigateTo = vi.fn().mockResolvedValue(undefined)
const confirmMock = vi.fn(() => true)
const routeMock = { params: { id: '15' } }

vi.stubGlobal('useApi', () => ({ request: mockRequest }))
vi.stubGlobal('navigateTo', mockNavigateTo)
vi.stubGlobal('confirm', confirmMock)
vi.stubGlobal('useRoute', () => routeMock)

describe('Cabins edit page', () => {
  beforeEach(() => {
    mockRequest.mockReset()
    mockNavigateTo.mockClear()
    confirmMock.mockClear()
    confirmMock.mockReturnValue(true)

    mockRequest.mockImplementation((url: string, options?: any) => {
      if (url === '/cabins/15' && !options) {
        return Promise.resolve({
          data: {
            voyage_id: 8,
            cabin_type_id: 3,
            code: 'A305',
            deck: '11F',
            area: 26,
            max_guests: 3,
            position: 'mid',
            orientation: 'port',
            has_window: true,
            has_balcony: false,
            bed_type: 'queen',
            amenities: '浴缸,智能电视',
            status: 1,
          },
        })
      }
      if (url === '/cabins/15' && options?.method === 'PUT') return Promise.resolve({ data: { ok: true } })
      if (url === '/cabins/15' && options?.method === 'DELETE') return Promise.resolve({ data: { ok: true } })
      return Promise.resolve({ data: {} })
    })
  })

  it('loads detail and submits update', async () => {
    const wrapper = mount(Page)
    await flushPromises()

    expect(wrapper.text()).toContain('编辑舱位 #15')
    expect(wrapper.text()).toContain('浴缸')

    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(mockRequest).toHaveBeenCalledWith('/cabins/15', expect.objectContaining({ method: 'PUT' }))
    expect(mockNavigateTo).toHaveBeenCalledWith('/cabins')
  })

  it('does not delete when user cancels confirm', async () => {
    confirmMock.mockReturnValueOnce(false)
    const wrapper = mount(Page)
    await flushPromises()

    const deleteBtn = wrapper.findAll('button').find((b) => b.text().includes('删除'))
    await deleteBtn!.trigger('click')
    await flushPromises()

    expect(mockRequest).not.toHaveBeenCalledWith('/cabins/15', { method: 'DELETE' })
  })
})
