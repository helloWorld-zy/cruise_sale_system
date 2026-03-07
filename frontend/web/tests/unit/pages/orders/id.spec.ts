import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'

const routeMock = {
  params: { id: '88' },
}

vi.mock('vue-router', () => ({
  useRoute: () => routeMock,
}))

const mockRequest = vi.fn()
vi.stubGlobal('useApi', () => ({ request: mockRequest }))
vi.stubGlobal('sessionStorage', { getItem: () => 'fake-token' })

import Page from '../../../../app/pages/orders/[id].vue'

describe('Order Detail Page', () => {
  beforeEach(() => {
    mockRequest.mockReset()
    routeMock.params.id = '88'
  })

  it('renders order detail on successful fetch', async () => {
    mockRequest.mockResolvedValueOnce({
      data: { id: '88', status: 'PAID', amount: 12345 },
    })

    const wrapper = mount(Page)
    await flushPromises()

    expect(wrapper.text()).toContain('Order #88')
    expect(wrapper.text()).toContain('PAID')
    expect(wrapper.text()).toContain('¥123.45')
  })

  it('shows generic error message on unknown exception', async () => {
    mockRequest.mockRejectedValueOnce(new Error('加载订单失败'))

    const wrapper = mount(Page)
    await flushPromises()

    expect(wrapper.find('.error').text()).toContain('加载订单失败')
  })

  it('renders not found state when api returns null', async () => {
    mockRequest.mockResolvedValueOnce(null)

    const wrapper = mount(Page)
    await flushPromises()

    expect(wrapper.text()).toContain('订单未找到')
  })
})
