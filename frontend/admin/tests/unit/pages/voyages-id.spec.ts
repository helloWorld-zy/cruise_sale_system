import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/voyages/[id].vue'

const mockRequest = vi.fn()
const mockNavigateTo = vi.fn().mockResolvedValue(undefined)
const routeMock = { params: { id: '13' } }

vi.stubGlobal('useApi', () => ({ request: mockRequest }))
vi.stubGlobal('navigateTo', mockNavigateTo)
vi.stubGlobal('useRoute', () => routeMock)

describe('Voyages edit page', () => {
  beforeEach(() => {
    mockRequest.mockReset()
    mockNavigateTo.mockClear()

    mockRequest.mockImplementation((url: string, options?: any) => {
      if (url === '/cruises' && options?.query) {
        return Promise.resolve({ data: { list: [{ id: 6, name: 'Ocean Nova' }] } })
      }
      if (url === '/voyages/13' && !options) {
        return Promise.resolve({
          data: {
            cruise_id: 6,
            code: 'V13',
            image_url: 'http://127.0.0.1:8080/uploads/v13.jpg',
            brief_info: '日韩春季线 5晚6天',
            depart_date: '2026-05-01T00:00:00Z',
            return_date: '2026-05-06T00:00:00Z',
            itineraries: [
              { day_no: 1, stop_index: 1, city: '天津' },
            ],
          },
        })
      }
      if (url === '/voyages/13' && options?.method === 'PUT') return Promise.resolve({ data: { ok: true } })
      if (url === '/voyages/13' && options?.method === 'DELETE') return Promise.resolve({ data: { ok: true } })
      return Promise.resolve({ data: {} })
    })
  })

  it('邮轮列表为空时显示提示并禁用保存', async () => {
    mockRequest.mockImplementation((url: string, options?: any) => {
      if (url === '/cruises' && options?.query) {
        return Promise.resolve({ data: { list: [] } })
      }
      if (url === '/voyages/13' && !options) {
        return Promise.resolve({
          data: {
            cruise_id: 6,
            code: 'V13',
            image_url: 'http://127.0.0.1:8080/uploads/v13.jpg',
            brief_info: '日韩春季线 5晚6天',
            depart_date: '2026-05-01T00:00:00Z',
            return_date: '2026-05-06T00:00:00Z',
            itineraries: [{ day_no: 1, stop_index: 1, city: '天津' }],
          },
        })
      }
      return Promise.resolve({ data: {} })
    })

    const wrapper = mount(Page)
    await flushPromises()

    expect(wrapper.find('[data-test="cruise-empty-hint"]').exists()).toBe(true)
    const saveBtn = wrapper.find('button[type="submit"]')
    expect(saveBtn.attributes('disabled')).toBeDefined()
  })

  it('loads detail and submits update', async () => {
    const wrapper = mount(Page)
    await flushPromises()

    expect(wrapper.find('input[placeholder="所属邮轮 ID"]').exists()).toBe(false)
    const cruiseSelect = wrapper.find('[data-test="cruise-select"]')
    expect(cruiseSelect.exists()).toBe(true)
    expect((cruiseSelect.element as HTMLSelectElement).value).toBe('6')

    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(mockRequest).toHaveBeenCalledWith(
      '/voyages/13',
      expect.objectContaining({
        method: 'PUT',
        body: expect.objectContaining({
          image_url: 'http://127.0.0.1:8080/uploads/v13.jpg',
        }),
      }),
    )
    expect(mockNavigateTo).toHaveBeenCalledWith('/voyages')
  })

  it('does not silently fallback to the first cruise when detail misses cruise_id', async () => {
    mockRequest.mockImplementation((url: string, options?: any) => {
      if (url === '/cruises' && options?.query) {
        return Promise.resolve({ data: { list: [{ id: 6, name: 'Ocean Nova' }] } })
      }
      if (url === '/voyages/13' && !options) {
        return Promise.resolve({
          data: {
            code: 'V13',
            brief_info: '日韩春季线 5晚6天',
            depart_date: '2026-05-01T00:00:00Z',
            return_date: '2026-05-06T00:00:00Z',
            itineraries: [{ day_no: 1, stop_index: 1, city: '天津' }],
          },
        })
      }
      return Promise.resolve({ data: {} })
    })

    const wrapper = mount(Page)
    await flushPromises()

    const cruiseSelect = wrapper.find('[data-test="cruise-select"]')
    expect((cruiseSelect.element as HTMLSelectElement).value).toBe('0')
    expect(wrapper.find('button[type="submit"]').attributes('disabled')).toBeDefined()
  })

  it('deletes and navigates when confirmed', async () => {
    const wrapper = mount(Page, { attachTo: document.body })
    await flushPromises()

    const deleteBtn = wrapper.findAll('button').find((b) => b.text().trim() === '删除')
    expect(deleteBtn).toBeTruthy()
    await deleteBtn!.trigger('click')
    await flushPromises()

    const confirmBtn = Array.from(document.body.querySelectorAll('button')).find((btn) =>
      (btn.textContent || '').includes('确认删除'),
    )
    expect(confirmBtn).toBeTruthy()
    ;(confirmBtn as HTMLButtonElement).click()
    await flushPromises()

    expect(mockRequest).toHaveBeenCalledWith('/voyages/13', { method: 'DELETE' })
    expect(mockNavigateTo).toHaveBeenCalledWith('/voyages')
    wrapper.unmount()
  })
})
