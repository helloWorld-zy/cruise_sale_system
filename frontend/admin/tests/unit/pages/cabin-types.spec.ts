import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/cabin-types/index.vue'

const mockRequest = vi.fn()
vi.stubGlobal('useApi', () => ({ request: mockRequest }))

beforeEach(() => {
  mockRequest.mockReset()
  mockRequest.mockImplementation((url: string, options?: any) => {
    if (url === '/cruises') {
      return Promise.resolve({ data: [{ id: 1, name: 'Ocean Nova' }, { id: 2, name: 'Ocean Star' }] })
    }
    if (url === '/cabin-types') {
      const cruiseID = Number(options?.query?.cruise_id || 0)
      if (cruiseID === 1) {
        return Promise.resolve({ data: { list: [{ id: 11, name: '阳台房', code: 'BAL', area_min: 24, area_max: 28, max_capacity: 3, status: 1 }], total: 1 } })
      }
      if (cruiseID === 2) {
        return Promise.resolve({ data: { list: [{ id: 12, name: '标准内舱', code: 'INS', area_min: 18, area_max: 22, max_capacity: 2, status: 1 }], total: 1 } })
      }
      return Promise.resolve({ data: { list: [], total: 0 } })
    }
    return Promise.resolve({ data: [] })
  })
})

describe('CabinType list', () => {
  it('renders cruise filter and table', async () => {
    const wrapper = mount(Page)
    await flushPromises()
    expect(wrapper.find('[data-test="cruise-filter"]').exists()).toBe(true)
    expect((wrapper.find('[data-test="cruise-filter"]').element as HTMLSelectElement).value).toBe('0')
    expect(wrapper.find('table').exists()).toBe(true)
  })

  it('renders all cabin types when cruise is not selected', async () => {
    const wrapper = mount(Page)
    await flushPromises()
    expect(wrapper.text()).toContain('阳台房')
    expect(wrapper.text()).toContain('标准内舱')
    const calls = mockRequest.mock.calls.filter((call) => call[0] === '/cabin-types')
    expect(calls.some((call) => Number(call[1]?.query?.cruise_id) === 1)).toBe(true)
    expect(calls.some((call) => Number(call[1]?.query?.cruise_id) === 2)).toBe(true)
  })
})
