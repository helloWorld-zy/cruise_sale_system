import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/facilities/index.vue'

const mockRequest = vi.fn()
vi.stubGlobal('useApi', () => ({ request: mockRequest }))

beforeEach(() => {
  mockRequest.mockReset()
  mockRequest.mockImplementation((url: string, options?: any) => {
    if (url === '/cruises') {
      return Promise.resolve({ data: [{ id: 1, name: 'Ocean Nova' }] })
    }
    if (url === '/facility-categories') {
      return Promise.resolve({ data: [{ id: 2, name: '娱乐' }] })
    }
    if (url === '/facilities') {
      const query = options?.query || {}
      if (query.cruise_id !== 1) return Promise.resolve({ data: [] })
      return Promise.resolve({ data: [{ id: 9, name: '海上剧院', location: '5层船中', open_hours: '10:00-22:00', extra_charge: true, status: 1 }] })
    }
    return Promise.resolve({ data: [] })
  })
})

describe('Facilities list', () => {
  it('renders cruise/category filters and table', async () => {
    const wrapper = mount(Page, {
      global: { stubs: { NuxtLink: { template: '<a><slot /></a>' } } },
    })
    await flushPromises()
    expect(wrapper.find('[data-test="facility-cruise-filter"]').exists()).toBe(true)
    expect(wrapper.find('table').exists()).toBe(true)
  })

  it('renders charge and status tags from api', async () => {
    const wrapper = mount(Page, {
      global: { stubs: { NuxtLink: { template: '<a><slot /></a>' } } },
    })
    await flushPromises()
    expect(wrapper.text()).toContain('海上剧院')
    expect(wrapper.text()).toContain('收费')
    expect(wrapper.text()).toContain('开放')
  })
})
