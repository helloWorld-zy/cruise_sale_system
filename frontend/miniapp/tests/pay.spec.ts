import { describe, it, expect, beforeEach, vi } from 'vitest'
import { render } from '@testing-library/vue'
import Pay from '../pages/pay/pay.vue'

const mockRequest = vi.fn().mockResolvedValue({ id: 42, total_cents: 19900 })
vi.mock('../src/utils/request', () => ({
  request: (...args: any[]) => mockRequest(...args),
}))

beforeEach(() => {
    mockRequest.mockClear()
})

describe('Miniapp Pay', () => {
    it('shows pay amount', async () => {
        const { findByText } = render(Pay, { props: { bookingId: 42 } })
        expect(await findByText('金额：19900')).toBeTruthy()
    })
})
