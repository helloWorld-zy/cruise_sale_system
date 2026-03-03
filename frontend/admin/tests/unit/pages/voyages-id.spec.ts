import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/voyages/[id].vue'

const mockRequest = vi.fn()
const mockNavigateTo = vi.fn().mockResolvedValue(undefined)
const confirmMock = vi.fn(() => true)
const routeMock = { params: { id: '13' } }

vi.stubGlobal('useApi', () => ({ request: mockRequest }))
vi.stubGlobal('navigateTo', mockNavigateTo)
vi.stubGlobal('confirm', confirmMock)
vi.stubGlobal('useRoute', () => routeMock)

describe('Voyages edit page', () => {
  beforeEach(() => {
    mockRequest.mockReset()
    mockNavigateTo.mockClear()
    confirmMock.mockClear()
    confirmMock.mockReturnValue(true)

    mockRequest.mockImplementation((url: string, options?: any) => {
      if (url === '/voyages/13' && !options) {
        return Promise.resolve({
          data: {
            code: 'V13',
            name: '日韩春季线',
            route_id: 5,
            departure_date: '2026-05-01T12:00:00Z',
            arrival_date: '2026-05-06T12:00:00Z',
          },
        })
      }
      if (url === '/voyages/13' && options?.method === 'PUT') return Promise.resolve({ data: { ok: true } })
      if (url === '/voyages/13' && options?.method === 'DELETE') return Promise.resolve({ data: { ok: true } })
      return Promise.resolve({ data: {} })
    })
  })

  it('loads detail and submits update', async () => {
    const wrapper = mount(Page)
    await flushPromises()

    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(mockRequest).toHaveBeenCalledWith('/voyages/13', expect.objectContaining({ method: 'PUT' }))
    expect(mockNavigateTo).toHaveBeenCalledWith('/voyages')
  })

  it('deletes and navigates when confirmed', async () => {
    const wrapper = mount(Page)
    await flushPromises()

    const deleteBtn = wrapper.findAll('button').find((b) => b.text().includes('删除'))
    await deleteBtn!.trigger('click')
    await flushPromises()

    expect(confirmMock).toHaveBeenCalled()
    expect(mockRequest).toHaveBeenCalledWith('/voyages/13', { method: 'DELETE' })
    expect(mockNavigateTo).toHaveBeenCalledWith('/voyages')
  })
})
