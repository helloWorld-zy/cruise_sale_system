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
            .mockResolvedValueOnce({ id: 8, code: 'A801', bed_type: 'King', price_cents: 19900, has_window: true, has_balcony: false })
            .mockResolvedValueOnce([{ date: '2026-05-01', price_cents: 19900 }])
            .mockResolvedValueOnce({ total: 10, locked: 3, sold: 2, available: 5 })
                const { getByText, findByText } = render(Page, { props: { cabinSkuId: 8 } })
    expect(await findByText('立即预订')).toBeTruthy()
        expect(await findByText('¥199起')).toBeTruthy()
        expect(getByText('含税费')).toBeTruthy()
    })
})
