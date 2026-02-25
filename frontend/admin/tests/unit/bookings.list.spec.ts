import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../app/pages/bookings/index.vue'

// 模拟 useApi：返回两条订单
const mockRequest = vi.fn().mockResolvedValue({ data: [
    { id: 1, status: 'created', total: 19900 },
    { id: 2, status: 'paid',    total: 38000 },
] })

vi.stubGlobal('useApi', () => ({ request: mockRequest }))

beforeEach(() => mockRequest.mockClear())

describe('Admin Bookings', () => {
    it('渲染订单标题', () => {
        const wrapper = mount(Page, {
            global: { stubs: { BookingRow: true } }
        })
        expect(wrapper.text()).toContain('Bookings')
    })

    it('调用 API 并显示加载后的行', async () => {
        const wrapper = mount(Page, {
            global: { stubs: { BookingRow: true } }
        })
        await flushPromises()
        expect(mockRequest).toHaveBeenCalledWith('/bookings')
        expect(wrapper.findAllComponents({ name: 'BookingRow' })).toHaveLength(2)
    })

    it('失败时显示错误信息', async () => {
        mockRequest.mockRejectedValueOnce(new Error('network error'))
        const wrapper = mount(Page, {
            global: { stubs: { BookingRow: true } }
        })
        await flushPromises()
        expect(wrapper.text()).toContain('network error')
    })
})
