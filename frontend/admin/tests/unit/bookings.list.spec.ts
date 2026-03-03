import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../app/pages/bookings/index.vue'

// 模拟 useApi：返回两条订单
const mockRequest = vi.fn().mockResolvedValue({ data: {
    list: [
        { id: 1, status: 'created', total_cents: 19900 },
        { id: 2, status: 'paid', total_cents: 38000 },
    ],
    total: 2,
} })
const confirmMock = vi.fn(() => true)

vi.stubGlobal('useApi', () => ({ request: mockRequest }))
vi.stubGlobal('confirm', confirmMock)

beforeEach(() => {
    mockRequest.mockClear()
    confirmMock.mockClear()
    confirmMock.mockReturnValue(true)
    mockRequest.mockResolvedValue({ data: {
        list: [
            { id: 1, status: 'created', total_cents: 19900 },
            { id: 2, status: 'paid', total_cents: 38000 },
        ],
        total: 2,
    } })
})

describe('Admin Bookings', () => {
    it('渲染订单标题', () => {
        const wrapper = mount(Page)
        expect(wrapper.text()).toContain('Bookings')
    })

    it('调用 API 并显示加载后的行', async () => {
        const wrapper = mount(Page)
        await flushPromises()
        expect(mockRequest).toHaveBeenCalledWith('/bookings')
        expect(wrapper.findAll('tbody tr')).toHaveLength(2)
    })

    it('失败时显示错误信息', async () => {
        mockRequest.mockRejectedValueOnce(new Error('network error'))
        const wrapper = mount(Page)
        await flushPromises()
        expect(wrapper.text()).toContain('network error')
    })

    it('空列表时显示总数 0', async () => {
        mockRequest.mockResolvedValueOnce({ data: { list: [], total: 0 } })
        const wrapper = mount(Page)
        await flushPromises()
        expect(wrapper.text()).toContain('总数：0')
        expect(wrapper.findAll('tbody tr')).toHaveLength(0)
    })

    it('无效 ID 删除时显示错误', async () => {
        const wrapper = mount(Page)
        await flushPromises()
        await (wrapper.vm as any).handleDelete('bad')
        await flushPromises()
        expect(wrapper.text()).toContain('无效记录 ID，无法删除')
    })

    it('删除确认取消时不发 DELETE 请求', async () => {
        confirmMock.mockReturnValueOnce(false)
        const wrapper = mount(Page)
        await flushPromises()
        await (wrapper.vm as any).handleDelete(1)
        await flushPromises()
        expect(mockRequest).not.toHaveBeenCalledWith('/bookings/1', { method: 'DELETE' })
    })

    it('删除成功时调用接口并刷新', async () => {
        mockRequest.mockImplementation((url: string, options?: any) => {
            if (url === '/bookings/1' && options?.method === 'DELETE') return Promise.resolve({ data: { ok: true } })
            return Promise.resolve({ data: {
                list: [{ id: 2, status: 'paid', total_cents: 38000 }],
                total: 1,
            } })
        })
        const wrapper = mount(Page)
        await flushPromises()
        await (wrapper.vm as any).handleDelete(1)
        await flushPromises()
        expect(mockRequest).toHaveBeenCalledWith('/bookings/1', { method: 'DELETE' })
    })
})
