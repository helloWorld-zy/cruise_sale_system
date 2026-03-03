import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/routes/new.vue'

const mockFetch = vi.fn().mockResolvedValue({})
const mockNavigate = vi.fn().mockResolvedValue(undefined)
vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: '/api/v1' } }))
vi.stubGlobal('$fetch', mockFetch)
vi.stubGlobal('navigateTo', mockNavigate)

beforeEach(() => {
  mockFetch.mockClear()
  mockNavigate.mockClear()
})

describe('RoutesNewPage', () => {
  it('提交表单会调用创建接口并跳转', async () => {
    const wrapper = mount(Page)
    const inputs = wrapper.findAll('input')
    await inputs[0]!.setValue('R-100')
    await inputs[1]!.setValue('Ocean')
    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(mockFetch).toHaveBeenCalledWith('/api/v1/admin/routes', expect.objectContaining({ method: 'POST' }))
    expect(mockNavigate).toHaveBeenCalledWith('/routes')
  })
})
