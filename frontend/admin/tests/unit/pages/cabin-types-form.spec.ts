import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import NewPage from '../../../app/pages/cabin-types/new.vue'
import EditPage from '../../../app/pages/cabin-types/[id].vue'

const mockRequest = vi.fn()
const mockNavigateTo = vi.fn().mockResolvedValue(undefined)
const confirmMock = vi.fn(() => true)
const routeMock = { params: { id: '11' } }

vi.stubGlobal('useApi', () => ({ request: mockRequest }))
vi.stubGlobal('navigateTo', mockNavigateTo)
vi.stubGlobal('useRoute', () => routeMock)
vi.stubGlobal('confirm', confirmMock)

describe('Cabin type form pages', () => {
  beforeEach(() => {
    mockRequest.mockReset()
    mockNavigateTo.mockClear()
    confirmMock.mockClear()
    confirmMock.mockReturnValue(true)

    mockRequest.mockImplementation((url: string, options?: any) => {
      if (url === '/cruises') {
        return Promise.resolve({ data: [{ id: 1, name: 'Ocean Nova' }] })
      }
      if (url === '/cabin-types' && !options) {
        return Promise.resolve({
          data: {
            list: [
              { id: 11, cruise_id: 1, name: '阳台房', code: 'BAL', tags: '亲子优选', amenities: '独立卫浴', status: 1 },
            ],
          },
        })
      }
      if (url === '/cabin-types' && options?.method === 'POST') {
        return Promise.resolve({ data: { id: 12 } })
      }
      if (url === '/cabin-types/11' && options?.method === 'PUT') {
        return Promise.resolve({ data: { ok: true } })
      }
      if (url === '/cabin-types/11' && options?.method === 'DELETE') {
        return Promise.resolve({ data: { ok: true } })
      }
      return Promise.resolve({ data: {} })
    })
  })

  it('new page submits with merged tags and amenities', async () => {
    const wrapper = mount(NewPage, {
      global: { stubs: { NuxtLink: { template: '<a><slot /></a>' } } },
    })
    await flushPromises()

    const textInputs = wrapper.findAll('input')
    await textInputs[0]!.setValue('豪华阳台房')

    const checkboxes = wrapper.findAll('input[type="checkbox"]')
    if (checkboxes[0]) await checkboxes[0].setValue(true)
    if (checkboxes[5]) await checkboxes[5].setValue(true)

    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(mockRequest).toHaveBeenCalledWith('/cabin-types', expect.objectContaining({ method: 'POST' }))
    expect(mockNavigateTo).toHaveBeenCalledWith('/cabin-types')
  })

  it('edit page loads detail and updates', async () => {
    const wrapper = mount(EditPage, {
      global: { stubs: { NuxtLink: { template: '<a><slot /></a>' } } },
    })
    await flushPromises()

    expect(wrapper.text()).toContain('编辑舱房类型')
    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(mockRequest).toHaveBeenCalledWith('/cabin-types/11', expect.objectContaining({ method: 'PUT' }))
    expect(mockNavigateTo).toHaveBeenCalledWith('/cabin-types')
  })

  it('edit page delete respects confirm', async () => {
    confirmMock.mockReturnValueOnce(false)
    const wrapper = mount(EditPage, {
      global: { stubs: { NuxtLink: { template: '<a><slot /></a>' } } },
    })
    await flushPromises()

    const deleteBtn = wrapper.findAll('button').find((b) => b.text().includes('删除'))
    await deleteBtn!.trigger('click')
    await flushPromises()

    expect(mockRequest).not.toHaveBeenCalledWith('/cabin-types/11', { method: 'DELETE' })
  })
})
