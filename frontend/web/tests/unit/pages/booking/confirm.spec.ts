import { describe, it, expect, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../../pages/booking/confirm.vue'

// Nuxt 自动导入 stub
vi.stubGlobal('useRoute', () => ({ query: { voyage_id: '2', cabin_sku_id: '3', guests: '2' } }))
const mockFetch = vi.fn().mockResolvedValue({})
vi.stubGlobal('$fetch', mockFetch)

describe('Booking Confirm', () => {
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

    it('提交调用 API 并显示成功', async () => {
        const wrapper = mount(Page)
        await wrapper.find('input#guests').setValue('4')
        await wrapper.find('form').trigger('submit')
        await flushPromises()
        expect(mockFetch).toHaveBeenCalledWith(
            '/api/v1/bookings',
            expect.objectContaining({
                method: 'POST',
                body: expect.objectContaining({ guests: 4 }),
            })
        )
        expect(wrapper.text()).toContain('预订成功')
    })

    it('提交进行中显示 loading 文案', async () => {
        vi.stubGlobal('useRoute', () => ({ query: { voyage_id: '2', cabin_sku_id: '3', guests: '2' } }))
        let resolver: (() => void) | null = null
        mockFetch.mockImplementationOnce(() => new Promise((resolve) => {
            resolver = () => resolve({})
        }))

        const wrapper = mount(Page)
        await wrapper.find('form').trigger('submit')
        expect(wrapper.text()).toContain('提交中…')

        resolver?.()
        await flushPromises()
    })

    it('失败时显示错误信息', async () => {
        mockFetch.mockRejectedValueOnce({ message: 'cabin unavailable' })
        const wrapper = mount(Page)
        await wrapper.find('form').trigger('submit')
        await flushPromises()
        expect(wrapper.find('.error').text()).toContain('cabin unavailable')
    })

    it('缺少参数时展示提示并禁止提交', () => {
        vi.stubGlobal('useRoute', () => ({ query: {} }))
        const wrapper = mount(Page)
        expect(wrapper.text()).toContain('缺少航次信息')
        expect(wrapper.text()).toContain('缺少舶房信息')
        expect(wrapper.find('button[type="submit"]').attributes('disabled')).toBeDefined()
    })

    it('失败时优先展示 data.message', async () => {
        vi.stubGlobal('useRoute', () => ({ query: { voyage_id: '2', cabin_sku_id: '3', guests: '2' } }))
        mockFetch.mockRejectedValueOnce({ data: { message: 'from data message' } })
        const wrapper = mount(Page)
        await wrapper.find('form').trigger('submit')
        await flushPromises()
        expect(wrapper.find('.error').text()).toContain('from data message')
    })

    it('失败时无 message 字段使用默认文案', async () => {
        vi.stubGlobal('useRoute', () => ({ query: { voyage_id: '2', cabin_sku_id: '3', guests: '2' } }))
        mockFetch.mockRejectedValueOnce({})
        const wrapper = mount(Page)
        await wrapper.find('form').trigger('submit')
        await flushPromises()
        expect(wrapper.find('.error').text()).toContain('预订失败，请重试')
    })
})
