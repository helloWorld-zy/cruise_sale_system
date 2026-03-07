import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/cabins/index.vue'

const mockRequest = vi.fn()
vi.stubGlobal('useApi', () => ({ request: mockRequest }))
vi.stubGlobal('confirm', () => true)

beforeEach(() => {
  mockRequest.mockReset()
  mockRequest.mockImplementation((url: string, options?: any) => {
    if (url === '/cruises') return Promise.resolve({ data: [{ id: 1, name: 'Ocean Nova' }, { id: 2, name: 'Ocean Star' }] })
    if (url === '/voyages') return Promise.resolve({ data: [{ id: 2, code: 'VY-001', cruise_id: 1 }, { id: 3, code: 'VY-002', cruise_id: 2 }] })
    if (url === '/cabin-types') return Promise.resolve({ data: { list: [{ id: 3, name: '阳台房' }] } })
    if (url === '/cabins') return Promise.resolve({ data: { list: [{ id: 10, code: 'A101', area: 24, total: 10, available: 8, status: 1 }], total: 1 } })
    return Promise.resolve({ data: [] })
  })
})

describe('CabinsIndexPage', () => {
  const globalStubs = {
    NuxtLink: { template: '<a><slot /></a>' },
    AdminActionLink: { template: '<a><slot /></a>' },
    AdminPageHeader: { props: ['title'], template: '<div>{{ title }}<slot /><slot name="actions" /></div>' },
    AdminFilterBar: { template: '<div><slot /></div>' },
    AdminDataCard: { template: '<div><slot /></div>' },
    AdminStatusTag: { props: ['text'], template: '<span>{{ text }}</span>' },
  }

  it('renders triple filters and table', async () => {
    const wrapper = mount(Page, { global: { stubs: globalStubs } })
    await flushPromises()
    expect(wrapper.find('[data-test="filter-cruise"]').exists()).toBe(true)
    expect(wrapper.find('[data-test="filter-voyage"]').exists()).toBe(true)
    expect(wrapper.find('[data-test="filter-cabin-type"]').exists()).toBe(true)
    expect(wrapper.find('table').exists()).toBe(true)
  })

  it('shows batch action after row selection', async () => {
    const wrapper = mount(Page, { global: { stubs: globalStubs } })
    await flushPromises()
    const checkbox = wrapper.find('tbody input[type="checkbox"]')
    await checkbox.setValue(true)
    expect(wrapper.find('[data-test="batch-action"]').exists()).toBe(true)
  })

  it('reloads cabin type options when cruise filter changes', async () => {
    const wrapper = mount(Page, { global: { stubs: globalStubs } })
    await flushPromises()

    const cruiseFilter = wrapper.find('[data-test="filter-cruise"]')
    await cruiseFilter.setValue('2')
    await flushPromises()

    const calls = mockRequest.mock.calls.filter((call) => call[0] === '/cabin-types')
    expect(calls.some((call) => Number(call[1]?.query?.cruise_id) === 2)).toBe(true)
  })
})
