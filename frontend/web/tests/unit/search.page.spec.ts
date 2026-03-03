import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../app/pages/search/index.vue'

const mockRequest = vi.fn()
vi.stubGlobal('useApi', () => ({ request: mockRequest }))
vi.stubGlobal('useRoute', () => ({ query: {} }))

beforeEach(() => {
  mockRequest.mockReset()
  mockRequest.mockResolvedValue({ data: [] })
})

describe('Search Page', () => {
  it('renders search form with filters', async () => {
    const wrapper = mount(Page, { global: { stubs: { NuxtLink: { template: '<a><slot /></a>' } } } })
    await flushPromises()
    expect(wrapper.text()).toContain('搜索舱位')
    expect(wrapper.text()).toContain('选择航线')
    expect(wrapper.text()).toContain('选择舱型')
  })

  it('renders search button', async () => {
    const wrapper = mount(Page, { global: { stubs: { NuxtLink: { template: '<a><slot /></a>' } } } })
    await flushPromises()
    expect(wrapper.text()).toContain('搜索')
  })

  it('renders empty state when no results', async () => {
    const wrapper = mount(Page, { global: { stubs: { NuxtLink: { template: '<a><slot /></a>' } } } })
    await flushPromises()
    expect(wrapper.text()).toContain('暂无搜索结果')
  })
})
