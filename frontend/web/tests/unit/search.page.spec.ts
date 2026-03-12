import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../app/pages/search/index.vue'

const mockRequest = vi.fn()
vi.stubGlobal('useApi', () => ({ request: mockRequest }))
vi.stubGlobal('useRoute', () => ({ query: {} }))

// Mock vue-router for useRouter
vi.mock('vue-router', () => ({
  useRouter: () => ({ push: vi.fn() }),
  useRoute: () => ({ query: {} })
}))

beforeEach(() => {
  mockRequest.mockReset()
  mockRequest.mockResolvedValue({ data: [] })
})

describe('Search Page', () => {
  it('renders search form with filters', async () => {
    const wrapper = mount(Page, { global: { stubs: { NuxtLink: { template: '<a><slot /></a>' } } } })
    await flushPromises()
    expect(wrapper.text()).toContain('探寻下一个航海传奇')
    expect(wrapper.text()).toContain('目的地 Where')
    expect(wrapper.text()).toContain('出发地 Departure')
    expect(wrapper.text()).toContain('日期 Dates')
  })

  it('renders search button', async () => {
    const wrapper = mount(Page, { global: { stubs: { NuxtLink: { template: '<a><slot /></a>' } } } })
    await flushPromises()
    expect(wrapper.text()).toContain('查找航线')
  })

  it('renders empty state when no results', async () => {
    const wrapper = mount(Page, { global: { stubs: { NuxtLink: { template: '<a><slot /></a>' } } } })
    await flushPromises()
    expect(wrapper.text()).toContain('未找到匹配航线')
  })
})
