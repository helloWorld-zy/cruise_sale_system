import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/voyages/index.vue'

const mockRequest = vi.fn()
vi.stubGlobal('useApi', () => ({ request: mockRequest }))

beforeEach(() => {
  mockRequest.mockReset()
  mockRequest.mockResolvedValue({ data: [{ id: 1, code: 'V-1', name: 'Voyage', route_id: 10, departure_date: '2026-03-01' }] })
})

describe('VoyagesIndexPage', () => {
  it('调用 voyages API', async () => {
    mount(Page)
    await flushPromises()
    expect(mockRequest).toHaveBeenCalledWith('/voyages')
  })

  it('渲染航次行', async () => {
    const wrapper = mount(Page)
    await flushPromises()
    expect(wrapper.text()).toContain('V-1')
    expect(wrapper.text()).toContain('2026-03-01')
  })
})
