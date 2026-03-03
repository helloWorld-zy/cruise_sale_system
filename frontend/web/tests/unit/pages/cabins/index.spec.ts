import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../../app/pages/cabins/index.vue'

const mockRequest = vi.fn()
vi.stubGlobal('useApi', () => ({ request: mockRequest }))

beforeEach(() => {
  mockRequest.mockReset()
  mockRequest.mockResolvedValue({ data: [] })
})

describe('Cabins Index Page', () => {
  it('renders page title', async () => {
    const wrapper = mount(Page, { global: { stubs: { NuxtLink: { template: '<a><slot /></a>' } } } })
    await flushPromises()
    expect(wrapper.text()).toContain('舱位浏览')
  })

  it('renders filter dropdowns', async () => {
    const wrapper = mount(Page, { global: { stubs: { NuxtLink: { template: '<a><slot /></a>' } } } })
    await flushPromises()
    expect(wrapper.text()).toContain('选择航线')
    expect(wrapper.text()).toContain('选择舱型')
    expect(wrapper.text()).toContain('默认排序')
  })

  it('renders search button', async () => {
    const wrapper = mount(Page, { global: { stubs: { NuxtLink: { template: '<a><slot /></a>' } } } })
    await flushPromises()
    expect(wrapper.text()).toContain('搜索')
  })

  it('renders empty state when no cabins', async () => {
    const wrapper = mount(Page, { global: { stubs: { NuxtLink: { template: '<a><slot /></a>' } } } })
    await flushPromises()
    expect(wrapper.text()).toContain('暂无可用舱位')
  })
})
