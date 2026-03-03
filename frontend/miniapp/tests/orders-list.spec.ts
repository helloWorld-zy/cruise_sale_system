import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { cleanup, render } from '@testing-library/vue'
import Page from '../pages/orders/list.vue'

const mockRequest = vi.fn()

vi.mock('../src/utils/request', () => ({
  request: (...args: any[]) => mockRequest(...args),
}))

describe('Miniapp Orders List', () => {
  beforeEach(() => {
    mockRequest.mockReset()
  })

  afterEach(() => {
    cleanup()
  })

  it('renders list and actions', async () => {
    mockRequest.mockResolvedValue({
      data: {
        list: [
          { id: 1, status: 'pending_payment', total_cents: 18800 },
          { id: 2, status: 'paid', total_cents: 28800 },
        ],
      },
    })

    const { findByText, findAllByText } = render(Page)
    expect(await findByText('我的订单')).toBeTruthy()
    expect(await findByText('去支付')).toBeTruthy()
    expect((await findAllByText('申请退改')).length).toBeGreaterThan(0)
  })

  it('shows empty state', async () => {
    mockRequest.mockResolvedValue({ data: { list: [] } })
    const { findByText } = render(Page)

    expect(await findByText('暂无订单')).toBeTruthy()
  })

  it('shows error state', async () => {
    mockRequest.mockRejectedValue(new Error('network down'))
    const { findByText } = render(Page)

    expect(await findByText('network down')).toBeTruthy()
  })
})
