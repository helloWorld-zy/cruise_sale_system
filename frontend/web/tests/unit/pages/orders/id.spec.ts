import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'

const routeMock = {
  params: { id: '88' },
}

vi.mock('vue-router', () => ({
  useRoute: () => routeMock,
}))

import Page from '../../../../app/pages/orders/[id].vue'

describe('Order Detail Page', () => {
  beforeEach(() => {
    vi.unstubAllGlobals()
    routeMock.params.id = '88'
  })

  it('renders order detail on successful fetch', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ id: '88', status: 'PAID', amount: 12345 }),
    }))

    const wrapper = mount(Page)
    await flushPromises()

    expect(wrapper.text()).toContain('订单号 #88')
    expect(wrapper.text()).toContain('状态: PAID')
    expect(wrapper.text()).toContain('金额: ¥123.45')
  })

  it('shows http error message on non-ok response', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: false, status: 500 }))

    const wrapper = mount(Page)
    await flushPromises()

    expect(wrapper.find('.error').text()).toContain('HTTP 500')
  })

  it('shows generic error message on unknown exception', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue('network-down'))

    const wrapper = mount(Page)
    await flushPromises()

    expect(wrapper.find('.error').text()).toContain('加载订单失败')
  })

  it('renders not found state when api returns null', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: async () => null,
    }))

    const wrapper = mount(Page)
    await flushPromises()

    expect(wrapper.text()).toContain('订单未找到')
  })
})
