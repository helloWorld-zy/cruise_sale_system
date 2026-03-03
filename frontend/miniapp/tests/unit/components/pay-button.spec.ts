import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { cleanup, fireEvent, render } from '@testing-library/vue'
import PayButton from '../../../components/PayButton.vue'

const mockRequest = vi.fn()
vi.mock('../../../src/utils/request', () => ({
  request: (...args: any[]) => mockRequest(...args),
}))

declare global {
  // eslint-disable-next-line no-var
  var uni: any
}

describe('Miniapp PayButton', () => {
  beforeEach(() => {
    mockRequest.mockReset()
    globalThis.uni = { requestPayment: vi.fn() }
  })

  afterEach(() => cleanup())

  it('emits paid when payment url returned without pay params', async () => {
    mockRequest.mockResolvedValueOnce({ data: { pay_url: 'https://pay.test/url' } })
    const { getByText, emitted } = render(PayButton, { props: { bookingId: 9, amountCents: 18800 } })

    await fireEvent.click(getByText('Pay Now'))

    expect(mockRequest).toHaveBeenCalledWith('/payments', expect.objectContaining({ method: 'POST' }))
    expect(emitted().paid).toEqual([['https://pay.test/url']])
  })

  it('calls uni.requestPayment and emits paid on success', async () => {
    mockRequest.mockResolvedValueOnce({
      data: {
        pay_url: 'wx://pay',
        pay_params: { timeStamp: '1', nonceStr: 'n', package: 'p', signType: 'MD5', paySign: 's' },
      },
    })

    const { getByText, emitted } = render(PayButton, { props: { bookingId: 9, amountCents: 18800 } })
    await fireEvent.click(getByText('Pay Now'))

    const paymentCall = globalThis.uni.requestPayment.mock.calls[0][0]
    paymentCall.success()

    expect(globalThis.uni.requestPayment).toHaveBeenCalledTimes(1)
    expect(emitted().paid).toEqual([['wx://pay']])
  })

  it('emits error when requestPayment fails', async () => {
    mockRequest.mockResolvedValueOnce({
      data: {
        pay_url: 'wx://pay',
        pay_params: { timeStamp: '1', nonceStr: 'n', package: 'p', signType: 'MD5', paySign: 's' },
      },
    })

    const { getByText, emitted } = render(PayButton, { props: { bookingId: 9, amountCents: 18800 } })
    await fireEvent.click(getByText('Pay Now'))

    const paymentCall = globalThis.uni.requestPayment.mock.calls[0][0]
    paymentCall.fail({ errMsg: 'user cancel' })

    expect(emitted().error).toEqual([['user cancel']])
  })

  it('emits error on request exception', async () => {
    mockRequest.mockRejectedValueOnce(new Error('gateway down'))
    const { getByText, emitted } = render(PayButton, { props: { bookingId: 9, amountCents: 18800 } })

    await fireEvent.click(getByText('Pay Now'))

    expect(emitted().error).toEqual([['gateway down']])
  })
})
