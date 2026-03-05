import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/companies/[id].vue'

const mockRequest = vi.fn()
const mockNavigateTo = vi.fn().mockResolvedValue(undefined)

vi.stubGlobal('useApi', () => ({ request: mockRequest }))
vi.stubGlobal('navigateTo', mockNavigateTo)
vi.stubGlobal('confirm', vi.fn(() => true))
vi.stubGlobal('useRoute', () => ({ params: { id: '1' } }))

describe('Companies edit page', () => {
  const globalStubs = {
    NuxtLink: { template: '<a><slot /></a>' },
    AdminActionLink: { template: '<a class="admin-btn"><slot /></a>' },
  }

  beforeEach(() => {
    mockRequest.mockReset()
    mockNavigateTo.mockClear()
    mockRequest.mockImplementation((url: string, options?: any) => {
      if (url === '/companies' && !options) {
        return Promise.resolve({ data: { list: [{ id: 1, name: '皇家加勒比', english_name: 'Royal Caribbean', logo_url: 'https://img/logo.png', description: 'desc' }] } })
      }
      if (url === '/companies/1' && options?.method === 'PUT') {
        return Promise.resolve({ data: { ok: true } })
      }
      return Promise.resolve({ data: {} })
    })
  })

  it('loads company detail and updates', async () => {
    const wrapper = mount(Page, { global: { stubs: globalStubs } })
    await flushPromises()

    expect(wrapper.text()).toContain('编辑邮轮公司')
    expect(wrapper.find('.company-logo-preview-frame').exists()).toBe(true)
    expect(wrapper.find('.company-logo-preview-image').exists()).toBe(true)
    const cancelAction = wrapper.findAll('a').find((el) => el.text().trim() === '取消')
    expect(cancelAction?.classes()).toContain('admin-btn')
    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(mockRequest).toHaveBeenCalledWith('/companies/1', expect.objectContaining({ method: 'PUT' }))
    expect(mockNavigateTo).toHaveBeenCalledWith('/companies')
  })
})
