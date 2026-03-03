import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'
import Page from '../../../../app/pages/booking/success.vue'

const mockNavigateTo = vi.fn()
const mockRoute = {
  query: {
    order_id: '42',
  },
}

vi.stubGlobal('useRoute', () => mockRoute)
vi.stubGlobal('navigateTo', mockNavigateTo)

describe('Booking Success Page', () => {
  beforeEach(() => {
    vi.useFakeTimers()
    mockNavigateTo.mockReset()
    mockRoute.query.order_id = '42'
  })

  afterEach(() => {
    vi.runOnlyPendingTimers()
    vi.useRealTimers()
  })

  it('renders success state with order id', () => {
    const wrapper = mount(Page)
    expect(wrapper.text()).toContain('预订成功')
    expect(wrapper.text()).toContain('订单号 #42')
  })

  it('clicking button navigates to order detail', async () => {
    const wrapper = mount(Page)
    await wrapper.find('button').trigger('click')
    expect(mockNavigateTo).toHaveBeenCalledWith('/orders/42')
  })

  it('auto redirects after 5 seconds', () => {
    mount(Page)
    vi.advanceTimersByTime(5000)
    expect(mockNavigateTo).toHaveBeenCalledWith('/orders/42')
  })

  it('disables button and shows fallback when order id missing', () => {
    mockRoute.query.order_id = 'abc'
    const wrapper = mount(Page)
    const button = wrapper.find('button')
    expect(button.attributes('disabled')).toBeDefined()
    expect(wrapper.text()).toContain('未获取到订单号')

    vi.advanceTimersByTime(5000)
    expect(mockNavigateTo).not.toHaveBeenCalled()
  })
})
