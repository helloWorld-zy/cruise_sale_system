import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import Page from '../../../app/pages/staff/index.vue'

const mockRequest = vi.fn()
vi.stubGlobal('useApi', () => ({ request: mockRequest }))

const mountOptions = {
  global: {
    stubs: {
      AdminPageHeader: { props: ['title', 'subtitle'], template: '<div>{{ title }} {{ subtitle }}<slot /><slot name="actions" /></div>' },
      AdminDataCard: { props: ['flush'], template: '<div><slot /></div>' },
      AdminStatusTag: { props: ['text'], template: '<span>{{ text }}</span>' },
    },
  },
}

describe('Staff page', () => {
  beforeEach(() => {
    mockRequest.mockReset()
  })

  it('shows loading and then table data', async () => {
    let resolveRequest: any
    mockRequest.mockImplementationOnce(
      () =>
        new Promise((resolve) => {
          resolveRequest = resolve
        }),
    )
    const wrapper = mount(Page, mountOptions)
    await flushPromises()
    expect(wrapper.find('[data-test="loading"]').exists()).toBe(true)
    resolveRequest({ data: [{ id: 1, real_name: 'Alice', email: 'a@test.com', role: 'operator', status: 1 }] })
    await flushPromises()
    expect(mockRequest).toHaveBeenCalledWith('/staffs')
    expect(wrapper.find('[data-test="table"]').exists()).toBe(true)
    expect(wrapper.text()).toContain('Alice')
  })

  it('shows empty state', async () => {
    mockRequest.mockResolvedValueOnce({ data: [] })
    const wrapper = mount(Page, mountOptions)
    await flushPromises()
    expect(wrapper.find('[data-test="empty"]').exists()).toBe(true)
  })

  it('shows error state', async () => {
    mockRequest.mockRejectedValueOnce(new Error('staff api failed'))
    const wrapper = mount(Page, mountOptions)
    await flushPromises()
    expect(wrapper.find('[data-test="error"]').text()).toContain('staff api failed')
  })
})
