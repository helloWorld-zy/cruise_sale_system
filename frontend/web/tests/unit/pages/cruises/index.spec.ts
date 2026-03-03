import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../../app/pages/cruises/index.vue'

const mockRequest = vi.fn()
vi.stubGlobal('useApi', () => ({ request: mockRequest }))

beforeEach(() => {
  mockRequest.mockReset()
  mockRequest.mockResolvedValue({ data: { list: [{ id: 1, name: 'Ocean Nova', tonnage: 120000, passenger_capacity: 2800, length: 290 }] } })
})

describe('Web cruises index page', () => {
  it('renders hero search card and cruise grid', async () => {
    const wrapper = mount(Page, { global: { stubs: { NuxtLink: { template: '<a><slot /></a>' } } } })
    await flushPromises()
    expect(wrapper.find('[data-test="search-card"]').exists()).toBe(true)
    expect(wrapper.find('[data-test="cruise-grid"]').exists()).toBe(true)
    expect(wrapper.text()).toContain('Ocean Nova')
  })

  it('sends trimmed keyword when searching', async () => {
    const wrapper = mount(Page, { global: { stubs: { NuxtLink: { template: '<a><slot /></a>' } } } })
    await flushPromises()

    await wrapper.find('input[placeholder="搜索邮轮名称/代码"]').setValue('  Ocean  ')
    await wrapper.find('button').trigger('click')
    await flushPromises()

    expect(mockRequest).toHaveBeenLastCalledWith('/cruises', {
      query: expect.objectContaining({ keyword: 'Ocean' }),
    })
  })

  it('renders error text when api fails', async () => {
    mockRequest.mockRejectedValueOnce(new Error('cruises load failed'))
    const wrapper = mount(Page, { global: { stubs: { NuxtLink: { template: '<a><slot /></a>' } } } })
    await flushPromises()

    expect(wrapper.text()).toContain('cruises load failed')
  })
})
