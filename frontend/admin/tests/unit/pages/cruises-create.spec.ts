import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/cruises/create.vue'

const mockFetch = vi.fn()
const mockNavigateTo = vi.fn().mockResolvedValue(undefined)

vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: '/api/v1' } }))
vi.stubGlobal('$fetch', mockFetch)
vi.stubGlobal('navigateTo', mockNavigateTo)

describe('Cruises create page', () => {
  beforeEach(() => {
    mockFetch.mockReset()
    mockNavigateTo.mockClear()
  })

  it('creates cruise and navigates back to list', async () => {
    mockFetch.mockResolvedValue({ data: { id: 5 } })
    const wrapper = mount(Page)

    const inputs = wrapper.findAll('input')
    await inputs[0]!.setValue('Ocean Star')
    await inputs[1]!.setValue('Ocean Star EN')
    await inputs[2]!.setValue('OS-01')
    await inputs[3]!.setValue('2')
    await inputs[4]!.setValue('90000')

    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(mockFetch).toHaveBeenCalledWith(
      '/api/v1/admin/cruises',
      expect.objectContaining({
        method: 'POST',
        headers: expect.objectContaining({ 'Content-Type': 'application/json' }),
      }),
    )
    expect(mockNavigateTo).toHaveBeenCalledWith('/cruises')
  })

  it('shows create error message when request fails', async () => {
    mockFetch.mockRejectedValueOnce(new Error('create cruise failed'))
    const wrapper = mount(Page)

    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(wrapper.text()).toContain('create cruise failed')
  })
})
