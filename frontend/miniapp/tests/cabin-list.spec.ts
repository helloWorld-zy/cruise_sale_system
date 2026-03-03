import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { cleanup, render } from '@testing-library/vue'
import List from '../pages/cabin/list.vue'

const mockRequest = vi.fn()

vi.mock('../src/utils/request', () => ({
    request: (...args: any[]) => mockRequest(...args),
}))

function deferred<T>() {
    let resolve!: (value: T) => void
    let reject!: (reason?: unknown) => void
    const promise = new Promise<T>((res, rej) => {
        resolve = res
        reject = rej
    })
    return { promise, resolve, reject }
}

describe('Cabin List', () => {
    beforeEach(() => {
        mockRequest.mockReset()
    })

    afterEach(() => {
        cleanup()
    })

    it('shows loading state while request is pending', async () => {
        const pending = deferred<any>()
        mockRequest.mockImplementation(() => pending.promise)

        const { getByText, findByText } = render(List)

        expect(getByText('精选舱房')).toBeTruthy()
        expect(await findByText('Loading...')).toBeTruthy()
        pending.resolve({ data: { list: [], total: 0 } })
    })

    it('shows empty state when no cabins returned', async () => {
        mockRequest.mockResolvedValue({ data: { list: [], total: 0 } })

        const { findByText } = render(List)

        expect(await findByText('暂无可预订舱位')).toBeTruthy()
        expect(mockRequest).toHaveBeenCalledWith('/cabins', expect.objectContaining({
            data: expect.objectContaining({ page: 1, page_size: 20 }),
        }))
    })

    it('shows error state when request fails', async () => {
        mockRequest.mockRejectedValue(new Error('network down'))

        const { findByText } = render(List)

        expect(await findByText('network down')).toBeTruthy()
    })

    it('renders cabin cards from api data', async () => {
        mockRequest.mockResolvedValue({
            data: {
                list: [
                    { id: 101, name: '星海阳台舱', bed_type: '双床', area: 28, price_cents: 23900, route_name: '日韩航线', departure_port: '上海', travel_date: '2026-03-10' },
                    { id: 102, code: 'A-102', amenities: '早餐,窗景', price_cents: 19900, route_name: '东南亚航线', departure_port: '深圳', travel_date: '2026-03-20' },
                ],
            },
        })

        const { findByText } = render(List)

        expect(await findByText('星海阳台舱')).toBeTruthy()
        expect(await findByText('A-102')).toBeTruthy()
    })

    it('filters by route and sorts by price asc', async () => {
        mockRequest.mockResolvedValue({
            data: {
                list: [
                    { id: 201, code: 'Z-201', price_cents: 30000, route_name: '日韩航线', departure_port: '上海', travel_date: '2026-04-01' },
                    { id: 202, code: 'A-202', price_cents: 20000, route_name: '日韩航线', departure_port: '上海', travel_date: '2026-04-02' },
                    { id: 203, code: 'B-203', price_cents: 10000, route_name: '欧美航线', departure_port: '天津', travel_date: '2026-04-03' },
                ],
            },
        })

        const { findByPlaceholderText, findByText, queryByText, getByText, container } = render(List)

        const routeInput = await findByPlaceholderText('按航线筛选')
        await routeInput.focus()
        await routeInput.setSelectionRange(0, 0)
        routeInput.value = '日韩'
        routeInput.dispatchEvent(new Event('input'))

        await findByText('Z-201')
        expect(queryByText('B-203')).toBeNull()

        getByText('价格升序').dispatchEvent(new Event('click'))
        await findByText('A-202')

        expect(container.textContent?.indexOf('A-202')).toBeLessThan(container.textContent?.indexOf('Z-201') || 0)
    })
})
