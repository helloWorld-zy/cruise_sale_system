import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/bookings/new.vue'

const mockFetch = vi.fn()
vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: '/api/v1' } }))
vi.stubGlobal('$fetch', mockFetch)

describe('Bookings new page', () => {
  beforeEach(() => {
    mockFetch.mockReset()
  })

  it('submits create booking request with numeric payload', async () => {
    mockFetch.mockResolvedValue({ data: { id: 100 } })
    const wrapper = mount(Page)

    const inputs = wrapper.findAll('input')
    await inputs[0]!.setValue('3')
    await inputs[1]!.setValue('11')
    await inputs[2]!.setValue('22')
    await inputs[3]!.setValue('2')

    await wrapper.find('button').trigger('click')
    await flushPromises()

    expect(mockFetch).toHaveBeenCalledWith('/api/v1/admin/bookings', {
      method: 'POST',
      body: {
        user_id: 3,
        voyage_id: 11,
        cabin_sku_id: 22,
        guests: 2,
      },
      headers: expect.objectContaining({ 'Content-Type': 'application/json' }),
    })
    expect(wrapper.text()).toContain('创建成功')
  })

  it('shows request error message', async () => {
    mockFetch.mockRejectedValueOnce(new Error('create failed'))
    const wrapper = mount(Page)

    await wrapper.find('button').trigger('click')
    await flushPromises()

    expect(wrapper.text()).toContain('create failed')
  })
})
