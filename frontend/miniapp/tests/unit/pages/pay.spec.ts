import { afterEach, beforeEach, describe, it, expect, vi } from 'vitest'
import { cleanup, render } from '@testing-library/vue'
import Page from '../../../pages/pay/pay.vue'

const mockRequest = vi.fn()
vi.mock('../../../src/utils/request', () => ({
  request: (...args: any[]) => mockRequest(...args),
}))

beforeEach(() => {
  mockRequest.mockReset()
  mockRequest.mockResolvedValue({ id: 42, total_cents: 19900 })
})

afterEach(() => cleanup())

describe('Miniapp Pay Page', () => {
  it('加载订单信息', async () => {
    render(Page, { props: { bookingId: 42 } })
    expect(mockRequest).toHaveBeenCalledWith('/bookings/42')
  })

  it('渲染订单金额', async () => {
    const { findByText } = render(Page, { props: { bookingId: 42 } })
    expect(await findByText('金额：19900')).toBeTruthy()
  })
})
