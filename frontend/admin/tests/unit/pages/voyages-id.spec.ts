import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/voyages/[id].vue'

async function waitForCitySearch() {
  await new Promise((resolve) => setTimeout(resolve, 380))
  await flushPromises()
}

const mockRequest = vi.fn()
const mockFetch = vi.fn()
const mockNavigateTo = vi.fn().mockResolvedValue(undefined)
const routeMock = { params: { id: '13' } }

vi.stubGlobal('useApi', () => ({ request: mockRequest }))
vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: '/api/v1' } }))
vi.stubGlobal('$fetch', mockFetch)
vi.stubGlobal('navigateTo', mockNavigateTo)
vi.stubGlobal('useRoute', () => routeMock)

describe('Voyages edit page', () => {
  beforeEach(() => {
    mockRequest.mockReset()
    mockFetch.mockReset()
    mockNavigateTo.mockClear()

    mockFetch.mockImplementation((url: string, options?: any) => {
      if (url === '/api/v1/admin/port-cities' && options?.query?.keyword === '仁川') {
        return Promise.resolve({ data: [{ label: '仁川（韩国）', city_name: '仁川', country_name: '韩国' }] })
      }
      return Promise.resolve({ data: [] })
    })

    mockRequest.mockImplementation((url: string, options?: any) => {
      if (url === '/cruises' && options?.query) {
        return Promise.resolve({ data: { list: [{ id: 6, name: 'Ocean Nova' }] } })
      }
      if (url === '/content-templates' && !options) {
        return Promise.resolve({ data: [
          { id: 8, name: '默认费用说明', kind: 'fee_note', status: 1, content: { included: [{ text: '船票' }, { text: '港务费' }], excluded: [{ text: '签证费' }] } },
          { id: 9, name: '默认预订须知', kind: 'booking_notice', status: 1, content: { sections: [{ key: 'documents', title: '出行证件', items: [{ text: '请携带护照' }] }, { key: 'cancel', title: '退改规则', items: [{ text: '开航前 7 天内不可退' }] }] } },
        ] })
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
            fee_note_template_id: 8,
            fee_note_mode: 'template',
            booking_notice_template_id: 9,
            booking_notice_mode: 'snapshot',
            booking_notice_content: { sections: [{ key: 'documents', title: '出行证件', items: [{ text: '请携带护照' }] }] },
            itineraries: [
              { day_no: 1, stop_index: 1, city: '天津（中国）' },
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
  expect(wrapper.find('input[placeholder="城市（必填）"]').exists()).toBe(false)
    const cruiseSelect = wrapper.find('[data-test="cruise-select"]')
    expect(cruiseSelect.exists()).toBe(true)
    expect((cruiseSelect.element as HTMLSelectElement).value).toBe('6')
    expect((wrapper.find('[data-test="fee-note-template-select"]').element as HTMLSelectElement).value).toBe('8')
    expect(wrapper.find('[data-test="copy-fee-note-to-snapshot"]').exists()).toBe(true)
    expect(wrapper.find('[data-test="booking-notice-text-0"]').exists()).toBe(true)

    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(mockRequest).toHaveBeenCalledWith(
      '/voyages/13',
      expect.objectContaining({
        method: 'PUT',
        body: expect.objectContaining({
          image_url: 'http://127.0.0.1:8080/uploads/v13.jpg',
          fee_note_template_id: 8,
          fee_note_mode: 'template',
          booking_notice_template_id: 9,
          booking_notice_mode: 'snapshot',
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

    const deleteBtn = wrapper.find('[data-test="delete-voyage"]')
    expect(deleteBtn.exists()).toBe(true)
    await deleteBtn.trigger('click')
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

  it('复制模板后支持多条费用说明和多分组预订须知快照编辑', async () => {
    const wrapper = mount(Page)
    await flushPromises()

    await wrapper.find('[data-test="copy-fee-note-to-snapshot"]').trigger('click')
    expect(wrapper.find('[data-test="fee-note-included-1"]').exists()).toBe(true)

    await wrapper.find('[data-test="copy-booking-notice-to-snapshot"]').trigger('click')
    expect(wrapper.find('[data-test="booking-notice-section-title-1"]').exists()).toBe(true)
  })

  it('编辑页支持搜索并选择城市候选', async () => {
    const wrapper = mount(Page)
    await flushPromises()

    await wrapper.find('[data-test="itinerary-city-1-1-input"]').setValue('仁川')
    await waitForCitySearch()
    await wrapper.find('[data-test="itinerary-city-1-1-option-0"]').trigger('click')
    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(mockRequest).toHaveBeenCalledWith(
      '/voyages/13',
      expect.objectContaining({
        method: 'PUT',
        body: expect.objectContaining({
          itineraries: expect.arrayContaining([
            expect.objectContaining({ city: '仁川（韩国）' }),
          ]),
        }),
      }),
    )
  })
})
