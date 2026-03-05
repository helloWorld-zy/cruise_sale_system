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
    expect(await findByText('金额：¥199.00')).toBeTruthy()
  })

  it('缺少 bookingId 时提示错误', async () => {
    const { findByText } = render(Page)
    expect(await findByText('缺少 bookingId 参数')).toBeTruthy()
  })

  it('加载订单失败时展示错误信息', async () => {
    mockRequest.mockRejectedValueOnce(new Error('加载失败'))
    const { findByText } = render(Page, { props: { bookingId: 42 } })
    expect(await findByText('加载失败')).toBeTruthy()
  })
})
