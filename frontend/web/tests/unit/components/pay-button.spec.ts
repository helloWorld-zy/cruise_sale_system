import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import PayButton from '../../../components/PayButton.vue'

const mockRequest = vi.fn()
vi.stubGlobal('useApi', () => ({ request: mockRequest }))

beforeEach(() => {
  mockRequest.mockReset()
})

describe('PayButton', () => {
  it('点击按钮调用支付 API', async () => {
    mockRequest.mockResolvedValue({ pay_url: 'https://pay.example.com' })
    const wrapper = mount(PayButton, { props: { bookingId: 12, amountCents: 19900 } })
    await wrapper.find('button').trigger('click')
    await flushPromises()
    expect(mockRequest).toHaveBeenCalledWith('/payments', expect.objectContaining({ method: 'POST' }))
    expect(wrapper.emitted('paid')).toBeTruthy()
  })

  it('失败时抛出错误事件', async () => {
    mockRequest.mockRejectedValueOnce(new Error('pay fail'))
    const wrapper = mount(PayButton, { props: { bookingId: 12, amountCents: 19900 } })
    await wrapper.find('button').trigger('click')
    await flushPromises()
    expect(wrapper.emitted('error')?.[0]?.[0]).toContain('pay fail')
  })
})
