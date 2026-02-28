import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../pages/pay/[id].vue'

const mockRequest = vi.fn()
vi.stubGlobal('useApi', () => ({ request: mockRequest }))
vi.stubGlobal('useRoute', () => ({ params: { id: '42' } }))

beforeEach(() => {
    mockRequest.mockReset()
    mockRequest.mockResolvedValue({ data: { id: 42, total_cents: 19900, status: 'created' } })
})

describe('Pay Page', () => {
    it('加载订单信息', async () => {
        mount(Page)
        await flushPromises()
        expect(mockRequest).toHaveBeenCalledWith('/bookings/42')
    })

    it('显示订单金额', async () => {
        const wrapper = mount(Page)
        await flushPromises()
        expect(wrapper.text()).toContain('19900')
    })

    it('点击支付按钮触发支付 API', async () => {
        mockRequest
            .mockResolvedValueOnce({ data: { id: 42, total_cents: 19900 } })
            .mockResolvedValueOnce({ pay_url: 'https://pay.example.com' })
        const wrapper = mount(Page)
        await flushPromises()
        await wrapper.find('button').trigger('click')
        await flushPromises()
        expect(mockRequest).toHaveBeenCalledTimes(2)
    })

    it('支付失败时显示错误', async () => {
        mockRequest
            .mockResolvedValueOnce({ data: { id: 42, total_cents: 19900 } })
            .mockRejectedValueOnce(new Error('payment failed'))
        const wrapper = mount(Page)
        await flushPromises()
        await wrapper.find('button').trigger('click')
        await flushPromises()
        expect(wrapper.text()).toContain('payment failed')
    })
})
