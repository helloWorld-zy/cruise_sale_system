import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import Page from '../../../../pages/booking/index.vue'

// Nuxt 自动导入 stub
const mockPush = vi.fn()
const mockQuery = {}
vi.stubGlobal('useRouter', () => ({ push: mockPush }))
vi.stubGlobal('useRoute', () => ({ query: mockQuery }))

describe('Booking Step 1', () => {
    it('渲染标题和输入框', () => {
        const wrapper = mount(Page)
        expect(wrapper.text()).toContain('Booking Step 1')
        expect(wrapper.find('input#voyage-id').exists()).toBe(true)
        expect(wrapper.find('input#cabin-sku-id').exists()).toBe(true)
        expect(wrapper.find('input#guests').exists()).toBe(true)
    })

    it('尚未填写时提交按钮禁用', () => {
        const wrapper = mount(Page)
        expect(wrapper.find('button[type="submit"]').attributes('disabled')).toBeDefined()
    })

    it('未填写时提交不会触发跳转', async () => {
        const wrapper = mount(Page)
        await wrapper.find('form').trigger('submit')
        expect(mockPush).not.toHaveBeenCalled()
    })

    it('填写后提交按钮可用并跳转', async () => {
        mockPush.mockResolvedValue(undefined)
        const wrapper = mount(Page)
        await wrapper.find('input#voyage-id').setValue('2')
        await wrapper.find('input#cabin-sku-id').setValue('3')
        await wrapper.find('input#guests').setValue('2')
        await wrapper.find('form').trigger('submit')
        expect(mockPush).toHaveBeenCalledWith(expect.objectContaining({
            path: '/booking/confirm',
        }))
    })

    it('跳转失败时显示错误提示', async () => {
        mockPush.mockRejectedValueOnce(new Error('route push failed'))
        const wrapper = mount(Page)
        await wrapper.find('input#voyage-id').setValue('2')
        await wrapper.find('input#cabin-sku-id').setValue('3')
        await wrapper.find('input#guests').setValue('2')
        await wrapper.find('form').trigger('submit')
        expect(wrapper.text()).toContain('route push failed')
    })
})
