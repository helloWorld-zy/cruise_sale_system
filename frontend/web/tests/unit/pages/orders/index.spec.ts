import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../../app/pages/orders/index.vue'

const mockRequest = vi.fn()
vi.stubGlobal('useApi', () => ({ request: mockRequest }))
vi.stubGlobal('sessionStorage', { getItem: () => 'fake-token' })
const globalMountOptions = {
  stubs: { NuxtLink: { template: '<a><slot /></a>' } }
}

beforeEach(() => {
  mockRequest.mockReset()
  mockRequest.mockResolvedValue({
    data: {
      list: [
        { id: 101, status: 'pending_payment', total_cents: 19900 },
        { id: 102, status: 'paid', total_cents: 29900 },
      ],
    },
  })
})

describe('Web Orders Index', () => {
  it('renders tabs and list', async () => {
    const wrapper = mount(Page, { global: globalMountOptions })
    await flushPromises()

    expect(wrapper.text()).toContain('我的订单')
    expect(wrapper.text()).toContain('待支付')
    expect(wrapper.text()).toContain('查看详情')
  })

  it('filters by status tab', async () => {
    const wrapper = mount(Page, { global: globalMountOptions })
    await flushPromises()

    const paidTab = wrapper.findAll('button').find((b) => b.text() === '已支付')
    await paidTab?.trigger('click')
    await flushPromises()

    expect(wrapper.text()).not.toContain('订单 #101')
    expect(wrapper.text()).toContain('订单 #102')
  })

  it('shows empty state', async () => {
    mockRequest.mockResolvedValueOnce({ data: { list: [] } })
    const wrapper = mount(Page, { global: globalMountOptions })
    await flushPromises()

    expect(wrapper.find('[data-test="empty"]').exists()).toBe(true)
  })

  it('shows error state', async () => {
    mockRequest.mockRejectedValueOnce(new Error('load failed'))
    const wrapper = mount(Page, { global: globalMountOptions })
    await flushPromises()

    expect(wrapper.find('[data-test="error"]').text()).toContain('load failed')
  })
})
