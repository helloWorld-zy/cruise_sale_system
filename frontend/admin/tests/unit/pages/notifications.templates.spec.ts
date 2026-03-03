import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import Page from '../../../app/pages/notifications/templates.vue'

const mockRequest = vi.fn()
vi.stubGlobal('useApi', () => ({ request: mockRequest }))

describe('Notification templates page', () => {
  beforeEach(() => {
    mockRequest.mockReset()
  })

  it('shows table after loading data', async () => {
    mockRequest.mockResolvedValueOnce({ data: [{ id: 1, event_type: 'order_paid', channel: 'sms', template: 'ok' }] })
    const wrapper = mount(Page)
    await flushPromises()
    expect(mockRequest).toHaveBeenCalledWith('/notification-templates')
    expect(wrapper.find('[data-test="table"]').exists()).toBe(true)
    expect(wrapper.text()).toContain('order_paid')
  })

  it('shows empty state', async () => {
    mockRequest.mockResolvedValueOnce({ data: [] })
    const wrapper = mount(Page)
    await flushPromises()
    expect(wrapper.find('[data-test="empty"]').exists()).toBe(true)
  })

  it('creates template and reloads list', async () => {
    mockRequest.mockResolvedValueOnce({ data: [] })
    mockRequest.mockResolvedValueOnce({ data: { id: 3 } })
    mockRequest.mockResolvedValueOnce({ data: [{ id: 3, event_type: 'order_paid', channel: 'sms', template: 'paid' }] })
    const wrapper = mount(Page)
    await flushPromises()
    await wrapper.find('[data-test="event"]').setValue('order_paid')
    await wrapper.find('[data-test="template"]').setValue('paid')
    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()
    expect(mockRequest).toHaveBeenCalledWith('/notification-templates', expect.objectContaining({ method: 'POST' }))
    expect(mockRequest).toHaveBeenCalledWith('/notification-templates')
  })
})
