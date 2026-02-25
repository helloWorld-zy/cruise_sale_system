import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import Page from '../../app/pages/bookings/[id].vue'

const mockRequest = vi.fn().mockResolvedValue({
    data: { id: 12, status: 'created', total: 19900 },
})

vi.stubGlobal('useRoute', () => ({ params: { id: '12' } }))
vi.stubGlobal('useApi', () => ({ request: mockRequest }))

beforeEach(() => {
    mockRequest.mockClear()
})

describe('Admin Booking Detail', () => {
    it('loads detail with route id', async () => {
        const wrapper = mount(Page)
        await flushPromises()

        expect(mockRequest).toHaveBeenCalledWith('/bookings/12')
        expect(wrapper.text()).toContain('Status: created')
        expect(wrapper.text()).toContain('Total: 19900')
    })

    it('shows error text when request fails', async () => {
        mockRequest.mockRejectedValueOnce(new Error('network error'))
        const wrapper = mount(Page)
        await flushPromises()

        expect(wrapper.text()).toContain('network error')
    })
})
