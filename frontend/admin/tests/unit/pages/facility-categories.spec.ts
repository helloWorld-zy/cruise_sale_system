import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/facility-categories/index.vue'

const mockRequest = vi.fn()
const confirmMock = vi.fn(() => true)
vi.stubGlobal('useApi', () => ({ request: mockRequest }))
vi.stubGlobal('confirm', confirmMock)

beforeEach(() => {
  mockRequest.mockReset()
  confirmMock.mockClear()
  confirmMock.mockReturnValue(true)
  mockRequest.mockImplementation((url: string, options?: any) => {
    if (url === '/facility-categories' && !options) {
      return Promise.resolve({
        data: [{ id: 1, name: '餐饮', icon: 'utensils', sort_order: 1, status: 1 }],
      })
    }
    if (url === '/facility-categories' && options?.method === 'POST') return Promise.resolve({ data: { id: 2 } })
    if (url === '/facility-categories/1' && options?.method === 'PUT') return Promise.resolve({ data: { ok: true } })
    if (url === '/facility-categories/1' && options?.method === 'DELETE') return Promise.resolve({ data: { ok: true } })
    return Promise.resolve({ data: [] })
  })
})

describe('FacilityCategory list', () => {
  it('renders compact table with inline inputs', async () => {
    const wrapper = mount(Page, {
      global: { stubs: { NuxtLink: { template: '<a><slot /></a>' } } },
    })
    await flushPromises()
    expect(wrapper.find('table').exists()).toBe(true)
    expect(wrapper.find('tbody input').exists()).toBe(true)
  })

  it('loads facility categories from api', async () => {
    const wrapper = mount(Page, {
      global: { stubs: { NuxtLink: { template: '<a><slot /></a>' } } },
    })
    await flushPromises()
    expect(mockRequest).toHaveBeenCalledWith('/facility-categories')
    const nameInput = wrapper.find('tbody tr td input')
    expect((nameInput.element as HTMLInputElement).value).toBe('餐饮')
  })

  it('appends a new editable row', async () => {
    const wrapper = mount(Page, {
      global: { stubs: { NuxtLink: { template: '<a><slot /></a>' } } },
    })
    await flushPromises()

    const addBtn = wrapper.findAll('button').find((b) => b.text().includes('新增分类'))
    await addBtn!.trigger('click')
    await flushPromises()

    const rows = wrapper.findAll('tbody tr')
    expect(rows.length).toBeGreaterThan(1)
  })

  it('saves existing row with PUT', async () => {
    const wrapper = mount(Page, {
      global: { stubs: { NuxtLink: { template: '<a><slot /></a>' } } },
    })
    await flushPromises()

    const saveBtn = wrapper.findAll('button').find((b) => b.text().includes('保存'))
    await saveBtn!.trigger('click')
    await flushPromises()

    expect(mockRequest).toHaveBeenCalledWith('/facility-categories/1', expect.objectContaining({ method: 'PUT' }))
  })

  it('opens icon picker and selects icon', async () => {
    const wrapper = mount(Page, {
      global: { stubs: { NuxtLink: { template: '<a><slot /></a>' } } },
    })
    await flushPromises()

    const selectBtn = wrapper.findAll('button').find((b) => b.text().includes('选择'))
    await selectBtn!.trigger('click')
    await flushPromises()

    const iconBtn = wrapper.findAll('button').find((b) => b.text() === 'music')
    await iconBtn!.trigger('click')
    await flushPromises()

    const iconInput = wrapper.findAll('tbody tr td input')[1]
    expect((iconInput.element as HTMLInputElement).value).toBe('music')
  })

  it('does not delete when confirm is cancelled', async () => {
    confirmMock.mockReturnValueOnce(false)
    const wrapper = mount(Page, {
      global: { stubs: { NuxtLink: { template: '<a><slot /></a>' } } },
    })
    await flushPromises()

    const deleteBtn = wrapper.findAll('button').find((b) => b.text().includes('删除'))
    await deleteBtn!.trigger('click')
    await flushPromises()

    expect(mockRequest).not.toHaveBeenCalledWith('/facility-categories/1', { method: 'DELETE' })
  })
})
