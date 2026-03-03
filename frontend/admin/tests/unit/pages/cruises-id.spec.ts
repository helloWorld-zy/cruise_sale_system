import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/cruises/[id].vue'

const mockRequest = vi.fn()
const mockNavigateTo = vi.fn().mockResolvedValue(undefined)
const routeMock = { params: { id: '5' } }
const confirmMock = vi.fn(() => true)

vi.stubGlobal('useApi', () => ({ request: mockRequest }))
vi.stubGlobal('navigateTo', mockNavigateTo)
vi.stubGlobal('useRoute', () => routeMock)
vi.stubGlobal('confirm', confirmMock)

describe('Cruise edit page', () => {
  beforeEach(() => {
    mockRequest.mockReset()
    mockNavigateTo.mockClear()
    confirmMock.mockClear()
    confirmMock.mockReturnValue(true)

    mockRequest.mockImplementation((url: string, options?: any) => {
      if (url === '/cruises/5' && !options) {
        return Promise.resolve({
          data: {
            id: 5,
            name: 'Ocean Nova',
            code: 'ONOVA',
            company_id: 1,
            tonnage: 90000,
            passenger_capacity: 2400,
            crew_count: 900,
            build_year: 2018,
            refurbish_year: 2024,
            length: 310,
            width: 38,
            deck_count: 14,
            status: 1,
          },
        })
      }
      if (url === '/cruises/5' && options?.method === 'PUT') {
        return Promise.resolve({ data: { ok: true } })
      }
      if (url === '/cruises/5' && options?.method === 'DELETE') {
        return Promise.resolve({ data: { ok: true } })
      }
      return Promise.resolve({ data: {} })
    })
  })

  it('loads detail on mounted', async () => {
    const wrapper = mount(Page)
    await flushPromises()

    expect(wrapper.text()).toContain('编辑邮轮 #5')
    expect(mockRequest).toHaveBeenCalledWith('/cruises/5')
    expect((wrapper.find('input').element as HTMLInputElement).value).toContain('Ocean Nova')
  })

  it('saves form then navigates', async () => {
    const wrapper = mount(Page)
    await flushPromises()

    const form = wrapper.find('form')
    await form.trigger('submit.prevent')
    await flushPromises()

    expect(mockRequest).toHaveBeenCalledWith('/cruises/5', expect.objectContaining({ method: 'PUT' }))
    expect(mockNavigateTo).toHaveBeenCalledWith('/cruises')
  })

  it('deletes after confirm then navigates', async () => {
    const wrapper = mount(Page)
    await flushPromises()

    const deleteButton = wrapper.findAll('button').find((btn) => btn.text().includes('删除'))
    await deleteButton!.trigger('click')
    await flushPromises()

    expect(confirmMock).toHaveBeenCalled()
    expect(mockRequest).toHaveBeenCalledWith('/cruises/5', { method: 'DELETE' })
    expect(mockNavigateTo).toHaveBeenCalledWith('/cruises')
  })

  it('does not delete when confirm is cancelled', async () => {
    confirmMock.mockReturnValueOnce(false)
    const wrapper = mount(Page)
    await flushPromises()

    const deleteButton = wrapper.findAll('button').find((btn) => btn.text().includes('删除'))
    await deleteButton!.trigger('click')
    await flushPromises()

    expect(mockRequest).not.toHaveBeenCalledWith('/cruises/5', { method: 'DELETE' })
  })

  it('shows save error when update request fails', async () => {
    mockRequest.mockImplementation((url: string, options?: any) => {
      if (url === '/cruises/5' && !options) {
        return Promise.resolve({ data: { id: 5, name: 'Ocean Nova', code: 'ONOVA' } })
      }
      if (url === '/cruises/5' && options?.method === 'PUT') {
        return Promise.reject(new Error('update failed'))
      }
      return Promise.resolve({ data: {} })
    })

    const wrapper = mount(Page)
    await flushPromises()

    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(wrapper.text()).toContain('update failed')
  })

  it('handles load error', async () => {
    mockRequest.mockRejectedValueOnce(new Error('load failed'))
    const wrapper = mount(Page)
    await flushPromises()

    expect(wrapper.text()).toContain('load failed')
  })
})
