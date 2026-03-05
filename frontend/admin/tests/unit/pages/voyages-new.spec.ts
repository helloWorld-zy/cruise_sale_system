import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/voyages/new.vue'

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
    if (url === '/api/v1/admin/voyages' && options?.method === 'POST') {
      return Promise.resolve({ data: { ok: true } })
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
    const inputs = wrapper.findAll('input')
    await inputs[0]!.setValue('V-100')
    await inputs[1]!.setValue('天津-济州-釜山-天津 5晚6天')
    await inputs[2]!.setValue('2026-05-01')
    await inputs[3]!.setValue('2026-05-10')
    await inputs[4]!.setValue('天津')
    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(mockFetch).toHaveBeenCalledWith('/api/v1/admin/voyages', expect.objectContaining({ method: 'POST' }))
    expect(mockNavigate).toHaveBeenCalledWith('/voyages')
  })
})
