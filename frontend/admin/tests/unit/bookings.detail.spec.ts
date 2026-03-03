import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import Page from '../../app/pages/bookings/[id].vue'

const mockRequest = vi.fn().mockResolvedValue({
    data: { id: 12, status: 'created', total_cents: 19900 },
})
const mockNavigateTo = vi.fn().mockResolvedValue(undefined)
const confirmMock = vi.fn(() => true)

vi.stubGlobal('useRoute', () => ({ params: { id: '12' } }))
vi.stubGlobal('useApi', () => ({ request: mockRequest }))
vi.stubGlobal('navigateTo', mockNavigateTo)
vi.stubGlobal('confirm', confirmMock)

beforeEach(() => {
    mockRequest.mockClear()
    mockNavigateTo.mockClear()
    confirmMock.mockClear()
    confirmMock.mockReturnValue(true)
    mockRequest.mockResolvedValue({
        data: { id: 12, status: 'created', total_cents: 19900 },
    })
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

    it('updates status when save button is clicked', async () => {
        const wrapper = mount(Page)
        await flushPromises()

        const input = wrapper.find('input')
        await input.setValue('paid')
        const saveBtn = wrapper.findAll('button').find((b) => b.text().includes('保存状态'))
        await saveBtn!.trigger('click')
        await flushPromises()

        expect(mockRequest).toHaveBeenCalledWith('/bookings/12', {
            method: 'PUT',
            body: { status: 'paid' },
        })
        expect(wrapper.text()).toContain('Status: paid')
    })

    it('does not delete when confirm is cancelled', async () => {
        confirmMock.mockReturnValueOnce(false)
        const wrapper = mount(Page)
        await flushPromises()

        const deleteBtn = wrapper.findAll('button').find((b) => b.text().includes('删除订单'))
        await deleteBtn!.trigger('click')
        await flushPromises()

        expect(mockRequest).not.toHaveBeenCalledWith('/bookings/12', { method: 'DELETE' })
        expect(mockNavigateTo).not.toHaveBeenCalled()
    })
})
