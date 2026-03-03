import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/cabin-types/index.vue'

const mockRequest = vi.fn()
vi.stubGlobal('useApi', () => ({ request: mockRequest }))

beforeEach(() => {
  mockRequest.mockReset()
  mockRequest.mockImplementation((url: string) => {
    if (url === '/cruises') {
      return Promise.resolve({ data: [{ id: 1, name: 'Ocean Nova' }] })
    }
    if (url === '/cabin-types') {
      return Promise.resolve({ data: { list: [{ id: 11, name: '阳台房', code: 'BAL', area_min: 24, area_max: 28, max_capacity: 3, status: 1 }], total: 1 } })
    }
    return Promise.resolve({ data: [] })
  })
})

describe('CabinType list', () => {
  it('renders cruise filter and table', async () => {
    const wrapper = mount(Page)
    await flushPromises()
    expect(wrapper.find('[data-test="cruise-filter"]').exists()).toBe(true)
    expect(wrapper.find('table').exists()).toBe(true)
  })

  it('renders rows from api', async () => {
    const wrapper = mount(Page)
    await flushPromises()
    expect(wrapper.text()).toContain('阳台房')
    expect(mockRequest).toHaveBeenCalledWith('/cabin-types', expect.objectContaining({ query: expect.objectContaining({ cruise_id: 1 }) }))
  })
})
