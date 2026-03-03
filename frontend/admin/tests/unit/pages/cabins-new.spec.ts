import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/cabins/new.vue'

const mockFetch = vi.fn()
const mockNavigateTo = vi.fn().mockResolvedValue(undefined)

vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: '/api/v1' } }))
vi.stubGlobal('$fetch', mockFetch)
vi.stubGlobal('navigateTo', mockNavigateTo)

describe('Cabins new page', () => {
  beforeEach(() => {
    mockFetch.mockReset()
    mockNavigateTo.mockClear()
  })

  it('creates cabin and navigates', async () => {
    mockFetch.mockResolvedValue({ data: { id: 21 } })
    const wrapper = mount(Page)

    const inputs = wrapper.findAll('input')
    await inputs[0]!.setValue('8')
    await inputs[1]!.setValue('3')
    await inputs[2]!.setValue('A201')
    await inputs[3]!.setValue('3')
    await inputs[4]!.setValue('10F')

    await wrapper.find('button').trigger('click')
    await flushPromises()

    expect(mockFetch).toHaveBeenCalledWith('/api/v1/admin/cabins', {
      method: 'POST',
      body: {
        voyage_id: 8,
        cabin_type_id: 3,
        code: 'A201',
        max_guests: 3,
        deck: '10F',
      },
      headers: expect.objectContaining({ 'Content-Type': 'application/json' }),
    })
    expect(mockNavigateTo).toHaveBeenCalledWith('/cabins')
  })

  it('renders error when create fails', async () => {
    mockFetch.mockRejectedValueOnce(new Error('create cabin failed'))
    const wrapper = mount(Page)

    await wrapper.find('button').trigger('click')
    await flushPromises()

    expect(wrapper.text()).toContain('create cabin failed')
  })
})
