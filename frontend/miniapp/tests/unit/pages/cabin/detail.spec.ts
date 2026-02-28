import { describe, it, expect, vi } from 'vitest'
import { render } from '@testing-library/vue'
import Page from '../../../../pages/cabin/detail.vue'

const mockRequest = vi.fn()
vi.mock('../../../../src/utils/request', () => ({
    request: (...args: any[]) => mockRequest(...args),
}))

describe('Cabin Detail', () => {
        it('renders view', async () => {
                mockRequest
                    .mockResolvedValueOnce({ id: 8, name: 'Deluxe', price_cents: 19900, available: 3 })
                    .mockResolvedValueOnce([{ date: '2026-05-01', price_cents: 19900 }])
                const { getByText, findByText } = render(Page, { props: { cabinSkuId: 8 } })
        expect(getByText('Cabin Detail')).toBeTruthy()
                expect(await findByText('名称：Deluxe')).toBeTruthy()
    })
})
