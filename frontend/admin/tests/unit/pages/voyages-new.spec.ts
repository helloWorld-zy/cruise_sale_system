import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/voyages/new.vue'

async function waitForCitySearch() {
  await new Promise((resolve) => setTimeout(resolve, 380))
  await flushPromises()
}

const mockFetch = vi.fn()
const mockNavigate = vi.fn().mockResolvedValue(undefined)
vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: '/api/v1' } }))
vi.stubGlobal('$fetch', mockFetch)
vi.stubGlobal('navigateTo', mockNavigate)

beforeEach(() => {
  mockFetch.mockClear()
  mockNavigate.mockClear()
  mockFetch.mockImplementation((url: string, options?: any) => {
    if (url === '/api/v1/admin/cruises' && options?.query) {
      return Promise.resolve({ data: { list: [{ id: 12, name: 'Ocean Nova' }] } })
    }
    if (url === '/api/v1/admin/content-templates') {
      return Promise.resolve({ data: [
        { id: 8, name: '默认费用说明', kind: 'fee_note', status: 1, content: { included: [{ text: '船票' }, { text: '港务费' }], excluded: [{ text: '签证费' }] } },
        { id: 9, name: '默认预订须知', kind: 'booking_notice', status: 1, content: { sections: [{ key: 'documents', title: '出行证件', items: [{ text: '请携带护照' }] }, { key: 'cancel', title: '退改规则', items: [{ text: '开航前 7 天内不可退' }] }] } },
      ] })
    }
    if (url === '/api/v1/admin/port-cities' && options?.query?.keyword === '仁川') {
      return Promise.resolve({ data: [{ label: '仁川（韩国）', city_name: '仁川', country_name: '韩国' }] })
    }
    if (url === '/api/v1/admin/port-cities' && options?.query?.keyword === '布') {
      return Promise.resolve({ data: [{ label: '布宜诺斯艾利斯（阿根廷）', city_name: '布宜诺斯艾利斯', country_name: '阿根廷' }] })
    }
    if (url === '/api/v1/admin/port-cities' && options?.query?.keyword === '海上') {
      return Promise.resolve({ data: [{ label: '海上巡游', is_special: true }] })
    }
    if (url === '/api/v1/admin/voyages' && options?.method === 'POST') {
      return Promise.resolve({ data: { ok: true } })
    }
    if (url === '/api/v1/admin/upload/image' && options?.method === 'POST') {
      return Promise.resolve({ data: { url: 'http://127.0.0.1:8080/uploads/voyage.jpg' } })
    }
    return Promise.resolve({})
  })
})

describe('VoyagesNewPage', () => {
  it('邮轮列表为空时显示提示并禁用提交', async () => {
    mockFetch.mockImplementation((url: string, options?: any) => {
      if (url === '/api/v1/admin/cruises' && options?.query) {
        return Promise.resolve({ data: { list: [] } })
      }
      return Promise.resolve({})
    })

    const wrapper = mount(Page)
    await flushPromises()

    expect(wrapper.find('[data-test="cruise-empty-hint"]').exists()).toBe(true)
    const submitBtn = wrapper.find('button[type="submit"]')
    expect(submitBtn.attributes('disabled')).toBeDefined()

    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()
    expect(mockFetch).not.toHaveBeenCalledWith('/api/v1/admin/voyages', expect.anything())
  })

  it('不展示手工输入的所属邮轮 ID 字段，展示邮轮下拉绑定', async () => {
    const wrapper = mount(Page)
    await flushPromises()
    expect(wrapper.text()).not.toContain('所属邮轮 ID')
    expect(wrapper.find('input[placeholder="所属邮轮 ID"]').exists()).toBe(false)
    const cruiseSelect = wrapper.find('[data-test="cruise-select"]')
    expect(cruiseSelect.exists()).toBe(true)
    expect((cruiseSelect.element as HTMLSelectElement).value).toBe('12')
  })

  it('提交表单触发创建并跳转', async () => {
    const wrapper = mount(Page)
    await flushPromises()

    await wrapper.find('input[placeholder="航次代码（如 TJ-20260701）"]').setValue('V-100')
    await wrapper.find('input[placeholder="航次简介（手动输入）"]').setValue('天津-济州-釜山-天津 5晚6天')

    const dateInputs = wrapper.findAll('input[type="date"]')
    await dateInputs[0]!.setValue('2026-05-01')
    await dateInputs[1]!.setValue('2026-05-10')

    await wrapper.find('[data-test="itinerary-city-1-1-input"]').setValue('仁川')
    await waitForCitySearch()
    await wrapper.find('[data-test="itinerary-city-1-1-option-0"]').trigger('click')
    await wrapper.find('[data-test="fee-note-template-select"]').setValue('8')
    await wrapper.find('[data-test="booking-notice-template-select"]').setValue('9')
    await wrapper.find('[data-test="copy-booking-notice-to-snapshot"]').trigger('click')
    await wrapper.find('[data-test="booking-notice-text-0"]').setValue('请至少提前 1 天到港')

    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(mockFetch).toHaveBeenCalledWith(
      '/api/v1/admin/voyages',
      expect.objectContaining({
        method: 'POST',
        body: expect.objectContaining({
          image_url: '',
          fee_note_template_id: 8,
          fee_note_mode: 'template',
          booking_notice_template_id: 9,
          booking_notice_mode: 'snapshot',
          booking_notice_content: expect.objectContaining({
            sections: expect.arrayContaining([
              expect.objectContaining({ items: expect.arrayContaining([{ text: '请至少提前 1 天到港' }]) }),
            ]),
          }),
          itineraries: expect.arrayContaining([
            expect.objectContaining({ city: '仁川（韩国）' }),
          ]),
        }),
      }),
    )
    expect(mockNavigate).toHaveBeenCalledWith('/voyages')
  })

  it('支持将海上巡游作为特殊站点从搜索结果中选择', async () => {
    const wrapper = mount(Page)
    await flushPromises()

    await wrapper.find('[data-test="itinerary-city-1-1-input"]').setValue('海上')
    await waitForCitySearch()
    await wrapper.find('[data-test="itinerary-city-1-1-option-0"]').trigger('click')

    expect((wrapper.find('[data-test="itinerary-city-1-1-input"]').element as HTMLInputElement).value).toBe('海上巡游')
  })

  it('输入首字时就会发起城市搜索并显示结果', async () => {
    const wrapper = mount(Page)
    await flushPromises()

    await wrapper.find('[data-test="itinerary-city-1-1-input"]').setValue('布')
    await waitForCitySearch()

    expect(mockFetch).toHaveBeenCalledWith(
      '/api/v1/admin/port-cities',
      expect.objectContaining({ query: { keyword: '布' } }),
    )
    expect(wrapper.find('[data-test="itinerary-city-1-1-option-0"]').text()).toContain('布宜诺斯艾利斯（阿根廷）')
  })

	it('显示模板选择器和复制为航次专属内容入口', async () => {
		const wrapper = mount(Page)
		await flushPromises()
		expect(wrapper.find('[data-test="fee-note-template-select"]').exists()).toBe(true)
		expect(wrapper.find('[data-test="booking-notice-template-select"]').exists()).toBe(true)
    expect(wrapper.find('[data-test="copy-fee-note-to-snapshot"]').exists()).toBe(true)
		expect(wrapper.text()).toContain('复制为航次专属内容')
	})

  it('复制模板后支持多条费用说明和多分组预订须知快照编辑', async () => {
    const wrapper = mount(Page)
    await flushPromises()

    await wrapper.find('[data-test="fee-note-template-select"]').setValue('8')
    await wrapper.find('[data-test="copy-fee-note-to-snapshot"]').trigger('click')
    expect(wrapper.find('[data-test="fee-note-included-1"]').exists()).toBe(true)

    await wrapper.find('[data-test="booking-notice-template-select"]').setValue('9')
    await wrapper.find('[data-test="copy-booking-notice-to-snapshot"]').trigger('click')
    expect(wrapper.find('[data-test="booking-notice-section-title-1"]').exists()).toBe(true)
  })
})
