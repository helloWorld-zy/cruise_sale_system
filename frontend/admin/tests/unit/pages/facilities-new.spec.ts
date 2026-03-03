import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/facilities/new.vue'

const mockRequest = vi.fn()
const mockNavigateTo = vi.fn().mockResolvedValue(undefined)

vi.stubGlobal('useApi', () => ({ request: mockRequest }))
vi.stubGlobal('navigateTo', mockNavigateTo)

describe('Facilities new page', () => {
  beforeEach(() => {
    mockRequest.mockReset()
    mockNavigateTo.mockClear()

    mockRequest.mockImplementation((url: string, options?: any) => {
      if (url === '/cruises') return Promise.resolve({ data: [{ id: 1, name: 'Ocean Nova' }] })
      if (url === '/facility-categories') return Promise.resolve({ data: [{ id: 2, name: '娱乐' }] })
      if (url === '/facilities' && options?.method === 'POST') return Promise.resolve({ data: { id: 9 } })
      return Promise.resolve({ data: [] })
    })
  })

  it('loads options and submits form', async () => {
    const wrapper = mount(Page, {
      global: { stubs: { NuxtLink: { template: '<a><slot /></a>' } } },
    })
    await flushPromises()

    const textInputs = wrapper.findAll('input[type="text"], input:not([type])')
    await textInputs[0]!.setValue('海上剧院')

    const audienceCheckboxes = wrapper.findAll('input[type="checkbox"]')
    if (audienceCheckboxes[1]) await audienceCheckboxes[1].setValue(true)

    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(mockRequest).toHaveBeenCalledWith('/facilities', expect.objectContaining({ method: 'POST' }))
    expect(mockNavigateTo).toHaveBeenCalledWith('/facilities')
  })

  it('falls back to empty options when load fails', async () => {
    mockRequest.mockImplementationOnce(() => Promise.reject(new Error('load failed')))
    const wrapper = mount(Page, {
      global: { stubs: { NuxtLink: { template: '<a><slot /></a>' } } },
    })
    await flushPromises()

    const selects = wrapper.findAll('select')
    expect(selects[0]?.findAll('option').length).toBe(1)
    expect(selects[1]?.findAll('option').length).toBe(1)
  })
})
