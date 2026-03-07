import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/cruises/create.vue'

const mockFetch = vi.fn()
const mockNavigateTo = vi.fn().mockResolvedValue(undefined)

vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: '/api/v1' } }))
vi.stubGlobal('$fetch', mockFetch)
vi.stubGlobal('navigateTo', mockNavigateTo)

describe('Cruises create page', () => {
  beforeEach(() => {
    mockFetch.mockReset()
    mockNavigateTo.mockClear()
  })

  it('creates cruise and navigates back to list', async () => {
    mockFetch.mockImplementation((url: string) => {
      if (url.includes('/companies')) {
        return Promise.resolve({ data: { list: [{ id: 2, name: '皇家加勒比' }] } })
      }
      return Promise.resolve({ data: { id: 5 } })
    })
    const wrapper = mount(Page)
    await flushPromises()

    const inputs = wrapper.findAll('input')
    await inputs[0]!.setValue('Ocean Star')
    await inputs[1]!.setValue('Ocean Star EN')
    await inputs[2]!.setValue('OS-01')
    await wrapper.get('[data-test="company-select-trigger"]').trigger('click')
    await wrapper.get('[data-test="company-option-2"]').trigger('click')
    await inputs[3]!.setValue('90000')

    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(mockFetch).toHaveBeenCalledWith(
      '/api/v1/admin/cruises',
      expect.objectContaining({
        method: 'POST',
        headers: expect.objectContaining({ 'Content-Type': 'application/json' }),
      }),
    )
    expect(mockNavigateTo).toHaveBeenCalledWith('/cruises')
  })

  it('shows create error message when request fails', async () => {
    mockFetch.mockImplementation((url: string, options?: any) => {
      if (url.includes('/companies')) {
        return Promise.resolve({ data: { list: [{ id: 2, name: '皇家加勒比' }] } })
      }
      if (options?.method === 'POST') {
        return Promise.reject(new Error('create cruise failed'))
      }
      return Promise.resolve({ data: {} })
    })
    const wrapper = mount(Page)
    await flushPromises()

    const inputs = wrapper.findAll('input')
    await inputs[0]!.setValue('Ocean Star')
    await wrapper.get('[data-test="company-select-trigger"]').trigger('click')
    await wrapper.get('[data-test="company-option-2"]').trigger('click')

    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(wrapper.text()).toContain('create cruise failed')
  })

  it('blocks submit and shows validation hints when required fields are invalid', async () => {
    mockFetch.mockImplementation((url: string) => {
      if (url.includes('/companies')) {
        return Promise.resolve({ data: { list: [] } })
      }
      return Promise.resolve({ data: {} })
    })

    const wrapper = mount(Page)
    await flushPromises()

    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(wrapper.text()).toContain('请先修正表单校验错误')
    expect(wrapper.text()).toContain('请填写邮轮名称')
    expect(wrapper.text()).toContain('请选择所属公司')
    expect(mockFetch).not.toHaveBeenCalledWith(
      '/api/v1/admin/cruises',
      expect.objectContaining({ method: 'POST' }),
    )
  })
})
