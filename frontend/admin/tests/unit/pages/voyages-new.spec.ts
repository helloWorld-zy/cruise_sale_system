import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/voyages/new.vue'

const mockRequest = vi.fn().mockResolvedValue({})
const mockNavigate = vi.fn().mockResolvedValue(undefined)
vi.stubGlobal('useApi', () => ({ request: mockRequest }))
vi.stubGlobal('navigateTo', mockNavigate)

beforeEach(() => {
  mockRequest.mockClear()
  mockNavigate.mockClear()
})

describe('VoyagesNewPage', () => {
  it('提交表单触发创建并跳转', async () => {
    const wrapper = mount(Page)
    const inputs = wrapper.findAll('input')
    await inputs[0]!.setValue('V-100')
    await inputs[1]!.setValue('Voyage 100')
    await inputs[2]!.setValue('12')
    await inputs[3]!.setValue('2026-05-01')
    await inputs[4]!.setValue('2026-05-10')
    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(mockRequest).toHaveBeenCalledWith('/voyages', expect.objectContaining({ method: 'POST' }))
    expect(mockNavigate).toHaveBeenCalledWith('/voyages')
  })
})
