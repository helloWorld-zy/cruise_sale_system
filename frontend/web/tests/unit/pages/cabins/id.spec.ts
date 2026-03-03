import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { flushPromises } from '@vue/test-utils'
import Page from '../../../../app/pages/cabins/[id].vue'

const mockRequest = vi.fn()
vi.stubGlobal('useApi', () => ({ request: mockRequest }))
vi.stubGlobal('useRoute', () => ({ params: { id: '7' } }))

beforeEach(() => {
    mockRequest.mockReset()
    mockRequest.mockImplementation((url: string) => {
        if (url === '/cabins/7') return Promise.resolve({ data: { id: 7, code: 'A701', bed_type: 'King', area: 28, max_guests: 3, has_window: true, has_balcony: true, position: 'mid', orientation: 'port' } })
        if (url === '/cabins/7/prices') return Promise.resolve({ data: [{ date: '2026-02-01', price_cents: 199000, price_type: 'base' }] })
        if (url === '/cabins/7/inventory') return Promise.resolve({ data: { total: 10, locked: 2, sold: 1, available: 7 } })
        return Promise.resolve({ data: [] })
    })
})

describe('Cabin Detail', () => {
    it('renders enhanced cabin detail view', async () => {
        const wrapper = mount(Page)
        await flushPromises()
        expect(wrapper.find('[data-test="attr-tags"]').exists()).toBe(true)
        expect(wrapper.find('[data-test="price-calendar"]').exists()).toBe(true)
        expect(wrapper.text()).toContain('立即预订')
    })
})
