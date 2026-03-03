import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import Page from '../../../app/pages/settings/shop.vue'

const mockRequest = vi.fn()
vi.stubGlobal('useApi', () => ({ request: mockRequest }))

describe('Shop settings page', () => {
  beforeEach(() => {
    mockRequest.mockReset()
  })

  it('shows empty state when no shop config', async () => {
    mockRequest.mockResolvedValueOnce({ data: {} })
    const wrapper = mount(Page)
    await flushPromises()
    expect(mockRequest).toHaveBeenCalledWith('/shop-info')
    expect(wrapper.find('[data-test="empty"]').exists()).toBe(true)
  })

  it('shows error state', async () => {
    mockRequest.mockRejectedValueOnce(new Error('shop load error'))
    const wrapper = mount(Page)
    await flushPromises()
    expect(wrapper.find('[data-test="error"]').text()).toContain('shop load error')
  })

  it('saves via api and shows success', async () => {
    mockRequest.mockResolvedValueOnce({ data: {} })
    mockRequest.mockResolvedValueOnce({ data: { ok: true } })
    const wrapper = mount(Page)
    await flushPromises()
    await wrapper.find('[data-test="name"]').setValue('CruiseX')
    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()
    expect(mockRequest).toHaveBeenCalledWith('/shop-info', expect.objectContaining({ method: 'PUT' }))
    expect(wrapper.find('[data-test="success"]').text()).toContain('保存成功')
  })
})
