import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/cabins/alerts.vue'

const mockRequest = vi.fn()
vi.stubGlobal('useApi', () => ({ request: mockRequest }))

beforeEach(() => {
  mockRequest.mockReset()
  mockRequest.mockResolvedValue({ data: [{ cabin_sku_id: 7, cabin_code: 'A701', available: 1, alert_threshold: 3 }] })
})

describe('CabinsAlertsPage', () => {
  it('loads alerts and renders warning table', async () => {
    const wrapper = mount(Page, { global: { stubs: { NuxtLink: { template: '<a><slot /></a>' } } } })
    await flushPromises()
    expect(mockRequest).toHaveBeenCalledWith('/cabins/alerts')
    expect(wrapper.text()).toContain('库存预警')
    expect(wrapper.text()).toContain('A701')
    expect(wrapper.text()).toContain('3')
  })
})
