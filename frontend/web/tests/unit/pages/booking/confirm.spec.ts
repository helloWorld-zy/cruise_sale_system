import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../../app/pages/booking/confirm.vue'

// Nuxt 自动导入 stub
vi.stubGlobal('useRoute', () => ({ query: { voyage_id: '2', cabin_sku_id: '3', guests: '2' } }))
const mockRequest = vi.fn().mockResolvedValue({ data: { id: 42 } })
const mockNavigateTo = vi.fn().mockResolvedValue(undefined)
vi.stubGlobal('useApi', () => ({ request: mockRequest }))
vi.stubGlobal('navigateTo', mockNavigateTo)

describe('Booking Confirm', () => {
    beforeEach(() => {
        mockRequest.mockClear()
        mockNavigateTo.mockClear()
    })

    it('渲染标题和表单', () => {
        const wrapper = mount(Page)
        expect(wrapper.text()).toContain('Confirm Booking')
        expect(wrapper.find('form[data-testid="booking-form"]').exists()).toBe(true)
        expect(wrapper.find('input#guests').exists()).toBe(true)
    })

    it('正确场景下按钮可用', () => {
        const wrapper = mount(Page)
        // 从 query 初始化后 voyageId 和 cabinSkuId > 0
        expect(wrapper.find('button[type="submit"]').attributes('disabled')).toBeUndefined()
    })

    it('提交调用 API 后跳转到成功页', async () => {
        const wrapper = mount(Page)
        await wrapper.find('input#guests').setValue('4')
        await wrapper.find('form').trigger('submit')
        await flushPromises()
        expect(mockRequest).toHaveBeenCalledWith(
            '/bookings',
            expect.objectContaining({
                method: 'POST',
                body: expect.objectContaining({ guests: 4 }),
            })
        )
        expect(mockNavigateTo).toHaveBeenCalledWith({
            path: '/booking/success',
            query: { order_id: '42' },
        })
    })

    it('提交进行中显示 loading 文案', async () => {
        vi.stubGlobal('useRoute', () => ({ query: { voyage_id: '2', cabin_sku_id: '3', guests: '2' } }))
        let resolver: (() => void) | null = null
        mockRequest.mockImplementationOnce(() => new Promise((resolve) => {
            resolver = () => resolve({})
        }))

        const wrapper = mount(Page)
        await wrapper.find('form').trigger('submit')
        expect(wrapper.text()).toContain('提交中…')

        resolver?.()
        await flushPromises()
    })

    it('失败时显示错误信息', async () => {
        mockRequest.mockRejectedValueOnce({ message: 'cabin unavailable' })
        const wrapper = mount(Page)
        await wrapper.find('form').trigger('submit')
        await flushPromises()
        expect(wrapper.find('.error').text()).toContain('cabin unavailable')
    })

    it('缺少参数时展示提示并禁止提交', () => {
        vi.stubGlobal('useRoute', () => ({ query: {} }))
        const wrapper = mount(Page)
        expect(wrapper.text()).toContain('缺少航次信息')
        expect(wrapper.text()).toContain('缺少舱房信息')
        expect(wrapper.find('button[type="submit"]').attributes('disabled')).toBeDefined()
    })

    it('失败时优先展示 data.message', async () => {
        vi.stubGlobal('useRoute', () => ({ query: { voyage_id: '2', cabin_sku_id: '3', guests: '2' } }))
        mockRequest.mockRejectedValueOnce({ data: { message: 'from data message' } })
        const wrapper = mount(Page)
        await wrapper.find('form').trigger('submit')
        await flushPromises()
        expect(wrapper.find('.error').text()).toContain('from data message')
    })

    it('失败时无 message 字段使用默认文案', async () => {
        vi.stubGlobal('useRoute', () => ({ query: { voyage_id: '2', cabin_sku_id: '3', guests: '2' } }))
        mockRequest.mockRejectedValueOnce({})
        const wrapper = mount(Page)
        await wrapper.find('form').trigger('submit')
        await flushPromises()
        expect(wrapper.find('.error').text()).toContain('预订失败，请重试')
    })

    it('成功但缺少订单号时展示错误', async () => {
        mockRequest.mockResolvedValueOnce({ data: {} })
        const wrapper = mount(Page)
        await wrapper.find('form').trigger('submit')
        await flushPromises()
        expect(wrapper.find('.error').text()).toContain('订单号缺失')
    })
})
