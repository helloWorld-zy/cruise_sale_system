import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import EditPage from '../../../app/pages/facility-categories/[id].vue'

const mockRequest = vi.fn()
const mockNavigateTo = vi.fn().mockResolvedValue(undefined)
const confirmMock = vi.fn(() => true)
const routeMock = { params: { id: '3' } }

vi.stubGlobal('useApi', () => ({ request: mockRequest }))
vi.stubGlobal('navigateTo', mockNavigateTo)
vi.stubGlobal('confirm', confirmMock)
vi.stubGlobal('useRoute', () => routeMock)

describe('Facility category edit page', () => {
  beforeEach(() => {
    mockRequest.mockReset()
    mockNavigateTo.mockClear()
    confirmMock.mockClear()
    confirmMock.mockReturnValue(true)

    mockRequest.mockImplementation((url: string, options?: any) => {
      if (url === '/facility-categories' && !options) {
        return Promise.resolve({
          data: [
            { id: 3, name: '餐饮', icon: 'utensils', sort_order: 1, status: 1 },
          ],
        })
      }
      if (url === '/facility-categories/3' && options?.method === 'PUT') {
        return Promise.resolve({ data: { ok: true } })
      }
      if (url === '/facility-categories/3' && options?.method === 'DELETE') {
        return Promise.resolve({ data: { ok: true } })
      }
      return Promise.resolve({ data: {} })
    })
  })

  it('loads detail and updates category', async () => {
    const wrapper = mount(EditPage, {
      global: { stubs: { NuxtLink: { template: '<a><slot /></a>' } } },
    })
    await flushPromises()

    expect(wrapper.text()).toContain('编辑设施分类')
    expect(wrapper.find('[data-test="facility-category-edit-icon-option-included-dining"]').exists()).toBe(true)
    expect(wrapper.find('[data-test="facility-category-edit-icon-option-specialty-dining"]').exists()).toBe(true)
    expect(wrapper.find('[data-test="facility-category-edit-icon-option-bar-lounge"]').exists()).toBe(true)
    expect(wrapper.find('[data-test="facility-category-edit-icon-option-leisure-entertainment"]').exists()).toBe(true)
    expect(wrapper.find('[data-test="facility-category-edit-icon-option-kids-family"]').exists()).toBe(true)
    expect(wrapper.find('[data-test="facility-category-edit-icon-option-suite-privilege"]').exists()).toBe(true)
    expect(wrapper.find('[data-test="facility-category-edit-icon-option-sports-fitness"]').exists()).toBe(true)
    expect(wrapper.find('[data-test="facility-category-edit-icon-option-other"]').exists()).toBe(true)
    await wrapper.find('[data-test="facility-category-edit-icon-option-bar-lounge"]').trigger('click')
    await flushPromises()
    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(mockRequest).toHaveBeenCalledWith('/facility-categories/3', expect.objectContaining({ method: 'PUT' }))
    expect(mockRequest).toHaveBeenCalledWith('/facility-categories/3', expect.objectContaining({
      body: expect.objectContaining({ icon: 'bar-lounge' }),
    }))
    expect(mockNavigateTo).toHaveBeenCalledWith('/facility-categories')
  })

  it('delete canceled by confirm false', async () => {
    confirmMock.mockReturnValueOnce(false)
    const wrapper = mount(EditPage, {
      global: { stubs: { NuxtLink: { template: '<a><slot /></a>' } } },
    })
    await flushPromises()

    const deleteBtn = wrapper.findAll('button').find((b) => b.text().includes('删除'))
    await deleteBtn!.trigger('click')
    await flushPromises()

    expect(mockRequest).not.toHaveBeenCalledWith('/facility-categories/3', { method: 'DELETE' })
  })
})
