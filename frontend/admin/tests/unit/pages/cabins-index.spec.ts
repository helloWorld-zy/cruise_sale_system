import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/cabins/index.vue'

const mockRequest = vi.fn()
vi.stubGlobal('useApi', () => ({ request: mockRequest }))
vi.stubGlobal('confirm', () => true)

beforeEach(() => {
  mockRequest.mockReset()
  mockRequest.mockImplementation((url: string) => {
    if (url === '/cruises') return Promise.resolve({ data: [{ id: 1, name: 'Ocean Nova' }] })
    if (url === '/routes') return Promise.resolve({ data: [{ id: 2, name: '地中海航线', cruise_id: 1 }] })
    if (url === '/cabin-types') return Promise.resolve({ data: { list: [{ id: 3, name: '阳台房' }] } })
    if (url === '/cabins') return Promise.resolve({ data: { list: [{ id: 10, code: 'A101', area: 24, total: 10, available: 8, status: 1 }], total: 1 } })
    return Promise.resolve({ data: [] })
  })
})

describe('CabinsIndexPage', () => {
  it('renders triple filters and table', async () => {
    const wrapper = mount(Page, { global: { stubs: { NuxtLink: { template: '<a><slot /></a>' } } } })
    await flushPromises()
    expect(wrapper.find('[data-test="filter-cruise"]').exists()).toBe(true)
    expect(wrapper.find('[data-test="filter-route"]').exists()).toBe(true)
    expect(wrapper.find('[data-test="filter-cabin-type"]').exists()).toBe(true)
    expect(wrapper.find('table').exists()).toBe(true)
  })

  it('shows batch action after row selection', async () => {
    const wrapper = mount(Page, { global: { stubs: { NuxtLink: { template: '<a><slot /></a>' } } } })
    await flushPromises()
    const checkbox = wrapper.find('tbody input[type="checkbox"]')
    await checkbox.setValue(true)
    expect(wrapper.find('[data-test="batch-action"]').exists()).toBe(true)
  })
})
